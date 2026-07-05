package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func (s *Store) CreatePerformanceCheckIn(ctx context.Context, item *domain.PerformanceCheckIn, actorID *uuid.UUID) (*domain.PerformanceCheckIn, error) {
	row, err := s.getQueries(ctx).CreatePerformanceCheckIn(ctx, performanceCreateCheckInParams(item, actorID))
	if err != nil {
		return nil, s.logDBError(ctx, "create performance check-in", err, tenantIDField(item.TenantID), stringField("worker_profile_id", item.WorkerProfileID.String()))
	}
	return mapPerformanceCheckIn(row), nil
}

func (s *Store) UpdatePerformanceCheckIn(ctx context.Context, item *domain.PerformanceCheckIn, actorID *uuid.UUID) (*domain.PerformanceCheckIn, error) {
	row, err := s.getQueries(ctx).UpdatePerformanceCheckIn(ctx, performanceUpdateCheckInParams(item, actorID))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrPerformanceCheckInNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "update performance check-in", err, tenantIDField(item.TenantID), stringField("checkin_id", item.ID.String()))
	}
	return mapPerformanceCheckIn(row), nil
}

func (s *Store) ReviewPerformanceCheckIn(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, managerComment *string, score *float64, calibrationBucket *string, actorID *uuid.UUID) (*domain.PerformanceCheckIn, error) {
	row, err := s.getQueries(ctx).ReviewPerformanceCheckIn(ctx, sqlc.ReviewPerformanceCheckInParams{TenantID: tenantID, ID: id, Status: status, ManagerComment: textFromPtr(managerComment), Score: numericFromFloatPtr(score), CalibrationBucket: textFromPtr(calibrationBucket), ReviewedBy: uuidFromPtr(actorID)})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrPerformanceCheckInNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "review performance check-in", err, tenantIDField(tenantID), stringField("checkin_id", id.String()))
	}
	return mapPerformanceCheckIn(row), nil
}

func (s *Store) UpdatePerformanceCheckInStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, actorID *uuid.UUID) (*domain.PerformanceCheckIn, error) {
	row, err := s.getQueries(ctx).UpdatePerformanceCheckInStatus(ctx, sqlc.UpdatePerformanceCheckInStatusParams{TenantID: tenantID, ID: id, Status: status, UpdatedBy: uuidFromPtr(actorID)})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrPerformanceCheckInNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "update performance check-in status", err, tenantIDField(tenantID), stringField("checkin_id", id.String()))
	}
	return mapPerformanceCheckIn(row), nil
}

func (s *Store) GetPerformanceCheckIn(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.PerformanceCheckIn, error) {
	row, err := s.getQueries(ctx).GetPerformanceCheckIn(ctx, sqlc.GetPerformanceCheckInParams{TenantID: tenantID, ID: id})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrPerformanceCheckInNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "get performance check-in", err, tenantIDField(tenantID), stringField("checkin_id", id.String()))
	}
	return mapPerformanceCheckInRow(row), nil
}

func (s *Store) ListPerformanceCheckIns(ctx context.Context, filter domain.PerformanceCheckInFilter) ([]*domain.PerformanceCheckIn, error) {
	rows, err := s.getQueries(ctx).ListPerformanceCheckIns(ctx, sqlc.ListPerformanceCheckInsParams{TenantID: filter.TenantID, Column2: uuidValueFromPtr(filter.WorkerProfileID), Column3: uuidValueFromPtr(filter.ReviewerWorkerProfileID), Column4: uuidValueFromPtr(filter.CycleID), Column5: stringFromPtr(filter.Status), Column6: stringFromPtr(filter.Mood)})
	if err != nil {
		return nil, s.logDBError(ctx, "list performance check-ins", err, tenantIDField(filter.TenantID))
	}
	return mapPerformanceCheckInRows(rows), nil
}

func (s *Store) DeletePerformanceCheckIn(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeletePerformanceCheckIn(ctx, sqlc.SoftDeletePerformanceCheckInParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete performance check-in", err, tenantIDField(tenantID), stringField("checkin_id", id.String()))
	}
	return nil
}

func (s *Store) GetPerformanceCheckInSummary(ctx context.Context, tenantID uuid.UUID, cycleID *uuid.UUID) ([]*domain.PerformanceCheckInSummaryRow, error) {
	rows, err := s.getQueries(ctx).GetPerformanceCheckInSummary(ctx, sqlc.GetPerformanceCheckInSummaryParams{TenantID: tenantID, Column2: uuidValueFromPtr(cycleID)})
	if err != nil {
		return nil, s.logDBError(ctx, "get performance check-in summary", err, tenantIDField(tenantID))
	}
	return mapPerformanceCheckInSummaryRows(rows), nil
}

func (s *Store) CreateFeedbackRequest(ctx context.Context, item *domain.FeedbackRequest, actorID *uuid.UUID) (*domain.FeedbackRequest, error) {
	row, err := s.getQueries(ctx).CreateFeedbackRequest(ctx, sqlc.CreateFeedbackRequestParams{TenantID: item.TenantID, SubjectWorkerProfileID: item.SubjectWorkerProfileID, RequesterWorkerProfileID: uuidFromPtr(item.RequesterWorkerProfileID), ObjectiveID: uuidFromPtr(item.ObjectiveID), Relationship: item.Relationship, FeedbackType: item.FeedbackType, Status: item.Status, IsAnonymous: item.IsAnonymous, Visibility: item.Visibility, DueDate: dateFromPtr(item.DueDate), Prompt: textFromPtr(item.Prompt), Column12: jsonBytesFromRaw(item.Metadata), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create feedback request", err, tenantIDField(item.TenantID), stringField("subject_worker_profile_id", item.SubjectWorkerProfileID.String()))
	}
	return mapFeedbackRequest(row), nil
}

func (s *Store) UpdateFeedbackRequest(ctx context.Context, item *domain.FeedbackRequest, actorID *uuid.UUID) (*domain.FeedbackRequest, error) {
	row, err := s.getQueries(ctx).UpdateFeedbackRequest(ctx, sqlc.UpdateFeedbackRequestParams{TenantID: item.TenantID, ID: item.ID, SubjectWorkerProfileID: item.SubjectWorkerProfileID, RequesterWorkerProfileID: uuidFromPtr(item.RequesterWorkerProfileID), ObjectiveID: uuidFromPtr(item.ObjectiveID), Relationship: item.Relationship, FeedbackType: item.FeedbackType, Status: item.Status, IsAnonymous: item.IsAnonymous, Visibility: item.Visibility, DueDate: dateFromPtr(item.DueDate), Prompt: textFromPtr(item.Prompt), Metadata: jsonBytesFromRaw(item.Metadata), UpdatedBy: uuidFromPtr(actorID)})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrFeedbackRequestNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "update feedback request", err, tenantIDField(item.TenantID), stringField("feedback_request_id", item.ID.String()))
	}
	return mapFeedbackRequest(row), nil
}

func (s *Store) UpdateFeedbackRequestStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, actorID *uuid.UUID) (*domain.FeedbackRequest, error) {
	row, err := s.getQueries(ctx).UpdateFeedbackRequestStatus(ctx, sqlc.UpdateFeedbackRequestStatusParams{TenantID: tenantID, ID: id, Status: status, UpdatedBy: uuidFromPtr(actorID)})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrFeedbackRequestNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "update feedback request status", err, tenantIDField(tenantID), stringField("feedback_request_id", id.String()))
	}
	return mapFeedbackRequest(row), nil
}

func (s *Store) GetFeedbackRequest(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.FeedbackRequest, error) {
	row, err := s.getQueries(ctx).GetFeedbackRequest(ctx, sqlc.GetFeedbackRequestParams{TenantID: tenantID, ID: id})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrFeedbackRequestNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "get feedback request", err, tenantIDField(tenantID), stringField("feedback_request_id", id.String()))
	}
	return mapFeedbackRequestRow(row), nil
}

func (s *Store) ListFeedbackRequests(ctx context.Context, filter domain.FeedbackRequestFilter) ([]*domain.FeedbackRequest, error) {
	rows, err := s.getQueries(ctx).ListFeedbackRequests(ctx, sqlc.ListFeedbackRequestsParams{TenantID: filter.TenantID, Column2: uuidValueFromPtr(filter.SubjectWorkerProfileID), Column3: uuidValueFromPtr(filter.RequesterWorkerProfileID), Column4: stringFromPtr(filter.Status), Column5: stringFromPtr(filter.FeedbackType)})
	if err != nil {
		return nil, s.logDBError(ctx, "list feedback requests", err, tenantIDField(filter.TenantID))
	}
	return mapFeedbackRequestRows(rows), nil
}

func (s *Store) CreateFeedbackResponse(ctx context.Context, item *domain.FeedbackResponse, actorID *uuid.UUID) (*domain.FeedbackResponse, error) {
	row, err := s.getQueries(ctx).CreateFeedbackResponse(ctx, sqlc.CreateFeedbackResponseParams{TenantID: item.TenantID, RequestID: item.RequestID, RespondentWorkerProfileID: uuidFromPtr(item.RespondentWorkerProfileID), Rating: numericFromFloatPtr(item.Rating), Strengths: textFromPtr(item.Strengths), Improvements: textFromPtr(item.Improvements), Comments: textFromPtr(item.Comments), Column8: jsonBytesFromRaw(item.Metadata), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create feedback response", err, tenantIDField(item.TenantID), stringField("feedback_request_id", item.RequestID.String()))
	}
	return mapFeedbackResponse(row), nil
}

func (s *Store) ListFeedbackResponses(ctx context.Context, filter domain.FeedbackResponseFilter) ([]*domain.FeedbackResponse, error) {
	rows, err := s.getQueries(ctx).ListFeedbackResponses(ctx, sqlc.ListFeedbackResponsesParams{TenantID: filter.TenantID, Column2: uuidValueFromPtr(filter.RequestID), Column3: uuidValueFromPtr(filter.SubjectWorkerProfileID), Column4: uuidValueFromPtr(filter.RespondentWorkerProfileID)})
	if err != nil {
		return nil, s.logDBError(ctx, "list feedback responses", err, tenantIDField(filter.TenantID))
	}
	return mapFeedbackResponseRows(rows), nil
}

func (s *Store) CreatePerformanceTimelineEvent(ctx context.Context, item *domain.PerformanceTimelineEvent, actorID *uuid.UUID) (*domain.PerformanceTimelineEvent, error) {
	row, err := s.getQueries(ctx).CreatePerformanceTimelineEvent(ctx, sqlc.CreatePerformanceTimelineEventParams{TenantID: item.TenantID, WorkerProfileID: item.WorkerProfileID, EventType: item.EventType, CheckinID: uuidFromPtr(item.CheckInID), FeedbackRequestID: uuidFromPtr(item.FeedbackRequestID), FeedbackResponseID: uuidFromPtr(item.FeedbackResponseID), ObjectiveID: uuidFromPtr(item.ObjectiveID), ActorWorkerProfileID: uuidFromPtr(item.ActorWorkerProfileID), Title: item.Title, Notes: textFromPtr(item.Notes), Column11: jsonBytesFromRaw(item.Metadata), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create performance timeline event", err, tenantIDField(item.TenantID), stringField("worker_profile_id", item.WorkerProfileID.String()))
	}
	return mapPerformanceTimelineEvent(row), nil
}

func (s *Store) ListPerformanceTimelineEvents(ctx context.Context, filter domain.PerformanceTimelineFilter) ([]*domain.PerformanceTimelineEvent, error) {
	rows, err := s.getQueries(ctx).ListPerformanceTimelineEvents(ctx, sqlc.ListPerformanceTimelineEventsParams{TenantID: filter.TenantID, Column2: uuidValueFromPtr(filter.WorkerProfileID), Column3: stringFromPtr(filter.EventType)})
	if err != nil {
		return nil, s.logDBError(ctx, "list performance timeline events", err, tenantIDField(filter.TenantID))
	}
	return mapPerformanceTimelineEventRows(rows), nil
}

func (s *Store) ListPerformanceCalibrationRows(ctx context.Context, filter domain.PerformanceCalibrationFilter) ([]*domain.PerformanceCalibrationRow, error) {
	rows, err := s.getQueries(ctx).ListPerformanceCalibrationRows(ctx, sqlc.ListPerformanceCalibrationRowsParams{TenantID: filter.TenantID, Column2: uuidValueFromPtr(filter.CycleID), Column3: uuidValueFromPtr(filter.WorkerProfileID)})
	if err != nil {
		return nil, s.logDBError(ctx, "list performance calibration rows", err, tenantIDField(filter.TenantID))
	}
	return mapPerformanceCalibrationRows(rows), nil
}

func performanceCreateCheckInParams(item *domain.PerformanceCheckIn, actorID *uuid.UUID) sqlc.CreatePerformanceCheckInParams {
	return sqlc.CreatePerformanceCheckInParams{TenantID: item.TenantID, WorkerProfileID: item.WorkerProfileID, ReviewerWorkerProfileID: uuidFromPtr(item.ReviewerWorkerProfileID), CycleID: uuidFromPtr(item.CycleID), CheckinDate: dateFromPtr(&item.CheckInDate), PeriodStart: dateFromPtr(&item.PeriodStart), PeriodEnd: dateFromPtr(&item.PeriodEnd), Mood: item.Mood, Status: item.Status, Visibility: item.Visibility, Highlights: textFromPtr(item.Highlights), Blockers: textFromPtr(item.Blockers), NextPlan: textFromPtr(item.NextPlan), EmployeeComment: textFromPtr(item.EmployeeComment), ManagerComment: textFromPtr(item.ManagerComment), Score: numericFromFloatPtr(item.Score), CalibrationBucket: textFromPtr(item.CalibrationBucket), Column18: jsonBytesFromRaw(item.Metadata), CreatedBy: uuidFromPtr(actorID)}
}

func performanceUpdateCheckInParams(item *domain.PerformanceCheckIn, actorID *uuid.UUID) sqlc.UpdatePerformanceCheckInParams {
	return sqlc.UpdatePerformanceCheckInParams{TenantID: item.TenantID, ID: item.ID, WorkerProfileID: item.WorkerProfileID, ReviewerWorkerProfileID: uuidFromPtr(item.ReviewerWorkerProfileID), CycleID: uuidFromPtr(item.CycleID), CheckinDate: dateFromPtr(&item.CheckInDate), PeriodStart: dateFromPtr(&item.PeriodStart), PeriodEnd: dateFromPtr(&item.PeriodEnd), Mood: item.Mood, Status: item.Status, Visibility: item.Visibility, Highlights: textFromPtr(item.Highlights), Blockers: textFromPtr(item.Blockers), NextPlan: textFromPtr(item.NextPlan), EmployeeComment: textFromPtr(item.EmployeeComment), ManagerComment: textFromPtr(item.ManagerComment), Score: numericFromFloatPtr(item.Score), CalibrationBucket: textFromPtr(item.CalibrationBucket), Metadata: jsonBytesFromRaw(item.Metadata), UpdatedBy: uuidFromPtr(actorID)}
}
