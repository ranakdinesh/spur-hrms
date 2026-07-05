package domain

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

const (
	WorkItemActionApprove    = "approve"
	WorkItemActionReject     = "reject"
	WorkItemActionComplete   = "complete"
	WorkItemActionRespond    = "respond"
	WorkItemActionOpenRecord = "open_record"
)

var (
	ErrOperationsWorkbenchCardNotFound      = errors.New("operations workbench card not found")
	ErrUnsupportedOperationsWorkbenchAction = errors.New("operations workbench action is not supported")
)

type OperationsWorkbenchFilter struct {
	TenantID  uuid.UUID
	Lane      *string
	Category  *string
	Severity  *string
	Search    *string
	Limit     int32
	Offset    int32
	Generated time.Time
}

type OperationsWorkbenchCard struct {
	TenantID       uuid.UUID        `json:"tenant_id"`
	CardKey        string           `json:"card_key"`
	Lane           string           `json:"lane"`
	Category       string           `json:"category"`
	SourceModule   string           `json:"source_module"`
	SourceType     string           `json:"source_type"`
	SourceID       uuid.UUID        `json:"source_id"`
	EmployeeUserID *uuid.UUID       `json:"employee_user_id,omitempty"`
	Title          string           `json:"title"`
	Summary        string           `json:"summary"`
	Status         string           `json:"status"`
	Severity       string           `json:"severity"`
	Priority       int32            `json:"priority"`
	DueAt          *time.Time       `json:"due_at,omitempty"`
	DetectedAt     time.Time        `json:"detected_at"`
	ActionLabel    string           `json:"action_label"`
	RouteSection   string           `json:"route_section"`
	RouteRecordID  *uuid.UUID       `json:"route_record_id,omitempty"`
	Metadata       json.RawMessage  `json:"metadata,omitempty"`
	Actions        []WorkItemAction `json:"actions,omitempty"`
}

type WorkItemAction struct {
	Key                string `json:"key"`
	Label              string `json:"label"`
	Tone               string `json:"tone,omitempty"`
	Primary            bool   `json:"primary,omitempty"`
	Inline             bool   `json:"inline"`
	RequiresRemarks    bool   `json:"requires_remarks,omitempty"`
	RemarksPlaceholder string `json:"remarks_placeholder,omitempty"`
	CompletionBadge    string `json:"completion_badge,omitempty"`
}

type OperationsWorkbenchSummary struct {
	Total             int32            `json:"total"`
	HighPriority      int32            `json:"high_priority"`
	Overdue           int32            `json:"overdue"`
	DueToday          int32            `json:"due_today"`
	PayrollBlockers   int32            `json:"payroll_blockers"`
	Approvals         int32            `json:"approvals"`
	Exceptions        int32            `json:"exceptions"`
	Joining           int32            `json:"joining"`
	Exits             int32            `json:"exits"`
	Compliance        int32            `json:"compliance"`
	Documents         int32            `json:"documents"`
	AIRecommendations int32            `json:"ai_recommendations"`
	EmployeeRequests  int32            `json:"employee_requests"`
	ByLane            map[string]int32 `json:"by_lane"`
	ByCategory        map[string]int32 `json:"by_category"`
	BySeverity        map[string]int32 `json:"by_severity"`
}

type OperationsWorkbench struct {
	GeneratedAt time.Time                  `json:"generated_at"`
	Summary     OperationsWorkbenchSummary `json:"summary"`
	Cards       []*OperationsWorkbenchCard `json:"cards"`
}

type OperationsWorkbenchActionResult struct {
	CardKey string                   `json:"card_key"`
	Action  string                   `json:"action"`
	Status  string                   `json:"status"`
	Badge   string                   `json:"badge,omitempty"`
	Card    *OperationsWorkbenchCard `json:"card,omitempty"`
	Source  json.RawMessage          `json:"source,omitempty"`
}
