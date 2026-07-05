package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) CreateWorkLog(ctx context.Context, cmd ports.WorkLogCommand) (*domain.WorkLog, error) {
	input, err := s.prepareWorkLogCommand(ctx, cmd, nil)
	if err != nil {
		return nil, err
	}
	item, err := domain.NewWorkLog(input)
	if err != nil {
		s.logError("validate work log create", err, serviceTenantIDField(cmd.TenantID), serviceStringField("engagement_id", cmd.EngagementID.String()))
		return nil, err
	}
	if err := s.ensureWorkLogWithinBudget(ctx, item, nil); err != nil {
		return nil, err
	}
	result, err := s.workLogs.CreateWorkLog(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create work log", err, serviceTenantIDField(cmd.TenantID), serviceStringField("engagement_id", item.EngagementID.String()))
		return nil, err
	}
	s.log.Info().Str("tenant_id", result.TenantID.String()).Str("work_log_id", result.ID.String()).Msg("hrms: work log created")
	return result, nil
}

func (s *TenantService) UpdateWorkLog(ctx context.Context, cmd ports.WorkLogCommand) (*domain.WorkLog, error) {
	if cmd.ID == uuid.Nil {
		err := domain.ErrInvalidWorkLogID
		s.logError("validate work log update id", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	existing, err := s.GetWorkLog(ctx, cmd.TenantID, cmd.ID)
	if err != nil {
		return nil, err
	}
	if existing.Status == domain.WorkLogStatusApproved {
		err := domain.ErrInvalidWorkLogWorkflow
		s.logError("validate approved work log update", err, serviceTenantIDField(cmd.TenantID), serviceStringField("work_log_id", cmd.ID.String()))
		return nil, err
	}
	if cmd.EngagementID == uuid.Nil {
		cmd.EngagementID = existing.EngagementID
	}
	if cmd.Status == "" {
		cmd.Status = existing.Status
	}
	input, err := s.prepareWorkLogCommand(ctx, cmd, existing)
	if err != nil {
		return nil, err
	}
	item, err := domain.NewWorkLog(input)
	if err != nil {
		s.logError("validate work log update", err, serviceTenantIDField(cmd.TenantID), serviceStringField("work_log_id", cmd.ID.String()))
		return nil, err
	}
	item.ID = cmd.ID
	if err := s.ensureWorkLogWithinBudget(ctx, item, &cmd.ID); err != nil {
		return nil, err
	}
	result, err := s.workLogs.UpdateWorkLog(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("update work log", err, serviceTenantIDField(cmd.TenantID), serviceStringField("work_log_id", cmd.ID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) SubmitWorkLog(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) (*domain.WorkLog, error) {
	item, err := s.GetWorkLog(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}
	if item.Status != domain.WorkLogStatusDraft && item.Status != domain.WorkLogStatusRejected {
		err := domain.ErrInvalidWorkLogWorkflow
		s.logError("validate work log submit state", err, serviceTenantIDField(tenantID), serviceStringField("work_log_id", id.String()))
		return nil, err
	}
	now := time.Now().UTC()
	item.Status = domain.WorkLogStatusSubmitted
	item.SubmittedAt = &now
	item.SubmittedBy = actorID
	item.ReviewedAt = nil
	item.ReviewedBy = nil
	item.ReviewComment = nil
	if err := s.ensureWorkLogWithinBudget(ctx, item, &id); err != nil {
		return nil, err
	}
	result, err := s.workLogs.UpdateWorkLogStatus(ctx, item, actorID)
	if err != nil {
		s.logError("submit work log", err, serviceTenantIDField(tenantID), serviceStringField("work_log_id", id.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) ReviewWorkLog(ctx context.Context, cmd ports.WorkLogStatusCommand) (*domain.WorkLog, error) {
	if cmd.TenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate work log review tenant", err)
		return nil, err
	}
	if cmd.WorkLogID == uuid.Nil {
		err := domain.ErrInvalidWorkLogID
		s.logError("validate work log review id", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	status, err := domain.ValidateWorkLogStatus(cmd.Status)
	if err != nil {
		s.logError("validate work log review status", err, serviceTenantIDField(cmd.TenantID), serviceStringField("work_log_id", cmd.WorkLogID.String()))
		return nil, err
	}
	if status != domain.WorkLogStatusApproved && status != domain.WorkLogStatusRejected {
		err := domain.ErrInvalidWorkLogWorkflow
		s.logError("validate work log review action", err, serviceTenantIDField(cmd.TenantID), serviceStringField("work_log_id", cmd.WorkLogID.String()))
		return nil, err
	}
	item, err := s.GetWorkLog(ctx, cmd.TenantID, cmd.WorkLogID)
	if err != nil {
		return nil, err
	}
	if item.Status != domain.WorkLogStatusSubmitted {
		err := domain.ErrInvalidWorkLogWorkflow
		s.logError("validate work log review state", err, serviceTenantIDField(cmd.TenantID), serviceStringField("work_log_id", cmd.WorkLogID.String()))
		return nil, err
	}
	now := time.Now().UTC()
	item.Status = status
	item.ReviewedAt = &now
	item.ReviewedBy = cmd.ActorID
	item.ReviewComment = cmd.ReviewComment
	if status == domain.WorkLogStatusApproved {
		if err := s.ensureWorkLogWithinBudget(ctx, item, &cmd.WorkLogID); err != nil {
			return nil, err
		}
	}
	result, err := s.workLogs.UpdateWorkLogStatus(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("review work log", err, serviceTenantIDField(cmd.TenantID), serviceStringField("work_log_id", cmd.WorkLogID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) GetWorkLog(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.WorkLog, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate work log get tenant", err)
		return nil, err
	}
	if id == uuid.Nil {
		err := domain.ErrInvalidWorkLogID
		s.logError("validate work log get id", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	result, err := s.workLogs.GetWorkLog(ctx, tenantID, id)
	if err != nil {
		s.logError("get work log", err, serviceTenantIDField(tenantID), serviceStringField("work_log_id", id.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) ListWorkLogs(ctx context.Context, filter domain.WorkLogFilter) ([]*domain.WorkLogListItem, error) {
	if filter.TenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate work log list tenant", err)
		return nil, err
	}
	result, err := s.workLogs.ListWorkLogs(ctx, filter)
	if err != nil {
		s.logError("list work logs", err, serviceTenantIDField(filter.TenantID))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) ListWorkLogRollups(ctx context.Context, filter domain.WorkLogFilter) ([]*domain.WorkLogRollup, error) {
	if filter.TenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate work log rollup tenant", err)
		return nil, err
	}
	result, err := s.workLogs.ListWorkLogRollups(ctx, filter)
	if err != nil {
		s.logError("list work log rollups", err, serviceTenantIDField(filter.TenantID))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) DeleteWorkLog(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate work log delete tenant", err)
		return err
	}
	if id == uuid.Nil {
		err := domain.ErrInvalidWorkLogID
		s.logError("validate work log delete id", err, serviceTenantIDField(tenantID))
		return err
	}
	item, err := s.GetWorkLog(ctx, tenantID, id)
	if err != nil {
		return err
	}
	if item.Status == domain.WorkLogStatusApproved {
		err := domain.ErrInvalidWorkLogWorkflow
		s.logError("validate approved work log delete", err, serviceTenantIDField(tenantID), serviceStringField("work_log_id", id.String()))
		return err
	}
	if err := s.workLogs.DeleteWorkLog(ctx, tenantID, id, actorID); err != nil {
		s.logError("delete work log", err, serviceTenantIDField(tenantID), serviceStringField("work_log_id", id.String()))
		return err
	}
	return nil
}

func (s *TenantService) prepareWorkLogCommand(ctx context.Context, cmd ports.WorkLogCommand, existing *domain.WorkLog) (domain.WorkLogInput, error) {
	if cmd.TenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate work log tenant", err)
		return domain.WorkLogInput{}, err
	}
	if cmd.EngagementID == uuid.Nil {
		err := domain.ErrInvalidEngagementID
		s.logError("validate work log engagement", err, serviceTenantIDField(cmd.TenantID))
		return domain.WorkLogInput{}, err
	}
	engagement, err := s.GetEngagement(ctx, cmd.TenantID, cmd.EngagementID)
	if err != nil {
		return domain.WorkLogInput{}, err
	}
	cmd.WorkerProfileID = engagement.WorkerProfileID
	logDate, err := parseRequiredWorkLogDate(cmd.LogDate)
	if err != nil {
		s.logError("parse work log date", err, serviceTenantIDField(cmd.TenantID), serviceStringField("engagement_id", cmd.EngagementID.String()))
		return domain.WorkLogInput{}, err
	}
	status := cmd.Status
	if status == "" {
		status = domain.WorkLogStatusDraft
	}
	var submittedAt *time.Time
	var submittedBy *uuid.UUID
	var reviewedAt *time.Time
	var reviewedBy *uuid.UUID
	var reviewComment *string
	if existing != nil {
		submittedAt = existing.SubmittedAt
		submittedBy = existing.SubmittedBy
		reviewedAt = existing.ReviewedAt
		reviewedBy = existing.ReviewedBy
		reviewComment = existing.ReviewComment
	}
	now := time.Now().UTC()
	if status == domain.WorkLogStatusSubmitted && submittedAt == nil {
		submittedAt = &now
		submittedBy = cmd.ActorID
		reviewedAt = nil
		reviewedBy = nil
		reviewComment = nil
	}
	if status == domain.WorkLogStatusDraft {
		submittedAt = nil
		submittedBy = nil
		reviewedAt = nil
		reviewedBy = nil
		reviewComment = nil
	}
	return domain.WorkLogInput{
		TenantID:             cmd.TenantID,
		EngagementID:         cmd.EngagementID,
		WorkerProfileID:      cmd.WorkerProfileID,
		LogDate:              logDate,
		HoursWorked:          cmd.HoursWorked,
		BillableHours:        cmd.BillableHours,
		WorkSummary:          cmd.WorkSummary,
		DeliverableReference: cmd.DeliverableReference,
		Status:               status,
		SubmittedAt:          submittedAt,
		SubmittedBy:          submittedBy,
		ReviewedAt:           reviewedAt,
		ReviewedBy:           reviewedBy,
		ReviewComment:        reviewComment,
		Metadata:             cmd.Metadata,
	}, nil
}

func (s *TenantService) ensureWorkLogWithinBudget(ctx context.Context, item *domain.WorkLog, excludeID *uuid.UUID) error {
	if item.Status != domain.WorkLogStatusSubmitted && item.Status != domain.WorkLogStatusApproved {
		return nil
	}
	usage, err := s.workLogs.GetWorkLogBudgetUsage(ctx, item.TenantID, item.EngagementID, excludeID)
	if err != nil {
		s.logError("get work log budget usage", err, serviceTenantIDField(item.TenantID), serviceStringField("engagement_id", item.EngagementID.String()))
		return err
	}
	if usage.HoursBudget == nil {
		return nil
	}
	if usage.UsedHours+item.HoursWorked > *usage.HoursBudget+0.0001 {
		err := domain.ErrWorkLogBudgetExceeded
		s.logError("validate work log budget", err, serviceTenantIDField(item.TenantID), serviceStringField("engagement_id", item.EngagementID.String()))
		return err
	}
	return nil
}

func parseRequiredWorkLogDate(value string) (*time.Time, error) {
	parsed, err := parseEngagementDate(value)
	if err != nil {
		return nil, err
	}
	if parsed == nil {
		return nil, domain.ErrInvalidWorkLogDate
	}
	return parsed, nil
}
