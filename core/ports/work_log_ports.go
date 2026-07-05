package ports

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

type WorkLogRepo interface {
	CreateWorkLog(ctx context.Context, item *domain.WorkLog, actorID *uuid.UUID) (*domain.WorkLog, error)
	UpdateWorkLog(ctx context.Context, item *domain.WorkLog, actorID *uuid.UUID) (*domain.WorkLog, error)
	UpdateWorkLogStatus(ctx context.Context, item *domain.WorkLog, actorID *uuid.UUID) (*domain.WorkLog, error)
	GetWorkLog(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.WorkLog, error)
	ListWorkLogs(ctx context.Context, filter domain.WorkLogFilter) ([]*domain.WorkLogListItem, error)
	ListWorkLogRollups(ctx context.Context, filter domain.WorkLogFilter) ([]*domain.WorkLogRollup, error)
	GetWorkLogBudgetUsage(ctx context.Context, tenantID uuid.UUID, engagementID uuid.UUID, excludeWorkLogID *uuid.UUID) (*domain.WorkLogBudgetUsage, error)
	DeleteWorkLog(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
}

type WorkLogCommand struct {
	ID                   uuid.UUID       `json:"id,omitempty"`
	TenantID             uuid.UUID       `json:"tenant_id"`
	EngagementID         uuid.UUID       `json:"engagement_id"`
	WorkerProfileID      uuid.UUID       `json:"worker_profile_id,omitempty"`
	LogDate              string          `json:"log_date"`
	HoursWorked          float64         `json:"hours_worked"`
	BillableHours        *float64        `json:"billable_hours,omitempty"`
	WorkSummary          *string         `json:"work_summary,omitempty"`
	DeliverableReference *string         `json:"deliverable_reference,omitempty"`
	Status               string          `json:"status"`
	Metadata             json.RawMessage `json:"metadata,omitempty"`
	ActorID              *uuid.UUID      `json:"-"`
}

type WorkLogStatusCommand struct {
	TenantID      uuid.UUID  `json:"tenant_id"`
	WorkLogID     uuid.UUID  `json:"work_log_id"`
	Status        string     `json:"status"`
	ReviewComment *string    `json:"review_comment,omitempty"`
	ActorID       *uuid.UUID `json:"-"`
}
