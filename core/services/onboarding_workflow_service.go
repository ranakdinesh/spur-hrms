package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) CreateOnboardingWorkflow(ctx context.Context, cmd ports.OnboardingWorkflowCommand) (*domain.OnboardingWorkflow, error) {
	item, err := domain.NewOnboardingWorkflow(onboardingWorkflowInput(cmd))
	if err != nil {
		s.logError("validate onboarding workflow create", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	result, err := s.onboardingWorkflows.CreateOnboardingWorkflow(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create onboarding workflow", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) ListOnboardingWorkflows(ctx context.Context, tenantID uuid.UUID) ([]*domain.OnboardingWorkflow, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate onboarding workflow list tenant", err)
		return nil, err
	}
	items, err := s.onboardingWorkflows.ListOnboardingWorkflows(ctx, tenantID)
	if err != nil {
		s.logError("list onboarding workflows", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	for _, item := range items {
		tasks, err := s.onboardingWorkflows.ListOnboardingTasks(ctx, tenantID, item.ID)
		if err == nil {
			item.Tasks = tasks
		}
	}
	return items, nil
}

func (s *TenantService) GetOnboardingWorkflow(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.OnboardingWorkflow, error) {
	if tenantID == uuid.Nil || id == uuid.Nil {
		err := domain.ErrInvalidOnboardingWorkflowID
		s.logError("validate onboarding workflow get", err)
		return nil, err
	}
	item, err := s.onboardingWorkflows.GetOnboardingWorkflow(ctx, tenantID, id)
	if err != nil {
		s.logError("get onboarding workflow", err, serviceTenantIDField(tenantID), serviceStringField("onboarding_workflow_id", id.String()))
		return nil, err
	}
	item.Tasks, _ = s.onboardingWorkflows.ListOnboardingTasks(ctx, tenantID, id)
	return item, nil
}

func (s *TenantService) UpdateOnboardingWorkflow(ctx context.Context, cmd ports.OnboardingWorkflowCommand) (*domain.OnboardingWorkflow, error) {
	if _, err := s.GetOnboardingWorkflow(ctx, cmd.TenantID, cmd.ID); err != nil {
		return nil, err
	}
	item, err := domain.NewOnboardingWorkflow(onboardingWorkflowInput(cmd))
	if err != nil {
		s.logError("validate onboarding workflow update", err, serviceTenantIDField(cmd.TenantID), serviceStringField("onboarding_workflow_id", cmd.ID.String()))
		return nil, err
	}
	item.ID = cmd.ID
	result, err := s.onboardingWorkflows.UpdateOnboardingWorkflow(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("update onboarding workflow", err, serviceTenantIDField(cmd.TenantID), serviceStringField("onboarding_workflow_id", cmd.ID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) DeleteOnboardingWorkflow(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if _, err := s.GetOnboardingWorkflow(ctx, tenantID, id); err != nil {
		return err
	}
	if err := s.onboardingWorkflows.DeleteOnboardingWorkflow(ctx, tenantID, id, actorID); err != nil {
		s.logError("delete onboarding workflow", err, serviceTenantIDField(tenantID), serviceStringField("onboarding_workflow_id", id.String()))
		return err
	}
	return nil
}

func (s *TenantService) CreateOnboardingTask(ctx context.Context, cmd ports.OnboardingTaskCommand) (*domain.OnboardingTask, error) {
	if _, err := s.GetOnboardingWorkflow(ctx, cmd.TenantID, cmd.WorkflowID); err != nil {
		return nil, err
	}
	item, err := domain.NewOnboardingTask(onboardingTaskInput(cmd))
	if err != nil {
		s.logError("validate onboarding task create", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	return s.onboardingWorkflows.CreateOnboardingTask(ctx, item, cmd.ActorID)
}

func (s *TenantService) ListOnboardingTasks(ctx context.Context, tenantID uuid.UUID, workflowID uuid.UUID) ([]*domain.OnboardingTask, error) {
	if _, err := s.GetOnboardingWorkflow(ctx, tenantID, workflowID); err != nil {
		return nil, err
	}
	return s.onboardingWorkflows.ListOnboardingTasks(ctx, tenantID, workflowID)
}

func (s *TenantService) GetOnboardingTask(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.OnboardingTask, error) {
	if tenantID == uuid.Nil || id == uuid.Nil {
		return nil, domain.ErrInvalidOnboardingTaskID
	}
	return s.onboardingWorkflows.GetOnboardingTask(ctx, tenantID, id)
}

func (s *TenantService) UpdateOnboardingTask(ctx context.Context, cmd ports.OnboardingTaskCommand) (*domain.OnboardingTask, error) {
	if _, err := s.GetOnboardingTask(ctx, cmd.TenantID, cmd.ID); err != nil {
		return nil, err
	}
	item, err := domain.NewOnboardingTask(onboardingTaskInput(cmd))
	if err != nil {
		s.logError("validate onboarding task update", err, serviceTenantIDField(cmd.TenantID), serviceStringField("onboarding_task_id", cmd.ID.String()))
		return nil, err
	}
	item.ID = cmd.ID
	return s.onboardingWorkflows.UpdateOnboardingTask(ctx, item, cmd.ActorID)
}

func (s *TenantService) DeleteOnboardingTask(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if _, err := s.GetOnboardingTask(ctx, tenantID, id); err != nil {
		return err
	}
	return s.onboardingWorkflows.DeleteOnboardingTask(ctx, tenantID, id, actorID)
}

func (s *TenantService) CreateOnboardingWorkflowAssignment(ctx context.Context, cmd ports.OnboardingAssignmentCommand) (*domain.OnboardingWorkflowAssignment, error) {
	if _, err := s.GetOnboardingWorkflow(ctx, cmd.TenantID, cmd.WorkflowID); err != nil {
		return nil, err
	}
	item, err := domain.NewOnboardingWorkflowAssignment(onboardingAssignmentInput(cmd))
	if err != nil {
		s.logError("validate onboarding assignment create", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	return s.onboardingWorkflows.CreateOnboardingWorkflowAssignment(ctx, item, cmd.ActorID)
}

func (s *TenantService) ListOnboardingWorkflowAssignments(ctx context.Context, tenantID uuid.UUID) ([]*domain.OnboardingWorkflowAssignment, error) {
	if tenantID == uuid.Nil {
		return nil, domain.ErrInvalidTenantID
	}
	return s.onboardingWorkflows.ListOnboardingWorkflowAssignments(ctx, tenantID)
}

func (s *TenantService) UpdateOnboardingWorkflowAssignment(ctx context.Context, cmd ports.OnboardingAssignmentCommand) (*domain.OnboardingWorkflowAssignment, error) {
	if _, err := s.onboardingWorkflows.GetOnboardingWorkflowAssignment(ctx, cmd.TenantID, cmd.ID); err != nil {
		return nil, err
	}
	item, err := domain.NewOnboardingWorkflowAssignment(onboardingAssignmentInput(cmd))
	if err != nil {
		return nil, err
	}
	item.ID = cmd.ID
	return s.onboardingWorkflows.UpdateOnboardingWorkflowAssignment(ctx, item, cmd.ActorID)
}

func (s *TenantService) DeleteOnboardingWorkflowAssignment(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if _, err := s.onboardingWorkflows.GetOnboardingWorkflowAssignment(ctx, tenantID, id); err != nil {
		return err
	}
	return s.onboardingWorkflows.DeleteOnboardingWorkflowAssignment(ctx, tenantID, id, actorID)
}

func onboardingWorkflowInput(cmd ports.OnboardingWorkflowCommand) domain.OnboardingWorkflowInput {
	return domain.OnboardingWorkflowInput{TenantID: cmd.TenantID, Name: cmd.Name, Description: cmd.Description, IsDefault: cmd.IsDefault, IsActive: cmd.IsActive}
}

func onboardingTaskInput(cmd ports.OnboardingTaskCommand) domain.OnboardingTaskInput {
	return domain.OnboardingTaskInput{TenantID: cmd.TenantID, WorkflowID: cmd.WorkflowID, Title: cmd.Title, Description: cmd.Description, DueDays: cmd.DueDays, IsRequired: cmd.IsRequired, SortOrder: cmd.SortOrder}
}

func onboardingAssignmentInput(cmd ports.OnboardingAssignmentCommand) domain.OnboardingAssignmentInput {
	return domain.OnboardingAssignmentInput{TenantID: cmd.TenantID, WorkflowID: cmd.WorkflowID, Name: cmd.Name, JobPostingID: cmd.JobPostingID, JobPositionID: cmd.JobPositionID, DepartmentID: cmd.DepartmentID, EmploymentTypeID: cmd.EmploymentTypeID, Priority: cmd.Priority}
}
