package postgres

import (
	"encoding/json"

	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func mapCandidateOnboardingRow(row sqlc.GetCandidateOnboardingRow) *domain.CandidateOnboarding {
	workflowName := row.WorkflowName
	return &domain.CandidateOnboarding{
		ID: row.ID, TenantID: row.TenantID, CandidateID: row.CandidateID, CandidateFirstname: ptrFromText(row.CandidateFirstname), CandidateLastname: ptrFromText(row.CandidateLastname), CandidateEmail: ptrFromText(row.CandidateEmail),
		WorkflowID: row.WorkflowID, WorkflowName: &workflowName, OnboardingStatus: row.OnboardingStatus, ProgressPercentage: row.ProgressPercentage,
		TotalTasks: row.TotalTasks, CompletedTasks: row.CompletedTasks, RequiredTasks: row.RequiredTasks, CompletedRequiredTasks: row.CompletedRequiredTasks, OverdueTasks: row.OverdueTasks,
		StartedAt: ptrFromTimestamptz(row.StartedAt), CompletedAt: ptrFromTimestamptz(row.CompletedAt), Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy),
	}
}

func mapCandidateOnboardingByCandidateRow(row sqlc.GetCandidateOnboardingByCandidateRow) *domain.CandidateOnboarding {
	workflowName := row.WorkflowName
	return &domain.CandidateOnboarding{
		ID: row.ID, TenantID: row.TenantID, CandidateID: row.CandidateID, CandidateFirstname: ptrFromText(row.CandidateFirstname), CandidateLastname: ptrFromText(row.CandidateLastname), CandidateEmail: ptrFromText(row.CandidateEmail),
		WorkflowID: row.WorkflowID, WorkflowName: &workflowName, OnboardingStatus: row.OnboardingStatus, ProgressPercentage: row.ProgressPercentage,
		TotalTasks: row.TotalTasks, CompletedTasks: row.CompletedTasks, RequiredTasks: row.RequiredTasks, CompletedRequiredTasks: row.CompletedRequiredTasks, OverdueTasks: row.OverdueTasks,
		StartedAt: ptrFromTimestamptz(row.StartedAt), CompletedAt: ptrFromTimestamptz(row.CompletedAt), Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy),
	}
}

func mapCandidateOnboardingListRow(row sqlc.ListCandidateOnboardingsRow) *domain.CandidateOnboarding {
	workflowName := row.WorkflowName
	return &domain.CandidateOnboarding{
		ID: row.ID, TenantID: row.TenantID, CandidateID: row.CandidateID, CandidateFirstname: ptrFromText(row.CandidateFirstname), CandidateLastname: ptrFromText(row.CandidateLastname), CandidateEmail: ptrFromText(row.CandidateEmail),
		WorkflowID: row.WorkflowID, WorkflowName: &workflowName, OnboardingStatus: row.OnboardingStatus, ProgressPercentage: row.ProgressPercentage,
		TotalTasks: row.TotalTasks, CompletedTasks: row.CompletedTasks, RequiredTasks: row.RequiredTasks, CompletedRequiredTasks: row.CompletedRequiredTasks, OverdueTasks: row.OverdueTasks,
		StartedAt: ptrFromTimestamptz(row.StartedAt), CompletedAt: ptrFromTimestamptz(row.CompletedAt), Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy),
	}
}

func mapCandidateOnboardingBase(row sqlc.HrmsCandidateOnboarding) *domain.CandidateOnboarding {
	return &domain.CandidateOnboarding{ID: row.ID, TenantID: row.TenantID, CandidateID: row.CandidateID, WorkflowID: row.WorkflowID, OnboardingStatus: row.OnboardingStatus, ProgressPercentage: row.ProgressPercentage, StartedAt: ptrFromTimestamptz(row.StartedAt), CompletedAt: ptrFromTimestamptz(row.CompletedAt), Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy)}
}

func mapCandidateOnboardingTasks(rows []sqlc.ListCandidateOnboardingTasksRow) []*domain.CandidateOnboardingTask {
	items := make([]*domain.CandidateOnboardingTask, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapCandidateOnboardingTaskRow(row))
	}
	return items
}

func mapCandidateOnboardingTaskRow(row sqlc.ListCandidateOnboardingTasksRow) *domain.CandidateOnboardingTask {
	taskTitle := row.TaskTitle
	return &domain.CandidateOnboardingTask{
		ID: row.ID, TenantID: row.TenantID, CandidateOnboardingID: row.CandidateOnboardingID, OnboardingTaskID: row.OnboardingTaskID,
		TaskTitle: &taskTitle, TaskDescription: ptrFromText(row.TaskDescription), TaskDueDays: row.TaskDueDays, TaskIsRequired: row.TaskIsRequired, TaskSortOrder: row.TaskSortOrder,
		Status: row.Status, DueAt: ptrFromTimestamptz(row.DueAt), StartedAt: ptrFromTimestamptz(row.StartedAt), CompletedAt: ptrFromTimestamptz(row.CompletedAt), CompletedBy: ptrFromUUID(row.CompletedBy), ReviewedBy: ptrFromUUID(row.ReviewedBy), Remarks: ptrFromText(row.Remarks), IsOverdue: row.IsOverdue,
		Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy),
	}
}

func mapCandidateOnboardingTaskGetRow(row sqlc.GetCandidateOnboardingTaskRow) *domain.CandidateOnboardingTask {
	taskTitle := row.TaskTitle
	return &domain.CandidateOnboardingTask{
		ID: row.ID, TenantID: row.TenantID, CandidateOnboardingID: row.CandidateOnboardingID, OnboardingTaskID: row.OnboardingTaskID,
		TaskTitle: &taskTitle, TaskDescription: ptrFromText(row.TaskDescription), TaskDueDays: row.TaskDueDays, TaskIsRequired: row.TaskIsRequired, TaskSortOrder: row.TaskSortOrder,
		Status: row.Status, DueAt: ptrFromTimestamptz(row.DueAt), StartedAt: ptrFromTimestamptz(row.StartedAt), CompletedAt: ptrFromTimestamptz(row.CompletedAt), CompletedBy: ptrFromUUID(row.CompletedBy), ReviewedBy: ptrFromUUID(row.ReviewedBy), Remarks: ptrFromText(row.Remarks), IsOverdue: row.IsOverdue,
		Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy),
	}
}

func mapCandidateOnboardingTaskBase(row sqlc.HrmsCandidateOnboardingTask) *domain.CandidateOnboardingTask {
	return &domain.CandidateOnboardingTask{ID: row.ID, TenantID: row.TenantID, CandidateOnboardingID: row.CandidateOnboardingID, OnboardingTaskID: row.OnboardingTaskID, Status: row.Status, DueAt: ptrFromTimestamptz(row.DueAt), StartedAt: ptrFromTimestamptz(row.StartedAt), CompletedAt: ptrFromTimestamptz(row.CompletedAt), CompletedBy: ptrFromUUID(row.CompletedBy), ReviewedBy: ptrFromUUID(row.ReviewedBy), Remarks: ptrFromText(row.Remarks), Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy)}
}

func mapCandidateOnboardingEvent(row sqlc.HrmsCandidateOnboardingEvent) *domain.CandidateOnboardingEvent {
	metadata := json.RawMessage(row.Metadata)
	if len(metadata) == 0 {
		metadata = json.RawMessage(`{}`)
	}
	return &domain.CandidateOnboardingEvent{ID: row.ID, TenantID: row.TenantID, CandidateOnboardingID: row.CandidateOnboardingID, CandidateOnboardingTaskID: ptrFromUUID(row.CandidateOnboardingTaskID), Action: row.Action, FromStatus: ptrFromText(row.FromStatus), ToStatus: ptrFromText(row.ToStatus), Remarks: ptrFromText(row.Remarks), Metadata: metadata, Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy)}
}

func mapCandidateOnboardingEvents(rows []sqlc.HrmsCandidateOnboardingEvent) []*domain.CandidateOnboardingEvent {
	items := make([]*domain.CandidateOnboardingEvent, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapCandidateOnboardingEvent(row))
	}
	return items
}
