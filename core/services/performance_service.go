package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) CreatePerformanceCheckIn(ctx context.Context, cmd ports.PerformanceCheckInCommand) (*domain.PerformanceCheckIn, error) {
	item, err := s.preparePerformanceCheckIn(ctx, cmd)
	if err != nil {
		s.logError("validate performance check-in", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	result, err := s.performance.CreatePerformanceCheckIn(ctx, item, cmd.ActorID)
	if err != nil {
		return nil, err
	}
	s.createPerformanceTimeline(ctx, domain.PerformanceTimelineEventInput{TenantID: result.TenantID, WorkerProfileID: result.WorkerProfileID, EventType: domain.PerformanceTimelineCheckInCreated, CheckInID: &result.ID, ActorWorkerProfileID: &result.WorkerProfileID, Title: "Check-in created", Notes: result.Highlights}, cmd.ActorID)
	return result, nil
}

func (s *TenantService) UpdatePerformanceCheckIn(ctx context.Context, cmd ports.PerformanceCheckInCommand) (*domain.PerformanceCheckIn, error) {
	if cmd.ID == uuid.Nil {
		return nil, domain.ErrInvalidPerformanceCheckIn
	}
	item, err := s.preparePerformanceCheckIn(ctx, cmd)
	if err != nil {
		s.logError("validate performance check-in update", err, serviceTenantIDField(cmd.TenantID), serviceStringField("checkin_id", cmd.ID.String()))
		return nil, err
	}
	item.ID = cmd.ID
	return s.performance.UpdatePerformanceCheckIn(ctx, item, cmd.ActorID)
}

func (s *TenantService) SubmitPerformanceCheckIn(ctx context.Context, cmd ports.PerformanceStatusCommand) (*domain.PerformanceCheckIn, error) {
	if cmd.TenantID == uuid.Nil || cmd.ID == uuid.Nil {
		return nil, domain.ErrInvalidPerformanceCheckIn
	}
	result, err := s.performance.UpdatePerformanceCheckInStatus(ctx, cmd.TenantID, cmd.ID, domain.PerformanceCheckInStatusSubmitted, cmd.ActorID)
	if err != nil {
		return nil, err
	}
	s.createPerformanceTimeline(ctx, domain.PerformanceTimelineEventInput{TenantID: result.TenantID, WorkerProfileID: result.WorkerProfileID, EventType: domain.PerformanceTimelineCheckInSubmitted, CheckInID: &result.ID, ActorWorkerProfileID: &result.WorkerProfileID, Title: "Check-in submitted", Notes: result.EmployeeComment}, cmd.ActorID)
	return result, nil
}

func (s *TenantService) ReviewPerformanceCheckIn(ctx context.Context, cmd ports.PerformanceCheckInReviewCommand) (*domain.PerformanceCheckIn, error) {
	if cmd.TenantID == uuid.Nil || cmd.ID == uuid.Nil || !stringIn(cmd.Status, []string{domain.PerformanceCheckInStatusReviewed, domain.PerformanceCheckInStatusClosed}) {
		return nil, domain.ErrInvalidPerformanceCheckIn
	}
	result, err := s.performance.ReviewPerformanceCheckIn(ctx, cmd.TenantID, cmd.ID, cmd.Status, cmd.ManagerComment, cmd.Score, cmd.CalibrationBucket, cmd.ActorID)
	if err != nil {
		return nil, err
	}
	s.createPerformanceTimeline(ctx, domain.PerformanceTimelineEventInput{TenantID: result.TenantID, WorkerProfileID: result.WorkerProfileID, EventType: domain.PerformanceTimelineCheckInReviewed, CheckInID: &result.ID, ActorWorkerProfileID: result.ReviewerWorkerProfileID, Title: "Check-in reviewed", Notes: result.ManagerComment}, cmd.ActorID)
	return result, nil
}

func (s *TenantService) GetPerformanceCheckIn(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.PerformanceCheckIn, error) {
	if tenantID == uuid.Nil || id == uuid.Nil {
		return nil, domain.ErrInvalidPerformanceCheckIn
	}
	return s.performance.GetPerformanceCheckIn(ctx, tenantID, id)
}

func (s *TenantService) ListPerformanceCheckIns(ctx context.Context, filter domain.PerformanceCheckInFilter) ([]*domain.PerformanceCheckIn, error) {
	if filter.TenantID == uuid.Nil {
		return nil, domain.ErrInvalidTenantID
	}
	filter.Status = domain.NormalizePerformanceSearch(filter.Status)
	filter.Mood = domain.NormalizePerformanceSearch(filter.Mood)
	return s.performance.ListPerformanceCheckIns(ctx, filter)
}

func (s *TenantService) DeletePerformanceCheckIn(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if tenantID == uuid.Nil || id == uuid.Nil {
		return domain.ErrInvalidPerformanceCheckIn
	}
	return s.performance.DeletePerformanceCheckIn(ctx, tenantID, id, actorID)
}

func (s *TenantService) GetPerformanceCheckInSummary(ctx context.Context, tenantID uuid.UUID, cycleID *uuid.UUID) ([]*domain.PerformanceCheckInSummaryRow, error) {
	if tenantID == uuid.Nil {
		return nil, domain.ErrInvalidTenantID
	}
	return s.performance.GetPerformanceCheckInSummary(ctx, tenantID, cycleID)
}

func (s *TenantService) CreateFeedbackRequest(ctx context.Context, cmd ports.FeedbackRequestCommand) (*domain.FeedbackRequest, error) {
	item, err := s.prepareFeedbackRequest(ctx, cmd)
	if err != nil {
		s.logError("validate feedback request", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	result, err := s.performance.CreateFeedbackRequest(ctx, item, cmd.ActorID)
	if err != nil {
		return nil, err
	}
	s.createPerformanceTimeline(ctx, domain.PerformanceTimelineEventInput{TenantID: result.TenantID, WorkerProfileID: result.SubjectWorkerProfileID, EventType: domain.PerformanceTimelineFeedbackRequested, FeedbackRequestID: &result.ID, ObjectiveID: result.ObjectiveID, ActorWorkerProfileID: result.RequesterWorkerProfileID, Title: "Feedback requested", Notes: result.Prompt}, cmd.ActorID)
	return result, nil
}

func (s *TenantService) UpdateFeedbackRequest(ctx context.Context, cmd ports.FeedbackRequestCommand) (*domain.FeedbackRequest, error) {
	if cmd.ID == uuid.Nil {
		return nil, domain.ErrInvalidFeedbackRequest
	}
	item, err := s.prepareFeedbackRequest(ctx, cmd)
	if err != nil {
		s.logError("validate feedback request update", err, serviceTenantIDField(cmd.TenantID), serviceStringField("feedback_request_id", cmd.ID.String()))
		return nil, err
	}
	item.ID = cmd.ID
	return s.performance.UpdateFeedbackRequest(ctx, item, cmd.ActorID)
}

func (s *TenantService) UpdateFeedbackRequestStatus(ctx context.Context, cmd ports.FeedbackStatusCommand) (*domain.FeedbackRequest, error) {
	if cmd.TenantID == uuid.Nil || cmd.ID == uuid.Nil || !stringIn(cmd.Status, []string{domain.FeedbackStatusRequested, domain.FeedbackStatusSubmitted, domain.FeedbackStatusDeclined, domain.FeedbackStatusExpired, domain.FeedbackStatusCancelled}) {
		return nil, domain.ErrInvalidFeedbackRequest
	}
	return s.performance.UpdateFeedbackRequestStatus(ctx, cmd.TenantID, cmd.ID, cmd.Status, cmd.ActorID)
}

func (s *TenantService) GetFeedbackRequest(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.FeedbackRequest, error) {
	if tenantID == uuid.Nil || id == uuid.Nil {
		return nil, domain.ErrInvalidFeedbackRequest
	}
	return s.performance.GetFeedbackRequest(ctx, tenantID, id)
}

func (s *TenantService) ListFeedbackRequests(ctx context.Context, filter domain.FeedbackRequestFilter) ([]*domain.FeedbackRequest, error) {
	if filter.TenantID == uuid.Nil {
		return nil, domain.ErrInvalidTenantID
	}
	filter.Status = domain.NormalizePerformanceSearch(filter.Status)
	filter.FeedbackType = domain.NormalizePerformanceSearch(filter.FeedbackType)
	return s.performance.ListFeedbackRequests(ctx, filter)
}

func (s *TenantService) CreateFeedbackResponse(ctx context.Context, cmd ports.FeedbackResponseCommand) (*domain.FeedbackResponse, error) {
	request, err := s.GetFeedbackRequest(ctx, cmd.TenantID, cmd.RequestID)
	if err != nil {
		return nil, err
	}
	if cmd.RespondentWorkerProfileID != nil {
		if _, err := s.GetWorkerProfile(ctx, cmd.TenantID, *cmd.RespondentWorkerProfileID); err != nil {
			return nil, err
		}
	}
	item, err := domain.NewFeedbackResponse(domain.FeedbackResponseInput{TenantID: cmd.TenantID, RequestID: cmd.RequestID, RespondentWorkerProfileID: cmd.RespondentWorkerProfileID, Rating: cmd.Rating, Strengths: cmd.Strengths, Improvements: cmd.Improvements, Comments: cmd.Comments, Metadata: cmd.Metadata})
	if err != nil {
		s.logError("validate feedback response", err, serviceTenantIDField(cmd.TenantID), serviceStringField("feedback_request_id", cmd.RequestID.String()))
		return nil, err
	}
	result, err := s.performance.CreateFeedbackResponse(ctx, item, cmd.ActorID)
	if err != nil {
		return nil, err
	}
	_, _ = s.performance.UpdateFeedbackRequestStatus(ctx, cmd.TenantID, cmd.RequestID, domain.FeedbackStatusSubmitted, cmd.ActorID)
	s.createPerformanceTimeline(ctx, domain.PerformanceTimelineEventInput{TenantID: cmd.TenantID, WorkerProfileID: request.SubjectWorkerProfileID, EventType: domain.PerformanceTimelineFeedbackSubmitted, FeedbackRequestID: &cmd.RequestID, FeedbackResponseID: &result.ID, ObjectiveID: request.ObjectiveID, ActorWorkerProfileID: cmd.RespondentWorkerProfileID, Title: "Feedback submitted", Notes: result.Comments}, cmd.ActorID)
	return result, nil
}

func (s *TenantService) ListFeedbackResponses(ctx context.Context, filter domain.FeedbackResponseFilter) ([]*domain.FeedbackResponse, error) {
	if filter.TenantID == uuid.Nil {
		return nil, domain.ErrInvalidTenantID
	}
	return s.performance.ListFeedbackResponses(ctx, filter)
}

func (s *TenantService) ListPerformanceTimelineEvents(ctx context.Context, filter domain.PerformanceTimelineFilter) ([]*domain.PerformanceTimelineEvent, error) {
	if filter.TenantID == uuid.Nil {
		return nil, domain.ErrInvalidTenantID
	}
	filter.EventType = domain.NormalizePerformanceSearch(filter.EventType)
	return s.performance.ListPerformanceTimelineEvents(ctx, filter)
}

func (s *TenantService) ListPerformanceCalibrationRows(ctx context.Context, filter domain.PerformanceCalibrationFilter) ([]*domain.PerformanceCalibrationRow, error) {
	if filter.TenantID == uuid.Nil {
		return nil, domain.ErrInvalidTenantID
	}
	return s.performance.ListPerformanceCalibrationRows(ctx, filter)
}

func (s *TenantService) preparePerformanceCheckIn(ctx context.Context, cmd ports.PerformanceCheckInCommand) (*domain.PerformanceCheckIn, error) {
	if _, err := s.GetWorkerProfile(ctx, cmd.TenantID, cmd.WorkerProfileID); err != nil {
		return nil, err
	}
	if cmd.ReviewerWorkerProfileID != nil {
		if _, err := s.GetWorkerProfile(ctx, cmd.TenantID, *cmd.ReviewerWorkerProfileID); err != nil {
			return nil, err
		}
	}
	if cmd.CycleID != nil {
		if _, err := s.GetOKRCycle(ctx, cmd.TenantID, *cmd.CycleID); err != nil {
			return nil, err
		}
	}
	checkInDate, err := parseWorkerProfileDate(cmd.CheckInDate)
	if err != nil {
		return nil, err
	}
	periodStart, err := parseWorkerProfileDate(cmd.PeriodStart)
	if err != nil {
		return nil, err
	}
	periodEnd, err := parseWorkerProfileDate(cmd.PeriodEnd)
	if err != nil {
		return nil, err
	}
	return domain.NewPerformanceCheckIn(domain.PerformanceCheckInInput{TenantID: cmd.TenantID, WorkerProfileID: cmd.WorkerProfileID, ReviewerWorkerProfileID: cmd.ReviewerWorkerProfileID, CycleID: cmd.CycleID, CheckInDate: checkInDate, PeriodStart: periodStart, PeriodEnd: periodEnd, Mood: cmd.Mood, Status: cmd.Status, Visibility: cmd.Visibility, Highlights: cmd.Highlights, Blockers: cmd.Blockers, NextPlan: cmd.NextPlan, EmployeeComment: cmd.EmployeeComment, ManagerComment: cmd.ManagerComment, Score: cmd.Score, CalibrationBucket: cmd.CalibrationBucket, Metadata: cmd.Metadata})
}

func (s *TenantService) prepareFeedbackRequest(ctx context.Context, cmd ports.FeedbackRequestCommand) (*domain.FeedbackRequest, error) {
	if _, err := s.GetWorkerProfile(ctx, cmd.TenantID, cmd.SubjectWorkerProfileID); err != nil {
		return nil, err
	}
	if cmd.RequesterWorkerProfileID != nil {
		if _, err := s.GetWorkerProfile(ctx, cmd.TenantID, *cmd.RequesterWorkerProfileID); err != nil {
			return nil, err
		}
	}
	if cmd.ObjectiveID != nil {
		if _, err := s.GetObjective(ctx, cmd.TenantID, *cmd.ObjectiveID); err != nil {
			return nil, err
		}
	}
	dueDate, err := parseWorkerProfileDate(cmd.DueDate)
	if err != nil {
		return nil, err
	}
	return domain.NewFeedbackRequest(domain.FeedbackRequestInput{TenantID: cmd.TenantID, SubjectWorkerProfileID: cmd.SubjectWorkerProfileID, RequesterWorkerProfileID: cmd.RequesterWorkerProfileID, ObjectiveID: cmd.ObjectiveID, Relationship: cmd.Relationship, FeedbackType: cmd.FeedbackType, Status: cmd.Status, IsAnonymous: cmd.IsAnonymous, Visibility: cmd.Visibility, DueDate: dueDate, Prompt: cmd.Prompt, Metadata: cmd.Metadata})
}

func (s *TenantService) createPerformanceTimeline(ctx context.Context, input domain.PerformanceTimelineEventInput, actorID *uuid.UUID) {
	event, err := domain.NewPerformanceTimelineEvent(input)
	if err != nil {
		s.logError("validate performance timeline", err, serviceTenantIDField(input.TenantID))
		return
	}
	if _, err := s.performance.CreatePerformanceTimelineEvent(ctx, event, actorID); err != nil {
		s.logError("create performance timeline", err, serviceTenantIDField(input.TenantID))
	}
}
