package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

type JobPositionRepo interface {
	CreateJobPosition(ctx context.Context, item *domain.JobPosition, actorID *uuid.UUID) (*domain.JobPosition, error)
	ListJobPositions(ctx context.Context, filter domain.JobPositionFilter) ([]*domain.JobPosition, error)
	CountJobPositions(ctx context.Context, filter domain.JobPositionFilter) (int64, error)
	GetJobPosition(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.JobPosition, error)
	UpdateJobPosition(ctx context.Context, item *domain.JobPosition, actorID *uuid.UUID) (*domain.JobPosition, error)
	DeleteJobPosition(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	CreateJobPositionLocation(ctx context.Context, item *domain.JobPositionLocation, actorID *uuid.UUID) (*domain.JobPositionLocation, error)
	ListJobPositionLocations(ctx context.Context, tenantID uuid.UUID, jobPositionID uuid.UUID) ([]*domain.JobPositionLocation, error)
	GetJobPositionLocation(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.JobPositionLocation, error)
	UpdateJobPositionLocation(ctx context.Context, item *domain.JobPositionLocation, actorID *uuid.UUID) (*domain.JobPositionLocation, error)
	DeleteJobPositionLocation(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
}

type JobPositionCommand struct {
	ID               uuid.UUID  `json:"id,omitempty"`
	TenantID         uuid.UUID  `json:"tenant_id"`
	Code             *string    `json:"code,omitempty"`
	Title            string     `json:"title"`
	Level            *string    `json:"level,omitempty"`
	Category         *string    `json:"category,omitempty"`
	Description      *string    `json:"description,omitempty"`
	DepartmentID     *uuid.UUID `json:"department_id,omitempty"`
	EmploymentTypeID *uuid.UUID `json:"employment_type_id,omitempty"`
	WorkMode         *string    `json:"work_mode,omitempty"`
	TotalPosition    int32      `json:"total_position"`
	BudgetedCost     *float64   `json:"budgeted_cost,omitempty"`
	ActorID          *uuid.UUID `json:"-"`
}

type JobPositionLocationCommand struct {
	ID            uuid.UUID  `json:"id,omitempty"`
	TenantID      uuid.UUID  `json:"tenant_id"`
	JobPositionID uuid.UUID  `json:"job_position_id"`
	Location      *string    `json:"location,omitempty"`
	City          *string    `json:"city,omitempty"`
	State         *string    `json:"state,omitempty"`
	Country       *string    `json:"country,omitempty"`
	IsRemote      bool       `json:"is_remote"`
	ActorID       *uuid.UUID `json:"-"`
}
