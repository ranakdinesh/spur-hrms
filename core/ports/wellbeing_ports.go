package ports

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

type WellbeingRepo interface {
	CreatePulseSurvey(ctx context.Context, item *domain.PulseSurvey, actorID *uuid.UUID) (*domain.PulseSurvey, error)
	UpdatePulseSurvey(ctx context.Context, item *domain.PulseSurvey, actorID *uuid.UUID) (*domain.PulseSurvey, error)
	UpdatePulseSurveyStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, actorID *uuid.UUID) (*domain.PulseSurvey, error)
	GetPulseSurvey(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.PulseSurvey, error)
	ListPulseSurveys(ctx context.Context, filter domain.PulseSurveyFilter) ([]*domain.PulseSurvey, error)
	DeletePulseSurvey(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	CreatePulseQuestion(ctx context.Context, item *domain.PulseQuestion, actorID *uuid.UUID) (*domain.PulseQuestion, error)
	UpdatePulseQuestion(ctx context.Context, item *domain.PulseQuestion, actorID *uuid.UUID) (*domain.PulseQuestion, error)
	GetPulseQuestion(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.PulseQuestion, error)
	ListPulseQuestions(ctx context.Context, tenantID uuid.UUID, surveyID *uuid.UUID) ([]*domain.PulseQuestion, error)
	DeletePulseQuestion(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	CreatePulseResponse(ctx context.Context, item *domain.PulseResponse, actorID *uuid.UUID) (*domain.PulseResponse, error)
	ListPulseResponses(ctx context.Context, filter domain.PulseResponseFilter) ([]*domain.PulseResponse, error)
	UpsertWellbeingScore(ctx context.Context, item *domain.WellbeingScore, actorID *uuid.UUID) (*domain.WellbeingScore, error)
	ListWellbeingScores(ctx context.Context, filter domain.WellbeingScoreFilter) ([]*domain.WellbeingScore, error)
	CreateWellbeingAlert(ctx context.Context, item *domain.WellbeingAlert, actorID *uuid.UUID) (*domain.WellbeingAlert, error)
	UpdateWellbeingAlertStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, resolutionNote *string, actorID *uuid.UUID) (*domain.WellbeingAlert, error)
	ListWellbeingAlerts(ctx context.Context, filter domain.WellbeingAlertFilter) ([]*domain.WellbeingAlert, error)
	ListWellbeingAggregateRows(ctx context.Context, tenantID uuid.UUID, surveyID *uuid.UUID) ([]*domain.WellbeingAggregateRow, error)
}

type PulseSurveyCommand struct {
	ID                    uuid.UUID       `json:"id,omitempty"`
	TenantID              uuid.UUID       `json:"tenant_id"`
	Title                 string          `json:"title"`
	Description           *string         `json:"description,omitempty"`
	SurveyType            string          `json:"survey_type"`
	Status                string          `json:"status"`
	AudienceScope         string          `json:"audience_scope"`
	DepartmentID          *uuid.UUID      `json:"department_id,omitempty"`
	StartDate             string          `json:"start_date"`
	EndDate               string          `json:"end_date"`
	Frequency             string          `json:"frequency"`
	AnonymityThreshold    int32           `json:"anonymity_threshold"`
	ConsentRequired       bool            `json:"consent_required"`
	ManagerAggregateOnly  bool            `json:"manager_aggregate_only"`
	CriticalAlertsEnabled bool            `json:"critical_alerts_enabled"`
	Metadata              json.RawMessage `json:"metadata,omitempty"`
	ActorID               *uuid.UUID      `json:"-"`
}

type PulseSurveyStatusCommand struct {
	TenantID uuid.UUID  `json:"tenant_id"`
	ID       uuid.UUID  `json:"id"`
	Status   string     `json:"status"`
	ActorID  *uuid.UUID `json:"-"`
}

type PulseQuestionCommand struct {
	ID           uuid.UUID       `json:"id,omitempty"`
	TenantID     uuid.UUID       `json:"tenant_id"`
	SurveyID     uuid.UUID       `json:"survey_id"`
	QuestionText string          `json:"question_text"`
	QuestionType string          `json:"question_type"`
	Category     string          `json:"category"`
	IsRequired   bool            `json:"is_required"`
	SortOrder    int32           `json:"sort_order"`
	Options      json.RawMessage `json:"options,omitempty"`
	ActorID      *uuid.UUID      `json:"-"`
}

type PulseResponseCommand struct {
	TenantID        uuid.UUID       `json:"tenant_id"`
	SurveyID        uuid.UUID       `json:"survey_id"`
	QuestionID      uuid.UUID       `json:"question_id"`
	WorkerProfileID *uuid.UUID      `json:"worker_profile_id,omitempty"`
	ResponseDate    string          `json:"response_date"`
	Score           *float64        `json:"score,omitempty"`
	TextResponse    *string         `json:"text_response,omitempty"`
	BooleanResponse *bool           `json:"boolean_response,omitempty"`
	OptionValue     *string         `json:"option_value,omitempty"`
	ConsentGiven    bool            `json:"consent_given"`
	IsAnonymous     bool            `json:"is_anonymous"`
	RiskLevel       string          `json:"risk_level"`
	CriticalAlert   bool            `json:"critical_alert"`
	Metadata        json.RawMessage `json:"metadata,omitempty"`
	ActorID         *uuid.UUID      `json:"-"`
}

type WellbeingScoreCommand struct {
	TenantID        uuid.UUID       `json:"tenant_id"`
	WorkerProfileID uuid.UUID       `json:"worker_profile_id"`
	ScoreDate       string          `json:"score_date"`
	SourceSurveyID  *uuid.UUID      `json:"source_survey_id,omitempty"`
	WellbeingScore  float64         `json:"wellbeing_score"`
	MoodScore       *float64        `json:"mood_score,omitempty"`
	StressScore     *float64        `json:"stress_score,omitempty"`
	WorkloadScore   *float64        `json:"workload_score,omitempty"`
	RiskLevel       string          `json:"risk_level"`
	ConsentScope    string          `json:"consent_scope"`
	Notes           *string         `json:"notes,omitempty"`
	Metadata        json.RawMessage `json:"metadata,omitempty"`
	ActorID         *uuid.UUID      `json:"-"`
}

type WellbeingAlertStatusCommand struct {
	TenantID       uuid.UUID  `json:"tenant_id"`
	ID             uuid.UUID  `json:"id"`
	Status         string     `json:"status"`
	ResolutionNote *string    `json:"resolution_note,omitempty"`
	ActorID        *uuid.UUID `json:"-"`
}
