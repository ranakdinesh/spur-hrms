package ports

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

type InsightRepo interface {
	UpsertInsight(ctx context.Context, item *domain.Insight, actorID *uuid.UUID) (*domain.Insight, error)
	ListInsights(ctx context.Context, filter domain.InsightFilter) ([]*domain.Insight, error)
	GetInsight(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.Insight, error)
	UpdateInsightStatus(ctx context.Context, cmd InsightStatusCommand) (*domain.Insight, error)
	CreateInsightEvent(ctx context.Context, event *domain.InsightEvent, actorID *uuid.UUID) (*domain.InsightEvent, error)
	ListInsightEvents(ctx context.Context, tenantID uuid.UUID, insightID uuid.UUID) ([]*domain.InsightEvent, error)
}

type InsightScoringPort interface {
	ScoreInsights(ctx context.Context, tenantID uuid.UUID, candidates []*domain.Insight) ([]*domain.Insight, error)
}

type InsightStatusCommand struct {
	TenantID       uuid.UUID  `json:"tenant_id"`
	ID             uuid.UUID  `json:"id"`
	Status         string     `json:"status"`
	AssignedTo     *uuid.UUID `json:"assigned_to,omitempty"`
	Remarks        *string    `json:"remarks,omitempty"`
	ResolutionNote *string    `json:"resolution_note,omitempty"`
	ActorID        *uuid.UUID `json:"-"`
}

type InsightEventCommand struct {
	TenantID   uuid.UUID       `json:"tenant_id"`
	InsightID  uuid.UUID       `json:"insight_id"`
	Action     string          `json:"action"`
	FromStatus *string         `json:"from_status,omitempty"`
	ToStatus   *string         `json:"to_status,omitempty"`
	Remarks    *string         `json:"remarks,omitempty"`
	Metadata   json.RawMessage `json:"metadata,omitempty"`
	ActorID    *uuid.UUID      `json:"-"`
}
