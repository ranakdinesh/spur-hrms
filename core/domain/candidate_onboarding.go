package domain

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidCandidateOnboardingID     = errors.New("candidate_onboarding_id is required")
	ErrInvalidCandidateOnboardingTaskID = errors.New("candidate_onboarding_task_id is required")
	ErrInvalidCandidateOnboardingStatus = errors.New("candidate onboarding status is invalid")
	ErrCandidateOnboardingNotFound      = errors.New("candidate onboarding not found")
	ErrCandidateOnboardingTaskNotFound  = errors.New("candidate onboarding task not found")
	ErrOnboardingWorkflowNotFound       = errors.New("onboarding workflow not found")
)

type CandidateOnboarding struct {
	ID                     uuid.UUID                   `json:"id"`
	TenantID               uuid.UUID                   `json:"tenant_id"`
	CandidateID            uuid.UUID                   `json:"candidate_id"`
	CandidateFirstname     *string                     `json:"candidate_firstname,omitempty"`
	CandidateLastname      *string                     `json:"candidate_lastname,omitempty"`
	CandidateEmail         *string                     `json:"candidate_email,omitempty"`
	WorkflowID             uuid.UUID                   `json:"workflow_id"`
	WorkflowName           *string                     `json:"workflow_name,omitempty"`
	OnboardingStatus       string                      `json:"onboarding_status"`
	ProgressPercentage     int32                       `json:"progress_percentage"`
	TotalTasks             int32                       `json:"total_tasks"`
	CompletedTasks         int32                       `json:"completed_tasks"`
	RequiredTasks          int32                       `json:"required_tasks"`
	CompletedRequiredTasks int32                       `json:"completed_required_tasks"`
	OverdueTasks           int32                       `json:"overdue_tasks"`
	StartedAt              *time.Time                  `json:"started_at,omitempty"`
	CompletedAt            *time.Time                  `json:"completed_at,omitempty"`
	Inactive               bool                        `json:"inactive"`
	CreatedAt              time.Time                   `json:"created_at"`
	CreatedBy              *uuid.UUID                  `json:"created_by,omitempty"`
	UpdatedAt              time.Time                   `json:"updated_at"`
	UpdatedBy              *uuid.UUID                  `json:"updated_by,omitempty"`
	Tasks                  []*CandidateOnboardingTask  `json:"tasks,omitempty"`
	Events                 []*CandidateOnboardingEvent `json:"events,omitempty"`
}

type CandidateOnboardingTask struct {
	ID                    uuid.UUID  `json:"id"`
	TenantID              uuid.UUID  `json:"tenant_id"`
	CandidateOnboardingID uuid.UUID  `json:"candidate_onboarding_id"`
	OnboardingTaskID      uuid.UUID  `json:"onboarding_task_id"`
	TaskTitle             *string    `json:"task_title,omitempty"`
	TaskDescription       *string    `json:"task_description,omitempty"`
	TaskDueDays           int32      `json:"task_due_days"`
	TaskIsRequired        bool       `json:"task_is_required"`
	TaskSortOrder         int32      `json:"task_sort_order"`
	Status                string     `json:"status"`
	DueAt                 *time.Time `json:"due_at,omitempty"`
	StartedAt             *time.Time `json:"started_at,omitempty"`
	CompletedAt           *time.Time `json:"completed_at,omitempty"`
	CompletedBy           *uuid.UUID `json:"completed_by,omitempty"`
	ReviewedBy            *uuid.UUID `json:"reviewed_by,omitempty"`
	Remarks               *string    `json:"remarks,omitempty"`
	IsOverdue             bool       `json:"is_overdue"`
	Inactive              bool       `json:"inactive"`
	CreatedAt             time.Time  `json:"created_at"`
	CreatedBy             *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt             time.Time  `json:"updated_at"`
	UpdatedBy             *uuid.UUID `json:"updated_by,omitempty"`
}

type CandidateOnboardingEvent struct {
	ID                        uuid.UUID       `json:"id"`
	TenantID                  uuid.UUID       `json:"tenant_id"`
	CandidateOnboardingID     uuid.UUID       `json:"candidate_onboarding_id"`
	CandidateOnboardingTaskID *uuid.UUID      `json:"candidate_onboarding_task_id,omitempty"`
	Action                    string          `json:"action"`
	FromStatus                *string         `json:"from_status,omitempty"`
	ToStatus                  *string         `json:"to_status,omitempty"`
	Remarks                   *string         `json:"remarks,omitempty"`
	Metadata                  json.RawMessage `json:"metadata,omitempty"`
	Inactive                  bool            `json:"inactive"`
	CreatedAt                 time.Time       `json:"created_at"`
	CreatedBy                 *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt                 time.Time       `json:"updated_at"`
	UpdatedBy                 *uuid.UUID      `json:"updated_by,omitempty"`
}

type CandidateOnboardingFilter struct {
	TenantID uuid.UUID
	Status   *string
	Search   *string
	Limit    int32
	Offset   int32
}

type CandidateOnboardingPage struct {
	Items  []*CandidateOnboarding `json:"items"`
	Total  int64                  `json:"total"`
	Limit  int32                  `json:"limit"`
	Offset int32                  `json:"offset"`
}

type CandidateOnboardingInput struct {
	TenantID    uuid.UUID
	CandidateID uuid.UUID
	WorkflowID  uuid.UUID
	Status      string
}

func NewCandidateOnboarding(input CandidateOnboardingInput) (*CandidateOnboarding, error) {
	if input.TenantID == uuid.Nil {
		return nil, ErrInvalidTenantID
	}
	if input.CandidateID == uuid.Nil {
		return nil, ErrInvalidCandidateID
	}
	if input.WorkflowID == uuid.Nil {
		return nil, ErrInvalidOnboardingWorkflowID
	}
	status := strings.TrimSpace(input.Status)
	if status == "" {
		status = OnboardStatusInProgress
	}
	if _, err := ValidateOnboardingStatus(status); err != nil {
		return nil, ErrInvalidCandidateOnboardingStatus
	}
	now := time.Now().UTC()
	return &CandidateOnboarding{TenantID: input.TenantID, CandidateID: input.CandidateID, WorkflowID: input.WorkflowID, OnboardingStatus: status, StartedAt: &now, CreatedAt: now, UpdatedAt: now}, nil
}

func ValidateCandidateOnboardingTaskStatus(value string) (string, error) {
	status, err := ValidateOnboardingStatus(strings.TrimSpace(value))
	if err != nil {
		return "", ErrInvalidCandidateOnboardingStatus
	}
	return status, nil
}
