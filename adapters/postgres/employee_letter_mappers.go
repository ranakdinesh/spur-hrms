package postgres

import (
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func mapEmployeeLetterTemplate(row sqlc.HrmsEmployeeLetterTemplate) *domain.EmployeeLetterTemplate {
	return &domain.EmployeeLetterTemplate{ID: row.ID, TenantID: row.TenantID, LetterType: row.LetterType, Name: row.Name, Description: ptrFromText(row.Description), Subject: ptrFromText(row.Subject), BodyHTML: row.BodyHtml, FooterHTML: ptrFromText(row.FooterHtml), Locale: row.Locale, IsDefault: row.IsDefault, IsActive: row.IsActive, Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy)}
}

func mapEmployeeLetterTemplates(rows []sqlc.HrmsEmployeeLetterTemplate) []*domain.EmployeeLetterTemplate {
	items := make([]*domain.EmployeeLetterTemplate, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapEmployeeLetterTemplate(row))
	}
	return items
}

func mapEmployeeLetter(row sqlc.HrmsEmployeeLetter) *domain.EmployeeLetter {
	return &domain.EmployeeLetter{ID: row.ID, TenantID: row.TenantID, EmployeeID: row.EmployeeID, UserID: row.UserID, TemplateID: ptrFromUUID(row.TemplateID), DocumentTypeID: ptrFromUUID(row.DocumentTypeID), EmployeeDocumentID: ptrFromUUID(row.EmployeeDocumentID), LetterType: row.LetterType, Subject: ptrFromText(row.Subject), RenderedHTML: ptrFromText(row.RenderedHtml), Status: row.Status, IssueDate: ptrFromDate(row.IssueDate), EffectiveDate: ptrFromDate(row.EffectiveDate), EndDate: ptrFromDate(row.EndDate), PDFPath: ptrFromText(row.PdfPath), Version: row.Version, IsLatest: row.IsLatest, ApprovalRequestedAt: ptrFromTimestamptz(row.ApprovalRequestedAt), ApprovedAt: ptrFromTimestamptz(row.ApprovedAt), ApprovedBy: ptrFromUUID(row.ApprovedBy), SentAt: ptrFromTimestamptz(row.SentAt), RevokedAt: ptrFromTimestamptz(row.RevokedAt), SignatureToken: ptrFromText(row.SignatureToken), SignatureRequestedAt: ptrFromTimestamptz(row.SignatureRequestedAt), SignatureCompletedAt: ptrFromTimestamptz(row.SignatureCompletedAt), SignerName: ptrFromText(row.SignerName), SignerEmail: ptrFromText(row.SignerEmail), SignerIP: stringPtrFromAddr(row.SignerIp), SignerUserAgent: ptrFromText(row.SignerUserAgent), SignatureHash: ptrFromText(row.SignatureHash), AuditCertificateURL: ptrFromText(row.AuditCertificateUrl), Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy)}
}

func mapEmployeeLetterListRow(row sqlc.ListEmployeeLettersRow) *domain.EmployeeLetter {
	return &domain.EmployeeLetter{ID: row.ID, TenantID: row.TenantID, EmployeeID: row.EmployeeID, UserID: row.UserID, EmployeeCode: ptrFromText(row.EmployeeCode), EmployeeFirstname: &row.EmployeeFirstname, EmployeeLastname: ptrFromText(row.EmployeeLastname), EmployeeEmail: ptrFromText(row.EmployeeEmail), DepartmentName: ptrFromText(row.DepartmentName), BranchName: ptrFromText(row.BranchName), DesignationName: ptrFromText(row.DesignationName), TemplateID: ptrFromUUID(row.TemplateID), TemplateName: ptrFromText(row.TemplateName), DocumentTypeID: ptrFromUUID(row.DocumentTypeID), DocumentTypeName: ptrFromText(row.DocumentTypeName), EmployeeDocumentID: ptrFromUUID(row.EmployeeDocumentID), LetterType: row.LetterType, Subject: ptrFromText(row.Subject), RenderedHTML: ptrFromText(row.RenderedHtml), Status: row.Status, IssueDate: ptrFromDate(row.IssueDate), EffectiveDate: ptrFromDate(row.EffectiveDate), EndDate: ptrFromDate(row.EndDate), PDFPath: ptrFromText(row.PdfPath), Version: row.Version, IsLatest: row.IsLatest, ApprovalRequestedAt: ptrFromTimestamptz(row.ApprovalRequestedAt), ApprovedAt: ptrFromTimestamptz(row.ApprovedAt), ApprovedBy: ptrFromUUID(row.ApprovedBy), SentAt: ptrFromTimestamptz(row.SentAt), RevokedAt: ptrFromTimestamptz(row.RevokedAt), SignatureToken: ptrFromText(row.SignatureToken), SignatureRequestedAt: ptrFromTimestamptz(row.SignatureRequestedAt), SignatureCompletedAt: ptrFromTimestamptz(row.SignatureCompletedAt), SignerName: ptrFromText(row.SignerName), SignerEmail: ptrFromText(row.SignerEmail), SignerIP: stringPtrFromAddr(row.SignerIp), SignerUserAgent: ptrFromText(row.SignerUserAgent), SignatureHash: ptrFromText(row.SignatureHash), AuditCertificateURL: ptrFromText(row.AuditCertificateUrl), Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy)}
}

func mapEmployeeLetters(rows []sqlc.ListEmployeeLettersRow) []*domain.EmployeeLetter {
	items := make([]*domain.EmployeeLetter, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapEmployeeLetterListRow(row))
	}
	return items
}

func mapEmployeeLetterEvent(row sqlc.HrmsEmployeeLetterEvent) *domain.EmployeeLetterEvent {
	return &domain.EmployeeLetterEvent{ID: row.ID, TenantID: row.TenantID, EmployeeLetterID: row.EmployeeLetterID, FromStatus: ptrFromText(row.FromStatus), ToStatus: row.ToStatus, Action: row.Action, Remarks: ptrFromText(row.Remarks), ActorEmail: ptrFromText(row.ActorEmail), IPAddress: stringPtrFromAddr(row.IpAddress), UserAgent: ptrFromText(row.UserAgent), Metadata: mapFromJSONBytes(row.Metadata), Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy)}
}

func mapEmployeeLetterEvents(rows []sqlc.HrmsEmployeeLetterEvent) []*domain.EmployeeLetterEvent {
	items := make([]*domain.EmployeeLetterEvent, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapEmployeeLetterEvent(row))
	}
	return items
}
