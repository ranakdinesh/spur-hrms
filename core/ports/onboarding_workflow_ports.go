package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

type OnboardingWorkflowRepo interface {
	CreateOnboardingWorkflow(ctx context.Context, item *domain.OnboardingWorkflow, actorID *uuid.UUID) (*domain.OnboardingWorkflow, error)
	ListOnboardingWorkflows(ctx context.Context, tenantID uuid.UUID) ([]*domain.OnboardingWorkflow, error)
	GetOnboardingWorkflow(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.OnboardingWorkflow, error)
	UpdateOnboardingWorkflow(ctx context.Context, item *domain.OnboardingWorkflow, actorID *uuid.UUID) (*domain.OnboardingWorkflow, error)
	DeleteOnboardingWorkflow(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	CreateOnboardingTask(ctx context.Context, item *domain.OnboardingTask, actorID *uuid.UUID) (*domain.OnboardingTask, error)
	ListOnboardingTasks(ctx context.Context, tenantID uuid.UUID, workflowID uuid.UUID) ([]*domain.OnboardingTask, error)
	GetOnboardingTask(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.OnboardingTask, error)
	UpdateOnboardingTask(ctx context.Context, item *domain.OnboardingTask, actorID *uuid.UUID) (*domain.OnboardingTask, error)
	DeleteOnboardingTask(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	CreateOnboardingWorkflowAssignment(ctx context.Context, item *domain.OnboardingWorkflowAssignment, actorID *uuid.UUID) (*domain.OnboardingWorkflowAssignment, error)
	ListOnboardingWorkflowAssignments(ctx context.Context, tenantID uuid.UUID) ([]*domain.OnboardingWorkflowAssignment, error)
	GetOnboardingWorkflowAssignment(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.OnboardingWorkflowAssignment, error)
	UpdateOnboardingWorkflowAssignment(ctx context.Context, item *domain.OnboardingWorkflowAssignment, actorID *uuid.UUID) (*domain.OnboardingWorkflowAssignment, error)
	DeleteOnboardingWorkflowAssignment(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
}

type OnboardingWorkflowCommand struct {
	ID          uuid.UUID  `json:"id,omitempty"`
	TenantID    uuid.UUID  `json:"tenant_id"`
	Name        string     `json:"name"`
	Description *string    `json:"description,omitempty"`
	IsDefault   bool       `json:"is_default"`
	IsActive    bool       `json:"is_active"`
	ActorID     *uuid.UUID `json:"-"`
}

type OnboardingTaskCommand struct {
	ID          uuid.UUID  `json:"id,omitempty"`
	TenantID    uuid.UUID  `json:"tenant_id"`
	WorkflowID  uuid.UUID  `json:"workflow_id"`
	Title       string     `json:"title"`
	Description *string    `json:"description,omitempty"`
	DueDays     int32      `json:"due_days"`
	IsRequired  bool       `json:"is_required"`
	SortOrder   int32      `json:"sort_order"`
	ActorID     *uuid.UUID `json:"-"`
}

type OnboardingAssignmentCommand struct {
	ID               uuid.UUID  `json:"id,omitempty"`
	TenantID         uuid.UUID  `json:"tenant_id"`
	WorkflowID       uuid.UUID  `json:"workflow_id"`
	Name             string     `json:"name"`
	JobPostingID     *uuid.UUID `json:"job_posting_id,omitempty"`
	JobPositionID    *uuid.UUID `json:"job_position_id,omitempty"`
	DepartmentID     *uuid.UUID `json:"department_id,omitempty"`
	EmploymentTypeID *uuid.UUID `json:"employment_type_id,omitempty"`
	Priority         int32      `json:"priority"`
	ActorID          *uuid.UUID `json:"-"`
}
