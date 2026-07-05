package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func (s *Store) CreateEmploymentType(ctx context.Context, item *domain.EmploymentType, actorID *uuid.UUID) (*domain.EmploymentType, error) {
	row, err := s.getQueries(ctx).CreateEmploymentType(ctx, sqlc.CreateEmploymentTypeParams{
		TenantID:  item.TenantID,
		Name:      item.Name,
		CreatedBy: uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "create employment type", err, tenantIDField(item.TenantID), stringField("employment_type_name", item.Name))
	}
	return mapEmploymentType(row), nil
}

func (s *Store) ListEmploymentTypes(ctx context.Context, tenantID uuid.UUID) ([]*domain.EmploymentType, error) {
	rows, err := s.getQueries(ctx).ListEmploymentTypes(ctx, tenantID)
	if err != nil {
		return nil, s.logDBError(ctx, "list employment types", err, tenantIDField(tenantID))
	}
	return mapEmploymentTypes(rows), nil
}

func (s *Store) GetEmploymentType(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.EmploymentType, error) {
	row, err := s.getQueries(ctx).GetEmploymentType(ctx, sqlc.GetEmploymentTypeParams{
		TenantID: tenantID,
		ID:       id,
	})
	if err != nil {
		return nil, s.logDBError(ctx, "get employment type", err, tenantIDField(tenantID), stringField("employment_type_id", id.String()))
	}
	return mapEmploymentType(row), nil
}

func (s *Store) UpdateEmploymentType(ctx context.Context, item *domain.EmploymentType, actorID *uuid.UUID) (*domain.EmploymentType, error) {
	row, err := s.getQueries(ctx).UpdateEmploymentType(ctx, sqlc.UpdateEmploymentTypeParams{
		TenantID:  item.TenantID,
		ID:        item.ID,
		Name:      item.Name,
		UpdatedBy: uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "update employment type", err, tenantIDField(item.TenantID), stringField("employment_type_id", item.ID.String()))
	}
	return mapEmploymentType(row), nil
}

func (s *Store) DeleteEmploymentType(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteEmploymentType(ctx, sqlc.SoftDeleteEmploymentTypeParams{
		TenantID:  tenantID,
		ID:        id,
		UpdatedBy: uuidFromPtr(actorID),
	}); err != nil {
		return s.logDBError(ctx, "delete employment type", err, tenantIDField(tenantID), stringField("employment_type_id", id.String()))
	}
	return nil
}

func (s *Store) CreateMaritalStatus(ctx context.Context, item *domain.MaritalStatus, actorID *uuid.UUID) (*domain.MaritalStatus, error) {
	row, err := s.getQueries(ctx).CreateMaritalStatus(ctx, sqlc.CreateMaritalStatusParams{
		TenantID:  item.TenantID,
		Name:      item.Name,
		CreatedBy: uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "create marital status", err, tenantIDField(item.TenantID), stringField("marital_status_name", item.Name))
	}
	return mapMaritalStatus(row), nil
}

func (s *Store) ListMaritalStatuses(ctx context.Context, tenantID uuid.UUID) ([]*domain.MaritalStatus, error) {
	rows, err := s.getQueries(ctx).ListMaritalStatuses(ctx, tenantID)
	if err != nil {
		return nil, s.logDBError(ctx, "list marital statuses", err, tenantIDField(tenantID))
	}
	return mapMaritalStatuses(rows), nil
}

func (s *Store) GetMaritalStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.MaritalStatus, error) {
	row, err := s.getQueries(ctx).GetMaritalStatus(ctx, sqlc.GetMaritalStatusParams{
		TenantID: tenantID,
		ID:       id,
	})
	if err != nil {
		return nil, s.logDBError(ctx, "get marital status", err, tenantIDField(tenantID), stringField("marital_status_id", id.String()))
	}
	return mapMaritalStatus(row), nil
}

func (s *Store) UpdateMaritalStatus(ctx context.Context, item *domain.MaritalStatus, actorID *uuid.UUID) (*domain.MaritalStatus, error) {
	row, err := s.getQueries(ctx).UpdateMaritalStatus(ctx, sqlc.UpdateMaritalStatusParams{
		TenantID:  item.TenantID,
		ID:        item.ID,
		Name:      item.Name,
		UpdatedBy: uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "update marital status", err, tenantIDField(item.TenantID), stringField("marital_status_id", item.ID.String()))
	}
	return mapMaritalStatus(row), nil
}

func (s *Store) DeleteMaritalStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteMaritalStatus(ctx, sqlc.SoftDeleteMaritalStatusParams{
		TenantID:  tenantID,
		ID:        id,
		UpdatedBy: uuidFromPtr(actorID),
	}); err != nil {
		return s.logDBError(ctx, "delete marital status", err, tenantIDField(tenantID), stringField("marital_status_id", id.String()))
	}
	return nil
}
