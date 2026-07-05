package postgres

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func (s *Store) CreateTenantOperationRequest(ctx context.Context, item *domain.TenantOperationRequest, actorID *uuid.UUID) (*domain.TenantOperationRequest, error) {
	row, err := s.getQueries(ctx).CreateTenantOperationRequest(ctx, sqlc.CreateTenantOperationRequestParams{
		OperationNumber:   item.OperationNumber,
		OperationType:     item.OperationType,
		Title:             item.Title,
		TargetTenantID:    uuidFromPtr(item.TargetTenantID),
		TargetTenantName:  textFromPtr(item.TargetTenantName),
		TargetTenantCode:  textFromPtr(item.TargetTenantCode),
		Status:            item.Status,
		RiskLevel:         item.RiskLevel,
		Reason:            item.Reason,
		ActorID:           uuidFromPtr(actorID),
		ApprovalRequired:  item.ApprovalRequired,
		BackupRequired:    item.BackupRequired,
		BackupConfirmed:   item.BackupConfirmed,
		RetentionUntil:    timestamptzFromPtr(item.RetentionUntil),
		RequestPayload:    json.RawMessage(item.RequestPayload),
		ValidationResults: json.RawMessage(item.ValidationResults),
		RollbackMetadata:  json.RawMessage(item.RollbackMetadata),
		Metadata:          json.RawMessage(item.Metadata),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "create tenant operation request", err, optionalTenantIDField(item.TargetTenantID), stringField("operation_type", item.OperationType))
	}
	return mapTenantOperationRequest(row), nil
}

func (s *Store) GetTenantOperationRequest(ctx context.Context, id uuid.UUID) (*domain.TenantOperationRequest, error) {
	row, err := s.getQueries(ctx).GetTenantOperationRequest(ctx, id)
	if err != nil {
		return nil, s.logDBError(ctx, "get tenant operation request", err, stringField("request_id", id.String()))
	}
	return mapTenantOperationRequest(row), nil
}

func (s *Store) ListTenantOperationRequests(ctx context.Context, filter domain.TenantOperationFilter) ([]*domain.TenantOperationRequest, error) {
	rows, err := s.getQueries(ctx).ListTenantOperationRequests(ctx, sqlc.ListTenantOperationRequestsParams{
		Status:         textFromPtr(filter.Status),
		OperationType:  textFromPtr(filter.OperationType),
		RiskLevel:      textFromPtr(filter.RiskLevel),
		TargetTenantID: uuidFromPtr(filter.TargetTenantID),
		Search:         textFromPtr(filter.Search),
		RowOffset:      filter.Offset,
		RowLimit:       limitOrDefault(filter.Limit),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "list tenant operation requests", err, optionalTenantIDField(filter.TargetTenantID))
	}
	return mapTenantOperationRequests(rows), nil
}

func (s *Store) UpdateTenantOperationRequestStatus(ctx context.Context, id uuid.UUID, status string, approvedBy *uuid.UUID, completedBy *uuid.UUID, backupConfirmed *bool, validationResults json.RawMessage, rollbackMetadata json.RawMessage, metadata json.RawMessage, actorID *uuid.UUID) (*domain.TenantOperationRequest, error) {
	row, err := s.getQueries(ctx).UpdateTenantOperationRequestStatus(ctx, sqlc.UpdateTenantOperationRequestStatusParams{
		Status:            status,
		ApprovedBy:        uuidFromPtr(approvedBy),
		CompletedBy:       uuidFromPtr(completedBy),
		BackupConfirmed:   boolFromPtr(backupConfirmed),
		ValidationResults: nullableJSON(validationResults),
		RollbackMetadata:  nullableJSON(rollbackMetadata),
		Metadata:          nullableJSON(metadata),
		ActorID:           uuidFromPtr(actorID),
		ID:                id,
	})
	if err != nil {
		return nil, s.logDBError(ctx, "update tenant operation request status", err, stringField("request_id", id.String()), stringField("status", status))
	}
	return mapTenantOperationRequest(row), nil
}

func (s *Store) CreateTenantOperationEvent(ctx context.Context, event *domain.TenantOperationEvent, actorID *uuid.UUID) (*domain.TenantOperationEvent, error) {
	row, err := s.getQueries(ctx).CreateTenantOperationEvent(ctx, sqlc.CreateTenantOperationEventParams{
		RequestID:  event.RequestID,
		Action:     event.Action,
		FromStatus: textFromPtr(event.FromStatus),
		ToStatus:   textFromPtr(event.ToStatus),
		ActorID:    uuidFromPtr(actorID),
		Remarks:    textFromPtr(event.Remarks),
		Metadata:   json.RawMessage(event.Metadata),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "create tenant operation event", err, stringField("request_id", event.RequestID.String()), stringField("action", event.Action))
	}
	return mapTenantOperationEvent(row), nil
}

func (s *Store) ListTenantOperationEvents(ctx context.Context, requestID uuid.UUID) ([]*domain.TenantOperationEvent, error) {
	rows, err := s.getQueries(ctx).ListTenantOperationEvents(ctx, requestID)
	if err != nil {
		return nil, s.logDBError(ctx, "list tenant operation events", err, stringField("request_id", requestID.String()))
	}
	return mapTenantOperationEvents(rows), nil
}

func nullableJSON(value json.RawMessage) []byte {
	if len(value) == 0 {
		return nil
	}
	return value
}
