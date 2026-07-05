package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func (s *Store) UpsertLeaveBalance(ctx context.Context, item *domain.LeaveBalance, actorID *uuid.UUID) (*domain.LeaveBalance, error) {
	row, err := s.getQueries(ctx).UpsertLeaveBalance(ctx, sqlc.UpsertLeaveBalanceParams{
		TenantID:    item.TenantID,
		UserID:      item.UserID,
		LeaveTypeID: item.LeaveTypeID,
		FyID:        item.FYID,
		TotalDays:   numericFromFloat(item.TotalDays),
		UsedDays:    numericFromFloat(item.UsedDays),
		PendingDays: numericFromFloat(item.PendingDays),
		CreatedBy:   uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "upsert leave balance", err, tenantIDField(item.TenantID), stringField("user_id", item.UserID.String()))
	}
	return mapLeaveBalance(row), nil
}

func (s *Store) GetLeaveBalance(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID, leaveTypeID uuid.UUID, fyID uuid.UUID) (*domain.LeaveBalance, error) {
	row, err := s.getQueries(ctx).GetLeaveBalance(ctx, sqlc.GetLeaveBalanceParams{TenantID: tenantID, UserID: userID, LeaveTypeID: leaveTypeID, FyID: fyID})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrLeaveBalanceNotFound
		}
		return nil, s.logDBError(ctx, "get leave balance", err, tenantIDField(tenantID), stringField("user_id", userID.String()), stringField("leave_type_id", leaveTypeID.String()))
	}
	return mapLeaveBalance(row), nil
}

func (s *Store) ListLeaveBalancesByUser(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID) ([]*domain.LeaveBalance, error) {
	rows, err := s.getQueries(ctx).ListLeaveBalancesByUser(ctx, sqlc.ListLeaveBalancesByUserParams{TenantID: tenantID, UserID: userID})
	if err != nil {
		return nil, s.logDBError(ctx, "list leave balances by user", err, tenantIDField(tenantID), stringField("user_id", userID.String()))
	}
	return mapLeaveBalances(rows), nil
}

func (s *Store) ListLeaveBalancesByTenantFY(ctx context.Context, tenantID uuid.UUID, fyID uuid.UUID) ([]*domain.LeaveBalance, error) {
	rows, err := s.getQueries(ctx).ListLeaveBalancesByTenantFY(ctx, sqlc.ListLeaveBalancesByTenantFYParams{TenantID: tenantID, FyID: fyID})
	if err != nil {
		return nil, s.logDBError(ctx, "list leave balances by tenant fy", err, tenantIDField(tenantID), stringField("financial_year_id", fyID.String()))
	}
	return mapLeaveBalances(rows), nil
}

func (s *Store) AddLeaveBalanceCredit(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID, leaveTypeID uuid.UUID, fyID uuid.UUID, days float64, actorID *uuid.UUID) (*domain.LeaveBalance, error) {
	row, err := s.getQueries(ctx).AddLeaveBalanceCredit(ctx, sqlc.AddLeaveBalanceCreditParams{TenantID: tenantID, UserID: userID, LeaveTypeID: leaveTypeID, FyID: fyID, TotalDays: numericFromFloat(days), UpdatedBy: uuidFromPtr(actorID)})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrLeaveBalanceNotFound
		}
		return nil, s.logDBError(ctx, "add leave balance credit", err, tenantIDField(tenantID), stringField("user_id", userID.String()))
	}
	return mapLeaveBalance(row), nil
}

func (s *Store) UpdateLeaveBalancePending(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID, leaveTypeID uuid.UUID, fyID uuid.UUID, daysDelta float64, actorID *uuid.UUID) (*domain.LeaveBalance, error) {
	row, err := s.getQueries(ctx).UpdateLeaveBalancePending(ctx, sqlc.UpdateLeaveBalancePendingParams{TenantID: tenantID, UserID: userID, LeaveTypeID: leaveTypeID, FyID: fyID, PendingDays: numericFromFloat(daysDelta), UpdatedBy: uuidFromPtr(actorID)})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrLeaveBalanceNotFound
		}
		return nil, s.logDBError(ctx, "update leave balance pending", err, tenantIDField(tenantID), stringField("user_id", userID.String()))
	}
	return mapLeaveBalance(row), nil
}

func (s *Store) ReverseLeaveBalancePending(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID, leaveTypeID uuid.UUID, fyID uuid.UUID, days float64, actorID *uuid.UUID) (*domain.LeaveBalance, error) {
	row, err := s.getQueries(ctx).ReverseLeaveBalancePending(ctx, sqlc.ReverseLeaveBalancePendingParams{TenantID: tenantID, UserID: userID, LeaveTypeID: leaveTypeID, FyID: fyID, PendingDays: numericFromFloat(days), UpdatedBy: uuidFromPtr(actorID)})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrLeaveBalanceNotFound
		}
		return nil, s.logDBError(ctx, "reverse leave balance pending", err, tenantIDField(tenantID), stringField("user_id", userID.String()))
	}
	return mapLeaveBalance(row), nil
}

func (s *Store) MoveLeaveBalancePendingToUsed(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID, leaveTypeID uuid.UUID, fyID uuid.UUID, days float64, actorID *uuid.UUID) (*domain.LeaveBalance, error) {
	row, err := s.getQueries(ctx).MoveLeaveBalancePendingToUsed(ctx, sqlc.MoveLeaveBalancePendingToUsedParams{TenantID: tenantID, UserID: userID, LeaveTypeID: leaveTypeID, FyID: fyID, PendingDays: numericFromFloat(days), UpdatedBy: uuidFromPtr(actorID)})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrLeaveBalanceNotFound
		}
		return nil, s.logDBError(ctx, "move leave balance pending to used", err, tenantIDField(tenantID), stringField("user_id", userID.String()))
	}
	return mapLeaveBalance(row), nil
}

func (s *Store) CreateLeaveLedgerEntry(ctx context.Context, item *domain.LeaveLedgerEntry, actorID *uuid.UUID) (*domain.LeaveLedgerEntry, error) {
	row, err := s.getQueries(ctx).CreateLeaveLedgerEntry(ctx, sqlc.CreateLeaveLedgerEntryParams{
		TenantID:        item.TenantID,
		UserID:          item.UserID,
		LeaveTypeID:     item.LeaveTypeID,
		FyID:            item.FYID,
		LeaveID:         uuidFromPtr(item.LeaveID),
		TransactionType: item.TransactionType,
		Days:            numericFromFloat(item.Days),
		Remarks:         textFromPtr(item.Remarks),
		SourceType:      item.SourceType,
		SourceID:        uuidFromPtr(item.SourceID),
		BalanceBefore:   numericFromFloatPtr(item.BalanceBefore),
		BalanceAfter:    numericFromFloatPtr(item.BalanceAfter),
		PendingBefore:   numericFromFloatPtr(item.PendingBefore),
		PendingAfter:    numericFromFloatPtr(item.PendingAfter),
		UsedBefore:      numericFromFloatPtr(item.UsedBefore),
		UsedAfter:       numericFromFloatPtr(item.UsedAfter),
		Metadata:        jsonBytesFromMap(item.Metadata),
		CreatedBy:       uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "create leave ledger entry", err, tenantIDField(item.TenantID), stringField("user_id", item.UserID.String()), stringField("source_type", item.SourceType))
	}
	mapped, err := mapLeaveLedgerEntry(row)
	if err != nil {
		return nil, s.logDBError(ctx, "map leave ledger entry", fmt.Errorf("hrms: map leave ledger entry: %w", err), tenantIDField(item.TenantID))
	}
	return mapped, nil
}

func (s *Store) ListLeaveLedgerByUser(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID) ([]*domain.LeaveLedgerEntry, error) {
	rows, err := s.getQueries(ctx).ListLeaveLedgerByUser(ctx, sqlc.ListLeaveLedgerByUserParams{TenantID: tenantID, UserID: userID})
	if err != nil {
		return nil, s.logDBError(ctx, "list leave ledger by user", err, tenantIDField(tenantID), stringField("user_id", userID.String()))
	}
	items, err := mapLeaveLedgerEntries(rows)
	if err != nil {
		return nil, s.logDBError(ctx, "map leave ledger by user", fmt.Errorf("hrms: map leave ledger by user: %w", err), tenantIDField(tenantID))
	}
	return items, nil
}

func (s *Store) ListLeaveLedgerByLeave(ctx context.Context, tenantID uuid.UUID, leaveID uuid.UUID) ([]*domain.LeaveLedgerEntry, error) {
	rows, err := s.getQueries(ctx).ListLeaveLedgerByLeave(ctx, sqlc.ListLeaveLedgerByLeaveParams{TenantID: tenantID, LeaveID: uuidFromPtr(&leaveID)})
	if err != nil {
		return nil, s.logDBError(ctx, "list leave ledger by leave", err, tenantIDField(tenantID), stringField("leave_id", leaveID.String()))
	}
	items, err := mapLeaveLedgerEntries(rows)
	if err != nil {
		return nil, s.logDBError(ctx, "map leave ledger by leave", fmt.Errorf("hrms: map leave ledger by leave: %w", err), tenantIDField(tenantID))
	}
	return items, nil
}

func (s *Store) GetLeaveLedgerBySource(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID, leaveTypeID uuid.UUID, fyID uuid.UUID, sourceType string, sourceID uuid.UUID) (*domain.LeaveLedgerEntry, error) {
	row, err := s.getQueries(ctx).GetLeaveLedgerBySource(ctx, sqlc.GetLeaveLedgerBySourceParams{TenantID: tenantID, UserID: userID, LeaveTypeID: leaveTypeID, FyID: fyID, SourceType: sourceType, SourceID: uuidFromPtr(&sourceID)})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrLeaveLedgerEntryNotFound
		}
		return nil, s.logDBError(ctx, "get leave ledger by source", err, tenantIDField(tenantID), stringField("user_id", userID.String()), stringField("source_type", sourceType))
	}
	item, err := mapLeaveLedgerEntry(row)
	if err != nil {
		return nil, s.logDBError(ctx, "map leave ledger by source", fmt.Errorf("hrms: map leave ledger by source: %w", err), tenantIDField(tenantID))
	}
	return item, nil
}
