package postgres

import (
	"encoding/json"

	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func mapAgreementTemplate(row sqlc.HrmsAgreementTemplate) *domain.AgreementTemplate {
	metadata := json.RawMessage(row.Metadata)
	if len(metadata) == 0 {
		metadata = json.RawMessage(`{}`)
	}
	return &domain.AgreementTemplate{ID: row.ID, TenantID: row.TenantID, AgreementType: row.AgreementType, Name: row.Name, Description: ptrFromText(row.Description), Subject: ptrFromText(row.Subject), BodyHTML: row.BodyHtml, FooterHTML: ptrFromText(row.FooterHtml), Locale: row.Locale, IsDefault: row.IsDefault, IsActive: row.IsActive, Metadata: metadata, Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy)}
}

func mapAgreementTemplates(rows []sqlc.HrmsAgreementTemplate) []*domain.AgreementTemplate {
	items := make([]*domain.AgreementTemplate, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapAgreementTemplate(row))
	}
	return items
}

func mapAgreement(row sqlc.HrmsAgreement) *domain.Agreement {
	metadata := json.RawMessage(row.Metadata)
	if len(metadata) == 0 {
		metadata = json.RawMessage(`{}`)
	}
	return &domain.Agreement{ID: row.ID, TenantID: row.TenantID, AgreementType: row.AgreementType, Title: row.Title, TemplateID: ptrFromUUID(row.TemplateID), WorkerProfileID: ptrFromUUID(row.WorkerProfileID), EngagementID: ptrFromUUID(row.EngagementID), ProjectID: ptrFromUUID(row.ProjectID), Subject: ptrFromText(row.Subject), RenderedHTML: ptrFromText(row.RenderedHtml), Status: row.Status, IssueDate: ptrFromDate(row.IssueDate), EffectiveDate: ptrFromDate(row.EffectiveDate), EndDate: ptrFromDate(row.EndDate), PDFPath: ptrFromText(row.PdfPath), Version: row.Version, IsLatest: row.IsLatest, SentAt: ptrFromTimestamptz(row.SentAt), RevokedAt: ptrFromTimestamptz(row.RevokedAt), SignatureToken: ptrFromText(row.SignatureToken), SignatureRequestedAt: ptrFromTimestamptz(row.SignatureRequestedAt), SignatureCompletedAt: ptrFromTimestamptz(row.SignatureCompletedAt), SignerName: ptrFromText(row.SignerName), SignerEmail: ptrFromText(row.SignerEmail), SignerIP: stringPtrFromAddr(row.SignerIp), SignerUserAgent: ptrFromText(row.SignerUserAgent), SignatureHash: ptrFromText(row.SignatureHash), AuditCertificateURL: ptrFromText(row.AuditCertificateUrl), Metadata: metadata, Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy)}
}

func mapAgreementListItems(rows []sqlc.ListAgreementsRow) []*domain.Agreement {
	items := make([]*domain.Agreement, 0, len(rows))
	for _, row := range rows {
		metadata := json.RawMessage(row.Metadata)
		if len(metadata) == 0 {
			metadata = json.RawMessage(`{}`)
		}
		items = append(items, &domain.Agreement{ID: row.ID, TenantID: row.TenantID, AgreementType: row.AgreementType, Title: row.Title, TemplateID: ptrFromUUID(row.TemplateID), TemplateName: ptrFromText(row.TemplateName), WorkerProfileID: ptrFromUUID(row.WorkerProfileID), WorkerDisplayName: ptrFromText(row.WorkerDisplayName), WorkerCode: ptrFromText(row.WorkerCode), EngagementID: ptrFromUUID(row.EngagementID), EngagementTitle: ptrFromText(row.EngagementTitle), EngagementCode: ptrFromText(row.EngagementCode), ProjectID: ptrFromUUID(row.ProjectID), ProjectName: ptrFromText(row.ProjectName), ProjectCode: ptrFromText(row.ProjectCode), Subject: ptrFromText(row.Subject), RenderedHTML: ptrFromText(row.RenderedHtml), Status: row.Status, IssueDate: ptrFromDate(row.IssueDate), EffectiveDate: ptrFromDate(row.EffectiveDate), EndDate: ptrFromDate(row.EndDate), PDFPath: ptrFromText(row.PdfPath), Version: row.Version, IsLatest: row.IsLatest, SentAt: ptrFromTimestamptz(row.SentAt), RevokedAt: ptrFromTimestamptz(row.RevokedAt), SignatureToken: ptrFromText(row.SignatureToken), SignatureRequestedAt: ptrFromTimestamptz(row.SignatureRequestedAt), SignatureCompletedAt: ptrFromTimestamptz(row.SignatureCompletedAt), SignerName: ptrFromText(row.SignerName), SignerEmail: ptrFromText(row.SignerEmail), SignerIP: stringPtrFromAddr(row.SignerIp), SignerUserAgent: ptrFromText(row.SignerUserAgent), SignatureHash: ptrFromText(row.SignatureHash), AuditCertificateURL: ptrFromText(row.AuditCertificateUrl), Metadata: metadata, Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy)})
	}
	return items
}

func mapAgreementEvent(row sqlc.HrmsAgreementEvent) *domain.AgreementEvent {
	return &domain.AgreementEvent{ID: row.ID, TenantID: row.TenantID, AgreementID: row.AgreementID, FromStatus: ptrFromText(row.FromStatus), ToStatus: row.ToStatus, Action: row.Action, Remarks: ptrFromText(row.Remarks), ActorEmail: ptrFromText(row.ActorEmail), IPAddress: stringPtrFromAddr(row.IpAddress), UserAgent: ptrFromText(row.UserAgent), Metadata: mapFromJSONBytes(row.Metadata), Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy)}
}

func mapAgreementEvents(rows []sqlc.HrmsAgreementEvent) []*domain.AgreementEvent {
	items := make([]*domain.AgreementEvent, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapAgreementEvent(row))
	}
	return items
}
