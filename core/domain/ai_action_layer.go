package domain

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidAISignal       = errors.New("ai signal is invalid")
	ErrInvalidAIAction       = errors.New("ai action is invalid")
	ErrInvalidAIOverride     = errors.New("ai human override is invalid")
	ErrInvalidAIEvent        = errors.New("ai event is invalid")
	ErrAIActionNotFound      = errors.New("ai action not found")
	ErrInvalidAIStatus       = errors.New("ai status is invalid")
	ErrInvalidAIVisibility   = errors.New("ai visibility is invalid")
	ErrInvalidAISeverity     = errors.New("ai severity is invalid")
	ErrInvalidAIDecision     = errors.New("ai override decision is invalid")
	ErrInvalidAIEventStatus  = errors.New("ai event status is invalid")
	ErrInvalidAISignalStatus = errors.New("ai signal status is invalid")
)

const (
	AIVisibilityEmployee         = "employee"
	AIVisibilityManagerAggregate = "manager_aggregate"
	AIVisibilityHR               = "hr"
	AIVisibilityAdmin            = "admin"

	AISignalStatusNew       = "new"
	AISignalStatusQueued    = "queued"
	AISignalStatusProcessed = "processed"
	AISignalStatusIgnored   = "ignored"
	AISignalStatusFailed    = "failed"

	AIActionStatusProposed   = "proposed"
	AIActionStatusQueued     = "queued"
	AIActionStatusReviewing  = "reviewing"
	AIActionStatusApproved   = "approved"
	AIActionStatusRejected   = "rejected"
	AIActionStatusExecuted   = "executed"
	AIActionStatusFailed     = "failed"
	AIActionStatusOverridden = "overridden"
	AIActionStatusCancelled  = "cancelled"

	AIOverrideDecisionAccepted     = "accepted"
	AIOverrideDecisionRejected     = "rejected"
	AIOverrideDecisionReplaced     = "replaced"
	AIOverrideDecisionManualAction = "manual_action"

	AIEventStatusPending   = "pending"
	AIEventStatusPublished = "published"
	AIEventStatusFailed    = "failed"
	AIEventStatusSkipped   = "skipped"

	AIEventTargetRedisStream = "redis_stream"
)

type AISignalLog struct {
	ID               uuid.UUID       `json:"id"`
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
	OccurredAt       time.Time       `json:"occurred_at"`
	ProcessedAt      *time.Time      `json:"processed_at,omitempty"`
	ErrorMessage     *string         `json:"error_message,omitempty"`
	Inactive         bool            `json:"inactive"`
	CreatedAt        time.Time       `json:"created_at"`
	CreatedBy        *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt        time.Time       `json:"updated_at"`
	UpdatedBy        *uuid.UUID      `json:"updated_by,omitempty"`
}

type AIAgentActionLog struct {
	ID                  uuid.UUID       `json:"id"`
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
	ReviewedBy          *uuid.UUID      `json:"reviewed_by,omitempty"`
	ReviewedAt          *time.Time      `json:"reviewed_at,omitempty"`
	ExecutedAt          *time.Time      `json:"executed_at,omitempty"`
	FailedAt            *time.Time      `json:"failed_at,omitempty"`
	FailureMessage      *string         `json:"failure_message,omitempty"`
	Inactive            bool            `json:"inactive"`
	CreatedAt           time.Time       `json:"created_at"`
	CreatedBy           *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt           time.Time       `json:"updated_at"`
	UpdatedBy           *uuid.UUID      `json:"updated_by,omitempty"`
}

type AIHumanOverride struct {
	ID             uuid.UUID       `json:"id"`
	TenantID       uuid.UUID       `json:"tenant_id"`
	InsightID      *uuid.UUID      `json:"insight_id,omitempty"`
	ActionID       *uuid.UUID      `json:"action_id,omitempty"`
	OverrideType   string          `json:"override_type"`
	OriginalStatus *string         `json:"original_status,omitempty"`
	OverrideStatus string          `json:"override_status"`
	Reason         string          `json:"reason"`
	Decision       string          `json:"decision"`
	Metadata       json.RawMessage `json:"metadata,omitempty"`
	Inactive       bool            `json:"inactive"`
	CreatedAt      time.Time       `json:"created_at"`
	CreatedBy      *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt      time.Time       `json:"updated_at"`
	UpdatedBy      *uuid.UUID      `json:"updated_by,omitempty"`
}

type AIEventOutbox struct {
	ID            uuid.UUID       `json:"id"`
	TenantID      uuid.UUID       `json:"tenant_id"`
	EventKey      string          `json:"event_key"`
	EventType     string          `json:"event_type"`
	TargetBus     string          `json:"target_bus"`
	Status        string          `json:"status"`
	Payload       json.RawMessage `json:"payload,omitempty"`
	CorrelationID *string         `json:"correlation_id,omitempty"`
	Attempts      int32           `json:"attempts"`
	PublishedAt   *time.Time      `json:"published_at,omitempty"`
	LastError     *string         `json:"last_error,omitempty"`
	Inactive      bool            `json:"inactive"`
	CreatedAt     time.Time       `json:"created_at"`
	CreatedBy     *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt     time.Time       `json:"updated_at"`
	UpdatedBy     *uuid.UUID      `json:"updated_by,omitempty"`
}

type AIActionFilter struct {
	TenantID         uuid.UUID
	Status           *string
	Severity         *string
	AgentKey         *string
	InsightID        *uuid.UUID
	ActionID         *uuid.UUID
	VisibilityScope  *string
	ProcessingStatus *string
	SourceModule     *string
	EventType        *string
	Decision         *string
	Limit            int32
	Offset           int32
}

type AIActionSummary struct {
	SignalsPending  int32 `json:"signals_pending"`
	SignalsFailed   int32 `json:"signals_failed"`
	ActionsProposed int32 `json:"actions_proposed"`
	ActionsApproved int32 `json:"actions_approved"`
	ActionsExecuted int32 `json:"actions_executed"`
	ActionsFailed   int32 `json:"actions_failed"`
	Overrides       int32 `json:"overrides"`
	OutboxPending   int32 `json:"outbox_pending"`
	OutboxFailed    int32 `json:"outbox_failed"`
}

type AIActionWorkspace struct {
	Signals   []*AISignalLog      `json:"signals"`
	Actions   []*AIAgentActionLog `json:"actions"`
	Overrides []*AIHumanOverride  `json:"overrides"`
	Events    []*AIEventOutbox    `json:"events"`
	Summary   AIActionSummary     `json:"summary"`
}

func NewAISignalLog(input AISignalLog) (*AISignalLog, error) {
	input.SignalKey = strings.TrimSpace(input.SignalKey)
	input.SignalType = strings.TrimSpace(input.SignalType)
	input.SourceModule = strings.TrimSpace(input.SourceModule)
	input.SourceEvent = strings.TrimSpace(input.SourceEvent)
	if input.TenantID == uuid.Nil || input.SignalKey == "" || input.SignalType == "" || input.SourceModule == "" || input.SourceEvent == "" {
		return nil, ErrInvalidAISignal
	}
	severity, err := NormalizeAISeverity(input.Severity)
	if err != nil {
		return nil, err
	}
	status, err := NormalizeAISignalStatus(input.ProcessingStatus)
	if err != nil {
		return nil, err
	}
	visibility, err := NormalizeAIVisibility(input.VisibilityScope)
	if err != nil {
		return nil, err
	}
	if input.OccurredAt.IsZero() {
		input.OccurredAt = time.Now().UTC()
	}
	input.Severity = severity
	input.ProcessingStatus = status
	input.VisibilityScope = visibility
	input.Payload = rawDefault(input.Payload, "{}")
	input.Explainability = rawDefault(input.Explainability, "{}")
	return &input, nil
}

func NewAIAgentActionLog(input AIAgentActionLog) (*AIAgentActionLog, error) {
	input.ActionKey = strings.TrimSpace(input.ActionKey)
	input.AgentKey = strings.TrimSpace(input.AgentKey)
	input.AgentName = strings.TrimSpace(input.AgentName)
	input.ActionType = strings.TrimSpace(input.ActionType)
	input.Title = strings.TrimSpace(input.Title)
	input.Summary = strings.TrimSpace(input.Summary)
	if input.TenantID == uuid.Nil || input.ActionKey == "" || input.AgentKey == "" || input.AgentName == "" || input.ActionType == "" || input.Title == "" {
		return nil, ErrInvalidAIAction
	}
	status, err := NormalizeAIActionStatus(input.Status)
	if err != nil {
		return nil, err
	}
	severity, err := NormalizeAISeverity(input.Severity)
	if err != nil {
		return nil, err
	}
	visibility, err := NormalizeAIVisibility(input.VisibilityScope)
	if err != nil {
		return nil, err
	}
	if input.ConfidenceScore < 0 {
		input.ConfidenceScore = 0
	}
	if input.ConfidenceScore > 100 {
		input.ConfidenceScore = 100
	}
	input.Status = status
	input.Severity = severity
	input.VisibilityScope = visibility
	input.ProposedAction = rawDefault(input.ProposedAction, "{}")
	input.InputSnapshot = rawDefault(input.InputSnapshot, "{}")
	input.OutputSnapshot = rawDefault(input.OutputSnapshot, "{}")
	input.Explainability = rawDefault(input.Explainability, "{}")
	return &input, nil
}

func NewAIHumanOverride(input AIHumanOverride) (*AIHumanOverride, error) {
	input.OverrideType = strings.TrimSpace(input.OverrideType)
	input.Reason = strings.TrimSpace(input.Reason)
	if input.TenantID == uuid.Nil || input.OverrideType == "" || input.Reason == "" || (input.InsightID == nil && input.ActionID == nil) {
		return nil, ErrInvalidAIOverride
	}
	status, err := NormalizeAIOverrideStatus(input.OverrideStatus)
	if err != nil {
		return nil, err
	}
	decision, err := NormalizeAIOverrideDecision(input.Decision)
	if err != nil {
		return nil, err
	}
	input.OverrideStatus = status
	input.Decision = decision
	input.Metadata = rawDefault(input.Metadata, "{}")
	return &input, nil
}

func NewAIEventOutbox(input AIEventOutbox) (*AIEventOutbox, error) {
	input.EventKey = strings.TrimSpace(input.EventKey)
	input.EventType = strings.TrimSpace(input.EventType)
	input.TargetBus = strings.TrimSpace(input.TargetBus)
	if input.TenantID == uuid.Nil || input.EventKey == "" || input.EventType == "" {
		return nil, ErrInvalidAIEvent
	}
	if input.TargetBus == "" {
		input.TargetBus = AIEventTargetRedisStream
	}
	status, err := NormalizeAIEventStatus(input.Status)
	if err != nil {
		return nil, err
	}
	input.Status = status
	input.Payload = rawDefault(input.Payload, "{}")
	return &input, nil
}

func NormalizeAISeverity(value string) (string, error) {
	switch strings.TrimSpace(value) {
	case "", InsightSeverityMedium:
		return InsightSeverityMedium, nil
	case InsightSeverityLow, InsightSeverityHigh, InsightSeverityCritical:
		return strings.TrimSpace(value), nil
	default:
		return "", ErrInvalidAISeverity
	}
}

func NormalizeAIVisibility(value string) (string, error) {
	switch strings.TrimSpace(value) {
	case "", AIVisibilityHR:
		return AIVisibilityHR, nil
	case AIVisibilityEmployee, AIVisibilityManagerAggregate, AIVisibilityAdmin:
		return strings.TrimSpace(value), nil
	default:
		return "", ErrInvalidAIVisibility
	}
}

func NormalizeAISignalStatus(value string) (string, error) {
	switch strings.TrimSpace(value) {
	case "", AISignalStatusNew:
		return AISignalStatusNew, nil
	case AISignalStatusQueued, AISignalStatusProcessed, AISignalStatusIgnored, AISignalStatusFailed:
		return strings.TrimSpace(value), nil
	default:
		return "", ErrInvalidAISignalStatus
	}
}

func NormalizeAIActionStatus(value string) (string, error) {
	switch strings.TrimSpace(value) {
	case "", AIActionStatusProposed:
		return AIActionStatusProposed, nil
	case AIActionStatusQueued, AIActionStatusReviewing, AIActionStatusApproved, AIActionStatusRejected, AIActionStatusExecuted, AIActionStatusFailed, AIActionStatusOverridden, AIActionStatusCancelled:
		return strings.TrimSpace(value), nil
	default:
		return "", ErrInvalidAIStatus
	}
}

func NormalizeAIOverrideStatus(value string) (string, error) {
	switch strings.TrimSpace(value) {
	case "", InsightStatusOverridden:
		return InsightStatusOverridden, nil
	case InsightStatusDismissed, InsightStatusResolved, AIActionStatusRejected, AIActionStatusCancelled:
		return strings.TrimSpace(value), nil
	default:
		return "", ErrInvalidAIStatus
	}
}

func NormalizeAIOverrideDecision(value string) (string, error) {
	switch strings.TrimSpace(value) {
	case "", AIOverrideDecisionManualAction:
		return AIOverrideDecisionManualAction, nil
	case AIOverrideDecisionAccepted, AIOverrideDecisionRejected, AIOverrideDecisionReplaced:
		return strings.TrimSpace(value), nil
	default:
		return "", ErrInvalidAIDecision
	}
}

func NormalizeAIEventStatus(value string) (string, error) {
	switch strings.TrimSpace(value) {
	case "", AIEventStatusPending:
		return AIEventStatusPending, nil
	case AIEventStatusPublished, AIEventStatusFailed, AIEventStatusSkipped:
		return strings.TrimSpace(value), nil
	default:
		return "", ErrInvalidAIEventStatus
	}
}

func rawDefault(value json.RawMessage, fallback string) json.RawMessage {
	if len(value) == 0 {
		return json.RawMessage(fallback)
	}
	return value
}
