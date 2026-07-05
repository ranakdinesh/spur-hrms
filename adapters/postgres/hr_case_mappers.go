package postgres

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func mapHRCaseCategory(row sqlc.HrmsHrCaseCategory) *domain.HRCaseCategory {
	return &domain.HRCaseCategory{
		ID:                     row.ID,
		TenantID:               row.TenantID,
		Code:                   row.Code,
		Name:                   row.Name,
		Description:            ptrFromText(row.Description),
		ConfidentialityDefault: row.ConfidentialityDefault,
		DefaultOwnerRole:       ptrFromText(row.DefaultOwnerRole),
		IsActive:               row.IsActive,
		Inactive:               row.Inactive,
		CreatedAt:              timeFromTimestamptz(row.CreatedAt),
		CreatedBy:              ptrFromUUID(row.CreatedBy),
		UpdatedAt:              timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:              ptrFromUUID(row.UpdatedBy),
	}
}

func mapHRCaseCategories(rows []sqlc.HrmsHrCaseCategory) []*domain.HRCaseCategory {
	items := make([]*domain.HRCaseCategory, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapHRCaseCategory(row))
	}
	return items
}

func mapHRCaseSLAPolicy(row sqlc.HrmsHrCaseSlaPolicy) *domain.HRCaseSLAPolicy {
	return &domain.HRCaseSLAPolicy{
		ID:              row.ID,
		TenantID:        row.TenantID,
		CategoryID:      ptrFromUUID(row.CategoryID),
		Priority:        row.Priority,
		ResponseHours:   row.ResponseHours,
		ResolutionHours: row.ResolutionHours,
		EscalationHours: row.EscalationHours,
		IsActive:        row.IsActive,
		Inactive:        row.Inactive,
		CreatedAt:       timeFromTimestamptz(row.CreatedAt),
		CreatedBy:       ptrFromUUID(row.CreatedBy),
		UpdatedAt:       timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:       ptrFromUUID(row.UpdatedBy),
	}
}

func mapHRCaseSLAPolicyListRow(row sqlc.ListHRCaseSLAPoliciesRow) *domain.HRCaseSLAPolicy {
	item := &domain.HRCaseSLAPolicy{
		ID:              row.ID,
		TenantID:        row.TenantID,
		CategoryID:      ptrFromUUID(row.CategoryID),
		CategoryName:    ptrFromText(row.CategoryName),
		Priority:        row.Priority,
		ResponseHours:   row.ResponseHours,
		ResolutionHours: row.ResolutionHours,
		EscalationHours: row.EscalationHours,
		IsActive:        row.IsActive,
		Inactive:        row.Inactive,
		CreatedAt:       timeFromTimestamptz(row.CreatedAt),
		CreatedBy:       ptrFromUUID(row.CreatedBy),
		UpdatedAt:       timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:       ptrFromUUID(row.UpdatedBy),
	}
	return item
}

func mapHRCaseSLAPolicyList(rows []sqlc.ListHRCaseSLAPoliciesRow) []*domain.HRCaseSLAPolicy {
	items := make([]*domain.HRCaseSLAPolicy, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapHRCaseSLAPolicyListRow(row))
	}
	return items
}

func mapHRCase(row sqlc.HrmsHrCase) *domain.HRCase {
	return hrCaseFromParts(row.ID, row.TenantID, row.CaseNumber, row.CategoryID, row.CaseType, row.Title, row.Description, row.ConfidentialityLevel, row.RequesterUserID, row.SubjectEmployeeUserID, row.OwnerUserID, row.OwnerRole, row.Status, row.Priority, row.SourceChannel, row.FirstResponseDueAt, row.FirstRespondedAt, row.DueAt, row.ResolvedAt, row.ClosedAt, row.EscalatedAt, row.EscalationLevel, row.LastActivityAt, row.ResolutionSummary, row.Inactive, row.CreatedAt, row.CreatedBy, row.UpdatedAt, row.UpdatedBy, pgtype.Text{}, pgtype.Text{}, pgtype.Text{}, pgtype.Text{}, pgtype.Text{}, 0, 0)
}

func mapHRCaseListRow(row sqlc.ListHRCasesRow) *domain.HRCase {
	return hrCaseFromParts(row.ID, row.TenantID, row.CaseNumber, row.CategoryID, row.CaseType, row.Title, row.Description, row.ConfidentialityLevel, row.RequesterUserID, row.SubjectEmployeeUserID, row.OwnerUserID, row.OwnerRole, row.Status, row.Priority, row.SourceChannel, row.FirstResponseDueAt, row.FirstRespondedAt, row.DueAt, row.ResolvedAt, row.ClosedAt, row.EscalatedAt, row.EscalationLevel, row.LastActivityAt, row.ResolutionSummary, row.Inactive, row.CreatedAt, row.CreatedBy, row.UpdatedAt, row.UpdatedBy, row.CategoryName, row.CategoryCode, row.RequesterEmail, row.SubjectEmail, row.OwnerEmail, row.CommentCount, row.AttachmentCount)
}

func mapHRCaseGetRow(row sqlc.GetHRCaseRow) *domain.HRCase {
	return hrCaseFromParts(row.ID, row.TenantID, row.CaseNumber, row.CategoryID, row.CaseType, row.Title, row.Description, row.ConfidentialityLevel, row.RequesterUserID, row.SubjectEmployeeUserID, row.OwnerUserID, row.OwnerRole, row.Status, row.Priority, row.SourceChannel, row.FirstResponseDueAt, row.FirstRespondedAt, row.DueAt, row.ResolvedAt, row.ClosedAt, row.EscalatedAt, row.EscalationLevel, row.LastActivityAt, row.ResolutionSummary, row.Inactive, row.CreatedAt, row.CreatedBy, row.UpdatedAt, row.UpdatedBy, row.CategoryName, row.CategoryCode, row.RequesterEmail, row.SubjectEmail, row.OwnerEmail, 0, 0)
}

func mapHRCaseList(rows []sqlc.ListHRCasesRow) []*domain.HRCase {
	items := make([]*domain.HRCase, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapHRCaseListRow(row))
	}
	return items
}

func hrCaseFromParts(id uuid.UUID, tenantID uuid.UUID, caseNumber string, categoryID pgtype.UUID, caseType string, title string, description string, confidentiality string, requesterID pgtype.UUID, subjectID pgtype.UUID, ownerID pgtype.UUID, ownerRole pgtype.Text, status string, priority string, sourceChannel string, firstResponseDueAt pgtype.Timestamptz, firstRespondedAt pgtype.Timestamptz, dueAt pgtype.Timestamptz, resolvedAt pgtype.Timestamptz, closedAt pgtype.Timestamptz, escalatedAt pgtype.Timestamptz, escalationLevel int32, lastActivityAt pgtype.Timestamptz, resolutionSummary pgtype.Text, inactive bool, createdAt pgtype.Timestamptz, createdBy pgtype.UUID, updatedAt pgtype.Timestamptz, updatedBy pgtype.UUID, categoryName pgtype.Text, categoryCode pgtype.Text, requesterEmail pgtype.Text, subjectEmail pgtype.Text, ownerEmail pgtype.Text, commentCount int32, attachmentCount int32) *domain.HRCase {
	return &domain.HRCase{
		ID:                    id,
		TenantID:              tenantID,
		CaseNumber:            caseNumber,
		CategoryID:            ptrFromUUID(categoryID),
		CategoryName:          ptrFromText(categoryName),
		CategoryCode:          ptrFromText(categoryCode),
		CaseType:              caseType,
		Title:                 title,
		Description:           description,
		ConfidentialityLevel:  confidentiality,
		RequesterUserID:       ptrFromUUID(requesterID),
		SubjectEmployeeUserID: ptrFromUUID(subjectID),
		OwnerUserID:           ptrFromUUID(ownerID),
		OwnerRole:             ptrFromText(ownerRole),
		Status:                status,
		Priority:              priority,
		SourceChannel:         sourceChannel,
		FirstResponseDueAt:    ptrFromTimestamptz(firstResponseDueAt),
		FirstRespondedAt:      ptrFromTimestamptz(firstRespondedAt),
		DueAt:                 ptrFromTimestamptz(dueAt),
		ResolvedAt:            ptrFromTimestamptz(resolvedAt),
		ClosedAt:              ptrFromTimestamptz(closedAt),
		EscalatedAt:           ptrFromTimestamptz(escalatedAt),
		EscalationLevel:       escalationLevel,
		LastActivityAt:        timeFromTimestamptz(lastActivityAt),
		ResolutionSummary:     ptrFromText(resolutionSummary),
		Inactive:              inactive,
		CreatedAt:             timeFromTimestamptz(createdAt),
		CreatedBy:             ptrFromUUID(createdBy),
		UpdatedAt:             timeFromTimestamptz(updatedAt),
		UpdatedBy:             ptrFromUUID(updatedBy),
		RequesterEmail:        ptrFromText(requesterEmail),
		SubjectEmail:          ptrFromText(subjectEmail),
		OwnerEmail:            ptrFromText(ownerEmail),
		CommentCount:          commentCount,
		AttachmentCount:       attachmentCount,
	}
}

func mapHRCaseComment(row sqlc.HrmsHrCaseComment) *domain.HRCaseComment {
	return &domain.HRCaseComment{ID: row.ID, TenantID: row.TenantID, CaseID: row.CaseID, AuthorUserID: ptrFromUUID(row.AuthorUserID), Visibility: row.Visibility, Body: row.Body, Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy)}
}

func mapHRCaseComments(rows []sqlc.HrmsHrCaseComment) []*domain.HRCaseComment {
	items := make([]*domain.HRCaseComment, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapHRCaseComment(row))
	}
	return items
}

func mapHRCaseAttachment(row sqlc.HrmsHrCaseAttachment) *domain.HRCaseAttachment {
	return &domain.HRCaseAttachment{ID: row.ID, TenantID: row.TenantID, CaseID: row.CaseID, CommentID: ptrFromUUID(row.CommentID), FileName: row.FileName, ContentType: row.ContentType, ObjectKey: row.ObjectKey, Visibility: row.Visibility, UploadedBy: ptrFromUUID(row.UploadedBy), Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy)}
}

func mapHRCaseAttachments(rows []sqlc.HrmsHrCaseAttachment) []*domain.HRCaseAttachment {
	items := make([]*domain.HRCaseAttachment, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapHRCaseAttachment(row))
	}
	return items
}

func mapHRCaseEvent(row sqlc.HrmsHrCaseEvent) *domain.HRCaseEvent {
	metadata := json.RawMessage(row.Metadata)
	if len(metadata) == 0 {
		metadata = json.RawMessage(`{}`)
	}
	return &domain.HRCaseEvent{ID: row.ID, TenantID: row.TenantID, CaseID: row.CaseID, EventType: row.EventType, FromStatus: ptrFromText(row.FromStatus), ToStatus: ptrFromText(row.ToStatus), ActorUserID: ptrFromUUID(row.ActorUserID), Comment: ptrFromText(row.Comment), Metadata: metadata, CreatedAt: timeFromTimestamptz(row.CreatedAt)}
}

func mapHRCaseEvents(rows []sqlc.HrmsHrCaseEvent) []*domain.HRCaseEvent {
	items := make([]*domain.HRCaseEvent, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapHRCaseEvent(row))
	}
	return items
}

func mapHRCaseSummary(rows []sqlc.GetHRCaseSummaryRow) []*domain.HRCaseSummaryRow {
	items := make([]*domain.HRCaseSummaryRow, 0, len(rows))
	for _, row := range rows {
		items = append(items, &domain.HRCaseSummaryRow{Status: row.Status, Priority: row.Priority, CaseCount: row.CaseCount, OverdueCount: row.OverdueCount, EscalatedCount: row.EscalatedCount, RestrictedCount: row.RestrictedCount})
	}
	return items
}
