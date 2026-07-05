package ports

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

type AIActionLayerRepo interface {
	UpsertAISignalLog(ctx context.Context, item *domain.AISignalLog, actorID *uuid.UUID) (*domain.AISignalLog, error)
	ListAISignalLogs(ctx context.Context, filter domain.AIActionFilter) ([]*domain.AISignalLog, error)
	UpdateAISignalStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, errorMessage *string, actorID *uuid.UUID) (*domain.AISignalLog, error)
	UpsertAIAgentActionLog(ctx context.Context, item *domain.AIAgentActionLog, actorID *uuid.UUID) (*domain.AIAgentActionLog, error)
	ListAIAgentActionLogs(ctx context.Context, filter domain.AIActionFilter) ([]*domain.AIAgentActionLog, error)
	GetAIAgentActionLog(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.AIAgentActionLog, error)
	UpdateAIAgentActionStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, failureMessage *string, actorID *uuid.UUID) (*domain.AIAgentActionLog, error)
	CreateAIHumanOverride(ctx context.Context, item *domain.AIHumanOverride, actorID *uuid.UUID) (*domain.AIHumanOverride, error)
	ListAIHumanOverrides(ctx context.Context, filter domain.AIActionFilter) ([]*domain.AIHumanOverride, error)
	UpsertAIEventOutbox(ctx context.Context, item *domain.AIEventOutbox, actorID *uuid.UUID) (*domain.AIEventOutbox, error)
	ListAIEventOutbox(ctx context.Context, filter domain.AIActionFilter) ([]*domain.AIEventOutbox, error)
	UpdateAIEventOutboxStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, lastError *string, actorID *uuid.UUID) (*domain.AIEventOutbox, error)
}

type AIEventPublisherPort interface {
	PublishAIEvent(ctx context.Context, event *domain.AIEventOutbox) error
}

type AIActionSidecarPort interface {
	ProposeActions(ctx context.Context, tenantID uuid.UUID, signals []*domain.AISignalLog) ([]*domain.AIAgentActionLog, error)
}

type AISignalCommand struct {
	TenantID         uuid.UUID       `json:"tenant_id"`
	SignalKey        string          `json:"signal_key"`
	SignalType       string          `json:"signal_type"`
	SourceModule     string          `json:"source_module"`
	SourceEvent      string          `json:"source_event"`
	Severity         string          `json:"severity"`
	ProcessingStatus string          `json:"processing_status"`
	EntityType       *string         `json:"entity_type,omitempty"`
	EntityID         *uuid.UUID      `json:"entity_id,omitempty"`
	EmployeeUserID   *uuid.UUID      `json:"employee_user_id,omitempty"`
	VisibilityScope  string          `json:"visibility_scope"`
	IdempotencyKey   *string         `json:"idempotency_key,omitempty"`
	CorrelationID    *string         `json:"correlation_id,omitempty"`
	Payload          json.RawMessage `json:"payload,omitempty"`
	Explainability   json.RawMessage `json:"explainability,omitempty"`
	OccurredAt       string          `json:"occurred_at,omitempty"`
	ActorID          *uuid.UUID      `json:"-"`
}

type AIActionCommand struct {
	TenantID            uuid.UUID       `json:"tenant_id"`
	ActionKey           string          `json:"action_key"`
	AgentKey            string          `json:"agent_key"`
	AgentName           string          `json:"agent_name"`
	ActionType          string          `json:"action_type"`
	Status              string          `json:"status"`
	Severity            string          `json:"severity"`
	Title               string          `json:"title"`
	Summary             string          `json:"summary"`
	InsightID           *uuid.UUID      `json:"insight_id,omitempty"`
	SignalID            *uuid.UUID      `json:"signal_id,omitempty"`
	EntityType          *string         `json:"entity_type,omitempty"`
	EntityID            *uuid.UUID      `json:"entity_id,omitempty"`
	EmployeeUserID      *uuid.UUID      `json:"employee_user_id,omitempty"`
	VisibilityScope     string          `json:"visibility_scope"`
	ProposedAction      json.RawMessage `json:"proposed_action,omitempty"`
	InputSnapshot       json.RawMessage `json:"input_snapshot,omitempty"`
	OutputSnapshot      json.RawMessage `json:"output_snapshot,omitempty"`
	Explainability      json.RawMessage `json:"explainability,omitempty"`
	ConfidenceScore     float64         `json:"confidence_score"`
	ModelVersion        *string         `json:"model_version,omitempty"`
	SidecarRunID        *string         `json:"sidecar_run_id,omitempty"`
	RequiresHumanReview bool            `json:"requires_human_review"`
	ActorID             *uuid.UUID      `json:"-"`
}

type AIStatusCommand struct {
	TenantID       uuid.UUID  `json:"tenant_id"`
	ID             uuid.UUID  `json:"id"`
	Status         string     `json:"status"`
	FailureMessage *string    `json:"failure_message,omitempty"`
	ActorID        *uuid.UUID `json:"-"`
}

type AIOverrideCommand struct {
	TenantID       uuid.UUID       `json:"tenant_id"`
	InsightID      *uuid.UUID      `json:"insight_id,omitempty"`
	ActionID       *uuid.UUID      `json:"action_id,omitempty"`
	OverrideType   string          `json:"override_type"`
	OriginalStatus *string         `json:"original_status,omitempty"`
	OverrideStatus string          `json:"override_status"`
	Reason         string          `json:"reason"`
	Decision       string          `json:"decision"`
	Metadata       json.RawMessage `json:"metadata,omitempty"`
	ActorID        *uuid.UUID      `json:"-"`
}

type AIWorkflowEventCommand struct {
	TenantID        uuid.UUID       `json:"tenant_id"`
	EventKey        string          `json:"event_key"`
	EventType       string          `json:"event_type"`
	SourceModule    string          `json:"source_module"`
	SourceEvent     string          `json:"source_event"`
	SignalType      string          `json:"signal_type"`
	Severity        string          `json:"severity"`
	EntityType      *string         `json:"entity_type,omitempty"`
	EntityID        *uuid.UUID      `json:"entity_id,omitempty"`
	EmployeeUserID  *uuid.UUID      `json:"employee_user_id,omitempty"`
	VisibilityScope string          `json:"visibility_scope"`
	Payload         json.RawMessage `json:"payload,omitempty"`
	Explainability  json.RawMessage `json:"explainability,omitempty"`
	CorrelationID   *string         `json:"correlation_id,omitempty"`
	ActorID         *uuid.UUID      `json:"-"`
}
