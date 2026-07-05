package services

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) CreateWorkerProfile(ctx context.Context, cmd ports.WorkerProfileCommand) (*domain.WorkerProfile, error) {
	prepared, err := s.prepareWorkerProfileCommand(ctx, cmd)
	if err != nil {
		return nil, err
	}
	item, err := domain.NewWorkerProfile(prepared)
	if err != nil {
		s.logError("validate worker profile create", err, serviceTenantIDField(cmd.TenantID), serviceStringField("worker_type_id", cmd.WorkerTypeID.String()))
		return nil, err
	}
	result, err := s.workerProfiles.CreateWorkerProfile(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create worker profile", err, serviceTenantIDField(cmd.TenantID), serviceStringField("worker_type_id", item.WorkerTypeID.String()))
		return nil, err
	}
	s.log.Info().Str("tenant_id", result.TenantID.String()).Str("worker_profile_id", result.ID.String()).Msg("hrms: worker profile created")
	return result, nil
}

func (s *TenantService) UpdateWorkerProfile(ctx context.Context, cmd ports.WorkerProfileCommand) (*domain.WorkerProfile, error) {
	if cmd.ID == uuid.Nil {
		err := domain.ErrInvalidWorkerProfileID
		s.logError("validate worker profile update id", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	if _, err := s.GetWorkerProfile(ctx, cmd.TenantID, cmd.ID); err != nil {
		return nil, err
	}
	prepared, err := s.prepareWorkerProfileCommand(ctx, cmd)
	if err != nil {
		return nil, err
	}
	item, err := domain.NewWorkerProfile(prepared)
	if err != nil {
		s.logError("validate worker profile update", err, serviceTenantIDField(cmd.TenantID), serviceStringField("worker_profile_id", cmd.ID.String()))
		return nil, err
	}
	item.ID = cmd.ID
	result, err := s.workerProfiles.UpdateWorkerProfile(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("update worker profile", err, serviceTenantIDField(cmd.TenantID), serviceStringField("worker_profile_id", cmd.ID.String()))
		return nil, err
	}
	s.log.Info().Str("tenant_id", result.TenantID.String()).Str("worker_profile_id", result.ID.String()).Msg("hrms: worker profile updated")
	return result, nil
}

func (s *TenantService) GetWorkerProfile(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.WorkerProfile, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate worker profile get tenant", err)
		return nil, err
	}
	if id == uuid.Nil {
		err := domain.ErrInvalidWorkerProfileID
		s.logError("validate worker profile get id", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	result, err := s.workerProfiles.GetWorkerProfile(ctx, tenantID, id)
	if err != nil {
		s.logError("get worker profile", err, serviceTenantIDField(tenantID), serviceStringField("worker_profile_id", id.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) GetWorkerProfileByEmployeeID(ctx context.Context, tenantID uuid.UUID, employeeID uuid.UUID) (*domain.WorkerProfile, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate worker profile by employee tenant", err)
		return nil, err
	}
	if employeeID == uuid.Nil {
		err := domain.ErrInvalidEmployeeID
		s.logError("validate worker profile by employee id", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	result, err := s.workerProfiles.GetWorkerProfileByEmployeeID(ctx, tenantID, employeeID)
	if err != nil {
		s.logError("get worker profile by employee", err, serviceTenantIDField(tenantID), serviceStringField("employee_id", employeeID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) ListWorkerProfiles(ctx context.Context, filter domain.WorkerProfileFilter) ([]*domain.WorkerProfileListItem, error) {
	if filter.TenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate worker profile list tenant", err)
		return nil, err
	}
	result, err := s.workerProfiles.ListWorkerProfiles(ctx, filter)
	if err != nil {
		s.logError("list worker profiles", err, serviceTenantIDField(filter.TenantID))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) DeleteWorkerProfile(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate worker profile delete tenant", err)
		return err
	}
	if id == uuid.Nil {
		err := domain.ErrInvalidWorkerProfileID
		s.logError("validate worker profile delete id", err, serviceTenantIDField(tenantID))
		return err
	}
	if err := s.workerProfiles.DeleteWorkerProfile(ctx, tenantID, id, actorID); err != nil {
		s.logError("delete worker profile", err, serviceTenantIDField(tenantID), serviceStringField("worker_profile_id", id.String()))
		return err
	}
	return nil
}

func (s *TenantService) prepareWorkerProfileCommand(ctx context.Context, cmd ports.WorkerProfileCommand) (domain.WorkerProfileInput, error) {
	if _, err := s.GetWorkerType(ctx, cmd.TenantID, cmd.WorkerTypeID); err != nil {
		return domain.WorkerProfileInput{}, err
	}
	startDate, err := parseWorkerProfileDate(cmd.StartDate)
	if err != nil {
		s.logError("parse worker profile start date", err, serviceTenantIDField(cmd.TenantID))
		return domain.WorkerProfileInput{}, err
	}
	endDate, err := parseWorkerProfileDate(cmd.EndDate)
	if err != nil {
		s.logError("parse worker profile end date", err, serviceTenantIDField(cmd.TenantID))
		return domain.WorkerProfileInput{}, err
	}
	if cmd.EmployeeID != nil && *cmd.EmployeeID != uuid.Nil {
		profile, err := s.employees.GetEmployeeProfile(ctx, cmd.TenantID, *cmd.EmployeeID)
		if err != nil {
			s.logError("validate worker profile employee link", err, serviceTenantIDField(cmd.TenantID), serviceStringField("employee_id", cmd.EmployeeID.String()))
			return domain.WorkerProfileInput{}, err
		}
		employee := profile.Employee
		cmd.EmployeeUserID = &employee.UserID
		if strings.TrimSpace(cmd.DisplayName) == "" {
			cmd.DisplayName = employeeDisplayName(employee.Firstname, employee.MiddleName, employee.Lastname)
		}
		if cmd.WorkerCode == nil {
			cmd.WorkerCode = employee.EmployeeCode
		}
		if cmd.Email == nil {
			cmd.Email = employee.Email
		}
		if cmd.Mobile == nil {
			cmd.Mobile = employee.Mobile
		}
		if cmd.BranchID == nil {
			cmd.BranchID = employee.BranchID
		}
		if cmd.DepartmentID == nil {
			cmd.DepartmentID = employee.DepartmentID
		}
		if cmd.ReportingManagerID == nil {
			cmd.ReportingManagerID = employee.ReportingManagerID
		}
	}
	return domain.WorkerProfileInput{
		TenantID:           cmd.TenantID,
		WorkerTypeID:       cmd.WorkerTypeID,
		EmployeeID:         cmd.EmployeeID,
		EmployeeUserID:     cmd.EmployeeUserID,
		WorkerCode:         cmd.WorkerCode,
		DisplayName:        cmd.DisplayName,
		LegalName:          cmd.LegalName,
		Email:              cmd.Email,
		Mobile:             cmd.Mobile,
		ProfileStatus:      cmd.ProfileStatus,
		StartDate:          startDate,
		EndDate:            endDate,
		BranchID:           cmd.BranchID,
		DepartmentID:       cmd.DepartmentID,
		ReportingManagerID: cmd.ReportingManagerID,
		WorkLocationLabel:  cmd.WorkLocationLabel,
		SourcePartner:      cmd.SourcePartner,
		ExternalReference:  cmd.ExternalReference,
		ComplianceStatus:   cmd.ComplianceStatus,
		PayrollStatus:      cmd.PayrollStatus,
		Notes:              cmd.Notes,
		Metadata:           cmd.Metadata,
	}, nil
}

func parseWorkerProfileDate(value string) (*time.Time, error) {
	clean := strings.TrimSpace(value)
	if clean == "" {
		return nil, nil
	}
	parsed, err := time.Parse("2006-01-02", clean)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}
