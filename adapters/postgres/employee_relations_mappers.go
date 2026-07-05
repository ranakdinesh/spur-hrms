package postgres

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func mapERCaseCategory(row sqlc.HrmsErCaseCategory) *domain.ERCaseCategory {
	return &domain.ERCaseCategory{ID: row.ID, TenantID: row.TenantID, Code: row.Code, Name: row.Name, CaseFamily: row.CaseFamily, Description: ptrFromText(row.Description), DefaultSeverity: row.DefaultSeverity, DefaultOwnerRole: ptrFromText(row.DefaultOwnerRole), IsActive: row.IsActive, Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy)}
}

func mapERCaseCategories(rows []sqlc.HrmsErCaseCategory) []*domain.ERCaseCategory {
	items := make([]*domain.ERCaseCategory, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapERCaseCategory(row))
	}
	return items
}

func mapERCase(row sqlc.HrmsErCase) *domain.ERCase {
	return erCaseFromParts(row.ID, row.TenantID, row.CaseNumber, row.SourceHrCaseID, row.CategoryID, row.Title, row.IntakeSummary, row.CaseFamily, row.Severity, row.Status, row.ConfidentialityLevel, row.ComplainantUserID, row.SubjectEmployeeUserID, row.OwnerUserID, row.OwnerRole, row.InvestigationLeadUserID, row.LegalHold, row.LegalHoldReason, row.LegalHoldAt, row.LegalHoldBy, row.DueAt, row.ClosedAt, row.ResolutionSummary, row.PrivacyNotes, row.Inactive, row.CreatedAt, row.CreatedBy, row.UpdatedAt, row.UpdatedBy, pgtype.Text{}, pgtype.Text{}, 0, 0, 0)
}

func mapERCaseListRow(row sqlc.ListERCasesRow) *domain.ERCase {
	return erCaseFromParts(row.ID, row.TenantID, row.CaseNumber, row.SourceHrCaseID, row.CategoryID, row.Title, row.IntakeSummary, row.CaseFamily, row.Severity, row.Status, row.ConfidentialityLevel, row.ComplainantUserID, row.SubjectEmployeeUserID, row.OwnerUserID, row.OwnerRole, row.InvestigationLeadUserID, row.LegalHold, row.LegalHoldReason, row.LegalHoldAt, row.LegalHoldBy, row.DueAt, row.ClosedAt, row.ResolutionSummary, row.PrivacyNotes, row.Inactive, row.CreatedAt, row.CreatedBy, row.UpdatedAt, row.UpdatedBy, row.CategoryName, row.CategoryCode, row.AllegationCount, row.EvidenceCount, row.OpenActionCount)
}

func mapERCaseGetRow(row sqlc.GetERCaseRow) *domain.ERCase {
	return erCaseFromParts(row.ID, row.TenantID, row.CaseNumber, row.SourceHrCaseID, row.CategoryID, row.Title, row.IntakeSummary, row.CaseFamily, row.Severity, row.Status, row.ConfidentialityLevel, row.ComplainantUserID, row.SubjectEmployeeUserID, row.OwnerUserID, row.OwnerRole, row.InvestigationLeadUserID, row.LegalHold, row.LegalHoldReason, row.LegalHoldAt, row.LegalHoldBy, row.DueAt, row.ClosedAt, row.ResolutionSummary, row.PrivacyNotes, row.Inactive, row.CreatedAt, row.CreatedBy, row.UpdatedAt, row.UpdatedBy, row.CategoryName, row.CategoryCode, 0, 0, 0)
}

func erCaseFromParts(id uuid.UUID, tenantID uuid.UUID, caseNumber string, sourceHRCaseID pgtype.UUID, categoryID pgtype.UUID, title string, intakeSummary string, family string, severity string, status string, confidentiality string, complainantID pgtype.UUID, subjectID pgtype.UUID, ownerID pgtype.UUID, ownerRole pgtype.Text, leadID pgtype.UUID, legalHold bool, legalHoldReason pgtype.Text, legalHoldAt pgtype.Timestamptz, legalHoldBy pgtype.UUID, dueAt pgtype.Timestamptz, closedAt pgtype.Timestamptz, resolutionSummary pgtype.Text, privacyNotes pgtype.Text, inactive bool, createdAt pgtype.Timestamptz, createdBy pgtype.UUID, updatedAt pgtype.Timestamptz, updatedBy pgtype.UUID, categoryName pgtype.Text, categoryCode pgtype.Text, allegationCount int32, evidenceCount int32, openActionCount int32) *domain.ERCase {
	return &domain.ERCase{ID: id, TenantID: tenantID, CaseNumber: caseNumber, SourceHRCaseID: ptrFromUUID(sourceHRCaseID), CategoryID: ptrFromUUID(categoryID), CategoryName: ptrFromText(categoryName), CategoryCode: ptrFromText(categoryCode), Title: title, IntakeSummary: intakeSummary, CaseFamily: family, Severity: severity, Status: status, ConfidentialityLevel: confidentiality, ComplainantUserID: ptrFromUUID(complainantID), SubjectEmployeeUserID: ptrFromUUID(subjectID), OwnerUserID: ptrFromUUID(ownerID), OwnerRole: ptrFromText(ownerRole), InvestigationLeadUserID: ptrFromUUID(leadID), LegalHold: legalHold, LegalHoldReason: ptrFromText(legalHoldReason), LegalHoldAt: ptrFromTimestamptz(legalHoldAt), LegalHoldBy: ptrFromUUID(legalHoldBy), DueAt: ptrFromTimestamptz(dueAt), ClosedAt: ptrFromTimestamptz(closedAt), ResolutionSummary: ptrFromText(resolutionSummary), PrivacyNotes: ptrFromText(privacyNotes), Inactive: inactive, CreatedAt: timeFromTimestamptz(createdAt), CreatedBy: ptrFromUUID(createdBy), UpdatedAt: timeFromTimestamptz(updatedAt), UpdatedBy: ptrFromUUID(updatedBy), AllegationCount: allegationCount, EvidenceCount: evidenceCount, OpenActionCount: openActionCount}
}

func mapERCases(rows []sqlc.ListERCasesRow) []*domain.ERCase {
	items := make([]*domain.ERCase, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapERCaseListRow(row))
	}
	return items
}

func mapERCaseParty(row sqlc.HrmsErCaseParty) *domain.ERCaseParty {
	return &domain.ERCaseParty{ID: row.ID, TenantID: row.TenantID, ERCaseID: row.ErCaseID, PartyUserID: ptrFromUUID(row.PartyUserID), PartyName: ptrFromText(row.PartyName), PartyRole: row.PartyRole, RepresentationNotes: ptrFromText(row.RepresentationNotes), ContactNotes: ptrFromText(row.ContactNotes), CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy)}
}

func mapERCaseParties(rows []sqlc.HrmsErCaseParty) []*domain.ERCaseParty {
	items := make([]*domain.ERCaseParty, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapERCaseParty(row))
	}
	return items
}

func mapERAllegation(row sqlc.HrmsErAllegation) *domain.ERAllegation {
	return &domain.ERAllegation{ID: row.ID, TenantID: row.TenantID, ERCaseID: row.ErCaseID, AllegationType: row.AllegationType, IncidentDate: ptrFromDate(row.IncidentDate), IncidentLocation: ptrFromText(row.IncidentLocation), Description: row.Description, PolicyReference: ptrFromText(row.PolicyReference), Status: row.Status, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy)}
}

func mapERAllegations(rows []sqlc.HrmsErAllegation) []*domain.ERAllegation {
	items := make([]*domain.ERAllegation, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapERAllegation(row))
	}
	return items
}

func mapERInvestigationStep(row sqlc.HrmsErInvestigationStep) *domain.ERInvestigationStep {
	return &domain.ERInvestigationStep{ID: row.ID, TenantID: row.TenantID, ERCaseID: row.ErCaseID, StepType: row.StepType, Title: row.Title, Description: ptrFromText(row.Description), OwnerUserID: ptrFromUUID(row.OwnerUserID), DueAt: ptrFromTimestamptz(row.DueAt), CompletedAt: ptrFromTimestamptz(row.CompletedAt), Status: row.Status, OutcomeNotes: ptrFromText(row.OutcomeNotes), CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy)}
}

func mapERInvestigationSteps(rows []sqlc.HrmsErInvestigationStep) []*domain.ERInvestigationStep {
	items := make([]*domain.ERInvestigationStep, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapERInvestigationStep(row))
	}
	return items
}

func mapERWitnessNote(row sqlc.HrmsErWitnessNote) *domain.ERWitnessNote {
	return &domain.ERWitnessNote{ID: row.ID, TenantID: row.TenantID, ERCaseID: row.ErCaseID, WitnessUserID: ptrFromUUID(row.WitnessUserID), WitnessName: ptrFromText(row.WitnessName), InterviewAt: ptrFromTimestamptz(row.InterviewAt), InterviewerUserID: ptrFromUUID(row.InterviewerUserID), StatementSummary: row.StatementSummary, ConfidentialityLevel: row.ConfidentialityLevel, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy)}
}

func mapERWitnessNotes(rows []sqlc.HrmsErWitnessNote) []*domain.ERWitnessNote {
	items := make([]*domain.ERWitnessNote, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapERWitnessNote(row))
	}
	return items
}

func mapEREvidenceAttachment(row sqlc.HrmsErEvidenceAttachment) *domain.EREvidenceAttachment {
	return &domain.EREvidenceAttachment{ID: row.ID, TenantID: row.TenantID, ERCaseID: row.ErCaseID, AllegationID: ptrFromUUID(row.AllegationID), FileName: row.FileName, ContentType: row.ContentType, StoragePath: row.StoragePath, ChecksumSHA: ptrFromText(row.ChecksumSha256), SizeBytes: row.SizeBytes, EvidenceType: row.EvidenceType, Description: ptrFromText(row.Description), UploadedBy: ptrFromUUID(row.UploadedBy), LegalHold: row.LegalHold, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy)}
}

func mapEREvidenceAttachments(rows []sqlc.HrmsErEvidenceAttachment) []*domain.EREvidenceAttachment {
	items := make([]*domain.EREvidenceAttachment, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapEREvidenceAttachment(row))
	}
	return items
}

func mapERFinding(row sqlc.HrmsErFinding) *domain.ERFinding {
	return &domain.ERFinding{ID: row.ID, TenantID: row.TenantID, ERCaseID: row.ErCaseID, AllegationID: ptrFromUUID(row.AllegationID), Finding: row.Finding, Rationale: row.Rationale, RecommendedAction: ptrFromText(row.RecommendedAction), DecidedBy: ptrFromUUID(row.DecidedBy), DecidedAt: ptrFromTimestamptz(row.DecidedAt), CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy)}
}

func mapERFindings(rows []sqlc.HrmsErFinding) []*domain.ERFinding {
	items := make([]*domain.ERFinding, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapERFinding(row))
	}
	return items
}

func mapERActionPlan(row sqlc.HrmsErActionPlan) *domain.ERActionPlan {
	return &domain.ERActionPlan{ID: row.ID, TenantID: row.TenantID, ERCaseID: row.ErCaseID, ActionType: row.ActionType, Description: row.Description, AssignedToUserID: ptrFromUUID(row.AssignedToUserID), DueAt: ptrFromTimestamptz(row.DueAt), CompletedAt: ptrFromTimestamptz(row.CompletedAt), Status: row.Status, FollowUpNotes: ptrFromText(row.FollowUpNotes), CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy)}
}

func mapERActionPlans(rows []sqlc.HrmsErActionPlan) []*domain.ERActionPlan {
	items := make([]*domain.ERActionPlan, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapERActionPlan(row))
	}
	return items
}

func mapERCaseEvent(row sqlc.HrmsErCaseEvent) *domain.ERCaseEvent {
	metadata := json.RawMessage(row.Metadata)
	if len(metadata) == 0 {
		metadata = json.RawMessage(`{}`)
	}
	return &domain.ERCaseEvent{ID: row.ID, TenantID: row.TenantID, ERCaseID: row.ErCaseID, EventType: row.EventType, FromStatus: ptrFromText(row.FromStatus), ToStatus: ptrFromText(row.ToStatus), ActorUserID: ptrFromUUID(row.ActorUserID), Comment: ptrFromText(row.Comment), Metadata: metadata, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy)}
}

func mapERCaseEvents(rows []sqlc.HrmsErCaseEvent) []*domain.ERCaseEvent {
	items := make([]*domain.ERCaseEvent, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapERCaseEvent(row))
	}
	return items
}

func mapERCaseSummary(rows []sqlc.GetERCaseSummaryRow) []*domain.ERCaseSummaryRow {
	items := make([]*domain.ERCaseSummaryRow, 0, len(rows))
	for _, row := range rows {
		items = append(items, &domain.ERCaseSummaryRow{Status: row.Status, Severity: row.Severity, CaseCount: row.CaseCount, LegalHoldCount: row.LegalHoldCount, OverdueCount: row.OverdueCount})
	}
	return items
}
