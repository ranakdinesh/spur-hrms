package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func (s *Store) CreateWorkflowDefinition(ctx context.Context, item *domain.WorkflowDefinition, actorID *uuid.UUID) (*domain.WorkflowDefinition, error) {
	row, err := s.getQueries(ctx).CreateWorkflowDefinition(ctx, sqlc.CreateWorkflowDefinitionParams{TenantID: item.TenantID, WorkflowKey: item.WorkflowKey, Name: item.Name, ModuleKey: item.ModuleKey, Description: textFromPtr(item.Description), Status: item.Status, VisibilityScope: item.VisibilityScope, Metadata: jsonBytesFromRaw(item.Metadata), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create workflow definition", err, tenantIDField(item.TenantID), stringField("workflow_key", item.WorkflowKey))
	}
	return mapWorkflowDefinition(row), nil
}

func (s *Store) UpdateWorkflowDefinition(ctx context.Context, item *domain.WorkflowDefinition, actorID *uuid.UUID) (*domain.WorkflowDefinition, error) {
	row, err := s.getQueries(ctx).UpdateWorkflowDefinition(ctx, sqlc.UpdateWorkflowDefinitionParams{TenantID: item.TenantID, ID: item.ID, WorkflowKey: item.WorkflowKey, Name: item.Name, ModuleKey: item.ModuleKey, Description: textFromPtr(item.Description), Status: item.Status, VisibilityScope: item.VisibilityScope, Metadata: jsonBytesFromRaw(item.Metadata), UpdatedBy: uuidFromPtr(actorID)})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrWorkflowDefinitionNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "update workflow definition", err, tenantIDField(item.TenantID), stringField("workflow_definition_id", item.ID.String()))
	}
	return mapWorkflowDefinition(row), nil
}

func (s *Store) GetWorkflowDefinition(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.WorkflowDefinition, error) {
	row, err := s.getQueries(ctx).GetWorkflowDefinition(ctx, sqlc.GetWorkflowDefinitionParams{TenantID: tenantID, ID: id})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrWorkflowDefinitionNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "get workflow definition", err, tenantIDField(tenantID), stringField("workflow_definition_id", id.String()))
	}
	return mapWorkflowDefinition(row), nil
}

func (s *Store) ListWorkflowDefinitions(ctx context.Context, tenantID uuid.UUID, status *string, moduleKey *string, search *string, limit int32, offset int32) ([]*domain.WorkflowDefinition, error) {
	rows, err := s.getQueries(ctx).ListWorkflowDefinitions(ctx, sqlc.ListWorkflowDefinitionsParams{TenantID: tenantID, Status: textFromPtr(status), ModuleKey: textFromPtr(moduleKey), Search: textFromPtr(search), Limit: limitOrDefault(limit), Offset: offset})
	if err != nil {
		return nil, s.logDBError(ctx, "list workflow definitions", err, tenantIDField(tenantID))
	}
	return mapWorkflowDefinitions(rows), nil
}

func (s *Store) CreateWorkflowDefinitionStep(ctx context.Context, item *domain.WorkflowDefinitionStep, actorID *uuid.UUID) (*domain.WorkflowDefinitionStep, error) {
	row, err := s.getQueries(ctx).CreateWorkflowDefinitionStep(ctx, sqlc.CreateWorkflowDefinitionStepParams{TenantID: item.TenantID, WorkflowDefinitionID: item.WorkflowDefinitionID, StepOrder: item.StepOrder, StepKey: item.StepKey, Name: item.Name, StepType: item.StepType, AssignmentType: item.AssignmentType, AssignmentValue: textFromPtr(item.AssignmentValue), Required: item.Required, DueOffsetHours: item.DueOffsetHours, AllowedActions: jsonBytesFromRaw(item.AllowedActions), Metadata: jsonBytesFromRaw(item.Metadata), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create workflow definition step", err, tenantIDField(item.TenantID), stringField("workflow_definition_id", item.WorkflowDefinitionID.String()))
	}
	return mapWorkflowDefinitionStep(row), nil
}

func (s *Store) UpdateWorkflowDefinitionStep(ctx context.Context, item *domain.WorkflowDefinitionStep, actorID *uuid.UUID) (*domain.WorkflowDefinitionStep, error) {
	row, err := s.getQueries(ctx).UpdateWorkflowDefinitionStep(ctx, sqlc.UpdateWorkflowDefinitionStepParams{TenantID: item.TenantID, ID: item.ID, StepOrder: item.StepOrder, StepKey: item.StepKey, Name: item.Name, StepType: item.StepType, AssignmentType: item.AssignmentType, AssignmentValue: textFromPtr(item.AssignmentValue), Required: item.Required, DueOffsetHours: item.DueOffsetHours, AllowedActions: jsonBytesFromRaw(item.AllowedActions), Metadata: jsonBytesFromRaw(item.Metadata), UpdatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "update workflow definition step", err, tenantIDField(item.TenantID), stringField("workflow_step_id", item.ID.String()))
	}
	return mapWorkflowDefinitionStep(row), nil
}

func (s *Store) ListWorkflowDefinitionSteps(ctx context.Context, tenantID uuid.UUID, workflowDefinitionID uuid.UUID) ([]*domain.WorkflowDefinitionStep, error) {
	rows, err := s.getQueries(ctx).ListWorkflowDefinitionSteps(ctx, sqlc.ListWorkflowDefinitionStepsParams{TenantID: tenantID, WorkflowDefinitionID: workflowDefinitionID})
	if err != nil {
		return nil, s.logDBError(ctx, "list workflow definition steps", err, tenantIDField(tenantID), stringField("workflow_definition_id", workflowDefinitionID.String()))
	}
	return mapWorkflowDefinitionSteps(rows), nil
}

func (s *Store) CreateOperationTemplate(ctx context.Context, item *domain.OperationTemplate, actorID *uuid.UUID) (*domain.OperationTemplate, error) {
	row, err := s.getQueries(ctx).CreateOperationTemplate(ctx, sqlc.CreateOperationTemplateParams{TenantID: item.TenantID, TemplateKey: item.TemplateKey, Name: item.Name, Category: item.Category, SourceModule: item.SourceModule, SourceType: item.SourceType, WorkflowDefinitionID: uuidFromPtr(item.WorkflowDefinitionID), DefaultPriority: item.DefaultPriority, DefaultSeverity: item.DefaultSeverity, AllowedActions: jsonBytesFromRaw(item.AllowedActions), LaunchSchema: jsonBytesFromRaw(item.LaunchSchema), IsActive: item.IsActive, Metadata: jsonBytesFromRaw(item.Metadata), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create operation template", err, tenantIDField(item.TenantID), stringField("template_key", item.TemplateKey))
	}
	return mapOperationTemplate(row), nil
}

func (s *Store) UpdateOperationTemplate(ctx context.Context, item *domain.OperationTemplate, actorID *uuid.UUID) (*domain.OperationTemplate, error) {
	row, err := s.getQueries(ctx).UpdateOperationTemplate(ctx, sqlc.UpdateOperationTemplateParams{TenantID: item.TenantID, ID: item.ID, TemplateKey: item.TemplateKey, Name: item.Name, Category: item.Category, SourceModule: item.SourceModule, SourceType: item.SourceType, WorkflowDefinitionID: uuidFromPtr(item.WorkflowDefinitionID), DefaultPriority: item.DefaultPriority, DefaultSeverity: item.DefaultSeverity, AllowedActions: jsonBytesFromRaw(item.AllowedActions), LaunchSchema: jsonBytesFromRaw(item.LaunchSchema), IsActive: item.IsActive, Metadata: jsonBytesFromRaw(item.Metadata), UpdatedBy: uuidFromPtr(actorID)})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrOperationTemplateNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "update operation template", err, tenantIDField(item.TenantID), stringField("operation_template_id", item.ID.String()))
	}
	return mapOperationTemplate(row), nil
}

func (s *Store) GetOperationTemplate(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.OperationTemplate, error) {
	row, err := s.getQueries(ctx).GetOperationTemplate(ctx, sqlc.GetOperationTemplateParams{TenantID: tenantID, ID: id})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrOperationTemplateNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "get operation template", err, tenantIDField(tenantID), stringField("operation_template_id", id.String()))
	}
	return mapOperationTemplate(row), nil
}

func (s *Store) ListOperationTemplates(ctx context.Context, tenantID uuid.UUID, category *string, sourceModule *string, activeOnly *bool, search *string, limit int32, offset int32) ([]*domain.OperationTemplate, error) {
	rows, err := s.getQueries(ctx).ListOperationTemplates(ctx, sqlc.ListOperationTemplatesParams{TenantID: tenantID, Category: textFromPtr(category), SourceModule: textFromPtr(sourceModule), ActiveOnly: boolFromPtr(activeOnly), Search: textFromPtr(search), Limit: limitOrDefault(limit), Offset: offset})
	if err != nil {
		return nil, s.logDBError(ctx, "list operation templates", err, tenantIDField(tenantID))
	}
	return mapOperationTemplates(rows), nil
}

func (s *Store) CreateWorkflowTask(ctx context.Context, item *domain.WorkflowTask, actorID *uuid.UUID) (*domain.WorkflowTask, error) {
	row, err := s.getQueries(ctx).CreateWorkflowTask(ctx, sqlc.CreateWorkflowTaskParams{TenantID: item.TenantID, TaskNumber: item.TaskNumber, TemplateID: uuidFromPtr(item.TemplateID), WorkflowDefinitionID: uuidFromPtr(item.WorkflowDefinitionID), WorkflowStepID: uuidFromPtr(item.WorkflowStepID), ParentTaskID: uuidFromPtr(item.ParentTaskID), SourceModule: item.SourceModule, SourceType: item.SourceType, SourceID: uuidFromPtr(item.SourceID), SourceRecordLabel: textFromPtr(item.SourceRecordLabel), Title: item.Title, Description: textFromPtr(item.Description), RequesterUserID: uuidFromPtr(item.RequesterUserID), AssigneeUserID: uuidFromPtr(item.AssigneeUserID), AssigneeRole: textFromPtr(item.AssigneeRole), AssigneeTeam: textFromPtr(item.AssigneeTeam), DelegatedFromUserID: uuidFromPtr(item.DelegatedFromUserID), Status: item.Status, Priority: item.Priority, Severity: item.Severity, VisibilityScope: item.VisibilityScope, DueAt: timestamptzFromPtr(item.DueAt), ActionSchema: jsonBytesFromRaw(item.ActionSchema), Metadata: jsonBytesFromRaw(item.Metadata), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create workflow task", err, tenantIDField(item.TenantID), stringField("task_number", item.TaskNumber))
	}
	return mapWorkflowTask(row), nil
}

func (s *Store) UpdateWorkflowTask(ctx context.Context, item *domain.WorkflowTask, actorID *uuid.UUID) (*domain.WorkflowTask, error) {
	row, err := s.getQueries(ctx).UpdateWorkflowTask(ctx, sqlc.UpdateWorkflowTaskParams{TenantID: item.TenantID, ID: item.ID, Title: item.Title, Description: textFromPtr(item.Description), AssigneeUserID: uuidFromPtr(item.AssigneeUserID), AssigneeRole: textFromPtr(item.AssigneeRole), AssigneeTeam: textFromPtr(item.AssigneeTeam), Priority: item.Priority, Severity: item.Severity, VisibilityScope: item.VisibilityScope, DueAt: timestamptzFromPtr(item.DueAt), ActionSchema: jsonBytesFromRaw(item.ActionSchema), Metadata: jsonBytesFromRaw(item.Metadata), UpdatedBy: uuidFromPtr(actorID)})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrWorkflowTaskNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "update workflow task", err, tenantIDField(item.TenantID), stringField("workflow_task_id", item.ID.String()))
	}
	return mapWorkflowTask(row), nil
}

func (s *Store) UpdateWorkflowTaskStatus(ctx context.Context, item *domain.WorkflowTask, actorID *uuid.UUID) (*domain.WorkflowTask, error) {
	row, err := s.getQueries(ctx).UpdateWorkflowTaskStatus(ctx, sqlc.UpdateWorkflowTaskStatusParams{TenantID: item.TenantID, ID: item.ID, Status: item.Status, AssigneeUserID: uuidFromPtr(item.AssigneeUserID), AssigneeRole: textFromPtr(item.AssigneeRole), AssigneeTeam: textFromPtr(item.AssigneeTeam), DelegatedFromUserID: uuidFromPtr(item.DelegatedFromUserID), UpdatedBy: uuidFromPtr(actorID)})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrWorkflowTaskNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "update workflow task status", err, tenantIDField(item.TenantID), stringField("workflow_task_id", item.ID.String()), stringField("status", item.Status))
	}
	return mapWorkflowTask(row), nil
}

func (s *Store) GetWorkflowTask(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.WorkflowTask, error) {
	row, err := s.getQueries(ctx).GetWorkflowTask(ctx, sqlc.GetWorkflowTaskParams{TenantID: tenantID, ID: id})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrWorkflowTaskNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "get workflow task", err, tenantIDField(tenantID), stringField("workflow_task_id", id.String()))
	}
	return mapWorkflowTask(row), nil
}

func (s *Store) ListWorkflowTasks(ctx context.Context, filter domain.WorkflowTaskFilter) ([]*domain.WorkflowTask, error) {
	rows, err := s.getQueries(ctx).ListWorkflowTasks(ctx, sqlc.ListWorkflowTasksParams{TenantID: filter.TenantID, ViewerUserID: uuidFromPtr(filter.ViewerUserID), Status: textFromPtr(filter.Status), Severity: textFromPtr(filter.Severity), SourceModule: textFromPtr(filter.SourceModule), Search: textFromPtr(filter.Search), ViewKey: textFromPtr(filter.ViewKey), ViewerRole: textFromPtr(filter.ViewerRole), ViewerTeam: textFromPtr(filter.ViewerTeam), Limit: limitOrDefault(filter.Limit), Offset: filter.Offset})
	if err != nil {
		return nil, s.logDBError(ctx, "list workflow tasks", err, tenantIDField(filter.TenantID))
	}
	return mapWorkflowTasks(rows), nil
}

func (s *Store) CreateWorkflowTaskWatcher(ctx context.Context, item *domain.WorkflowTaskWatcher, actorID *uuid.UUID) (*domain.WorkflowTaskWatcher, error) {
	row, err := s.getQueries(ctx).CreateWorkflowTaskWatcher(ctx, sqlc.CreateWorkflowTaskWatcherParams{TenantID: item.TenantID, TaskID: item.TaskID, WatcherUserID: item.WatcherUserID, WatchReason: textFromPtr(item.WatchReason), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create workflow task watcher", err, tenantIDField(item.TenantID), stringField("workflow_task_id", item.TaskID.String()))
	}
	return mapWorkflowTaskWatcher(row), nil
}

func (s *Store) RemoveWorkflowTaskWatcher(ctx context.Context, tenantID uuid.UUID, taskID uuid.UUID, watcherUserID uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).RemoveWorkflowTaskWatcher(ctx, sqlc.RemoveWorkflowTaskWatcherParams{TenantID: tenantID, TaskID: taskID, WatcherUserID: watcherUserID, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "remove workflow task watcher", err, tenantIDField(tenantID), stringField("workflow_task_id", taskID.String()))
	}
	return nil
}

func (s *Store) ListWorkflowTaskWatchers(ctx context.Context, tenantID uuid.UUID, taskID uuid.UUID) ([]*domain.WorkflowTaskWatcher, error) {
	rows, err := s.getQueries(ctx).ListWorkflowTaskWatchers(ctx, sqlc.ListWorkflowTaskWatchersParams{TenantID: tenantID, TaskID: taskID})
	if err != nil {
		return nil, s.logDBError(ctx, "list workflow task watchers", err, tenantIDField(tenantID), stringField("workflow_task_id", taskID.String()))
	}
	return mapWorkflowTaskWatchers(rows), nil
}

func (s *Store) CreateWorkflowTaskComment(ctx context.Context, item *domain.WorkflowTaskComment, actorID *uuid.UUID) (*domain.WorkflowTaskComment, error) {
	row, err := s.getQueries(ctx).CreateWorkflowTaskComment(ctx, sqlc.CreateWorkflowTaskCommentParams{TenantID: item.TenantID, TaskID: item.TaskID, Visibility: item.Visibility, Body: item.Body, Metadata: jsonBytesFromRaw(item.Metadata), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create workflow task comment", err, tenantIDField(item.TenantID), stringField("workflow_task_id", item.TaskID.String()))
	}
	return mapWorkflowTaskComment(row), nil
}

func (s *Store) ListWorkflowTaskComments(ctx context.Context, tenantID uuid.UUID, taskID uuid.UUID) ([]*domain.WorkflowTaskComment, error) {
	rows, err := s.getQueries(ctx).ListWorkflowTaskComments(ctx, sqlc.ListWorkflowTaskCommentsParams{TenantID: tenantID, TaskID: taskID})
	if err != nil {
		return nil, s.logDBError(ctx, "list workflow task comments", err, tenantIDField(tenantID), stringField("workflow_task_id", taskID.String()))
	}
	return mapWorkflowTaskComments(rows), nil
}

func (s *Store) CreateWorkflowTaskAttachment(ctx context.Context, item *domain.WorkflowTaskAttachment, actorID *uuid.UUID) (*domain.WorkflowTaskAttachment, error) {
	row, err := s.getQueries(ctx).CreateWorkflowTaskAttachment(ctx, sqlc.CreateWorkflowTaskAttachmentParams{TenantID: item.TenantID, TaskID: item.TaskID, CommentID: uuidFromPtr(item.CommentID), FileName: item.FileName, ContentType: item.ContentType, StoragePath: item.StoragePath, ChecksumSha256: textFromPtr(item.ChecksumSHA256), SizeBytes: item.SizeBytes, Visibility: item.Visibility, Metadata: jsonBytesFromRaw(item.Metadata), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create workflow task attachment", err, tenantIDField(item.TenantID), stringField("workflow_task_id", item.TaskID.String()))
	}
	return mapWorkflowTaskAttachment(row), nil
}

func (s *Store) ListWorkflowTaskAttachments(ctx context.Context, tenantID uuid.UUID, taskID uuid.UUID) ([]*domain.WorkflowTaskAttachment, error) {
	rows, err := s.getQueries(ctx).ListWorkflowTaskAttachments(ctx, sqlc.ListWorkflowTaskAttachmentsParams{TenantID: tenantID, TaskID: taskID})
	if err != nil {
		return nil, s.logDBError(ctx, "list workflow task attachments", err, tenantIDField(tenantID), stringField("workflow_task_id", taskID.String()))
	}
	return mapWorkflowTaskAttachments(rows), nil
}

func (s *Store) CreateWorkflowTaskEvent(ctx context.Context, item *domain.WorkflowTaskEvent, actorID *uuid.UUID) (*domain.WorkflowTaskEvent, error) {
	row, err := s.getQueries(ctx).CreateWorkflowTaskEvent(ctx, sqlc.CreateWorkflowTaskEventParams{TenantID: item.TenantID, TaskID: item.TaskID, Action: item.Action, FromStatus: textFromPtr(item.FromStatus), ToStatus: textFromPtr(item.ToStatus), ActorUserID: uuidFromPtr(item.ActorUserID), Remarks: textFromPtr(item.Remarks), Metadata: jsonBytesFromRaw(item.Metadata), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create workflow task event", err, tenantIDField(item.TenantID), stringField("workflow_task_id", item.TaskID.String()), stringField("action", item.Action))
	}
	return mapWorkflowTaskEvent(row), nil
}

func (s *Store) ListWorkflowTaskEvents(ctx context.Context, tenantID uuid.UUID, taskID uuid.UUID) ([]*domain.WorkflowTaskEvent, error) {
	rows, err := s.getQueries(ctx).ListWorkflowTaskEvents(ctx, sqlc.ListWorkflowTaskEventsParams{TenantID: tenantID, TaskID: taskID})
	if err != nil {
		return nil, s.logDBError(ctx, "list workflow task events", err, tenantIDField(tenantID), stringField("workflow_task_id", taskID.String()))
	}
	return mapWorkflowTaskEvents(rows), nil
}

func (s *Store) GetWorkflowTaskSummary(ctx context.Context, tenantID uuid.UUID) ([]*domain.WorkflowTaskSummaryRow, error) {
	rows, err := s.getQueries(ctx).GetWorkflowTaskSummary(ctx, tenantID)
	if err != nil {
		return nil, s.logDBError(ctx, "get workflow task summary", err, tenantIDField(tenantID))
	}
	return mapWorkflowTaskSummary(rows), nil
}
