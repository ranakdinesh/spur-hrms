package services

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) CreateERCaseCategory(ctx context.Context, cmd ports.ERCaseCategoryCommand) (*domain.ERCaseCategory, error) {
	item, err := erCaseCategoryFromCommand(cmd)
	if err != nil {
		s.logError("validate er case category", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	result, err := s.employeeRelations.CreateERCaseCategory(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create er case category", err, serviceTenantIDField(cmd.TenantID), serviceStringField("code", item.Code))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) UpdateERCaseCategory(ctx context.Context, cmd ports.ERCaseCategoryCommand) (*domain.ERCaseCategory, error) {
	if cmd.ID == uuid.Nil {
		return nil, domain.ErrInvalidERCaseCategory
	}
	item, err := erCaseCategoryFromCommand(cmd)
	if err != nil {
		s.logError("validate er case category update", err, serviceTenantIDField(cmd.TenantID), serviceStringField("er_case_category_id", cmd.ID.String()))
		return nil, err
	}
	item.ID = cmd.ID
	result, err := s.employeeRelations.UpdateERCaseCategory(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("update er case category", err, serviceTenantIDField(cmd.TenantID), serviceStringField("er_case_category_id", cmd.ID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) ListERCaseCategories(ctx context.Context, tenantID uuid.UUID, activeOnly *bool) ([]*domain.ERCaseCategory, error) {
	items, err := s.employeeRelations.ListERCaseCategories(ctx, tenantID, activeOnly)
	if err != nil {
		s.logError("list er case categories", err, serviceTenantIDField(tenantID))
	}
	return items, err
}

func (s *TenantService) DeleteERCaseCategory(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if tenantID == uuid.Nil || id == uuid.Nil {
		return domain.ErrInvalidERCaseCategory
	}
	err := s.employeeRelations.DeleteERCaseCategory(ctx, tenantID, id, actorID)
	if err != nil {
		s.logError("delete er case category", err, serviceTenantIDField(tenantID), serviceStringField("er_case_category_id", id.String()))
	}
	return err
}

func (s *TenantService) CreateERCase(ctx context.Context, cmd ports.ERCaseCommand) (*domain.ERCase, error) {
	item, err := erCaseFromCommand(cmd)
	if err != nil {
		s.logError("validate er case", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	now := time.Now().UTC()
	item.CaseNumber = domain.NewERCaseNumber(now)
	if item.ComplainantUserID == nil && cmd.ActorID != nil {
		item.ComplainantUserID = cmd.ActorID
	}
	if item.CategoryID != nil {
		category, err := s.employeeRelations.GetERCaseCategory(ctx, item.TenantID, *item.CategoryID)
		if err != nil && !errors.Is(err, domain.ErrERCaseCategoryNotFound) {
			s.logError("get er case category for case create", err, serviceTenantIDField(item.TenantID), serviceStringField("er_case_category_id", item.CategoryID.String()))
			return nil, err
		}
		if category != nil {
			item.CaseFamily = category.CaseFamily
			if cmd.Severity == "" {
				item.Severity = category.DefaultSeverity
			}
			if item.OwnerRole == nil {
				item.OwnerRole = category.DefaultOwnerRole
			}
		}
	}
	result, err := s.employeeRelations.CreateERCase(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create er case", err, serviceTenantIDField(item.TenantID), serviceStringField("title", item.Title))
		return nil, err
	}
	s.recordERCaseEvent(ctx, result, "created", nil, &result.Status, nil, cmd.ActorID, json.RawMessage(`{"source":"employee_relations"}`))
	return result, nil
}

func (s *TenantService) UpdateERCase(ctx context.Context, cmd ports.ERCaseCommand) (*domain.ERCase, error) {
	if cmd.ID == uuid.Nil {
		return nil, domain.ErrInvalidERCase
	}
	item, err := erCaseFromCommand(cmd)
	if err != nil {
		s.logError("validate er case update", err, serviceTenantIDField(cmd.TenantID), serviceStringField("er_case_id", cmd.ID.String()))
		return nil, err
	}
	item.ID = cmd.ID
	result, err := s.employeeRelations.UpdateERCase(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("update er case", err, serviceTenantIDField(cmd.TenantID), serviceStringField("er_case_id", cmd.ID.String()))
		return nil, err
	}
	s.recordERCaseEvent(ctx, result, "updated", nil, nil, nil, cmd.ActorID, nil)
	return result, nil
}

func (s *TenantService) UpdateERCaseStatus(ctx context.Context, cmd ports.ERCaseStatusCommand) (*domain.ERCase, error) {
	status := cleanERCaseStatus(cmd.Status)
	if cmd.TenantID == uuid.Nil || cmd.ID == uuid.Nil || status == "" {
		return nil, domain.ErrInvalidERCase
	}
	before, _ := s.employeeRelations.GetERCase(ctx, cmd.TenantID, cmd.ID)
	if before != nil && !erStatusTransitionAllowed(before.Status, status) {
		s.logError("validate er case status transition", domain.ErrInvalidERCaseTransition, serviceTenantIDField(cmd.TenantID), serviceStringField("er_case_id", cmd.ID.String()), serviceStringField("from_status", before.Status), serviceStringField("to_status", status))
		return nil, domain.ErrInvalidERCaseTransition
	}
	result, err := s.employeeRelations.UpdateERCaseStatus(ctx, cmd.TenantID, cmd.ID, status, cmd.ResolutionSummary, cmd.ActorID)
	if err != nil {
		s.logError("update er case status", err, serviceTenantIDField(cmd.TenantID), serviceStringField("er_case_id", cmd.ID.String()), serviceStringField("status", status))
		return nil, err
	}
	var from *string
	if before != nil {
		from = &before.Status
	}
	s.recordERCaseEvent(ctx, result, "status_changed", from, &result.Status, cmd.Comment, cmd.ActorID, nil)
	return result, nil
}

func (s *TenantService) UpdateERCaseLegalHold(ctx context.Context, cmd ports.ERCaseLegalHoldCommand) (*domain.ERCase, error) {
	if cmd.TenantID == uuid.Nil || cmd.ID == uuid.Nil || (cmd.Enabled && strings.TrimSpace(valueOrEmpty(cmd.Reason)) == "") {
		return nil, domain.ErrInvalidERCase
	}
	result, err := s.employeeRelations.UpdateERCaseLegalHold(ctx, cmd.TenantID, cmd.ID, cmd.Enabled, cmd.Reason, cmd.ActorID)
	if err != nil {
		s.logError("update er case legal hold", err, serviceTenantIDField(cmd.TenantID), serviceStringField("er_case_id", cmd.ID.String()))
		return nil, err
	}
	eventType := "legal_hold_released"
	if cmd.Enabled {
		eventType = "legal_hold_enabled"
	}
	s.recordERCaseEvent(ctx, result, eventType, nil, &result.Status, cmd.Reason, cmd.ActorID, nil)
	return result, nil
}

func (s *TenantService) ListERCases(ctx context.Context, filter domain.ERCaseFilter) (*domain.ERCasePage, error) {
	items, err := s.employeeRelations.ListERCases(ctx, filter)
	if err != nil {
		s.logError("list er cases", err, serviceTenantIDField(filter.TenantID))
		return nil, err
	}
	total, err := s.employeeRelations.CountERCases(ctx, filter)
	if err != nil {
		s.logError("count er cases", err, serviceTenantIDField(filter.TenantID))
		return nil, err
	}
	summary, _ := s.employeeRelations.GetERCaseSummary(ctx, filter.TenantID)
	categories, _ := s.employeeRelations.ListERCaseCategories(ctx, filter.TenantID, nil)
	return &domain.ERCasePage{Items: items, Total: total, Summary: summary, Categories: categories}, nil
}

func (s *TenantService) GetERCaseWorkspace(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.ERCaseWorkspace, error) {
	item, err := s.employeeRelations.GetERCase(ctx, tenantID, id)
	if err != nil {
		s.logError("get er case workspace case", err, serviceTenantIDField(tenantID), serviceStringField("er_case_id", id.String()))
		return nil, err
	}
	parties, _ := s.employeeRelations.ListERCaseParties(ctx, tenantID, id)
	allegations, _ := s.employeeRelations.ListERAllegations(ctx, tenantID, id)
	steps, _ := s.employeeRelations.ListERInvestigationSteps(ctx, tenantID, id)
	witnesses, _ := s.employeeRelations.ListERWitnessNotes(ctx, tenantID, id)
	evidence, _ := s.employeeRelations.ListEREvidenceAttachments(ctx, tenantID, id)
	findings, _ := s.employeeRelations.ListERFindings(ctx, tenantID, id)
	actions, _ := s.employeeRelations.ListERActionPlans(ctx, tenantID, id)
	events, _ := s.employeeRelations.ListERCaseEvents(ctx, tenantID, id)
	return &domain.ERCaseWorkspace{Case: item, Parties: parties, Allegations: allegations, Steps: steps, Witnesses: witnesses, Evidence: evidence, Findings: findings, Actions: actions, Events: events}, nil
}

func (s *TenantService) CreateERCaseParty(ctx context.Context, cmd ports.ERCasePartyCommand) (*domain.ERCaseParty, error) {
	item, err := domain.NewERCaseParty(domain.ERCaseParty{TenantID: cmd.TenantID, ERCaseID: cmd.ERCaseID, PartyUserID: cmd.PartyUserID, PartyName: cmd.PartyName, PartyRole: cmd.PartyRole, RepresentationNotes: cmd.RepresentationNotes, ContactNotes: cmd.ContactNotes})
	if err != nil {
		return nil, err
	}
	result, err := s.employeeRelations.CreateERCaseParty(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create er case party", err, serviceTenantIDField(cmd.TenantID), serviceStringField("er_case_id", cmd.ERCaseID.String()))
		return nil, err
	}
	s.recordERCaseEvent(ctx, &domain.ERCase{TenantID: cmd.TenantID, ID: cmd.ERCaseID}, "party_added", nil, nil, nil, cmd.ActorID, nil)
	return result, nil
}

func (s *TenantService) CreateERAllegation(ctx context.Context, cmd ports.ERAllegationCommand) (*domain.ERAllegation, error) {
	incidentDate, err := parseEROptionalDate(cmd.IncidentDate)
	if err != nil {
		return nil, domain.ErrInvalidERAllegation
	}
	item, err := domain.NewERAllegation(domain.ERAllegation{TenantID: cmd.TenantID, ERCaseID: cmd.ERCaseID, AllegationType: cmd.AllegationType, IncidentDate: incidentDate, IncidentLocation: cmd.IncidentLocation, Description: cmd.Description, PolicyReference: cmd.PolicyReference, Status: cmd.Status})
	if err != nil {
		return nil, err
	}
	result, err := s.employeeRelations.CreateERAllegation(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create er allegation", err, serviceTenantIDField(cmd.TenantID), serviceStringField("er_case_id", cmd.ERCaseID.String()))
		return nil, err
	}
	s.recordERCaseEvent(ctx, &domain.ERCase{TenantID: cmd.TenantID, ID: cmd.ERCaseID}, "allegation_added", nil, nil, nil, cmd.ActorID, nil)
	return result, nil
}

func (s *TenantService) CreateERInvestigationStep(ctx context.Context, cmd ports.ERInvestigationStepCommand) (*domain.ERInvestigationStep, error) {
	dueAt, err := parseEROptionalTime(cmd.DueAt)
	if err != nil {
		return nil, domain.ErrInvalidERStep
	}
	completedAt, err := parseEROptionalTime(cmd.CompletedAt)
	if err != nil {
		return nil, domain.ErrInvalidERStep
	}
	item, err := domain.NewERInvestigationStep(domain.ERInvestigationStep{TenantID: cmd.TenantID, ERCaseID: cmd.ERCaseID, StepType: cmd.StepType, Title: cmd.Title, Description: cmd.Description, OwnerUserID: cmd.OwnerUserID, DueAt: dueAt, CompletedAt: completedAt, Status: cmd.Status, OutcomeNotes: cmd.OutcomeNotes})
	if err != nil {
		return nil, err
	}
	result, err := s.employeeRelations.CreateERInvestigationStep(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create er investigation step", err, serviceTenantIDField(cmd.TenantID), serviceStringField("er_case_id", cmd.ERCaseID.String()))
		return nil, err
	}
	s.recordERCaseEvent(ctx, &domain.ERCase{TenantID: cmd.TenantID, ID: cmd.ERCaseID}, "investigation_step_added", nil, nil, nil, cmd.ActorID, nil)
	return result, nil
}

func (s *TenantService) UpdateERInvestigationStepStatus(ctx context.Context, cmd ports.ERInvestigationStepCommand) (*domain.ERInvestigationStep, error) {
	completedAt, err := parseEROptionalTime(cmd.CompletedAt)
	if err != nil {
		return nil, domain.ErrInvalidERStep
	}
	result, err := s.employeeRelations.UpdateERInvestigationStepStatus(ctx, cmd.TenantID, cmd.ID, cmd.Status, completedAt, cmd.OutcomeNotes, cmd.ActorID)
	if err != nil {
		s.logError("update er investigation step status", err, serviceTenantIDField(cmd.TenantID), serviceStringField("er_step_id", cmd.ID.String()))
		return nil, err
	}
	s.recordERCaseEvent(ctx, &domain.ERCase{TenantID: result.TenantID, ID: result.ERCaseID}, "investigation_step_status_changed", nil, nil, &result.Status, cmd.ActorID, nil)
	return result, nil
}

func (s *TenantService) CreateERWitnessNote(ctx context.Context, cmd ports.ERWitnessNoteCommand) (*domain.ERWitnessNote, error) {
	interviewAt, err := parseEROptionalTime(cmd.InterviewAt)
	if err != nil {
		return nil, domain.ErrInvalidERWitnessNote
	}
	item, err := domain.NewERWitnessNote(domain.ERWitnessNote{TenantID: cmd.TenantID, ERCaseID: cmd.ERCaseID, WitnessUserID: cmd.WitnessUserID, WitnessName: cmd.WitnessName, InterviewAt: interviewAt, InterviewerUserID: cmd.InterviewerUserID, StatementSummary: cmd.StatementSummary, ConfidentialityLevel: cmd.ConfidentialityLevel})
	if err != nil {
		return nil, err
	}
	result, err := s.employeeRelations.CreateERWitnessNote(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create er witness note", err, serviceTenantIDField(cmd.TenantID), serviceStringField("er_case_id", cmd.ERCaseID.String()))
		return nil, err
	}
	s.recordERCaseEvent(ctx, &domain.ERCase{TenantID: cmd.TenantID, ID: cmd.ERCaseID}, "witness_note_added", nil, nil, nil, cmd.ActorID, nil)
	return result, nil
}

func (s *TenantService) CreateEREvidenceAttachment(ctx context.Context, cmd ports.EREvidenceAttachmentCommand) (*domain.EREvidenceAttachment, error) {
	if s.objectStorage == nil {
		return nil, domain.ErrStorageProviderSettingsNotFound
	}
	content, err := base64.StdEncoding.DecodeString(cmd.FileContentBase64)
	if err != nil || len(content) == 0 {
		return nil, domain.ErrInvalidEREvidence
	}
	settings, err := s.resolveWorkflowStorageSettings(ctx, cmd.TenantID)
	if err != nil {
		s.logError("resolve er evidence storage", err, serviceTenantIDField(cmd.TenantID), serviceStringField("er_case_id", cmd.ERCaseID.String()))
		return nil, err
	}
	entityID := uuid.New()
	storagePath, err := s.objectStorage.PutObject(ctx, settings, ports.StoreObjectInput{TenantID: cmd.TenantID, Category: ports.StorageCategoryEREvidence, OwnerID: cmd.ERCaseID, EntityID: entityID, FileName: cmd.FileName, ContentType: cmd.ContentType, Content: content})
	if err != nil {
		s.logError("store er evidence", err, serviceTenantIDField(cmd.TenantID), serviceStringField("er_case_id", cmd.ERCaseID.String()))
		return nil, err
	}
	sum := sha256.Sum256(content)
	checksum := hex.EncodeToString(sum[:])
	item, err := domain.NewEREvidenceAttachment(domain.EREvidenceAttachment{ID: entityID, TenantID: cmd.TenantID, ERCaseID: cmd.ERCaseID, AllegationID: cmd.AllegationID, FileName: cmd.FileName, ContentType: cmd.ContentType, StoragePath: storagePath, ChecksumSHA: &checksum, SizeBytes: int64(len(content)), EvidenceType: cmd.EvidenceType, Description: cmd.Description, UploadedBy: cmd.ActorID, LegalHold: cmd.LegalHold})
	if err != nil {
		return nil, err
	}
	result, err := s.employeeRelations.CreateEREvidenceAttachment(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create er evidence attachment", err, serviceTenantIDField(cmd.TenantID), serviceStringField("er_case_id", cmd.ERCaseID.String()))
		return nil, err
	}
	s.recordERCaseEvent(ctx, &domain.ERCase{TenantID: cmd.TenantID, ID: cmd.ERCaseID}, "evidence_added", nil, nil, nil, cmd.ActorID, json.RawMessage(`{"storage":"tenant_storage"}`))
	return result, nil
}

func (s *TenantService) CreateERFinding(ctx context.Context, cmd ports.ERFindingCommand) (*domain.ERFinding, error) {
	now := time.Now().UTC()
	item, err := domain.NewERFinding(domain.ERFinding{TenantID: cmd.TenantID, ERCaseID: cmd.ERCaseID, AllegationID: cmd.AllegationID, Finding: cmd.Finding, Rationale: cmd.Rationale, RecommendedAction: cmd.RecommendedAction, DecidedBy: cmd.ActorID, DecidedAt: &now})
	if err != nil {
		return nil, err
	}
	result, err := s.employeeRelations.CreateERFinding(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create er finding", err, serviceTenantIDField(cmd.TenantID), serviceStringField("er_case_id", cmd.ERCaseID.String()))
		return nil, err
	}
	s.recordERCaseEvent(ctx, &domain.ERCase{TenantID: cmd.TenantID, ID: cmd.ERCaseID}, "finding_added", nil, nil, nil, cmd.ActorID, nil)
	return result, nil
}

func (s *TenantService) CreateERActionPlan(ctx context.Context, cmd ports.ERActionPlanCommand) (*domain.ERActionPlan, error) {
	dueAt, err := parseEROptionalTime(cmd.DueAt)
	if err != nil {
		return nil, domain.ErrInvalidERActionPlan
	}
	completedAt, err := parseEROptionalTime(cmd.CompletedAt)
	if err != nil {
		return nil, domain.ErrInvalidERActionPlan
	}
	item, err := domain.NewERActionPlan(domain.ERActionPlan{TenantID: cmd.TenantID, ERCaseID: cmd.ERCaseID, ActionType: cmd.ActionType, Description: cmd.Description, AssignedToUserID: cmd.AssignedToUserID, DueAt: dueAt, CompletedAt: completedAt, Status: cmd.Status, FollowUpNotes: cmd.FollowUpNotes})
	if err != nil {
		return nil, err
	}
	result, err := s.employeeRelations.CreateERActionPlan(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create er action plan", err, serviceTenantIDField(cmd.TenantID), serviceStringField("er_case_id", cmd.ERCaseID.String()))
		return nil, err
	}
	s.recordERCaseEvent(ctx, &domain.ERCase{TenantID: cmd.TenantID, ID: cmd.ERCaseID}, "action_plan_added", nil, nil, nil, cmd.ActorID, nil)
	return result, nil
}

func (s *TenantService) UpdateERActionPlanStatus(ctx context.Context, cmd ports.ERActionPlanCommand) (*domain.ERActionPlan, error) {
	completedAt, err := parseEROptionalTime(cmd.CompletedAt)
	if err != nil {
		return nil, domain.ErrInvalidERActionPlan
	}
	result, err := s.employeeRelations.UpdateERActionPlanStatus(ctx, cmd.TenantID, cmd.ID, cmd.Status, completedAt, cmd.FollowUpNotes, cmd.ActorID)
	if err != nil {
		s.logError("update er action plan status", err, serviceTenantIDField(cmd.TenantID), serviceStringField("er_action_plan_id", cmd.ID.String()))
		return nil, err
	}
	s.recordERCaseEvent(ctx, &domain.ERCase{TenantID: result.TenantID, ID: result.ERCaseID}, "action_plan_status_changed", nil, nil, &result.Status, cmd.ActorID, nil)
	return result, nil
}

func (s *TenantService) recordERCaseEvent(ctx context.Context, item *domain.ERCase, eventType string, fromStatus *string, toStatus *string, comment *string, actorID *uuid.UUID, metadata json.RawMessage) {
	if item == nil || item.TenantID == uuid.Nil || item.ID == uuid.Nil {
		return
	}
	if len(metadata) == 0 {
		metadata = json.RawMessage(`{}`)
	}
	_, _ = s.employeeRelations.CreateERCaseEvent(ctx, &domain.ERCaseEvent{TenantID: item.TenantID, ERCaseID: item.ID, EventType: eventType, FromStatus: fromStatus, ToStatus: toStatus, ActorUserID: actorID, Comment: comment, Metadata: metadata}, actorID)
}

func erCaseCategoryFromCommand(cmd ports.ERCaseCategoryCommand) (*domain.ERCaseCategory, error) {
	return domain.NewERCaseCategory(domain.ERCaseCategory{TenantID: cmd.TenantID, Code: cmd.Code, Name: cmd.Name, CaseFamily: cmd.CaseFamily, Description: cmd.Description, DefaultSeverity: cmd.DefaultSeverity, DefaultOwnerRole: cmd.DefaultOwnerRole, IsActive: cmd.IsActive})
}

func erCaseFromCommand(cmd ports.ERCaseCommand) (*domain.ERCase, error) {
	dueAt, err := parseEROptionalTime(cmd.DueAt)
	if err != nil {
		return nil, domain.ErrInvalidERCase
	}
	return domain.NewERCase(domain.ERCase{TenantID: cmd.TenantID, SourceHRCaseID: cmd.SourceHRCaseID, CategoryID: cmd.CategoryID, Title: cmd.Title, IntakeSummary: cmd.IntakeSummary, CaseFamily: cmd.CaseFamily, Severity: cmd.Severity, Status: cmd.Status, ConfidentialityLevel: cmd.ConfidentialityLevel, ComplainantUserID: cmd.ComplainantUserID, SubjectEmployeeUserID: cmd.SubjectEmployeeUserID, OwnerUserID: cmd.OwnerUserID, OwnerRole: cmd.OwnerRole, InvestigationLeadUserID: cmd.InvestigationLeadUserID, DueAt: dueAt, PrivacyNotes: cmd.PrivacyNotes})
}

func cleanERCaseStatus(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case domain.ERCaseStatusIntake, domain.ERCaseStatusTriage, domain.ERCaseStatusInvestigation, domain.ERCaseStatusFindings, domain.ERCaseStatusActionPlan, domain.ERCaseStatusMonitoring, domain.ERCaseStatusClosed, domain.ERCaseStatusCancelled:
		return strings.ToLower(strings.TrimSpace(value))
	default:
		return ""
	}
}

func erStatusTransitionAllowed(from string, to string) bool {
	if from == to {
		return true
	}
	if from == domain.ERCaseStatusClosed || from == domain.ERCaseStatusCancelled {
		return false
	}
	return true
}

func parseEROptionalDate(value *string) (*time.Time, error) {
	if value == nil || strings.TrimSpace(*value) == "" {
		return nil, nil
	}
	parsed, err := time.Parse("2006-01-02", strings.TrimSpace(*value))
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

func parseEROptionalTime(value *string) (*time.Time, error) {
	if value == nil || strings.TrimSpace(*value) == "" {
		return nil, nil
	}
	raw := strings.TrimSpace(*value)
	if parsed, err := time.Parse(time.RFC3339, raw); err == nil {
		return &parsed, nil
	}
	parsed, err := time.Parse("2006-01-02", raw)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

func valueOrEmpty(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
