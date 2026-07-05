package postgres

import (
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func mapEmployeeCredentialEvent(row sqlc.HrmsEmployeeCredentialEvent) *domain.EmployeeCredentialEvent {
	return &domain.EmployeeCredentialEvent{
		ID:              row.ID,
		TenantID:        row.TenantID,
		EmployeeID:      row.EmployeeID,
		UserID:          row.UserID,
		EventType:       row.EventType,
		DeliveryChannel: row.DeliveryChannel,
		DeliveryTarget:  row.DeliveryTarget,
		Status:          row.Status,
		FailureReason:   ptrFromText(row.FailureReason),
		CreatedAt:       row.CreatedAt.Time,
		CreatedBy:       ptrFromUUID(row.CreatedBy),
	}
}

func mapEmployeeCredentialEvents(rows []sqlc.HrmsEmployeeCredentialEvent) []*domain.EmployeeCredentialEvent {
	items := make([]*domain.EmployeeCredentialEvent, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapEmployeeCredentialEvent(row))
	}
	return items
}
