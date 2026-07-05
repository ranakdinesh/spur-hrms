package postgres

import (
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func mapOnboardingWorkflow(row sqlc.HrmsOnboardingWorkflow) *domain.OnboardingWorkflow {
	return &domain.OnboardingWorkflow{ID: row.ID, TenantID: row.TenantID, Name: row.Name, Description: ptrFromText(row.Description), IsDefault: row.IsDefault, IsActive: row.IsActive, Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy)}
}

func mapOnboardingWorkflows(rows []sqlc.HrmsOnboardingWorkflow) []*domain.OnboardingWorkflow {
	items := make([]*domain.OnboardingWorkflow, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapOnboardingWorkflow(row))
	}
	return items
}

func mapOnboardingTask(row sqlc.HrmsOnboardingTask) *domain.OnboardingTask {
	return &domain.OnboardingTask{ID: row.ID, TenantID: row.TenantID, WorkflowID: row.WorkflowID, Title: row.Title, Description: ptrFromText(row.Description), DueDays: row.DueDays, IsRequired: row.IsRequired, SortOrder: row.SortOrder, Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy)}
}

func mapOnboardingTasks(rows []sqlc.HrmsOnboardingTask) []*domain.OnboardingTask {
	items := make([]*domain.OnboardingTask, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapOnboardingTask(row))
	}
	return items
}

func mapOnboardingAssignment(row sqlc.HrmsOnboardingWorkflowAssignment) *domain.OnboardingWorkflowAssignment {
	return &domain.OnboardingWorkflowAssignment{ID: row.ID, TenantID: row.TenantID, WorkflowID: row.WorkflowID, Name: row.Name, JobPostingID: ptrFromUUID(row.JobPostingID), JobPositionID: ptrFromUUID(row.JobPositionID), DepartmentID: ptrFromUUID(row.DepartmentID), EmploymentTypeID: ptrFromUUID(row.EmploymentTypeID), Priority: row.Priority, Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy)}
}

func mapOnboardingAssignmentListRow(row sqlc.ListOnboardingWorkflowAssignmentsRow) *domain.OnboardingWorkflowAssignment {
	workflowName := row.WorkflowName
	return &domain.OnboardingWorkflowAssignment{ID: row.ID, TenantID: row.TenantID, WorkflowID: row.WorkflowID, WorkflowName: &workflowName, Name: row.Name, JobPostingID: ptrFromUUID(row.JobPostingID), JobPostingTitle: ptrFromText(row.JobPostingTitle), JobPositionID: ptrFromUUID(row.JobPositionID), JobPositionTitle: ptrFromText(row.JobPositionTitle), DepartmentID: ptrFromUUID(row.DepartmentID), DepartmentName: ptrFromText(row.DepartmentName), EmploymentTypeID: ptrFromUUID(row.EmploymentTypeID), EmploymentTypeName: ptrFromText(row.EmploymentTypeName), Priority: row.Priority, Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy)}
}

func mapOnboardingAssignments(rows []sqlc.ListOnboardingWorkflowAssignmentsRow) []*domain.OnboardingWorkflowAssignment {
	items := make([]*domain.OnboardingWorkflowAssignment, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapOnboardingAssignmentListRow(row))
	}
	return items
}
