package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) CreateInterviewRound(ctx context.Context, cmd ports.InterviewRoundCommand) (*domain.InterviewRound, error) {
	if _, err := s.GetCandidateApplication(ctx, cmd.TenantID, cmd.ApplicationID); err != nil {
		return nil, err
	}
	item, err := domain.NewInterviewRound(interviewRoundInput(cmd))
	if err != nil {
		s.logError("validate interview round create", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	result, err := s.candidates.CreateInterviewRound(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create interview round", err, serviceTenantIDField(cmd.TenantID), serviceStringField("candidate_application_id", cmd.ApplicationID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) ListInterviewRounds(ctx context.Context, filter domain.InterviewRoundFilter) (*domain.InterviewRoundPage, error) {
	if filter.TenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate interview round list tenant", err)
		return nil, err
	}
	if filter.ApplicationID != nil && *filter.ApplicationID == uuid.Nil {
		err := domain.ErrInvalidInterviewApplicationID
		s.logError("validate interview round list application", err, serviceTenantIDField(filter.TenantID))
		return nil, err
	}
	if _, err := domain.ValidateInterviewStatus(filter.Status); err != nil {
		s.logError("validate interview round list status", err, serviceTenantIDField(filter.TenantID))
		return nil, err
	}
	filter.Search = cleanStringPtr(filter.Search)
	filter.Status = cleanStringPtr(filter.Status)
	limit, offset := normalizeListWindow(filter.Limit, filter.Offset)
	filter.Limit = limit
	filter.Offset = offset
	items, err := s.candidates.ListInterviewRounds(ctx, filter)
	if err != nil {
		s.logError("list interview rounds", err, serviceTenantIDField(filter.TenantID))
		return nil, err
	}
	total, err := s.candidates.CountInterviewRounds(ctx, filter)
	if err != nil {
		s.logError("count interview rounds", err, serviceTenantIDField(filter.TenantID))
		return nil, err
	}
	page := &domain.InterviewRoundPage{Items: items, Total: total, Limit: limit, Offset: offset}
	if int64(offset)+int64(len(items)) < total {
		next := offset + limit
		page.NextOffset = &next
	}
	return page, nil
}

func (s *TenantService) GetInterviewRound(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.InterviewRound, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate interview round get tenant", err)
		return nil, err
	}
	if id == uuid.Nil {
		err := domain.ErrInvalidInterviewRoundID
		s.logError("validate interview round get id", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	item, err := s.candidates.GetInterviewRound(ctx, tenantID, id)
	if err != nil {
		s.logError("get interview round", err, serviceTenantIDField(tenantID), serviceStringField("interview_round_id", id.String()))
		return nil, err
	}
	return item, nil
}

func (s *TenantService) UpdateInterviewRound(ctx context.Context, cmd ports.InterviewRoundCommand) (*domain.InterviewRound, error) {
	if cmd.ID == uuid.Nil {
		err := domain.ErrInvalidInterviewRoundID
		s.logError("validate interview round update id", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	if _, err := s.GetInterviewRound(ctx, cmd.TenantID, cmd.ID); err != nil {
		return nil, err
	}
	if _, err := s.GetCandidateApplication(ctx, cmd.TenantID, cmd.ApplicationID); err != nil {
		return nil, err
	}
	item, err := domain.NewInterviewRound(interviewRoundInput(cmd))
	if err != nil {
		s.logError("validate interview round update", err, serviceTenantIDField(cmd.TenantID), serviceStringField("interview_round_id", cmd.ID.String()))
		return nil, err
	}
	item.ID = cmd.ID
	result, err := s.candidates.UpdateInterviewRound(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("update interview round", err, serviceTenantIDField(cmd.TenantID), serviceStringField("interview_round_id", cmd.ID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) UpdateInterviewRoundStatus(ctx context.Context, cmd ports.InterviewRoundStatusCommand) (*domain.InterviewRound, error) {
	if cmd.ID == uuid.Nil {
		err := domain.ErrInvalidInterviewRoundID
		s.logError("validate interview round status id", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	if _, err := s.GetInterviewRound(ctx, cmd.TenantID, cmd.ID); err != nil {
		return nil, err
	}
	status, err := domain.ValidateInterviewStatus(&cmd.Status)
	if err != nil {
		s.logError("validate interview round status", err, serviceTenantIDField(cmd.TenantID), serviceStringField("interview_round_id", cmd.ID.String()))
		return nil, err
	}
	if status == nil {
		err := domain.ErrInvalidInterviewStatus
		s.logError("validate interview round status empty", err, serviceTenantIDField(cmd.TenantID), serviceStringField("interview_round_id", cmd.ID.String()))
		return nil, err
	}
	if cmd.Score != nil && (*cmd.Score < 0 || *cmd.Score > 5) {
		err := domain.ErrInvalidInterviewScore
		s.logError("validate interview round score", err, serviceTenantIDField(cmd.TenantID), serviceStringField("interview_round_id", cmd.ID.String()))
		return nil, err
	}
	if _, err := domain.ValidateInterviewDecision(cmd.Decision); err != nil {
		s.logError("validate interview round decision", err, serviceTenantIDField(cmd.TenantID), serviceStringField("interview_round_id", cmd.ID.String()))
		return nil, err
	}
	result, err := s.candidates.UpdateInterviewRoundStatus(ctx, cmd.TenantID, cmd.ID, *status, cmd.Remarks, cmd.Feedback, cmd.Score, cmd.Decision, cmd.CompletedAt, cmd.ActorID)
	if err != nil {
		s.logError("update interview round status", err, serviceTenantIDField(cmd.TenantID), serviceStringField("interview_round_id", cmd.ID.String()), serviceStringField("status", *status))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) DeleteInterviewRound(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if _, err := s.GetInterviewRound(ctx, tenantID, id); err != nil {
		return err
	}
	if err := s.candidates.DeleteInterviewRound(ctx, tenantID, id, actorID); err != nil {
		s.logError("delete interview round", err, serviceTenantIDField(tenantID), serviceStringField("interview_round_id", id.String()))
		return err
	}
	return nil
}

func interviewRoundInput(cmd ports.InterviewRoundCommand) domain.InterviewRoundInput {
	return domain.InterviewRoundInput{TenantID: cmd.TenantID, ApplicationID: cmd.ApplicationID, RoundName: cmd.RoundName, RoundNumber: cmd.RoundNumber, ScheduledDate: cmd.ScheduledDate, DurationMinutes: cmd.DurationMinutes, InterviewerUserID: cmd.InterviewerUserID, Mode: cmd.Mode, MeetingLink: cmd.MeetingLink, Location: cmd.Location, Status: cmd.Status, Remarks: cmd.Remarks, Timezone: cmd.Timezone, Feedback: cmd.Feedback, Score: cmd.Score, Decision: cmd.Decision, CompletedAt: cmd.CompletedAt}
}
