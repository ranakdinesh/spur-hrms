package postgres

import (
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func mapLeaveType(row sqlc.HrmsLeaveType) *domain.LeaveType {
	return &domain.LeaveType{
		ID:                   row.ID,
		TenantID:             row.TenantID,
		Name:                 row.Name,
		Shortcode:            ptrFromText(row.Shortcode),
		Description:          ptrFromText(row.Description),
		IsPaid:               row.IsPaid,
		IsCarryForward:       row.IsCarryForward,
		MaxCarryForward:      row.MaxCarryForward,
		IsConsecutiveLimit:   row.IsConsecutiveLimit,
		ConsecutiveDaysLimit: row.ConsecutiveDaysLimit,
		IsEnabled:            row.IsEnabled,
		IsSystem:             row.IsSystem,
		Inactive:             row.Inactive,
		CreatedAt:            timeFromTimestamptz(row.CreatedAt),
		CreatedBy:            ptrFromUUID(row.CreatedBy),
		UpdatedAt:            timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:            ptrFromUUID(row.UpdatedBy),
	}
}

func mapLeaveTypes(rows []sqlc.HrmsLeaveType) []*domain.LeaveType {
	items := make([]*domain.LeaveType, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapLeaveType(row))
	}
	return items
}
