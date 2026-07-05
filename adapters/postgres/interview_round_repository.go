package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func (s *Store) CreateInterviewRound(ctx context.Context, item *domain.InterviewRound, actorID *uuid.UUID) (*domain.InterviewRound, error) {
	row, err := s.getQueries(ctx).CreateInterviewRound(ctx, interviewRoundCreateParams(item, actorID))
	if err != nil {
		return nil, s.logDBError(ctx, "create interview round", fmt.Errorf("hrms: create interview round: %w", err), tenantIDField(item.TenantID), stringField("candidate_application_id", item.ApplicationID.String()))
	}
	return mapInterviewRound(row), nil
}

func (s *Store) ListInterviewRounds(ctx context.Context, filter domain.InterviewRoundFilter) ([]*domain.InterviewRound, error) {
	rows, err := s.getQueries(ctx).ListInterviewRounds(ctx, interviewRoundListParams(filter))
	if err != nil {
		return nil, s.logDBError(ctx, "list interview rounds", err, tenantIDField(filter.TenantID))
	}
	return mapInterviewRounds(rows), nil
}

func (s *Store) CountInterviewRounds(ctx context.Context, filter domain.InterviewRoundFilter) (int64, error) {
	count, err := s.getQueries(ctx).CountInterviewRounds(ctx, sqlc.CountInterviewRoundsParams{TenantID: filter.TenantID, ApplicationID: uuidFromPtr(filter.ApplicationID), Status: textFromPtr(filter.Status), InterviewerUserID: uuidFromPtr(filter.InterviewerUserID), DateFrom: timestamptzFromPtr(filter.DateFrom), DateTo: timestamptzFromPtr(filter.DateTo), Search: textFromPtr(filter.Search)})
	if err != nil {
		return 0, s.logDBError(ctx, "count interview rounds", err, tenantIDField(filter.TenantID))
	}
	return count, nil
}

func (s *Store) GetInterviewRound(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.InterviewRound, error) {
	row, err := s.getQueries(ctx).GetInterviewRound(ctx, sqlc.GetInterviewRoundParams{TenantID: tenantID, ID: id})
	if err != nil {
		return nil, s.logDBError(ctx, "get interview round", fmt.Errorf("hrms: get interview round: %w", err), tenantIDField(tenantID), stringField("interview_round_id", id.String()))
	}
	return mapInterviewRound(row), nil
}

func (s *Store) UpdateInterviewRound(ctx context.Context, item *domain.InterviewRound, actorID *uuid.UUID) (*domain.InterviewRound, error) {
	row, err := s.getQueries(ctx).UpdateInterviewRound(ctx, interviewRoundUpdateParams(item, actorID))
	if err != nil {
		return nil, s.logDBError(ctx, "update interview round", fmt.Errorf("hrms: update interview round: %w", err), tenantIDField(item.TenantID), stringField("interview_round_id", item.ID.String()))
	}
	return mapInterviewRound(row), nil
}

func (s *Store) UpdateInterviewRoundStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, remarks *string, feedback *string, score *float64, decision *string, completedAt *time.Time, actorID *uuid.UUID) (*domain.InterviewRound, error) {
	row, err := s.getQueries(ctx).UpdateInterviewRoundStatus(ctx, sqlc.UpdateInterviewRoundStatusParams{TenantID: tenantID, ID: id, Status: status, Remarks: textFromPtr(remarks), Feedback: textFromPtr(feedback), Score: numericFromFloatPtr(score), Decision: textFromPtr(decision), CompletedAt: timestamptzFromPtr(completedAt), UpdatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "update interview round status", fmt.Errorf("hrms: update interview round status: %w", err), tenantIDField(tenantID), stringField("interview_round_id", id.String()), stringField("status", status))
	}
	return mapInterviewRound(row), nil
}

func (s *Store) DeleteInterviewRound(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteInterviewRound(ctx, sqlc.SoftDeleteInterviewRoundParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete interview round", fmt.Errorf("hrms: delete interview round: %w", err), tenantIDField(tenantID), stringField("interview_round_id", id.String()))
	}
	return nil
}

func interviewRoundCreateParams(item *domain.InterviewRound, actorID *uuid.UUID) sqlc.CreateInterviewRoundParams {
	return sqlc.CreateInterviewRoundParams{TenantID: item.TenantID, ApplicationID: item.ApplicationID, RoundName: textFromPtr(item.RoundName), RoundNumber: int4FromPtr(item.RoundNumber), ScheduledDate: timestamptzFromPtr(item.ScheduledDate), DurationMinutes: int4FromPtr(item.DurationMinutes), InterviewerUserID: uuidFromPtr(item.InterviewerUserID), Mode: textFromPtr(item.Mode), MeetingLink: textFromPtr(item.MeetingLink), Location: textFromPtr(item.Location), Status: item.Status, Remarks: textFromPtr(item.Remarks), Timezone: item.Timezone, Feedback: textFromPtr(item.Feedback), Score: numericFromFloatPtr(item.Score), Decision: textFromPtr(item.Decision), CompletedAt: timestamptzFromPtr(item.CompletedAt), CreatedBy: uuidFromPtr(actorID)}
}

func interviewRoundUpdateParams(item *domain.InterviewRound, actorID *uuid.UUID) sqlc.UpdateInterviewRoundParams {
	return sqlc.UpdateInterviewRoundParams{TenantID: item.TenantID, ID: item.ID, ApplicationID: item.ApplicationID, RoundName: textFromPtr(item.RoundName), RoundNumber: int4FromPtr(item.RoundNumber), ScheduledDate: timestamptzFromPtr(item.ScheduledDate), DurationMinutes: int4FromPtr(item.DurationMinutes), InterviewerUserID: uuidFromPtr(item.InterviewerUserID), Mode: textFromPtr(item.Mode), MeetingLink: textFromPtr(item.MeetingLink), Location: textFromPtr(item.Location), Status: item.Status, Remarks: textFromPtr(item.Remarks), Timezone: item.Timezone, Feedback: textFromPtr(item.Feedback), Score: numericFromFloatPtr(item.Score), Decision: textFromPtr(item.Decision), CompletedAt: timestamptzFromPtr(item.CompletedAt), UpdatedBy: uuidFromPtr(actorID)}
}

func interviewRoundListParams(filter domain.InterviewRoundFilter) sqlc.ListInterviewRoundsParams {
	return sqlc.ListInterviewRoundsParams{TenantID: filter.TenantID, ApplicationID: uuidFromPtr(filter.ApplicationID), Status: textFromPtr(filter.Status), InterviewerUserID: uuidFromPtr(filter.InterviewerUserID), DateFrom: timestamptzFromPtr(filter.DateFrom), DateTo: timestamptzFromPtr(filter.DateTo), Search: textFromPtr(filter.Search), Offset: filter.Offset, Limit: filter.Limit}
}
