package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func (s *Store) CreateDepartment(ctx context.Context, department *domain.Department, actorID *uuid.UUID) (*domain.Department, error) {
	row, err := s.getQueries(ctx).CreateDepartment(ctx, sqlc.CreateDepartmentParams{
		TenantID:    department.TenantID,
		Name:        department.Name,
		ShortCode:   department.ShortCode,
		Description: textFromPtr(department.Description),
		CreatedBy:   uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "create department", err, tenantIDField(department.TenantID), stringField("department_name", department.Name))
	}
	return mapDepartment(row), nil
}

func (s *Store) ListDepartments(ctx context.Context, tenantID uuid.UUID) ([]*domain.Department, error) {
	rows, err := s.getQueries(ctx).ListDepartments(ctx, tenantID)
	if err != nil {
		return nil, s.logDBError(ctx, "list departments", err, tenantIDField(tenantID))
	}
	return mapDepartments(rows), nil
}

func (s *Store) GetDepartment(ctx context.Context, tenantID uuid.UUID, departmentID uuid.UUID) (*domain.Department, error) {
	row, err := s.getQueries(ctx).GetDepartment(ctx, sqlc.GetDepartmentParams{
		TenantID: tenantID,
		ID:       departmentID,
	})
	if err != nil {
		return nil, s.logDBError(ctx, "get department", err, tenantIDField(tenantID), stringField("department_id", departmentID.String()))
	}
	return mapDepartment(row), nil
}

func (s *Store) UpdateDepartment(ctx context.Context, department *domain.Department, actorID *uuid.UUID) (*domain.Department, error) {
	row, err := s.getQueries(ctx).UpdateDepartment(ctx, sqlc.UpdateDepartmentParams{
		TenantID:    department.TenantID,
		ID:          department.ID,
		Name:        department.Name,
		ShortCode:   department.ShortCode,
		Description: textFromPtr(department.Description),
		UpdatedBy:   uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "update department", err, tenantIDField(department.TenantID), stringField("department_id", department.ID.String()))
	}
	return mapDepartment(row), nil
}

func (s *Store) DeleteDepartment(ctx context.Context, tenantID uuid.UUID, departmentID uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteDepartment(ctx, sqlc.SoftDeleteDepartmentParams{
		TenantID:  tenantID,
		ID:        departmentID,
		UpdatedBy: uuidFromPtr(actorID),
	}); err != nil {
		return s.logDBError(ctx, "delete department", err, tenantIDField(tenantID), stringField("department_id", departmentID.String()))
	}
	return nil
}
