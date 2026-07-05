package ports

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

type TenantOperationGovernanceRepo interface {
	CreateTenantOperationRequest(ctx context.Context, item *domain.TenantOperationRequest, actorID *uuid.UUID) (*domain.TenantOperationRequest, error)
	GetTenantOperationRequest(ctx context.Context, id uuid.UUID) (*domain.TenantOperationRequest, error)
	ListTenantOperationRequests(ctx context.Context, filter domain.TenantOperationFilter) ([]*domain.TenantOperationRequest, error)
	UpdateTenantOperationRequestStatus(ctx context.Context, id uuid.UUID, status string, approvedBy *uuid.UUID, completedBy *uuid.UUID, backupConfirmed *bool, validationResults json.RawMessage, rollbackMetadata json.RawMessage, metadata json.RawMessage, actorID *uuid.UUID) (*domain.TenantOperationRequest, error)
	CreateTenantOperationEvent(ctx context.Context, event *domain.TenantOperationEvent, actorID *uuid.UUID) (*domain.TenantOperationEvent, error)
	ListTenantOperationEvents(ctx context.Context, requestID uuid.UUID) ([]*domain.TenantOperationEvent, error)
}

type TenantOperationCommand struct {
	ID                uuid.UUID       `json:"id,omitempty"`
	OperationType     string          `json:"operation_type"`
	Title             string          `json:"title"`
	TargetTenantID    *uuid.UUID      `json:"target_tenant_id,omitempty"`
	TargetTenantName  *string         `json:"target_tenant_name,omitempty"`
	TargetTenantCode  *string         `json:"target_tenant_code,omitempty"`
	RiskLevel         string          `json:"risk_level"`
	Reason            string          `json:"reason"`
	ApprovalRequired  bool            `json:"approval_required"`
	BackupConfirmed   bool            `json:"backup_confirmed"`
	RetentionUntil    *string         `json:"retention_until,omitempty"`
	RequestPayload    json.RawMessage `json:"request_payload,omitempty"`
	ValidationResults json.RawMessage `json:"validation_results,omitempty"`
	RollbackMetadata  json.RawMessage `json:"rollback_metadata,omitempty"`
	Metadata          json.RawMessage `json:"metadata,omitempty"`
	ActorID           *uuid.UUID      `json:"-"`
}

type TenantOperationActionCommand struct {
	ID                uuid.UUID       `json:"id"`
	Action            string          `json:"action"`
	Remarks           *string         `json:"remarks,omitempty"`
	BackupConfirmed   *bool           `json:"backup_confirmed,omitempty"`
	ValidationResults json.RawMessage `json:"validation_results,omitempty"`
	RollbackMetadata  json.RawMessage `json:"rollback_metadata,omitempty"`
	Metadata          json.RawMessage `json:"metadata,omitempty"`
	ActorID           *uuid.UUID      `json:"-"`
}
