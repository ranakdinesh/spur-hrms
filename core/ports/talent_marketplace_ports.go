package ports

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

type TalentMarketplaceRepo interface {
	CreateTalentMarketplaceOpportunity(ctx context.Context, item *domain.TalentMarketplaceOpportunity, actorID *uuid.UUID) (*domain.TalentMarketplaceOpportunity, error)
	UpdateTalentMarketplaceOpportunity(ctx context.Context, item *domain.TalentMarketplaceOpportunity, actorID *uuid.UUID) (*domain.TalentMarketplaceOpportunity, error)
	UpdateTalentMarketplaceOpportunityFallback(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, actorID *uuid.UUID) (*domain.TalentMarketplaceOpportunity, error)
	GetTalentMarketplaceOpportunity(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.TalentMarketplaceOpportunity, error)
	ListTalentMarketplaceOpportunities(ctx context.Context, filter domain.TalentMarketplaceOpportunityFilter) ([]*domain.TalentMarketplaceOpportunity, error)
	DeleteTalentMarketplaceOpportunity(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	CreateTalentMarketplaceApplication(ctx context.Context, item *domain.TalentMarketplaceApplication, actorID *uuid.UUID) (*domain.TalentMarketplaceApplication, error)
	UpdateTalentMarketplaceApplicationStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, workerNote *string, managerNote *string, actorID *uuid.UUID) (*domain.TalentMarketplaceApplication, error)
	GetTalentMarketplaceApplication(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.TalentMarketplaceApplication, error)
	ListTalentMarketplaceApplications(ctx context.Context, filter domain.TalentMarketplaceApplicationFilter) ([]*domain.TalentMarketplaceApplication, error)
	ListTalentMarketplaceRecommendations(ctx context.Context, tenantID uuid.UUID, opportunityID uuid.UUID) ([]*domain.TalentMarketplaceRecommendation, error)
	CreateTalentMarketplaceEvent(ctx context.Context, event *domain.TalentMarketplaceEvent) (*domain.TalentMarketplaceEvent, error)
	ListTalentMarketplaceEvents(ctx context.Context, filter domain.TalentMarketplaceEventFilter) ([]*domain.TalentMarketplaceEvent, error)
}

type TalentMarketplaceOpportunityCommand struct {
	ID                       uuid.UUID       `json:"id,omitempty"`
	TenantID                 uuid.UUID       `json:"tenant_id"`
	ProjectID                *uuid.UUID      `json:"project_id,omitempty"`
	EngagementID             *uuid.UUID      `json:"engagement_id,omitempty"`
	SourceRequirementID      *uuid.UUID      `json:"source_requirement_id,omitempty"`
	JobPostingID             *uuid.UUID      `json:"job_posting_id,omitempty"`
	Title                    string          `json:"title"`
	Description              *string         `json:"description,omitempty"`
	OpportunityType          string          `json:"opportunity_type"`
	Status                   string          `json:"status"`
	Visibility               string          `json:"visibility"`
	Priority                 string          `json:"priority"`
	Seats                    int32           `json:"seats"`
	LocationMode             string          `json:"location_mode"`
	MinAllocationPercent     *int32          `json:"min_allocation_percent,omitempty"`
	DurationLabel            *string         `json:"duration_label,omitempty"`
	StartDate                string          `json:"start_date"`
	DueDate                  string          `json:"due_date"`
	CandidateFallbackEnabled bool            `json:"candidate_fallback_enabled"`
	CandidateFallbackStatus  string          `json:"candidate_fallback_status"`
	Metadata                 json.RawMessage `json:"metadata,omitempty"`
	ActorID                  *uuid.UUID      `json:"-"`
}

type TalentMarketplaceApplicationCommand struct {
	ID              uuid.UUID       `json:"id,omitempty"`
	TenantID        uuid.UUID       `json:"tenant_id"`
	OpportunityID   uuid.UUID       `json:"opportunity_id"`
	WorkerProfileID uuid.UUID       `json:"worker_profile_id"`
	Status          string          `json:"status"`
	MatchScore      *float64        `json:"match_score,omitempty"`
	MatchReasons    json.RawMessage `json:"match_reasons,omitempty"`
	WorkerNote      *string         `json:"worker_note,omitempty"`
	ManagerNote     *string         `json:"manager_note,omitempty"`
	ActorID         *uuid.UUID      `json:"-"`
}

type TalentMarketplaceApplicationStatusCommand struct {
	TenantID    uuid.UUID  `json:"tenant_id"`
	ID          uuid.UUID  `json:"id"`
	Status      string     `json:"status"`
	WorkerNote  *string    `json:"worker_note,omitempty"`
	ManagerNote *string    `json:"manager_note,omitempty"`
	ActorID     *uuid.UUID `json:"-"`
}

type TalentMarketplaceFallbackCommand struct {
	TenantID uuid.UUID  `json:"tenant_id"`
	ID       uuid.UUID  `json:"id"`
	Status   string     `json:"status"`
	Notes    *string    `json:"notes,omitempty"`
	ActorID  *uuid.UUID `json:"-"`
}
