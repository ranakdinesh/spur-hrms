package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func (s *Store) CreateWorkingHour(ctx context.Context, item *domain.WorkingHour, actorID *uuid.UUID) (*domain.WorkingHour, error) {
	row, err := s.getQueries(ctx).CreateWorkingHour(ctx, sqlc.CreateWorkingHourParams{
		TenantID:     item.TenantID,
		BranchID:     uuidFromPtr(item.BranchID),
		UserID:       uuidFromPtr(item.UserID),
		DayOfWeek:    item.DayOfWeek,
		IsWorkingDay: item.IsWorkingDay,
		StartTime:    timeFromClockString(item.StartTime),
		EndTime:      timeFromClockString(item.EndTime),
		BreakMinutes: item.BreakMinutes,
		CreatedBy:    uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "create working hour", err, tenantIDField(item.TenantID), stringField("day_of_week", item.DayOfWeek))
	}
	return mapWorkingHour(row), nil
}

func (s *Store) ListWorkingHours(ctx context.Context, tenantID uuid.UUID) ([]*domain.WorkingHour, error) {
	rows, err := s.getQueries(ctx).ListWorkingHours(ctx, tenantID)
	if err != nil {
		return nil, s.logDBError(ctx, "list working hours", err, tenantIDField(tenantID))
	}
	return mapWorkingHours(rows), nil
}

func (s *Store) GetWorkingHour(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.WorkingHour, error) {
	row, err := s.getQueries(ctx).GetWorkingHour(ctx, sqlc.GetWorkingHourParams{TenantID: tenantID, ID: id})
	if err != nil {
		return nil, s.logDBError(ctx, "get working hour", err, tenantIDField(tenantID), stringField("working_hour_id", id.String()))
	}
	return mapWorkingHour(row), nil
}

func (s *Store) ResolveWorkingHour(ctx context.Context, tenantID uuid.UUID, dayOfWeek string, branchID *uuid.UUID, userID *uuid.UUID) (*domain.WorkingHour, error) {
	row, err := s.getQueries(ctx).ResolveWorkingHour(ctx, sqlc.ResolveWorkingHourParams{
		TenantID:  tenantID,
		DayOfWeek: dayOfWeek,
		UserID:    uuidFromPtr(userID),
		BranchID:  uuidFromPtr(branchID),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrWorkingHourNotFound
		}
		return nil, s.logDBError(ctx, "resolve working hour", err, tenantIDField(tenantID), stringField("day_of_week", dayOfWeek))
	}
	return mapWorkingHour(row), nil
}

func (s *Store) UpdateWorkingHour(ctx context.Context, item *domain.WorkingHour, actorID *uuid.UUID) (*domain.WorkingHour, error) {
	row, err := s.getQueries(ctx).UpdateWorkingHour(ctx, sqlc.UpdateWorkingHourParams{
		TenantID:     item.TenantID,
		ID:           item.ID,
		BranchID:     uuidFromPtr(item.BranchID),
		UserID:       uuidFromPtr(item.UserID),
		DayOfWeek:    item.DayOfWeek,
		IsWorkingDay: item.IsWorkingDay,
		StartTime:    timeFromClockString(item.StartTime),
		EndTime:      timeFromClockString(item.EndTime),
		BreakMinutes: item.BreakMinutes,
		UpdatedBy:    uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "update working hour", err, tenantIDField(item.TenantID), stringField("working_hour_id", item.ID.String()))
	}
	return mapWorkingHour(row), nil
}

func (s *Store) DeleteWorkingHour(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteWorkingHour(ctx, sqlc.SoftDeleteWorkingHourParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete working hour", err, tenantIDField(tenantID), stringField("working_hour_id", id.String()))
	}
	return nil
}

func (s *Store) CopyTenantWorkingHoursToBranch(ctx context.Context, tenantID uuid.UUID, branchID uuid.UUID, actorID *uuid.UUID) ([]*domain.WorkingHour, error) {
	queries := s.getQueries(ctx)
	if err := queries.SoftDeleteBranchWorkingHours(ctx, sqlc.SoftDeleteBranchWorkingHoursParams{TenantID: tenantID, BranchID: uuidFromPtr(&branchID), UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return nil, s.logDBError(ctx, "clear branch working hours", err, tenantIDField(tenantID), stringField("branch_id", branchID.String()))
	}
	rows, err := queries.CopyTenantWorkingHoursToBranch(ctx, sqlc.CopyTenantWorkingHoursToBranchParams{TenantID: tenantID, BranchID: uuidFromPtr(&branchID), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "copy tenant working hours to branch", err, tenantIDField(tenantID), stringField("branch_id", branchID.String()))
	}
	return mapWorkingHours(rows), nil
}
