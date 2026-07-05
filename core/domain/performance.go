package domain

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	PerformanceCheckInMoodGreat    = "great"
	PerformanceCheckInMoodGood     = "good"
	PerformanceCheckInMoodNeutral  = "neutral"
	PerformanceCheckInMoodLow      = "low"
	PerformanceCheckInMoodStressed = "stressed"

	PerformanceCheckInStatusDraft     = "draft"
	PerformanceCheckInStatusSubmitted = "submitted"
	PerformanceCheckInStatusReviewed  = "reviewed"
	PerformanceCheckInStatusClosed    = "closed"

	PerformanceVisibilityWorkerManagerHR = "worker_manager_hr"
	PerformanceVisibilityManagerHR       = "manager_hr"
	PerformanceVisibilityHROnly          = "hr_only"

	PerformanceCalibrationHigh    = "high"
	PerformanceCalibrationSolid   = "solid"
	PerformanceCalibrationWatch   = "watch"
	PerformanceCalibrationImprove = "improve"

	FeedbackRelationshipManager      = "manager"
	FeedbackRelationshipPeer         = "peer"
	FeedbackRelationshipDirectReport = "direct_report"
	FeedbackRelationshipSelf         = "self"
	FeedbackRelationshipHR           = "hr"
	FeedbackRelationshipClient       = "client"

	FeedbackType360           = "360"
	FeedbackTypeProject       = "project"
	FeedbackTypeGeneral       = "general"
	FeedbackTypeOKR           = "okr"
	FeedbackTypeManagerReview = "manager_review"

	FeedbackStatusRequested = "requested"
	FeedbackStatusSubmitted = "submitted"
	FeedbackStatusDeclined  = "declined"
	FeedbackStatusExpired   = "expired"
	FeedbackStatusCancelled = "cancelled"

	FeedbackVisibilitySubjectManagerHR = "subject_manager_hr"
	FeedbackVisibilityManagerHR        = "manager_hr"
	FeedbackVisibilityHROnly           = "hr_only"
	FeedbackVisibilitySubjectOnly      = "subject_only"

	PerformanceTimelineCheckInCreated    = "checkin_created"
	PerformanceTimelineCheckInSubmitted  = "checkin_submitted"
	PerformanceTimelineCheckInReviewed   = "checkin_reviewed"
	PerformanceTimelineFeedbackRequested = "feedback_requested"
	PerformanceTimelineFeedbackSubmitted = "feedback_submitted"
	PerformanceTimelineCalibrationNote   = "calibration_note"
	PerformanceTimelineObjectiveUpdate   = "objective_update"
)

var (
	ErrInvalidPerformanceCheckIn       = errors.New("performance check-in is invalid")
	ErrPerformanceCheckInNotFound      = errors.New("performance check-in not found")
	ErrInvalidFeedbackRequest          = errors.New("feedback request is invalid")
	ErrFeedbackRequestNotFound         = errors.New("feedback request not found")
	ErrInvalidFeedbackResponse         = errors.New("feedback response is invalid")
	ErrInvalidPerformanceTimelineEvent = errors.New("performance timeline event is invalid")
)

type PerformanceCheckIn struct {
	ID                      uuid.UUID       `json:"id"`
	TenantID                uuid.UUID       `json:"tenant_id"`
	WorkerProfileID         uuid.UUID       `json:"worker_profile_id"`
	ReviewerWorkerProfileID *uuid.UUID      `json:"reviewer_worker_profile_id,omitempty"`
	CycleID                 *uuid.UUID      `json:"cycle_id,omitempty"`
	CheckInDate             time.Time       `json:"checkin_date"`
	PeriodStart             time.Time       `json:"period_start"`
	PeriodEnd               time.Time       `json:"period_end"`
	Mood                    string          `json:"mood"`
	Status                  string          `json:"status"`
	Visibility              string          `json:"visibility"`
	Highlights              *string         `json:"highlights,omitempty"`
	Blockers                *string         `json:"blockers,omitempty"`
	NextPlan                *string         `json:"next_plan,omitempty"`
	EmployeeComment         *string         `json:"employee_comment,omitempty"`
	ManagerComment          *string         `json:"manager_comment,omitempty"`
	Score                   *float64        `json:"score,omitempty"`
	CalibrationBucket       *string         `json:"calibration_bucket,omitempty"`
	ReviewedAt              *time.Time      `json:"reviewed_at,omitempty"`
	ReviewedBy              *uuid.UUID      `json:"reviewed_by,omitempty"`
	Metadata                json.RawMessage `json:"metadata,omitempty"`
	Inactive                bool            `json:"inactive"`
	CreatedAt               time.Time       `json:"created_at"`
	CreatedBy               *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt               time.Time       `json:"updated_at"`
	UpdatedBy               *uuid.UUID      `json:"updated_by,omitempty"`
	WorkerDisplayName       *string         `json:"worker_display_name,omitempty"`
	WorkerCode              *string         `json:"worker_code,omitempty"`
	ReviewerDisplayName     *string         `json:"reviewer_display_name,omitempty"`
	CycleName               *string         `json:"cycle_name,omitempty"`
	FeedbackCount           int32           `json:"feedback_count"`
	AverageFeedbackRating   float64         `json:"average_feedback_rating"`
}

type PerformanceCheckInInput struct {
	TenantID                uuid.UUID
	WorkerProfileID         uuid.UUID
	ReviewerWorkerProfileID *uuid.UUID
	CycleID                 *uuid.UUID
	CheckInDate             *time.Time
	PeriodStart             *time.Time
	PeriodEnd               *time.Time
	Mood                    string
	Status                  string
	Visibility              string
	Highlights              *string
	Blockers                *string
	NextPlan                *string
	EmployeeComment         *string
	ManagerComment          *string
	Score                   *float64
	CalibrationBucket       *string
	Metadata                json.RawMessage
}

type PerformanceCheckInFilter struct {
	TenantID                uuid.UUID
	WorkerProfileID         *uuid.UUID
	ReviewerWorkerProfileID *uuid.UUID
	CycleID                 *uuid.UUID
	Status                  *string
	Mood                    *string
}

type PerformanceCheckInSummaryRow struct {
	Status       string  `json:"status"`
	Mood         string  `json:"mood"`
	CheckInCount int32   `json:"checkin_count"`
	AverageScore float64 `json:"average_score"`
}

type FeedbackRequest struct {
	ID                       uuid.UUID       `json:"id"`
	TenantID                 uuid.UUID       `json:"tenant_id"`
	SubjectWorkerProfileID   uuid.UUID       `json:"subject_worker_profile_id"`
	RequesterWorkerProfileID *uuid.UUID      `json:"requester_worker_profile_id,omitempty"`
	ObjectiveID              *uuid.UUID      `json:"objective_id,omitempty"`
	Relationship             string          `json:"relationship"`
	FeedbackType             string          `json:"feedback_type"`
	Status                   string          `json:"status"`
	IsAnonymous              bool            `json:"is_anonymous"`
	Visibility               string          `json:"visibility"`
	DueDate                  *time.Time      `json:"due_date,omitempty"`
	Prompt                   *string         `json:"prompt,omitempty"`
	Metadata                 json.RawMessage `json:"metadata,omitempty"`
	Inactive                 bool            `json:"inactive"`
	CreatedAt                time.Time       `json:"created_at"`
	CreatedBy                *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt                time.Time       `json:"updated_at"`
	UpdatedBy                *uuid.UUID      `json:"updated_by,omitempty"`
	SubjectDisplayName       *string         `json:"subject_display_name,omitempty"`
	SubjectWorkerCode        *string         `json:"subject_worker_code,omitempty"`
	RequesterDisplayName     *string         `json:"requester_display_name,omitempty"`
	ObjectiveTitle           *string         `json:"objective_title,omitempty"`
	ResponseCount            int32           `json:"response_count"`
}

type FeedbackRequestInput struct {
	TenantID                 uuid.UUID
	SubjectWorkerProfileID   uuid.UUID
	RequesterWorkerProfileID *uuid.UUID
	ObjectiveID              *uuid.UUID
	Relationship             string
	FeedbackType             string
	Status                   string
	IsAnonymous              bool
	Visibility               string
	DueDate                  *time.Time
	Prompt                   *string
	Metadata                 json.RawMessage
}

type FeedbackRequestFilter struct {
	TenantID                 uuid.UUID
	SubjectWorkerProfileID   *uuid.UUID
	RequesterWorkerProfileID *uuid.UUID
	Status                   *string
	FeedbackType             *string
}

type FeedbackResponse struct {
	ID                        uuid.UUID       `json:"id"`
	TenantID                  uuid.UUID       `json:"tenant_id"`
	RequestID                 uuid.UUID       `json:"request_id"`
	RespondentWorkerProfileID *uuid.UUID      `json:"respondent_worker_profile_id,omitempty"`
	Rating                    *float64        `json:"rating,omitempty"`
	Strengths                 *string         `json:"strengths,omitempty"`
	Improvements              *string         `json:"improvements,omitempty"`
	Comments                  *string         `json:"comments,omitempty"`
	SubmittedAt               time.Time       `json:"submitted_at"`
	Metadata                  json.RawMessage `json:"metadata,omitempty"`
	Inactive                  bool            `json:"inactive"`
	CreatedAt                 time.Time       `json:"created_at"`
	CreatedBy                 *uuid.UUID      `json:"created_by,omitempty"`
	SubjectWorkerProfileID    *uuid.UUID      `json:"subject_worker_profile_id,omitempty"`
	IsAnonymous               bool            `json:"is_anonymous"`
	SubjectDisplayName        *string         `json:"subject_display_name,omitempty"`
	SubjectWorkerCode         *string         `json:"subject_worker_code,omitempty"`
	RespondentDisplayName     *string         `json:"respondent_display_name,omitempty"`
	FeedbackType              *string         `json:"feedback_type,omitempty"`
	Relationship              *string         `json:"relationship,omitempty"`
}

type FeedbackResponseInput struct {
	TenantID                  uuid.UUID
	RequestID                 uuid.UUID
	RespondentWorkerProfileID *uuid.UUID
	Rating                    *float64
	Strengths                 *string
	Improvements              *string
	Comments                  *string
	Metadata                  json.RawMessage
}

type FeedbackResponseFilter struct {
	TenantID                  uuid.UUID
	RequestID                 *uuid.UUID
	SubjectWorkerProfileID    *uuid.UUID
	RespondentWorkerProfileID *uuid.UUID
}

type PerformanceTimelineEvent struct {
	ID                   uuid.UUID       `json:"id"`
	TenantID             uuid.UUID       `json:"tenant_id"`
	WorkerProfileID      uuid.UUID       `json:"worker_profile_id"`
	EventType            string          `json:"event_type"`
	CheckInID            *uuid.UUID      `json:"checkin_id,omitempty"`
	FeedbackRequestID    *uuid.UUID      `json:"feedback_request_id,omitempty"`
	FeedbackResponseID   *uuid.UUID      `json:"feedback_response_id,omitempty"`
	ObjectiveID          *uuid.UUID      `json:"objective_id,omitempty"`
	ActorWorkerProfileID *uuid.UUID      `json:"actor_worker_profile_id,omitempty"`
	Title                string          `json:"title"`
	Notes                *string         `json:"notes,omitempty"`
	Metadata             json.RawMessage `json:"metadata,omitempty"`
	Inactive             bool            `json:"inactive"`
	CreatedAt            time.Time       `json:"created_at"`
	CreatedBy            *uuid.UUID      `json:"created_by,omitempty"`
	WorkerDisplayName    *string         `json:"worker_display_name,omitempty"`
	ActorDisplayName     *string         `json:"actor_display_name,omitempty"`
	ObjectiveTitle       *string         `json:"objective_title,omitempty"`
}

type PerformanceTimelineEventInput struct {
	TenantID             uuid.UUID
	WorkerProfileID      uuid.UUID
	EventType            string
	CheckInID            *uuid.UUID
	FeedbackRequestID    *uuid.UUID
	FeedbackResponseID   *uuid.UUID
	ObjectiveID          *uuid.UUID
	ActorWorkerProfileID *uuid.UUID
	Title                string
	Notes                *string
	Metadata             json.RawMessage
}

type PerformanceTimelineFilter struct {
	TenantID        uuid.UUID
	WorkerProfileID *uuid.UUID
	EventType       *string
}

type PerformanceCalibrationRow struct {
	WorkerProfileID       uuid.UUID  `json:"worker_profile_id"`
	WorkerDisplayName     string     `json:"worker_display_name"`
	WorkerCode            *string    `json:"worker_code,omitempty"`
	CycleID               *uuid.UUID `json:"cycle_id,omitempty"`
	CycleName             *string    `json:"cycle_name,omitempty"`
	CheckInCount          int32      `json:"checkin_count"`
	SubmittedCheckInCount int32      `json:"submitted_checkin_count"`
	AverageScore          float64    `json:"average_score"`
	CalibrationBucket     *string    `json:"calibration_bucket,omitempty"`
	AverageOKRProgress    float64    `json:"average_okr_progress"`
	FeedbackCount         int32      `json:"feedback_count"`
	AverageFeedbackRating float64    `json:"average_feedback_rating"`
}

type PerformanceCalibrationFilter struct {
	TenantID        uuid.UUID
	CycleID         *uuid.UUID
	WorkerProfileID *uuid.UUID
}

func NewPerformanceCheckIn(input PerformanceCheckInInput) (*PerformanceCheckIn, error) {
	if input.TenantID == uuid.Nil || input.WorkerProfileID == uuid.Nil || input.PeriodStart == nil || input.PeriodEnd == nil || input.PeriodEnd.Before(*input.PeriodStart) {
		return nil, ErrInvalidPerformanceCheckIn
	}
	checkInDate := time.Now()
	if input.CheckInDate != nil {
		checkInDate = *input.CheckInDate
	}
	mood := normalizeWorkerProfileEnum(input.Mood, PerformanceCheckInMoodNeutral)
	if !containsString([]string{PerformanceCheckInMoodGreat, PerformanceCheckInMoodGood, PerformanceCheckInMoodNeutral, PerformanceCheckInMoodLow, PerformanceCheckInMoodStressed}, mood) {
		return nil, ErrInvalidPerformanceCheckIn
	}
	status := normalizeWorkerProfileEnum(input.Status, PerformanceCheckInStatusDraft)
	if !containsString([]string{PerformanceCheckInStatusDraft, PerformanceCheckInStatusSubmitted, PerformanceCheckInStatusReviewed, PerformanceCheckInStatusClosed}, status) {
		return nil, ErrInvalidPerformanceCheckIn
	}
	visibility := normalizeWorkerProfileEnum(input.Visibility, PerformanceVisibilityWorkerManagerHR)
	if !containsString([]string{PerformanceVisibilityWorkerManagerHR, PerformanceVisibilityManagerHR, PerformanceVisibilityHROnly}, visibility) {
		return nil, ErrInvalidPerformanceCheckIn
	}
	score := cleanPerformanceScore(input.Score)
	bucket := cleanPerformanceBucket(input.CalibrationBucket)
	return &PerformanceCheckIn{TenantID: input.TenantID, WorkerProfileID: input.WorkerProfileID, ReviewerWorkerProfileID: cleanOKRUUIDPtr(input.ReviewerWorkerProfileID), CycleID: cleanOKRUUIDPtr(input.CycleID), CheckInDate: checkInDate, PeriodStart: *input.PeriodStart, PeriodEnd: *input.PeriodEnd, Mood: mood, Status: status, Visibility: visibility, Highlights: cleanOKRStringPtr(input.Highlights), Blockers: cleanOKRStringPtr(input.Blockers), NextPlan: cleanOKRStringPtr(input.NextPlan), EmployeeComment: cleanOKRStringPtr(input.EmployeeComment), ManagerComment: cleanOKRStringPtr(input.ManagerComment), Score: score, CalibrationBucket: bucket, Metadata: cleanOKRObjectJSON(input.Metadata)}, nil
}

func NewFeedbackRequest(input FeedbackRequestInput) (*FeedbackRequest, error) {
	if input.TenantID == uuid.Nil || input.SubjectWorkerProfileID == uuid.Nil {
		return nil, ErrInvalidFeedbackRequest
	}
	relationship := normalizeWorkerProfileEnum(input.Relationship, FeedbackRelationshipPeer)
	if !containsString([]string{FeedbackRelationshipManager, FeedbackRelationshipPeer, FeedbackRelationshipDirectReport, FeedbackRelationshipSelf, FeedbackRelationshipHR, FeedbackRelationshipClient}, relationship) {
		return nil, ErrInvalidFeedbackRequest
	}
	feedbackType := normalizeWorkerProfileEnum(input.FeedbackType, FeedbackType360)
	if !containsString([]string{FeedbackType360, FeedbackTypeProject, FeedbackTypeGeneral, FeedbackTypeOKR, FeedbackTypeManagerReview}, feedbackType) {
		return nil, ErrInvalidFeedbackRequest
	}
	status := normalizeWorkerProfileEnum(input.Status, FeedbackStatusRequested)
	if !containsString([]string{FeedbackStatusRequested, FeedbackStatusSubmitted, FeedbackStatusDeclined, FeedbackStatusExpired, FeedbackStatusCancelled}, status) {
		return nil, ErrInvalidFeedbackRequest
	}
	visibility := normalizeWorkerProfileEnum(input.Visibility, FeedbackVisibilitySubjectManagerHR)
	if !containsString([]string{FeedbackVisibilitySubjectManagerHR, FeedbackVisibilityManagerHR, FeedbackVisibilityHROnly, FeedbackVisibilitySubjectOnly}, visibility) {
		return nil, ErrInvalidFeedbackRequest
	}
	return &FeedbackRequest{TenantID: input.TenantID, SubjectWorkerProfileID: input.SubjectWorkerProfileID, RequesterWorkerProfileID: cleanOKRUUIDPtr(input.RequesterWorkerProfileID), ObjectiveID: cleanOKRUUIDPtr(input.ObjectiveID), Relationship: relationship, FeedbackType: feedbackType, Status: status, IsAnonymous: input.IsAnonymous, Visibility: visibility, DueDate: input.DueDate, Prompt: cleanOKRStringPtr(input.Prompt), Metadata: cleanOKRObjectJSON(input.Metadata)}, nil
}

func NewFeedbackResponse(input FeedbackResponseInput) (*FeedbackResponse, error) {
	if input.TenantID == uuid.Nil || input.RequestID == uuid.Nil {
		return nil, ErrInvalidFeedbackResponse
	}
	return &FeedbackResponse{TenantID: input.TenantID, RequestID: input.RequestID, RespondentWorkerProfileID: cleanOKRUUIDPtr(input.RespondentWorkerProfileID), Rating: cleanPerformanceScore(input.Rating), Strengths: cleanOKRStringPtr(input.Strengths), Improvements: cleanOKRStringPtr(input.Improvements), Comments: cleanOKRStringPtr(input.Comments), Metadata: cleanOKRObjectJSON(input.Metadata)}, nil
}

func NewPerformanceTimelineEvent(input PerformanceTimelineEventInput) (*PerformanceTimelineEvent, error) {
	title := strings.TrimSpace(input.Title)
	if input.TenantID == uuid.Nil || input.WorkerProfileID == uuid.Nil || title == "" {
		return nil, ErrInvalidPerformanceTimelineEvent
	}
	eventType := normalizeWorkerProfileEnum(input.EventType, "")
	if !containsString([]string{PerformanceTimelineCheckInCreated, PerformanceTimelineCheckInSubmitted, PerformanceTimelineCheckInReviewed, PerformanceTimelineFeedbackRequested, PerformanceTimelineFeedbackSubmitted, PerformanceTimelineCalibrationNote, PerformanceTimelineObjectiveUpdate}, eventType) {
		return nil, ErrInvalidPerformanceTimelineEvent
	}
	return &PerformanceTimelineEvent{TenantID: input.TenantID, WorkerProfileID: input.WorkerProfileID, EventType: eventType, CheckInID: cleanOKRUUIDPtr(input.CheckInID), FeedbackRequestID: cleanOKRUUIDPtr(input.FeedbackRequestID), FeedbackResponseID: cleanOKRUUIDPtr(input.FeedbackResponseID), ObjectiveID: cleanOKRUUIDPtr(input.ObjectiveID), ActorWorkerProfileID: cleanOKRUUIDPtr(input.ActorWorkerProfileID), Title: title, Notes: cleanOKRStringPtr(input.Notes), Metadata: cleanOKRObjectJSON(input.Metadata)}, nil
}

func NormalizePerformanceSearch(value *string) *string {
	if value == nil {
		return nil
	}
	clean := strings.TrimSpace(*value)
	if clean == "" {
		return nil
	}
	return &clean
}

func cleanPerformanceScore(value *float64) *float64 {
	if value == nil {
		return nil
	}
	clean := *value
	if clean < 0 {
		clean = 0
	}
	if clean > 5 {
		clean = 5
	}
	return &clean
}

func cleanPerformanceBucket(value *string) *string {
	if value == nil {
		return nil
	}
	clean := strings.TrimSpace(*value)
	if clean == "" {
		return nil
	}
	if !containsString([]string{PerformanceCalibrationHigh, PerformanceCalibrationSolid, PerformanceCalibrationWatch, PerformanceCalibrationImprove}, clean) {
		return nil
	}
	return &clean
}
