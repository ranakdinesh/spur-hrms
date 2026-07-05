package domain

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidEmployeeExitID       = errors.New("employee_exit_id is required")
	ErrInvalidEmployeeExitStatus   = errors.New("employee exit status is invalid")
	ErrInvalidEmployeeExitType     = errors.New("employee exit type is invalid")
	ErrInvalidEmployeeExitDate     = errors.New("employee exit dates are invalid")
	ErrEmployeeExitNotFound        = errors.New("employee exit not found")
	ErrEmployeeExitAlreadyActive   = errors.New("employee already has an active exit workflow")
	ErrEmployeeExitTaskNotFound    = errors.New("employee exit task not found")
	ErrEmployeeExitCompletionBlock = errors.New("employee exit cannot be completed until all checklist tasks are completed or waived")
)

const (
	EmployeeExitStatusSubmitted = "submitted"
	EmployeeExitStatusApproved  = "approved"
	EmployeeExitStatusRejected  = "rejected"
	EmployeeExitStatusCompleted = "completed"
	EmployeeExitStatusCanceled  = "canceled"

	EmployeeExitTypeResignation = "resignation"
	EmployeeExitTypeTermination = "termination"
	EmployeeExitTypeRetirement  = "retirement"
	EmployeeExitTypeContractEnd = "contract_end"
	EmployeeExitTypeAbsconding  = "absconding"
	EmployeeExitTypeOther       = "other"

	EmployeeExitTaskPending    = "pending"
	EmployeeExitTaskInProgress = "in_progress"
	EmployeeExitTaskCompleted  = "completed"
	EmployeeExitTaskWaived     = "waived"
	EmployeeExitTaskBlocked    = "blocked"
)

type EmployeeExitRequest struct {
	ID                     uuid.UUID            `json:"id"`
	TenantID               uuid.UUID            `json:"tenant_id"`
	EmployeeID             uuid.UUID            `json:"employee_id"`
	EmployeeUserID         uuid.UUID            `json:"employee_user_id"`
	EmployeeFirstname      *string              `json:"employee_firstname,omitempty"`
	EmployeeLastname       *string              `json:"employee_lastname,omitempty"`
	EmployeeCode           *string              `json:"employee_code,omitempty"`
	EmployeeEmail          *string              `json:"employee_email,omitempty"`
	DepartmentName         *string              `json:"department_name,omitempty"`
	BranchName             *string              `json:"branch_name,omitempty"`
	InitiatedBy            *uuid.UUID           `json:"initiated_by,omitempty"`
	ApprovedBy             *uuid.UUID           `json:"approved_by,omitempty"`
	ApprovedAt             *time.Time           `json:"approved_at,omitempty"`
	CompletedBy            *uuid.UUID           `json:"completed_by,omitempty"`
	CompletedAt            *time.Time           `json:"completed_at,omitempty"`
	Status                 string               `json:"status"`
	ExitType               string               `json:"exit_type"`
	Reason                 *string              `json:"reason,omitempty"`
	ResignationDate        *time.Time           `json:"resignation_date,omitempty"`
	NoticeStartDate        *time.Time           `json:"notice_start_date,omitempty"`
	LastWorkingDate        time.Time            `json:"last_working_date"`
	RequestedRelievingDate *time.Time           `json:"requested_relieving_date,omitempty"`
	ApprovedRelievingDate  *time.Time           `json:"approved_relieving_date,omitempty"`
	FinalSettlementStatus  string               `json:"final_settlement_status"`
	AccessRevocationStatus string               `json:"access_revocation_status"`
	AssetClearanceStatus   string               `json:"asset_clearance_status"`
	HandoverStatus         string               `json:"handover_status"`
	ExitInterviewStatus    string               `json:"exit_interview_status"`
	Notes                  *string              `json:"notes,omitempty"`
	RejectionReason        *string              `json:"rejection_reason,omitempty"`
	CancelReason           *string              `json:"cancel_reason,omitempty"`
	TotalTasks             int32                `json:"total_tasks"`
	CompletedTasks         int32                `json:"completed_tasks"`
	BlockedTasks           int32                `json:"blocked_tasks"`
	Inactive               bool                 `json:"inactive"`
	CreatedAt              time.Time            `json:"created_at"`
	CreatedBy              *uuid.UUID           `json:"created_by,omitempty"`
	UpdatedAt              time.Time            `json:"updated_at"`
	UpdatedBy              *uuid.UUID           `json:"updated_by,omitempty"`
	Tasks                  []*EmployeeExitTask  `json:"tasks,omitempty"`
	Events                 []*EmployeeExitEvent `json:"events,omitempty"`
}

type EmployeeExitTask struct {
	ID             uuid.UUID  `json:"id"`
	TenantID       uuid.UUID  `json:"tenant_id"`
	ExitRequestID  uuid.UUID  `json:"exit_request_id"`
	EmployeeUserID uuid.UUID  `json:"employee_user_id"`
	TaskKey        string     `json:"task_key"`
	Title          string     `json:"title"`
	Description    *string    `json:"description,omitempty"`
	OwnerRole      *string    `json:"owner_role,omitempty"`
	OwnerUserID    *uuid.UUID `json:"owner_user_id,omitempty"`
	DueDate        *time.Time `json:"due_date,omitempty"`
	Status         string     `json:"status"`
	CompletedBy    *uuid.UUID `json:"completed_by,omitempty"`
	CompletedAt    *time.Time `json:"completed_at,omitempty"`
	Remarks        *string    `json:"remarks,omitempty"`
	SortOrder      int32      `json:"sort_order"`
	Inactive       bool       `json:"inactive"`
	CreatedAt      time.Time  `json:"created_at"`
	CreatedBy      *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt      time.Time  `json:"updated_at"`
	UpdatedBy      *uuid.UUID `json:"updated_by,omitempty"`
}

type EmployeeExitEvent struct {
	ID            uuid.UUID       `json:"id"`
	TenantID      uuid.UUID       `json:"tenant_id"`
	ExitRequestID uuid.UUID       `json:"exit_request_id"`
	ExitTaskID    *uuid.UUID      `json:"exit_task_id,omitempty"`
	Action        string          `json:"action"`
	FromStatus    *string         `json:"from_status,omitempty"`
	ToStatus      *string         `json:"to_status,omitempty"`
	Remarks       *string         `json:"remarks,omitempty"`
	Metadata      json.RawMessage `json:"metadata,omitempty"`
	Inactive      bool            `json:"inactive"`
	CreatedAt     time.Time       `json:"created_at"`
	CreatedBy     *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt     time.Time       `json:"updated_at"`
	UpdatedBy     *uuid.UUID      `json:"updated_by,omitempty"`
}

type EmployeeExitFilter struct {
	TenantID       uuid.UUID
	Status         *string
	EmployeeUserID *uuid.UUID
	Search         *string
	Limit          int32
	Offset         int32
}

type EmployeeExitPage struct {
	Items  []*EmployeeExitRequest `json:"items"`
	Total  int64                  `json:"total"`
	Limit  int32                  `json:"limit"`
	Offset int32                  `json:"offset"`
}

type EmployeeExitTaskTemplate struct {
	Key         string
	Title       string
	Description string
	OwnerRole   string
	DueDate     *time.Time
	SortOrder   int32
}

func ValidateEmployeeExitStatus(value string) (string, error) {
	status := strings.TrimSpace(value)
	switch status {
	case EmployeeExitStatusSubmitted, EmployeeExitStatusApproved, EmployeeExitStatusRejected, EmployeeExitStatusCompleted, EmployeeExitStatusCanceled:
		return status, nil
	default:
		return "", ErrInvalidEmployeeExitStatus
	}
}

func ValidateEmployeeExitTaskStatus(value string) (string, error) {
	status := strings.TrimSpace(value)
	switch status {
	case EmployeeExitTaskPending, EmployeeExitTaskInProgress, EmployeeExitTaskCompleted, EmployeeExitTaskWaived, EmployeeExitTaskBlocked:
		return status, nil
	default:
		return "", ErrInvalidEmployeeExitStatus
	}
}

func ValidateEmployeeExitType(value string) (string, error) {
	exitType := strings.TrimSpace(value)
	if exitType == "" {
		exitType = EmployeeExitTypeResignation
	}
	switch exitType {
	case EmployeeExitTypeResignation, EmployeeExitTypeTermination, EmployeeExitTypeRetirement, EmployeeExitTypeContractEnd, EmployeeExitTypeAbsconding, EmployeeExitTypeOther:
		return exitType, nil
	default:
		return "", ErrInvalidEmployeeExitType
	}
}

func DefaultEmployeeExitTaskTemplates(lastWorkingDate time.Time) []EmployeeExitTaskTemplate {
	due := func(daysBefore int) *time.Time {
		date := lastWorkingDate.AddDate(0, 0, -daysBefore)
		return &date
	}
	return []EmployeeExitTaskTemplate{
		{Key: "handover_plan", Title: "Handover plan", Description: "Capture open work, contacts, credentials ownership, and transfer owner.", OwnerRole: "manager", DueDate: due(7), SortOrder: 10},
		{Key: "knowledge_transfer", Title: "Knowledge transfer", Description: "Complete knowledge transfer sessions and attach/link notes where applicable.", OwnerRole: "manager", DueDate: due(5), SortOrder: 20},
		{Key: "asset_return", Title: "Asset return", Description: "Recover laptop, ID card, access cards, SIM, devices, and tenant-owned equipment.", OwnerRole: "hr_admin", DueDate: due(2), SortOrder: 30},
		{Key: "access_revocation", Title: "Access revocation", Description: "Schedule or complete revocation for HRMS, email, SaaS tools, VPN, and physical access.", OwnerRole: "it_admin", DueDate: &lastWorkingDate, SortOrder: 40},
		{Key: "exit_interview", Title: "Exit interview", Description: "Record interview completion or mark skipped with reason.", OwnerRole: "hr", DueDate: due(3), SortOrder: 50},
		{Key: "final_settlement", Title: "Final settlement readiness", Description: "Validate attendance, LOP/LWP, recoveries, reimbursements, gratuity/bonus, and F&F readiness.", OwnerRole: "payroll", DueDate: &lastWorkingDate, SortOrder: 60},
		{Key: "letters_documents", Title: "Letters and documents", Description: "Prepare experience/relieving letter readiness and final document handover.", OwnerRole: "hr", DueDate: &lastWorkingDate, SortOrder: 70},
	}
}
