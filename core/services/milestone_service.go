package services

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) CreateProjectMilestone(ctx context.Context, cmd ports.ProjectMilestoneCommand) (*domain.ProjectMilestone, error) {
	input, err := s.prepareProjectMilestoneCommand(ctx, cmd, nil)
	if err != nil {
		return nil, err
	}
	item, err := domain.NewProjectMilestone(input)
	if err != nil {
		s.logError("validate project milestone create", err, serviceTenantIDField(cmd.TenantID), serviceStringField("project_id", cmd.ProjectID.String()))
		return nil, err
	}
	var result *domain.ProjectMilestone
	err = s.system.RunAsSystem(ctx, func(txCtx context.Context) error {
		created, err := s.projects.CreateProjectMilestone(txCtx, item, cmd.ActorID)
		if err != nil {
			return err
		}
		event := milestoneEvent(created, domain.MilestoneEventCreated, nil, &created.Status, nil, cmd.ActorID, json.RawMessage(`{}`))
		if _, err := s.projects.CreateProjectMilestoneEvent(txCtx, event); err != nil {
			return err
		}
		result = created
		return nil
	})
	if err != nil {
		s.logError("create project milestone", err, serviceTenantIDField(cmd.TenantID), serviceStringField("project_id", item.ProjectID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) UpdateProjectMilestone(ctx context.Context, cmd ports.ProjectMilestoneCommand) (*domain.ProjectMilestone, error) {
	if cmd.ID == uuid.Nil {
		err := domain.ErrInvalidMilestoneID
		s.logError("validate project milestone update id", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	existing, err := s.GetProjectMilestone(ctx, cmd.TenantID, cmd.ID)
	if err != nil {
		return nil, err
	}
	if existing.Status == domain.MilestoneStatusAccepted || existing.Status == domain.MilestoneStatusCancelled {
		err := domain.ErrInvalidMilestoneWorkflow
		s.logError("validate project milestone update state", err, serviceTenantIDField(cmd.TenantID), serviceStringField("milestone_id", cmd.ID.String()))
		return nil, err
	}
	if cmd.ProjectID == uuid.Nil {
		cmd.ProjectID = existing.ProjectID
	}
	if cmd.Status == "" {
		cmd.Status = existing.Status
	}
	input, err := s.prepareProjectMilestoneCommand(ctx, cmd, existing)
	if err != nil {
		return nil, err
	}
	item, err := domain.NewProjectMilestone(input)
	if err != nil {
		s.logError("validate project milestone update", err, serviceTenantIDField(cmd.TenantID), serviceStringField("milestone_id", cmd.ID.String()))
		return nil, err
	}
	item.ID = cmd.ID
	var result *domain.ProjectMilestone
	err = s.system.RunAsSystem(ctx, func(txCtx context.Context) error {
		updated, err := s.projects.UpdateProjectMilestone(txCtx, item, cmd.ActorID)
		if err != nil {
			return err
		}
		event := milestoneEvent(updated, domain.MilestoneEventUpdated, &existing.Status, &updated.Status, nil, cmd.ActorID, json.RawMessage(`{}`))
		if _, err := s.projects.CreateProjectMilestoneEvent(txCtx, event); err != nil {
			return err
		}
		result = updated
		return nil
	})
	if err != nil {
		s.logError("update project milestone", err, serviceTenantIDField(cmd.TenantID), serviceStringField("milestone_id", cmd.ID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) SubmitProjectMilestone(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) (*domain.ProjectMilestone, error) {
	item, err := s.GetProjectMilestone(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}
	if item.Status != domain.MilestoneStatusDraft && item.Status != domain.MilestoneStatusOpen && item.Status != domain.MilestoneStatusRejected {
		err := domain.ErrInvalidMilestoneWorkflow
		s.logError("validate project milestone submit state", err, serviceTenantIDField(tenantID), serviceStringField("milestone_id", id.String()))
		return nil, err
	}
	fromStatus := item.Status
	now := time.Now().UTC()
	item.Status = domain.MilestoneStatusSubmitted
	item.SubmittedAt = &now
	item.SubmittedBy = actorID
	item.AcceptedAt = nil
	item.AcceptedBy = nil
	item.RejectedAt = nil
	item.RejectedBy = nil
	item.ReviewComment = nil
	var result *domain.ProjectMilestone
	err = s.system.RunAsSystem(ctx, func(txCtx context.Context) error {
		updated, err := s.projects.UpdateProjectMilestoneStatus(txCtx, item, actorID)
		if err != nil {
			return err
		}
		event := milestoneEvent(updated, domain.MilestoneEventSubmitted, &fromStatus, &updated.Status, nil, actorID, json.RawMessage(`{}`))
		if _, err := s.projects.CreateProjectMilestoneEvent(txCtx, event); err != nil {
			return err
		}
		result = updated
		return nil
	})
	if err != nil {
		s.logError("submit project milestone", err, serviceTenantIDField(tenantID), serviceStringField("milestone_id", id.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) ReviewProjectMilestone(ctx context.Context, cmd ports.ProjectMilestoneStatusCommand) (*domain.ProjectMilestone, error) {
	if cmd.TenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate project milestone review tenant", err)
		return nil, err
	}
	if cmd.MilestoneID == uuid.Nil {
		err := domain.ErrInvalidMilestoneID
		s.logError("validate project milestone review id", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	status, err := domain.ValidateMilestoneStatus(cmd.Status)
	if err != nil {
		s.logError("validate project milestone review status", err, serviceTenantIDField(cmd.TenantID), serviceStringField("milestone_id", cmd.MilestoneID.String()))
		return nil, err
	}
	if status != domain.MilestoneStatusAccepted && status != domain.MilestoneStatusRejected {
		err := domain.ErrInvalidMilestoneWorkflow
		s.logError("validate project milestone review action", err, serviceTenantIDField(cmd.TenantID), serviceStringField("milestone_id", cmd.MilestoneID.String()))
		return nil, err
	}
	item, err := s.GetProjectMilestone(ctx, cmd.TenantID, cmd.MilestoneID)
	if err != nil {
		return nil, err
	}
	if item.Status != domain.MilestoneStatusSubmitted {
		err := domain.ErrInvalidMilestoneWorkflow
		s.logError("validate project milestone review state", err, serviceTenantIDField(cmd.TenantID), serviceStringField("milestone_id", cmd.MilestoneID.String()))
		return nil, err
	}
	fromStatus := item.Status
	now := time.Now().UTC()
	item.Status = status
	item.ReviewComment = cmd.ReviewComment
	if status == domain.MilestoneStatusAccepted {
		item.AcceptedAt = &now
		item.AcceptedBy = cmd.ActorID
		item.RejectedAt = nil
		item.RejectedBy = nil
	} else {
		item.RejectedAt = &now
		item.RejectedBy = cmd.ActorID
		item.AcceptedAt = nil
		item.AcceptedBy = nil
	}
	eventType := domain.MilestoneEventAccepted
	if status == domain.MilestoneStatusRejected {
		eventType = domain.MilestoneEventRejected
	}
	var result *domain.ProjectMilestone
	err = s.system.RunAsSystem(ctx, func(txCtx context.Context) error {
		updated, err := s.projects.UpdateProjectMilestoneStatus(txCtx, item, cmd.ActorID)
		if err != nil {
			return err
		}
		event := milestoneEvent(updated, eventType, &fromStatus, &updated.Status, cmd.ReviewComment, cmd.ActorID, json.RawMessage(`{}`))
		if _, err := s.projects.CreateProjectMilestoneEvent(txCtx, event); err != nil {
			return err
		}
		result = updated
		return nil
	})
	if err != nil {
		s.logError("review project milestone", err, serviceTenantIDField(cmd.TenantID), serviceStringField("milestone_id", cmd.MilestoneID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) GetProjectMilestone(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.ProjectMilestone, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate project milestone get tenant", err)
		return nil, err
	}
	if id == uuid.Nil {
		err := domain.ErrInvalidMilestoneID
		s.logError("validate project milestone get id", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	result, err := s.projects.GetProjectMilestone(ctx, tenantID, id)
	if err != nil {
		s.logError("get project milestone", err, serviceTenantIDField(tenantID), serviceStringField("milestone_id", id.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) ListProjectMilestones(ctx context.Context, filter domain.ProjectMilestoneFilter) ([]*domain.ProjectMilestoneListItem, error) {
	if filter.TenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate project milestone list tenant", err)
		return nil, err
	}
	result, err := s.projects.ListProjectMilestones(ctx, filter)
	if err != nil {
		s.logError("list project milestones", err, serviceTenantIDField(filter.TenantID))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) DeleteProjectMilestone(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	item, err := s.GetProjectMilestone(ctx, tenantID, id)
	if err != nil {
		return err
	}
	if item.Status == domain.MilestoneStatusAccepted {
		err := domain.ErrInvalidMilestoneWorkflow
		s.logError("validate accepted project milestone delete", err, serviceTenantIDField(tenantID), serviceStringField("milestone_id", id.String()))
		return err
	}
	if err := s.projects.DeleteProjectMilestone(ctx, tenantID, id, actorID); err != nil {
		s.logError("delete project milestone", err, serviceTenantIDField(tenantID), serviceStringField("milestone_id", id.String()))
		return err
	}
	event := milestoneEvent(item, domain.MilestoneEventCancelled, &item.Status, stringPtr(domain.MilestoneStatusCancelled), nil, actorID, json.RawMessage(`{}`))
	if _, err := s.projects.CreateProjectMilestoneEvent(ctx, event); err != nil {
		s.logError("create project milestone cancel event", err, serviceTenantIDField(tenantID), serviceStringField("milestone_id", id.String()))
		return err
	}
	return nil
}

func (s *TenantService) ListProjectMilestoneEvents(ctx context.Context, tenantID uuid.UUID, milestoneID uuid.UUID) ([]*domain.ProjectMilestoneEvent, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate project milestone events tenant", err)
		return nil, err
	}
	if milestoneID == uuid.Nil {
		err := domain.ErrInvalidMilestoneID
		s.logError("validate project milestone events id", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	if _, err := s.GetProjectMilestone(ctx, tenantID, milestoneID); err != nil {
		return nil, err
	}
	result, err := s.projects.ListProjectMilestoneEvents(ctx, tenantID, milestoneID)
	if err != nil {
		s.logError("list project milestone events", err, serviceTenantIDField(tenantID), serviceStringField("milestone_id", milestoneID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) prepareProjectMilestoneCommand(ctx context.Context, cmd ports.ProjectMilestoneCommand, existing *domain.ProjectMilestone) (domain.ProjectMilestoneInput, error) {
	if cmd.TenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate project milestone tenant", err)
		return domain.ProjectMilestoneInput{}, err
	}
	if cmd.ProjectID == uuid.Nil {
		err := domain.ErrInvalidProjectID
		s.logError("validate project milestone project", err, serviceTenantIDField(cmd.TenantID))
		return domain.ProjectMilestoneInput{}, err
	}
	if _, err := s.GetProject(ctx, cmd.TenantID, cmd.ProjectID); err != nil {
		return domain.ProjectMilestoneInput{}, err
	}
	if cmd.EngagementID != nil {
		if _, err := s.GetEngagement(ctx, cmd.TenantID, *cmd.EngagementID); err != nil {
			return domain.ProjectMilestoneInput{}, err
		}
	}
	dueDate, err := parseEngagementDate(cmd.DueDate)
	if err != nil {
		s.logError("parse project milestone due date", err, serviceTenantIDField(cmd.TenantID), serviceStringField("project_id", cmd.ProjectID.String()))
		return domain.ProjectMilestoneInput{}, err
	}
	status := cmd.Status
	var submittedAt *time.Time
	var submittedBy *uuid.UUID
	var acceptedAt *time.Time
	var acceptedBy *uuid.UUID
	var rejectedAt *time.Time
	var rejectedBy *uuid.UUID
	var reviewComment *string
	if existing != nil {
		submittedAt = existing.SubmittedAt
		submittedBy = existing.SubmittedBy
		acceptedAt = existing.AcceptedAt
		acceptedBy = existing.AcceptedBy
		rejectedAt = existing.RejectedAt
		rejectedBy = existing.RejectedBy
		reviewComment = existing.ReviewComment
		if status == "" {
			status = existing.Status
		}
		if cmd.CurrencyCode == "" {
			cmd.CurrencyCode = existing.CurrencyCode
		}
	}
	return domain.ProjectMilestoneInput{
		TenantID:           cmd.TenantID,
		ProjectID:          cmd.ProjectID,
		EngagementID:       cmd.EngagementID,
		MilestoneCode:      cmd.MilestoneCode,
		Title:              cmd.Title,
		Description:        cmd.Description,
		AcceptanceCriteria: cmd.AcceptanceCriteria,
		DueDate:            dueDate,
		Status:             status,
		Amount:             cmd.Amount,
		CurrencyCode:       cmd.CurrencyCode,
		PaymentTrigger:     cmd.PaymentTrigger,
		SubmittedAt:        submittedAt,
		SubmittedBy:        submittedBy,
		AcceptedAt:         acceptedAt,
		AcceptedBy:         acceptedBy,
		RejectedAt:         rejectedAt,
		RejectedBy:         rejectedBy,
		ReviewComment:      reviewComment,
		Notes:              cmd.Notes,
		Metadata:           cmd.Metadata,
	}, nil
}

func milestoneEvent(item *domain.ProjectMilestone, eventType string, fromStatus *string, toStatus *string, comment *string, actorID *uuid.UUID, metadata json.RawMessage) *domain.ProjectMilestoneEvent {
	return &domain.ProjectMilestoneEvent{
		TenantID:    item.TenantID,
		ProjectID:   item.ProjectID,
		MilestoneID: item.ID,
		EventType:   eventType,
		FromStatus:  fromStatus,
		ToStatus:    toStatus,
		Comment:     comment,
		ActorID:     actorID,
		Metadata:    metadata,
	}
}
