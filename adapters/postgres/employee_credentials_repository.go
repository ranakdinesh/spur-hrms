package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func (s *Store) CreateEmployeeCredentialEvent(ctx context.Context, event *domain.EmployeeCredentialEvent) (*domain.EmployeeCredentialEvent, error) {
	row, err := s.getQueries(ctx).CreateEmployeeCredentialEvent(ctx, sqlc.CreateEmployeeCredentialEventParams{
		TenantID:        event.TenantID,
		EmployeeID:      event.EmployeeID,
		UserID:          event.UserID,
		EventType:       event.EventType,
		DeliveryChannel: event.DeliveryChannel,
		DeliveryTarget:  event.DeliveryTarget,
		Status:          event.Status,
		FailureReason:   textFromPtr(event.FailureReason),
		CreatedBy:       uuidFromPtr(event.CreatedBy),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "create employee credential event", err, tenantIDField(event.TenantID), stringField("employee_id", event.EmployeeID.String()), stringField("event_type", event.EventType))
	}
	return mapEmployeeCredentialEvent(row), nil
}

func (s *Store) ListEmployeeCredentialEvents(ctx context.Context, tenantID uuid.UUID, employeeID uuid.UUID, limit int32) ([]*domain.EmployeeCredentialEvent, error) {
	rows, err := s.getQueries(ctx).ListEmployeeCredentialEvents(ctx, sqlc.ListEmployeeCredentialEventsParams{TenantID: tenantID, EmployeeID: employeeID, LimitRows: limit})
	if err != nil {
		return nil, s.logDBError(ctx, "list employee credential events", err, tenantIDField(tenantID), stringField("employee_id", employeeID.String()))
	}
	return mapEmployeeCredentialEvents(rows), nil
}
