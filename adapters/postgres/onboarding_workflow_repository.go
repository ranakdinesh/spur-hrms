package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func (s *Store) CreateOnboardingWorkflow(ctx context.Context, item *domain.OnboardingWorkflow, actorID *uuid.UUID) (*domain.OnboardingWorkflow, error) {
	q := s.getQueries(ctx)
	if item.IsDefault {
		if err := q.ClearDefaultOnboardingWorkflows(ctx, sqlc.ClearDefaultOnboardingWorkflowsParams{TenantID: item.TenantID, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
			return nil, s.logDBError(ctx, "clear default onboarding workflows", err, tenantIDField(item.TenantID))
		}
	}
	row, err := q.CreateOnboardingWorkflow(ctx, sqlc.CreateOnboardingWorkflowParams{TenantID: item.TenantID, Name: item.Name, Description: textFromPtr(item.Description), IsDefault: item.IsDefault, IsActive: item.IsActive, CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create onboarding workflow", fmt.Errorf("hrms: create onboarding workflow: %w", err), tenantIDField(item.TenantID))
	}
	return mapOnboardingWorkflow(row), nil
}

func (s *Store) ListOnboardingWorkflows(ctx context.Context, tenantID uuid.UUID) ([]*domain.OnboardingWorkflow, error) {
	rows, err := s.getQueries(ctx).ListOnboardingWorkflows(ctx, tenantID)
	if err != nil {
		return nil, s.logDBError(ctx, "list onboarding workflows", err, tenantIDField(tenantID))
	}
	return mapOnboardingWorkflows(rows), nil
}

func (s *Store) GetOnboardingWorkflow(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.OnboardingWorkflow, error) {
	row, err := s.getQueries(ctx).GetOnboardingWorkflow(ctx, sqlc.GetOnboardingWorkflowParams{TenantID: tenantID, ID: id})
	if err != nil {
		return nil, s.logDBError(ctx, "get onboarding workflow", err, tenantIDField(tenantID), stringField("onboarding_workflow_id", id.String()))
	}
	return mapOnboardingWorkflow(row), nil
}

func (s *Store) UpdateOnboardingWorkflow(ctx context.Context, item *domain.OnboardingWorkflow, actorID *uuid.UUID) (*domain.OnboardingWorkflow, error) {
	q := s.getQueries(ctx)
	if item.IsDefault {
		if err := q.ClearDefaultOnboardingWorkflows(ctx, sqlc.ClearDefaultOnboardingWorkflowsParams{TenantID: item.TenantID, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
			return nil, s.logDBError(ctx, "clear default onboarding workflows", err, tenantIDField(item.TenantID))
		}
	}
	row, err := q.UpdateOnboardingWorkflow(ctx, sqlc.UpdateOnboardingWorkflowParams{TenantID: item.TenantID, ID: item.ID, Name: item.Name, Description: textFromPtr(item.Description), IsDefault: item.IsDefault, IsActive: item.IsActive, UpdatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "update onboarding workflow", err, tenantIDField(item.TenantID), stringField("onboarding_workflow_id", item.ID.String()))
	}
	return mapOnboardingWorkflow(row), nil
}

func (s *Store) DeleteOnboardingWorkflow(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteOnboardingWorkflow(ctx, sqlc.SoftDeleteOnboardingWorkflowParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete onboarding workflow", err, tenantIDField(tenantID), stringField("onboarding_workflow_id", id.String()))
	}
	return nil
}

func (s *Store) CreateOnboardingTask(ctx context.Context, item *domain.OnboardingTask, actorID *uuid.UUID) (*domain.OnboardingTask, error) {
	row, err := s.getQueries(ctx).CreateOnboardingTask(ctx, sqlc.CreateOnboardingTaskParams{TenantID: item.TenantID, WorkflowID: item.WorkflowID, Title: item.Title, Description: textFromPtr(item.Description), DueDays: item.DueDays, IsRequired: item.IsRequired, SortOrder: item.SortOrder, CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create onboarding task", err, tenantIDField(item.TenantID), stringField("onboarding_workflow_id", item.WorkflowID.String()))
	}
	return mapOnboardingTask(row), nil
}

func (s *Store) ListOnboardingTasks(ctx context.Context, tenantID uuid.UUID, workflowID uuid.UUID) ([]*domain.OnboardingTask, error) {
	rows, err := s.getQueries(ctx).ListOnboardingTasksByWorkflow(ctx, sqlc.ListOnboardingTasksByWorkflowParams{TenantID: tenantID, WorkflowID: workflowID})
	if err != nil {
		return nil, s.logDBError(ctx, "list onboarding tasks", err, tenantIDField(tenantID), stringField("onboarding_workflow_id", workflowID.String()))
	}
	return mapOnboardingTasks(rows), nil
}

func (s *Store) GetOnboardingTask(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.OnboardingTask, error) {
	row, err := s.getQueries(ctx).GetOnboardingTask(ctx, sqlc.GetOnboardingTaskParams{TenantID: tenantID, ID: id})
	if err != nil {
		return nil, s.logDBError(ctx, "get onboarding task", err, tenantIDField(tenantID), stringField("onboarding_task_id", id.String()))
	}
	return mapOnboardingTask(row), nil
}

func (s *Store) UpdateOnboardingTask(ctx context.Context, item *domain.OnboardingTask, actorID *uuid.UUID) (*domain.OnboardingTask, error) {
	row, err := s.getQueries(ctx).UpdateOnboardingTask(ctx, sqlc.UpdateOnboardingTaskParams{TenantID: item.TenantID, ID: item.ID, Title: item.Title, Description: textFromPtr(item.Description), DueDays: item.DueDays, IsRequired: item.IsRequired, SortOrder: item.SortOrder, UpdatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "update onboarding task", err, tenantIDField(item.TenantID), stringField("onboarding_task_id", item.ID.String()))
	}
	return mapOnboardingTask(row), nil
}

func (s *Store) DeleteOnboardingTask(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteOnboardingTask(ctx, sqlc.SoftDeleteOnboardingTaskParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete onboarding task", err, tenantIDField(tenantID), stringField("onboarding_task_id", id.String()))
	}
	return nil
}

func (s *Store) CreateOnboardingWorkflowAssignment(ctx context.Context, item *domain.OnboardingWorkflowAssignment, actorID *uuid.UUID) (*domain.OnboardingWorkflowAssignment, error) {
	row, err := s.getQueries(ctx).CreateOnboardingWorkflowAssignment(ctx, sqlc.CreateOnboardingWorkflowAssignmentParams{TenantID: item.TenantID, WorkflowID: item.WorkflowID, Name: item.Name, JobPostingID: uuidFromPtr(item.JobPostingID), JobPositionID: uuidFromPtr(item.JobPositionID), DepartmentID: uuidFromPtr(item.DepartmentID), EmploymentTypeID: uuidFromPtr(item.EmploymentTypeID), Priority: item.Priority, CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create onboarding assignment", err, tenantIDField(item.TenantID), stringField("onboarding_workflow_id", item.WorkflowID.String()))
	}
	return mapOnboardingAssignment(row), nil
}

func (s *Store) ListOnboardingWorkflowAssignments(ctx context.Context, tenantID uuid.UUID) ([]*domain.OnboardingWorkflowAssignment, error) {
	rows, err := s.getQueries(ctx).ListOnboardingWorkflowAssignments(ctx, tenantID)
	if err != nil {
		return nil, s.logDBError(ctx, "list onboarding assignments", err, tenantIDField(tenantID))
	}
	return mapOnboardingAssignments(rows), nil
}

func (s *Store) GetOnboardingWorkflowAssignment(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.OnboardingWorkflowAssignment, error) {
	row, err := s.getQueries(ctx).GetOnboardingWorkflowAssignment(ctx, sqlc.GetOnboardingWorkflowAssignmentParams{TenantID: tenantID, ID: id})
	if err != nil {
		return nil, s.logDBError(ctx, "get onboarding assignment", err, tenantIDField(tenantID), stringField("onboarding_assignment_id", id.String()))
	}
	return mapOnboardingAssignment(row), nil
}

func (s *Store) UpdateOnboardingWorkflowAssignment(ctx context.Context, item *domain.OnboardingWorkflowAssignment, actorID *uuid.UUID) (*domain.OnboardingWorkflowAssignment, error) {
	row, err := s.getQueries(ctx).UpdateOnboardingWorkflowAssignment(ctx, sqlc.UpdateOnboardingWorkflowAssignmentParams{TenantID: item.TenantID, ID: item.ID, WorkflowID: item.WorkflowID, Name: item.Name, JobPostingID: uuidFromPtr(item.JobPostingID), JobPositionID: uuidFromPtr(item.JobPositionID), DepartmentID: uuidFromPtr(item.DepartmentID), EmploymentTypeID: uuidFromPtr(item.EmploymentTypeID), Priority: item.Priority, UpdatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "update onboarding assignment", err, tenantIDField(item.TenantID), stringField("onboarding_assignment_id", item.ID.String()))
	}
	return mapOnboardingAssignment(row), nil
}

func (s *Store) DeleteOnboardingWorkflowAssignment(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteOnboardingWorkflowAssignment(ctx, sqlc.SoftDeleteOnboardingWorkflowAssignmentParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete onboarding assignment", err, tenantIDField(tenantID), stringField("onboarding_assignment_id", id.String()))
	}
	return nil
}
