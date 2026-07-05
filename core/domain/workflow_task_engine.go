package domain

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	WorkflowDefinitionDraft    = "draft"
	WorkflowDefinitionActive   = "active"
	WorkflowDefinitionPaused   = "paused"
	WorkflowDefinitionArchived = "archived"

	WorkflowVisibilityTenant       = "tenant"
	WorkflowVisibilityRestricted   = "restricted"
	WorkflowVisibilityConfidential = "confidential"

	WorkflowStepApproval     = "approval"
	WorkflowStepReview       = "review"
	WorkflowStepAction       = "action"
	WorkflowStepChecklist    = "checklist"
	WorkflowStepNotification = "notification"

	WorkflowAssignmentUser        = "user"
	WorkflowAssignmentRole        = "role"
	WorkflowAssignmentTeam        = "team"
	WorkflowAssignmentRequester   = "requester"
	WorkflowAssignmentManager     = "manager"
	WorkflowAssignmentSourceOwner = "source_owner"

	WorkflowTaskPending     = "pending"
	WorkflowTaskInProgress  = "in_progress"
	WorkflowTaskWaitingInfo = "waiting_info"
	WorkflowTaskApproved    = "approved"
	WorkflowTaskRejected    = "rejected"
	WorkflowTaskCompleted   = "completed"
	WorkflowTaskCancelled   = "cancelled"
	WorkflowTaskDelegated   = "delegated"
	WorkflowTaskBlocked     = "blocked"

	WorkflowSeverityLow      = "low"
	WorkflowSeverityMedium   = "medium"
	WorkflowSeverityHigh     = "high"
	WorkflowSeverityCritical = "critical"

	WorkflowActionApprove     = "approve"
	WorkflowActionReject      = "reject"
	WorkflowActionRequestInfo = "request_info"
	WorkflowActionDelegate    = "delegate"
	WorkflowActionComment     = "comment"
	WorkflowActionComplete    = "complete"
	WorkflowActionOpenRecord  = "open_record"
	WorkflowActionWatch       = "watch"
	WorkflowActionUnwatch     = "unwatch"
)

var (
	ErrInvalidWorkflowDefinition  = errors.New("workflow definition is invalid")
	ErrWorkflowDefinitionNotFound = errors.New("workflow definition not found")
	ErrInvalidWorkflowStep        = errors.New("workflow step is invalid")
	ErrInvalidOperationTemplate   = errors.New("operation template is invalid")
	ErrOperationTemplateNotFound  = errors.New("operation template not found")
	ErrInvalidWorkflowTask        = errors.New("workflow task is invalid")
	ErrWorkflowTaskNotFound       = errors.New("workflow task not found")
	ErrInvalidWorkflowComment     = errors.New("workflow task comment is invalid")
	ErrInvalidWorkflowAttachment  = errors.New("workflow task attachment is invalid")
	ErrInvalidWorkflowAction      = errors.New("workflow task action is invalid")
)

type WorkflowDefinition struct {
	ID              uuid.UUID       `json:"id"`
	TenantID        uuid.UUID       `json:"tenant_id"`
	WorkflowKey     string          `json:"workflow_key"`
	Name            string          `json:"name"`
	ModuleKey       string          `json:"module_key"`
	Description     *string         `json:"description,omitempty"`
	Status          string          `json:"status"`
	VisibilityScope string          `json:"visibility_scope"`
	Metadata        json.RawMessage `json:"metadata"`
	Inactive        bool            `json:"inactive"`
	CreatedAt       time.Time       `json:"created_at"`
	CreatedBy       *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt       time.Time       `json:"updated_at"`
	UpdatedBy       *uuid.UUID      `json:"updated_by,omitempty"`
}

type WorkflowDefinitionStep struct {
	ID                   uuid.UUID       `json:"id"`
	TenantID             uuid.UUID       `json:"tenant_id"`
	WorkflowDefinitionID uuid.UUID       `json:"workflow_definition_id"`
	StepOrder            int32           `json:"step_order"`
	StepKey              string          `json:"step_key"`
	Name                 string          `json:"name"`
	StepType             string          `json:"step_type"`
	AssignmentType       string          `json:"assignment_type"`
	AssignmentValue      *string         `json:"assignment_value,omitempty"`
	Required             bool            `json:"required"`
	DueOffsetHours       int32           `json:"due_offset_hours"`
	AllowedActions       json.RawMessage `json:"allowed_actions"`
	Metadata             json.RawMessage `json:"metadata"`
	Inactive             bool            `json:"inactive"`
	CreatedAt            time.Time       `json:"created_at"`
	CreatedBy            *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt            time.Time       `json:"updated_at"`
	UpdatedBy            *uuid.UUID      `json:"updated_by,omitempty"`
}

type OperationTemplate struct {
	ID                   uuid.UUID       `json:"id"`
	TenantID             uuid.UUID       `json:"tenant_id"`
	TemplateKey          string          `json:"template_key"`
	Name                 string          `json:"name"`
	Category             string          `json:"category"`
	SourceModule         string          `json:"source_module"`
	SourceType           string          `json:"source_type"`
	WorkflowDefinitionID *uuid.UUID      `json:"workflow_definition_id,omitempty"`
	DefaultPriority      int32           `json:"default_priority"`
	DefaultSeverity      string          `json:"default_severity"`
	AllowedActions       json.RawMessage `json:"allowed_actions"`
	LaunchSchema         json.RawMessage `json:"launch_schema"`
	IsActive             bool            `json:"is_active"`
	Metadata             json.RawMessage `json:"metadata"`
	Inactive             bool            `json:"inactive"`
	CreatedAt            time.Time       `json:"created_at"`
	CreatedBy            *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt            time.Time       `json:"updated_at"`
	UpdatedBy            *uuid.UUID      `json:"updated_by,omitempty"`
}

type WorkflowTask struct {
	ID                   uuid.UUID       `json:"id"`
	TenantID             uuid.UUID       `json:"tenant_id"`
	TaskNumber           string          `json:"task_number"`
	TemplateID           *uuid.UUID      `json:"template_id,omitempty"`
	WorkflowDefinitionID *uuid.UUID      `json:"workflow_definition_id,omitempty"`
	WorkflowStepID       *uuid.UUID      `json:"workflow_step_id,omitempty"`
	ParentTaskID         *uuid.UUID      `json:"parent_task_id,omitempty"`
	SourceModule         string          `json:"source_module"`
	SourceType           string          `json:"source_type"`
	SourceID             *uuid.UUID      `json:"source_id,omitempty"`
	SourceRecordLabel    *string         `json:"source_record_label,omitempty"`
	Title                string          `json:"title"`
	Description          *string         `json:"description,omitempty"`
	RequesterUserID      *uuid.UUID      `json:"requester_user_id,omitempty"`
	AssigneeUserID       *uuid.UUID      `json:"assignee_user_id,omitempty"`
	AssigneeRole         *string         `json:"assignee_role,omitempty"`
	AssigneeTeam         *string         `json:"assignee_team,omitempty"`
	DelegatedFromUserID  *uuid.UUID      `json:"delegated_from_user_id,omitempty"`
	Status               string          `json:"status"`
	Priority             int32           `json:"priority"`
	Severity             string          `json:"severity"`
	VisibilityScope      string          `json:"visibility_scope"`
	DueAt                *time.Time      `json:"due_at,omitempty"`
	CompletedAt          *time.Time      `json:"completed_at,omitempty"`
	CompletedBy          *uuid.UUID      `json:"completed_by,omitempty"`
	ActionSchema         json.RawMessage `json:"action_schema"`
	Metadata             json.RawMessage `json:"metadata"`
	Inactive             bool            `json:"inactive"`
	CreatedAt            time.Time       `json:"created_at"`
	CreatedBy            *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt            time.Time       `json:"updated_at"`
	UpdatedBy            *uuid.UUID      `json:"updated_by,omitempty"`
	CommentCount         int32           `json:"comment_count,omitempty"`
	AttachmentCount      int32           `json:"attachment_count,omitempty"`
	WatchedByViewer      bool            `json:"watched_by_viewer,omitempty"`
}

type WorkflowTaskWatcher struct {
	ID            uuid.UUID  `json:"id"`
	TenantID      uuid.UUID  `json:"tenant_id"`
	TaskID        uuid.UUID  `json:"task_id"`
	WatcherUserID uuid.UUID  `json:"watcher_user_id"`
	WatchReason   *string    `json:"watch_reason,omitempty"`
	Inactive      bool       `json:"inactive"`
	CreatedAt     time.Time  `json:"created_at"`
	CreatedBy     *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt     time.Time  `json:"updated_at"`
	UpdatedBy     *uuid.UUID `json:"updated_by,omitempty"`
}

type WorkflowTaskComment struct {
	ID         uuid.UUID       `json:"id"`
	TenantID   uuid.UUID       `json:"tenant_id"`
	TaskID     uuid.UUID       `json:"task_id"`
	Visibility string          `json:"visibility"`
	Body       string          `json:"body"`
	Metadata   json.RawMessage `json:"metadata"`
	Inactive   bool            `json:"inactive"`
	CreatedAt  time.Time       `json:"created_at"`
	CreatedBy  *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt  time.Time       `json:"updated_at"`
	UpdatedBy  *uuid.UUID      `json:"updated_by,omitempty"`
}

type WorkflowTaskAttachment struct {
	ID             uuid.UUID       `json:"id"`
	TenantID       uuid.UUID       `json:"tenant_id"`
	TaskID         uuid.UUID       `json:"task_id"`
	CommentID      *uuid.UUID      `json:"comment_id,omitempty"`
	FileName       string          `json:"file_name"`
	ContentType    string          `json:"content_type"`
	StoragePath    string          `json:"storage_path"`
	ChecksumSHA256 *string         `json:"checksum_sha256,omitempty"`
	SizeBytes      int64           `json:"size_bytes"`
	Visibility     string          `json:"visibility"`
	Metadata       json.RawMessage `json:"metadata"`
	Inactive       bool            `json:"inactive"`
	CreatedAt      time.Time       `json:"created_at"`
	CreatedBy      *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt      time.Time       `json:"updated_at"`
	UpdatedBy      *uuid.UUID      `json:"updated_by,omitempty"`
}

type WorkflowTaskEvent struct {
	ID          uuid.UUID       `json:"id"`
	TenantID    uuid.UUID       `json:"tenant_id"`
	TaskID      uuid.UUID       `json:"task_id"`
	Action      string          `json:"action"`
	FromStatus  *string         `json:"from_status,omitempty"`
	ToStatus    *string         `json:"to_status,omitempty"`
	ActorUserID *uuid.UUID      `json:"actor_user_id,omitempty"`
	Remarks     *string         `json:"remarks,omitempty"`
	Metadata    json.RawMessage `json:"metadata"`
	Inactive    bool            `json:"inactive"`
	CreatedAt   time.Time       `json:"created_at"`
	CreatedBy   *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt   time.Time       `json:"updated_at"`
	UpdatedBy   *uuid.UUID      `json:"updated_by,omitempty"`
}

type WorkflowTaskSummaryRow struct {
	Metric      string `json:"metric"`
	MetricCount int64  `json:"metric_count"`
}

type WorkflowTaskWorkspace struct {
	Task        *WorkflowTask             `json:"task"`
	Comments    []*WorkflowTaskComment    `json:"comments"`
	Attachments []*WorkflowTaskAttachment `json:"attachments"`
	Events      []*WorkflowTaskEvent      `json:"events"`
	Watchers    []*WorkflowTaskWatcher    `json:"watchers"`
}

type WorkflowTaskFilter struct {
	TenantID     uuid.UUID
	ViewKey      *string
	Status       *string
	Severity     *string
	SourceModule *string
	Search       *string
	ViewerUserID *uuid.UUID
	ViewerRole   *string
	ViewerTeam   *string
	Limit        int32
	Offset       int32
}

func NormalizeWorkflowDefinition(item *WorkflowDefinition) error {
	if item == nil || item.TenantID == uuid.Nil || strings.TrimSpace(item.WorkflowKey) == "" || strings.TrimSpace(item.Name) == "" || strings.TrimSpace(item.ModuleKey) == "" {
		return ErrInvalidWorkflowDefinition
	}
	item.WorkflowKey = strings.ToLower(strings.TrimSpace(item.WorkflowKey))
	item.Name = strings.TrimSpace(item.Name)
	item.ModuleKey = strings.TrimSpace(item.ModuleKey)
	item.Status = normalizeWorkflowAllowed(item.Status, WorkflowDefinitionDraft, WorkflowDefinitionDraft, WorkflowDefinitionActive, WorkflowDefinitionPaused, WorkflowDefinitionArchived)
	item.VisibilityScope = normalizeVisibility(item.VisibilityScope)
	item.Metadata = ensureJSONObject(item.Metadata)
	return nil
}

func NormalizeWorkflowStep(item *WorkflowDefinitionStep) error {
	if item == nil || item.TenantID == uuid.Nil || item.WorkflowDefinitionID == uuid.Nil || strings.TrimSpace(item.StepKey) == "" || strings.TrimSpace(item.Name) == "" {
		return ErrInvalidWorkflowStep
	}
	item.StepKey = strings.ToLower(strings.TrimSpace(item.StepKey))
	item.Name = strings.TrimSpace(item.Name)
	item.StepType = normalizeWorkflowAllowed(item.StepType, WorkflowStepApproval, WorkflowStepApproval, WorkflowStepReview, WorkflowStepAction, WorkflowStepChecklist, WorkflowStepNotification)
	item.AssignmentType = normalizeWorkflowAllowed(item.AssignmentType, WorkflowAssignmentRole, WorkflowAssignmentUser, WorkflowAssignmentRole, WorkflowAssignmentTeam, WorkflowAssignmentRequester, WorkflowAssignmentManager, WorkflowAssignmentSourceOwner)
	item.AllowedActions = ensureJSONArray(item.AllowedActions)
	item.Metadata = ensureJSONObject(item.Metadata)
	if item.StepOrder <= 0 {
		item.StepOrder = 1
	}
	return nil
}

func NormalizeOperationTemplate(item *OperationTemplate) error {
	if item == nil || item.TenantID == uuid.Nil || strings.TrimSpace(item.TemplateKey) == "" || strings.TrimSpace(item.Name) == "" || strings.TrimSpace(item.SourceModule) == "" || strings.TrimSpace(item.SourceType) == "" {
		return ErrInvalidOperationTemplate
	}
	item.TemplateKey = strings.ToLower(strings.TrimSpace(item.TemplateKey))
	item.Name = strings.TrimSpace(item.Name)
	item.Category = strings.TrimSpace(item.Category)
	if item.Category == "" {
		item.Category = "general"
	}
	item.SourceModule = strings.TrimSpace(item.SourceModule)
	item.SourceType = strings.TrimSpace(item.SourceType)
	item.DefaultSeverity = normalizeSeverity(item.DefaultSeverity)
	if item.DefaultPriority <= 0 || item.DefaultPriority > 100 {
		item.DefaultPriority = 50
	}
	item.AllowedActions = ensureJSONArray(item.AllowedActions)
	item.LaunchSchema = ensureJSONObject(item.LaunchSchema)
	item.Metadata = ensureJSONObject(item.Metadata)
	return nil
}

func NormalizeWorkflowTask(item *WorkflowTask) error {
	if item == nil || item.TenantID == uuid.Nil || strings.TrimSpace(item.Title) == "" {
		return ErrInvalidWorkflowTask
	}
	item.Title = strings.TrimSpace(item.Title)
	item.SourceModule = strings.TrimSpace(item.SourceModule)
	if item.SourceModule == "" {
		item.SourceModule = "hrms"
	}
	item.SourceType = strings.TrimSpace(item.SourceType)
	if item.SourceType == "" {
		item.SourceType = "manual"
	}
	item.Status = normalizeTaskStatus(item.Status)
	item.Severity = normalizeSeverity(item.Severity)
	item.VisibilityScope = normalizeVisibility(item.VisibilityScope)
	if item.Priority <= 0 || item.Priority > 100 {
		item.Priority = 50
	}
	item.ActionSchema = ensureJSONArray(item.ActionSchema)
	item.Metadata = ensureJSONObject(item.Metadata)
	return nil
}

func IsWorkflowTerminalStatus(status string) bool {
	switch status {
	case WorkflowTaskApproved, WorkflowTaskRejected, WorkflowTaskCompleted, WorkflowTaskCancelled:
		return true
	default:
		return false
	}
}

func normalizeTaskStatus(status string) string {
	return normalizeWorkflowAllowed(status, WorkflowTaskPending, WorkflowTaskPending, WorkflowTaskInProgress, WorkflowTaskWaitingInfo, WorkflowTaskApproved, WorkflowTaskRejected, WorkflowTaskCompleted, WorkflowTaskCancelled, WorkflowTaskDelegated, WorkflowTaskBlocked)
}

func normalizeSeverity(value string) string {
	return normalizeWorkflowAllowed(value, WorkflowSeverityMedium, WorkflowSeverityLow, WorkflowSeverityMedium, WorkflowSeverityHigh, WorkflowSeverityCritical)
}

func normalizeVisibility(value string) string {
	return normalizeWorkflowAllowed(value, WorkflowVisibilityTenant, WorkflowVisibilityTenant, WorkflowVisibilityRestricted, WorkflowVisibilityConfidential)
}

func normalizeWorkflowAllowed(value string, fallback string, allowed ...string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	for _, candidate := range allowed {
		if value == candidate {
			return value
		}
	}
	return fallback
}

func ensureJSONObject(value json.RawMessage) json.RawMessage {
	if len(value) == 0 || !json.Valid(value) {
		return json.RawMessage(`{}`)
	}
	return value
}

func ensureJSONArray(value json.RawMessage) json.RawMessage {
	if len(value) == 0 || !json.Valid(value) {
		return json.RawMessage(`[]`)
	}
	return value
}
