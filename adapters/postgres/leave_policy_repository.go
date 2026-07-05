package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func (s *Store) CreateLeavePolicy(ctx context.Context, item *domain.LeavePolicy, actorID *uuid.UUID) (*domain.LeavePolicy, error) {
	row, err := s.getQueries(ctx).CreateLeavePolicy(ctx, sqlc.CreateLeavePolicyParams{
		TenantID:             item.TenantID,
		LeaveTypeID:          item.LeaveTypeID,
		FyID:                 item.FYID,
		TotalDays:            numericFromFloat(item.TotalDays),
		AllocationType:       item.AllocationType,
		Jan:                  item.Jan,
		Feb:                  item.Feb,
		Mar:                  item.Mar,
		Apr:                  item.Apr,
		May:                  item.May,
		Jun:                  item.Jun,
		Jul:                  item.Jul,
		Aug:                  item.Aug,
		Sep:                  item.Sep,
		Oct:                  item.Oct,
		Nov:                  item.Nov,
		Dec:                  item.Dec,
		IsSandwichApplicable: item.IsSandwichApplicable,
		CreatedBy:            uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "create leave policy", err, tenantIDField(item.TenantID), stringField("leave_type_id", item.LeaveTypeID.String()), stringField("financial_year_id", item.FYID.String()))
	}
	return mapLeavePolicy(row), nil
}

func (s *Store) ListLeavePolicies(ctx context.Context, tenantID uuid.UUID) ([]*domain.LeavePolicy, error) {
	rows, err := s.getQueries(ctx).ListLeavePolicies(ctx, tenantID)
	if err != nil {
		return nil, s.logDBError(ctx, "list leave policies", err, tenantIDField(tenantID))
	}
	return mapLeavePolicies(rows), nil
}

func (s *Store) GetLeavePolicy(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.LeavePolicy, error) {
	row, err := s.getQueries(ctx).GetLeavePolicy(ctx, sqlc.GetLeavePolicyParams{TenantID: tenantID, ID: id})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrLeavePolicyNotFound
		}
		return nil, s.logDBError(ctx, "get leave policy", err, tenantIDField(tenantID), stringField("leave_policy_id", id.String()))
	}
	return mapLeavePolicy(row), nil
}

func (s *Store) GetLeavePolicyByTypeAndFY(ctx context.Context, tenantID uuid.UUID, leaveTypeID uuid.UUID, fyID uuid.UUID) (*domain.LeavePolicy, error) {
	row, err := s.getQueries(ctx).GetLeavePolicyByTypeAndFY(ctx, sqlc.GetLeavePolicyByTypeAndFYParams{TenantID: tenantID, LeaveTypeID: leaveTypeID, FyID: fyID})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrLeavePolicyNotFound
		}
		return nil, s.logDBError(ctx, "get leave policy by type and financial year", err, tenantIDField(tenantID), stringField("leave_type_id", leaveTypeID.String()), stringField("financial_year_id", fyID.String()))
	}
	return mapLeavePolicy(row), nil
}

func (s *Store) UpdateLeavePolicy(ctx context.Context, item *domain.LeavePolicy, actorID *uuid.UUID) (*domain.LeavePolicy, error) {
	row, err := s.getQueries(ctx).UpdateLeavePolicy(ctx, sqlc.UpdateLeavePolicyParams{
		TenantID:             item.TenantID,
		ID:                   item.ID,
		LeaveTypeID:          item.LeaveTypeID,
		FyID:                 item.FYID,
		TotalDays:            numericFromFloat(item.TotalDays),
		AllocationType:       item.AllocationType,
		Jan:                  item.Jan,
		Feb:                  item.Feb,
		Mar:                  item.Mar,
		Apr:                  item.Apr,
		May:                  item.May,
		Jun:                  item.Jun,
		Jul:                  item.Jul,
		Aug:                  item.Aug,
		Sep:                  item.Sep,
		Oct:                  item.Oct,
		Nov:                  item.Nov,
		Dec:                  item.Dec,
		IsSandwichApplicable: item.IsSandwichApplicable,
		UpdatedBy:            uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "update leave policy", err, tenantIDField(item.TenantID), stringField("leave_policy_id", item.ID.String()))
	}
	return mapLeavePolicy(row), nil
}

func (s *Store) DeleteLeavePolicy(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteLeavePolicy(ctx, sqlc.SoftDeleteLeavePolicyParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete leave policy", err, tenantIDField(tenantID), stringField("leave_policy_id", id.String()))
	}
	return nil
}
