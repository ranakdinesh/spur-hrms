package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *Store) CreateEmployee(ctx context.Context, employee *domain.Employee, actorID *uuid.UUID) (*domain.Employee, error) {
	row, err := s.getQueries(ctx).CreateEmployee(ctx, sqlc.CreateEmployeeParams{
		TenantID:              employee.TenantID,
		UserID:                employee.UserID,
		EmployeeCode:          textFromPtr(employee.EmployeeCode),
		Firstname:             employee.Firstname,
		MiddleName:            textFromPtr(employee.MiddleName),
		Lastname:              textFromPtr(employee.Lastname),
		Email:                 textFromPtr(employee.Email),
		Mobile:                textFromPtr(employee.Mobile),
		Dob:                   dateFromPtr(employee.DOB),
		Gender:                textFromPtr(employee.Gender),
		MaritalStatus:         textFromPtr(employee.MaritalStatus),
		BloodGroup:            textFromPtr(employee.BloodGroup),
		ProfilePhotoPath:      textFromPtr(employee.ProfilePhotoPath),
		Address:               textFromPtr(employee.Address),
		City:                  textFromPtr(employee.City),
		State:                 textFromPtr(employee.State),
		Country:               textFromPtr(employee.Country),
		Pincode:               textFromPtr(employee.Pincode),
		EmergencyContact:      textFromPtr(employee.EmergencyContact),
		JoiningDate:           dateFromPtr(employee.JoiningDate),
		DepartmentID:          uuidFromPtr(employee.DepartmentID),
		BranchID:              uuidFromPtr(employee.BranchID),
		DesignationID:         uuidFromPtr(employee.DesignationID),
		ReportingManagerID:    uuidFromPtr(employee.ReportingManagerID),
		EmploymentTypeID:      uuidFromPtr(employee.EmploymentTypeID),
		Role:                  textFromPtr(employee.Role),
		Grade:                 textFromPtr(employee.Grade),
		ExperienceYear:        employee.ExperienceYear,
		ExperienceMonth:       employee.ExperienceMonth,
		ProbationStatus:       employee.ProbationStatus,
		ProbationStartDate:    dateFromPtr(employee.ProbationStartDate),
		ProbationEndDate:      dateFromPtr(employee.ProbationEndDate),
		ProbationDurationDays: employee.ProbationDurationDays,
		ProbationConfirmedAt:  dateFromPtr(employee.ProbationConfirmedAt),
		IsPayrollStaff:        employee.IsPayrollStaff,
		CreatedBy:             uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "create employee", err, tenantIDField(employee.TenantID), stringField("user_id", employee.UserID.String()))
	}
	return mapEmployee(row), nil
}

func (s *Store) CountActiveEmployees(ctx context.Context, tenantID uuid.UUID) (int32, error) {
	count, err := s.getQueries(ctx).CountActiveEmployees(ctx, tenantID)
	if err != nil {
		return 0, s.logDBError(ctx, "count active employees", err, tenantIDField(tenantID))
	}
	return int32(count), nil
}

func (s *Store) UpdateEmployee(ctx context.Context, employee *domain.Employee, actorID *uuid.UUID) (*domain.Employee, error) {
	row, err := s.getQueries(ctx).UpdateEmployee(ctx, sqlc.UpdateEmployeeParams{
		TenantID:              employee.TenantID,
		ID:                    employee.ID,
		EmployeeCode:          textFromPtr(employee.EmployeeCode),
		Firstname:             employee.Firstname,
		MiddleName:            textFromPtr(employee.MiddleName),
		Lastname:              textFromPtr(employee.Lastname),
		Email:                 textFromPtr(employee.Email),
		Mobile:                textFromPtr(employee.Mobile),
		Dob:                   dateFromPtr(employee.DOB),
		Gender:                textFromPtr(employee.Gender),
		MaritalStatus:         textFromPtr(employee.MaritalStatus),
		BloodGroup:            textFromPtr(employee.BloodGroup),
		ProfilePhotoPath:      textFromPtr(employee.ProfilePhotoPath),
		Address:               textFromPtr(employee.Address),
		City:                  textFromPtr(employee.City),
		State:                 textFromPtr(employee.State),
		Country:               textFromPtr(employee.Country),
		Pincode:               textFromPtr(employee.Pincode),
		EmergencyContact:      textFromPtr(employee.EmergencyContact),
		JoiningDate:           dateFromPtr(employee.JoiningDate),
		ResignationDate:       dateFromPtr(employee.ResignationDate),
		DepartmentID:          uuidFromPtr(employee.DepartmentID),
		BranchID:              uuidFromPtr(employee.BranchID),
		DesignationID:         uuidFromPtr(employee.DesignationID),
		ReportingManagerID:    uuidFromPtr(employee.ReportingManagerID),
		EmploymentTypeID:      uuidFromPtr(employee.EmploymentTypeID),
		Role:                  textFromPtr(employee.Role),
		Grade:                 textFromPtr(employee.Grade),
		ExperienceYear:        employee.ExperienceYear,
		ExperienceMonth:       employee.ExperienceMonth,
		ProbationStatus:       employee.ProbationStatus,
		ProbationStartDate:    dateFromPtr(employee.ProbationStartDate),
		ProbationEndDate:      dateFromPtr(employee.ProbationEndDate),
		ProbationDurationDays: employee.ProbationDurationDays,
		ProbationConfirmedAt:  dateFromPtr(employee.ProbationConfirmedAt),
		IsPayrollStaff:        employee.IsPayrollStaff,
		UpdatedBy:             uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "update employee", err, tenantIDField(employee.TenantID), stringField("employee_id", employee.ID.String()))
	}
	return mapEmployee(row), nil
}

func (s *Store) DeactivateEmployee(ctx context.Context, tenantID uuid.UUID, employeeID uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteEmployee(ctx, sqlc.SoftDeleteEmployeeParams{
		TenantID:  tenantID,
		ID:        employeeID,
		UpdatedBy: uuidFromPtr(actorID),
	}); err != nil {
		return s.logDBError(ctx, "soft delete employee", err, tenantIDField(tenantID), stringField("employee_id", employeeID.String()))
	}
	return nil
}

func (s *Store) ListEmployees(ctx context.Context, tenantID uuid.UUID) ([]*domain.EmployeeListItem, error) {
	rows, err := s.getQueries(ctx).ListEmployees(ctx, tenantID)
	if err != nil {
		return nil, s.logDBError(ctx, "list employees", err, tenantIDField(tenantID))
	}
	items := make([]*domain.EmployeeListItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapEmployeeListItem(row))
	}
	return items, nil
}

func (s *Store) GetEmployeeAttendanceRequired(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID) (bool, error) {
	required, err := s.getQueries(ctx).GetEmployeeAttendanceRequired(ctx, sqlc.GetEmployeeAttendanceRequiredParams{
		TenantID: tenantID,
		UserID:   userID,
	})
	if err != nil {
		return false, s.logDBError(ctx, "get employee attendance requirement", err, tenantIDField(tenantID), stringField("user_id", userID.String()))
	}
	return required, nil
}

func (s *Store) GetEmployeeProfile(ctx context.Context, tenantID uuid.UUID, employeeID uuid.UUID) (*domain.EmployeeProfile, error) {
	queries := s.getQueries(ctx)
	employeeRow, err := queries.GetEmployeeProfileItem(ctx, sqlc.GetEmployeeProfileItemParams{TenantID: tenantID, ID: employeeID})
	if err != nil {
		return nil, s.logDBError(ctx, "get employee profile item", err, tenantIDField(tenantID), stringField("employee_id", employeeID.String()))
	}
	employee := mapEmployeeProfileItem(employeeRow)
	profile := &domain.EmployeeProfile{
		Employee:  employee,
		Banks:     []*domain.EmployeeBank{},
		Documents: []*domain.EmployeeDocument{},
		Lookups: domain.EmployeeProfileLookups{
			DocumentTypes: []*domain.DocumentType{},
		},
	}

	statutoryRow, err := queries.GetEmployeeStatutoryByUserID(ctx, sqlc.GetEmployeeStatutoryByUserIDParams{TenantID: tenantID, UserID: employee.UserID})
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, s.logDBError(ctx, "get employee statutory", err, tenantIDField(tenantID), stringField("user_id", employee.UserID.String()))
	}
	if err == nil {
		profile.Statutory = mapEmployeeStatutory(statutoryRow)
	}

	bankRows, err := queries.ListEmployeeBanksByUserID(ctx, sqlc.ListEmployeeBanksByUserIDParams{TenantID: tenantID, UserID: employee.UserID})
	if err != nil {
		return nil, s.logDBError(ctx, "list employee banks", err, tenantIDField(tenantID), stringField("user_id", employee.UserID.String()))
	}
	for _, row := range bankRows {
		profile.Banks = append(profile.Banks, mapEmployeeBank(row))
	}

	documentRows, err := queries.ListEmployeeDocumentsByUserID(ctx, sqlc.ListEmployeeDocumentsByUserIDParams{TenantID: tenantID, UserID: employee.UserID})
	if err != nil {
		return nil, s.logDBError(ctx, "list employee documents", err, tenantIDField(tenantID), stringField("user_id", employee.UserID.String()))
	}
	for _, row := range documentRows {
		profile.Documents = append(profile.Documents, mapEmployeeDocument(row))
	}

	documentTypeRows, err := queries.ListDocumentTypes(ctx, tenantID)
	if err != nil {
		return nil, s.logDBError(ctx, "list employee document types", err, tenantIDField(tenantID))
	}
	for _, row := range documentTypeRows {
		profile.Lookups.DocumentTypes = append(profile.Lookups.DocumentTypes, mapDocumentType(row))
	}

	return profile, nil
}

func (s *Store) UpsertEmployeeStatutory(ctx context.Context, item ports.EmployeeStatutoryCommand) (*domain.EmployeeStatutory, error) {
	row, err := s.getQueries(ctx).UpsertEmployeeStatutory(ctx, sqlc.UpsertEmployeeStatutoryParams{
		TenantID:       item.TenantID,
		UserID:         item.UserID,
		PfNo:           textFromPtr(item.PFNo),
		UanNo:          textFromPtr(item.UANNo),
		EsicNo:         textFromPtr(item.ESICNo),
		Pan:            textFromPtr(item.PAN),
		Aadhaar:        textFromPtr(item.Aadhaar),
		PtApplicable:   item.PTApplicable,
		PfApplicable:   item.PFApplicable,
		EsicApplicable: item.ESICApplicable,
		LwfApplicable:  item.LWFApplicable,
		CreatedBy:      uuidFromPtr(item.ActorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "upsert employee statutory", err, tenantIDField(item.TenantID), stringField("user_id", item.UserID.String()))
	}
	return mapEmployeeStatutory(row), nil
}

func (s *Store) UpsertPrimaryEmployeeBank(ctx context.Context, item ports.EmployeeBankCommand) (*domain.EmployeeBank, error) {
	row, err := s.getQueries(ctx).UpsertPrimaryEmployeeBank(ctx, sqlc.UpsertPrimaryEmployeeBankParams{
		TenantID:      item.TenantID,
		UserID:        item.UserID,
		BankName:      textFromPtr(item.BankName),
		AccountNumber: textFromPtr(item.AccountNumber),
		IfscCode:      textFromPtr(item.IFSCCode),
		AccountType:   textFromPtr(item.AccountType),
		BranchName:    textFromPtr(item.BranchName),
		CreatedBy:     uuidFromPtr(item.ActorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "upsert primary employee bank", err, tenantIDField(item.TenantID), stringField("user_id", item.UserID.String()))
	}
	return mapEmployeeBank(row), nil
}

func (s *Store) EmployeeCodeExists(ctx context.Context, tenantID uuid.UUID, employeeCode string) (bool, error) {
	exists, err := s.getQueries(ctx).EmployeeCodeExists(ctx, sqlc.EmployeeCodeExistsParams{TenantID: tenantID, EmployeeCode: textFromString(employeeCode)})
	if err != nil {
		return false, s.logDBError(ctx, "employee code exists", err, tenantIDField(tenantID), stringField("employee_code", employeeCode))
	}
	return exists, nil
}

func (s *Store) EmployeeCodeExistsForOtherEmployee(ctx context.Context, tenantID uuid.UUID, employeeID uuid.UUID, employeeCode string) (bool, error) {
	exists, err := s.getQueries(ctx).EmployeeCodeExistsForOtherEmployee(ctx, sqlc.EmployeeCodeExistsForOtherEmployeeParams{TenantID: tenantID, EmployeeCode: textFromString(employeeCode), ID: employeeID})
	if err != nil {
		return false, s.logDBError(ctx, "employee code exists for other employee", err, tenantIDField(tenantID), stringField("employee_id", employeeID.String()), stringField("employee_code", employeeCode))
	}
	return exists, nil
}

func (s *Store) GetEmployeeByCode(ctx context.Context, tenantID uuid.UUID, employeeCode string) (*domain.Employee, error) {
	row, err := s.getQueries(ctx).GetEmployeeByCode(ctx, sqlc.GetEmployeeByCodeParams{TenantID: tenantID, EmployeeCode: textFromString(employeeCode)})
	if err != nil {
		return nil, s.logDBError(ctx, "get employee by code", err, tenantIDField(tenantID), stringField("employee_code", employeeCode))
	}
	return mapEmployee(row), nil
}

func (s *Store) GetEmployeeByUserID(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID) (*domain.Employee, error) {
	row, err := s.getQueries(ctx).GetEmployeeByUserID(ctx, sqlc.GetEmployeeByUserIDParams{TenantID: tenantID, UserID: userID})
	if err != nil {
		return nil, s.logDBError(ctx, "get employee by user", err, tenantIDField(tenantID), stringField("user_id", userID.String()))
	}
	return mapEmployee(row), nil
}
