package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

type WorkflowTaskEngineRepo interface {
	CreateWorkflowDefinition(ctx context.Context, item *domain.WorkflowDefinition, actorID *uuid.UUID) (*domain.WorkflowDefinition, error)
	UpdateWorkflowDefinition(ctx context.Context, item *domain.WorkflowDefinition, actorID *uuid.UUID) (*domain.WorkflowDefinition, error)
	GetWorkflowDefinition(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.WorkflowDefinition, error)
	ListWorkflowDefinitions(ctx context.Context, tenantID uuid.UUID, status *string, moduleKey *string, search *string, limit int32, offset int32) ([]*domain.WorkflowDefinition, error)
	CreateWorkflowDefinitionStep(ctx context.Context, item *domain.WorkflowDefinitionStep, actorID *uuid.UUID) (*domain.WorkflowDefinitionStep, error)
	UpdateWorkflowDefinitionStep(ctx context.Context, item *domain.WorkflowDefinitionStep, actorID *uuid.UUID) (*domain.WorkflowDefinitionStep, error)
	ListWorkflowDefinitionSteps(ctx context.Context, tenantID uuid.UUID, workflowDefinitionID uuid.UUID) ([]*domain.WorkflowDefinitionStep, error)
	CreateOperationTemplate(ctx context.Context, item *domain.OperationTemplate, actorID *uuid.UUID) (*domain.OperationTemplate, error)
	UpdateOperationTemplate(ctx context.Context, item *domain.OperationTemplate, actorID *uuid.UUID) (*domain.OperationTemplate, error)
	GetOperationTemplate(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.OperationTemplate, error)
	ListOperationTemplates(ctx context.Context, tenantID uuid.UUID, category *string, sourceModule *string, activeOnly *bool, search *string, limit int32, offset int32) ([]*domain.OperationTemplate, error)
	CreateWorkflowTask(ctx context.Context, item *domain.WorkflowTask, actorID *uuid.UUID) (*domain.WorkflowTask, error)
	UpdateWorkflowTask(ctx context.Context, item *domain.WorkflowTask, actorID *uuid.UUID) (*domain.WorkflowTask, error)
	UpdateWorkflowTaskStatus(ctx context.Context, item *domain.WorkflowTask, actorID *uuid.UUID) (*domain.WorkflowTask, error)
	GetWorkflowTask(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.WorkflowTask, error)
	ListWorkflowTasks(ctx context.Context, filter domain.WorkflowTaskFilter) ([]*domain.WorkflowTask, error)
	CreateWorkflowTaskWatcher(ctx context.Context, item *domain.WorkflowTaskWatcher, actorID *uuid.UUID) (*domain.WorkflowTaskWatcher, error)
	RemoveWorkflowTaskWatcher(ctx context.Context, tenantID uuid.UUID, taskID uuid.UUID, watcherUserID uuid.UUID, actorID *uuid.UUID) error
	ListWorkflowTaskWatchers(ctx context.Context, tenantID uuid.UUID, taskID uuid.UUID) ([]*domain.WorkflowTaskWatcher, error)
	CreateWorkflowTaskComment(ctx context.Context, item *domain.WorkflowTaskComment, actorID *uuid.UUID) (*domain.WorkflowTaskComment, error)
	ListWorkflowTaskComments(ctx context.Context, tenantID uuid.UUID, taskID uuid.UUID) ([]*domain.WorkflowTaskComment, error)
	CreateWorkflowTaskAttachment(ctx context.Context, item *domain.WorkflowTaskAttachment, actorID *uuid.UUID) (*domain.WorkflowTaskAttachment, error)
	ListWorkflowTaskAttachments(ctx context.Context, tenantID uuid.UUID, taskID uuid.UUID) ([]*domain.WorkflowTaskAttachment, error)
	CreateWorkflowTaskEvent(ctx context.Context, item *domain.WorkflowTaskEvent, actorID *uuid.UUID) (*domain.WorkflowTaskEvent, error)
	ListWorkflowTaskEvents(ctx context.Context, tenantID uuid.UUID, taskID uuid.UUID) ([]*domain.WorkflowTaskEvent, error)
	GetWorkflowTaskSummary(ctx context.Context, tenantID uuid.UUID) ([]*domain.WorkflowTaskSummaryRow, error)
}

type WorkflowDefinitionCommand struct {
	TenantID        uuid.UUID  `json:"tenant_id"`
	ID              uuid.UUID  `json:"id,omitempty"`
	WorkflowKey     string     `json:"workflow_key"`
	Name            string     `json:"name"`
	ModuleKey       string     `json:"module_key"`
	Description     *string    `json:"description,omitempty"`
	Status          string     `json:"status"`
	VisibilityScope string     `json:"visibility_scope"`
	Metadata        []byte     `json:"metadata,omitempty"`
	ActorID         *uuid.UUID `json:"-"`
}

type WorkflowDefinitionStepCommand struct {
	TenantID             uuid.UUID  `json:"tenant_id"`
	ID                   uuid.UUID  `json:"id,omitempty"`
	WorkflowDefinitionID uuid.UUID  `json:"workflow_definition_id"`
	StepOrder            int32      `json:"step_order"`
	StepKey              string     `json:"step_key"`
	Name                 string     `json:"name"`
	StepType             string     `json:"step_type"`
	AssignmentType       string     `json:"assignment_type"`
	AssignmentValue      *string    `json:"assignment_value,omitempty"`
	Required             bool       `json:"required"`
	DueOffsetHours       int32      `json:"due_offset_hours"`
	AllowedActions       []byte     `json:"allowed_actions,omitempty"`
	Metadata             []byte     `json:"metadata,omitempty"`
	ActorID              *uuid.UUID `json:"-"`
}

type OperationTemplateCommand struct {
	TenantID             uuid.UUID  `json:"tenant_id"`
	ID                   uuid.UUID  `json:"id,omitempty"`
	TemplateKey          string     `json:"template_key"`
	Name                 string     `json:"name"`
	Category             string     `json:"category"`
	SourceModule         string     `json:"source_module"`
	SourceType           string     `json:"source_type"`
	WorkflowDefinitionID *uuid.UUID `json:"workflow_definition_id,omitempty"`
	DefaultPriority      int32      `json:"default_priority"`
	DefaultSeverity      string     `json:"default_severity"`
	AllowedActions       []byte     `json:"allowed_actions,omitempty"`
	LaunchSchema         []byte     `json:"launch_schema,omitempty"`
	IsActive             bool       `json:"is_active"`
	Metadata             []byte     `json:"metadata,omitempty"`
	ActorID              *uuid.UUID `json:"-"`
}

type WorkflowTaskCommand struct {
	TenantID             uuid.UUID  `json:"tenant_id"`
	ID                   uuid.UUID  `json:"id,omitempty"`
	TaskNumber           string     `json:"task_number,omitempty"`
	TemplateID           *uuid.UUID `json:"template_id,omitempty"`
	WorkflowDefinitionID *uuid.UUID `json:"workflow_definition_id,omitempty"`
	WorkflowStepID       *uuid.UUID `json:"workflow_step_id,omitempty"`
	ParentTaskID         *uuid.UUID `json:"parent_task_id,omitempty"`
	SourceModule         string     `json:"source_module"`
	SourceType           string     `json:"source_type"`
	SourceID             *uuid.UUID `json:"source_id,omitempty"`
	SourceRecordLabel    *string    `json:"source_record_label,omitempty"`
	Title                string     `json:"title"`
	Description          *string    `json:"description,omitempty"`
	RequesterUserID      *uuid.UUID `json:"requester_user_id,omitempty"`
	AssigneeUserID       *uuid.UUID `json:"assignee_user_id,omitempty"`
	AssigneeRole         *string    `json:"assignee_role,omitempty"`
	AssigneeTeam         *string    `json:"assignee_team,omitempty"`
	Status               string     `json:"status"`
	Priority             int32      `json:"priority"`
	Severity             string     `json:"severity"`
	VisibilityScope      string     `json:"visibility_scope"`
	DueAt                *string    `json:"due_at,omitempty"`
	ActionSchema         []byte     `json:"action_schema,omitempty"`
	Metadata             []byte     `json:"metadata,omitempty"`
	ActorID              *uuid.UUID `json:"-"`
}

type WorkflowTaskActionCommand struct {
	TenantID       uuid.UUID  `json:"tenant_id"`
	TaskID         uuid.UUID  `json:"task_id"`
	Action         string     `json:"action"`
	Remarks        *string    `json:"remarks,omitempty"`
	AssigneeUserID *uuid.UUID `json:"assignee_user_id,omitempty"`
	AssigneeRole   *string    `json:"assignee_role,omitempty"`
	AssigneeTeam   *string    `json:"assignee_team,omitempty"`
	Metadata       []byte     `json:"metadata,omitempty"`
	ActorID        *uuid.UUID `json:"-"`
}

type WorkflowTaskCommentCommand struct {
	TenantID   uuid.UUID  `json:"tenant_id"`
	TaskID     uuid.UUID  `json:"task_id"`
	Visibility string     `json:"visibility"`
	Body       string     `json:"body"`
	Metadata   []byte     `json:"metadata,omitempty"`
	ActorID    *uuid.UUID `json:"-"`
}

type WorkflowTaskAttachmentCommand struct {
	TenantID          uuid.UUID  `json:"tenant_id"`
	TaskID            uuid.UUID  `json:"task_id"`
	CommentID         *uuid.UUID `json:"comment_id,omitempty"`
	FileName          string     `json:"file_name"`
	ContentType       string     `json:"content_type"`
	FileContentBase64 string     `json:"file_content_base64"`
	Visibility        string     `json:"visibility"`
	Metadata          []byte     `json:"metadata,omitempty"`
	ActorID           *uuid.UUID `json:"-"`
}

type WorkflowTaskWatchCommand struct {
	TenantID      uuid.UUID  `json:"tenant_id"`
	TaskID        uuid.UUID  `json:"task_id"`
	WatcherUserID uuid.UUID  `json:"watcher_user_id"`
	WatchReason   *string    `json:"watch_reason,omitempty"`
	ActorID       *uuid.UUID `json:"-"`
}
