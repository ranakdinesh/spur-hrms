package postgres

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func mapWorkflowDefinition(row sqlc.HrmsWorkflowDefinition) *domain.WorkflowDefinition {
	return &domain.WorkflowDefinition{ID: row.ID, TenantID: row.TenantID, WorkflowKey: row.WorkflowKey, Name: row.Name, ModuleKey: row.ModuleKey, Description: ptrFromText(row.Description), Status: row.Status, VisibilityScope: row.VisibilityScope, Metadata: jsonRaw(row.Metadata, `{}`), Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy)}
}

func mapWorkflowDefinitions(rows []sqlc.HrmsWorkflowDefinition) []*domain.WorkflowDefinition {
	items := make([]*domain.WorkflowDefinition, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapWorkflowDefinition(row))
	}
	return items
}

func mapWorkflowDefinitionStep(row sqlc.HrmsWorkflowDefinitionStep) *domain.WorkflowDefinitionStep {
	return &domain.WorkflowDefinitionStep{ID: row.ID, TenantID: row.TenantID, WorkflowDefinitionID: row.WorkflowDefinitionID, StepOrder: row.StepOrder, StepKey: row.StepKey, Name: row.Name, StepType: row.StepType, AssignmentType: row.AssignmentType, AssignmentValue: ptrFromText(row.AssignmentValue), Required: row.Required, DueOffsetHours: row.DueOffsetHours, AllowedActions: jsonRaw(row.AllowedActions, `[]`), Metadata: jsonRaw(row.Metadata, `{}`), Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy)}
}

func mapWorkflowDefinitionSteps(rows []sqlc.HrmsWorkflowDefinitionStep) []*domain.WorkflowDefinitionStep {
	items := make([]*domain.WorkflowDefinitionStep, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapWorkflowDefinitionStep(row))
	}
	return items
}

func mapOperationTemplate(row sqlc.HrmsOperationTemplate) *domain.OperationTemplate {
	return &domain.OperationTemplate{ID: row.ID, TenantID: row.TenantID, TemplateKey: row.TemplateKey, Name: row.Name, Category: row.Category, SourceModule: row.SourceModule, SourceType: row.SourceType, WorkflowDefinitionID: ptrFromUUID(row.WorkflowDefinitionID), DefaultPriority: row.DefaultPriority, DefaultSeverity: row.DefaultSeverity, AllowedActions: jsonRaw(row.AllowedActions, `[]`), LaunchSchema: jsonRaw(row.LaunchSchema, `{}`), IsActive: row.IsActive, Metadata: jsonRaw(row.Metadata, `{}`), Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy)}
}

func mapOperationTemplates(rows []sqlc.HrmsOperationTemplate) []*domain.OperationTemplate {
	items := make([]*domain.OperationTemplate, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapOperationTemplate(row))
	}
	return items
}

func mapWorkflowTask(row sqlc.HrmsWorkflowTask) *domain.WorkflowTask {
	return workflowTaskFromParts(row.ID, row.TenantID, row.TaskNumber, row.TemplateID, row.WorkflowDefinitionID, row.WorkflowStepID, row.ParentTaskID, row.SourceModule, row.SourceType, row.SourceID, row.SourceRecordLabel, row.Title, row.Description, row.RequesterUserID, row.AssigneeUserID, row.AssigneeRole, row.AssigneeTeam, row.DelegatedFromUserID, row.Status, row.Priority, row.Severity, row.VisibilityScope, row.DueAt, row.CompletedAt, row.CompletedBy, row.ActionSchema, row.Metadata, row.Inactive, row.CreatedAt, row.CreatedBy, row.UpdatedAt, row.UpdatedBy, 0, 0, false)
}

func mapWorkflowTaskListRow(row sqlc.ListWorkflowTasksRow) *domain.WorkflowTask {
	return workflowTaskFromParts(row.ID, row.TenantID, row.TaskNumber, row.TemplateID, row.WorkflowDefinitionID, row.WorkflowStepID, row.ParentTaskID, row.SourceModule, row.SourceType, row.SourceID, row.SourceRecordLabel, row.Title, row.Description, row.RequesterUserID, row.AssigneeUserID, row.AssigneeRole, row.AssigneeTeam, row.DelegatedFromUserID, row.Status, row.Priority, row.Severity, row.VisibilityScope, row.DueAt, row.CompletedAt, row.CompletedBy, row.ActionSchema, row.Metadata, row.Inactive, row.CreatedAt, row.CreatedBy, row.UpdatedAt, row.UpdatedBy, row.CommentCount, row.AttachmentCount, row.WatchedByViewer)
}

func workflowTaskFromParts(id uuid.UUID, tenantID uuid.UUID, taskNumber string, templateID pgtype.UUID, definitionID pgtype.UUID, stepID pgtype.UUID, parentID pgtype.UUID, sourceModule string, sourceType string, sourceID pgtype.UUID, sourceLabel pgtype.Text, title string, description pgtype.Text, requesterID pgtype.UUID, assigneeID pgtype.UUID, assigneeRole pgtype.Text, assigneeTeam pgtype.Text, delegatedFrom pgtype.UUID, status string, priority int32, severity string, visibility string, dueAt pgtype.Timestamptz, completedAt pgtype.Timestamptz, completedBy pgtype.UUID, actionSchema []byte, metadata []byte, inactive bool, createdAt pgtype.Timestamptz, createdBy pgtype.UUID, updatedAt pgtype.Timestamptz, updatedBy pgtype.UUID, commentCount int32, attachmentCount int32, watched bool) *domain.WorkflowTask {
	return &domain.WorkflowTask{ID: id, TenantID: tenantID, TaskNumber: taskNumber, TemplateID: ptrFromUUID(templateID), WorkflowDefinitionID: ptrFromUUID(definitionID), WorkflowStepID: ptrFromUUID(stepID), ParentTaskID: ptrFromUUID(parentID), SourceModule: sourceModule, SourceType: sourceType, SourceID: ptrFromUUID(sourceID), SourceRecordLabel: ptrFromText(sourceLabel), Title: title, Description: ptrFromText(description), RequesterUserID: ptrFromUUID(requesterID), AssigneeUserID: ptrFromUUID(assigneeID), AssigneeRole: ptrFromText(assigneeRole), AssigneeTeam: ptrFromText(assigneeTeam), DelegatedFromUserID: ptrFromUUID(delegatedFrom), Status: status, Priority: priority, Severity: severity, VisibilityScope: visibility, DueAt: ptrFromTimestamptz(dueAt), CompletedAt: ptrFromTimestamptz(completedAt), CompletedBy: ptrFromUUID(completedBy), ActionSchema: jsonRaw(actionSchema, `[]`), Metadata: jsonRaw(metadata, `{}`), Inactive: inactive, CreatedAt: timeFromTimestamptz(createdAt), CreatedBy: ptrFromUUID(createdBy), UpdatedAt: timeFromTimestamptz(updatedAt), UpdatedBy: ptrFromUUID(updatedBy), CommentCount: commentCount, AttachmentCount: attachmentCount, WatchedByViewer: watched}
}

func mapWorkflowTasks(rows []sqlc.ListWorkflowTasksRow) []*domain.WorkflowTask {
	items := make([]*domain.WorkflowTask, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapWorkflowTaskListRow(row))
	}
	return items
}

func mapWorkflowTaskWatcher(row sqlc.HrmsWorkflowTaskWatcher) *domain.WorkflowTaskWatcher {
	return &domain.WorkflowTaskWatcher{ID: row.ID, TenantID: row.TenantID, TaskID: row.TaskID, WatcherUserID: row.WatcherUserID, WatchReason: ptrFromText(row.WatchReason), Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy)}
}

func mapWorkflowTaskWatchers(rows []sqlc.HrmsWorkflowTaskWatcher) []*domain.WorkflowTaskWatcher {
	items := make([]*domain.WorkflowTaskWatcher, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapWorkflowTaskWatcher(row))
	}
	return items
}

func mapWorkflowTaskComment(row sqlc.HrmsWorkflowTaskComment) *domain.WorkflowTaskComment {
	return &domain.WorkflowTaskComment{ID: row.ID, TenantID: row.TenantID, TaskID: row.TaskID, Visibility: row.Visibility, Body: row.Body, Metadata: jsonRaw(row.Metadata, `{}`), Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy)}
}

func mapWorkflowTaskComments(rows []sqlc.HrmsWorkflowTaskComment) []*domain.WorkflowTaskComment {
	items := make([]*domain.WorkflowTaskComment, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapWorkflowTaskComment(row))
	}
	return items
}

func mapWorkflowTaskAttachment(row sqlc.HrmsWorkflowTaskAttachment) *domain.WorkflowTaskAttachment {
	return &domain.WorkflowTaskAttachment{ID: row.ID, TenantID: row.TenantID, TaskID: row.TaskID, CommentID: ptrFromUUID(row.CommentID), FileName: row.FileName, ContentType: row.ContentType, StoragePath: row.StoragePath, ChecksumSHA256: ptrFromText(row.ChecksumSha256), SizeBytes: row.SizeBytes, Visibility: row.Visibility, Metadata: jsonRaw(row.Metadata, `{}`), Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy)}
}

func mapWorkflowTaskAttachments(rows []sqlc.HrmsWorkflowTaskAttachment) []*domain.WorkflowTaskAttachment {
	items := make([]*domain.WorkflowTaskAttachment, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapWorkflowTaskAttachment(row))
	}
	return items
}

func mapWorkflowTaskEvent(row sqlc.HrmsWorkflowTaskEvent) *domain.WorkflowTaskEvent {
	return &domain.WorkflowTaskEvent{ID: row.ID, TenantID: row.TenantID, TaskID: row.TaskID, Action: row.Action, FromStatus: ptrFromText(row.FromStatus), ToStatus: ptrFromText(row.ToStatus), ActorUserID: ptrFromUUID(row.ActorUserID), Remarks: ptrFromText(row.Remarks), Metadata: jsonRaw(row.Metadata, `{}`), Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy)}
}

func mapWorkflowTaskEvents(rows []sqlc.HrmsWorkflowTaskEvent) []*domain.WorkflowTaskEvent {
	items := make([]*domain.WorkflowTaskEvent, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapWorkflowTaskEvent(row))
	}
	return items
}

func mapWorkflowTaskSummary(rows []sqlc.GetWorkflowTaskSummaryRow) []*domain.WorkflowTaskSummaryRow {
	items := make([]*domain.WorkflowTaskSummaryRow, 0, len(rows))
	for _, row := range rows {
		items = append(items, &domain.WorkflowTaskSummaryRow{Metric: row.Metric, MetricCount: row.MetricCount})
	}
	return items
}

func jsonRaw(value []byte, fallback string) json.RawMessage {
	if len(value) == 0 || !json.Valid(value) {
		return json.RawMessage(fallback)
	}
	return json.RawMessage(value)
}
