package ports

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

type SkillGapRepo interface {
	CreateProjectSkillRequirement(ctx context.Context, item *domain.ProjectSkillRequirement, actorID *uuid.UUID) (*domain.ProjectSkillRequirement, error)
	UpdateProjectSkillRequirement(ctx context.Context, item *domain.ProjectSkillRequirement, actorID *uuid.UUID) (*domain.ProjectSkillRequirement, error)
	GetProjectSkillRequirement(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.ProjectSkillRequirement, error)
	ListProjectSkillRequirements(ctx context.Context, filter domain.ProjectSkillRequirementFilter) ([]*domain.ProjectSkillRequirement, error)
	DeleteProjectSkillRequirement(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	ListProjectSkillGapRows(ctx context.Context, filter domain.ProjectSkillRequirementFilter) ([]*domain.ProjectSkillGapRow, error)
	ListSkillGapSummary(ctx context.Context, tenantID uuid.UUID, projectID *uuid.UUID) ([]*domain.SkillGapSummaryRow, error)
	ListSinglePersonSkillDependencies(ctx context.Context, tenantID uuid.UUID, projectID *uuid.UUID) ([]*domain.SinglePersonSkillDependency, error)
}

type ProjectSkillRequirementCommand struct {
	ID                  uuid.UUID       `json:"id,omitempty"`
	TenantID            uuid.UUID       `json:"tenant_id"`
	ProjectID           *uuid.UUID      `json:"project_id,omitempty"`
	EngagementID        *uuid.UUID      `json:"engagement_id,omitempty"`
	SkillID             uuid.UUID       `json:"skill_id"`
	RequiredProficiency string          `json:"required_proficiency"`
	MinYearsExperience  *float64        `json:"min_years_experience,omitempty"`
	RequiredCount       int32           `json:"required_count"`
	Importance          string          `json:"importance"`
	RequirementSource   string          `json:"requirement_source"`
	Notes               *string         `json:"notes,omitempty"`
	Metadata            json.RawMessage `json:"metadata,omitempty"`
	ActorID             *uuid.UUID      `json:"-"`
}
