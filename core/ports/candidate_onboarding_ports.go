package ports

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

type CandidateOnboardingRepo interface {
	CreateCandidateOnboarding(ctx context.Context, item *domain.CandidateOnboarding, actorID *uuid.UUID) (*domain.CandidateOnboarding, error)
	ListCandidateOnboardings(ctx context.Context, filter domain.CandidateOnboardingFilter) (*domain.CandidateOnboardingPage, error)
	GetCandidateOnboarding(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.CandidateOnboarding, error)
	GetCandidateOnboardingByCandidate(ctx context.Context, tenantID uuid.UUID, candidateID uuid.UUID) (*domain.CandidateOnboarding, error)
	GetDefaultOnboardingWorkflow(ctx context.Context, tenantID uuid.UUID) (*domain.OnboardingWorkflow, error)
	ResolveOnboardingWorkflowForCandidate(ctx context.Context, tenantID uuid.UUID, candidateID uuid.UUID) (*domain.OnboardingWorkflow, error)
	CreateCandidateOnboardingTasksFromWorkflow(ctx context.Context, tenantID uuid.UUID, candidateOnboardingID uuid.UUID, workflowID uuid.UUID, actorID *uuid.UUID) ([]*domain.CandidateOnboardingTask, error)
	RecalculateCandidateOnboardingProgress(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) (*domain.CandidateOnboarding, error)
	ListCandidateOnboardingTasks(ctx context.Context, tenantID uuid.UUID, candidateOnboardingID uuid.UUID) ([]*domain.CandidateOnboardingTask, error)
	GetCandidateOnboardingTask(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.CandidateOnboardingTask, error)
	UpdateCandidateOnboardingTaskStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, remarks *string, actorID *uuid.UUID) (*domain.CandidateOnboardingTask, error)
	CreateCandidateOnboardingEvent(ctx context.Context, event *domain.CandidateOnboardingEvent, actorID *uuid.UUID) (*domain.CandidateOnboardingEvent, error)
	ListCandidateOnboardingEvents(ctx context.Context, tenantID uuid.UUID, candidateOnboardingID uuid.UUID) ([]*domain.CandidateOnboardingEvent, error)
	DeleteCandidateOnboarding(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
}

type StartCandidateOnboardingCommand struct {
	TenantID    uuid.UUID  `json:"tenant_id"`
	CandidateID uuid.UUID  `json:"candidate_id"`
	WorkflowID  *uuid.UUID `json:"workflow_id,omitempty"`
	ActorID     *uuid.UUID `json:"-"`
}

type CandidateOnboardingTaskStatusCommand struct {
	TenantID uuid.UUID  `json:"tenant_id"`
	TaskID   uuid.UUID  `json:"task_id"`
	Status   string     `json:"status"`
	Remarks  *string    `json:"remarks,omitempty"`
	ActorID  *uuid.UUID `json:"-"`
}

type CandidateOnboardingEventCommand struct {
	TenantID                  uuid.UUID       `json:"tenant_id"`
	CandidateOnboardingID     uuid.UUID       `json:"candidate_onboarding_id"`
	CandidateOnboardingTaskID *uuid.UUID      `json:"candidate_onboarding_task_id,omitempty"`
	Action                    string          `json:"action"`
	FromStatus                *string         `json:"from_status,omitempty"`
	ToStatus                  *string         `json:"to_status,omitempty"`
	Remarks                   *string         `json:"remarks,omitempty"`
	Metadata                  json.RawMessage `json:"metadata,omitempty"`
	ActorID                   *uuid.UUID      `json:"-"`
}
