package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) CreateDepartment(ctx context.Context, cmd ports.DepartmentCommand) (*domain.Department, error) {
	department, err := domain.NewDepartment(domain.DepartmentInput{
		TenantID:    cmd.TenantID,
		Name:        cmd.Name,
		ShortCode:   cmd.ShortCode,
		Description: cmd.Description,
	})
	if err != nil {
		s.logError("validate department create", err, serviceTenantIDField(cmd.TenantID), serviceStringField("department_name", cmd.Name), serviceStringField("department_short_code", cmd.ShortCode))
		return nil, err
	}
	result, err := s.departments.CreateDepartment(ctx, department, cmd.ActorID)
	if err != nil {
		s.logError("create department", err, serviceTenantIDField(cmd.TenantID), serviceStringField("department_name", department.Name))
		return nil, err
	}
	s.log.Info().Str("tenant_id", result.TenantID.String()).Str("department_id", result.ID.String()).Msg("hrms: department created")
	return result, nil
}

func (s *TenantService) ListDepartments(ctx context.Context, tenantID uuid.UUID) ([]*domain.Department, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate department list tenant", err)
		return nil, err
	}
	result, err := s.departments.ListDepartments(ctx, tenantID)
	if err != nil {
		s.logError("list departments", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) GetDepartment(ctx context.Context, tenantID uuid.UUID, departmentID uuid.UUID) (*domain.Department, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate department get tenant", err)
		return nil, err
	}
	if departmentID == uuid.Nil {
		err := domain.ErrInvalidDepartmentID
		s.logError("validate department get id", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	result, err := s.departments.GetDepartment(ctx, tenantID, departmentID)
	if err != nil {
		s.logError("get department", err, serviceTenantIDField(tenantID), serviceStringField("department_id", departmentID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) UpdateDepartment(ctx context.Context, cmd ports.DepartmentCommand) (*domain.Department, error) {
	if cmd.ID == uuid.Nil {
		err := domain.ErrInvalidDepartmentID
		s.logError("validate department update id", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	department, err := domain.NewDepartment(domain.DepartmentInput{
		TenantID:    cmd.TenantID,
		Name:        cmd.Name,
		ShortCode:   cmd.ShortCode,
		Description: cmd.Description,
	})
	if err != nil {
		s.logError("validate department update", err, serviceTenantIDField(cmd.TenantID), serviceStringField("department_id", cmd.ID.String()), serviceStringField("department_name", cmd.Name), serviceStringField("department_short_code", cmd.ShortCode))
		return nil, err
	}
	department.ID = cmd.ID
	result, err := s.departments.UpdateDepartment(ctx, department, cmd.ActorID)
	if err != nil {
		s.logError("update department", err, serviceTenantIDField(cmd.TenantID), serviceStringField("department_id", cmd.ID.String()))
		return nil, err
	}
	s.log.Info().Str("tenant_id", result.TenantID.String()).Str("department_id", result.ID.String()).Msg("hrms: department updated")
	return result, nil
}

func (s *TenantService) DeleteDepartment(ctx context.Context, tenantID uuid.UUID, departmentID uuid.UUID, actorID *uuid.UUID) error {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate department delete tenant", err)
		return err
	}
	if departmentID == uuid.Nil {
		err := domain.ErrInvalidDepartmentID
		s.logError("validate department delete id", err, serviceTenantIDField(tenantID))
		return err
	}
	if err := s.departments.DeleteDepartment(ctx, tenantID, departmentID, actorID); err != nil {
		s.logError("delete department", err, serviceTenantIDField(tenantID), serviceStringField("department_id", departmentID.String()))
		return err
	}
	s.log.Info().Str("tenant_id", tenantID.String()).Str("department_id", departmentID.String()).Msg("hrms: department deactivated")
	return nil
}
