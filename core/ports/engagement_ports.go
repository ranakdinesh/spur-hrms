package ports

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

type EngagementRepo interface {
	CreateEngagement(ctx context.Context, item *domain.Engagement, actorID *uuid.UUID) (*domain.Engagement, error)
	UpdateEngagement(ctx context.Context, item *domain.Engagement, actorID *uuid.UUID) (*domain.Engagement, error)
	UpdateEngagementStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, terminationReason *string, terminatedAt *time.Time, actorID *uuid.UUID) (*domain.Engagement, error)
	GetEngagement(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.Engagement, error)
	ListEngagements(ctx context.Context, filter domain.EngagementFilter) ([]*domain.EngagementListItem, error)
	DeleteEngagement(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
}

type EngagementCommand struct {
	ID                 uuid.UUID       `json:"id,omitempty"`
	TenantID           uuid.UUID       `json:"tenant_id"`
	WorkerProfileID    uuid.UUID       `json:"worker_profile_id"`
	EngagementCode     *string         `json:"engagement_code,omitempty"`
	Title              string          `json:"title"`
	Description        *string         `json:"description,omitempty"`
	EngagementType     string          `json:"engagement_type"`
	Status             string          `json:"status"`
	StartDate          string          `json:"start_date"`
	EndDate            string          `json:"end_date,omitempty"`
	HoursBudget        *float64        `json:"hours_budget,omitempty"`
	RateAmount         *float64        `json:"rate_amount,omitempty"`
	CurrencyCode       string          `json:"currency_code"`
	RateUnit           string          `json:"rate_unit"`
	BranchID           *uuid.UUID      `json:"branch_id,omitempty"`
	DepartmentID       *uuid.UUID      `json:"department_id,omitempty"`
	ReportingManagerID *uuid.UUID      `json:"reporting_manager_id,omitempty"`
	ProjectLabel       *string         `json:"project_label,omitempty"`
	ProjectCode        *string         `json:"project_code,omitempty"`
	CostCenter         *string         `json:"cost_center,omitempty"`
	RenewalDueDate     string          `json:"renewal_due_date,omitempty"`
	RenewalStatus      string          `json:"renewal_status"`
	TerminationReason  *string         `json:"termination_reason,omitempty"`
	TerminatedAt       string          `json:"terminated_at,omitempty"`
	Notes              *string         `json:"notes,omitempty"`
	Metadata           json.RawMessage `json:"metadata,omitempty"`
	ActorID            *uuid.UUID      `json:"-"`
}

type EngagementStatusCommand struct {
	TenantID          uuid.UUID  `json:"tenant_id"`
	EngagementID      uuid.UUID  `json:"engagement_id"`
	Status            string     `json:"status"`
	TerminationReason *string    `json:"termination_reason,omitempty"`
	TerminatedAt      string     `json:"terminated_at,omitempty"`
	ActorID           *uuid.UUID `json:"-"`
}
