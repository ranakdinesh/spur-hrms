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
	WorkLogStatusDraft     = "draft"
	WorkLogStatusSubmitted = "submitted"
	WorkLogStatusApproved  = "approved"
	WorkLogStatusRejected  = "rejected"
	WorkLogStatusCancelled = "cancelled"
)

var (
	ErrInvalidWorkLogID       = errors.New("work_log_id is required")
	ErrInvalidWorkLog         = errors.New("work log is invalid")
	ErrInvalidWorkLogDate     = errors.New("work log date is invalid")
	ErrInvalidWorkLogHours    = errors.New("work log hours must be greater than 0 and no more than 24")
	ErrInvalidWorkLogStatus   = errors.New("work log status is invalid")
	ErrInvalidWorkLogMetadata = errors.New("work log metadata must be a valid JSON object")
	ErrWorkLogBudgetExceeded  = errors.New("work log hours exceed engagement hours budget")
	ErrInvalidWorkLogWorkflow = errors.New("work log workflow transition is invalid")
)

type WorkLog struct {
	ID                   uuid.UUID       `json:"id"`
	TenantID             uuid.UUID       `json:"tenant_id"`
	EngagementID         uuid.UUID       `json:"engagement_id"`
	WorkerProfileID      uuid.UUID       `json:"worker_profile_id"`
	LogDate              time.Time       `json:"log_date"`
	HoursWorked          float64         `json:"hours_worked"`
	BillableHours        *float64        `json:"billable_hours,omitempty"`
	WorkSummary          *string         `json:"work_summary,omitempty"`
	DeliverableReference *string         `json:"deliverable_reference,omitempty"`
	Status               string          `json:"status"`
	SubmittedAt          *time.Time      `json:"submitted_at,omitempty"`
	SubmittedBy          *uuid.UUID      `json:"submitted_by,omitempty"`
	ReviewedAt           *time.Time      `json:"reviewed_at,omitempty"`
	ReviewedBy           *uuid.UUID      `json:"reviewed_by,omitempty"`
	ReviewComment        *string         `json:"review_comment,omitempty"`
	Metadata             json.RawMessage `json:"metadata,omitempty"`
	Inactive             bool            `json:"inactive"`
	CreatedAt            time.Time       `json:"created_at"`
	CreatedBy            *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt            time.Time       `json:"updated_at"`
	UpdatedBy            *uuid.UUID      `json:"updated_by,omitempty"`
}

type WorkLogListItem struct {
	WorkLog
	EngagementTitle    string     `json:"engagement_title"`
	EngagementCode     *string    `json:"engagement_code,omitempty"`
	ProjectLabel       *string    `json:"project_label,omitempty"`
	ProjectCode        *string    `json:"project_code,omitempty"`
	CostCenter         *string    `json:"cost_center,omitempty"`
	WorkerDisplayName  string     `json:"worker_display_name"`
	WorkerCode         *string    `json:"worker_code,omitempty"`
	EmployeeID         *uuid.UUID `json:"employee_id,omitempty"`
	ReportingManagerID *uuid.UUID `json:"reporting_manager_id,omitempty"`
	DepartmentID       *uuid.UUID `json:"department_id,omitempty"`
	DepartmentName     *string    `json:"department_name,omitempty"`
	BranchID           *uuid.UUID `json:"branch_id,omitempty"`
	BranchName         *string    `json:"branch_name,omitempty"`
}

type WorkLogRollup struct {
	TenantID          uuid.UUID `json:"tenant_id"`
	EngagementID      uuid.UUID `json:"engagement_id"`
	EngagementTitle   string    `json:"engagement_title"`
	EngagementCode    *string   `json:"engagement_code,omitempty"`
	WorkerProfileID   uuid.UUID `json:"worker_profile_id"`
	WorkerDisplayName string    `json:"worker_display_name"`
	LogCount          int32     `json:"log_count"`
	TotalHours        float64   `json:"total_hours"`
	BillableHours     float64   `json:"billable_hours"`
	ApprovedHours     float64   `json:"approved_hours"`
	SubmittedHours    float64   `json:"submitted_hours"`
	RejectedHours     float64   `json:"rejected_hours"`
	HoursBudget       *float64  `json:"hours_budget,omitempty"`
	RemainingHours    *float64  `json:"remaining_hours,omitempty"`
}

type WorkLogBudgetUsage struct {
	HoursBudget *float64
	UsedHours   float64
}

type WorkLogInput struct {
	TenantID             uuid.UUID
	EngagementID         uuid.UUID
	WorkerProfileID      uuid.UUID
	LogDate              *time.Time
	HoursWorked          float64
	BillableHours        *float64
	WorkSummary          *string
	DeliverableReference *string
	Status               string
	SubmittedAt          *time.Time
	SubmittedBy          *uuid.UUID
	ReviewedAt           *time.Time
	ReviewedBy           *uuid.UUID
	ReviewComment        *string
	Metadata             json.RawMessage
}

type WorkLogFilter struct {
	TenantID        uuid.UUID
	EngagementID    *uuid.UUID
	WorkerProfileID *uuid.UUID
	Status          *string
	DateFrom        *time.Time
	DateTo          *time.Time
	Search          *string
}

func NewWorkLog(input WorkLogInput) (*WorkLog, error) {
	if input.TenantID == uuid.Nil || input.EngagementID == uuid.Nil || input.WorkerProfileID == uuid.Nil {
		return nil, ErrInvalidWorkLog
	}
	if input.LogDate == nil || input.LogDate.IsZero() {
		return nil, ErrInvalidWorkLogDate
	}
	logDate := datePtrUTC(input.LogDate)
	if invalidWorkLogHours(input.HoursWorked) {
		return nil, ErrInvalidWorkLogHours
	}
	if input.BillableHours != nil && (*input.BillableHours < 0 || *input.BillableHours > input.HoursWorked || math.IsNaN(*input.BillableHours) || math.IsInf(*input.BillableHours, 0)) {
		return nil, ErrInvalidWorkLogHours
	}
	status := normalizeWorkerProfileEnum(input.Status, WorkLogStatusDraft)
	if !containsString(workLogStatuses(), status) {
		return nil, ErrInvalidWorkLogStatus
	}
	if (status == WorkLogStatusSubmitted || status == WorkLogStatusApproved) && strings.TrimSpace(valueOrEmpty(input.WorkSummary)) == "" {
		return nil, ErrInvalidWorkLog
	}
	metadata := normalizeWorkerJSONObject(input.Metadata, "{}")
	if !json.Valid(metadata) || !jsonObject(metadata) {
		return nil, ErrInvalidWorkLogMetadata
	}
	now := time.Now().UTC()
	return &WorkLog{
		TenantID:             input.TenantID,
		EngagementID:         input.EngagementID,
		WorkerProfileID:      input.WorkerProfileID,
		LogDate:              *logDate,
		HoursWorked:          input.HoursWorked,
		BillableHours:        cleanFloatOptional(input.BillableHours),
		WorkSummary:          cleanOptional(input.WorkSummary),
		DeliverableReference: cleanOptional(input.DeliverableReference),
		Status:               status,
		SubmittedAt:          input.SubmittedAt,
		SubmittedBy:          cleanUUIDOptional(input.SubmittedBy),
		ReviewedAt:           input.ReviewedAt,
		ReviewedBy:           cleanUUIDOptional(input.ReviewedBy),
		ReviewComment:        cleanOptional(input.ReviewComment),
		Metadata:             metadata,
		CreatedAt:            now,
		UpdatedAt:            now,
	}, nil
}

func ValidateWorkLogStatus(value string) (string, error) {
	status := normalizeWorkerProfileEnum(value, "")
	if !containsString(workLogStatuses(), status) {
		return "", ErrInvalidWorkLogStatus
	}
	return status, nil
}

func WorkLogBudgetRemaining(budget *float64, used float64) *float64 {
	if budget == nil {
		return nil
	}
	remaining := *budget - used
	return &remaining
}

func workLogStatuses() []string {
	return []string{WorkLogStatusDraft, WorkLogStatusSubmitted, WorkLogStatusApproved, WorkLogStatusRejected, WorkLogStatusCancelled}
}

func invalidWorkLogHours(value float64) bool {
	return value <= 0 || value > 24 || math.IsNaN(value) || math.IsInf(value, 0)
}

func valueOrEmpty(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
