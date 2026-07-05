package services

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

const maxWorkflowTasks = 500

func (s *TenantService) CreateWorkflowDefinition(ctx context.Context, cmd ports.WorkflowDefinitionCommand) (*domain.WorkflowDefinition, error) {
	item := &domain.WorkflowDefinition{TenantID: cmd.TenantID, ID: cmd.ID, WorkflowKey: cmd.WorkflowKey, Name: cmd.Name, ModuleKey: cmd.ModuleKey, Description: cleanOptionalString(cmd.Description), Status: cmd.Status, VisibilityScope: cmd.VisibilityScope, Metadata: json.RawMessage(cmd.Metadata)}
	if err := domain.NormalizeWorkflowDefinition(item); err != nil {
		s.log.Warn().Err(err).Str("operation", "create workflow definition").Str("tenant_id", cmd.TenantID.String()).Msg("invalid workflow definition")
		return nil, err
	}
	return s.workflowTasks.CreateWorkflowDefinition(ctx, item, cmd.ActorID)
}

func (s *TenantService) UpdateWorkflowDefinition(ctx context.Context, cmd ports.WorkflowDefinitionCommand) (*domain.WorkflowDefinition, error) {
	item := &domain.WorkflowDefinition{TenantID: cmd.TenantID, ID: cmd.ID, WorkflowKey: cmd.WorkflowKey, Name: cmd.Name, ModuleKey: cmd.ModuleKey, Description: cleanOptionalString(cmd.Description), Status: cmd.Status, VisibilityScope: cmd.VisibilityScope, Metadata: json.RawMessage(cmd.Metadata)}
	if item.ID == uuid.Nil {
		return nil, domain.ErrInvalidWorkflowDefinition
	}
	if err := domain.NormalizeWorkflowDefinition(item); err != nil {
		s.log.Warn().Err(err).Str("operation", "update workflow definition").Str("tenant_id", cmd.TenantID.String()).Msg("invalid workflow definition")
		return nil, err
	}
	return s.workflowTasks.UpdateWorkflowDefinition(ctx, item, cmd.ActorID)
}

func (s *TenantService) ListWorkflowDefinitions(ctx context.Context, tenantID uuid.UUID, status *string, moduleKey *string, search *string, limit int32, offset int32) ([]*domain.WorkflowDefinition, error) {
	if tenantID == uuid.Nil {
		return nil, domain.ErrInvalidTenantID
	}
	return s.workflowTasks.ListWorkflowDefinitions(ctx, tenantID, cleanOptionalString(status), cleanOptionalString(moduleKey), cleanOptionalString(search), boundedWorkflowLimit(limit), offset)
}

func (s *TenantService) CreateWorkflowDefinitionStep(ctx context.Context, cmd ports.WorkflowDefinitionStepCommand) (*domain.WorkflowDefinitionStep, error) {
	item := &domain.WorkflowDefinitionStep{TenantID: cmd.TenantID, ID: cmd.ID, WorkflowDefinitionID: cmd.WorkflowDefinitionID, StepOrder: cmd.StepOrder, StepKey: cmd.StepKey, Name: cmd.Name, StepType: cmd.StepType, AssignmentType: cmd.AssignmentType, AssignmentValue: cleanOptionalString(cmd.AssignmentValue), Required: cmd.Required, DueOffsetHours: cmd.DueOffsetHours, AllowedActions: json.RawMessage(cmd.AllowedActions), Metadata: json.RawMessage(cmd.Metadata)}
	if err := domain.NormalizeWorkflowStep(item); err != nil {
		s.log.Warn().Err(err).Str("operation", "create workflow definition step").Str("tenant_id", cmd.TenantID.String()).Msg("invalid workflow step")
		return nil, err
	}
	return s.workflowTasks.CreateWorkflowDefinitionStep(ctx, item, cmd.ActorID)
}

func (s *TenantService) UpdateWorkflowDefinitionStep(ctx context.Context, cmd ports.WorkflowDefinitionStepCommand) (*domain.WorkflowDefinitionStep, error) {
	item := &domain.WorkflowDefinitionStep{TenantID: cmd.TenantID, ID: cmd.ID, WorkflowDefinitionID: cmd.WorkflowDefinitionID, StepOrder: cmd.StepOrder, StepKey: cmd.StepKey, Name: cmd.Name, StepType: cmd.StepType, AssignmentType: cmd.AssignmentType, AssignmentValue: cleanOptionalString(cmd.AssignmentValue), Required: cmd.Required, DueOffsetHours: cmd.DueOffsetHours, AllowedActions: json.RawMessage(cmd.AllowedActions), Metadata: json.RawMessage(cmd.Metadata)}
	if item.ID == uuid.Nil {
		return nil, domain.ErrInvalidWorkflowStep
	}
	if err := domain.NormalizeWorkflowStep(item); err != nil {
		s.log.Warn().Err(err).Str("operation", "update workflow definition step").Str("tenant_id", cmd.TenantID.String()).Msg("invalid workflow step")
		return nil, err
	}
	return s.workflowTasks.UpdateWorkflowDefinitionStep(ctx, item, cmd.ActorID)
}

func (s *TenantService) ListWorkflowDefinitionSteps(ctx context.Context, tenantID uuid.UUID, workflowDefinitionID uuid.UUID) ([]*domain.WorkflowDefinitionStep, error) {
	if tenantID == uuid.Nil || workflowDefinitionID == uuid.Nil {
		return nil, domain.ErrInvalidWorkflowStep
	}
	return s.workflowTasks.ListWorkflowDefinitionSteps(ctx, tenantID, workflowDefinitionID)
}

func (s *TenantService) CreateOperationTemplate(ctx context.Context, cmd ports.OperationTemplateCommand) (*domain.OperationTemplate, error) {
	item := operationTemplateFromCommand(cmd)
	if err := domain.NormalizeOperationTemplate(item); err != nil {
		s.log.Warn().Err(err).Str("operation", "create operation template").Str("tenant_id", cmd.TenantID.String()).Msg("invalid operation template")
		return nil, err
	}
	return s.workflowTasks.CreateOperationTemplate(ctx, item, cmd.ActorID)
}

func (s *TenantService) UpdateOperationTemplate(ctx context.Context, cmd ports.OperationTemplateCommand) (*domain.OperationTemplate, error) {
	item := operationTemplateFromCommand(cmd)
	if item.ID == uuid.Nil {
		return nil, domain.ErrInvalidOperationTemplate
	}
	if err := domain.NormalizeOperationTemplate(item); err != nil {
		s.log.Warn().Err(err).Str("operation", "update operation template").Str("tenant_id", cmd.TenantID.String()).Msg("invalid operation template")
		return nil, err
	}
	return s.workflowTasks.UpdateOperationTemplate(ctx, item, cmd.ActorID)
}

func (s *TenantService) ListOperationTemplates(ctx context.Context, tenantID uuid.UUID, category *string, sourceModule *string, activeOnly *bool, search *string, limit int32, offset int32) ([]*domain.OperationTemplate, error) {
	if tenantID == uuid.Nil {
		return nil, domain.ErrInvalidTenantID
	}
	return s.workflowTasks.ListOperationTemplates(ctx, tenantID, cleanOptionalString(category), cleanOptionalString(sourceModule), activeOnly, cleanOptionalString(search), boundedWorkflowLimit(limit), offset)
}

func (s *TenantService) CreateWorkflowTask(ctx context.Context, cmd ports.WorkflowTaskCommand) (*domain.WorkflowTask, error) {
	item, err := s.workflowTaskFromCommand(ctx, cmd)
	if err != nil {
		return nil, err
	}
	if item.TaskNumber == "" {
		item.TaskNumber = generateWorkflowTaskNumber()
	}
	if item.RequesterUserID == nil {
		item.RequesterUserID = cmd.ActorID
	}
	if err := domain.NormalizeWorkflowTask(item); err != nil {
		s.log.Warn().Err(err).Str("operation", "create workflow task").Str("tenant_id", cmd.TenantID.String()).Msg("invalid workflow task")
		return nil, err
	}
	created, err := s.workflowTasks.CreateWorkflowTask(ctx, item, cmd.ActorID)
	if err != nil {
		return nil, err
	}
	_, _ = s.workflowTasks.CreateWorkflowTaskEvent(ctx, &domain.WorkflowTaskEvent{TenantID: created.TenantID, TaskID: created.ID, Action: "created", ToStatus: &created.Status, ActorUserID: cmd.ActorID, Metadata: json.RawMessage(`{}`)}, cmd.ActorID)
	return created, nil
}

func (s *TenantService) UpdateWorkflowTask(ctx context.Context, cmd ports.WorkflowTaskCommand) (*domain.WorkflowTask, error) {
	item, err := s.workflowTaskFromCommand(ctx, cmd)
	if err != nil {
		return nil, err
	}
	if item.ID == uuid.Nil {
		return nil, domain.ErrInvalidWorkflowTask
	}
	if err := domain.NormalizeWorkflowTask(item); err != nil {
		s.log.Warn().Err(err).Str("operation", "update workflow task").Str("tenant_id", cmd.TenantID.String()).Msg("invalid workflow task")
		return nil, err
	}
	updated, err := s.workflowTasks.UpdateWorkflowTask(ctx, item, cmd.ActorID)
	if err != nil {
		return nil, err
	}
	_, _ = s.workflowTasks.CreateWorkflowTaskEvent(ctx, &domain.WorkflowTaskEvent{TenantID: updated.TenantID, TaskID: updated.ID, Action: "updated", ActorUserID: cmd.ActorID, Metadata: json.RawMessage(`{}`)}, cmd.ActorID)
	return updated, nil
}

func (s *TenantService) ListWorkflowTasks(ctx context.Context, filter domain.WorkflowTaskFilter) ([]*domain.WorkflowTask, error) {
	if filter.TenantID == uuid.Nil {
		return nil, domain.ErrInvalidTenantID
	}
	filter.ViewKey = cleanOptionalString(filter.ViewKey)
	filter.Status = cleanOptionalString(filter.Status)
	filter.Severity = cleanOptionalString(filter.Severity)
	filter.SourceModule = cleanOptionalString(filter.SourceModule)
	filter.Search = cleanOptionalString(filter.Search)
	filter.ViewerRole = cleanOptionalString(filter.ViewerRole)
	filter.ViewerTeam = cleanOptionalString(filter.ViewerTeam)
	filter.Limit = boundedWorkflowLimit(filter.Limit)
	return s.workflowTasks.ListWorkflowTasks(ctx, filter)
}

func (s *TenantService) GetWorkflowTaskWorkspace(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.WorkflowTaskWorkspace, error) {
	if tenantID == uuid.Nil || id == uuid.Nil {
		return nil, domain.ErrInvalidWorkflowTask
	}
	task, err := s.workflowTasks.GetWorkflowTask(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}
	comments, err := s.workflowTasks.ListWorkflowTaskComments(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}
	attachments, err := s.workflowTasks.ListWorkflowTaskAttachments(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}
	events, err := s.workflowTasks.ListWorkflowTaskEvents(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}
	watchers, err := s.workflowTasks.ListWorkflowTaskWatchers(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}
	return &domain.WorkflowTaskWorkspace{Task: task, Comments: comments, Attachments: attachments, Events: events, Watchers: watchers}, nil
}

func (s *TenantService) ActWorkflowTask(ctx context.Context, cmd ports.WorkflowTaskActionCommand) (*domain.WorkflowTask, error) {
	task, err := s.workflowTasks.GetWorkflowTask(ctx, cmd.TenantID, cmd.TaskID)
	if err != nil {
		return nil, err
	}
	action := strings.ToLower(strings.TrimSpace(cmd.Action))
	nextStatus := task.Status
	switch action {
	case domain.WorkflowActionApprove:
		nextStatus = domain.WorkflowTaskApproved
	case domain.WorkflowActionReject:
		nextStatus = domain.WorkflowTaskRejected
	case domain.WorkflowActionRequestInfo:
		nextStatus = domain.WorkflowTaskWaitingInfo
	case domain.WorkflowActionComplete:
		nextStatus = domain.WorkflowTaskCompleted
	case domain.WorkflowActionDelegate:
		nextStatus = domain.WorkflowTaskDelegated
		task.AssigneeUserID = cmd.AssigneeUserID
		task.AssigneeRole = cleanOptionalString(cmd.AssigneeRole)
		task.AssigneeTeam = cleanOptionalString(cmd.AssigneeTeam)
		task.DelegatedFromUserID = cmd.ActorID
		if task.AssigneeUserID == nil && task.AssigneeRole == nil && task.AssigneeTeam == nil {
			return nil, domain.ErrInvalidWorkflowAction
		}
	case domain.WorkflowActionOpenRecord:
		nextStatus = task.Status
	default:
		return nil, domain.ErrInvalidWorkflowAction
	}
	fromStatus := task.Status
	task.Status = nextStatus
	updated := task
	if nextStatus != fromStatus || action == domain.WorkflowActionDelegate {
		updated, err = s.workflowTasks.UpdateWorkflowTaskStatus(ctx, task, cmd.ActorID)
		if err != nil {
			return nil, err
		}
	}
	metadata := json.RawMessage(cmd.Metadata)
	if len(metadata) == 0 {
		metadata = json.RawMessage(`{"source_action_recorded":true}`)
	}
	_, err = s.workflowTasks.CreateWorkflowTaskEvent(ctx, &domain.WorkflowTaskEvent{TenantID: cmd.TenantID, TaskID: cmd.TaskID, Action: action, FromStatus: &fromStatus, ToStatus: &nextStatus, ActorUserID: cmd.ActorID, Remarks: cleanOptionalString(cmd.Remarks), Metadata: metadata}, cmd.ActorID)
	return updated, err
}

func (s *TenantService) CreateWorkflowTaskComment(ctx context.Context, cmd ports.WorkflowTaskCommentCommand) (*domain.WorkflowTaskComment, error) {
	body := strings.TrimSpace(cmd.Body)
	if cmd.TenantID == uuid.Nil || cmd.TaskID == uuid.Nil || body == "" {
		return nil, domain.ErrInvalidWorkflowComment
	}
	visibility := strings.TrimSpace(cmd.Visibility)
	if visibility == "" {
		visibility = domain.WorkflowVisibilityTenant
	}
	comment, err := s.workflowTasks.CreateWorkflowTaskComment(ctx, &domain.WorkflowTaskComment{TenantID: cmd.TenantID, TaskID: cmd.TaskID, Visibility: visibility, Body: body, Metadata: json.RawMessage(cmd.Metadata)}, cmd.ActorID)
	if err != nil {
		return nil, err
	}
	_, _ = s.workflowTasks.CreateWorkflowTaskEvent(ctx, &domain.WorkflowTaskEvent{TenantID: cmd.TenantID, TaskID: cmd.TaskID, Action: domain.WorkflowActionComment, ActorUserID: cmd.ActorID, Metadata: json.RawMessage(`{}`)}, cmd.ActorID)
	return comment, nil
}

func (s *TenantService) CreateWorkflowTaskAttachment(ctx context.Context, cmd ports.WorkflowTaskAttachmentCommand) (*domain.WorkflowTaskAttachment, error) {
	if cmd.TenantID == uuid.Nil || cmd.TaskID == uuid.Nil || strings.TrimSpace(cmd.FileName) == "" {
		return nil, domain.ErrInvalidWorkflowAttachment
	}
	content, err := base64.StdEncoding.DecodeString(cmd.FileContentBase64)
	if err != nil || len(content) == 0 {
		return nil, domain.ErrInvalidWorkflowAttachment
	}
	if s.objectStorage == nil {
		return nil, domain.ErrStorageProviderSettingsNotFound
	}
	settings, err := s.resolveWorkflowStorageSettings(ctx, cmd.TenantID)
	if err != nil {
		return nil, err
	}
	entityID := uuid.New()
	contentType := strings.TrimSpace(cmd.ContentType)
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	storagePath, err := s.objectStorage.PutObject(ctx, settings, ports.StoreObjectInput{TenantID: cmd.TenantID, Category: ports.StorageCategoryWorkflowAttachment, OwnerID: cmd.TaskID, EntityID: entityID, FileName: cmd.FileName, ContentType: contentType, Content: content})
	if err != nil {
		return nil, err
	}
	hash := sha256.Sum256(content)
	checksum := hex.EncodeToString(hash[:])
	visibility := strings.TrimSpace(cmd.Visibility)
	if visibility == "" {
		visibility = domain.WorkflowVisibilityTenant
	}
	item := &domain.WorkflowTaskAttachment{TenantID: cmd.TenantID, TaskID: cmd.TaskID, CommentID: cmd.CommentID, FileName: strings.TrimSpace(cmd.FileName), ContentType: contentType, StoragePath: storagePath, ChecksumSHA256: &checksum, SizeBytes: int64(len(content)), Visibility: visibility, Metadata: json.RawMessage(cmd.Metadata)}
	created, err := s.workflowTasks.CreateWorkflowTaskAttachment(ctx, item, cmd.ActorID)
	if err != nil {
		return nil, err
	}
	_, _ = s.workflowTasks.CreateWorkflowTaskEvent(ctx, &domain.WorkflowTaskEvent{TenantID: cmd.TenantID, TaskID: cmd.TaskID, Action: "attachment_added", ActorUserID: cmd.ActorID, Metadata: json.RawMessage(`{}`)}, cmd.ActorID)
	return created, nil
}

func (s *TenantService) WatchWorkflowTask(ctx context.Context, cmd ports.WorkflowTaskWatchCommand) (*domain.WorkflowTaskWatcher, error) {
	if cmd.TenantID == uuid.Nil || cmd.TaskID == uuid.Nil || cmd.WatcherUserID == uuid.Nil {
		return nil, domain.ErrInvalidWorkflowTask
	}
	return s.workflowTasks.CreateWorkflowTaskWatcher(ctx, &domain.WorkflowTaskWatcher{TenantID: cmd.TenantID, TaskID: cmd.TaskID, WatcherUserID: cmd.WatcherUserID, WatchReason: cleanOptionalString(cmd.WatchReason)}, cmd.ActorID)
}

func (s *TenantService) UnwatchWorkflowTask(ctx context.Context, cmd ports.WorkflowTaskWatchCommand) error {
	if cmd.TenantID == uuid.Nil || cmd.TaskID == uuid.Nil || cmd.WatcherUserID == uuid.Nil {
		return domain.ErrInvalidWorkflowTask
	}
	return s.workflowTasks.RemoveWorkflowTaskWatcher(ctx, cmd.TenantID, cmd.TaskID, cmd.WatcherUserID, cmd.ActorID)
}

func (s *TenantService) GetWorkflowTaskSummary(ctx context.Context, tenantID uuid.UUID) ([]*domain.WorkflowTaskSummaryRow, error) {
	if tenantID == uuid.Nil {
		return nil, domain.ErrInvalidTenantID
	}
	return s.workflowTasks.GetWorkflowTaskSummary(ctx, tenantID)
}

func operationTemplateFromCommand(cmd ports.OperationTemplateCommand) *domain.OperationTemplate {
	return &domain.OperationTemplate{TenantID: cmd.TenantID, ID: cmd.ID, TemplateKey: cmd.TemplateKey, Name: cmd.Name, Category: cmd.Category, SourceModule: cmd.SourceModule, SourceType: cmd.SourceType, WorkflowDefinitionID: cmd.WorkflowDefinitionID, DefaultPriority: cmd.DefaultPriority, DefaultSeverity: cmd.DefaultSeverity, AllowedActions: json.RawMessage(cmd.AllowedActions), LaunchSchema: json.RawMessage(cmd.LaunchSchema), IsActive: cmd.IsActive, Metadata: json.RawMessage(cmd.Metadata)}
}

func (s *TenantService) workflowTaskFromCommand(ctx context.Context, cmd ports.WorkflowTaskCommand) (*domain.WorkflowTask, error) {
	dueAt, err := parseOptionalWorkflowTime(cmd.DueAt)
	if err != nil {
		return nil, domain.ErrInvalidWorkflowTask
	}
	item := &domain.WorkflowTask{TenantID: cmd.TenantID, ID: cmd.ID, TaskNumber: strings.TrimSpace(cmd.TaskNumber), TemplateID: cmd.TemplateID, WorkflowDefinitionID: cmd.WorkflowDefinitionID, WorkflowStepID: cmd.WorkflowStepID, ParentTaskID: cmd.ParentTaskID, SourceModule: cmd.SourceModule, SourceType: cmd.SourceType, SourceID: cmd.SourceID, SourceRecordLabel: cleanOptionalString(cmd.SourceRecordLabel), Title: cmd.Title, Description: cleanOptionalString(cmd.Description), RequesterUserID: cmd.RequesterUserID, AssigneeUserID: cmd.AssigneeUserID, AssigneeRole: cleanOptionalString(cmd.AssigneeRole), AssigneeTeam: cleanOptionalString(cmd.AssigneeTeam), Status: cmd.Status, Priority: cmd.Priority, Severity: cmd.Severity, VisibilityScope: cmd.VisibilityScope, DueAt: dueAt, ActionSchema: json.RawMessage(cmd.ActionSchema), Metadata: json.RawMessage(cmd.Metadata)}
	if cmd.TemplateID != nil {
		tpl, err := s.workflowTasks.GetOperationTemplate(ctx, cmd.TenantID, *cmd.TemplateID)
		if err != nil {
			return nil, err
		}
		if item.WorkflowDefinitionID == nil {
			item.WorkflowDefinitionID = tpl.WorkflowDefinitionID
		}
		if item.SourceModule == "" {
			item.SourceModule = tpl.SourceModule
		}
		if item.SourceType == "" {
			item.SourceType = tpl.SourceType
		}
		if item.Priority == 0 {
			item.Priority = tpl.DefaultPriority
		}
		if item.Severity == "" {
			item.Severity = tpl.DefaultSeverity
		}
		if len(item.ActionSchema) == 0 {
			item.ActionSchema = tpl.AllowedActions
		}
	}
	return item, nil
}

func parseOptionalWorkflowTime(value *string) (*time.Time, error) {
	if value == nil || strings.TrimSpace(*value) == "" {
		return nil, nil
	}
	parsed, err := time.Parse(time.RFC3339, strings.TrimSpace(*value))
	if err != nil {
		return nil, err
	}
	utc := parsed.UTC()
	return &utc, nil
}

func generateWorkflowTaskNumber() string {
	return "WF-" + time.Now().UTC().Format("20060102-150405") + "-" + strings.ToUpper(uuid.NewString()[:8])
}

func boundedWorkflowLimit(limit int32) int32 {
	if limit <= 0 {
		return 50
	}
	if limit > maxWorkflowTasks {
		return maxWorkflowTasks
	}
	return limit
}

func (s *TenantService) resolveWorkflowStorageSettings(ctx context.Context, tenantID uuid.UUID) (*domain.StorageProviderSettings, error) {
	if s.storageProviders != nil {
		settings, err := s.storageProviders.GetStorageProviderSettings(ctx, tenantID)
		if err == nil {
			return settings, nil
		}
		if err != domain.ErrStorageProviderSettingsNotFound {
			return nil, err
		}
	}
	if s.defaultStorageProvider == nil {
		return nil, domain.ErrStorageProviderSettingsNotFound
	}
	copy := *s.defaultStorageProvider
	copy.TenantID = tenantID
	return &copy, nil
}
