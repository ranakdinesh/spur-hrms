package postgres

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func mapPerformanceCheckIn(row sqlc.HrmsPerformanceCheckin) *domain.PerformanceCheckIn {
	return performanceCheckInFromParts(row.ID, row.TenantID, row.WorkerProfileID, row.ReviewerWorkerProfileID, row.CycleID, row.CheckinDate, row.PeriodStart, row.PeriodEnd, row.Mood, row.Status, row.Visibility, row.Highlights, row.Blockers, row.NextPlan, row.EmployeeComment, row.ManagerComment, row.Score, row.CalibrationBucket, row.ReviewedAt, row.ReviewedBy, row.Metadata, row.Inactive, row.CreatedAt, row.CreatedBy, row.UpdatedAt, row.UpdatedBy, nil, nil, nil, nil, 0, 0)
}

func mapPerformanceCheckInRow(row sqlc.GetPerformanceCheckInRow) *domain.PerformanceCheckIn {
	return performanceCheckInFromParts(row.ID, row.TenantID, row.WorkerProfileID, row.ReviewerWorkerProfileID, row.CycleID, row.CheckinDate, row.PeriodStart, row.PeriodEnd, row.Mood, row.Status, row.Visibility, row.Highlights, row.Blockers, row.NextPlan, row.EmployeeComment, row.ManagerComment, row.Score, row.CalibrationBucket, row.ReviewedAt, row.ReviewedBy, row.Metadata, row.Inactive, row.CreatedAt, row.CreatedBy, row.UpdatedAt, row.UpdatedBy, &row.WorkerDisplayName, ptrFromText(row.WorkerCode), ptrFromText(row.ReviewerDisplayName), ptrFromText(row.CycleName), row.FeedbackCount, floatFromNumeric(row.AverageFeedbackRating))
}

func mapPerformanceCheckInRows(rows []sqlc.ListPerformanceCheckInsRow) []*domain.PerformanceCheckIn {
	items := make([]*domain.PerformanceCheckIn, 0, len(rows))
	for _, row := range rows {
		items = append(items, performanceCheckInFromParts(row.ID, row.TenantID, row.WorkerProfileID, row.ReviewerWorkerProfileID, row.CycleID, row.CheckinDate, row.PeriodStart, row.PeriodEnd, row.Mood, row.Status, row.Visibility, row.Highlights, row.Blockers, row.NextPlan, row.EmployeeComment, row.ManagerComment, row.Score, row.CalibrationBucket, row.ReviewedAt, row.ReviewedBy, row.Metadata, row.Inactive, row.CreatedAt, row.CreatedBy, row.UpdatedAt, row.UpdatedBy, &row.WorkerDisplayName, ptrFromText(row.WorkerCode), ptrFromText(row.ReviewerDisplayName), ptrFromText(row.CycleName), row.FeedbackCount, floatFromNumeric(row.AverageFeedbackRating)))
	}
	return items
}

func performanceCheckInFromParts(id uuid.UUID, tenantID uuid.UUID, workerProfileID uuid.UUID, reviewerWorkerProfileID pgtype.UUID, cycleID pgtype.UUID, checkInDate pgtype.Date, periodStart pgtype.Date, periodEnd pgtype.Date, mood string, status string, visibility string, highlights pgtype.Text, blockers pgtype.Text, nextPlan pgtype.Text, employeeComment pgtype.Text, managerComment pgtype.Text, score pgtype.Numeric, calibrationBucket pgtype.Text, reviewedAt pgtype.Timestamptz, reviewedBy pgtype.UUID, metadata []byte, inactive bool, createdAt pgtype.Timestamptz, createdBy pgtype.UUID, updatedAt pgtype.Timestamptz, updatedBy pgtype.UUID, workerDisplayName *string, workerCode *string, reviewerDisplayName *string, cycleName *string, feedbackCount int32, averageFeedbackRating float64) *domain.PerformanceCheckIn {
	return &domain.PerformanceCheckIn{ID: id, TenantID: tenantID, WorkerProfileID: workerProfileID, ReviewerWorkerProfileID: ptrFromUUID(reviewerWorkerProfileID), CycleID: ptrFromUUID(cycleID), CheckInDate: timeFromDate(checkInDate), PeriodStart: timeFromDate(periodStart), PeriodEnd: timeFromDate(periodEnd), Mood: mood, Status: status, Visibility: visibility, Highlights: ptrFromText(highlights), Blockers: ptrFromText(blockers), NextPlan: ptrFromText(nextPlan), EmployeeComment: ptrFromText(employeeComment), ManagerComment: ptrFromText(managerComment), Score: floatPtrFromNumeric(score), CalibrationBucket: ptrFromText(calibrationBucket), ReviewedAt: ptrFromTimestamptz(reviewedAt), ReviewedBy: ptrFromUUID(reviewedBy), Metadata: okrRaw(metadata), Inactive: inactive, CreatedAt: timeFromTimestamptz(createdAt), CreatedBy: ptrFromUUID(createdBy), UpdatedAt: timeFromTimestamptz(updatedAt), UpdatedBy: ptrFromUUID(updatedBy), WorkerDisplayName: workerDisplayName, WorkerCode: workerCode, ReviewerDisplayName: reviewerDisplayName, CycleName: cycleName, FeedbackCount: feedbackCount, AverageFeedbackRating: averageFeedbackRating}
}

func mapPerformanceCheckInSummaryRows(rows []sqlc.GetPerformanceCheckInSummaryRow) []*domain.PerformanceCheckInSummaryRow {
	items := make([]*domain.PerformanceCheckInSummaryRow, 0, len(rows))
	for _, row := range rows {
		items = append(items, &domain.PerformanceCheckInSummaryRow{Status: row.Status, Mood: row.Mood, CheckInCount: row.CheckinCount, AverageScore: floatFromNumeric(row.AverageScore)})
	}
	return items
}

func mapFeedbackRequest(row sqlc.HrmsFeedbackRequest) *domain.FeedbackRequest {
	return feedbackRequestFromParts(row.ID, row.TenantID, row.SubjectWorkerProfileID, row.RequesterWorkerProfileID, row.ObjectiveID, row.Relationship, row.FeedbackType, row.Status, row.IsAnonymous, row.Visibility, row.DueDate, row.Prompt, row.Metadata, row.Inactive, row.CreatedAt, row.CreatedBy, row.UpdatedAt, row.UpdatedBy, nil, nil, nil, nil, 0)
}

func mapFeedbackRequestRow(row sqlc.GetFeedbackRequestRow) *domain.FeedbackRequest {
	return feedbackRequestFromParts(row.ID, row.TenantID, row.SubjectWorkerProfileID, row.RequesterWorkerProfileID, row.ObjectiveID, row.Relationship, row.FeedbackType, row.Status, row.IsAnonymous, row.Visibility, row.DueDate, row.Prompt, row.Metadata, row.Inactive, row.CreatedAt, row.CreatedBy, row.UpdatedAt, row.UpdatedBy, &row.SubjectDisplayName, ptrFromText(row.SubjectWorkerCode), ptrFromText(row.RequesterDisplayName), ptrFromText(row.ObjectiveTitle), row.ResponseCount)
}

func mapFeedbackRequestRows(rows []sqlc.ListFeedbackRequestsRow) []*domain.FeedbackRequest {
	items := make([]*domain.FeedbackRequest, 0, len(rows))
	for _, row := range rows {
		items = append(items, feedbackRequestFromParts(row.ID, row.TenantID, row.SubjectWorkerProfileID, row.RequesterWorkerProfileID, row.ObjectiveID, row.Relationship, row.FeedbackType, row.Status, row.IsAnonymous, row.Visibility, row.DueDate, row.Prompt, row.Metadata, row.Inactive, row.CreatedAt, row.CreatedBy, row.UpdatedAt, row.UpdatedBy, &row.SubjectDisplayName, ptrFromText(row.SubjectWorkerCode), ptrFromText(row.RequesterDisplayName), ptrFromText(row.ObjectiveTitle), row.ResponseCount))
	}
	return items
}

func feedbackRequestFromParts(id uuid.UUID, tenantID uuid.UUID, subjectWorkerProfileID uuid.UUID, requesterWorkerProfileID pgtype.UUID, objectiveID pgtype.UUID, relationship string, feedbackType string, status string, isAnonymous bool, visibility string, dueDate pgtype.Date, prompt pgtype.Text, metadata []byte, inactive bool, createdAt pgtype.Timestamptz, createdBy pgtype.UUID, updatedAt pgtype.Timestamptz, updatedBy pgtype.UUID, subjectDisplayName *string, subjectWorkerCode *string, requesterDisplayName *string, objectiveTitle *string, responseCount int32) *domain.FeedbackRequest {
	return &domain.FeedbackRequest{ID: id, TenantID: tenantID, SubjectWorkerProfileID: subjectWorkerProfileID, RequesterWorkerProfileID: ptrFromUUID(requesterWorkerProfileID), ObjectiveID: ptrFromUUID(objectiveID), Relationship: relationship, FeedbackType: feedbackType, Status: status, IsAnonymous: isAnonymous, Visibility: visibility, DueDate: ptrFromDate(dueDate), Prompt: ptrFromText(prompt), Metadata: okrRaw(metadata), Inactive: inactive, CreatedAt: timeFromTimestamptz(createdAt), CreatedBy: ptrFromUUID(createdBy), UpdatedAt: timeFromTimestamptz(updatedAt), UpdatedBy: ptrFromUUID(updatedBy), SubjectDisplayName: subjectDisplayName, SubjectWorkerCode: subjectWorkerCode, RequesterDisplayName: requesterDisplayName, ObjectiveTitle: objectiveTitle, ResponseCount: responseCount}
}

func mapFeedbackResponse(row sqlc.HrmsFeedbackResponse) *domain.FeedbackResponse {
	return &domain.FeedbackResponse{ID: row.ID, TenantID: row.TenantID, RequestID: row.RequestID, RespondentWorkerProfileID: ptrFromUUID(row.RespondentWorkerProfileID), Rating: floatPtrFromNumeric(row.Rating), Strengths: ptrFromText(row.Strengths), Improvements: ptrFromText(row.Improvements), Comments: ptrFromText(row.Comments), SubmittedAt: timeFromTimestamptz(row.SubmittedAt), Metadata: okrRaw(row.Metadata), Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy)}
}

func mapFeedbackResponseRows(rows []sqlc.ListFeedbackResponsesRow) []*domain.FeedbackResponse {
	items := make([]*domain.FeedbackResponse, 0, len(rows))
	for _, row := range rows {
		subjectID := row.SubjectWorkerProfileID
		var respondentName *string
		if row.RespondentDisplayName != "" {
			value := row.RespondentDisplayName
			respondentName = &value
		}
		items = append(items, &domain.FeedbackResponse{ID: row.ID, TenantID: row.TenantID, RequestID: row.RequestID, RespondentWorkerProfileID: ptrFromUUID(row.RespondentWorkerProfileID), Rating: floatPtrFromNumeric(row.Rating), Strengths: ptrFromText(row.Strengths), Improvements: ptrFromText(row.Improvements), Comments: ptrFromText(row.Comments), SubmittedAt: timeFromTimestamptz(row.SubmittedAt), Metadata: okrRaw(row.Metadata), Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), SubjectWorkerProfileID: &subjectID, IsAnonymous: row.IsAnonymous, SubjectDisplayName: &row.SubjectDisplayName, SubjectWorkerCode: ptrFromText(row.SubjectWorkerCode), RespondentDisplayName: respondentName, FeedbackType: &row.FeedbackType, Relationship: &row.Relationship})
	}
	return items
}

func mapPerformanceTimelineEvent(row sqlc.HrmsPerformanceTimelineEvent) *domain.PerformanceTimelineEvent {
	return performanceTimelineEventFromParts(row.ID, row.TenantID, row.WorkerProfileID, row.EventType, row.CheckinID, row.FeedbackRequestID, row.FeedbackResponseID, row.ObjectiveID, row.ActorWorkerProfileID, row.Title, row.Notes, row.Metadata, row.Inactive, row.CreatedAt, row.CreatedBy, nil, nil, nil)
}

func mapPerformanceTimelineEventRows(rows []sqlc.ListPerformanceTimelineEventsRow) []*domain.PerformanceTimelineEvent {
	items := make([]*domain.PerformanceTimelineEvent, 0, len(rows))
	for _, row := range rows {
		items = append(items, performanceTimelineEventFromParts(row.ID, row.TenantID, row.WorkerProfileID, row.EventType, row.CheckinID, row.FeedbackRequestID, row.FeedbackResponseID, row.ObjectiveID, row.ActorWorkerProfileID, row.Title, row.Notes, row.Metadata, row.Inactive, row.CreatedAt, row.CreatedBy, &row.WorkerDisplayName, ptrFromText(row.ActorDisplayName), ptrFromText(row.ObjectiveTitle)))
	}
	return items
}

func performanceTimelineEventFromParts(id uuid.UUID, tenantID uuid.UUID, workerProfileID uuid.UUID, eventType string, checkInID pgtype.UUID, feedbackRequestID pgtype.UUID, feedbackResponseID pgtype.UUID, objectiveID pgtype.UUID, actorWorkerProfileID pgtype.UUID, title string, notes pgtype.Text, metadata []byte, inactive bool, createdAt pgtype.Timestamptz, createdBy pgtype.UUID, workerDisplayName *string, actorDisplayName *string, objectiveTitle *string) *domain.PerformanceTimelineEvent {
	return &domain.PerformanceTimelineEvent{ID: id, TenantID: tenantID, WorkerProfileID: workerProfileID, EventType: eventType, CheckInID: ptrFromUUID(checkInID), FeedbackRequestID: ptrFromUUID(feedbackRequestID), FeedbackResponseID: ptrFromUUID(feedbackResponseID), ObjectiveID: ptrFromUUID(objectiveID), ActorWorkerProfileID: ptrFromUUID(actorWorkerProfileID), Title: title, Notes: ptrFromText(notes), Metadata: okrRaw(metadata), Inactive: inactive, CreatedAt: timeFromTimestamptz(createdAt), CreatedBy: ptrFromUUID(createdBy), WorkerDisplayName: workerDisplayName, ActorDisplayName: actorDisplayName, ObjectiveTitle: objectiveTitle}
}

func mapPerformanceCalibrationRows(rows []sqlc.ListPerformanceCalibrationRowsRow) []*domain.PerformanceCalibrationRow {
	items := make([]*domain.PerformanceCalibrationRow, 0, len(rows))
	for _, row := range rows {
		bucket := row.CalibrationBucket
		item := &domain.PerformanceCalibrationRow{WorkerProfileID: row.WorkerProfileID, WorkerDisplayName: row.WorkerDisplayName, WorkerCode: ptrFromText(row.WorkerCode), CycleID: ptrFromUUID(row.CycleID), CycleName: ptrFromText(row.CycleName), CheckInCount: row.CheckinCount, SubmittedCheckInCount: row.SubmittedCheckinCount, AverageScore: floatFromNumeric(row.AverageScore), AverageOKRProgress: floatFromNumeric(row.AverageOkrProgress), FeedbackCount: row.FeedbackCount, AverageFeedbackRating: floatFromNumeric(row.AverageFeedbackRating)}
		if bucket != "" {
			item.CalibrationBucket = &bucket
		}
		items = append(items, item)
	}
	return items
}
