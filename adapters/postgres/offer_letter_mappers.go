package postgres

import (
	"encoding/json"
	"net/netip"

	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func mapOfferLetterTemplate(row sqlc.HrmsOfferLetterTemplate) *domain.OfferLetterTemplate {
	return &domain.OfferLetterTemplate{ID: row.ID, TenantID: row.TenantID, Name: row.Name, Description: ptrFromText(row.Description), Subject: ptrFromText(row.Subject), BodyHTML: row.BodyHtml, FooterHTML: ptrFromText(row.FooterHtml), Locale: row.Locale, IsDefault: row.IsDefault, IsActive: row.IsActive, Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy)}
}

func mapOfferLetterTemplates(rows []sqlc.HrmsOfferLetterTemplate) []*domain.OfferLetterTemplate {
	items := make([]*domain.OfferLetterTemplate, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapOfferLetterTemplate(row))
	}
	return items
}

func mapOfferLetter(row sqlc.HrmsOfferLetter) *domain.OfferLetter {
	return &domain.OfferLetter{ID: row.ID, TenantID: row.TenantID, ApplicationID: row.ApplicationID, CandidateID: ptrFromUUID(row.CandidateID), TemplateID: ptrFromUUID(row.TemplateID), OfferedCTC: floatPtrFromNumeric(row.OfferedCtc), Currency: row.Currency, SalaryBreakdown: mapFromJSONBytes(row.SalaryBreakdown), JoiningDate: ptrFromDate(row.JoiningDate), ValidUntilDate: ptrFromDate(row.ValidUntilDate), Status: row.Status, OfferLetterURL: ptrFromText(row.OfferLetterUrl), CandidateReactionDate: ptrFromTimestamptz(row.CandidateReactionDate), CandidateRejectionReason: ptrFromText(row.CandidateRejectionReason), Version: row.Version, IsLatest: row.IsLatest, Subject: ptrFromText(row.Subject), RenderedHTML: ptrFromText(row.RenderedHtml), SentAt: ptrFromTimestamptz(row.SentAt), RevokedAt: ptrFromTimestamptz(row.RevokedAt), SignatureToken: ptrFromText(row.SignatureToken), SignatureRequestedAt: ptrFromTimestamptz(row.SignatureRequestedAt), SignatureCompletedAt: ptrFromTimestamptz(row.SignatureCompletedAt), SignerName: ptrFromText(row.SignerName), SignerEmail: ptrFromText(row.SignerEmail), SignerIP: stringPtrFromAddr(row.SignerIp), SignerUserAgent: ptrFromText(row.SignerUserAgent), SignatureHash: ptrFromText(row.SignatureHash), AuditCertificateURL: ptrFromText(row.AuditCertificateUrl), Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy)}
}

func mapOfferLetterListRow(row sqlc.ListOfferLettersRow) *domain.OfferLetter {
	item := &domain.OfferLetter{ID: row.ID, TenantID: row.TenantID, ApplicationID: row.ApplicationID, CandidateID: ptrFromUUID(row.CandidateID), CandidateFirstname: ptrFromText(row.CandidateFirstname), CandidateLastname: ptrFromText(row.CandidateLastname), CandidateEmail: ptrFromText(row.CandidateEmail), JobPostingTitle: ptrFromText(row.JobPostingTitle), JobPostingCode: ptrFromText(row.JobPostingCode), TemplateID: ptrFromUUID(row.TemplateID), TemplateName: ptrFromText(row.TemplateName), OfferedCTC: floatPtrFromNumeric(row.OfferedCtc), Currency: row.Currency, SalaryBreakdown: mapFromJSONBytes(row.SalaryBreakdown), JoiningDate: ptrFromDate(row.JoiningDate), ValidUntilDate: ptrFromDate(row.ValidUntilDate), Status: row.Status, OfferLetterURL: ptrFromText(row.OfferLetterUrl), CandidateReactionDate: ptrFromTimestamptz(row.CandidateReactionDate), CandidateRejectionReason: ptrFromText(row.CandidateRejectionReason), Version: row.Version, IsLatest: row.IsLatest, Subject: ptrFromText(row.Subject), RenderedHTML: ptrFromText(row.RenderedHtml), SentAt: ptrFromTimestamptz(row.SentAt), RevokedAt: ptrFromTimestamptz(row.RevokedAt), SignatureToken: ptrFromText(row.SignatureToken), SignatureRequestedAt: ptrFromTimestamptz(row.SignatureRequestedAt), SignatureCompletedAt: ptrFromTimestamptz(row.SignatureCompletedAt), SignerName: ptrFromText(row.SignerName), SignerEmail: ptrFromText(row.SignerEmail), SignerIP: stringPtrFromAddr(row.SignerIp), SignerUserAgent: ptrFromText(row.SignerUserAgent), SignatureHash: ptrFromText(row.SignatureHash), AuditCertificateURL: ptrFromText(row.AuditCertificateUrl), Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy)}
	return item
}

func mapOfferLetters(rows []sqlc.ListOfferLettersRow) []*domain.OfferLetter {
	items := make([]*domain.OfferLetter, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapOfferLetterListRow(row))
	}
	return items
}

func mapOfferLettersBase(rows []sqlc.HrmsOfferLetter) []*domain.OfferLetter {
	items := make([]*domain.OfferLetter, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapOfferLetter(row))
	}
	return items
}

func mapOfferLetterEvent(row sqlc.HrmsOfferLetterEvent) *domain.OfferLetterEvent {
	return &domain.OfferLetterEvent{ID: row.ID, TenantID: row.TenantID, OfferLetterID: row.OfferLetterID, FromStatus: ptrFromText(row.FromStatus), ToStatus: row.ToStatus, Action: row.Action, Remarks: ptrFromText(row.Remarks), ActorEmail: ptrFromText(row.ActorEmail), IPAddress: stringPtrFromAddr(row.IpAddress), UserAgent: ptrFromText(row.UserAgent), Metadata: mapFromJSONBytes(row.Metadata), Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy)}
}

func mapOfferLetterEvents(rows []sqlc.HrmsOfferLetterEvent) []*domain.OfferLetterEvent {
	items := make([]*domain.OfferLetterEvent, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapOfferLetterEvent(row))
	}
	return items
}

func mapFromJSONBytes(data []byte) map[string]any {
	if len(data) == 0 {
		return map[string]any{}
	}
	var value map[string]any
	if err := json.Unmarshal(data, &value); err != nil || value == nil {
		return map[string]any{}
	}
	return value
}

func jsonBytesFromAnyMap(value map[string]any) []byte {
	if value == nil {
		return []byte("{}")
	}
	data, err := json.Marshal(value)
	if err != nil {
		return []byte("{}")
	}
	return data
}

func stringPtrFromAddr(value *netip.Addr) *string {
	if value == nil {
		return nil
	}
	text := value.String()
	return &text
}

func addrFromStringPtr(value *string) *netip.Addr {
	if value == nil || *value == "" {
		return nil
	}
	addr, err := netip.ParseAddr(*value)
	if err != nil {
		return nil
	}
	return &addr
}
