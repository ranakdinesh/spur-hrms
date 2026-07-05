package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func (s *Store) CreateERCaseCategory(ctx context.Context, item *domain.ERCaseCategory, actorID *uuid.UUID) (*domain.ERCaseCategory, error) {
	row, err := s.getQueries(ctx).CreateERCaseCategory(ctx, sqlc.CreateERCaseCategoryParams{TenantID: item.TenantID, Code: item.Code, Name: item.Name, CaseFamily: item.CaseFamily, Description: textFromPtr(item.Description), DefaultSeverity: item.DefaultSeverity, DefaultOwnerRole: textFromPtr(item.DefaultOwnerRole), IsActive: item.IsActive, ActorID: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create er case category", err, tenantIDField(item.TenantID), stringField("code", item.Code))
	}
	return mapERCaseCategory(row), nil
}

func (s *Store) UpdateERCaseCategory(ctx context.Context, item *domain.ERCaseCategory, actorID *uuid.UUID) (*domain.ERCaseCategory, error) {
	row, err := s.getQueries(ctx).UpdateERCaseCategory(ctx, sqlc.UpdateERCaseCategoryParams{TenantID: item.TenantID, ID: item.ID, Code: item.Code, Name: item.Name, CaseFamily: item.CaseFamily, Description: textFromPtr(item.Description), DefaultSeverity: item.DefaultSeverity, DefaultOwnerRole: textFromPtr(item.DefaultOwnerRole), IsActive: item.IsActive, ActorID: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "update er case category", err, tenantIDField(item.TenantID), stringField("er_case_category_id", item.ID.String()))
	}
	return mapERCaseCategory(row), nil
}

func (s *Store) ListERCaseCategories(ctx context.Context, tenantID uuid.UUID, activeOnly *bool) ([]*domain.ERCaseCategory, error) {
	rows, err := s.getQueries(ctx).ListERCaseCategories(ctx, sqlc.ListERCaseCategoriesParams{TenantID: tenantID, IsActive: boolFromPtr(activeOnly)})
	if err != nil {
		return nil, s.logDBError(ctx, "list er case categories", err, tenantIDField(tenantID))
	}
	return mapERCaseCategories(rows), nil
}

func (s *Store) GetERCaseCategory(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.ERCaseCategory, error) {
	row, err := s.getQueries(ctx).GetERCaseCategory(ctx, sqlc.GetERCaseCategoryParams{TenantID: tenantID, ID: id})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrERCaseCategoryNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "get er case category", err, tenantIDField(tenantID), stringField("er_case_category_id", id.String()))
	}
	return mapERCaseCategory(row), nil
}

func (s *Store) DeleteERCaseCategory(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteERCaseCategory(ctx, sqlc.SoftDeleteERCaseCategoryParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete er case category", err, tenantIDField(tenantID), stringField("er_case_category_id", id.String()))
	}
	return nil
}

func (s *Store) CreateERCase(ctx context.Context, item *domain.ERCase, actorID *uuid.UUID) (*domain.ERCase, error) {
	row, err := s.getQueries(ctx).CreateERCase(ctx, erCaseParams(item, actorID))
	if err != nil {
		return nil, s.logDBError(ctx, "create er case", err, tenantIDField(item.TenantID), stringField("case_number", item.CaseNumber))
	}
	return mapERCase(row), nil
}

func (s *Store) UpdateERCase(ctx context.Context, item *domain.ERCase, actorID *uuid.UUID) (*domain.ERCase, error) {
	row, err := s.getQueries(ctx).UpdateERCase(ctx, sqlc.UpdateERCaseParams{TenantID: item.TenantID, ID: item.ID, Title: item.Title, IntakeSummary: item.IntakeSummary, CaseFamily: item.CaseFamily, Severity: item.Severity, ConfidentialityLevel: item.ConfidentialityLevel, SourceHrCaseID: uuidFromPtr(item.SourceHRCaseID), CategoryID: uuidFromPtr(item.CategoryID), ComplainantUserID: uuidFromPtr(item.ComplainantUserID), SubjectEmployeeUserID: uuidFromPtr(item.SubjectEmployeeUserID), OwnerUserID: uuidFromPtr(item.OwnerUserID), OwnerRole: textFromPtr(item.OwnerRole), InvestigationLeadUserID: uuidFromPtr(item.InvestigationLeadUserID), DueAt: timestamptzFromPtr(item.DueAt), PrivacyNotes: textFromPtr(item.PrivacyNotes), ActorID: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "update er case", err, tenantIDField(item.TenantID), stringField("er_case_id", item.ID.String()))
	}
	return mapERCase(row), nil
}

func (s *Store) UpdateERCaseStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, resolutionSummary *string, actorID *uuid.UUID) (*domain.ERCase, error) {
	row, err := s.getQueries(ctx).UpdateERCaseStatus(ctx, sqlc.UpdateERCaseStatusParams{TenantID: tenantID, ID: id, Status: status, ResolutionSummary: textFromPtr(resolutionSummary), ActorID: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "update er case status", err, tenantIDField(tenantID), stringField("er_case_id", id.String()), stringField("status", status))
	}
	return mapERCase(row), nil
}

func (s *Store) UpdateERCaseLegalHold(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, enabled bool, reason *string, actorID *uuid.UUID) (*domain.ERCase, error) {
	row, err := s.getQueries(ctx).UpdateERCaseLegalHold(ctx, sqlc.UpdateERCaseLegalHoldParams{TenantID: tenantID, ID: id, LegalHold: enabled, LegalHoldReason: textFromPtr(reason), ActorID: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "update er case legal hold", err, tenantIDField(tenantID), stringField("er_case_id", id.String()))
	}
	return mapERCase(row), nil
}

func (s *Store) ListERCases(ctx context.Context, filter domain.ERCaseFilter) ([]*domain.ERCase, error) {
	params := erCaseFilterParams(filter)
	rows, err := s.getQueries(ctx).ListERCases(ctx, params)
	if err != nil {
		return nil, s.logDBError(ctx, "list er cases", err, tenantIDField(filter.TenantID))
	}
	return mapERCases(rows), nil
}

func (s *Store) CountERCases(ctx context.Context, filter domain.ERCaseFilter) (int64, error) {
	params := erCaseFilterParams(filter)
	total, err := s.getQueries(ctx).CountERCases(ctx, sqlc.CountERCasesParams{TenantID: params.TenantID, Status: params.Status, Severity: params.Severity, CaseFamily: params.CaseFamily, CategoryID: params.CategoryID, OwnerUserID: params.OwnerUserID, SubjectEmployeeUserID: params.SubjectEmployeeUserID, ComplainantUserID: params.ComplainantUserID, LegalHold: params.LegalHold, Search: params.Search})
	if err != nil {
		return 0, s.logDBError(ctx, "count er cases", err, tenantIDField(filter.TenantID))
	}
	return total, nil
}

func (s *Store) GetERCase(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.ERCase, error) {
	row, err := s.getQueries(ctx).GetERCase(ctx, sqlc.GetERCaseParams{TenantID: tenantID, ID: id})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrERCaseNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "get er case", err, tenantIDField(tenantID), stringField("er_case_id", id.String()))
	}
	return mapERCaseGetRow(row), nil
}

func (s *Store) CreateERCaseParty(ctx context.Context, item *domain.ERCaseParty, actorID *uuid.UUID) (*domain.ERCaseParty, error) {
	row, err := s.getQueries(ctx).CreateERCaseParty(ctx, sqlc.CreateERCasePartyParams{TenantID: item.TenantID, ErCaseID: item.ERCaseID, PartyUserID: uuidFromPtr(item.PartyUserID), PartyName: textFromPtr(item.PartyName), PartyRole: item.PartyRole, RepresentationNotes: textFromPtr(item.RepresentationNotes), ContactNotes: textFromPtr(item.ContactNotes), ActorID: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create er case party", err, tenantIDField(item.TenantID), stringField("er_case_id", item.ERCaseID.String()))
	}
	return mapERCaseParty(row), nil
}

func (s *Store) ListERCaseParties(ctx context.Context, tenantID uuid.UUID, caseID uuid.UUID) ([]*domain.ERCaseParty, error) {
	rows, err := s.getQueries(ctx).ListERCaseParties(ctx, sqlc.ListERCasePartiesParams{TenantID: tenantID, ErCaseID: caseID})
	if err != nil {
		return nil, s.logDBError(ctx, "list er case parties", err, tenantIDField(tenantID), stringField("er_case_id", caseID.String()))
	}
	return mapERCaseParties(rows), nil
}

func (s *Store) CreateERAllegation(ctx context.Context, item *domain.ERAllegation, actorID *uuid.UUID) (*domain.ERAllegation, error) {
	row, err := s.getQueries(ctx).CreateERAllegation(ctx, sqlc.CreateERAllegationParams{TenantID: item.TenantID, ErCaseID: item.ERCaseID, AllegationType: item.AllegationType, Description: item.Description, Status: item.Status, IncidentDate: dateFromPtr(item.IncidentDate), IncidentLocation: textFromPtr(item.IncidentLocation), PolicyReference: textFromPtr(item.PolicyReference), ActorID: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create er allegation", err, tenantIDField(item.TenantID), stringField("er_case_id", item.ERCaseID.String()))
	}
	return mapERAllegation(row), nil
}

func (s *Store) ListERAllegations(ctx context.Context, tenantID uuid.UUID, caseID uuid.UUID) ([]*domain.ERAllegation, error) {
	rows, err := s.getQueries(ctx).ListERAllegations(ctx, sqlc.ListERAllegationsParams{TenantID: tenantID, ErCaseID: caseID})
	if err != nil {
		return nil, s.logDBError(ctx, "list er allegations", err, tenantIDField(tenantID), stringField("er_case_id", caseID.String()))
	}
	return mapERAllegations(rows), nil
}

func (s *Store) CreateERInvestigationStep(ctx context.Context, item *domain.ERInvestigationStep, actorID *uuid.UUID) (*domain.ERInvestigationStep, error) {
	row, err := s.getQueries(ctx).CreateERInvestigationStep(ctx, sqlc.CreateERInvestigationStepParams{TenantID: item.TenantID, ErCaseID: item.ERCaseID, StepType: item.StepType, Title: item.Title, Status: item.Status, Description: textFromPtr(item.Description), OwnerUserID: uuidFromPtr(item.OwnerUserID), DueAt: timestamptzFromPtr(item.DueAt), CompletedAt: timestamptzFromPtr(item.CompletedAt), OutcomeNotes: textFromPtr(item.OutcomeNotes), ActorID: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create er investigation step", err, tenantIDField(item.TenantID), stringField("er_case_id", item.ERCaseID.String()))
	}
	return mapERInvestigationStep(row), nil
}

func (s *Store) UpdateERInvestigationStepStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, completedAt *time.Time, outcomeNotes *string, actorID *uuid.UUID) (*domain.ERInvestigationStep, error) {
	row, err := s.getQueries(ctx).UpdateERInvestigationStepStatus(ctx, sqlc.UpdateERInvestigationStepStatusParams{TenantID: tenantID, ID: id, Status: status, CompletedAt: timestamptzFromPtr(completedAt), OutcomeNotes: textFromPtr(outcomeNotes), ActorID: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "update er investigation step status", err, tenantIDField(tenantID), stringField("er_step_id", id.String()), stringField("status", status))
	}
	return mapERInvestigationStep(row), nil
}

func (s *Store) ListERInvestigationSteps(ctx context.Context, tenantID uuid.UUID, caseID uuid.UUID) ([]*domain.ERInvestigationStep, error) {
	rows, err := s.getQueries(ctx).ListERInvestigationSteps(ctx, sqlc.ListERInvestigationStepsParams{TenantID: tenantID, ErCaseID: caseID})
	if err != nil {
		return nil, s.logDBError(ctx, "list er investigation steps", err, tenantIDField(tenantID), stringField("er_case_id", caseID.String()))
	}
	return mapERInvestigationSteps(rows), nil
}

func (s *Store) CreateERWitnessNote(ctx context.Context, item *domain.ERWitnessNote, actorID *uuid.UUID) (*domain.ERWitnessNote, error) {
	row, err := s.getQueries(ctx).CreateERWitnessNote(ctx, sqlc.CreateERWitnessNoteParams{TenantID: item.TenantID, ErCaseID: item.ERCaseID, StatementSummary: item.StatementSummary, ConfidentialityLevel: item.ConfidentialityLevel, WitnessUserID: uuidFromPtr(item.WitnessUserID), WitnessName: textFromPtr(item.WitnessName), InterviewAt: timestamptzFromPtr(item.InterviewAt), InterviewerUserID: uuidFromPtr(item.InterviewerUserID), ActorID: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create er witness note", err, tenantIDField(item.TenantID), stringField("er_case_id", item.ERCaseID.String()))
	}
	return mapERWitnessNote(row), nil
}

func (s *Store) ListERWitnessNotes(ctx context.Context, tenantID uuid.UUID, caseID uuid.UUID) ([]*domain.ERWitnessNote, error) {
	rows, err := s.getQueries(ctx).ListERWitnessNotes(ctx, sqlc.ListERWitnessNotesParams{TenantID: tenantID, ErCaseID: caseID})
	if err != nil {
		return nil, s.logDBError(ctx, "list er witness notes", err, tenantIDField(tenantID), stringField("er_case_id", caseID.String()))
	}
	return mapERWitnessNotes(rows), nil
}

func (s *Store) CreateEREvidenceAttachment(ctx context.Context, item *domain.EREvidenceAttachment, actorID *uuid.UUID) (*domain.EREvidenceAttachment, error) {
	row, err := s.getQueries(ctx).CreateEREvidenceAttachment(ctx, sqlc.CreateEREvidenceAttachmentParams{TenantID: item.TenantID, ErCaseID: item.ERCaseID, AllegationID: uuidFromPtr(item.AllegationID), FileName: item.FileName, ContentType: item.ContentType, StoragePath: item.StoragePath, ChecksumSha256: textFromPtr(item.ChecksumSHA), SizeBytes: item.SizeBytes, EvidenceType: item.EvidenceType, Description: textFromPtr(item.Description), UploadedBy: uuidFromPtr(item.UploadedBy), LegalHold: item.LegalHold, ActorID: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create er evidence attachment", err, tenantIDField(item.TenantID), stringField("er_case_id", item.ERCaseID.String()))
	}
	return mapEREvidenceAttachment(row), nil
}

func (s *Store) ListEREvidenceAttachments(ctx context.Context, tenantID uuid.UUID, caseID uuid.UUID) ([]*domain.EREvidenceAttachment, error) {
	rows, err := s.getQueries(ctx).ListEREvidenceAttachments(ctx, sqlc.ListEREvidenceAttachmentsParams{TenantID: tenantID, ErCaseID: caseID})
	if err != nil {
		return nil, s.logDBError(ctx, "list er evidence attachments", err, tenantIDField(tenantID), stringField("er_case_id", caseID.String()))
	}
	return mapEREvidenceAttachments(rows), nil
}

func (s *Store) CreateERFinding(ctx context.Context, item *domain.ERFinding, actorID *uuid.UUID) (*domain.ERFinding, error) {
	decidedBy := item.DecidedBy
	if decidedBy == nil {
		decidedBy = actorID
	}
	row, err := s.getQueries(ctx).CreateERFinding(ctx, sqlc.CreateERFindingParams{TenantID: item.TenantID, ErCaseID: item.ERCaseID, AllegationID: uuidFromPtr(item.AllegationID), Finding: item.Finding, Rationale: item.Rationale, RecommendedAction: textFromPtr(item.RecommendedAction), DecidedBy: uuidFromPtr(decidedBy), DecidedAt: timestamptzFromPtr(item.DecidedAt), ActorID: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create er finding", err, tenantIDField(item.TenantID), stringField("er_case_id", item.ERCaseID.String()))
	}
	return mapERFinding(row), nil
}

func (s *Store) ListERFindings(ctx context.Context, tenantID uuid.UUID, caseID uuid.UUID) ([]*domain.ERFinding, error) {
	rows, err := s.getQueries(ctx).ListERFindings(ctx, sqlc.ListERFindingsParams{TenantID: tenantID, ErCaseID: caseID})
	if err != nil {
		return nil, s.logDBError(ctx, "list er findings", err, tenantIDField(tenantID), stringField("er_case_id", caseID.String()))
	}
	return mapERFindings(rows), nil
}

func (s *Store) CreateERActionPlan(ctx context.Context, item *domain.ERActionPlan, actorID *uuid.UUID) (*domain.ERActionPlan, error) {
	row, err := s.getQueries(ctx).CreateERActionPlan(ctx, sqlc.CreateERActionPlanParams{TenantID: item.TenantID, ErCaseID: item.ERCaseID, ActionType: item.ActionType, Description: item.Description, AssignedToUserID: uuidFromPtr(item.AssignedToUserID), DueAt: timestamptzFromPtr(item.DueAt), CompletedAt: timestamptzFromPtr(item.CompletedAt), Status: item.Status, FollowUpNotes: textFromPtr(item.FollowUpNotes), ActorID: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create er action plan", err, tenantIDField(item.TenantID), stringField("er_case_id", item.ERCaseID.String()))
	}
	return mapERActionPlan(row), nil
}

func (s *Store) UpdateERActionPlanStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, completedAt *time.Time, followUpNotes *string, actorID *uuid.UUID) (*domain.ERActionPlan, error) {
	row, err := s.getQueries(ctx).UpdateERActionPlanStatus(ctx, sqlc.UpdateERActionPlanStatusParams{TenantID: tenantID, ID: id, Status: status, CompletedAt: timestamptzFromPtr(completedAt), FollowUpNotes: textFromPtr(followUpNotes), ActorID: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "update er action plan status", err, tenantIDField(tenantID), stringField("er_action_plan_id", id.String()), stringField("status", status))
	}
	return mapERActionPlan(row), nil
}

func (s *Store) ListERActionPlans(ctx context.Context, tenantID uuid.UUID, caseID uuid.UUID) ([]*domain.ERActionPlan, error) {
	rows, err := s.getQueries(ctx).ListERActionPlans(ctx, sqlc.ListERActionPlansParams{TenantID: tenantID, ErCaseID: caseID})
	if err != nil {
		return nil, s.logDBError(ctx, "list er action plans", err, tenantIDField(tenantID), stringField("er_case_id", caseID.String()))
	}
	return mapERActionPlans(rows), nil
}

func (s *Store) CreateERCaseEvent(ctx context.Context, item *domain.ERCaseEvent, actorID *uuid.UUID) (*domain.ERCaseEvent, error) {
	row, err := s.getQueries(ctx).CreateERCaseEvent(ctx, sqlc.CreateERCaseEventParams{TenantID: item.TenantID, ErCaseID: item.ERCaseID, EventType: item.EventType, FromStatus: textFromPtr(item.FromStatus), ToStatus: textFromPtr(item.ToStatus), ActorUserID: uuidFromPtr(item.ActorUserID), Comment: textFromPtr(item.Comment), Metadata: jsonBytesFromRaw(item.Metadata), ActorID: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create er case event", err, tenantIDField(item.TenantID), stringField("er_case_id", item.ERCaseID.String()), stringField("event_type", item.EventType))
	}
	return mapERCaseEvent(row), nil
}

func (s *Store) ListERCaseEvents(ctx context.Context, tenantID uuid.UUID, caseID uuid.UUID) ([]*domain.ERCaseEvent, error) {
	rows, err := s.getQueries(ctx).ListERCaseEvents(ctx, sqlc.ListERCaseEventsParams{TenantID: tenantID, ErCaseID: caseID})
	if err != nil {
		return nil, s.logDBError(ctx, "list er case events", err, tenantIDField(tenantID), stringField("er_case_id", caseID.String()))
	}
	return mapERCaseEvents(rows), nil
}

func (s *Store) GetERCaseSummary(ctx context.Context, tenantID uuid.UUID) ([]*domain.ERCaseSummaryRow, error) {
	rows, err := s.getQueries(ctx).GetERCaseSummary(ctx, tenantID)
	if err != nil {
		return nil, s.logDBError(ctx, "get er case summary", err, tenantIDField(tenantID))
	}
	return mapERCaseSummary(rows), nil
}

func erCaseParams(item *domain.ERCase, actorID *uuid.UUID) sqlc.CreateERCaseParams {
	return sqlc.CreateERCaseParams{TenantID: item.TenantID, CaseNumber: item.CaseNumber, Title: item.Title, IntakeSummary: item.IntakeSummary, CaseFamily: item.CaseFamily, Severity: item.Severity, Status: item.Status, ConfidentialityLevel: item.ConfidentialityLevel, LegalHold: item.LegalHold, SourceHrCaseID: uuidFromPtr(item.SourceHRCaseID), CategoryID: uuidFromPtr(item.CategoryID), ComplainantUserID: uuidFromPtr(item.ComplainantUserID), SubjectEmployeeUserID: uuidFromPtr(item.SubjectEmployeeUserID), OwnerUserID: uuidFromPtr(item.OwnerUserID), OwnerRole: textFromPtr(item.OwnerRole), InvestigationLeadUserID: uuidFromPtr(item.InvestigationLeadUserID), LegalHoldReason: textFromPtr(item.LegalHoldReason), LegalHoldAt: timestamptzFromPtr(item.LegalHoldAt), LegalHoldBy: uuidFromPtr(item.LegalHoldBy), DueAt: timestamptzFromPtr(item.DueAt), PrivacyNotes: textFromPtr(item.PrivacyNotes), ActorID: uuidFromPtr(actorID)}
}

func erCaseFilterParams(filter domain.ERCaseFilter) sqlc.ListERCasesParams {
	limit := filter.Limit
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	offset := filter.Offset
	if offset < 0 {
		offset = 0
	}
	return sqlc.ListERCasesParams{TenantID: filter.TenantID, Status: textFromPtr(filter.Status), Severity: textFromPtr(filter.Severity), CaseFamily: textFromPtr(filter.CaseFamily), CategoryID: uuidFromPtr(filter.CategoryID), OwnerUserID: uuidFromPtr(filter.OwnerUserID), SubjectEmployeeUserID: uuidFromPtr(filter.SubjectEmployeeUserID), ComplainantUserID: uuidFromPtr(filter.ComplainantUserID), LegalHold: boolFromPtr(filter.LegalHold), Search: textFromPtr(filter.Search), Limit: limit, Offset: offset}
}
