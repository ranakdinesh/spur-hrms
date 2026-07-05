package ports

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

type PerformanceRepo interface {
	CreatePerformanceCheckIn(ctx context.Context, item *domain.PerformanceCheckIn, actorID *uuid.UUID) (*domain.PerformanceCheckIn, error)
	UpdatePerformanceCheckIn(ctx context.Context, item *domain.PerformanceCheckIn, actorID *uuid.UUID) (*domain.PerformanceCheckIn, error)
	ReviewPerformanceCheckIn(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, managerComment *string, score *float64, calibrationBucket *string, actorID *uuid.UUID) (*domain.PerformanceCheckIn, error)
	UpdatePerformanceCheckInStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, actorID *uuid.UUID) (*domain.PerformanceCheckIn, error)
	GetPerformanceCheckIn(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.PerformanceCheckIn, error)
	ListPerformanceCheckIns(ctx context.Context, filter domain.PerformanceCheckInFilter) ([]*domain.PerformanceCheckIn, error)
	DeletePerformanceCheckIn(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	GetPerformanceCheckInSummary(ctx context.Context, tenantID uuid.UUID, cycleID *uuid.UUID) ([]*domain.PerformanceCheckInSummaryRow, error)
	CreateFeedbackRequest(ctx context.Context, item *domain.FeedbackRequest, actorID *uuid.UUID) (*domain.FeedbackRequest, error)
	UpdateFeedbackRequest(ctx context.Context, item *domain.FeedbackRequest, actorID *uuid.UUID) (*domain.FeedbackRequest, error)
	UpdateFeedbackRequestStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, actorID *uuid.UUID) (*domain.FeedbackRequest, error)
	GetFeedbackRequest(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.FeedbackRequest, error)
	ListFeedbackRequests(ctx context.Context, filter domain.FeedbackRequestFilter) ([]*domain.FeedbackRequest, error)
	CreateFeedbackResponse(ctx context.Context, item *domain.FeedbackResponse, actorID *uuid.UUID) (*domain.FeedbackResponse, error)
	ListFeedbackResponses(ctx context.Context, filter domain.FeedbackResponseFilter) ([]*domain.FeedbackResponse, error)
	CreatePerformanceTimelineEvent(ctx context.Context, item *domain.PerformanceTimelineEvent, actorID *uuid.UUID) (*domain.PerformanceTimelineEvent, error)
	ListPerformanceTimelineEvents(ctx context.Context, filter domain.PerformanceTimelineFilter) ([]*domain.PerformanceTimelineEvent, error)
	ListPerformanceCalibrationRows(ctx context.Context, filter domain.PerformanceCalibrationFilter) ([]*domain.PerformanceCalibrationRow, error)
}

type PerformanceCheckInCommand struct {
	ID                      uuid.UUID       `json:"id,omitempty"`
	TenantID                uuid.UUID       `json:"tenant_id"`
	WorkerProfileID         uuid.UUID       `json:"worker_profile_id"`
	ReviewerWorkerProfileID *uuid.UUID      `json:"reviewer_worker_profile_id,omitempty"`
	CycleID                 *uuid.UUID      `json:"cycle_id,omitempty"`
	CheckInDate             string          `json:"checkin_date"`
	PeriodStart             string          `json:"period_start"`
	PeriodEnd               string          `json:"period_end"`
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
	Metadata                json.RawMessage `json:"metadata,omitempty"`
	ActorID                 *uuid.UUID      `json:"-"`
}

type PerformanceCheckInReviewCommand struct {
	TenantID          uuid.UUID  `json:"tenant_id"`
	ID                uuid.UUID  `json:"id"`
	Status            string     `json:"status"`
	ManagerComment    *string    `json:"manager_comment,omitempty"`
	Score             *float64   `json:"score,omitempty"`
	CalibrationBucket *string    `json:"calibration_bucket,omitempty"`
	ActorID           *uuid.UUID `json:"-"`
}

type PerformanceStatusCommand struct {
	TenantID uuid.UUID  `json:"tenant_id"`
	ID       uuid.UUID  `json:"id"`
	Status   string     `json:"status"`
	ActorID  *uuid.UUID `json:"-"`
}

type FeedbackRequestCommand struct {
	ID                       uuid.UUID       `json:"id,omitempty"`
	TenantID                 uuid.UUID       `json:"tenant_id"`
	SubjectWorkerProfileID   uuid.UUID       `json:"subject_worker_profile_id"`
	RequesterWorkerProfileID *uuid.UUID      `json:"requester_worker_profile_id,omitempty"`
	ObjectiveID              *uuid.UUID      `json:"objective_id,omitempty"`
	Relationship             string          `json:"relationship"`
	FeedbackType             string          `json:"feedback_type"`
	Status                   string          `json:"status"`
	IsAnonymous              bool            `json:"is_anonymous"`
	Visibility               string          `json:"visibility"`
	DueDate                  string          `json:"due_date"`
	Prompt                   *string         `json:"prompt,omitempty"`
	Metadata                 json.RawMessage `json:"metadata,omitempty"`
	ActorID                  *uuid.UUID      `json:"-"`
}

type FeedbackStatusCommand struct {
	TenantID uuid.UUID  `json:"tenant_id"`
	ID       uuid.UUID  `json:"id"`
	Status   string     `json:"status"`
	ActorID  *uuid.UUID `json:"-"`
}

type FeedbackResponseCommand struct {
	TenantID                  uuid.UUID       `json:"tenant_id"`
	RequestID                 uuid.UUID       `json:"request_id"`
	RespondentWorkerProfileID *uuid.UUID      `json:"respondent_worker_profile_id,omitempty"`
	Rating                    *float64        `json:"rating,omitempty"`
	Strengths                 *string         `json:"strengths,omitempty"`
	Improvements              *string         `json:"improvements,omitempty"`
	Comments                  *string         `json:"comments,omitempty"`
	Metadata                  json.RawMessage `json:"metadata,omitempty"`
	ActorID                   *uuid.UUID      `json:"-"`
}
