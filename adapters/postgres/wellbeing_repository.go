package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func (s *Store) CreatePulseSurvey(ctx context.Context, item *domain.PulseSurvey, actorID *uuid.UUID) (*domain.PulseSurvey, error) {
	row, err := s.getQueries(ctx).CreatePulseSurvey(ctx, pulseSurveyCreateParams(item, actorID))
	if err != nil {
		return nil, s.logDBError(ctx, "create pulse survey", err, tenantIDField(item.TenantID), stringField("title", item.Title))
	}
	return mapPulseSurvey(row), nil
}

func (s *Store) UpdatePulseSurvey(ctx context.Context, item *domain.PulseSurvey, actorID *uuid.UUID) (*domain.PulseSurvey, error) {
	row, err := s.getQueries(ctx).UpdatePulseSurvey(ctx, pulseSurveyUpdateParams(item, actorID))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrPulseSurveyNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "update pulse survey", err, tenantIDField(item.TenantID), stringField("survey_id", item.ID.String()))
	}
	return mapPulseSurvey(row), nil
}

func (s *Store) UpdatePulseSurveyStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, actorID *uuid.UUID) (*domain.PulseSurvey, error) {
	row, err := s.getQueries(ctx).UpdatePulseSurveyStatus(ctx, sqlc.UpdatePulseSurveyStatusParams{TenantID: tenantID, ID: id, Status: status, UpdatedBy: uuidFromPtr(actorID)})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrPulseSurveyNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "update pulse survey status", err, tenantIDField(tenantID), stringField("survey_id", id.String()))
	}
	return mapPulseSurvey(row), nil
}

func (s *Store) GetPulseSurvey(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.PulseSurvey, error) {
	row, err := s.getQueries(ctx).GetPulseSurvey(ctx, sqlc.GetPulseSurveyParams{TenantID: tenantID, ID: id})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrPulseSurveyNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "get pulse survey", err, tenantIDField(tenantID), stringField("survey_id", id.String()))
	}
	return mapPulseSurveyRow(row), nil
}

func (s *Store) ListPulseSurveys(ctx context.Context, filter domain.PulseSurveyFilter) ([]*domain.PulseSurvey, error) {
	rows, err := s.getQueries(ctx).ListPulseSurveys(ctx, sqlc.ListPulseSurveysParams{TenantID: filter.TenantID, Column2: stringFromPtr(filter.Status), Column3: stringFromPtr(filter.SurveyType), Column4: stringFromPtr(filter.Search)})
	if err != nil {
		return nil, s.logDBError(ctx, "list pulse surveys", err, tenantIDField(filter.TenantID))
	}
	return mapPulseSurveyRows(rows), nil
}

func (s *Store) DeletePulseSurvey(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeletePulseSurvey(ctx, sqlc.SoftDeletePulseSurveyParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete pulse survey", err, tenantIDField(tenantID), stringField("survey_id", id.String()))
	}
	return nil
}

func (s *Store) CreatePulseQuestion(ctx context.Context, item *domain.PulseQuestion, actorID *uuid.UUID) (*domain.PulseQuestion, error) {
	row, err := s.getQueries(ctx).CreatePulseQuestion(ctx, sqlc.CreatePulseQuestionParams{TenantID: item.TenantID, SurveyID: item.SurveyID, QuestionText: item.QuestionText, QuestionType: item.QuestionType, Category: item.Category, IsRequired: item.IsRequired, SortOrder: item.SortOrder, Column8: jsonBytesFromRaw(item.Options), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create pulse question", err, tenantIDField(item.TenantID), stringField("survey_id", item.SurveyID.String()))
	}
	return mapPulseQuestion(row), nil
}

func (s *Store) UpdatePulseQuestion(ctx context.Context, item *domain.PulseQuestion, actorID *uuid.UUID) (*domain.PulseQuestion, error) {
	row, err := s.getQueries(ctx).UpdatePulseQuestion(ctx, sqlc.UpdatePulseQuestionParams{TenantID: item.TenantID, ID: item.ID, QuestionText: item.QuestionText, QuestionType: item.QuestionType, Category: item.Category, IsRequired: item.IsRequired, SortOrder: item.SortOrder, Options: jsonBytesFromRaw(item.Options), UpdatedBy: uuidFromPtr(actorID)})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrPulseQuestionNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "update pulse question", err, tenantIDField(item.TenantID), stringField("question_id", item.ID.String()))
	}
	return mapPulseQuestion(row), nil
}

func (s *Store) GetPulseQuestion(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.PulseQuestion, error) {
	row, err := s.getQueries(ctx).GetPulseQuestion(ctx, sqlc.GetPulseQuestionParams{TenantID: tenantID, ID: id})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrPulseQuestionNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "get pulse question", err, tenantIDField(tenantID), stringField("question_id", id.String()))
	}
	return mapPulseQuestion(row), nil
}

func (s *Store) ListPulseQuestions(ctx context.Context, tenantID uuid.UUID, surveyID *uuid.UUID) ([]*domain.PulseQuestion, error) {
	rows, err := s.getQueries(ctx).ListPulseQuestions(ctx, sqlc.ListPulseQuestionsParams{TenantID: tenantID, Column2: uuidValueFromPtr(surveyID)})
	if err != nil {
		return nil, s.logDBError(ctx, "list pulse questions", err, tenantIDField(tenantID))
	}
	return mapPulseQuestions(rows), nil
}

func (s *Store) DeletePulseQuestion(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeletePulseQuestion(ctx, sqlc.SoftDeletePulseQuestionParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete pulse question", err, tenantIDField(tenantID), stringField("question_id", id.String()))
	}
	return nil
}

func (s *Store) CreatePulseResponse(ctx context.Context, item *domain.PulseResponse, actorID *uuid.UUID) (*domain.PulseResponse, error) {
	row, err := s.getQueries(ctx).CreatePulseResponse(ctx, sqlc.CreatePulseResponseParams{TenantID: item.TenantID, SurveyID: item.SurveyID, QuestionID: item.QuestionID, WorkerProfileID: uuidFromPtr(item.WorkerProfileID), ResponseDate: dateFromPtr(&item.ResponseDate), Score: numericFromFloatPtr(item.Score), TextResponse: textFromPtr(item.TextResponse), BooleanResponse: boolFromPtr(item.BooleanResponse), OptionValue: textFromPtr(item.OptionValue), ConsentGiven: item.ConsentGiven, IsAnonymous: item.IsAnonymous, RiskLevel: item.RiskLevel, CriticalAlert: item.CriticalAlert, Column14: jsonBytesFromRaw(item.Metadata), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create pulse response", err, tenantIDField(item.TenantID), stringField("survey_id", item.SurveyID.String()))
	}
	return mapPulseResponse(row), nil
}

func (s *Store) ListPulseResponses(ctx context.Context, filter domain.PulseResponseFilter) ([]*domain.PulseResponse, error) {
	rows, err := s.getQueries(ctx).ListPulseResponses(ctx, sqlc.ListPulseResponsesParams{TenantID: filter.TenantID, Column2: uuidValueFromPtr(filter.SurveyID), Column3: uuidValueFromPtr(filter.WorkerProfileID), Column4: stringFromPtr(filter.RiskLevel)})
	if err != nil {
		return nil, s.logDBError(ctx, "list pulse responses", err, tenantIDField(filter.TenantID))
	}
	return mapPulseResponseRows(rows), nil
}

func (s *Store) UpsertWellbeingScore(ctx context.Context, item *domain.WellbeingScore, actorID *uuid.UUID) (*domain.WellbeingScore, error) {
	row, err := s.getQueries(ctx).UpsertWellbeingScore(ctx, sqlc.UpsertWellbeingScoreParams{TenantID: item.TenantID, WorkerProfileID: item.WorkerProfileID, ScoreDate: dateFromPtr(&item.ScoreDate), SourceSurveyID: uuidFromPtr(item.SourceSurveyID), WellbeingScore: numericFromFloat(item.WellbeingScore), MoodScore: numericFromFloatPtr(item.MoodScore), StressScore: numericFromFloatPtr(item.StressScore), WorkloadScore: numericFromFloatPtr(item.WorkloadScore), RiskLevel: item.RiskLevel, ConsentScope: item.ConsentScope, Notes: textFromPtr(item.Notes), Column12: jsonBytesFromRaw(item.Metadata), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "upsert wellbeing score", err, tenantIDField(item.TenantID), stringField("worker_profile_id", item.WorkerProfileID.String()))
	}
	return mapWellbeingScore(row), nil
}

func (s *Store) ListWellbeingScores(ctx context.Context, filter domain.WellbeingScoreFilter) ([]*domain.WellbeingScore, error) {
	rows, err := s.getQueries(ctx).ListWellbeingScores(ctx, sqlc.ListWellbeingScoresParams{TenantID: filter.TenantID, Column2: uuidValueFromPtr(filter.WorkerProfileID), Column3: stringFromPtr(filter.RiskLevel)})
	if err != nil {
		return nil, s.logDBError(ctx, "list wellbeing scores", err, tenantIDField(filter.TenantID))
	}
	return mapWellbeingScoreRows(rows), nil
}

func (s *Store) CreateWellbeingAlert(ctx context.Context, item *domain.WellbeingAlert, actorID *uuid.UUID) (*domain.WellbeingAlert, error) {
	row, err := s.getQueries(ctx).CreateWellbeingAlert(ctx, sqlc.CreateWellbeingAlertParams{TenantID: item.TenantID, WorkerProfileID: uuidFromPtr(item.WorkerProfileID), SurveyID: uuidFromPtr(item.SurveyID), ResponseID: uuidFromPtr(item.ResponseID), AlertType: item.AlertType, Severity: item.Severity, Status: item.Status, Message: item.Message, Column9: jsonBytesFromRaw(item.Metadata), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create wellbeing alert", err, tenantIDField(item.TenantID))
	}
	return mapWellbeingAlert(row), nil
}

func (s *Store) UpdateWellbeingAlertStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, resolutionNote *string, actorID *uuid.UUID) (*domain.WellbeingAlert, error) {
	row, err := s.getQueries(ctx).UpdateWellbeingAlertStatus(ctx, sqlc.UpdateWellbeingAlertStatusParams{TenantID: tenantID, ID: id, Status: status, ResolutionNote: textFromPtr(resolutionNote), UpdatedBy: uuidFromPtr(actorID)})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrWellbeingAlertNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "update wellbeing alert status", err, tenantIDField(tenantID), stringField("alert_id", id.String()))
	}
	return mapWellbeingAlert(row), nil
}

func (s *Store) ListWellbeingAlerts(ctx context.Context, filter domain.WellbeingAlertFilter) ([]*domain.WellbeingAlert, error) {
	rows, err := s.getQueries(ctx).ListWellbeingAlerts(ctx, sqlc.ListWellbeingAlertsParams{TenantID: filter.TenantID, Column2: stringFromPtr(filter.Status), Column3: stringFromPtr(filter.Severity)})
	if err != nil {
		return nil, s.logDBError(ctx, "list wellbeing alerts", err, tenantIDField(filter.TenantID))
	}
	return mapWellbeingAlertRows(rows), nil
}

func (s *Store) ListWellbeingAggregateRows(ctx context.Context, tenantID uuid.UUID, surveyID *uuid.UUID) ([]*domain.WellbeingAggregateRow, error) {
	rows, err := s.getQueries(ctx).ListWellbeingAggregateRows(ctx, sqlc.ListWellbeingAggregateRowsParams{TenantID: tenantID, Column2: uuidValueFromPtr(surveyID)})
	if err != nil {
		return nil, s.logDBError(ctx, "list wellbeing aggregates", err, tenantIDField(tenantID))
	}
	return mapWellbeingAggregateRows(rows), nil
}

func pulseSurveyCreateParams(item *domain.PulseSurvey, actorID *uuid.UUID) sqlc.CreatePulseSurveyParams {
	return sqlc.CreatePulseSurveyParams{TenantID: item.TenantID, Title: item.Title, Description: textFromPtr(item.Description), SurveyType: item.SurveyType, Status: item.Status, AudienceScope: item.AudienceScope, DepartmentID: uuidFromPtr(item.DepartmentID), StartDate: dateFromPtr(&item.StartDate), EndDate: dateFromPtr(item.EndDate), Frequency: item.Frequency, AnonymityThreshold: item.AnonymityThreshold, ConsentRequired: item.ConsentRequired, ManagerAggregateOnly: item.ManagerAggregateOnly, CriticalAlertsEnabled: item.CriticalAlertsEnabled, Column15: jsonBytesFromRaw(item.Metadata), CreatedBy: uuidFromPtr(actorID)}
}

func pulseSurveyUpdateParams(item *domain.PulseSurvey, actorID *uuid.UUID) sqlc.UpdatePulseSurveyParams {
	return sqlc.UpdatePulseSurveyParams{TenantID: item.TenantID, ID: item.ID, Title: item.Title, Description: textFromPtr(item.Description), SurveyType: item.SurveyType, Status: item.Status, AudienceScope: item.AudienceScope, DepartmentID: uuidFromPtr(item.DepartmentID), StartDate: dateFromPtr(&item.StartDate), EndDate: dateFromPtr(item.EndDate), Frequency: item.Frequency, AnonymityThreshold: item.AnonymityThreshold, ConsentRequired: item.ConsentRequired, ManagerAggregateOnly: item.ManagerAggregateOnly, CriticalAlertsEnabled: item.CriticalAlertsEnabled, Metadata: jsonBytesFromRaw(item.Metadata), UpdatedBy: uuidFromPtr(actorID)}
}
