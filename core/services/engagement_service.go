package services

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) CreateEngagement(ctx context.Context, cmd ports.EngagementCommand) (*domain.Engagement, error) {
	input, err := s.prepareEngagementCommand(ctx, cmd)
	if err != nil {
		return nil, err
	}
	item, err := domain.NewEngagement(input)
	if err != nil {
		s.logError("validate engagement create", err, serviceTenantIDField(cmd.TenantID), serviceStringField("worker_profile_id", cmd.WorkerProfileID.String()))
		return nil, err
	}
	result, err := s.engagements.CreateEngagement(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create engagement", err, serviceTenantIDField(cmd.TenantID), serviceStringField("worker_profile_id", item.WorkerProfileID.String()))
		return nil, err
	}
	s.log.Info().Str("tenant_id", result.TenantID.String()).Str("engagement_id", result.ID.String()).Msg("hrms: engagement created")
	return result, nil
}

func (s *TenantService) UpdateEngagement(ctx context.Context, cmd ports.EngagementCommand) (*domain.Engagement, error) {
	if cmd.ID == uuid.Nil {
		err := domain.ErrInvalidEngagementID
		s.logError("validate engagement update id", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	if _, err := s.GetEngagement(ctx, cmd.TenantID, cmd.ID); err != nil {
		return nil, err
	}
	input, err := s.prepareEngagementCommand(ctx, cmd)
	if err != nil {
		return nil, err
	}
	item, err := domain.NewEngagement(input)
	if err != nil {
		s.logError("validate engagement update", err, serviceTenantIDField(cmd.TenantID), serviceStringField("engagement_id", cmd.ID.String()))
		return nil, err
	}
	item.ID = cmd.ID
	result, err := s.engagements.UpdateEngagement(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("update engagement", err, serviceTenantIDField(cmd.TenantID), serviceStringField("engagement_id", cmd.ID.String()))
		return nil, err
	}
	s.log.Info().Str("tenant_id", result.TenantID.String()).Str("engagement_id", result.ID.String()).Msg("hrms: engagement updated")
	return result, nil
}

func (s *TenantService) UpdateEngagementStatus(ctx context.Context, cmd ports.EngagementStatusCommand) (*domain.Engagement, error) {
	if cmd.TenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate engagement status tenant", err)
		return nil, err
	}
	if cmd.EngagementID == uuid.Nil {
		err := domain.ErrInvalidEngagementID
		s.logError("validate engagement status id", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	status, err := domain.ValidateEngagementStatus(cmd.Status)
	if err != nil {
		s.logError("validate engagement status", err, serviceTenantIDField(cmd.TenantID), serviceStringField("engagement_id", cmd.EngagementID.String()))
		return nil, err
	}
	terminatedAt, err := parseEngagementTimestamp(cmd.TerminatedAt)
	if err != nil {
		s.logError("parse engagement terminated at", err, serviceTenantIDField(cmd.TenantID), serviceStringField("engagement_id", cmd.EngagementID.String()))
		return nil, err
	}
	terminationReason := cmd.TerminationReason
	if status == domain.EngagementStatusTerminated && terminatedAt == nil {
		now := time.Now().UTC()
		terminatedAt = &now
	}
	if status != domain.EngagementStatusTerminated {
		terminatedAt = nil
		terminationReason = nil
	}
	result, err := s.engagements.UpdateEngagementStatus(ctx, cmd.TenantID, cmd.EngagementID, status, terminationReason, terminatedAt, cmd.ActorID)
	if err != nil {
		s.logError("update engagement status", err, serviceTenantIDField(cmd.TenantID), serviceStringField("engagement_id", cmd.EngagementID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) GetEngagement(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.Engagement, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate engagement get tenant", err)
		return nil, err
	}
	if id == uuid.Nil {
		err := domain.ErrInvalidEngagementID
		s.logError("validate engagement get id", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	result, err := s.engagements.GetEngagement(ctx, tenantID, id)
	if err != nil {
		s.logError("get engagement", err, serviceTenantIDField(tenantID), serviceStringField("engagement_id", id.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) ListEngagements(ctx context.Context, filter domain.EngagementFilter) ([]*domain.EngagementListItem, error) {
	if filter.TenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate engagement list tenant", err)
		return nil, err
	}
	result, err := s.engagements.ListEngagements(ctx, filter)
	if err != nil {
		s.logError("list engagements", err, serviceTenantIDField(filter.TenantID))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) DeleteEngagement(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate engagement delete tenant", err)
		return err
	}
	if id == uuid.Nil {
		err := domain.ErrInvalidEngagementID
		s.logError("validate engagement delete id", err, serviceTenantIDField(tenantID))
		return err
	}
	if err := s.engagements.DeleteEngagement(ctx, tenantID, id, actorID); err != nil {
		s.logError("delete engagement", err, serviceTenantIDField(tenantID), serviceStringField("engagement_id", id.String()))
		return err
	}
	return nil
}

func (s *TenantService) prepareEngagementCommand(ctx context.Context, cmd ports.EngagementCommand) (domain.EngagementInput, error) {
	if cmd.TenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate engagement tenant", err)
		return domain.EngagementInput{}, err
	}
	if cmd.WorkerProfileID == uuid.Nil {
		err := domain.ErrInvalidWorkerProfileID
		s.logError("validate engagement worker profile", err, serviceTenantIDField(cmd.TenantID))
		return domain.EngagementInput{}, err
	}
	worker, err := s.GetWorkerProfile(ctx, cmd.TenantID, cmd.WorkerProfileID)
	if err != nil {
		return domain.EngagementInput{}, err
	}
	if cmd.BranchID == nil {
		cmd.BranchID = worker.BranchID
	}
	if cmd.DepartmentID == nil {
		cmd.DepartmentID = worker.DepartmentID
	}
	if cmd.ReportingManagerID == nil {
		cmd.ReportingManagerID = worker.ReportingManagerID
	}
	if cmd.BranchID != nil && *cmd.BranchID != uuid.Nil {
		if _, err := s.GetBranch(ctx, cmd.TenantID, *cmd.BranchID); err != nil {
			return domain.EngagementInput{}, err
		}
	}
	if cmd.DepartmentID != nil && *cmd.DepartmentID != uuid.Nil {
		if _, err := s.GetDepartment(ctx, cmd.TenantID, *cmd.DepartmentID); err != nil {
			return domain.EngagementInput{}, err
		}
	}
	startDate, err := parseEngagementDate(cmd.StartDate)
	if err != nil {
		s.logError("parse engagement start date", err, serviceTenantIDField(cmd.TenantID))
		return domain.EngagementInput{}, err
	}
	endDate, err := parseEngagementDate(cmd.EndDate)
	if err != nil {
		s.logError("parse engagement end date", err, serviceTenantIDField(cmd.TenantID))
		return domain.EngagementInput{}, err
	}
	renewalDueDate, err := parseEngagementDate(cmd.RenewalDueDate)
	if err != nil {
		s.logError("parse engagement renewal due date", err, serviceTenantIDField(cmd.TenantID))
		return domain.EngagementInput{}, err
	}
	terminatedAt, err := parseEngagementTimestamp(cmd.TerminatedAt)
	if err != nil {
		s.logError("parse engagement terminated at", err, serviceTenantIDField(cmd.TenantID))
		return domain.EngagementInput{}, err
	}
	return domain.EngagementInput{
		TenantID:           cmd.TenantID,
		WorkerProfileID:    cmd.WorkerProfileID,
		EngagementCode:     cmd.EngagementCode,
		Title:              cmd.Title,
		Description:        cmd.Description,
		EngagementType:     cmd.EngagementType,
		Status:             cmd.Status,
		StartDate:          startDate,
		EndDate:            endDate,
		HoursBudget:        cmd.HoursBudget,
		RateAmount:         cmd.RateAmount,
		CurrencyCode:       cmd.CurrencyCode,
		RateUnit:           cmd.RateUnit,
		BranchID:           cmd.BranchID,
		DepartmentID:       cmd.DepartmentID,
		ReportingManagerID: cmd.ReportingManagerID,
		ProjectLabel:       cmd.ProjectLabel,
		ProjectCode:        cmd.ProjectCode,
		CostCenter:         cmd.CostCenter,
		RenewalDueDate:     renewalDueDate,
		RenewalStatus:      cmd.RenewalStatus,
		TerminationReason:  cmd.TerminationReason,
		TerminatedAt:       terminatedAt,
		Notes:              cmd.Notes,
		Metadata:           cmd.Metadata,
	}, nil
}

func parseEngagementDate(value string) (*time.Time, error) {
	clean := strings.TrimSpace(value)
	if clean == "" {
		return nil, nil
	}
	parsed, err := time.Parse("2006-01-02", clean)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

func parseEngagementTimestamp(value string) (*time.Time, error) {
	clean := strings.TrimSpace(value)
	if clean == "" {
		return nil, nil
	}
	if parsed, err := time.Parse(time.RFC3339, clean); err == nil {
		utc := parsed.UTC()
		return &utc, nil
	}
	parsed, err := time.Parse("2006-01-02", clean)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}
