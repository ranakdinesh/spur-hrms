package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) CreatePulseSurvey(ctx context.Context, cmd ports.PulseSurveyCommand) (*domain.PulseSurvey, error) {
	item, err := s.preparePulseSurvey(ctx, cmd)
	if err != nil {
		s.logError("validate pulse survey", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	return s.wellbeing.CreatePulseSurvey(ctx, item, cmd.ActorID)
}

func (s *TenantService) UpdatePulseSurvey(ctx context.Context, cmd ports.PulseSurveyCommand) (*domain.PulseSurvey, error) {
	if cmd.ID == uuid.Nil {
		return nil, domain.ErrInvalidPulseSurvey
	}
	item, err := s.preparePulseSurvey(ctx, cmd)
	if err != nil {
		s.logError("validate pulse survey update", err, serviceTenantIDField(cmd.TenantID), serviceStringField("survey_id", cmd.ID.String()))
		return nil, err
	}
	item.ID = cmd.ID
	return s.wellbeing.UpdatePulseSurvey(ctx, item, cmd.ActorID)
}

func (s *TenantService) UpdatePulseSurveyStatus(ctx context.Context, cmd ports.PulseSurveyStatusCommand) (*domain.PulseSurvey, error) {
	if cmd.TenantID == uuid.Nil || cmd.ID == uuid.Nil || !stringIn(cmd.Status, []string{domain.PulseSurveyStatusDraft, domain.PulseSurveyStatusActive, domain.PulseSurveyStatusClosed, domain.PulseSurveyStatusArchived}) {
		return nil, domain.ErrInvalidPulseSurvey
	}
	return s.wellbeing.UpdatePulseSurveyStatus(ctx, cmd.TenantID, cmd.ID, cmd.Status, cmd.ActorID)
}

func (s *TenantService) GetPulseSurvey(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.PulseSurvey, error) {
	if tenantID == uuid.Nil || id == uuid.Nil {
		return nil, domain.ErrInvalidPulseSurvey
	}
	return s.wellbeing.GetPulseSurvey(ctx, tenantID, id)
}

func (s *TenantService) ListPulseSurveys(ctx context.Context, filter domain.PulseSurveyFilter) ([]*domain.PulseSurvey, error) {
	if filter.TenantID == uuid.Nil {
		return nil, domain.ErrInvalidTenantID
	}
	filter.Status = domain.NormalizeWellbeingSearch(filter.Status)
	filter.SurveyType = domain.NormalizeWellbeingSearch(filter.SurveyType)
	filter.Search = domain.NormalizeWellbeingSearch(filter.Search)
	return s.wellbeing.ListPulseSurveys(ctx, filter)
}

func (s *TenantService) DeletePulseSurvey(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if tenantID == uuid.Nil || id == uuid.Nil {
		return domain.ErrInvalidPulseSurvey
	}
	return s.wellbeing.DeletePulseSurvey(ctx, tenantID, id, actorID)
}

func (s *TenantService) CreatePulseQuestion(ctx context.Context, cmd ports.PulseQuestionCommand) (*domain.PulseQuestion, error) {
	item, err := s.preparePulseQuestion(ctx, cmd)
	if err != nil {
		s.logError("validate pulse question", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	return s.wellbeing.CreatePulseQuestion(ctx, item, cmd.ActorID)
}

func (s *TenantService) UpdatePulseQuestion(ctx context.Context, cmd ports.PulseQuestionCommand) (*domain.PulseQuestion, error) {
	if cmd.ID == uuid.Nil {
		return nil, domain.ErrInvalidPulseQuestion
	}
	item, err := s.preparePulseQuestion(ctx, cmd)
	if err != nil {
		s.logError("validate pulse question update", err, serviceTenantIDField(cmd.TenantID), serviceStringField("question_id", cmd.ID.String()))
		return nil, err
	}
	item.ID = cmd.ID
	return s.wellbeing.UpdatePulseQuestion(ctx, item, cmd.ActorID)
}

func (s *TenantService) GetPulseQuestion(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.PulseQuestion, error) {
	if tenantID == uuid.Nil || id == uuid.Nil {
		return nil, domain.ErrInvalidPulseQuestion
	}
	return s.wellbeing.GetPulseQuestion(ctx, tenantID, id)
}

func (s *TenantService) ListPulseQuestions(ctx context.Context, tenantID uuid.UUID, surveyID *uuid.UUID) ([]*domain.PulseQuestion, error) {
	if tenantID == uuid.Nil {
		return nil, domain.ErrInvalidTenantID
	}
	return s.wellbeing.ListPulseQuestions(ctx, tenantID, surveyID)
}

func (s *TenantService) DeletePulseQuestion(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if tenantID == uuid.Nil || id == uuid.Nil {
		return domain.ErrInvalidPulseQuestion
	}
	return s.wellbeing.DeletePulseQuestion(ctx, tenantID, id, actorID)
}

func (s *TenantService) CreatePulseResponse(ctx context.Context, cmd ports.PulseResponseCommand) (*domain.PulseResponse, error) {
	survey, err := s.GetPulseSurvey(ctx, cmd.TenantID, cmd.SurveyID)
	if err != nil {
		return nil, err
	}
	question, err := s.GetPulseQuestion(ctx, cmd.TenantID, cmd.QuestionID)
	if err != nil {
		return nil, err
	}
	if cmd.WorkerProfileID != nil {
		if _, err := s.GetWorkerProfile(ctx, cmd.TenantID, *cmd.WorkerProfileID); err != nil {
			return nil, err
		}
	}
	if survey.ConsentRequired && !cmd.ConsentGiven {
		return nil, domain.ErrInvalidPulseResponse
	}
	responseDate, err := parseWorkerProfileDate(cmd.ResponseDate)
	if err != nil {
		return nil, err
	}
	item, err := domain.NewPulseResponse(domain.PulseResponseInput{TenantID: cmd.TenantID, SurveyID: cmd.SurveyID, QuestionID: cmd.QuestionID, WorkerProfileID: cmd.WorkerProfileID, ResponseDate: responseDate, Score: cmd.Score, TextResponse: cmd.TextResponse, BooleanResponse: cmd.BooleanResponse, OptionValue: cmd.OptionValue, ConsentGiven: cmd.ConsentGiven, IsAnonymous: cmd.IsAnonymous, RiskLevel: cmd.RiskLevel, CriticalAlert: cmd.CriticalAlert, Metadata: cmd.Metadata})
	if err != nil {
		s.logError("validate pulse response", err, serviceTenantIDField(cmd.TenantID), serviceStringField("survey_id", cmd.SurveyID.String()))
		return nil, err
	}
	result, err := s.wellbeing.CreatePulseResponse(ctx, item, cmd.ActorID)
	if err != nil {
		return nil, err
	}
	s.captureWellbeingSignals(ctx, survey, question, result, cmd.ActorID)
	return result, nil
}

func (s *TenantService) ListPulseResponses(ctx context.Context, filter domain.PulseResponseFilter) ([]*domain.PulseResponse, error) {
	if filter.TenantID == uuid.Nil {
		return nil, domain.ErrInvalidTenantID
	}
	filter.RiskLevel = domain.NormalizeWellbeingSearch(filter.RiskLevel)
	return s.wellbeing.ListPulseResponses(ctx, filter)
}

func (s *TenantService) UpsertWellbeingScore(ctx context.Context, cmd ports.WellbeingScoreCommand) (*domain.WellbeingScore, error) {
	if _, err := s.GetWorkerProfile(ctx, cmd.TenantID, cmd.WorkerProfileID); err != nil {
		return nil, err
	}
	if cmd.SourceSurveyID != nil {
		if _, err := s.GetPulseSurvey(ctx, cmd.TenantID, *cmd.SourceSurveyID); err != nil {
			return nil, err
		}
	}
	scoreDate, err := parseWorkerProfileDate(cmd.ScoreDate)
	if err != nil {
		return nil, err
	}
	item, err := domain.NewWellbeingScore(domain.WellbeingScoreInput{TenantID: cmd.TenantID, WorkerProfileID: cmd.WorkerProfileID, ScoreDate: scoreDate, SourceSurveyID: cmd.SourceSurveyID, WellbeingScore: cmd.WellbeingScore, MoodScore: cmd.MoodScore, StressScore: cmd.StressScore, WorkloadScore: cmd.WorkloadScore, RiskLevel: cmd.RiskLevel, ConsentScope: cmd.ConsentScope, Notes: cmd.Notes, Metadata: cmd.Metadata})
	if err != nil {
		s.logError("validate wellbeing score", err, serviceTenantIDField(cmd.TenantID), serviceStringField("worker_profile_id", cmd.WorkerProfileID.String()))
		return nil, err
	}
	return s.wellbeing.UpsertWellbeingScore(ctx, item, cmd.ActorID)
}

func (s *TenantService) ListWellbeingScores(ctx context.Context, filter domain.WellbeingScoreFilter) ([]*domain.WellbeingScore, error) {
	if filter.TenantID == uuid.Nil {
		return nil, domain.ErrInvalidTenantID
	}
	filter.RiskLevel = domain.NormalizeWellbeingSearch(filter.RiskLevel)
	return s.wellbeing.ListWellbeingScores(ctx, filter)
}

func (s *TenantService) UpdateWellbeingAlertStatus(ctx context.Context, cmd ports.WellbeingAlertStatusCommand) (*domain.WellbeingAlert, error) {
	if cmd.TenantID == uuid.Nil || cmd.ID == uuid.Nil || !stringIn(cmd.Status, []string{domain.WellbeingAlertStatusOpen, domain.WellbeingAlertStatusAcknowledged, domain.WellbeingAlertStatusResolved, domain.WellbeingAlertStatusDismissed}) {
		return nil, domain.ErrInvalidWellbeingAlert
	}
	return s.wellbeing.UpdateWellbeingAlertStatus(ctx, cmd.TenantID, cmd.ID, cmd.Status, cmd.ResolutionNote, cmd.ActorID)
}

func (s *TenantService) ListWellbeingAlerts(ctx context.Context, filter domain.WellbeingAlertFilter) ([]*domain.WellbeingAlert, error) {
	if filter.TenantID == uuid.Nil {
		return nil, domain.ErrInvalidTenantID
	}
	filter.Status = domain.NormalizeWellbeingSearch(filter.Status)
	filter.Severity = domain.NormalizeWellbeingSearch(filter.Severity)
	return s.wellbeing.ListWellbeingAlerts(ctx, filter)
}

func (s *TenantService) ListWellbeingAggregateRows(ctx context.Context, tenantID uuid.UUID, surveyID *uuid.UUID) ([]*domain.WellbeingAggregateRow, error) {
	if tenantID == uuid.Nil {
		return nil, domain.ErrInvalidTenantID
	}
	return s.wellbeing.ListWellbeingAggregateRows(ctx, tenantID, surveyID)
}

func (s *TenantService) preparePulseSurvey(ctx context.Context, cmd ports.PulseSurveyCommand) (*domain.PulseSurvey, error) {
	if cmd.DepartmentID != nil {
		if _, err := s.GetDepartment(ctx, cmd.TenantID, *cmd.DepartmentID); err != nil {
			return nil, err
		}
	}
	startDate, err := parseWorkerProfileDate(cmd.StartDate)
	if err != nil {
		return nil, err
	}
	endDate, err := parseWorkerProfileDate(cmd.EndDate)
	if err != nil {
		return nil, err
	}
	return domain.NewPulseSurvey(domain.PulseSurveyInput{TenantID: cmd.TenantID, Title: cmd.Title, Description: cmd.Description, SurveyType: cmd.SurveyType, Status: cmd.Status, AudienceScope: cmd.AudienceScope, DepartmentID: cmd.DepartmentID, StartDate: startDate, EndDate: endDate, Frequency: cmd.Frequency, AnonymityThreshold: cmd.AnonymityThreshold, ConsentRequired: cmd.ConsentRequired, ManagerAggregateOnly: cmd.ManagerAggregateOnly, CriticalAlertsEnabled: cmd.CriticalAlertsEnabled, Metadata: cmd.Metadata})
}

func (s *TenantService) preparePulseQuestion(ctx context.Context, cmd ports.PulseQuestionCommand) (*domain.PulseQuestion, error) {
	if _, err := s.GetPulseSurvey(ctx, cmd.TenantID, cmd.SurveyID); err != nil {
		return nil, err
	}
	return domain.NewPulseQuestion(domain.PulseQuestionInput{TenantID: cmd.TenantID, SurveyID: cmd.SurveyID, QuestionText: cmd.QuestionText, QuestionType: cmd.QuestionType, Category: cmd.Category, IsRequired: cmd.IsRequired, SortOrder: cmd.SortOrder, Options: cmd.Options})
}

func (s *TenantService) captureWellbeingSignals(ctx context.Context, survey *domain.PulseSurvey, question *domain.PulseQuestion, response *domain.PulseResponse, actorID *uuid.UUID) {
	if response.WorkerProfileID != nil && response.Score != nil {
		wellbeingScore := (*response.Score / 5) * 100
		var moodScore, stressScore, workloadScore *float64
		switch question.Category {
		case domain.WellbeingCategoryMood:
			moodScore = response.Score
		case domain.WellbeingCategoryStress:
			stressScore = response.Score
		case domain.WellbeingCategoryWorkload:
			workloadScore = response.Score
		}
		score, err := domain.NewWellbeingScore(domain.WellbeingScoreInput{TenantID: response.TenantID, WorkerProfileID: *response.WorkerProfileID, ScoreDate: &response.ResponseDate, SourceSurveyID: &response.SurveyID, WellbeingScore: wellbeingScore, MoodScore: moodScore, StressScore: stressScore, WorkloadScore: workloadScore, RiskLevel: response.RiskLevel, ConsentScope: domain.WellbeingConsentAggregate, Notes: response.TextResponse, Metadata: response.Metadata})
		if err == nil {
			_, _ = s.wellbeing.UpsertWellbeingScore(ctx, score, actorID)
		}
	}
	if !survey.CriticalAlertsEnabled || !response.CriticalAlert {
		return
	}
	severity := domain.WellbeingAlertSeverityHigh
	if response.RiskLevel == domain.WellbeingRiskCritical {
		severity = domain.WellbeingAlertSeverityCritical
	}
	alert, err := domain.NewWellbeingAlert(domain.WellbeingAlertInput{TenantID: response.TenantID, WorkerProfileID: response.WorkerProfileID, SurveyID: &response.SurveyID, ResponseID: &response.ID, AlertType: domain.WellbeingAlertCriticalResponse, Severity: severity, Status: domain.WellbeingAlertStatusOpen, Message: fmt.Sprintf("Critical wellbeing signal recorded for %s", survey.Title), Metadata: response.Metadata})
	if err != nil {
		s.logError("validate wellbeing alert", err, serviceTenantIDField(response.TenantID))
		return
	}
	if _, err := s.wellbeing.CreateWellbeingAlert(ctx, alert, actorID); err != nil {
		s.logError("create wellbeing alert", err, serviceTenantIDField(response.TenantID))
	}
}
