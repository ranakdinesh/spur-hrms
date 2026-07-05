package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidOnboardingWorkflowID   = errors.New("onboarding_workflow_id is required")
	ErrInvalidOnboardingWorkflowName = errors.New("onboarding workflow name is required")
	ErrInvalidOnboardingTaskID       = errors.New("onboarding_task_id is required")
	ErrInvalidOnboardingTaskTitle    = errors.New("onboarding task title is required")
	ErrInvalidOnboardingDueDays      = errors.New("onboarding task due days is invalid")
	ErrInvalidOnboardingAssignmentID = errors.New("onboarding_assignment_id is required")
	ErrInvalidOnboardingAssignment   = errors.New("onboarding assignment is invalid")
)

type OnboardingWorkflow struct {
	ID          uuid.UUID         `json:"id"`
	TenantID    uuid.UUID         `json:"tenant_id"`
	Name        string            `json:"name"`
	Description *string           `json:"description,omitempty"`
	IsDefault   bool              `json:"is_default"`
	IsActive    bool              `json:"is_active"`
	Inactive    bool              `json:"inactive"`
	CreatedAt   time.Time         `json:"created_at"`
	CreatedBy   *uuid.UUID        `json:"created_by,omitempty"`
	UpdatedAt   time.Time         `json:"updated_at"`
	UpdatedBy   *uuid.UUID        `json:"updated_by,omitempty"`
	Tasks       []*OnboardingTask `json:"tasks,omitempty"`
}

type OnboardingTask struct {
	ID          uuid.UUID  `json:"id"`
	TenantID    uuid.UUID  `json:"tenant_id"`
	WorkflowID  uuid.UUID  `json:"workflow_id"`
	Title       string     `json:"title"`
	Description *string    `json:"description,omitempty"`
	DueDays     int32      `json:"due_days"`
	IsRequired  bool       `json:"is_required"`
	SortOrder   int32      `json:"sort_order"`
	Inactive    bool       `json:"inactive"`
	CreatedAt   time.Time  `json:"created_at"`
	CreatedBy   *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt   time.Time  `json:"updated_at"`
	UpdatedBy   *uuid.UUID `json:"updated_by,omitempty"`
}

type OnboardingWorkflowAssignment struct {
	ID                 uuid.UUID  `json:"id"`
	TenantID           uuid.UUID  `json:"tenant_id"`
	WorkflowID         uuid.UUID  `json:"workflow_id"`
	WorkflowName       *string    `json:"workflow_name,omitempty"`
	Name               string     `json:"name"`
	JobPostingID       *uuid.UUID `json:"job_posting_id,omitempty"`
	JobPostingTitle    *string    `json:"job_posting_title,omitempty"`
	JobPositionID      *uuid.UUID `json:"job_position_id,omitempty"`
	JobPositionTitle   *string    `json:"job_position_title,omitempty"`
	DepartmentID       *uuid.UUID `json:"department_id,omitempty"`
	DepartmentName     *string    `json:"department_name,omitempty"`
	EmploymentTypeID   *uuid.UUID `json:"employment_type_id,omitempty"`
	EmploymentTypeName *string    `json:"employment_type_name,omitempty"`
	Priority           int32      `json:"priority"`
	Inactive           bool       `json:"inactive"`
	CreatedAt          time.Time  `json:"created_at"`
	CreatedBy          *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt          time.Time  `json:"updated_at"`
	UpdatedBy          *uuid.UUID `json:"updated_by,omitempty"`
}

type OnboardingWorkflowInput struct {
	TenantID    uuid.UUID
	Name        string
	Description *string
	IsDefault   bool
	IsActive    bool
}

type OnboardingTaskInput struct {
	TenantID    uuid.UUID
	WorkflowID  uuid.UUID
	Title       string
	Description *string
	DueDays     int32
	IsRequired  bool
	SortOrder   int32
}

type OnboardingAssignmentInput struct {
	TenantID         uuid.UUID
	WorkflowID       uuid.UUID
	Name             string
	JobPostingID     *uuid.UUID
	JobPositionID    *uuid.UUID
	DepartmentID     *uuid.UUID
	EmploymentTypeID *uuid.UUID
	Priority         int32
}

func NewOnboardingWorkflow(input OnboardingWorkflowInput) (*OnboardingWorkflow, error) {
	if input.TenantID == uuid.Nil {
		return nil, ErrInvalidTenantID
	}
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, ErrInvalidOnboardingWorkflowName
	}
	now := time.Now().UTC()
	return &OnboardingWorkflow{TenantID: input.TenantID, Name: name, Description: cleanOptional(input.Description), IsDefault: input.IsDefault, IsActive: input.IsActive, CreatedAt: now, UpdatedAt: now}, nil
}

func NewOnboardingTask(input OnboardingTaskInput) (*OnboardingTask, error) {
	if input.TenantID == uuid.Nil {
		return nil, ErrInvalidTenantID
	}
	if input.WorkflowID == uuid.Nil {
		return nil, ErrInvalidOnboardingWorkflowID
	}
	title := strings.TrimSpace(input.Title)
	if title == "" {
		return nil, ErrInvalidOnboardingTaskTitle
	}
	if input.DueDays < 0 {
		return nil, ErrInvalidOnboardingDueDays
	}
	now := time.Now().UTC()
	return &OnboardingTask{TenantID: input.TenantID, WorkflowID: input.WorkflowID, Title: title, Description: cleanOptional(input.Description), DueDays: input.DueDays, IsRequired: input.IsRequired, SortOrder: input.SortOrder, CreatedAt: now, UpdatedAt: now}, nil
}

func NewOnboardingWorkflowAssignment(input OnboardingAssignmentInput) (*OnboardingWorkflowAssignment, error) {
	if input.TenantID == uuid.Nil {
		return nil, ErrInvalidTenantID
	}
	if input.WorkflowID == uuid.Nil {
		return nil, ErrInvalidOnboardingWorkflowID
	}
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, ErrInvalidOnboardingAssignment
	}
	if input.Priority <= 0 {
		input.Priority = 100
	}
	now := time.Now().UTC()
	return &OnboardingWorkflowAssignment{TenantID: input.TenantID, WorkflowID: input.WorkflowID, Name: name, JobPostingID: cleanUUIDOptional(input.JobPostingID), JobPositionID: cleanUUIDOptional(input.JobPositionID), DepartmentID: cleanUUIDOptional(input.DepartmentID), EmploymentTypeID: cleanUUIDOptional(input.EmploymentTypeID), Priority: input.Priority, CreatedAt: now, UpdatedAt: now}, nil
}
