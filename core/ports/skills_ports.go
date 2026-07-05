package ports

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

type SkillsRepo interface {
	CreateSkillCategory(ctx context.Context, item *domain.SkillCategory, actorID *uuid.UUID) (*domain.SkillCategory, error)
	UpdateSkillCategory(ctx context.Context, item *domain.SkillCategory, actorID *uuid.UUID) (*domain.SkillCategory, error)
	GetSkillCategory(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.SkillCategory, error)
	ListSkillCategories(ctx context.Context, filter domain.SkillCategoryFilter) ([]*domain.SkillCategory, error)
	DeleteSkillCategory(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	CreateSkill(ctx context.Context, item *domain.Skill, actorID *uuid.UUID) (*domain.Skill, error)
	UpdateSkill(ctx context.Context, item *domain.Skill, actorID *uuid.UUID) (*domain.Skill, error)
	GetSkill(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.Skill, error)
	GetSkillByCode(ctx context.Context, tenantID uuid.UUID, code string) (*domain.Skill, error)
	ListSkills(ctx context.Context, filter domain.SkillFilter) ([]*domain.Skill, error)
	DeleteSkill(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	CreateWorkerSkill(ctx context.Context, item *domain.WorkerSkill, actorID *uuid.UUID) (*domain.WorkerSkill, error)
	UpdateWorkerSkill(ctx context.Context, item *domain.WorkerSkill, actorID *uuid.UUID) (*domain.WorkerSkill, error)
	GetWorkerSkill(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.WorkerSkill, error)
	ListWorkerSkills(ctx context.Context, filter domain.WorkerSkillFilter) ([]*domain.WorkerSkill, error)
	UpdateWorkerSkillVerification(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, notes *string, actorID *uuid.UUID) (*domain.WorkerSkill, error)
	DeleteWorkerSkill(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	CreateWorkerSkillAssessment(ctx context.Context, item *domain.WorkerSkillAssessment, actorID *uuid.UUID) (*domain.WorkerSkillAssessment, error)
	ListWorkerSkillAssessments(ctx context.Context, tenantID uuid.UUID, workerSkillID *uuid.UUID) ([]*domain.WorkerSkillAssessment, error)
	GetSkillsSummary(ctx context.Context, tenantID uuid.UUID) ([]*domain.SkillsSummaryRow, error)
}

type SkillCategoryCommand struct {
	TenantID    uuid.UUID       `json:"tenant_id"`
	ID          uuid.UUID       `json:"id,omitempty"`
	ParentID    *uuid.UUID      `json:"parent_id,omitempty"`
	Code        string          `json:"code"`
	Name        string          `json:"name"`
	Description *string         `json:"description,omitempty"`
	SortOrder   int32           `json:"sort_order"`
	Metadata    json.RawMessage `json:"metadata,omitempty"`
	ActorID     *uuid.UUID      `json:"actor_id,omitempty"`
}

type SkillCommand struct {
	TenantID            uuid.UUID       `json:"tenant_id"`
	ID                  uuid.UUID       `json:"id,omitempty"`
	CategoryID          *uuid.UUID      `json:"category_id,omitempty"`
	Code                string          `json:"code"`
	Name                string          `json:"name"`
	Description         *string         `json:"description,omitempty"`
	SkillType           string          `json:"skill_type"`
	CertificateRequired bool            `json:"certificate_required"`
	AssessmentRequired  bool            `json:"assessment_required"`
	IsActive            bool            `json:"is_active"`
	Metadata            json.RawMessage `json:"metadata,omitempty"`
	ActorID             *uuid.UUID      `json:"actor_id,omitempty"`
}

type WorkerSkillCommand struct {
	TenantID             uuid.UUID       `json:"tenant_id"`
	ID                   uuid.UUID       `json:"id,omitempty"`
	WorkerProfileID      uuid.UUID       `json:"worker_profile_id"`
	SkillID              uuid.UUID       `json:"skill_id"`
	Proficiency          string          `json:"proficiency"`
	YearsExperience      *float64        `json:"years_experience,omitempty"`
	LastUsedOn           *time.Time      `json:"last_used_on,omitempty"`
	VerificationStatus   string          `json:"verification_status"`
	CertificateURL       *string         `json:"certificate_url,omitempty"`
	CertificateExpiresOn *time.Time      `json:"certificate_expires_on,omitempty"`
	AssessmentScore      *float64        `json:"assessment_score,omitempty"`
	AssessedOn           *time.Time      `json:"assessed_on,omitempty"`
	Notes                *string         `json:"notes,omitempty"`
	Metadata             json.RawMessage `json:"metadata,omitempty"`
	ActorID              *uuid.UUID      `json:"actor_id,omitempty"`
}

type WorkerSkillVerificationCommand struct {
	TenantID uuid.UUID  `json:"tenant_id"`
	ID       uuid.UUID  `json:"id"`
	Status   string     `json:"status"`
	Notes    *string    `json:"notes,omitempty"`
	ActorID  *uuid.UUID `json:"actor_id,omitempty"`
}

type WorkerSkillAssessmentCommand struct {
	TenantID       uuid.UUID       `json:"tenant_id"`
	WorkerSkillID  uuid.UUID       `json:"worker_skill_id"`
	AssessmentType string          `json:"assessment_type"`
	ResultStatus   string          `json:"result_status"`
	Score          *float64        `json:"score,omitempty"`
	MaxScore       *float64        `json:"max_score,omitempty"`
	AssessedBy     *uuid.UUID      `json:"assessed_by,omitempty"`
	AssessedOn     *time.Time      `json:"assessed_on,omitempty"`
	EvidenceURL    *string         `json:"evidence_url,omitempty"`
	Notes          *string         `json:"notes,omitempty"`
	Metadata       json.RawMessage `json:"metadata,omitempty"`
	ActorID        *uuid.UUID      `json:"actor_id,omitempty"`
}
