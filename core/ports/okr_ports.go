package ports

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

type OKRRepo interface {
	CreateOKRCycle(ctx context.Context, item *domain.OKRCycle, actorID *uuid.UUID) (*domain.OKRCycle, error)
	UpdateOKRCycle(ctx context.Context, item *domain.OKRCycle, actorID *uuid.UUID) (*domain.OKRCycle, error)
	UpdateOKRCycleStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, actorID *uuid.UUID) (*domain.OKRCycle, error)
	GetOKRCycle(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.OKRCycle, error)
	ListOKRCycles(ctx context.Context, filter domain.OKRCycleFilter) ([]*domain.OKRCycle, error)
	DeleteOKRCycle(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	CreateObjective(ctx context.Context, item *domain.Objective, actorID *uuid.UUID) (*domain.Objective, error)
	UpdateObjective(ctx context.Context, item *domain.Objective, actorID *uuid.UUID) (*domain.Objective, error)
	UpdateObjectiveStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, actorID *uuid.UUID) (*domain.Objective, error)
	RefreshObjectiveProgress(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) (*domain.Objective, error)
	GetObjective(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.Objective, error)
	ListObjectives(ctx context.Context, filter domain.ObjectiveFilter) ([]*domain.Objective, error)
	DeleteObjective(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	CreateKeyResult(ctx context.Context, item *domain.KeyResult, actorID *uuid.UUID) (*domain.KeyResult, error)
	UpdateKeyResult(ctx context.Context, item *domain.KeyResult, actorID *uuid.UUID) (*domain.KeyResult, error)
	UpdateKeyResultProgress(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, value float64, progress float64, confidence string, status string, actorID *uuid.UUID) (*domain.KeyResult, error)
	GetKeyResult(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.KeyResult, error)
	ListKeyResults(ctx context.Context, filter domain.KeyResultFilter) ([]*domain.KeyResult, error)
	DeleteKeyResult(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	CreateKeyResultCheckIn(ctx context.Context, item *domain.KeyResultCheckIn, actorID *uuid.UUID) (*domain.KeyResultCheckIn, error)
	ListKeyResultCheckIns(ctx context.Context, filter domain.KeyResultCheckInFilter) ([]*domain.KeyResultCheckIn, error)
	GetOKRSummary(ctx context.Context, tenantID uuid.UUID, cycleID *uuid.UUID) ([]*domain.OKRSummaryRow, error)
}

type OKRCycleCommand struct {
	ID            uuid.UUID       `json:"id,omitempty"`
	TenantID      uuid.UUID       `json:"tenant_id"`
	Name          string          `json:"name"`
	CycleCode     string          `json:"cycle_code"`
	Description   *string         `json:"description,omitempty"`
	StartDate     string          `json:"start_date"`
	EndDate       string          `json:"end_date"`
	Status        string          `json:"status"`
	ReviewCadence string          `json:"review_cadence"`
	Metadata      json.RawMessage `json:"metadata,omitempty"`
	ActorID       *uuid.UUID      `json:"-"`
}

type OKRStatusCommand struct {
	TenantID uuid.UUID  `json:"tenant_id"`
	ID       uuid.UUID  `json:"id"`
	Status   string     `json:"status"`
	ActorID  *uuid.UUID `json:"-"`
}

type ObjectiveCommand struct {
	ID                   uuid.UUID       `json:"id,omitempty"`
	TenantID             uuid.UUID       `json:"tenant_id"`
	CycleID              uuid.UUID       `json:"cycle_id"`
	ParentObjectiveID    *uuid.UUID      `json:"parent_objective_id,omitempty"`
	OwnerType            string          `json:"owner_type"`
	OwnerWorkerProfileID *uuid.UUID      `json:"owner_worker_profile_id,omitempty"`
	OwnerDepartmentID    *uuid.UUID      `json:"owner_department_id,omitempty"`
	OwnerProjectID       *uuid.UUID      `json:"owner_project_id,omitempty"`
	Title                string          `json:"title"`
	Description          *string         `json:"description,omitempty"`
	Status               string          `json:"status"`
	Priority             string          `json:"priority"`
	ProgressPercent      *float64        `json:"progress_percent,omitempty"`
	Weight               *float64        `json:"weight,omitempty"`
	StartDate            string          `json:"start_date"`
	DueDate              string          `json:"due_date"`
	Metadata             json.RawMessage `json:"metadata,omitempty"`
	ActorID              *uuid.UUID      `json:"-"`
}

type KeyResultCommand struct {
	ID              uuid.UUID       `json:"id,omitempty"`
	TenantID        uuid.UUID       `json:"tenant_id"`
	ObjectiveID     uuid.UUID       `json:"objective_id"`
	Title           string          `json:"title"`
	Description     *string         `json:"description,omitempty"`
	MetricType      string          `json:"metric_type"`
	StartValue      float64         `json:"start_value"`
	TargetValue     float64         `json:"target_value"`
	CurrentValue    float64         `json:"current_value"`
	ProgressPercent *float64        `json:"progress_percent,omitempty"`
	Confidence      string          `json:"confidence"`
	Status          string          `json:"status"`
	Weight          *float64        `json:"weight,omitempty"`
	UnitLabel       *string         `json:"unit_label,omitempty"`
	DueDate         string          `json:"due_date"`
	Metadata        json.RawMessage `json:"metadata,omitempty"`
	ActorID         *uuid.UUID      `json:"-"`
}

type KeyResultCheckInCommand struct {
	TenantID        uuid.UUID       `json:"tenant_id"`
	KeyResultID     uuid.UUID       `json:"key_result_id"`
	CheckInDate     string          `json:"checkin_date"`
	Value           float64         `json:"value"`
	ProgressPercent *float64        `json:"progress_percent,omitempty"`
	Confidence      string          `json:"confidence"`
	Status          string          `json:"status"`
	Note            *string         `json:"note,omitempty"`
	Metadata        json.RawMessage `json:"metadata,omitempty"`
	ActorID         *uuid.UUID      `json:"-"`
}
