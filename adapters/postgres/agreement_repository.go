package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func (s *Store) CreateAgreementTemplate(ctx context.Context, item *domain.AgreementTemplate, actorID *uuid.UUID) (*domain.AgreementTemplate, error) {
	row, err := s.getQueries(ctx).CreateAgreementTemplate(ctx, sqlc.CreateAgreementTemplateParams{TenantID: item.TenantID, AgreementType: item.AgreementType, Name: item.Name, Description: textFromPtr(item.Description), Subject: textFromPtr(item.Subject), BodyHtml: item.BodyHTML, FooterHtml: textFromPtr(item.FooterHTML), Locale: item.Locale, IsDefault: item.IsDefault, IsActive: item.IsActive, Metadata: jsonBytesFromRaw(item.Metadata), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create agreement template", err, tenantIDField(item.TenantID), stringField("agreement_type", item.AgreementType))
	}
	return mapAgreementTemplate(row), nil
}

func (s *Store) ListAgreementTemplates(ctx context.Context, tenantID uuid.UUID, agreementType *string) ([]*domain.AgreementTemplate, error) {
	rows, err := s.getQueries(ctx).ListAgreementTemplates(ctx, sqlc.ListAgreementTemplatesParams{TenantID: tenantID, AgreementType: textFromPtr(agreementType)})
	if err != nil {
		return nil, s.logDBError(ctx, "list agreement templates", err, tenantIDField(tenantID))
	}
	return mapAgreementTemplates(rows), nil
}

func (s *Store) GetDefaultAgreementTemplate(ctx context.Context, tenantID uuid.UUID, agreementType string) (*domain.AgreementTemplate, error) {
	row, err := s.getQueries(ctx).GetDefaultAgreementTemplate(ctx, sqlc.GetDefaultAgreementTemplateParams{TenantID: tenantID, AgreementType: agreementType})
	if err != nil {
		return nil, s.logDBError(ctx, "get default agreement template", err, tenantIDField(tenantID), stringField("agreement_type", agreementType))
	}
	return mapAgreementTemplate(row), nil
}

func (s *Store) GetAgreementTemplate(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.AgreementTemplate, error) {
	row, err := s.getQueries(ctx).GetAgreementTemplate(ctx, sqlc.GetAgreementTemplateParams{TenantID: tenantID, ID: id})
	if err != nil {
		return nil, s.logDBError(ctx, "get agreement template", err, tenantIDField(tenantID), stringField("agreement_template_id", id.String()))
	}
	return mapAgreementTemplate(row), nil
}

func (s *Store) UpdateAgreementTemplate(ctx context.Context, item *domain.AgreementTemplate, actorID *uuid.UUID) (*domain.AgreementTemplate, error) {
	row, err := s.getQueries(ctx).UpdateAgreementTemplate(ctx, sqlc.UpdateAgreementTemplateParams{TenantID: item.TenantID, ID: item.ID, AgreementType: item.AgreementType, Name: item.Name, Description: textFromPtr(item.Description), Subject: textFromPtr(item.Subject), BodyHtml: item.BodyHTML, FooterHtml: textFromPtr(item.FooterHTML), Locale: item.Locale, IsDefault: item.IsDefault, IsActive: item.IsActive, Metadata: jsonBytesFromRaw(item.Metadata), UpdatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "update agreement template", err, tenantIDField(item.TenantID), stringField("agreement_template_id", item.ID.String()))
	}
	return mapAgreementTemplate(row), nil
}

func (s *Store) DeleteAgreementTemplate(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteAgreementTemplate(ctx, sqlc.SoftDeleteAgreementTemplateParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete agreement template", err, tenantIDField(tenantID), stringField("agreement_template_id", id.String()))
	}
	return nil
}

func (s *Store) CreateAgreement(ctx context.Context, item *domain.Agreement, actorID *uuid.UUID) (*domain.Agreement, error) {
	row, err := s.getQueries(ctx).CreateAgreement(ctx, sqlc.CreateAgreementParams{TenantID: item.TenantID, AgreementType: item.AgreementType, Title: item.Title, TemplateID: uuidFromPtr(item.TemplateID), WorkerProfileID: uuidFromPtr(item.WorkerProfileID), EngagementID: uuidFromPtr(item.EngagementID), ProjectID: uuidFromPtr(item.ProjectID), Subject: textFromPtr(item.Subject), RenderedHtml: textFromPtr(item.RenderedHTML), Status: item.Status, IssueDate: dateFromPtr(item.IssueDate), EffectiveDate: dateFromPtr(item.EffectiveDate), EndDate: dateFromPtr(item.EndDate), PdfPath: textFromPtr(item.PDFPath), SignatureToken: textFromPtr(item.SignatureToken), SignatureRequestedAt: timestamptzFromPtr(item.SignatureRequestedAt), SignerEmail: textFromPtr(item.SignerEmail), Metadata: jsonBytesFromRaw(item.Metadata), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create agreement", err, tenantIDField(item.TenantID), stringField("agreement_type", item.AgreementType))
	}
	return mapAgreement(row), nil
}

func (s *Store) ListAgreements(ctx context.Context, filter domain.AgreementFilter) ([]*domain.Agreement, error) {
	rows, err := s.getQueries(ctx).ListAgreements(ctx, sqlc.ListAgreementsParams{TenantID: filter.TenantID, AgreementType: textFromPtr(filter.AgreementType), Status: textFromPtr(filter.Status), WorkerProfileID: uuidFromPtr(filter.WorkerProfileID), EngagementID: uuidFromPtr(filter.EngagementID), ProjectID: uuidFromPtr(filter.ProjectID), Search: textFromPtr(filter.Search)})
	if err != nil {
		return nil, s.logDBError(ctx, "list agreements", err, tenantIDField(filter.TenantID))
	}
	return mapAgreementListItems(rows), nil
}

func (s *Store) GetAgreement(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.Agreement, error) {
	row, err := s.getQueries(ctx).GetAgreement(ctx, sqlc.GetAgreementParams{TenantID: tenantID, ID: id})
	if err != nil {
		return nil, s.logDBError(ctx, "get agreement", err, tenantIDField(tenantID), stringField("agreement_id", id.String()))
	}
	return mapAgreement(row), nil
}

func (s *Store) GetAgreementBySignatureToken(ctx context.Context, token string) (*domain.Agreement, error) {
	row, err := s.getQueries(ctx).GetAgreementBySignatureToken(ctx, textFromPtr(&token))
	if err != nil {
		return nil, s.logDBError(ctx, "get agreement by signature token", err)
	}
	return mapAgreement(row), nil
}

func (s *Store) UpdateAgreementPDF(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, pdfPath *string, actorID *uuid.UUID) (*domain.Agreement, error) {
	row, err := s.getQueries(ctx).UpdateAgreementPDF(ctx, sqlc.UpdateAgreementPDFParams{TenantID: tenantID, ID: id, PdfPath: textFromPtr(pdfPath), UpdatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "update agreement pdf", err, tenantIDField(tenantID), stringField("agreement_id", id.String()))
	}
	return mapAgreement(row), nil
}

func (s *Store) UpdateAgreementStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, actorID *uuid.UUID) (*domain.Agreement, error) {
	row, err := s.getQueries(ctx).UpdateAgreementStatus(ctx, sqlc.UpdateAgreementStatusParams{TenantID: tenantID, ID: id, Status: status, UpdatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "update agreement status", err, tenantIDField(tenantID), stringField("agreement_id", id.String()))
	}
	return mapAgreement(row), nil
}

func (s *Store) SignAgreement(ctx context.Context, token string, signerName string, signerEmail string, signerIP *string, userAgent *string, signatureHash string) (*domain.Agreement, error) {
	row, err := s.getQueries(ctx).SignAgreement(ctx, sqlc.SignAgreementParams{SignatureToken: textFromPtr(&token), SignerName: textFromPtr(&signerName), SignerEmail: textFromPtr(&signerEmail), SignerIp: addrFromStringPtr(signerIP), SignerUserAgent: textFromPtr(userAgent), SignatureHash: textFromPtr(&signatureHash)})
	if err != nil {
		return nil, s.logDBError(ctx, "sign agreement", err)
	}
	return mapAgreement(row), nil
}

func (s *Store) DeleteAgreement(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteAgreement(ctx, sqlc.SoftDeleteAgreementParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete agreement", err, tenantIDField(tenantID), stringField("agreement_id", id.String()))
	}
	return nil
}

func (s *Store) CreateAgreementEvent(ctx context.Context, event *domain.AgreementEvent, actorID *uuid.UUID) (*domain.AgreementEvent, error) {
	row, err := s.getQueries(ctx).CreateAgreementEvent(ctx, sqlc.CreateAgreementEventParams{TenantID: event.TenantID, AgreementID: event.AgreementID, FromStatus: textFromPtr(event.FromStatus), ToStatus: event.ToStatus, Action: event.Action, Remarks: textFromPtr(event.Remarks), ActorEmail: textFromPtr(event.ActorEmail), IpAddress: addrFromStringPtr(event.IPAddress), UserAgent: textFromPtr(event.UserAgent), Metadata: jsonBytesFromAnyMap(event.Metadata), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create agreement event", err, tenantIDField(event.TenantID), stringField("agreement_id", event.AgreementID.String()))
	}
	return mapAgreementEvent(row), nil
}

func (s *Store) ListAgreementEvents(ctx context.Context, tenantID uuid.UUID, agreementID uuid.UUID) ([]*domain.AgreementEvent, error) {
	rows, err := s.getQueries(ctx).ListAgreementEvents(ctx, sqlc.ListAgreementEventsParams{TenantID: tenantID, AgreementID: agreementID})
	if err != nil {
		return nil, s.logDBError(ctx, "list agreement events", err, tenantIDField(tenantID), stringField("agreement_id", agreementID.String()))
	}
	return mapAgreementEvents(rows), nil
}
