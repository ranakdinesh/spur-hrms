package domain

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidInsightID       = errors.New("insight_id is required")
	ErrInvalidInsightKey      = errors.New("insight key is required")
	ErrInvalidInsightStatus   = errors.New("insight status is invalid")
	ErrInvalidInsightSeverity = errors.New("insight severity is invalid")
	ErrInsightNotFound        = errors.New("insight not found")
)

const (
	InsightStatusOpen       = "open"
	InsightStatusReviewing  = "reviewing"
	InsightStatusResolved   = "resolved"
	InsightStatusDismissed  = "dismissed"
	InsightStatusOverridden = "overridden"

	InsightSeverityLow      = "low"
	InsightSeverityMedium   = "medium"
	InsightSeverityHigh     = "high"
	InsightSeverityCritical = "critical"

	InsightSourceDeterministic = "deterministic"
	InsightSourceAI            = "ai_assisted"
)

type Insight struct {
	ID              uuid.UUID       `json:"id"`
	TenantID        uuid.UUID       `json:"tenant_id"`
	InsightKey      string          `json:"insight_key"`
	InsightType     string          `json:"insight_type"`
	Category        string          `json:"category"`
	Severity        string          `json:"severity"`
	Status          string          `json:"status"`
	Title           string          `json:"title"`
	Summary         string          `json:"summary"`
	ConfidenceScore float64         `json:"confidence_score"`
	Score           float64         `json:"score"`
	Source          string          `json:"source"`
	ModelVersion    *string         `json:"model_version,omitempty"`
	EntityType      *string         `json:"entity_type,omitempty"`
	EntityID        *uuid.UUID      `json:"entity_id,omitempty"`
	EmployeeUserID  *uuid.UUID      `json:"employee_user_id,omitempty"`
	Reasons         json.RawMessage `json:"reasons,omitempty"`
	Recommendations json.RawMessage `json:"recommendations,omitempty"`
	Context         json.RawMessage `json:"context,omitempty"`
	Explainability  json.RawMessage `json:"explainability,omitempty"`
	DetectedAt      time.Time       `json:"detected_at"`
	DueAt           *time.Time      `json:"due_at,omitempty"`
	AssignedTo      *uuid.UUID      `json:"assigned_to,omitempty"`
	ReviewedBy      *uuid.UUID      `json:"reviewed_by,omitempty"`
	ReviewedAt      *time.Time      `json:"reviewed_at,omitempty"`
	ResolvedAt      *time.Time      `json:"resolved_at,omitempty"`
	ResolutionNote  *string         `json:"resolution_note,omitempty"`
	Inactive        bool            `json:"inactive"`
	CreatedAt       time.Time       `json:"created_at"`
	CreatedBy       *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt       time.Time       `json:"updated_at"`
	UpdatedBy       *uuid.UUID      `json:"updated_by,omitempty"`
}

type InsightEvent struct {
	ID         uuid.UUID       `json:"id"`
	TenantID   uuid.UUID       `json:"tenant_id"`
	InsightID  uuid.UUID       `json:"insight_id"`
	Action     string          `json:"action"`
	FromStatus *string         `json:"from_status,omitempty"`
	ToStatus   *string         `json:"to_status,omitempty"`
	Remarks    *string         `json:"remarks,omitempty"`
	Metadata   json.RawMessage `json:"metadata,omitempty"`
	Inactive   bool            `json:"inactive"`
	CreatedAt  time.Time       `json:"created_at"`
	CreatedBy  *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt  time.Time       `json:"updated_at"`
	UpdatedBy  *uuid.UUID      `json:"updated_by,omitempty"`
}

type InsightFilter struct {
	TenantID    uuid.UUID
	Status      *string
	Severity    *string
	Category    *string
	InsightType *string
	AssignedTo  *uuid.UUID
	Limit       int32
	Offset      int32
}

type InsightSummary struct {
	Total      int32            `json:"total"`
	Open       int32            `json:"open"`
	Reviewing  int32            `json:"reviewing"`
	Resolved   int32            `json:"resolved"`
	Dismissed  int32            `json:"dismissed"`
	Overridden int32            `json:"overridden"`
	Critical   int32            `json:"critical"`
	High       int32            `json:"high"`
	Medium     int32            `json:"medium"`
	Low        int32            `json:"low"`
	ByCategory map[string]int32 `json:"by_category"`
}

type InsightWorkspace struct {
	Items   []*Insight      `json:"items"`
	Summary InsightSummary  `json:"summary"`
	Events  []*InsightEvent `json:"events,omitempty"`
}

func ValidateInsightStatus(value string) (string, error) {
	status := strings.TrimSpace(value)
	switch status {
	case InsightStatusOpen, InsightStatusReviewing, InsightStatusResolved, InsightStatusDismissed, InsightStatusOverridden:
		return status, nil
	default:
		return "", ErrInvalidInsightStatus
	}
}

func ValidateInsightSeverity(value string) (string, error) {
	severity := strings.TrimSpace(value)
	switch severity {
	case InsightSeverityLow, InsightSeverityMedium, InsightSeverityHigh, InsightSeverityCritical:
		return severity, nil
	default:
		return "", ErrInvalidInsightSeverity
	}
}
