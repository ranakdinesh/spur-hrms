package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) StartCandidateOnboarding(ctx context.Context, cmd ports.StartCandidateOnboardingCommand) (*domain.CandidateOnboarding, error) {
	if cmd.TenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate start candidate onboarding tenant", err)
		return nil, err
	}
	if cmd.CandidateID == uuid.Nil {
		err := domain.ErrInvalidCandidateID
		s.logError("validate start candidate onboarding candidate", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	if existing, err := s.candidateOnboardings.GetCandidateOnboardingByCandidate(ctx, cmd.TenantID, cmd.CandidateID); err == nil {
		return s.enrichCandidateOnboarding(ctx, existing)
	} else if !errors.Is(err, domain.ErrCandidateOnboardingNotFound) {
		s.logError("lookup existing candidate onboarding", err, serviceTenantIDField(cmd.TenantID), serviceStringField("candidate_id", cmd.CandidateID.String()))
		return nil, err
	}

	workflow, err := s.resolveCandidateOnboardingWorkflow(ctx, cmd)
	if err != nil {
		s.logError("resolve candidate onboarding workflow", err, serviceTenantIDField(cmd.TenantID), serviceStringField("candidate_id", cmd.CandidateID.String()))
		return nil, err
	}

	item, err := domain.NewCandidateOnboarding(domain.CandidateOnboardingInput{TenantID: cmd.TenantID, CandidateID: cmd.CandidateID, WorkflowID: workflow.ID, Status: domain.OnboardStatusInProgress})
	if err != nil {
		s.logError("validate candidate onboarding create", err, serviceTenantIDField(cmd.TenantID), serviceStringField("candidate_id", cmd.CandidateID.String()))
		return nil, err
	}
	created, err := s.candidateOnboardings.CreateCandidateOnboarding(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create candidate onboarding", err, serviceTenantIDField(cmd.TenantID), serviceStringField("candidate_id", cmd.CandidateID.String()), serviceStringField("workflow_id", workflow.ID.String()))
		return nil, err
	}
	if _, err := s.candidateOnboardings.CreateCandidateOnboardingTasksFromWorkflow(ctx, cmd.TenantID, created.ID, workflow.ID, cmd.ActorID); err != nil {
		s.logError("create candidate onboarding task snapshots", err, serviceTenantIDField(cmd.TenantID), serviceStringField("candidate_onboarding_id", created.ID.String()))
		return nil, err
	}
	if _, err := s.candidateOnboardings.RecalculateCandidateOnboardingProgress(ctx, cmd.TenantID, created.ID, cmd.ActorID); err != nil {
		s.logError("recalculate started candidate onboarding", err, serviceTenantIDField(cmd.TenantID), serviceStringField("candidate_onboarding_id", created.ID.String()))
		return nil, err
	}
	if _, err := s.createCandidateOnboardingEvent(ctx, ports.CandidateOnboardingEventCommand{TenantID: cmd.TenantID, CandidateOnboardingID: created.ID, Action: "started", ToStatus: &created.OnboardingStatus, ActorID: cmd.ActorID}); err != nil {
		s.logError("record candidate onboarding start event", err, serviceTenantIDField(cmd.TenantID), serviceStringField("candidate_onboarding_id", created.ID.String()))
	}
	return s.GetCandidateOnboarding(ctx, cmd.TenantID, created.ID)
}

func (s *TenantService) ListCandidateOnboardings(ctx context.Context, filter domain.CandidateOnboardingFilter) (*domain.CandidateOnboardingPage, error) {
	if filter.TenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate list candidate onboardings tenant", err)
		return nil, err
	}
	page, err := s.candidateOnboardings.ListCandidateOnboardings(ctx, filter)
	if err != nil {
		s.logError("list candidate onboardings", err, serviceTenantIDField(filter.TenantID))
		return nil, err
	}
	return page, nil
}

func (s *TenantService) GetCandidateOnboarding(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.CandidateOnboarding, error) {
	if tenantID == uuid.Nil || id == uuid.Nil {
		err := domain.ErrInvalidCandidateOnboardingID
		s.logError("validate get candidate onboarding", err)
		return nil, err
	}
	item, err := s.candidateOnboardings.GetCandidateOnboarding(ctx, tenantID, id)
	if err != nil {
		s.logError("get candidate onboarding", err, serviceTenantIDField(tenantID), serviceStringField("candidate_onboarding_id", id.String()))
		return nil, err
	}
	return s.enrichCandidateOnboarding(ctx, item)
}

func (s *TenantService) GetCandidateOnboardingByCandidate(ctx context.Context, tenantID uuid.UUID, candidateID uuid.UUID) (*domain.CandidateOnboarding, error) {
	if tenantID == uuid.Nil || candidateID == uuid.Nil {
		err := domain.ErrInvalidCandidateID
		s.logError("validate get candidate onboarding by candidate", err)
		return nil, err
	}
	item, err := s.candidateOnboardings.GetCandidateOnboardingByCandidate(ctx, tenantID, candidateID)
	if err != nil {
		s.logError("get candidate onboarding by candidate", err, serviceTenantIDField(tenantID), serviceStringField("candidate_id", candidateID.String()))
		return nil, err
	}
	return s.enrichCandidateOnboarding(ctx, item)
}

func (s *TenantService) UpdateCandidateOnboardingTaskStatus(ctx context.Context, cmd ports.CandidateOnboardingTaskStatusCommand) (*domain.CandidateOnboarding, error) {
	if cmd.TenantID == uuid.Nil || cmd.TaskID == uuid.Nil {
		err := domain.ErrInvalidCandidateOnboardingTaskID
		s.logError("validate candidate onboarding task status", err)
		return nil, err
	}
	status, err := domain.ValidateCandidateOnboardingTaskStatus(cmd.Status)
	if err != nil {
		s.logError("validate candidate onboarding task status value", err, serviceTenantIDField(cmd.TenantID), serviceStringField("candidate_onboarding_task_id", cmd.TaskID.String()))
		return nil, err
	}
	before, err := s.candidateOnboardings.GetCandidateOnboardingTask(ctx, cmd.TenantID, cmd.TaskID)
	if err != nil {
		s.logError("get candidate onboarding task before status update", err, serviceTenantIDField(cmd.TenantID), serviceStringField("candidate_onboarding_task_id", cmd.TaskID.String()))
		return nil, err
	}
	updatedTask, err := s.candidateOnboardings.UpdateCandidateOnboardingTaskStatus(ctx, cmd.TenantID, cmd.TaskID, status, cmd.Remarks, cmd.ActorID)
	if err != nil {
		s.logError("update candidate onboarding task status", err, serviceTenantIDField(cmd.TenantID), serviceStringField("candidate_onboarding_task_id", cmd.TaskID.String()))
		return nil, err
	}
	recalculated, err := s.candidateOnboardings.RecalculateCandidateOnboardingProgress(ctx, cmd.TenantID, updatedTask.CandidateOnboardingID, cmd.ActorID)
	if err != nil {
		s.logError("recalculate candidate onboarding after task status", err, serviceTenantIDField(cmd.TenantID), serviceStringField("candidate_onboarding_id", updatedTask.CandidateOnboardingID.String()))
		return nil, err
	}
	if _, err := s.createCandidateOnboardingEvent(ctx, ports.CandidateOnboardingEventCommand{TenantID: cmd.TenantID, CandidateOnboardingID: updatedTask.CandidateOnboardingID, CandidateOnboardingTaskID: &updatedTask.ID, Action: "task_status_changed", FromStatus: &before.Status, ToStatus: &updatedTask.Status, Remarks: cmd.Remarks, ActorID: cmd.ActorID}); err != nil {
		s.logError("record candidate onboarding task event", err, serviceTenantIDField(cmd.TenantID), serviceStringField("candidate_onboarding_task_id", cmd.TaskID.String()))
	}
	if recalculated.OnboardingStatus == domain.OnboardStatusCompleted {
		status := domain.OnboardStatusCompleted
		if _, err := s.createCandidateOnboardingEvent(ctx, ports.CandidateOnboardingEventCommand{TenantID: cmd.TenantID, CandidateOnboardingID: updatedTask.CandidateOnboardingID, Action: "completed", ToStatus: &status, ActorID: cmd.ActorID}); err != nil {
			s.logError("record candidate onboarding completed event", err, serviceTenantIDField(cmd.TenantID), serviceStringField("candidate_onboarding_id", updatedTask.CandidateOnboardingID.String()))
		}
	}
	return s.GetCandidateOnboarding(ctx, cmd.TenantID, updatedTask.CandidateOnboardingID)
}

func (s *TenantService) DeleteCandidateOnboarding(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if tenantID == uuid.Nil || id == uuid.Nil {
		err := domain.ErrInvalidCandidateOnboardingID
		s.logError("validate delete candidate onboarding", err)
		return err
	}
	if _, err := s.candidateOnboardings.GetCandidateOnboarding(ctx, tenantID, id); err != nil {
		return err
	}
	if err := s.candidateOnboardings.DeleteCandidateOnboarding(ctx, tenantID, id, actorID); err != nil {
		s.logError("delete candidate onboarding", err, serviceTenantIDField(tenantID), serviceStringField("candidate_onboarding_id", id.String()))
		return err
	}
	return nil
}

func (s *TenantService) ListCandidateOnboardingEvents(ctx context.Context, tenantID uuid.UUID, candidateOnboardingID uuid.UUID) ([]*domain.CandidateOnboardingEvent, error) {
	if tenantID == uuid.Nil || candidateOnboardingID == uuid.Nil {
		err := domain.ErrInvalidCandidateOnboardingID
		s.logError("validate list candidate onboarding events", err)
		return nil, err
	}
	events, err := s.candidateOnboardings.ListCandidateOnboardingEvents(ctx, tenantID, candidateOnboardingID)
	if err != nil {
		s.logError("list candidate onboarding events", err, serviceTenantIDField(tenantID), serviceStringField("candidate_onboarding_id", candidateOnboardingID.String()))
		return nil, err
	}
	return events, nil
}

func (s *TenantService) resolveCandidateOnboardingWorkflow(ctx context.Context, cmd ports.StartCandidateOnboardingCommand) (*domain.OnboardingWorkflow, error) {
	if cmd.WorkflowID != nil && *cmd.WorkflowID != uuid.Nil {
		return s.onboardingWorkflows.GetOnboardingWorkflow(ctx, cmd.TenantID, *cmd.WorkflowID)
	}
	if workflow, err := s.candidateOnboardings.ResolveOnboardingWorkflowForCandidate(ctx, cmd.TenantID, cmd.CandidateID); err == nil {
		return workflow, nil
	}
	return s.candidateOnboardings.GetDefaultOnboardingWorkflow(ctx, cmd.TenantID)
}

func (s *TenantService) enrichCandidateOnboarding(ctx context.Context, item *domain.CandidateOnboarding) (*domain.CandidateOnboarding, error) {
	if item == nil {
		return nil, domain.ErrCandidateOnboardingNotFound
	}
	tasks, err := s.candidateOnboardings.ListCandidateOnboardingTasks(ctx, item.TenantID, item.ID)
	if err != nil {
		s.logError("list candidate onboarding tasks for detail", err, serviceTenantIDField(item.TenantID), serviceStringField("candidate_onboarding_id", item.ID.String()))
		return nil, err
	}
	events, err := s.candidateOnboardings.ListCandidateOnboardingEvents(ctx, item.TenantID, item.ID)
	if err != nil {
		s.logError("list candidate onboarding events for detail", err, serviceTenantIDField(item.TenantID), serviceStringField("candidate_onboarding_id", item.ID.String()))
		return nil, err
	}
	item.Tasks = tasks
	item.Events = events
	return item, nil
}

func (s *TenantService) createCandidateOnboardingEvent(ctx context.Context, cmd ports.CandidateOnboardingEventCommand) (*domain.CandidateOnboardingEvent, error) {
	event := &domain.CandidateOnboardingEvent{TenantID: cmd.TenantID, CandidateOnboardingID: cmd.CandidateOnboardingID, CandidateOnboardingTaskID: cmd.CandidateOnboardingTaskID, Action: cmd.Action, FromStatus: cmd.FromStatus, ToStatus: cmd.ToStatus, Remarks: cmd.Remarks, Metadata: cmd.Metadata}
	return s.candidateOnboardings.CreateCandidateOnboardingEvent(ctx, event, cmd.ActorID)
}
