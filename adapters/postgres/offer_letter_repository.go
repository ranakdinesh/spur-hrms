package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func (s *Store) CreateOfferLetterTemplate(ctx context.Context, item *domain.OfferLetterTemplate, actorID *uuid.UUID) (*domain.OfferLetterTemplate, error) {
	q := s.getQueries(ctx)
	if item.IsDefault {
		if err := q.ClearDefaultOfferLetterTemplates(ctx, sqlc.ClearDefaultOfferLetterTemplatesParams{TenantID: item.TenantID, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
			return nil, s.logDBError(ctx, "clear default offer templates", err, tenantIDField(item.TenantID))
		}
	}
	row, err := q.CreateOfferLetterTemplate(ctx, sqlc.CreateOfferLetterTemplateParams{TenantID: item.TenantID, Name: item.Name, Description: textFromPtr(item.Description), Subject: textFromPtr(item.Subject), BodyHtml: item.BodyHTML, FooterHtml: textFromPtr(item.FooterHTML), Locale: item.Locale, IsDefault: item.IsDefault, IsActive: item.IsActive, CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create offer template", fmt.Errorf("hrms: create offer template: %w", err), tenantIDField(item.TenantID))
	}
	return mapOfferLetterTemplate(row), nil
}

func (s *Store) ListOfferLetterTemplates(ctx context.Context, tenantID uuid.UUID) ([]*domain.OfferLetterTemplate, error) {
	rows, err := s.getQueries(ctx).ListOfferLetterTemplates(ctx, tenantID)
	if err != nil {
		return nil, s.logDBError(ctx, "list offer templates", err, tenantIDField(tenantID))
	}
	return mapOfferLetterTemplates(rows), nil
}

func (s *Store) GetDefaultOfferLetterTemplate(ctx context.Context, tenantID uuid.UUID) (*domain.OfferLetterTemplate, error) {
	row, err := s.getQueries(ctx).GetDefaultOfferLetterTemplate(ctx, tenantID)
	if err != nil {
		return nil, s.logDBError(ctx, "get default offer template", err, tenantIDField(tenantID))
	}
	return mapOfferLetterTemplate(row), nil
}

func (s *Store) GetOfferLetterTemplate(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.OfferLetterTemplate, error) {
	row, err := s.getQueries(ctx).GetOfferLetterTemplate(ctx, sqlc.GetOfferLetterTemplateParams{TenantID: tenantID, ID: id})
	if err != nil {
		return nil, s.logDBError(ctx, "get offer template", err, tenantIDField(tenantID), stringField("offer_template_id", id.String()))
	}
	return mapOfferLetterTemplate(row), nil
}

func (s *Store) UpdateOfferLetterTemplate(ctx context.Context, item *domain.OfferLetterTemplate, actorID *uuid.UUID) (*domain.OfferLetterTemplate, error) {
	q := s.getQueries(ctx)
	if item.IsDefault {
		if err := q.ClearDefaultOfferLetterTemplates(ctx, sqlc.ClearDefaultOfferLetterTemplatesParams{TenantID: item.TenantID, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
			return nil, s.logDBError(ctx, "clear default offer templates", err, tenantIDField(item.TenantID))
		}
	}
	row, err := q.UpdateOfferLetterTemplate(ctx, sqlc.UpdateOfferLetterTemplateParams{TenantID: item.TenantID, ID: item.ID, Name: item.Name, Description: textFromPtr(item.Description), Subject: textFromPtr(item.Subject), BodyHtml: item.BodyHTML, FooterHtml: textFromPtr(item.FooterHTML), Locale: item.Locale, IsDefault: item.IsDefault, IsActive: item.IsActive, UpdatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "update offer template", err, tenantIDField(item.TenantID), stringField("offer_template_id", item.ID.String()))
	}
	return mapOfferLetterTemplate(row), nil
}

func (s *Store) DeleteOfferLetterTemplate(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteOfferLetterTemplate(ctx, sqlc.SoftDeleteOfferLetterTemplateParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete offer template", err, tenantIDField(tenantID), stringField("offer_template_id", id.String()))
	}
	return nil
}

func (s *Store) CreateOfferLetter(ctx context.Context, item *domain.OfferLetter, actorID *uuid.UUID) (*domain.OfferLetter, error) {
	q := s.getQueries(ctx)
	if err := q.ClearLatestOfferLetters(ctx, sqlc.ClearLatestOfferLettersParams{TenantID: item.TenantID, ApplicationID: item.ApplicationID, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return nil, s.logDBError(ctx, "clear latest offers", err, tenantIDField(item.TenantID), stringField("candidate_application_id", item.ApplicationID.String()))
	}
	version, err := q.NextOfferLetterVersion(ctx, sqlc.NextOfferLetterVersionParams{TenantID: item.TenantID, ApplicationID: item.ApplicationID})
	if err != nil {
		return nil, s.logDBError(ctx, "next offer version", err, tenantIDField(item.TenantID), stringField("candidate_application_id", item.ApplicationID.String()))
	}
	row, err := q.CreateOfferLetter(ctx, sqlc.CreateOfferLetterParams{TenantID: item.TenantID, ApplicationID: item.ApplicationID, CandidateID: uuidFromPtr(item.CandidateID), TemplateID: uuidFromPtr(item.TemplateID), OfferedCtc: numericFromFloatPtr(item.OfferedCTC), Currency: item.Currency, SalaryBreakdown: jsonBytesFromAnyMap(item.SalaryBreakdown), JoiningDate: dateFromPtr(item.JoiningDate), ValidUntilDate: dateFromPtr(item.ValidUntilDate), Status: item.Status, OfferLetterUrl: textFromPtr(item.OfferLetterURL), Version: version, IsLatest: true, Subject: textFromPtr(item.Subject), RenderedHtml: textFromPtr(item.RenderedHTML), SignatureToken: textFromPtr(item.SignatureToken), SignatureRequestedAt: timestamptzFromPtr(item.SignatureRequestedAt), SignerEmail: textFromPtr(item.SignerEmail), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create offer", fmt.Errorf("hrms: create offer: %w", err), tenantIDField(item.TenantID), stringField("candidate_application_id", item.ApplicationID.String()))
	}
	return mapOfferLetter(row), nil
}

func (s *Store) ListOfferLetters(ctx context.Context, filter domain.OfferLetterFilter) ([]*domain.OfferLetter, error) {
	rows, err := s.getQueries(ctx).ListOfferLetters(ctx, sqlc.ListOfferLettersParams{TenantID: filter.TenantID, ApplicationID: uuidFromPtr(filter.ApplicationID), Status: textFromPtr(filter.Status), Search: textFromPtr(filter.Search), Offset: filter.Offset, Limit: filter.Limit})
	if err != nil {
		return nil, s.logDBError(ctx, "list offers", err, tenantIDField(filter.TenantID))
	}
	return mapOfferLetters(rows), nil
}

func (s *Store) CountOfferLetters(ctx context.Context, filter domain.OfferLetterFilter) (int64, error) {
	count, err := s.getQueries(ctx).CountOfferLetters(ctx, sqlc.CountOfferLettersParams{TenantID: filter.TenantID, ApplicationID: uuidFromPtr(filter.ApplicationID), Status: textFromPtr(filter.Status), Search: textFromPtr(filter.Search)})
	if err != nil {
		return 0, s.logDBError(ctx, "count offers", err, tenantIDField(filter.TenantID))
	}
	return count, nil
}

func (s *Store) ListOfferLettersByApplication(ctx context.Context, tenantID uuid.UUID, applicationID uuid.UUID) ([]*domain.OfferLetter, error) {
	rows, err := s.getQueries(ctx).ListOfferLettersByApplication(ctx, sqlc.ListOfferLettersByApplicationParams{TenantID: tenantID, ApplicationID: applicationID})
	if err != nil {
		return nil, s.logDBError(ctx, "list offers by application", err, tenantIDField(tenantID), stringField("candidate_application_id", applicationID.String()))
	}
	return mapOfferLettersBase(rows), nil
}

func (s *Store) GetLatestOfferLetterByApplication(ctx context.Context, tenantID uuid.UUID, applicationID uuid.UUID) (*domain.OfferLetter, error) {
	row, err := s.getQueries(ctx).GetLatestOfferLetterByApplication(ctx, sqlc.GetLatestOfferLetterByApplicationParams{TenantID: tenantID, ApplicationID: applicationID})
	if err != nil {
		return nil, s.logDBError(ctx, "get latest offer", err, tenantIDField(tenantID), stringField("candidate_application_id", applicationID.String()))
	}
	return mapOfferLetter(row), nil
}

func (s *Store) GetOfferLetter(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.OfferLetter, error) {
	row, err := s.getQueries(ctx).GetOfferLetter(ctx, sqlc.GetOfferLetterParams{TenantID: tenantID, ID: id})
	if err != nil {
		return nil, s.logDBError(ctx, "get offer", err, tenantIDField(tenantID), stringField("offer_letter_id", id.String()))
	}
	return mapOfferLetter(row), nil
}

func (s *Store) GetOfferLetterBySignatureToken(ctx context.Context, token string) (*domain.OfferLetter, error) {
	row, err := s.getQueries(ctx).GetOfferLetterBySignatureToken(ctx, textFromPtr(&token))
	if err != nil {
		return nil, s.logDBError(ctx, "get offer by signature token", err)
	}
	return mapOfferLetter(row), nil
}

func (s *Store) UpdateOfferLetter(ctx context.Context, item *domain.OfferLetter, actorID *uuid.UUID) (*domain.OfferLetter, error) {
	row, err := s.getQueries(ctx).UpdateOfferLetter(ctx, sqlc.UpdateOfferLetterParams{TenantID: item.TenantID, ID: item.ID, TemplateID: uuidFromPtr(item.TemplateID), OfferedCtc: numericFromFloatPtr(item.OfferedCTC), Currency: item.Currency, SalaryBreakdown: jsonBytesFromAnyMap(item.SalaryBreakdown), JoiningDate: dateFromPtr(item.JoiningDate), ValidUntilDate: dateFromPtr(item.ValidUntilDate), Status: item.Status, OfferLetterUrl: textFromPtr(item.OfferLetterURL), Subject: textFromPtr(item.Subject), RenderedHtml: textFromPtr(item.RenderedHTML), SignerEmail: textFromPtr(item.SignerEmail), UpdatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "update offer", err, tenantIDField(item.TenantID), stringField("offer_letter_id", item.ID.String()))
	}
	return mapOfferLetter(row), nil
}

func (s *Store) UpdateOfferLetterStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, reason *string, actorID *uuid.UUID) (*domain.OfferLetter, error) {
	row, err := s.getQueries(ctx).UpdateOfferLetterStatus(ctx, sqlc.UpdateOfferLetterStatusParams{TenantID: tenantID, ID: id, Status: status, CandidateRejectionReason: textFromPtr(reason), UpdatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "update offer status", err, tenantIDField(tenantID), stringField("offer_letter_id", id.String()), stringField("status", status))
	}
	return mapOfferLetter(row), nil
}

func (s *Store) SignOfferLetter(ctx context.Context, token string, signerName string, signerEmail string, signerIP *string, userAgent *string, signatureHash string) (*domain.OfferLetter, error) {
	row, err := s.getQueries(ctx).SignOfferLetter(ctx, sqlc.SignOfferLetterParams{SignatureToken: textFromPtr(&token), SignerName: textFromPtr(&signerName), SignerEmail: textFromPtr(&signerEmail), SignerIp: addrFromStringPtr(signerIP), SignerUserAgent: textFromPtr(userAgent), SignatureHash: textFromPtr(&signatureHash)})
	if err != nil {
		return nil, s.logDBError(ctx, "sign offer", err)
	}
	return mapOfferLetter(row), nil
}

func (s *Store) DeleteOfferLetter(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteOfferLetter(ctx, sqlc.SoftDeleteOfferLetterParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete offer", err, tenantIDField(tenantID), stringField("offer_letter_id", id.String()))
	}
	return nil
}

func (s *Store) CreateOfferLetterEvent(ctx context.Context, event *domain.OfferLetterEvent, actorID *uuid.UUID) (*domain.OfferLetterEvent, error) {
	row, err := s.getQueries(ctx).CreateOfferLetterEvent(ctx, sqlc.CreateOfferLetterEventParams{TenantID: event.TenantID, OfferLetterID: event.OfferLetterID, FromStatus: textFromPtr(event.FromStatus), ToStatus: event.ToStatus, Action: event.Action, Remarks: textFromPtr(event.Remarks), ActorEmail: textFromPtr(event.ActorEmail), IpAddress: addrFromStringPtr(event.IPAddress), UserAgent: textFromPtr(event.UserAgent), Column10: jsonBytesFromAnyMap(event.Metadata), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create offer event", err, tenantIDField(event.TenantID), stringField("offer_letter_id", event.OfferLetterID.String()))
	}
	return mapOfferLetterEvent(row), nil
}

func (s *Store) ListOfferLetterEvents(ctx context.Context, tenantID uuid.UUID, offerLetterID uuid.UUID) ([]*domain.OfferLetterEvent, error) {
	rows, err := s.getQueries(ctx).ListOfferLetterEvents(ctx, sqlc.ListOfferLetterEventsParams{TenantID: tenantID, OfferLetterID: offerLetterID})
	if err != nil {
		return nil, s.logDBError(ctx, "list offer events", err, tenantIDField(tenantID), stringField("offer_letter_id", offerLetterID.String()))
	}
	return mapOfferLetterEvents(rows), nil
}
