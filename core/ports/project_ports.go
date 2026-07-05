package ports

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

type ProjectRepo interface {
	CreateProject(ctx context.Context, item *domain.Project, actorID *uuid.UUID) (*domain.Project, error)
	UpdateProject(ctx context.Context, item *domain.Project, actorID *uuid.UUID) (*domain.Project, error)
	UpdateProjectStatus(ctx context.Context, item *domain.Project, actorID *uuid.UUID) (*domain.Project, error)
	GetProject(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.Project, error)
	ListProjects(ctx context.Context, filter domain.ProjectFilter) ([]*domain.ProjectListItem, error)
	DeleteProject(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	CreateProjectMilestone(ctx context.Context, item *domain.ProjectMilestone, actorID *uuid.UUID) (*domain.ProjectMilestone, error)
	UpdateProjectMilestone(ctx context.Context, item *domain.ProjectMilestone, actorID *uuid.UUID) (*domain.ProjectMilestone, error)
	UpdateProjectMilestoneStatus(ctx context.Context, item *domain.ProjectMilestone, actorID *uuid.UUID) (*domain.ProjectMilestone, error)
	GetProjectMilestone(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.ProjectMilestone, error)
	ListProjectMilestones(ctx context.Context, filter domain.ProjectMilestoneFilter) ([]*domain.ProjectMilestoneListItem, error)
	DeleteProjectMilestone(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	CreateProjectMilestoneEvent(ctx context.Context, event *domain.ProjectMilestoneEvent) (*domain.ProjectMilestoneEvent, error)
	ListProjectMilestoneEvents(ctx context.Context, tenantID uuid.UUID, milestoneID uuid.UUID) ([]*domain.ProjectMilestoneEvent, error)
}

type ProjectCommand struct {
	ID               uuid.UUID       `json:"id,omitempty"`
	TenantID         uuid.UUID       `json:"tenant_id"`
	ProjectCode      *string         `json:"project_code,omitempty"`
	Name             string          `json:"name"`
	Description      *string         `json:"description,omitempty"`
	Status           string          `json:"status"`
	DepartmentID     *uuid.UUID      `json:"department_id,omitempty"`
	BranchID         *uuid.UUID      `json:"branch_id,omitempty"`
	ProjectManagerID *uuid.UUID      `json:"project_manager_id,omitempty"`
	StartDate        string          `json:"start_date"`
	DueDate          string          `json:"due_date"`
	BudgetAmount     *float64        `json:"budget_amount,omitempty"`
	CurrencyCode     string          `json:"currency_code"`
	BillingType      string          `json:"billing_type"`
	ClientLabel      *string         `json:"client_label,omitempty"`
	Priority         string          `json:"priority"`
	Notes            *string         `json:"notes,omitempty"`
	Metadata         json.RawMessage `json:"metadata,omitempty"`
	ActorID          *uuid.UUID      `json:"-"`
}

type ProjectStatusCommand struct {
	TenantID  uuid.UUID  `json:"tenant_id"`
	ProjectID uuid.UUID  `json:"project_id"`
	Status    string     `json:"status"`
	ActorID   *uuid.UUID `json:"-"`
}

type ProjectMilestoneCommand struct {
	ID                 uuid.UUID       `json:"id,omitempty"`
	TenantID           uuid.UUID       `json:"tenant_id"`
	ProjectID          uuid.UUID       `json:"project_id"`
	EngagementID       *uuid.UUID      `json:"engagement_id,omitempty"`
	MilestoneCode      *string         `json:"milestone_code,omitempty"`
	Title              string          `json:"title"`
	Description        *string         `json:"description,omitempty"`
	AcceptanceCriteria *string         `json:"acceptance_criteria,omitempty"`
	DueDate            string          `json:"due_date"`
	Status             string          `json:"status"`
	Amount             *float64        `json:"amount,omitempty"`
	CurrencyCode       string          `json:"currency_code"`
	PaymentTrigger     json.RawMessage `json:"payment_trigger,omitempty"`
	Notes              *string         `json:"notes,omitempty"`
	Metadata           json.RawMessage `json:"metadata,omitempty"`
	ActorID            *uuid.UUID      `json:"-"`
}

type ProjectMilestoneStatusCommand struct {
	TenantID      uuid.UUID  `json:"tenant_id"`
	MilestoneID   uuid.UUID  `json:"milestone_id"`
	Status        string     `json:"status"`
	ReviewComment *string    `json:"review_comment,omitempty"`
	ActorID       *uuid.UUID `json:"-"`
}
