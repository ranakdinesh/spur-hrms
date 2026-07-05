package ports

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

type JobRequisitionRepo interface {
	CreateJobRequisition(ctx context.Context, item *domain.JobRequisition, actorID *uuid.UUID) (*domain.JobRequisition, error)
	ListJobRequisitions(ctx context.Context, filter domain.JobRequisitionFilter) ([]*domain.JobRequisition, error)
	CountJobRequisitions(ctx context.Context, filter domain.JobRequisitionFilter) (int64, error)
	GetJobRequisition(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.JobRequisition, error)
	UpdateJobRequisition(ctx context.Context, item *domain.JobRequisition, actorID *uuid.UUID) (*domain.JobRequisition, error)
	UpdateJobRequisitionStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, notes *string, actorID *uuid.UUID) (*domain.JobRequisition, error)
	DeleteJobRequisition(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	CreateJobRequisitionLog(ctx context.Context, item *domain.JobRequisitionLog, actorID *uuid.UUID) (*domain.JobRequisitionLog, error)
	ListJobRequisitionLogs(ctx context.Context, tenantID uuid.UUID, jobRequisitionID uuid.UUID) ([]*domain.JobRequisitionLog, error)
}

type JobRequisitionCommand struct {
	ID                  uuid.UUID  `json:"id,omitempty"`
	TenantID            uuid.UUID  `json:"tenant_id"`
	JobPositionID       uuid.UUID  `json:"job_position_id"`
	Code                *string    `json:"code,omitempty"`
	Title               string     `json:"title"`
	Level               *string    `json:"level,omitempty"`
	Category            *string    `json:"category,omitempty"`
	DepartmentID        *uuid.UUID `json:"department_id,omitempty"`
	EmploymentTypeID    *uuid.UUID `json:"employment_type_id,omitempty"`
	Description         *string    `json:"description,omitempty"`
	WorkMode            *string    `json:"work_mode,omitempty"`
	TotalOpenings       int32      `json:"total_openings"`
	ReasonForHire       *string    `json:"reason_for_hire,omitempty"`
	MinSalary           *float64   `json:"min_salary,omitempty"`
	MaxSalary           *float64   `json:"max_salary,omitempty"`
	Currency            string     `json:"currency"`
	TargetHireDate      *time.Time `json:"target_hire_date,omitempty"`
	ExpectedClosureDate *time.Time `json:"expected_closure_date,omitempty"`
	RequestedBy         uuid.UUID  `json:"requested_by"`
	RequestedDate       *time.Time `json:"requested_date,omitempty"`
	Priority            *string    `json:"priority,omitempty"`
	Status              string     `json:"status"`
	Notes               *string    `json:"notes,omitempty"`
	ActorID             *uuid.UUID `json:"-"`
}

type JobRequisitionActionCommand struct {
	TenantID uuid.UUID  `json:"tenant_id"`
	ID       uuid.UUID  `json:"id"`
	Remarks  *string    `json:"remarks,omitempty"`
	ActorID  *uuid.UUID `json:"-"`
}
