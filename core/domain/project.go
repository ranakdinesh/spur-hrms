package domain

import (
	"encoding/json"
	"errors"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	ProjectStatusDraft     = "draft"
	ProjectStatusActive    = "active"
	ProjectStatusOnHold    = "on_hold"
	ProjectStatusCompleted = "completed"
	ProjectStatusCancelled = "cancelled"

	ProjectBillingNone      = "none"
	ProjectBillingFixed     = "fixed"
	ProjectBillingHourly    = "hourly"
	ProjectBillingMilestone = "milestone"
	ProjectBillingRetainer  = "retainer"

	ProjectPriorityLow      = "low"
	ProjectPriorityNormal   = "normal"
	ProjectPriorityHigh     = "high"
	ProjectPriorityCritical = "critical"
)

var (
	ErrInvalidProjectID       = errors.New("project_id is required")
	ErrInvalidProject         = errors.New("project is invalid")
	ErrInvalidProjectStatus   = errors.New("project status is invalid")
	ErrInvalidProjectBilling  = errors.New("project billing_type is invalid")
	ErrInvalidProjectPriority = errors.New("project priority is invalid")
	ErrInvalidProjectDates    = errors.New("project dates are invalid")
	ErrInvalidProjectBudget   = errors.New("project budget_amount must be non-negative")
	ErrInvalidProjectCurrency = errors.New("project currency_code must be a 3-letter code")
	ErrInvalidProjectMetadata = errors.New("project metadata must be a valid JSON object")
)

type Project struct {
	ID               uuid.UUID       `json:"id"`
	TenantID         uuid.UUID       `json:"tenant_id"`
	ProjectCode      *string         `json:"project_code,omitempty"`
	Name             string          `json:"name"`
	Description      *string         `json:"description,omitempty"`
	Status           string          `json:"status"`
	DepartmentID     *uuid.UUID      `json:"department_id,omitempty"`
	BranchID         *uuid.UUID      `json:"branch_id,omitempty"`
	ProjectManagerID *uuid.UUID      `json:"project_manager_id,omitempty"`
	StartDate        *time.Time      `json:"start_date,omitempty"`
	DueDate          *time.Time      `json:"due_date,omitempty"`
	CompletedAt      *time.Time      `json:"completed_at,omitempty"`
	BudgetAmount     *float64        `json:"budget_amount,omitempty"`
	CurrencyCode     string          `json:"currency_code"`
	BillingType      string          `json:"billing_type"`
	ClientLabel      *string         `json:"client_label,omitempty"`
	Priority         string          `json:"priority"`
	Notes            *string         `json:"notes,omitempty"`
	Metadata         json.RawMessage `json:"metadata,omitempty"`
	Inactive         bool            `json:"inactive"`
	CreatedAt        time.Time       `json:"created_at"`
	CreatedBy        *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt        time.Time       `json:"updated_at"`
	UpdatedBy        *uuid.UUID      `json:"updated_by,omitempty"`
}

type ProjectListItem struct {
	Project
	DepartmentName          *string  `json:"department_name,omitempty"`
	BranchName              *string  `json:"branch_name,omitempty"`
	MilestoneCount          int32    `json:"milestone_count"`
	SubmittedMilestoneCount int32    `json:"submitted_milestone_count"`
	AcceptedMilestoneCount  int32    `json:"accepted_milestone_count"`
	RejectedMilestoneCount  int32    `json:"rejected_milestone_count"`
	MilestoneAmount         float64  `json:"milestone_amount"`
	AcceptedAmount          float64  `json:"accepted_amount"`
	RemainingBudgetAmount   *float64 `json:"remaining_budget_amount,omitempty"`
}

type ProjectInput struct {
	TenantID         uuid.UUID
	ProjectCode      *string
	Name             string
	Description      *string
	Status           string
	DepartmentID     *uuid.UUID
	BranchID         *uuid.UUID
	ProjectManagerID *uuid.UUID
	StartDate        *time.Time
	DueDate          *time.Time
	CompletedAt      *time.Time
	BudgetAmount     *float64
	CurrencyCode     string
	BillingType      string
	ClientLabel      *string
	Priority         string
	Notes            *string
	Metadata         json.RawMessage
}

type ProjectFilter struct {
	TenantID         uuid.UUID
	Status           *string
	DepartmentID     *uuid.UUID
	BranchID         *uuid.UUID
	ProjectManagerID *uuid.UUID
	Search           *string
}

func NewProject(input ProjectInput) (*Project, error) {
	if input.TenantID == uuid.Nil {
		return nil, ErrInvalidTenantID
	}
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, ErrInvalidProject
	}
	status := normalizeWorkerProfileEnum(input.Status, ProjectStatusDraft)
	if !containsString(projectStatuses(), status) {
		return nil, ErrInvalidProjectStatus
	}
	startDate := datePtrUTC(input.StartDate)
	dueDate := datePtrUTC(input.DueDate)
	if startDate != nil && dueDate != nil && dueDate.Before(*startDate) {
		return nil, ErrInvalidProjectDates
	}
	if invalidPositiveFloat(input.BudgetAmount) {
		return nil, ErrInvalidProjectBudget
	}
	currency := normalizeCurrencyCode(input.CurrencyCode)
	if len(currency) != 3 {
		return nil, ErrInvalidProjectCurrency
	}
	billingType := normalizeWorkerProfileEnum(input.BillingType, ProjectBillingNone)
	if !containsString(projectBillingTypes(), billingType) {
		return nil, ErrInvalidProjectBilling
	}
	priority := normalizeWorkerProfileEnum(input.Priority, ProjectPriorityNormal)
	if !containsString(projectPriorities(), priority) {
		return nil, ErrInvalidProjectPriority
	}
	metadata := normalizeWorkerJSONObject(input.Metadata, "{}")
	if !json.Valid(metadata) || !jsonObject(metadata) {
		return nil, ErrInvalidProjectMetadata
	}
	completedAt := input.CompletedAt
	if status == ProjectStatusCompleted && completedAt == nil {
		now := time.Now().UTC()
		completedAt = &now
	}
	if status != ProjectStatusCompleted {
		completedAt = nil
	}
	now := time.Now().UTC()
	return &Project{
		TenantID:         input.TenantID,
		ProjectCode:      cleanOptional(input.ProjectCode),
		Name:             name,
		Description:      cleanOptional(input.Description),
		Status:           status,
		DepartmentID:     cleanUUIDOptional(input.DepartmentID),
		BranchID:         cleanUUIDOptional(input.BranchID),
		ProjectManagerID: cleanUUIDOptional(input.ProjectManagerID),
		StartDate:        startDate,
		DueDate:          dueDate,
		CompletedAt:      completedAt,
		BudgetAmount:     cleanFloatOptional(input.BudgetAmount),
		CurrencyCode:     currency,
		BillingType:      billingType,
		ClientLabel:      cleanOptional(input.ClientLabel),
		Priority:         priority,
		Notes:            cleanOptional(input.Notes),
		Metadata:         metadata,
		CreatedAt:        now,
		UpdatedAt:        now,
	}, nil
}

func ValidateProjectStatus(value string) (string, error) {
	status := normalizeWorkerProfileEnum(value, "")
	if !containsString(projectStatuses(), status) {
		return "", ErrInvalidProjectStatus
	}
	return status, nil
}

func ProjectBudgetRemaining(budget *float64, accepted float64) *float64 {
	if budget == nil || math.IsNaN(accepted) || math.IsInf(accepted, 0) {
		return nil
	}
	remaining := *budget - accepted
	return &remaining
}

func projectStatuses() []string {
	return []string{ProjectStatusDraft, ProjectStatusActive, ProjectStatusOnHold, ProjectStatusCompleted, ProjectStatusCancelled}
}

func projectBillingTypes() []string {
	return []string{ProjectBillingNone, ProjectBillingFixed, ProjectBillingHourly, ProjectBillingMilestone, ProjectBillingRetainer}
}

func projectPriorities() []string {
	return []string{ProjectPriorityLow, ProjectPriorityNormal, ProjectPriorityHigh, ProjectPriorityCritical}
}

func normalizeCurrencyCode(value string) string {
	currency := strings.ToUpper(strings.TrimSpace(value))
	if currency == "" {
		return "INR"
	}
	return currency
}
