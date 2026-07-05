package postgres

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func mapPulseSurvey(row sqlc.HrmsPulseSurvey) *domain.PulseSurvey {
	return pulseSurveyFromParts(row.ID, row.TenantID, row.Title, row.Description, row.SurveyType, row.Status, row.AudienceScope, row.DepartmentID, row.StartDate, row.EndDate, row.Frequency, row.AnonymityThreshold, row.ConsentRequired, row.ManagerAggregateOnly, row.CriticalAlertsEnabled, row.Metadata, row.Inactive, row.CreatedAt, row.CreatedBy, row.UpdatedAt, row.UpdatedBy, nil, 0, 0, 0)
}

func mapPulseSurveyRow(row sqlc.GetPulseSurveyRow) *domain.PulseSurvey {
	return pulseSurveyFromParts(row.ID, row.TenantID, row.Title, row.Description, row.SurveyType, row.Status, row.AudienceScope, row.DepartmentID, row.StartDate, row.EndDate, row.Frequency, row.AnonymityThreshold, row.ConsentRequired, row.ManagerAggregateOnly, row.CriticalAlertsEnabled, row.Metadata, row.Inactive, row.CreatedAt, row.CreatedBy, row.UpdatedAt, row.UpdatedBy, ptrFromText(row.DepartmentName), row.QuestionCount, row.RespondentCount, row.ResponseCount)
}

func mapPulseSurveyRows(rows []sqlc.ListPulseSurveysRow) []*domain.PulseSurvey {
	items := make([]*domain.PulseSurvey, 0, len(rows))
	for _, row := range rows {
		items = append(items, pulseSurveyFromParts(row.ID, row.TenantID, row.Title, row.Description, row.SurveyType, row.Status, row.AudienceScope, row.DepartmentID, row.StartDate, row.EndDate, row.Frequency, row.AnonymityThreshold, row.ConsentRequired, row.ManagerAggregateOnly, row.CriticalAlertsEnabled, row.Metadata, row.Inactive, row.CreatedAt, row.CreatedBy, row.UpdatedAt, row.UpdatedBy, ptrFromText(row.DepartmentName), row.QuestionCount, row.RespondentCount, row.ResponseCount))
	}
	return items
}

func pulseSurveyFromParts(id uuid.UUID, tenantID uuid.UUID, title string, description pgtype.Text, surveyType string, status string, audienceScope string, departmentID pgtype.UUID, startDate pgtype.Date, endDate pgtype.Date, frequency string, threshold int32, consentRequired bool, managerAggregateOnly bool, criticalAlertsEnabled bool, metadata []byte, inactive bool, createdAt pgtype.Timestamptz, createdBy pgtype.UUID, updatedAt pgtype.Timestamptz, updatedBy pgtype.UUID, departmentName *string, questionCount int32, respondentCount int32, responseCount int32) *domain.PulseSurvey {
	return &domain.PulseSurvey{ID: id, TenantID: tenantID, Title: title, Description: ptrFromText(description), SurveyType: surveyType, Status: status, AudienceScope: audienceScope, DepartmentID: ptrFromUUID(departmentID), StartDate: timeFromDate(startDate), EndDate: ptrFromDate(endDate), Frequency: frequency, AnonymityThreshold: threshold, ConsentRequired: consentRequired, ManagerAggregateOnly: managerAggregateOnly, CriticalAlertsEnabled: criticalAlertsEnabled, Metadata: okrRaw(metadata), Inactive: inactive, CreatedAt: timeFromTimestamptz(createdAt), CreatedBy: ptrFromUUID(createdBy), UpdatedAt: timeFromTimestamptz(updatedAt), UpdatedBy: ptrFromUUID(updatedBy), DepartmentName: departmentName, QuestionCount: questionCount, RespondentCount: respondentCount, ResponseCount: responseCount}
}

func mapPulseQuestion(row sqlc.HrmsPulseSurveyQuestion) *domain.PulseQuestion {
	return &domain.PulseQuestion{ID: row.ID, TenantID: row.TenantID, SurveyID: row.SurveyID, QuestionText: row.QuestionText, QuestionType: row.QuestionType, Category: row.Category, IsRequired: row.IsRequired, SortOrder: row.SortOrder, Options: jsonRawDefault(row.Options, `[]`), Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy)}
}

func mapPulseQuestions(rows []sqlc.HrmsPulseSurveyQuestion) []*domain.PulseQuestion {
	items := make([]*domain.PulseQuestion, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapPulseQuestion(row))
	}
	return items
}

func mapPulseResponse(row sqlc.HrmsPulseResponse) *domain.PulseResponse {
	return pulseResponseFromParts(row.ID, row.TenantID, row.SurveyID, row.QuestionID, row.WorkerProfileID, row.ResponseDate, row.Score, row.TextResponse, row.BooleanResponse, row.OptionValue, row.ConsentGiven, row.IsAnonymous, row.RiskLevel, row.CriticalAlert, row.Metadata, row.Inactive, row.CreatedAt, row.CreatedBy, nil, nil, nil, nil, nil)
}

func mapPulseResponseRows(rows []sqlc.ListPulseResponsesRow) []*domain.PulseResponse {
	items := make([]*domain.PulseResponse, 0, len(rows))
	for _, row := range rows {
		workerName := ptrFromText(row.WorkerDisplayName)
		workerCode := ptrFromText(row.WorkerCode)
		if row.IsAnonymous && !row.CriticalAlert {
			workerName = nil
			workerCode = nil
		}
		items = append(items, pulseResponseFromParts(row.ID, row.TenantID, row.SurveyID, row.QuestionID, row.WorkerProfileID, row.ResponseDate, row.Score, row.TextResponse, row.BooleanResponse, row.OptionValue, row.ConsentGiven, row.IsAnonymous, row.RiskLevel, row.CriticalAlert, row.Metadata, row.Inactive, row.CreatedAt, row.CreatedBy, &row.SurveyTitle, &row.QuestionText, &row.Category, workerName, workerCode))
	}
	return items
}

func pulseResponseFromParts(id uuid.UUID, tenantID uuid.UUID, surveyID uuid.UUID, questionID uuid.UUID, workerProfileID pgtype.UUID, responseDate pgtype.Date, score pgtype.Numeric, textResponse pgtype.Text, booleanResponse pgtype.Bool, optionValue pgtype.Text, consentGiven bool, isAnonymous bool, riskLevel string, criticalAlert bool, metadata []byte, inactive bool, createdAt pgtype.Timestamptz, createdBy pgtype.UUID, surveyTitle *string, questionText *string, category *string, workerDisplayName *string, workerCode *string) *domain.PulseResponse {
	return &domain.PulseResponse{ID: id, TenantID: tenantID, SurveyID: surveyID, QuestionID: questionID, WorkerProfileID: ptrFromUUID(workerProfileID), ResponseDate: timeFromDate(responseDate), Score: floatPtrFromNumeric(score), TextResponse: ptrFromText(textResponse), BooleanResponse: ptrFromBool(booleanResponse), OptionValue: ptrFromText(optionValue), ConsentGiven: consentGiven, IsAnonymous: isAnonymous, RiskLevel: riskLevel, CriticalAlert: criticalAlert, Metadata: okrRaw(metadata), Inactive: inactive, CreatedAt: timeFromTimestamptz(createdAt), CreatedBy: ptrFromUUID(createdBy), SurveyTitle: surveyTitle, QuestionText: questionText, Category: category, WorkerDisplayName: workerDisplayName, WorkerCode: workerCode}
}

func mapWellbeingScore(row sqlc.HrmsWellbeingScore) *domain.WellbeingScore {
	return wellbeingScoreFromParts(row.ID, row.TenantID, row.WorkerProfileID, row.ScoreDate, row.SourceSurveyID, row.WellbeingScore, row.MoodScore, row.StressScore, row.WorkloadScore, row.RiskLevel, row.ConsentScope, row.Notes, row.Metadata, row.Inactive, row.CreatedAt, row.CreatedBy, row.UpdatedAt, row.UpdatedBy, nil, nil, nil)
}

func mapWellbeingScoreRows(rows []sqlc.ListWellbeingScoresRow) []*domain.WellbeingScore {
	items := make([]*domain.WellbeingScore, 0, len(rows))
	for _, row := range rows {
		items = append(items, wellbeingScoreFromParts(row.ID, row.TenantID, row.WorkerProfileID, row.ScoreDate, row.SourceSurveyID, row.WellbeingScore, row.MoodScore, row.StressScore, row.WorkloadScore, row.RiskLevel, row.ConsentScope, row.Notes, row.Metadata, row.Inactive, row.CreatedAt, row.CreatedBy, row.UpdatedAt, row.UpdatedBy, &row.WorkerDisplayName, ptrFromText(row.WorkerCode), ptrFromText(row.SurveyTitle)))
	}
	return items
}

func wellbeingScoreFromParts(id uuid.UUID, tenantID uuid.UUID, workerProfileID uuid.UUID, scoreDate pgtype.Date, sourceSurveyID pgtype.UUID, score pgtype.Numeric, mood pgtype.Numeric, stress pgtype.Numeric, workload pgtype.Numeric, riskLevel string, consentScope string, notes pgtype.Text, metadata []byte, inactive bool, createdAt pgtype.Timestamptz, createdBy pgtype.UUID, updatedAt pgtype.Timestamptz, updatedBy pgtype.UUID, workerDisplayName *string, workerCode *string, surveyTitle *string) *domain.WellbeingScore {
	return &domain.WellbeingScore{ID: id, TenantID: tenantID, WorkerProfileID: workerProfileID, ScoreDate: timeFromDate(scoreDate), SourceSurveyID: ptrFromUUID(sourceSurveyID), WellbeingScore: floatFromNumeric(score), MoodScore: floatPtrFromNumeric(mood), StressScore: floatPtrFromNumeric(stress), WorkloadScore: floatPtrFromNumeric(workload), RiskLevel: riskLevel, ConsentScope: consentScope, Notes: ptrFromText(notes), Metadata: okrRaw(metadata), Inactive: inactive, CreatedAt: timeFromTimestamptz(createdAt), CreatedBy: ptrFromUUID(createdBy), UpdatedAt: timeFromTimestamptz(updatedAt), UpdatedBy: ptrFromUUID(updatedBy), WorkerDisplayName: workerDisplayName, WorkerCode: workerCode, SurveyTitle: surveyTitle}
}

func mapWellbeingAlert(row sqlc.HrmsWellbeingAlert) *domain.WellbeingAlert {
	return wellbeingAlertFromParts(row.ID, row.TenantID, row.WorkerProfileID, row.SurveyID, row.ResponseID, row.AlertType, row.Severity, row.Status, row.VisibleTo, row.Message, row.ResolutionNote, row.Metadata, row.Inactive, row.CreatedAt, row.CreatedBy, row.UpdatedAt, row.UpdatedBy, nil, nil, nil)
}

func mapWellbeingAlertRows(rows []sqlc.ListWellbeingAlertsRow) []*domain.WellbeingAlert {
	items := make([]*domain.WellbeingAlert, 0, len(rows))
	for _, row := range rows {
		items = append(items, wellbeingAlertFromParts(row.ID, row.TenantID, row.WorkerProfileID, row.SurveyID, row.ResponseID, row.AlertType, row.Severity, row.Status, row.VisibleTo, row.Message, row.ResolutionNote, row.Metadata, row.Inactive, row.CreatedAt, row.CreatedBy, row.UpdatedAt, row.UpdatedBy, ptrFromText(row.WorkerDisplayName), ptrFromText(row.WorkerCode), ptrFromText(row.SurveyTitle)))
	}
	return items
}

func wellbeingAlertFromParts(id uuid.UUID, tenantID uuid.UUID, workerProfileID pgtype.UUID, surveyID pgtype.UUID, responseID pgtype.UUID, alertType string, severity string, status string, visibleTo string, message string, resolutionNote pgtype.Text, metadata []byte, inactive bool, createdAt pgtype.Timestamptz, createdBy pgtype.UUID, updatedAt pgtype.Timestamptz, updatedBy pgtype.UUID, workerDisplayName *string, workerCode *string, surveyTitle *string) *domain.WellbeingAlert {
	return &domain.WellbeingAlert{ID: id, TenantID: tenantID, WorkerProfileID: ptrFromUUID(workerProfileID), SurveyID: ptrFromUUID(surveyID), ResponseID: ptrFromUUID(responseID), AlertType: alertType, Severity: severity, Status: status, VisibleTo: visibleTo, Message: message, ResolutionNote: ptrFromText(resolutionNote), Metadata: okrRaw(metadata), Inactive: inactive, CreatedAt: timeFromTimestamptz(createdAt), CreatedBy: ptrFromUUID(createdBy), UpdatedAt: timeFromTimestamptz(updatedAt), UpdatedBy: ptrFromUUID(updatedBy), WorkerDisplayName: workerDisplayName, WorkerCode: workerCode, SurveyTitle: surveyTitle}
}

func mapWellbeingAggregateRows(rows []sqlc.ListWellbeingAggregateRowsRow) []*domain.WellbeingAggregateRow {
	items := make([]*domain.WellbeingAggregateRow, 0, len(rows))
	for _, row := range rows {
		items = append(items, &domain.WellbeingAggregateRow{SurveyID: row.SurveyID, SurveyTitle: row.SurveyTitle, DepartmentID: row.DepartmentID, DepartmentName: row.DepartmentName, Category: row.Category, ResponseCount: row.ResponseCount, RespondentCount: row.RespondentCount, Suppressed: row.Suppressed, AverageScore: floatPtrFromNumeric(row.AverageScore), RiskCount: row.RiskCount, AnonymityThreshold: row.AnonymityThreshold})
	}
	return items
}

func ptrFromBool(value pgtype.Bool) *bool {
	if !value.Valid {
		return nil
	}
	return &value.Bool
}

func jsonRawDefault(value []byte, fallback string) json.RawMessage {
	if len(value) == 0 {
		return json.RawMessage(fallback)
	}
	return json.RawMessage(value)
}
