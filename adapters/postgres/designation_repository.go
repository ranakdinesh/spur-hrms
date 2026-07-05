package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func (s *Store) CreateDesignation(ctx context.Context, designation *domain.Designation, actorID *uuid.UUID) (*domain.Designation, error) {
	row, err := s.getQueries(ctx).CreateDesignation(ctx, sqlc.CreateDesignationParams{
		TenantID:           designation.TenantID,
		Name:               designation.Name,
		LevelCode:          designation.LevelCode,
		SeniorityRank:      designation.SeniorityRank,
		Description:        textFromPtr(designation.Description),
		AttendanceRequired: designation.AttendanceRequired,
		CreatedBy:          uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "create designation", err, tenantIDField(designation.TenantID), stringField("designation_name", designation.Name))
	}
	return mapDesignation(row), nil
}

func (s *Store) ListDesignations(ctx context.Context, tenantID uuid.UUID) ([]*domain.Designation, error) {
	rows, err := s.getQueries(ctx).ListDesignations(ctx, tenantID)
	if err != nil {
		return nil, s.logDBError(ctx, "list designations", err, tenantIDField(tenantID))
	}
	return mapDesignations(rows), nil
}

func (s *Store) GetDesignation(ctx context.Context, tenantID uuid.UUID, designationID uuid.UUID) (*domain.Designation, error) {
	row, err := s.getQueries(ctx).GetDesignation(ctx, sqlc.GetDesignationParams{
		TenantID: tenantID,
		ID:       designationID,
	})
	if err != nil {
		return nil, s.logDBError(ctx, "get designation", err, tenantIDField(tenantID), stringField("designation_id", designationID.String()))
	}
	return mapDesignation(row), nil
}

func (s *Store) UpdateDesignation(ctx context.Context, designation *domain.Designation, actorID *uuid.UUID) (*domain.Designation, error) {
	row, err := s.getQueries(ctx).UpdateDesignation(ctx, sqlc.UpdateDesignationParams{
		TenantID:           designation.TenantID,
		ID:                 designation.ID,
		Name:               designation.Name,
		LevelCode:          designation.LevelCode,
		SeniorityRank:      designation.SeniorityRank,
		Description:        textFromPtr(designation.Description),
		AttendanceRequired: designation.AttendanceRequired,
		UpdatedBy:          uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "update designation", err, tenantIDField(designation.TenantID), stringField("designation_id", designation.ID.String()))
	}
	return mapDesignation(row), nil
}

func (s *Store) DeleteDesignation(ctx context.Context, tenantID uuid.UUID, designationID uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteDesignation(ctx, sqlc.SoftDeleteDesignationParams{
		TenantID:  tenantID,
		ID:        designationID,
		UpdatedBy: uuidFromPtr(actorID),
	}); err != nil {
		return s.logDBError(ctx, "delete designation", err, tenantIDField(tenantID), stringField("designation_id", designationID.String()))
	}
	return nil
}
