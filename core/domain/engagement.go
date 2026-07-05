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
	EngagementTypeEmployeeAssignment = "employee_assignment"
	EngagementTypeFixedTerm          = "fixed_term"
	EngagementTypeProject            = "project"
	EngagementTypeHourly             = "hourly"
	EngagementTypeRetainer           = "retainer"
	EngagementTypeStipend            = "stipend"
	EngagementTypeAgency             = "agency"
	EngagementTypeConsulting         = "consulting"

	EngagementStatusDraft      = "draft"
	EngagementStatusActive     = "active"
	EngagementStatusPaused     = "paused"
	EngagementStatusCompleted  = "completed"
	EngagementStatusTerminated = "terminated"
	EngagementStatusCancelled  = "cancelled"

	EngagementRateNone      = "none"
	EngagementRateHour      = "hour"
	EngagementRateDay       = "day"
	EngagementRateMonth     = "month"
	EngagementRateMilestone = "milestone"
	EngagementRateRetainer  = "retainer"
	EngagementRateStipend   = "stipend"

	EngagementRenewalNotRequired = "not_required"
	EngagementRenewalPending     = "pending"
	EngagementRenewalRenewed     = "renewed"
	EngagementRenewalNotRenewed  = "not_renewed"
)

var (
	ErrInvalidEngagementID       = errors.New("engagement_id is required")
	ErrInvalidEngagement         = errors.New("engagement is invalid")
	ErrInvalidEngagementType     = errors.New("engagement type is invalid")
	ErrInvalidEngagementStatus   = errors.New("engagement status is invalid")
	ErrInvalidEngagementRateUnit = errors.New("engagement rate_unit is invalid")
	ErrInvalidEngagementCurrency = errors.New("engagement currency_code must be a 3-letter code")
	ErrInvalidEngagementDates    = errors.New("engagement dates are invalid")
	ErrInvalidEngagementBudget   = errors.New("engagement hours_budget and rate_amount must be non-negative")
	ErrInvalidEngagementMetadata = errors.New("engagement metadata must be a valid JSON object")
)

type Engagement struct {
	ID                 uuid.UUID       `json:"id"`
	TenantID           uuid.UUID       `json:"tenant_id"`
	WorkerProfileID    uuid.UUID       `json:"worker_profile_id"`
	EngagementCode     *string         `json:"engagement_code,omitempty"`
	Title              string          `json:"title"`
	Description        *string         `json:"description,omitempty"`
	EngagementType     string          `json:"engagement_type"`
	Status             string          `json:"status"`
	StartDate          time.Time       `json:"start_date"`
	EndDate            *time.Time      `json:"end_date,omitempty"`
	HoursBudget        *float64        `json:"hours_budget,omitempty"`
	RateAmount         *float64        `json:"rate_amount,omitempty"`
	CurrencyCode       string          `json:"currency_code"`
	RateUnit           string          `json:"rate_unit"`
	BranchID           *uuid.UUID      `json:"branch_id,omitempty"`
	DepartmentID       *uuid.UUID      `json:"department_id,omitempty"`
	ReportingManagerID *uuid.UUID      `json:"reporting_manager_id,omitempty"`
	ProjectLabel       *string         `json:"project_label,omitempty"`
	ProjectCode        *string         `json:"project_code,omitempty"`
	CostCenter         *string         `json:"cost_center,omitempty"`
	RenewalDueDate     *time.Time      `json:"renewal_due_date,omitempty"`
	RenewalStatus      string          `json:"renewal_status"`
	TerminationReason  *string         `json:"termination_reason,omitempty"`
	TerminatedAt       *time.Time      `json:"terminated_at,omitempty"`
	Notes              *string         `json:"notes,omitempty"`
	Metadata           json.RawMessage `json:"metadata,omitempty"`
	Inactive           bool            `json:"inactive"`
	CreatedAt          time.Time       `json:"created_at"`
	CreatedBy          *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt          time.Time       `json:"updated_at"`
	UpdatedBy          *uuid.UUID      `json:"updated_by,omitempty"`
}

type EngagementListItem struct {
	Engagement
	WorkerDisplayName   string     `json:"worker_display_name"`
	WorkerCode          *string    `json:"worker_code,omitempty"`
	EmployeeID          *uuid.UUID `json:"employee_id,omitempty"`
	WorkerTypeName      string     `json:"worker_type_name"`
	ClassificationGroup string     `json:"classification_group"`
	BranchName          *string    `json:"branch_name,omitempty"`
	DepartmentName      *string    `json:"department_name,omitempty"`
}

type EngagementInput struct {
	TenantID           uuid.UUID
	WorkerProfileID    uuid.UUID
	EngagementCode     *string
	Title              string
	Description        *string
	EngagementType     string
	Status             string
	StartDate          *time.Time
	EndDate            *time.Time
	HoursBudget        *float64
	RateAmount         *float64
	CurrencyCode       string
	RateUnit           string
	BranchID           *uuid.UUID
	DepartmentID       *uuid.UUID
	ReportingManagerID *uuid.UUID
	ProjectLabel       *string
	ProjectCode        *string
	CostCenter         *string
	RenewalDueDate     *time.Time
	RenewalStatus      string
	TerminationReason  *string
	TerminatedAt       *time.Time
	Notes              *string
	Metadata           json.RawMessage
}

type EngagementFilter struct {
	TenantID        uuid.UUID
	WorkerProfileID *uuid.UUID
	EngagementType  *string
	Status          *string
	DepartmentID    *uuid.UUID
	Search          *string
}

func NewEngagement(input EngagementInput) (*Engagement, error) {
	if input.TenantID == uuid.Nil || input.WorkerProfileID == uuid.Nil {
		return nil, ErrInvalidEngagement
	}
	title := strings.TrimSpace(input.Title)
	if title == "" {
		return nil, ErrInvalidEngagement
	}
	engagementType := normalizeWorkerProfileEnum(input.EngagementType, "")
	if !containsString(engagementTypes(), engagementType) {
		return nil, ErrInvalidEngagementType
	}
	status := normalizeWorkerProfileEnum(input.Status, EngagementStatusDraft)
	if !containsString(engagementStatuses(), status) {
		return nil, ErrInvalidEngagementStatus
	}
	if input.StartDate == nil || input.StartDate.IsZero() {
		return nil, ErrInvalidEngagementDates
	}
	startDate := datePtrUTC(input.StartDate)
	endDate := datePtrUTC(input.EndDate)
	if endDate != nil && endDate.Before(*startDate) {
		return nil, ErrInvalidEngagementDates
	}
	renewalDueDate := datePtrUTC(input.RenewalDueDate)
	if renewalDueDate != nil && endDate != nil && renewalDueDate.After(*endDate) {
		return nil, ErrInvalidEngagementDates
	}
	if invalidPositiveFloat(input.HoursBudget) || invalidPositiveFloat(input.RateAmount) {
		return nil, ErrInvalidEngagementBudget
	}
	currency := strings.ToUpper(strings.TrimSpace(input.CurrencyCode))
	if currency == "" {
		currency = "INR"
	}
	if len(currency) != 3 {
		return nil, ErrInvalidEngagementCurrency
	}
	rateUnit := normalizeWorkerProfileEnum(input.RateUnit, EngagementRateNone)
	if !containsString([]string{EngagementRateNone, EngagementRateHour, EngagementRateDay, EngagementRateMonth, EngagementRateMilestone, EngagementRateRetainer, EngagementRateStipend}, rateUnit) {
		return nil, ErrInvalidEngagementRateUnit
	}
	renewalStatus := normalizeWorkerProfileEnum(input.RenewalStatus, EngagementRenewalNotRequired)
	if !containsString([]string{EngagementRenewalNotRequired, EngagementRenewalPending, EngagementRenewalRenewed, EngagementRenewalNotRenewed}, renewalStatus) {
		return nil, ErrInvalidEngagement
	}
	metadata := normalizeWorkerJSONObject(input.Metadata, "{}")
	if !json.Valid(metadata) || !jsonObject(metadata) {
		return nil, ErrInvalidEngagementMetadata
	}
	terminatedAt := input.TerminatedAt
	if status == EngagementStatusTerminated && terminatedAt == nil {
		now := time.Now().UTC()
		terminatedAt = &now
	}
	if status != EngagementStatusTerminated {
		terminatedAt = nil
	}
	now := time.Now().UTC()
	return &Engagement{
		TenantID:           input.TenantID,
		WorkerProfileID:    input.WorkerProfileID,
		EngagementCode:     cleanOptional(input.EngagementCode),
		Title:              title,
		Description:        cleanOptional(input.Description),
		EngagementType:     engagementType,
		Status:             status,
		StartDate:          *startDate,
		EndDate:            endDate,
		HoursBudget:        cleanFloatOptional(input.HoursBudget),
		RateAmount:         cleanFloatOptional(input.RateAmount),
		CurrencyCode:       currency,
		RateUnit:           rateUnit,
		BranchID:           cleanUUIDOptional(input.BranchID),
		DepartmentID:       cleanUUIDOptional(input.DepartmentID),
		ReportingManagerID: cleanUUIDOptional(input.ReportingManagerID),
		ProjectLabel:       cleanOptional(input.ProjectLabel),
		ProjectCode:        cleanOptional(input.ProjectCode),
		CostCenter:         cleanOptional(input.CostCenter),
		RenewalDueDate:     renewalDueDate,
		RenewalStatus:      renewalStatus,
		TerminationReason:  cleanOptional(input.TerminationReason),
		TerminatedAt:       terminatedAt,
		Notes:              cleanOptional(input.Notes),
		Metadata:           metadata,
		CreatedAt:          now,
		UpdatedAt:          now,
	}, nil
}

func ValidateEngagementStatus(value string) (string, error) {
	status := normalizeWorkerProfileEnum(value, "")
	if !containsString(engagementStatuses(), status) {
		return "", ErrInvalidEngagementStatus
	}
	return status, nil
}

func engagementTypes() []string {
	return []string{EngagementTypeEmployeeAssignment, EngagementTypeFixedTerm, EngagementTypeProject, EngagementTypeHourly, EngagementTypeRetainer, EngagementTypeStipend, EngagementTypeAgency, EngagementTypeConsulting}
}

func engagementStatuses() []string {
	return []string{EngagementStatusDraft, EngagementStatusActive, EngagementStatusPaused, EngagementStatusCompleted, EngagementStatusTerminated, EngagementStatusCancelled}
}

func invalidPositiveFloat(value *float64) bool {
	return value != nil && (*value < 0 || math.IsNaN(*value) || math.IsInf(*value, 0))
}

func cleanFloatOptional(value *float64) *float64 {
	if value == nil || invalidPositiveFloat(value) {
		return nil
	}
	return value
}
