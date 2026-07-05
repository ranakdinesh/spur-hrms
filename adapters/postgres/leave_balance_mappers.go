package postgres

import (
	"encoding/json"

	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func mapLeaveBalance(row sqlc.HrmsLeaveBalance) *domain.LeaveBalance {
	return &domain.LeaveBalance{
		ID:          row.ID,
		TenantID:    row.TenantID,
		UserID:      row.UserID,
		LeaveTypeID: row.LeaveTypeID,
		FYID:        row.FyID,
		TotalDays:   floatFromNumeric(row.TotalDays),
		UsedDays:    floatFromNumeric(row.UsedDays),
		PendingDays: floatFromNumeric(row.PendingDays),
		BalanceDays: floatFromNumeric(row.BalanceDays),
		Inactive:    row.Inactive,
		CreatedAt:   timeFromTimestamptz(row.CreatedAt),
		CreatedBy:   ptrFromUUID(row.CreatedBy),
		UpdatedAt:   timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:   ptrFromUUID(row.UpdatedBy),
	}
}

func mapLeaveBalances(rows []sqlc.HrmsLeaveBalance) []*domain.LeaveBalance {
	items := make([]*domain.LeaveBalance, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapLeaveBalance(row))
	}
	return items
}

func mapLeaveLedgerEntry(row sqlc.HrmsLeaveLedger) (*domain.LeaveLedgerEntry, error) {
	metadata := map[string]any{}
	if len(row.Metadata) > 0 {
		if err := json.Unmarshal(row.Metadata, &metadata); err != nil {
			return nil, err
		}
	}
	return &domain.LeaveLedgerEntry{
		ID:              row.ID,
		TenantID:        row.TenantID,
		UserID:          row.UserID,
		LeaveTypeID:     row.LeaveTypeID,
		FYID:            row.FyID,
		LeaveID:         ptrFromUUID(row.LeaveID),
		TransactionType: row.TransactionType,
		Days:            floatFromNumeric(row.Days),
		Remarks:         ptrFromText(row.Remarks),
		SourceType:      row.SourceType,
		SourceID:        ptrFromUUID(row.SourceID),
		BalanceBefore:   floatPtrFromNumeric(row.BalanceBefore),
		BalanceAfter:    floatPtrFromNumeric(row.BalanceAfter),
		PendingBefore:   floatPtrFromNumeric(row.PendingBefore),
		PendingAfter:    floatPtrFromNumeric(row.PendingAfter),
		UsedBefore:      floatPtrFromNumeric(row.UsedBefore),
		UsedAfter:       floatPtrFromNumeric(row.UsedAfter),
		Metadata:        metadata,
		Inactive:        row.Inactive,
		CreatedAt:       timeFromTimestamptz(row.CreatedAt),
		CreatedBy:       ptrFromUUID(row.CreatedBy),
	}, nil
}

func mapLeaveLedgerEntries(rows []sqlc.HrmsLeaveLedger) ([]*domain.LeaveLedgerEntry, error) {
	items := make([]*domain.LeaveLedgerEntry, 0, len(rows))
	for _, row := range rows {
		item, err := mapLeaveLedgerEntry(row)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}
