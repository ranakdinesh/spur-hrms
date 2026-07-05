package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func (s *Store) CreateLeaveType(ctx context.Context, item *domain.LeaveType, actorID *uuid.UUID) (*domain.LeaveType, error) {
	row, err := s.getQueries(ctx).CreateLeaveType(ctx, sqlc.CreateLeaveTypeParams{
		TenantID:             item.TenantID,
		Name:                 item.Name,
		Shortcode:            textFromPtr(item.Shortcode),
		Description:          textFromPtr(item.Description),
		IsPaid:               item.IsPaid,
		IsCarryForward:       item.IsCarryForward,
		MaxCarryForward:      item.MaxCarryForward,
		IsConsecutiveLimit:   item.IsConsecutiveLimit,
		ConsecutiveDaysLimit: item.ConsecutiveDaysLimit,
		IsEnabled:            item.IsEnabled,
		IsSystem:             item.IsSystem,
		CreatedBy:            uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "create leave type", err, tenantIDField(item.TenantID), stringField("leave_type_name", item.Name))
	}
	return mapLeaveType(row), nil
}

func (s *Store) UpsertSystemLeaveType(ctx context.Context, item *domain.LeaveType, actorID *uuid.UUID) (*domain.LeaveType, error) {
	row, err := s.getQueries(ctx).UpsertSystemLeaveType(ctx, sqlc.UpsertSystemLeaveTypeParams{
		TenantID:             item.TenantID,
		Name:                 item.Name,
		Shortcode:            textFromPtr(item.Shortcode),
		Description:          textFromPtr(item.Description),
		IsPaid:               item.IsPaid,
		IsCarryForward:       item.IsCarryForward,
		MaxCarryForward:      item.MaxCarryForward,
		IsConsecutiveLimit:   item.IsConsecutiveLimit,
		ConsecutiveDaysLimit: item.ConsecutiveDaysLimit,
		IsEnabled:            item.IsEnabled,
		CreatedBy:            uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "upsert system leave type", err, tenantIDField(item.TenantID), stringField("leave_type_name", item.Name))
	}
	return mapLeaveType(row), nil
}

func (s *Store) ListLeaveTypes(ctx context.Context, tenantID uuid.UUID) ([]*domain.LeaveType, error) {
	rows, err := s.getQueries(ctx).ListLeaveTypes(ctx, tenantID)
	if err != nil {
		return nil, s.logDBError(ctx, "list leave types", err, tenantIDField(tenantID))
	}
	return mapLeaveTypes(rows), nil
}

func (s *Store) GetLeaveType(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.LeaveType, error) {
	row, err := s.getQueries(ctx).GetLeaveType(ctx, sqlc.GetLeaveTypeParams{TenantID: tenantID, ID: id})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrLeaveTypeNotFound
		}
		return nil, s.logDBError(ctx, "get leave type", err, tenantIDField(tenantID), stringField("leave_type_id", id.String()))
	}
	return mapLeaveType(row), nil
}

func (s *Store) GetLeaveTypeByShortcode(ctx context.Context, tenantID uuid.UUID, shortcode string) (*domain.LeaveType, error) {
	row, err := s.getQueries(ctx).GetLeaveTypeByShortcode(ctx, sqlc.GetLeaveTypeByShortcodeParams{TenantID: tenantID, Lower: shortcode})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrLeaveTypeNotFound
		}
		return nil, s.logDBError(ctx, "get leave type by shortcode", err, tenantIDField(tenantID), stringField("leave_type_shortcode", shortcode))
	}
	return mapLeaveType(row), nil
}

func (s *Store) UpdateLeaveType(ctx context.Context, item *domain.LeaveType, actorID *uuid.UUID) (*domain.LeaveType, error) {
	row, err := s.getQueries(ctx).UpdateLeaveType(ctx, sqlc.UpdateLeaveTypeParams{
		TenantID:             item.TenantID,
		ID:                   item.ID,
		Name:                 item.Name,
		Shortcode:            textFromPtr(item.Shortcode),
		Description:          textFromPtr(item.Description),
		IsPaid:               item.IsPaid,
		IsCarryForward:       item.IsCarryForward,
		MaxCarryForward:      item.MaxCarryForward,
		IsConsecutiveLimit:   item.IsConsecutiveLimit,
		ConsecutiveDaysLimit: item.ConsecutiveDaysLimit,
		IsEnabled:            item.IsEnabled,
		UpdatedBy:            uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "update leave type", err, tenantIDField(item.TenantID), stringField("leave_type_id", item.ID.String()))
	}
	return mapLeaveType(row), nil
}

func (s *Store) DeleteLeaveType(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteLeaveType(ctx, sqlc.SoftDeleteLeaveTypeParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete leave type", err, tenantIDField(tenantID), stringField("leave_type_id", id.String()))
	}
	return nil
}
