package domain

import (
	"encoding/json"
	"errors"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	PulseSurveyTypePulse      = "pulse"
	PulseSurveyTypeWellbeing  = "wellbeing"
	PulseSurveyTypeEngagement = "engagement"

	PulseSurveyStatusDraft    = "draft"
	PulseSurveyStatusActive   = "active"
	PulseSurveyStatusClosed   = "closed"
	PulseSurveyStatusArchived = "archived"

	PulseAudienceAll        = "all"
	PulseAudienceDepartment = "department"
	PulseAudienceCustom     = "custom"

	PulseFrequencyOneTime  = "one_time"
	PulseFrequencyWeekly   = "weekly"
	PulseFrequencyBiweekly = "biweekly"
	PulseFrequencyMonthly  = "monthly"

	PulseQuestionScale        = "scale_1_5"
	PulseQuestionText         = "text"
	PulseQuestionBoolean      = "boolean"
	PulseQuestionSingleChoice = "single_choice"

	WellbeingCategoryMood           = "mood"
	WellbeingCategoryWorkload       = "workload"
	WellbeingCategoryStress         = "stress"
	WellbeingCategoryBelonging      = "belonging"
	WellbeingCategoryManagerSupport = "manager_support"
	WellbeingCategorySafety         = "safety"
	WellbeingCategoryGeneral        = "general"

	WellbeingRiskNone     = "none"
	WellbeingRiskLow      = "low"
	WellbeingRiskMedium   = "medium"
	WellbeingRiskHigh     = "high"
	WellbeingRiskCritical = "critical"

	WellbeingConsentPrivate   = "private"
	WellbeingConsentAggregate = "aggregate"
	WellbeingConsentHRAlert   = "hr_alert"

	WellbeingAlertCriticalResponse = "critical_response"
	WellbeingAlertLowScore         = "low_score"
	WellbeingAlertConsentIssue     = "consent_issue"

	WellbeingAlertSeverityMedium   = "medium"
	WellbeingAlertSeverityHigh     = "high"
	WellbeingAlertSeverityCritical = "critical"

	WellbeingAlertStatusOpen         = "open"
	WellbeingAlertStatusAcknowledged = "acknowledged"
	WellbeingAlertStatusResolved     = "resolved"
	WellbeingAlertStatusDismissed    = "dismissed"
)

var (
	ErrInvalidPulseSurvey     = errors.New("pulse survey is invalid")
	ErrPulseSurveyNotFound    = errors.New("pulse survey not found")
	ErrInvalidPulseQuestion   = errors.New("pulse survey question is invalid")
	ErrPulseQuestionNotFound  = errors.New("pulse survey question not found")
	ErrInvalidPulseResponse   = errors.New("pulse response is invalid")
	ErrInvalidWellbeingScore  = errors.New("wellbeing score is invalid")
	ErrInvalidWellbeingAlert  = errors.New("wellbeing alert is invalid")
	ErrWellbeingAlertNotFound = errors.New("wellbeing alert not found")
)

type PulseSurvey struct {
	ID                    uuid.UUID       `json:"id"`
	TenantID              uuid.UUID       `json:"tenant_id"`
	Title                 string          `json:"title"`
	Description           *string         `json:"description,omitempty"`
	SurveyType            string          `json:"survey_type"`
	Status                string          `json:"status"`
	AudienceScope         string          `json:"audience_scope"`
	DepartmentID          *uuid.UUID      `json:"department_id,omitempty"`
	StartDate             time.Time       `json:"start_date"`
	EndDate               *time.Time      `json:"end_date,omitempty"`
	Frequency             string          `json:"frequency"`
	AnonymityThreshold    int32           `json:"anonymity_threshold"`
	ConsentRequired       bool            `json:"consent_required"`
	ManagerAggregateOnly  bool            `json:"manager_aggregate_only"`
	CriticalAlertsEnabled bool            `json:"critical_alerts_enabled"`
	Metadata              json.RawMessage `json:"metadata,omitempty"`
	Inactive              bool            `json:"inactive"`
	CreatedAt             time.Time       `json:"created_at"`
	CreatedBy             *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt             time.Time       `json:"updated_at"`
	UpdatedBy             *uuid.UUID      `json:"updated_by,omitempty"`
	DepartmentName        *string         `json:"department_name,omitempty"`
	QuestionCount         int32           `json:"question_count"`
	RespondentCount       int32           `json:"respondent_count"`
	ResponseCount         int32           `json:"response_count"`
}

type PulseSurveyInput struct {
	TenantID              uuid.UUID
	Title                 string
	Description           *string
	SurveyType            string
	Status                string
	AudienceScope         string
	DepartmentID          *uuid.UUID
	StartDate             *time.Time
	EndDate               *time.Time
	Frequency             string
	AnonymityThreshold    int32
	ConsentRequired       bool
	ManagerAggregateOnly  bool
	CriticalAlertsEnabled bool
	Metadata              json.RawMessage
}

type PulseSurveyFilter struct {
	TenantID   uuid.UUID
	Status     *string
	SurveyType *string
	Search     *string
}

type PulseQuestion struct {
	ID           uuid.UUID       `json:"id"`
	TenantID     uuid.UUID       `json:"tenant_id"`
	SurveyID     uuid.UUID       `json:"survey_id"`
	QuestionText string          `json:"question_text"`
	QuestionType string          `json:"question_type"`
	Category     string          `json:"category"`
	IsRequired   bool            `json:"is_required"`
	SortOrder    int32           `json:"sort_order"`
	Options      json.RawMessage `json:"options,omitempty"`
	Inactive     bool            `json:"inactive"`
	CreatedAt    time.Time       `json:"created_at"`
	CreatedBy    *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt    time.Time       `json:"updated_at"`
	UpdatedBy    *uuid.UUID      `json:"updated_by,omitempty"`
}

type PulseQuestionInput struct {
	TenantID     uuid.UUID
	SurveyID     uuid.UUID
	QuestionText string
	QuestionType string
	Category     string
	IsRequired   bool
	SortOrder    int32
	Options      json.RawMessage
}

type PulseResponse struct {
	ID                uuid.UUID       `json:"id"`
	TenantID          uuid.UUID       `json:"tenant_id"`
	SurveyID          uuid.UUID       `json:"survey_id"`
	QuestionID        uuid.UUID       `json:"question_id"`
	WorkerProfileID   *uuid.UUID      `json:"worker_profile_id,omitempty"`
	ResponseDate      time.Time       `json:"response_date"`
	Score             *float64        `json:"score,omitempty"`
	TextResponse      *string         `json:"text_response,omitempty"`
	BooleanResponse   *bool           `json:"boolean_response,omitempty"`
	OptionValue       *string         `json:"option_value,omitempty"`
	ConsentGiven      bool            `json:"consent_given"`
	IsAnonymous       bool            `json:"is_anonymous"`
	RiskLevel         string          `json:"risk_level"`
	CriticalAlert     bool            `json:"critical_alert"`
	Metadata          json.RawMessage `json:"metadata,omitempty"`
	Inactive          bool            `json:"inactive"`
	CreatedAt         time.Time       `json:"created_at"`
	CreatedBy         *uuid.UUID      `json:"created_by,omitempty"`
	SurveyTitle       *string         `json:"survey_title,omitempty"`
	QuestionText      *string         `json:"question_text,omitempty"`
	Category          *string         `json:"category,omitempty"`
	WorkerDisplayName *string         `json:"worker_display_name,omitempty"`
	WorkerCode        *string         `json:"worker_code,omitempty"`
}

type PulseResponseInput struct {
	TenantID        uuid.UUID
	SurveyID        uuid.UUID
	QuestionID      uuid.UUID
	WorkerProfileID *uuid.UUID
	ResponseDate    *time.Time
	Score           *float64
	TextResponse    *string
	BooleanResponse *bool
	OptionValue     *string
	ConsentGiven    bool
	IsAnonymous     bool
	RiskLevel       string
	CriticalAlert   bool
	Metadata        json.RawMessage
}

type PulseResponseFilter struct {
	TenantID        uuid.UUID
	SurveyID        *uuid.UUID
	WorkerProfileID *uuid.UUID
	RiskLevel       *string
}

type WellbeingScore struct {
	ID                uuid.UUID       `json:"id"`
	TenantID          uuid.UUID       `json:"tenant_id"`
	WorkerProfileID   uuid.UUID       `json:"worker_profile_id"`
	ScoreDate         time.Time       `json:"score_date"`
	SourceSurveyID    *uuid.UUID      `json:"source_survey_id,omitempty"`
	WellbeingScore    float64         `json:"wellbeing_score"`
	MoodScore         *float64        `json:"mood_score,omitempty"`
	StressScore       *float64        `json:"stress_score,omitempty"`
	WorkloadScore     *float64        `json:"workload_score,omitempty"`
	RiskLevel         string          `json:"risk_level"`
	ConsentScope      string          `json:"consent_scope"`
	Notes             *string         `json:"notes,omitempty"`
	Metadata          json.RawMessage `json:"metadata,omitempty"`
	Inactive          bool            `json:"inactive"`
	CreatedAt         time.Time       `json:"created_at"`
	CreatedBy         *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt         time.Time       `json:"updated_at"`
	UpdatedBy         *uuid.UUID      `json:"updated_by,omitempty"`
	WorkerDisplayName *string         `json:"worker_display_name,omitempty"`
	WorkerCode        *string         `json:"worker_code,omitempty"`
	SurveyTitle       *string         `json:"survey_title,omitempty"`
}

type WellbeingScoreInput struct {
	TenantID        uuid.UUID
	WorkerProfileID uuid.UUID
	ScoreDate       *time.Time
	SourceSurveyID  *uuid.UUID
	WellbeingScore  float64
	MoodScore       *float64
	StressScore     *float64
	WorkloadScore   *float64
	RiskLevel       string
	ConsentScope    string
	Notes           *string
	Metadata        json.RawMessage
}

type WellbeingScoreFilter struct {
	TenantID        uuid.UUID
	WorkerProfileID *uuid.UUID
	RiskLevel       *string
}

type WellbeingAlert struct {
	ID                uuid.UUID       `json:"id"`
	TenantID          uuid.UUID       `json:"tenant_id"`
	WorkerProfileID   *uuid.UUID      `json:"worker_profile_id,omitempty"`
	SurveyID          *uuid.UUID      `json:"survey_id,omitempty"`
	ResponseID        *uuid.UUID      `json:"response_id,omitempty"`
	AlertType         string          `json:"alert_type"`
	Severity          string          `json:"severity"`
	Status            string          `json:"status"`
	VisibleTo         string          `json:"visible_to"`
	Message           string          `json:"message"`
	ResolutionNote    *string         `json:"resolution_note,omitempty"`
	Metadata          json.RawMessage `json:"metadata,omitempty"`
	Inactive          bool            `json:"inactive"`
	CreatedAt         time.Time       `json:"created_at"`
	CreatedBy         *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt         time.Time       `json:"updated_at"`
	UpdatedBy         *uuid.UUID      `json:"updated_by,omitempty"`
	WorkerDisplayName *string         `json:"worker_display_name,omitempty"`
	WorkerCode        *string         `json:"worker_code,omitempty"`
	SurveyTitle       *string         `json:"survey_title,omitempty"`
}

type WellbeingAlertInput struct {
	TenantID        uuid.UUID
	WorkerProfileID *uuid.UUID
	SurveyID        *uuid.UUID
	ResponseID      *uuid.UUID
	AlertType       string
	Severity        string
	Status          string
	Message         string
	Metadata        json.RawMessage
}

type WellbeingAlertFilter struct {
	TenantID uuid.UUID
	Status   *string
	Severity *string
}

type WellbeingAggregateRow struct {
	SurveyID           uuid.UUID `json:"survey_id"`
	SurveyTitle        string    `json:"survey_title"`
	DepartmentID       uuid.UUID `json:"department_id"`
	DepartmentName     string    `json:"department_name"`
	Category           string    `json:"category"`
	ResponseCount      int32     `json:"response_count"`
	RespondentCount    int32     `json:"respondent_count"`
	Suppressed         bool      `json:"suppressed"`
	AverageScore       *float64  `json:"average_score,omitempty"`
	RiskCount          int32     `json:"risk_count"`
	AnonymityThreshold int32     `json:"anonymity_threshold"`
}

func NewPulseSurvey(input PulseSurveyInput) (*PulseSurvey, error) {
	if input.TenantID == uuid.Nil || input.StartDate == nil {
		return nil, ErrInvalidPulseSurvey
	}
	title := strings.TrimSpace(input.Title)
	if title == "" {
		return nil, ErrInvalidPulseSurvey
	}
	if input.EndDate != nil && input.EndDate.Before(*input.StartDate) {
		return nil, ErrInvalidPulseSurvey
	}
	surveyType := normalizeWorkerProfileEnum(input.SurveyType, PulseSurveyTypePulse)
	if !containsString([]string{PulseSurveyTypePulse, PulseSurveyTypeWellbeing, PulseSurveyTypeEngagement}, surveyType) {
		return nil, ErrInvalidPulseSurvey
	}
	status := normalizeWorkerProfileEnum(input.Status, PulseSurveyStatusDraft)
	if !containsString([]string{PulseSurveyStatusDraft, PulseSurveyStatusActive, PulseSurveyStatusClosed, PulseSurveyStatusArchived}, status) {
		return nil, ErrInvalidPulseSurvey
	}
	audience := normalizeWorkerProfileEnum(input.AudienceScope, PulseAudienceAll)
	if !containsString([]string{PulseAudienceAll, PulseAudienceDepartment, PulseAudienceCustom}, audience) {
		return nil, ErrInvalidPulseSurvey
	}
	if audience == PulseAudienceDepartment && (input.DepartmentID == nil || *input.DepartmentID == uuid.Nil) {
		return nil, ErrInvalidPulseSurvey
	}
	frequency := normalizeWorkerProfileEnum(input.Frequency, PulseFrequencyOneTime)
	if !containsString([]string{PulseFrequencyOneTime, PulseFrequencyWeekly, PulseFrequencyBiweekly, PulseFrequencyMonthly}, frequency) {
		return nil, ErrInvalidPulseSurvey
	}
	threshold := input.AnonymityThreshold
	if threshold < 3 {
		threshold = 5
	}
	return &PulseSurvey{TenantID: input.TenantID, Title: title, Description: cleanOKRStringPtr(input.Description), SurveyType: surveyType, Status: status, AudienceScope: audience, DepartmentID: cleanOKRUUIDPtr(input.DepartmentID), StartDate: *input.StartDate, EndDate: input.EndDate, Frequency: frequency, AnonymityThreshold: threshold, ConsentRequired: input.ConsentRequired, ManagerAggregateOnly: input.ManagerAggregateOnly, CriticalAlertsEnabled: input.CriticalAlertsEnabled, Metadata: cleanOKRObjectJSON(input.Metadata)}, nil
}

func NewPulseQuestion(input PulseQuestionInput) (*PulseQuestion, error) {
	text := strings.TrimSpace(input.QuestionText)
	if input.TenantID == uuid.Nil || input.SurveyID == uuid.Nil || text == "" {
		return nil, ErrInvalidPulseQuestion
	}
	questionType := normalizeWorkerProfileEnum(input.QuestionType, PulseQuestionScale)
	if !containsString([]string{PulseQuestionScale, PulseQuestionText, PulseQuestionBoolean, PulseQuestionSingleChoice}, questionType) {
		return nil, ErrInvalidPulseQuestion
	}
	category := normalizeWorkerProfileEnum(input.Category, WellbeingCategoryGeneral)
	if !containsString([]string{WellbeingCategoryMood, WellbeingCategoryWorkload, WellbeingCategoryStress, WellbeingCategoryBelonging, WellbeingCategoryManagerSupport, WellbeingCategorySafety, WellbeingCategoryGeneral}, category) {
		return nil, ErrInvalidPulseQuestion
	}
	return &PulseQuestion{TenantID: input.TenantID, SurveyID: input.SurveyID, QuestionText: text, QuestionType: questionType, Category: category, IsRequired: input.IsRequired, SortOrder: input.SortOrder, Options: cleanJSONArray(input.Options)}, nil
}

func NewPulseResponse(input PulseResponseInput) (*PulseResponse, error) {
	if input.TenantID == uuid.Nil || input.SurveyID == uuid.Nil || input.QuestionID == uuid.Nil {
		return nil, ErrInvalidPulseResponse
	}
	responseDate := time.Now()
	if input.ResponseDate != nil {
		responseDate = *input.ResponseDate
	}
	risk := normalizeWorkerProfileEnum(input.RiskLevel, deriveWellbeingRisk(input.Score))
	if !containsString([]string{WellbeingRiskNone, WellbeingRiskLow, WellbeingRiskMedium, WellbeingRiskHigh, WellbeingRiskCritical}, risk) {
		return nil, ErrInvalidPulseResponse
	}
	return &PulseResponse{TenantID: input.TenantID, SurveyID: input.SurveyID, QuestionID: input.QuestionID, WorkerProfileID: cleanOKRUUIDPtr(input.WorkerProfileID), ResponseDate: responseDate, Score: cleanScaleScore(input.Score), TextResponse: cleanOKRStringPtr(input.TextResponse), BooleanResponse: input.BooleanResponse, OptionValue: cleanOKRStringPtr(input.OptionValue), ConsentGiven: input.ConsentGiven, IsAnonymous: input.IsAnonymous, RiskLevel: risk, CriticalAlert: input.CriticalAlert || risk == WellbeingRiskCritical, Metadata: cleanOKRObjectJSON(input.Metadata)}, nil
}

func NewWellbeingScore(input WellbeingScoreInput) (*WellbeingScore, error) {
	if input.TenantID == uuid.Nil || input.WorkerProfileID == uuid.Nil {
		return nil, ErrInvalidWellbeingScore
	}
	scoreDate := time.Now()
	if input.ScoreDate != nil {
		scoreDate = *input.ScoreDate
	}
	risk := normalizeWorkerProfileEnum(input.RiskLevel, deriveWellbeingScoreRisk(input.WellbeingScore))
	if !containsString([]string{WellbeingRiskNone, WellbeingRiskLow, WellbeingRiskMedium, WellbeingRiskHigh, WellbeingRiskCritical}, risk) {
		return nil, ErrInvalidWellbeingScore
	}
	scope := normalizeWorkerProfileEnum(input.ConsentScope, WellbeingConsentAggregate)
	if !containsString([]string{WellbeingConsentPrivate, WellbeingConsentAggregate, WellbeingConsentHRAlert}, scope) {
		return nil, ErrInvalidWellbeingScore
	}
	return &WellbeingScore{TenantID: input.TenantID, WorkerProfileID: input.WorkerProfileID, ScoreDate: scoreDate, SourceSurveyID: cleanOKRUUIDPtr(input.SourceSurveyID), WellbeingScore: clampFloat(input.WellbeingScore, 0, 100), MoodScore: cleanScaleScore(input.MoodScore), StressScore: cleanScaleScore(input.StressScore), WorkloadScore: cleanScaleScore(input.WorkloadScore), RiskLevel: risk, ConsentScope: scope, Notes: cleanOKRStringPtr(input.Notes), Metadata: cleanOKRObjectJSON(input.Metadata)}, nil
}

func NewWellbeingAlert(input WellbeingAlertInput) (*WellbeingAlert, error) {
	message := strings.TrimSpace(input.Message)
	if input.TenantID == uuid.Nil || message == "" {
		return nil, ErrInvalidWellbeingAlert
	}
	alertType := normalizeWorkerProfileEnum(input.AlertType, WellbeingAlertCriticalResponse)
	if !containsString([]string{WellbeingAlertCriticalResponse, WellbeingAlertLowScore, WellbeingAlertConsentIssue}, alertType) {
		return nil, ErrInvalidWellbeingAlert
	}
	severity := normalizeWorkerProfileEnum(input.Severity, WellbeingAlertSeverityMedium)
	if !containsString([]string{WellbeingAlertSeverityMedium, WellbeingAlertSeverityHigh, WellbeingAlertSeverityCritical}, severity) {
		return nil, ErrInvalidWellbeingAlert
	}
	status := normalizeWorkerProfileEnum(input.Status, WellbeingAlertStatusOpen)
	if !containsString([]string{WellbeingAlertStatusOpen, WellbeingAlertStatusAcknowledged, WellbeingAlertStatusResolved, WellbeingAlertStatusDismissed}, status) {
		return nil, ErrInvalidWellbeingAlert
	}
	return &WellbeingAlert{TenantID: input.TenantID, WorkerProfileID: cleanOKRUUIDPtr(input.WorkerProfileID), SurveyID: cleanOKRUUIDPtr(input.SurveyID), ResponseID: cleanOKRUUIDPtr(input.ResponseID), AlertType: alertType, Severity: severity, Status: status, VisibleTo: "hr", Message: message, Metadata: cleanOKRObjectJSON(input.Metadata)}, nil
}

func NormalizeWellbeingSearch(value *string) *string {
	if value == nil {
		return nil
	}
	clean := strings.TrimSpace(*value)
	if clean == "" {
		return nil
	}
	return &clean
}

func cleanScaleScore(value *float64) *float64 {
	if value == nil {
		return nil
	}
	clean := clampFloat(*value, 0, 5)
	return &clean
}

func clampFloat(value float64, min float64, max float64) float64 {
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return min
	}
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return math.Round(value*100) / 100
}

func deriveWellbeingRisk(score *float64) string {
	if score == nil {
		return WellbeingRiskNone
	}
	switch {
	case *score <= 1:
		return WellbeingRiskCritical
	case *score <= 2:
		return WellbeingRiskHigh
	case *score <= 3:
		return WellbeingRiskMedium
	default:
		return WellbeingRiskLow
	}
}

func deriveWellbeingScoreRisk(score float64) string {
	switch {
	case score <= 20:
		return WellbeingRiskCritical
	case score <= 40:
		return WellbeingRiskHigh
	case score <= 60:
		return WellbeingRiskMedium
	case score <= 75:
		return WellbeingRiskLow
	default:
		return WellbeingRiskNone
	}
}

func cleanJSONArray(value json.RawMessage) json.RawMessage {
	if len(value) == 0 || !json.Valid(value) {
		return json.RawMessage(`[]`)
	}
	var parsed any
	if err := json.Unmarshal(value, &parsed); err != nil {
		return json.RawMessage(`[]`)
	}
	if _, ok := parsed.([]any); !ok {
		return json.RawMessage(`[]`)
	}
	return value
}
