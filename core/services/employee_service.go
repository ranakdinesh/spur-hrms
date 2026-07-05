package services

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

type employeeProbationValues struct {
	Status         string
	StartDate      *time.Time
	EndDate        *time.Time
	DurationDays   int32
	ConfirmedAt    *time.Time
	IsPayrollStaff bool
}

func (s *TenantService) CreateEmployee(ctx context.Context, cmd ports.CreateEmployeeCommand) (*domain.Employee, error) {
	if s.employeeIdentity == nil {
		err := domain.ErrEmployeeIdentityPortMissing
		s.logError("create employee identity port missing", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	role, err := employeeRole(cmd.Role)
	if err != nil {
		s.logError("validate employee role", err, serviceTenantIDField(cmd.TenantID), serviceStringField("role", cmd.Role))
		return nil, err
	}
	if err := s.ensureTenantEmployeeCapacity(ctx, cmd.TenantID); err != nil {
		s.logError("validate tenant employee subscription capacity", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	email := strings.ToLower(strings.TrimSpace(valueFromPtr(cmd.Email)))
	mobile := strings.TrimSpace(valueFromPtr(cmd.Mobile))
	if err := s.employeeIdentity.CheckEmployeeIdentityAvailability(ctx, ports.EmployeeIdentityAvailabilityCommand{TenantID: cmd.TenantID, Email: email, Mobile: mobile}); err != nil {
		s.logError("check employee identity availability", err, serviceTenantIDField(cmd.TenantID), serviceStringField("email", email))
		return nil, err
	}
	if err := s.validateEmployeeReferences(ctx, cmd); err != nil {
		return nil, err
	}
	identity, err := s.employeeIdentity.CreateEmployeeIdentity(ctx, ports.CreateEmployeeIdentityCommand{
		TenantID:  cmd.TenantID,
		FirstName: cmd.FirstName,
		LastName:  valueFromPtr(cmd.LastName),
		Email:     email,
		Mobile:    mobile,
		Password:  cmd.Password,
		Role:      role,
		ActorID:   cmd.ActorID,
	})
	if err != nil {
		s.logError("create employee identity", err, serviceTenantIDField(cmd.TenantID), serviceStringField("email", email))
		return nil, err
	}
	if identity == nil || identity.UserID == uuid.Nil {
		err := domain.ErrInvalidEmployeeUserID
		s.logError("create employee identity missing user id", err, serviceTenantIDField(cmd.TenantID), serviceStringField("email", email))
		return nil, err
	}
	if err := s.employeeIdentity.AssignEmployeeRole(ctx, ports.AssignEmployeeRoleCommand{TenantID: cmd.TenantID, UserID: identity.UserID, Role: role, ActorID: cmd.ActorID}); err != nil {
		s.logError("assign employee identity role", err, serviceTenantIDField(cmd.TenantID), serviceStringField("user_id", identity.UserID.String()), serviceStringField("role", role))
		return nil, err
	}
	dob, err := parseOptionalDate(cmd.DOB)
	if err != nil {
		s.logError("validate employee dob", domain.ErrInvalidEmployeeResignation, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	joiningDate, err := parseOptionalDate(cmd.JoiningDate)
	if err != nil {
		s.logError("validate employee joining date", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	probation, err := normalizeCreateEmployeeProbation(cmd, joiningDate)
	if err != nil {
		s.logError("validate employee probation", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	employeeCode := cmd.EmployeeCode
	if employeeCode == nil || strings.TrimSpace(*employeeCode) == "" {
		generated := domain.DefaultEmployeeCode(identity.UserID)
		employeeCode = &generated
	}
	if err := s.ensureEmployeeCodeAvailable(ctx, cmd.TenantID, *employeeCode); err != nil {
		return nil, err
	}
	employee, err := domain.NewEmployee(domain.EmployeeInput{
		TenantID:              cmd.TenantID,
		UserID:                identity.UserID,
		EmployeeCode:          employeeCode,
		Firstname:             cmd.FirstName,
		MiddleName:            cmd.MiddleName,
		Lastname:              cmd.LastName,
		Email:                 &email,
		Mobile:                &mobile,
		DOB:                   dob,
		Gender:                cmd.Gender,
		MaritalStatus:         cmd.MaritalStatus,
		BloodGroup:            cmd.BloodGroup,
		ProfilePhotoPath:      cmd.ProfilePhotoPath,
		Address:               cmd.Address,
		City:                  cmd.City,
		State:                 cmd.State,
		Country:               cmd.Country,
		Pincode:               cmd.Pincode,
		EmergencyContact:      cmd.EmergencyContact,
		JoiningDate:           joiningDate,
		DepartmentID:          cmd.DepartmentID,
		BranchID:              cmd.BranchID,
		DesignationID:         cmd.DesignationID,
		ReportingManagerID:    cmd.ReportingManagerID,
		EmploymentTypeID:      cmd.EmploymentTypeID,
		Role:                  &role,
		Grade:                 cmd.Grade,
		ExperienceYear:        cmd.ExperienceYear,
		ExperienceMonth:       cmd.ExperienceMonth,
		ProbationStatus:       probation.Status,
		ProbationStartDate:    probation.StartDate,
		ProbationEndDate:      probation.EndDate,
		ProbationDurationDays: probation.DurationDays,
		ProbationConfirmedAt:  probation.ConfirmedAt,
		IsPayrollStaff:        probation.IsPayrollStaff,
	})
	if err != nil {
		s.logError("validate employee create", err, serviceTenantIDField(cmd.TenantID), serviceStringField("email", email))
		return nil, err
	}
	var result *domain.Employee
	if err := s.system.RunAsSystem(ctx, func(txCtx context.Context) error {
		var createErr error
		result, createErr = s.employees.CreateEmployee(txCtx, employee, cmd.ActorID)
		if createErr != nil {
			return createErr
		}
		if cmd.Statutory != nil {
			cmd.Statutory.TenantID = cmd.TenantID
			cmd.Statutory.UserID = result.UserID
			cmd.Statutory.ActorID = cmd.ActorID
			if _, err := s.employees.UpsertEmployeeStatutory(txCtx, *cmd.Statutory); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		s.logError("create employee record", err, serviceTenantIDField(cmd.TenantID), serviceStringField("user_id", identity.UserID.String()))
		return nil, err
	}
	s.log.Info().Str("tenant_id", result.TenantID.String()).Str("employee_id", result.ID.String()).Str("user_id", result.UserID.String()).Msg("hrms: employee created")
	return result, nil
}

func (s *TenantService) ensureTenantEmployeeCapacity(ctx context.Context, tenantID uuid.UUID) error {
	if tenantID == uuid.Nil {
		return domain.ErrInvalidTenantID
	}
	current, err := s.subscriptions.GetCurrentTenantSubscription(ctx, tenantID)
	if err != nil {
		return domain.ErrTenantSubscriptionRequired
	}
	if current.MaxEmployees <= 0 {
		return domain.ErrSubscriptionEmployeeLimitExceeded
	}
	count, err := s.employees.CountActiveEmployees(ctx, tenantID)
	if err != nil {
		return err
	}
	if count >= current.MaxEmployees {
		return domain.ErrSubscriptionEmployeeLimitExceeded
	}
	return nil
}

func (s *TenantService) ListEmployees(ctx context.Context, tenantID uuid.UUID) ([]*domain.EmployeeListItem, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate employee list tenant", err)
		return nil, err
	}
	items, err := s.employees.ListEmployees(ctx, tenantID)
	if err != nil {
		s.logError("list employees", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	return items, nil
}

func (s *TenantService) GetEmployeeProfile(ctx context.Context, tenantID uuid.UUID, employeeID uuid.UUID, actorID *uuid.UUID) (*domain.EmployeeProfile, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate employee profile tenant", err)
		return nil, err
	}
	if employeeID == uuid.Nil {
		err := domain.ErrInvalidEmployeeID
		s.logError("validate employee profile id", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	profile, err := s.employees.GetEmployeeProfile(ctx, tenantID, employeeID)
	if err != nil {
		s.logError("get employee profile", err, serviceTenantIDField(tenantID), serviceStringField("employee_id", employeeID.String()))
		return nil, err
	}
	if branches, err := s.branches.ListBranches(ctx, tenantID); err != nil {
		s.logError("list employee profile branches", err, serviceTenantIDField(tenantID))
		return nil, err
	} else {
		profile.Lookups.Branches = branches
	}
	if departments, err := s.departments.ListDepartments(ctx, tenantID); err != nil {
		s.logError("list employee profile departments", err, serviceTenantIDField(tenantID))
		return nil, err
	} else {
		profile.Lookups.Departments = departments
	}
	if designations, err := s.designations.ListDesignations(ctx, tenantID); err != nil {
		s.logError("list employee profile designations", err, serviceTenantIDField(tenantID))
		return nil, err
	} else {
		profile.Lookups.Designations = designations
	}
	if employmentTypes, err := s.ListEmploymentTypes(ctx, tenantID, actorID); err != nil {
		s.logError("list employee profile employment types", err, serviceTenantIDField(tenantID))
		return nil, err
	} else {
		profile.Lookups.EmploymentTypes = employmentTypes
	}
	profile.Onboarding = calculateEmployeeOnboardingStatus(profile)
	return profile, nil
}

func (s *TenantService) GetEmployeeSelfProfile(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID, actorID *uuid.UUID) (*domain.EmployeeProfile, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate employee self profile tenant", err)
		return nil, err
	}
	if userID == uuid.Nil {
		err := domain.ErrInvalidEmployeeUserID
		s.logError("validate employee self profile user", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	employee, err := s.employees.GetEmployeeByUserID(ctx, tenantID, userID)
	if err != nil {
		s.logError("get employee self record", err, serviceTenantIDField(tenantID), serviceStringField("user_id", userID.String()))
		return nil, err
	}
	return s.GetEmployeeProfile(ctx, tenantID, employee.ID, actorID)
}

func (s *TenantService) UpdateEmployee(ctx context.Context, cmd ports.UpdateEmployeeCommand) (*domain.EmployeeProfile, error) {
	if cmd.TenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate employee update tenant", err)
		return nil, err
	}
	if cmd.ID == uuid.Nil {
		err := domain.ErrInvalidEmployeeID
		s.logError("validate employee update id", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	current, err := s.employees.GetEmployeeProfile(ctx, cmd.TenantID, cmd.ID)
	if err != nil {
		s.logError("get employee before update", err, serviceTenantIDField(cmd.TenantID), serviceStringField("employee_id", cmd.ID.String()))
		return nil, err
	}
	if err := s.validateUpdateEmployeeReferences(ctx, cmd); err != nil {
		return nil, err
	}
	employee, err := s.employeeFromUpdateCommand(cmd, current.Employee.UserID)
	if err != nil {
		s.logError("validate employee update", err, serviceTenantIDField(cmd.TenantID), serviceStringField("employee_id", cmd.ID.String()))
		return nil, err
	}
	if employee.EmployeeCode != nil {
		if err := s.ensureEmployeeCodeAvailableForUpdate(ctx, cmd.TenantID, cmd.ID, *employee.EmployeeCode); err != nil {
			return nil, err
		}
	}
	if cmd.Bank != nil {
		cmd.Bank.TenantID = cmd.TenantID
		cmd.Bank.UserID = current.Employee.UserID
		cmd.Bank.ActorID = cmd.ActorID
	}
	if cmd.Statutory != nil {
		cmd.Statutory.TenantID = cmd.TenantID
		cmd.Statutory.UserID = current.Employee.UserID
		cmd.Statutory.ActorID = cmd.ActorID
	}
	if err := s.system.RunAsSystem(ctx, func(txCtx context.Context) error {
		if _, err := s.employees.UpdateEmployee(txCtx, employee, cmd.ActorID); err != nil {
			return err
		}
		if cmd.Bank != nil {
			if _, err := s.employees.UpsertPrimaryEmployeeBank(txCtx, *cmd.Bank); err != nil {
				return err
			}
		}
		if cmd.Statutory != nil {
			if _, err := s.employees.UpsertEmployeeStatutory(txCtx, *cmd.Statutory); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		s.logError("update employee transaction", err, serviceTenantIDField(cmd.TenantID), serviceStringField("employee_id", cmd.ID.String()))
		return nil, err
	}
	s.log.Info().Str("tenant_id", cmd.TenantID.String()).Str("employee_id", cmd.ID.String()).Msg("hrms: employee updated")
	return s.GetEmployeeProfile(ctx, cmd.TenantID, cmd.ID, cmd.ActorID)
}

func (s *TenantService) DeactivateEmployee(ctx context.Context, tenantID uuid.UUID, employeeID uuid.UUID, actorID *uuid.UUID) error {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate employee deactivate tenant", err)
		return err
	}
	if employeeID == uuid.Nil {
		err := domain.ErrInvalidEmployeeID
		s.logError("validate employee deactivate id", err, serviceTenantIDField(tenantID))
		return err
	}
	if s.employeeIdentity == nil {
		err := domain.ErrEmployeeIdentityPortMissing
		s.logError("deactivate employee identity port missing", err, serviceTenantIDField(tenantID), serviceStringField("employee_id", employeeID.String()))
		return err
	}
	profile, err := s.employees.GetEmployeeProfile(ctx, tenantID, employeeID)
	if err != nil {
		s.logError("get employee before deactivate", err, serviceTenantIDField(tenantID), serviceStringField("employee_id", employeeID.String()))
		return err
	}
	if err := s.system.RunAsSystem(ctx, func(txCtx context.Context) error {
		if err := s.employees.DeactivateEmployee(txCtx, tenantID, employeeID, actorID); err != nil {
			return err
		}
		return s.employeeIdentity.DeactivateEmployeeIdentity(txCtx, ports.DeactivateEmployeeIdentityCommand{
			TenantID: tenantID,
			UserID:   profile.Employee.UserID,
			ActorID:  actorID,
		})
	}); err != nil {
		s.logError("deactivate employee transaction", err, serviceTenantIDField(tenantID), serviceStringField("employee_id", employeeID.String()), serviceStringField("user_id", profile.Employee.UserID.String()))
		return err
	}
	s.log.Info().Str("tenant_id", tenantID.String()).Str("employee_id", employeeID.String()).Str("user_id", profile.Employee.UserID.String()).Msg("hrms: employee deactivated")
	return nil
}

func employeeRole(value string) (string, error) {
	role := strings.TrimSpace(value)
	if role == "" {
		role = domain.RoleEmployee
	}
	return domain.ValidateRole(role)
}

func (s *TenantService) employeeFromUpdateCommand(cmd ports.UpdateEmployeeCommand, userID uuid.UUID) (*domain.Employee, error) {
	role, err := employeeRole(cmd.Role)
	if err != nil {
		return nil, err
	}
	dob, err := parseOptionalDate(cmd.DOB)
	if err != nil {
		return nil, domain.ErrInvalidEmployeeDate
	}
	joiningDate, err := parseOptionalDate(cmd.JoiningDate)
	if err != nil {
		return nil, domain.ErrInvalidEmployeeDate
	}
	resignationDate, err := parseOptionalDate(cmd.ResignationDate)
	if err != nil {
		return nil, domain.ErrInvalidEmployeeDate
	}
	probation, err := normalizeUpdateEmployeeProbation(cmd, joiningDate)
	if err != nil {
		return nil, err
	}
	email := cleanCommandString(cmd.Email)
	if email != nil {
		value := strings.ToLower(*email)
		email = &value
	}
	employee, err := domain.NewEmployee(domain.EmployeeInput{
		TenantID:              cmd.TenantID,
		UserID:                userID,
		EmployeeCode:          cleanCommandString(cmd.EmployeeCode),
		Firstname:             cmd.FirstName,
		MiddleName:            cleanCommandString(cmd.MiddleName),
		Lastname:              cleanCommandString(cmd.LastName),
		Email:                 email,
		Mobile:                cleanCommandString(cmd.Mobile),
		DOB:                   dob,
		Gender:                cleanCommandString(cmd.Gender),
		MaritalStatus:         cleanCommandString(cmd.MaritalStatus),
		BloodGroup:            cleanCommandString(cmd.BloodGroup),
		ProfilePhotoPath:      cleanCommandString(cmd.ProfilePhotoPath),
		Address:               cleanCommandString(cmd.Address),
		City:                  cleanCommandString(cmd.City),
		State:                 cleanCommandString(cmd.State),
		Country:               cleanCommandString(cmd.Country),
		Pincode:               cleanCommandString(cmd.Pincode),
		EmergencyContact:      cleanCommandString(cmd.EmergencyContact),
		JoiningDate:           joiningDate,
		ResignationDate:       resignationDate,
		DepartmentID:          cmd.DepartmentID,
		BranchID:              cmd.BranchID,
		DesignationID:         cmd.DesignationID,
		ReportingManagerID:    cmd.ReportingManagerID,
		EmploymentTypeID:      cmd.EmploymentTypeID,
		Role:                  &role,
		Grade:                 cleanCommandString(cmd.Grade),
		ExperienceYear:        cmd.ExperienceYear,
		ExperienceMonth:       cmd.ExperienceMonth,
		ProbationStatus:       probation.Status,
		ProbationStartDate:    probation.StartDate,
		ProbationEndDate:      probation.EndDate,
		ProbationDurationDays: probation.DurationDays,
		ProbationConfirmedAt:  probation.ConfirmedAt,
		IsPayrollStaff:        probation.IsPayrollStaff,
	})
	if err != nil {
		return nil, err
	}
	employee.ID = cmd.ID
	return employee, nil
}

func (s *TenantService) ensureEmployeeCodeAvailable(ctx context.Context, tenantID uuid.UUID, employeeCode string) error {
	if strings.TrimSpace(employeeCode) == "" {
		return nil
	}
	exists, err := s.employees.EmployeeCodeExists(ctx, tenantID, employeeCode)
	if err != nil {
		s.logError("check employee code uniqueness", err, serviceTenantIDField(tenantID), serviceStringField("employee_code", employeeCode))
		return err
	}
	if exists {
		s.logError("validate employee code uniqueness", domain.ErrEmployeeCodeAlreadyExists, serviceTenantIDField(tenantID), serviceStringField("employee_code", employeeCode))
		return domain.ErrEmployeeCodeAlreadyExists
	}
	return nil
}

func (s *TenantService) ensureEmployeeCodeAvailableForUpdate(ctx context.Context, tenantID uuid.UUID, employeeID uuid.UUID, employeeCode string) error {
	if strings.TrimSpace(employeeCode) == "" {
		return nil
	}
	exists, err := s.employees.EmployeeCodeExistsForOtherEmployee(ctx, tenantID, employeeID, employeeCode)
	if err != nil {
		s.logError("check employee code update uniqueness", err, serviceTenantIDField(tenantID), serviceStringField("employee_code", employeeCode), serviceStringField("employee_id", employeeID.String()))
		return err
	}
	if exists {
		s.logError("validate employee code update uniqueness", domain.ErrEmployeeCodeAlreadyExists, serviceTenantIDField(tenantID), serviceStringField("employee_code", employeeCode), serviceStringField("employee_id", employeeID.String()))
		return domain.ErrEmployeeCodeAlreadyExists
	}
	return nil
}

func calculateEmployeeOnboardingStatus(profile *domain.EmployeeProfile) domain.EmployeeOnboardingStatus {
	status := domain.EmployeeOnboardingStatus{
		Status:                   "not_started",
		MissingRequiredDocuments: []string{},
	}
	if profile == nil {
		return status
	}
	latestByType := map[uuid.UUID]*domain.EmployeeDocument{}
	for _, document := range profile.Documents {
		if document == nil || document.DocumentTypeID == nil {
			continue
		}
		if _, exists := latestByType[*document.DocumentTypeID]; !exists {
			latestByType[*document.DocumentTypeID] = document
		}
		switch document.Status {
		case domain.EmployeeDocumentStatusPendingReview:
			status.PendingReviewDocuments++
		case domain.EmployeeDocumentStatusRejected, domain.EmployeeDocumentStatusResubmissionRequested:
			status.RejectedDocuments++
		}
	}
	for _, documentType := range profile.Lookups.DocumentTypes {
		if documentType == nil || !documentType.IsRequired {
			continue
		}
		status.RequiredDocuments++
		document := latestByType[documentType.ID]
		if document == nil {
			status.MissingRequiredDocuments = append(status.MissingRequiredDocuments, documentType.Name)
			continue
		}
		status.UploadedRequiredDocuments++
		if document.Status == domain.EmployeeDocumentStatusApproved {
			status.ApprovedRequiredDocuments++
		}
	}
	switch {
	case status.RequiredDocuments == 0:
		status.Status = "not_configured"
		status.IsComplete = true
	case len(status.MissingRequiredDocuments) > 0:
		status.Status = "documents_pending"
	case status.RejectedDocuments > 0:
		status.Status = "rework_required"
	case status.PendingReviewDocuments > 0 || status.ApprovedRequiredDocuments < status.RequiredDocuments:
		status.Status = "review_pending"
	default:
		status.Status = "complete"
		status.IsComplete = true
	}
	return status
}

func (s *TenantService) validateEmployeeReferences(ctx context.Context, cmd ports.CreateEmployeeCommand) error {
	if cmd.TenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate employee tenant", err)
		return err
	}
	if cmd.BranchID != nil && *cmd.BranchID != uuid.Nil {
		if _, err := s.branches.GetBranch(ctx, cmd.TenantID, *cmd.BranchID); err != nil {
			s.logError("validate employee branch", err, serviceTenantIDField(cmd.TenantID), serviceStringField("branch_id", cmd.BranchID.String()))
			return err
		}
	}
	if cmd.DepartmentID != nil && *cmd.DepartmentID != uuid.Nil {
		if _, err := s.departments.GetDepartment(ctx, cmd.TenantID, *cmd.DepartmentID); err != nil {
			s.logError("validate employee department", err, serviceTenantIDField(cmd.TenantID), serviceStringField("department_id", cmd.DepartmentID.String()))
			return err
		}
	}
	if cmd.DesignationID != nil && *cmd.DesignationID != uuid.Nil {
		if _, err := s.designations.GetDesignation(ctx, cmd.TenantID, *cmd.DesignationID); err != nil {
			s.logError("validate employee designation", err, serviceTenantIDField(cmd.TenantID), serviceStringField("designation_id", cmd.DesignationID.String()))
			return err
		}
	}
	if cmd.EmploymentTypeID != nil && *cmd.EmploymentTypeID != uuid.Nil {
		if _, err := s.lookups.GetEmploymentType(ctx, cmd.TenantID, *cmd.EmploymentTypeID); err != nil {
			s.logError("validate employee employment type", err, serviceTenantIDField(cmd.TenantID), serviceStringField("employment_type_id", cmd.EmploymentTypeID.String()))
			return err
		}
	}
	if cmd.ReportingManagerID != nil && *cmd.ReportingManagerID != uuid.Nil {
		if _, err := s.employees.GetEmployeeByUserID(ctx, cmd.TenantID, *cmd.ReportingManagerID); err != nil {
			s.logError("validate employee reporting manager", err, serviceTenantIDField(cmd.TenantID), serviceStringField("reporting_manager_id", cmd.ReportingManagerID.String()))
			return err
		}
	}
	return nil
}

func (s *TenantService) validateUpdateEmployeeReferences(ctx context.Context, cmd ports.UpdateEmployeeCommand) error {
	if cmd.BranchID != nil && *cmd.BranchID != uuid.Nil {
		if _, err := s.branches.GetBranch(ctx, cmd.TenantID, *cmd.BranchID); err != nil {
			s.logError("validate employee update branch", err, serviceTenantIDField(cmd.TenantID), serviceStringField("branch_id", cmd.BranchID.String()))
			return err
		}
	}
	if cmd.DepartmentID != nil && *cmd.DepartmentID != uuid.Nil {
		if _, err := s.departments.GetDepartment(ctx, cmd.TenantID, *cmd.DepartmentID); err != nil {
			s.logError("validate employee update department", err, serviceTenantIDField(cmd.TenantID), serviceStringField("department_id", cmd.DepartmentID.String()))
			return err
		}
	}
	if cmd.DesignationID != nil && *cmd.DesignationID != uuid.Nil {
		if _, err := s.designations.GetDesignation(ctx, cmd.TenantID, *cmd.DesignationID); err != nil {
			s.logError("validate employee update designation", err, serviceTenantIDField(cmd.TenantID), serviceStringField("designation_id", cmd.DesignationID.String()))
			return err
		}
	}
	if cmd.EmploymentTypeID != nil && *cmd.EmploymentTypeID != uuid.Nil {
		if _, err := s.lookups.GetEmploymentType(ctx, cmd.TenantID, *cmd.EmploymentTypeID); err != nil {
			s.logError("validate employee update employment type", err, serviceTenantIDField(cmd.TenantID), serviceStringField("employment_type_id", cmd.EmploymentTypeID.String()))
			return err
		}
	}
	if cmd.ReportingManagerID != nil && *cmd.ReportingManagerID != uuid.Nil {
		manager, err := s.employees.GetEmployeeByUserID(ctx, cmd.TenantID, *cmd.ReportingManagerID)
		if err != nil {
			s.logError("validate employee update reporting manager", err, serviceTenantIDField(cmd.TenantID), serviceStringField("reporting_manager_id", cmd.ReportingManagerID.String()))
			return err
		}
		if manager.ID == cmd.ID {
			err := domain.ErrInvalidEmployeeUserID
			s.logError("validate employee update reporting manager self", err, serviceTenantIDField(cmd.TenantID), serviceStringField("employee_id", cmd.ID.String()))
			return err
		}
	}
	return nil
}

func valueFromPtr(value *string) string {
	if value == nil {
		return ""
	}
	return strings.TrimSpace(*value)
}

func cleanCommandString(value *string) *string {
	if value == nil {
		return nil
	}
	clean := strings.TrimSpace(*value)
	if clean == "" {
		return nil
	}
	return &clean
}

func normalizeCreateEmployeeProbation(cmd ports.CreateEmployeeCommand, joiningDate *time.Time) (employeeProbationValues, error) {
	return normalizeEmployeeProbation(cmd.ProbationStatus, cmd.ProbationStartDate, cmd.ProbationEndDate, cmd.ProbationConfirmedAt, cmd.ProbationDurationDays, cmd.IsPayrollStaff, joiningDate)
}

func normalizeUpdateEmployeeProbation(cmd ports.UpdateEmployeeCommand, joiningDate *time.Time) (employeeProbationValues, error) {
	return normalizeEmployeeProbation(cmd.ProbationStatus, cmd.ProbationStartDate, cmd.ProbationEndDate, cmd.ProbationConfirmedAt, cmd.ProbationDurationDays, cmd.IsPayrollStaff, joiningDate)
}

func normalizeEmployeeProbation(statusValue string, startValue string, endValue string, confirmedValue string, durationDays int32, isPayrollStaff bool, joiningDate *time.Time) (employeeProbationValues, error) {
	startDate, err := parseOptionalDate(startValue)
	if err != nil {
		return employeeProbationValues{}, domain.ErrInvalidEmployeeDate
	}
	endDate, err := parseOptionalDate(endValue)
	if err != nil {
		return employeeProbationValues{}, domain.ErrInvalidEmployeeDate
	}
	confirmedAt, err := parseOptionalDate(confirmedValue)
	if err != nil {
		return employeeProbationValues{}, domain.ErrInvalidEmployeeDate
	}
	status := strings.TrimSpace(statusValue)
	if status == "" {
		if isPayrollStaff {
			status = domain.EmployeeProbationProbation
		} else {
			status = domain.EmployeeProbationConfirmed
		}
	}
	switch status {
	case domain.EmployeeProbationNotApplicable, domain.EmployeeProbationProbation, domain.EmployeeProbationConfirmed, domain.EmployeeProbationExtended:
	default:
		return employeeProbationValues{}, domain.ErrInvalidEmployeeProbation
	}
	if isPayrollStaff {
		if durationDays == 0 {
			durationDays = domain.PayrollStaffProbationDurationDays
		}
		if durationDays < domain.PayrollStaffProbationDurationDays {
			return employeeProbationValues{}, domain.ErrInvalidEmployeeProbation
		}
		if startDate == nil {
			startDate = joiningDate
		}
		if endDate == nil && startDate != nil {
			end := startDate.AddDate(0, 6, -1)
			endDate = &end
		}
		if status == domain.EmployeeProbationNotApplicable {
			status = domain.EmployeeProbationProbation
		}
	}
	if status == domain.EmployeeProbationConfirmed && confirmedAt == nil && endDate != nil {
		confirmed := *endDate
		confirmedAt = &confirmed
	}
	return employeeProbationValues{Status: status, StartDate: startDate, EndDate: endDate, DurationDays: durationDays, ConfirmedAt: confirmedAt, IsPayrollStaff: isPayrollStaff}, nil
}
