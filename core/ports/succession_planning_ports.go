package ports

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

type SuccessionPlanningRepo interface {
	CreateSuccessionReviewCycle(ctx context.Context, item *domain.SuccessionReviewCycle, actorID *uuid.UUID) (*domain.SuccessionReviewCycle, error)
	UpdateSuccessionReviewCycle(ctx context.Context, item *domain.SuccessionReviewCycle, actorID *uuid.UUID) (*domain.SuccessionReviewCycle, error)
	UpdateSuccessionReviewCycleStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, actorID *uuid.UUID) (*domain.SuccessionReviewCycle, error)
	GetSuccessionReviewCycle(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.SuccessionReviewCycle, error)
	ListSuccessionReviewCycles(ctx context.Context, filter domain.SuccessionFilter) ([]*domain.SuccessionReviewCycle, error)
	CreateSuccessionCriticalRole(ctx context.Context, item *domain.SuccessionCriticalRole, actorID *uuid.UUID) (*domain.SuccessionCriticalRole, error)
	UpdateSuccessionCriticalRole(ctx context.Context, item *domain.SuccessionCriticalRole, actorID *uuid.UUID) (*domain.SuccessionCriticalRole, error)
	UpdateSuccessionCriticalRoleStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, actorID *uuid.UUID) (*domain.SuccessionCriticalRole, error)
	GetSuccessionCriticalRole(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.SuccessionCriticalRole, error)
	ListSuccessionCriticalRoles(ctx context.Context, filter domain.SuccessionFilter) ([]*domain.SuccessionCriticalRole, error)
	DeleteSuccessionCriticalRole(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	CreateSuccessionSuccessorNomination(ctx context.Context, item *domain.SuccessionSuccessorNomination, actorID *uuid.UUID) (*domain.SuccessionSuccessorNomination, error)
	UpdateSuccessionSuccessorNomination(ctx context.Context, item *domain.SuccessionSuccessorNomination, actorID *uuid.UUID) (*domain.SuccessionSuccessorNomination, error)
	UpdateSuccessionSuccessorNominationStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, actorID *uuid.UUID) (*domain.SuccessionSuccessorNomination, error)
	ListSuccessionSuccessorNominations(ctx context.Context, filter domain.SuccessionFilter) ([]*domain.SuccessionSuccessorNomination, error)
	CreateSuccessionDevelopmentAction(ctx context.Context, item *domain.SuccessionDevelopmentAction, actorID *uuid.UUID) (*domain.SuccessionDevelopmentAction, error)
	UpdateSuccessionDevelopmentAction(ctx context.Context, item *domain.SuccessionDevelopmentAction, actorID *uuid.UUID) (*domain.SuccessionDevelopmentAction, error)
	UpdateSuccessionDevelopmentActionStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, actorID *uuid.UUID) (*domain.SuccessionDevelopmentAction, error)
	ListSuccessionDevelopmentActions(ctx context.Context, filter domain.SuccessionFilter) ([]*domain.SuccessionDevelopmentAction, error)
	CreateSuccessionEvent(ctx context.Context, item *domain.SuccessionEvent, actorID *uuid.UUID) (*domain.SuccessionEvent, error)
	ListSuccessionEvents(ctx context.Context, filter domain.SuccessionFilter, sourceType *string, sourceID *uuid.UUID) ([]*domain.SuccessionEvent, error)
	GetSuccessionSummary(ctx context.Context, tenantID uuid.UUID) ([]*domain.SuccessionSummaryRow, error)
}

type SuccessionReviewCycleCommand struct {
	TenantID             uuid.UUID       `json:"tenant_id"`
	ID                   uuid.UUID       `json:"id,omitempty"`
	Code                 string          `json:"code"`
	Name                 string          `json:"name"`
	Status               string          `json:"status"`
	StartsOn             *time.Time      `json:"starts_on,omitempty"`
	EndsOn               *time.Time      `json:"ends_on,omitempty"`
	ConfidentialityLevel string          `json:"confidentiality_level"`
	Notes                *string         `json:"notes,omitempty"`
	Metadata             json.RawMessage `json:"metadata,omitempty"`
	ActorID              *uuid.UUID      `json:"-"`
}

type SuccessionStatusCommand struct {
	TenantID uuid.UUID  `json:"tenant_id"`
	ID       uuid.UUID  `json:"id"`
	Status   string     `json:"status"`
	Remarks  *string    `json:"remarks,omitempty"`
	ActorID  *uuid.UUID `json:"-"`
}

type SuccessionCriticalRoleCommand struct {
	TenantID                      uuid.UUID       `json:"tenant_id"`
	ID                            uuid.UUID       `json:"id,omitempty"`
	CycleID                       *uuid.UUID      `json:"cycle_id,omitempty"`
	Code                          string          `json:"code"`
	Title                         string          `json:"title"`
	DepartmentID                  *uuid.UUID      `json:"department_id,omitempty"`
	DesignationID                 *uuid.UUID      `json:"designation_id,omitempty"`
	IncumbentWorkerProfileID      *uuid.UUID      `json:"incumbent_worker_profile_id,omitempty"`
	EmergencyCoverWorkerProfileID *uuid.UUID      `json:"emergency_cover_worker_profile_id,omitempty"`
	Criticality                   string          `json:"criticality"`
	ImpactLevel                   string          `json:"impact_level"`
	VacancyRisk                   string          `json:"vacancy_risk"`
	AttritionRisk                 string          `json:"attrition_risk"`
	ReadinessTarget               string          `json:"readiness_target"`
	SuccessorRequiredCount        int32           `json:"successor_required_count"`
	RoleSummary                   *string         `json:"role_summary,omitempty"`
	Status                        string          `json:"status"`
	Metadata                      json.RawMessage `json:"metadata,omitempty"`
	ActorID                       *uuid.UUID      `json:"-"`
}

type SuccessionSuccessorNominationCommand struct {
	TenantID           uuid.UUID       `json:"tenant_id"`
	ID                 uuid.UUID       `json:"id,omitempty"`
	CriticalRoleID     uuid.UUID       `json:"critical_role_id"`
	WorkerProfileID    uuid.UUID       `json:"worker_profile_id"`
	NominatedBy        *uuid.UUID      `json:"nominated_by,omitempty"`
	ReadinessLevel     string          `json:"readiness_level"`
	ReadinessMonths    int32           `json:"readiness_months"`
	PotentialRating    *string         `json:"potential_rating,omitempty"`
	PerformanceRating  *string         `json:"performance_rating,omitempty"`
	RetentionRisk      string          `json:"retention_risk"`
	MobilityPreference *string         `json:"mobility_preference,omitempty"`
	NominationStatus   string          `json:"nomination_status"`
	DevelopmentNotes   *string         `json:"development_notes,omitempty"`
	Metadata           json.RawMessage `json:"metadata,omitempty"`
	ActorID            *uuid.UUID      `json:"-"`
}

type SuccessionDevelopmentActionCommand struct {
	TenantID         uuid.UUID       `json:"tenant_id"`
	ID               uuid.UUID       `json:"id,omitempty"`
	NominationID     *uuid.UUID      `json:"nomination_id,omitempty"`
	CriticalRoleID   *uuid.UUID      `json:"critical_role_id,omitempty"`
	WorkerProfileID  uuid.UUID       `json:"worker_profile_id"`
	ActionType       string          `json:"action_type"`
	Title            string          `json:"title"`
	LearningCourseID *uuid.UUID      `json:"learning_course_id,omitempty"`
	LearningPathID   *uuid.UUID      `json:"learning_path_id,omitempty"`
	OwnerUserID      *uuid.UUID      `json:"owner_user_id,omitempty"`
	DueDate          *time.Time      `json:"due_date,omitempty"`
	Status           string          `json:"status"`
	Notes            *string         `json:"notes,omitempty"`
	Metadata         json.RawMessage `json:"metadata,omitempty"`
	ActorID          *uuid.UUID      `json:"-"`
}
