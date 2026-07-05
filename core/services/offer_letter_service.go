package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) CreateOfferLetterTemplate(ctx context.Context, cmd ports.OfferLetterTemplateCommand) (*domain.OfferLetterTemplate, error) {
	item, err := domain.NewOfferLetterTemplate(offerTemplateInput(cmd))
	if err != nil {
		s.logError("validate offer template create", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	result, err := s.offerLetters.CreateOfferLetterTemplate(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create offer template", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) ListOfferLetterTemplates(ctx context.Context, tenantID uuid.UUID) ([]*domain.OfferLetterTemplate, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate offer template list tenant", err)
		return nil, err
	}
	items, err := s.offerLetters.ListOfferLetterTemplates(ctx, tenantID)
	if err != nil {
		s.logError("list offer templates", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	return items, nil
}

func (s *TenantService) GetOfferLetterTemplate(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.OfferLetterTemplate, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate offer template get tenant", err)
		return nil, err
	}
	if id == uuid.Nil {
		err := domain.ErrInvalidOfferTemplateID
		s.logError("validate offer template get id", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	item, err := s.offerLetters.GetOfferLetterTemplate(ctx, tenantID, id)
	if err != nil {
		s.logError("get offer template", err, serviceTenantIDField(tenantID), serviceStringField("offer_template_id", id.String()))
		return nil, err
	}
	return item, nil
}

func (s *TenantService) UpdateOfferLetterTemplate(ctx context.Context, cmd ports.OfferLetterTemplateCommand) (*domain.OfferLetterTemplate, error) {
	if cmd.ID == uuid.Nil {
		err := domain.ErrInvalidOfferTemplateID
		s.logError("validate offer template update id", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	if _, err := s.GetOfferLetterTemplate(ctx, cmd.TenantID, cmd.ID); err != nil {
		return nil, err
	}
	item, err := domain.NewOfferLetterTemplate(offerTemplateInput(cmd))
	if err != nil {
		s.logError("validate offer template update", err, serviceTenantIDField(cmd.TenantID), serviceStringField("offer_template_id", cmd.ID.String()))
		return nil, err
	}
	item.ID = cmd.ID
	result, err := s.offerLetters.UpdateOfferLetterTemplate(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("update offer template", err, serviceTenantIDField(cmd.TenantID), serviceStringField("offer_template_id", cmd.ID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) DeleteOfferLetterTemplate(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if _, err := s.GetOfferLetterTemplate(ctx, tenantID, id); err != nil {
		return err
	}
	if err := s.offerLetters.DeleteOfferLetterTemplate(ctx, tenantID, id, actorID); err != nil {
		s.logError("delete offer template", err, serviceTenantIDField(tenantID), serviceStringField("offer_template_id", id.String()))
		return err
	}
	return nil
}

func (s *TenantService) CreateOfferLetter(ctx context.Context, cmd ports.OfferLetterCommand) (*domain.OfferLetter, error) {
	app, err := s.GetCandidateApplication(ctx, cmd.TenantID, cmd.ApplicationID)
	if err != nil {
		return nil, err
	}
	if cmd.TemplateID == nil {
		if tmpl, err := s.offerLetters.GetDefaultOfferLetterTemplate(ctx, cmd.TenantID); err == nil && tmpl != nil {
			cmd.TemplateID = &tmpl.ID
		}
	}
	item, err := domain.NewOfferLetter(offerLetterInput(cmd, app.CandidateID))
	if err != nil {
		s.logError("validate offer create", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	if item.TemplateID != nil {
		if tmpl, err := s.GetOfferLetterTemplate(ctx, cmd.TenantID, *item.TemplateID); err == nil {
			renderOfferLetter(item, tmpl, app)
		} else {
			return nil, err
		}
	}
	token, err := domain.NewOfferSignatureToken()
	if err != nil {
		s.logError("generate offer signature token", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	item.SignatureToken = &token
	result, err := s.offerLetters.CreateOfferLetter(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create offer", err, serviceTenantIDField(cmd.TenantID), serviceStringField("candidate_application_id", cmd.ApplicationID.String()))
		return nil, err
	}
	_, _ = s.offerLetters.CreateOfferLetterEvent(ctx, &domain.OfferLetterEvent{TenantID: result.TenantID, OfferLetterID: result.ID, ToStatus: result.Status, Action: "generated", Remarks: result.Subject, Metadata: map[string]any{"version": result.Version}}, cmd.ActorID)
	return result, nil
}

func (s *TenantService) ListOfferLetters(ctx context.Context, filter domain.OfferLetterFilter) (*domain.OfferLetterPage, error) {
	if filter.TenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate offer list tenant", err)
		return nil, err
	}
	if filter.Status != nil {
		if _, err := domain.ValidateOfferStatus(*filter.Status); err != nil {
			s.logError("validate offer list status", err, serviceTenantIDField(filter.TenantID))
			return nil, err
		}
	}
	filter.Search = cleanStringPtr(filter.Search)
	filter.Status = cleanStringPtr(filter.Status)
	limit, offset := normalizeListWindow(filter.Limit, filter.Offset)
	filter.Limit = limit
	filter.Offset = offset
	items, err := s.offerLetters.ListOfferLetters(ctx, filter)
	if err != nil {
		s.logError("list offers", err, serviceTenantIDField(filter.TenantID))
		return nil, err
	}
	total, err := s.offerLetters.CountOfferLetters(ctx, filter)
	if err != nil {
		s.logError("count offers", err, serviceTenantIDField(filter.TenantID))
		return nil, err
	}
	page := &domain.OfferLetterPage{Items: items, Total: total, Limit: limit, Offset: offset}
	if int64(offset)+int64(len(items)) < total {
		next := offset + limit
		page.NextOffset = &next
	}
	return page, nil
}

func (s *TenantService) ListOfferLettersByApplication(ctx context.Context, tenantID uuid.UUID, applicationID uuid.UUID) ([]*domain.OfferLetter, error) {
	if _, err := s.GetCandidateApplication(ctx, tenantID, applicationID); err != nil {
		return nil, err
	}
	return s.offerLetters.ListOfferLettersByApplication(ctx, tenantID, applicationID)
}

func (s *TenantService) GetLatestOfferLetterByApplication(ctx context.Context, tenantID uuid.UUID, applicationID uuid.UUID) (*domain.OfferLetter, error) {
	if _, err := s.GetCandidateApplication(ctx, tenantID, applicationID); err != nil {
		return nil, err
	}
	return s.offerLetters.GetLatestOfferLetterByApplication(ctx, tenantID, applicationID)
}

func (s *TenantService) GetOfferLetter(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.OfferLetter, error) {
	if tenantID == uuid.Nil || id == uuid.Nil {
		err := domain.ErrInvalidOfferLetterID
		s.logError("validate offer get", err)
		return nil, err
	}
	item, err := s.offerLetters.GetOfferLetter(ctx, tenantID, id)
	if err != nil {
		s.logError("get offer", err, serviceTenantIDField(tenantID), serviceStringField("offer_letter_id", id.String()))
		return nil, err
	}
	return item, nil
}

func (s *TenantService) GetOfferLetterBySignatureToken(ctx context.Context, token string) (*domain.OfferLetter, error) {
	if strings.TrimSpace(token) == "" {
		return nil, domain.ErrOfferSignatureTokenMissing
	}
	return s.offerLetters.GetOfferLetterBySignatureToken(ctx, token)
}

func (s *TenantService) UpdateOfferLetter(ctx context.Context, cmd ports.OfferLetterCommand) (*domain.OfferLetter, error) {
	existing, err := s.GetOfferLetter(ctx, cmd.TenantID, cmd.ID)
	if err != nil {
		return nil, err
	}
	if existing.Status != domain.OfferStatusGenerated {
		return nil, domain.ErrOfferCannotBeChanged
	}
	app, err := s.GetCandidateApplication(ctx, cmd.TenantID, cmd.ApplicationID)
	if err != nil {
		return nil, err
	}
	item, err := domain.NewOfferLetter(offerLetterInput(cmd, app.CandidateID))
	if err != nil {
		s.logError("validate offer update", err, serviceTenantIDField(cmd.TenantID), serviceStringField("offer_letter_id", cmd.ID.String()))
		return nil, err
	}
	item.ID = cmd.ID
	item.Version = existing.Version
	item.IsLatest = existing.IsLatest
	item.SignatureToken = existing.SignatureToken
	if item.TemplateID != nil {
		tmpl, err := s.GetOfferLetterTemplate(ctx, cmd.TenantID, *item.TemplateID)
		if err != nil {
			return nil, err
		}
		renderOfferLetter(item, tmpl, app)
	}
	result, err := s.offerLetters.UpdateOfferLetter(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("update offer", err, serviceTenantIDField(cmd.TenantID), serviceStringField("offer_letter_id", cmd.ID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) UpdateOfferLetterStatus(ctx context.Context, cmd ports.OfferLetterStatusCommand) (*domain.OfferLetter, error) {
	existing, err := s.GetOfferLetter(ctx, cmd.TenantID, cmd.ID)
	if err != nil {
		return nil, err
	}
	status, err := domain.ValidateOfferStatus(cmd.Status)
	if err != nil {
		s.logError("validate offer status", err, serviceTenantIDField(cmd.TenantID), serviceStringField("offer_letter_id", cmd.ID.String()))
		return nil, err
	}
	result, err := s.offerLetters.UpdateOfferLetterStatus(ctx, cmd.TenantID, cmd.ID, status, cmd.Reason, cmd.ActorID)
	if err != nil {
		s.logError("update offer status", err, serviceTenantIDField(cmd.TenantID), serviceStringField("offer_letter_id", cmd.ID.String()), serviceStringField("status", status))
		return nil, err
	}
	_, _ = s.offerLetters.CreateOfferLetterEvent(ctx, &domain.OfferLetterEvent{TenantID: result.TenantID, OfferLetterID: result.ID, FromStatus: &existing.Status, ToStatus: result.Status, Action: strings.ToLower(result.Status), Remarks: cmd.Reason, Metadata: map[string]any{"version": result.Version}}, cmd.ActorID)
	return result, nil
}

func (s *TenantService) SignOfferLetter(ctx context.Context, cmd ports.OfferLetterSignatureCommand) (*domain.OfferLetter, error) {
	if strings.TrimSpace(cmd.Token) == "" || strings.TrimSpace(cmd.SignerName) == "" || strings.TrimSpace(cmd.SignerEmail) == "" {
		err := domain.ErrInvalidOfferSignature
		s.logError("validate offer signature", err)
		return nil, err
	}
	existing, err := s.offerLetters.GetOfferLetterBySignatureToken(ctx, cmd.Token)
	if err != nil {
		s.logError("get offer by signature token", err)
		return nil, err
	}
	hash := signatureHash(existing, cmd)
	result, err := s.offerLetters.SignOfferLetter(ctx, cmd.Token, cmd.SignerName, cmd.SignerEmail, cmd.IPAddress, cmd.UserAgent, hash)
	if err != nil {
		s.logError("sign offer", err, serviceTenantIDField(existing.TenantID), serviceStringField("offer_letter_id", existing.ID.String()))
		return nil, err
	}
	_, _ = s.offerLetters.CreateOfferLetterEvent(ctx, &domain.OfferLetterEvent{TenantID: result.TenantID, OfferLetterID: result.ID, FromStatus: &existing.Status, ToStatus: result.Status, Action: "signed", ActorEmail: &cmd.SignerEmail, IPAddress: cmd.IPAddress, UserAgent: cmd.UserAgent, Metadata: map[string]any{"signature_hash": hash}}, nil)
	return result, nil
}

func (s *TenantService) DeleteOfferLetter(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if _, err := s.GetOfferLetter(ctx, tenantID, id); err != nil {
		return err
	}
	if err := s.offerLetters.DeleteOfferLetter(ctx, tenantID, id, actorID); err != nil {
		s.logError("delete offer", err, serviceTenantIDField(tenantID), serviceStringField("offer_letter_id", id.String()))
		return err
	}
	return nil
}

func (s *TenantService) ListOfferLetterEvents(ctx context.Context, tenantID uuid.UUID, offerLetterID uuid.UUID) ([]*domain.OfferLetterEvent, error) {
	if _, err := s.GetOfferLetter(ctx, tenantID, offerLetterID); err != nil {
		return nil, err
	}
	return s.offerLetters.ListOfferLetterEvents(ctx, tenantID, offerLetterID)
}

func offerTemplateInput(cmd ports.OfferLetterTemplateCommand) domain.OfferLetterTemplateInput {
	return domain.OfferLetterTemplateInput{TenantID: cmd.TenantID, Name: cmd.Name, Description: cmd.Description, Subject: cmd.Subject, BodyHTML: cmd.BodyHTML, FooterHTML: cmd.FooterHTML, Locale: cmd.Locale, IsDefault: cmd.IsDefault, IsActive: cmd.IsActive}
}

func offerLetterInput(cmd ports.OfferLetterCommand, candidateID *uuid.UUID) domain.OfferLetterInput {
	return domain.OfferLetterInput{TenantID: cmd.TenantID, ApplicationID: cmd.ApplicationID, CandidateID: candidateID, TemplateID: cmd.TemplateID, OfferedCTC: cmd.OfferedCTC, Currency: cmd.Currency, SalaryBreakdown: cmd.SalaryBreakdown, JoiningDate: cmd.JoiningDate, ValidUntilDate: cmd.ValidUntilDate, Status: cmd.Status, OfferLetterURL: cmd.OfferLetterURL, Subject: cmd.Subject, RenderedHTML: cmd.RenderedHTML, SignerEmail: cmd.SignerEmail}
}

func renderOfferLetter(item *domain.OfferLetter, tmpl *domain.OfferLetterTemplate, app *domain.CandidateApplication) {
	subject := valueOrDefault(tmpl.Subject, "Offer Letter")
	body := tmpl.BodyHTML
	replacements := map[string]string{
		"candidate_name":  strings.TrimSpace(strings.Join([]string{ptrString(app.CandidateFirstname), ptrString(app.CandidateLastname)}, " ")),
		"candidate_email": ptrString(app.CandidateEmail),
		"job_title":       ptrString(app.JobPostingTitle),
		"job_code":        ptrString(app.JobPostingCode),
		"offered_ctc":     fmt.Sprintf("%.2f", ptrFloat(item.OfferedCTC)),
		"currency":        item.Currency,
		"joining_date":    formatDate(item.JoiningDate),
		"valid_until":     formatDate(item.ValidUntilDate),
	}
	for key, value := range replacements {
		body = strings.ReplaceAll(body, "{{"+key+"}}", value)
		subject = strings.ReplaceAll(subject, "{{"+key+"}}", value)
	}
	if tmpl.FooterHTML != nil {
		body += "\n" + *tmpl.FooterHTML
	}
	item.Subject = &subject
	item.RenderedHTML = &body
}

func signatureHash(offer *domain.OfferLetter, cmd ports.OfferLetterSignatureCommand) string {
	sum := sha256.Sum256([]byte(strings.Join([]string{offer.ID.String(), ptrString(offer.SignatureToken), cmd.SignerName, cmd.SignerEmail, ptrString(cmd.IPAddress), ptrString(cmd.UserAgent)}, "|")))
	return hex.EncodeToString(sum[:])
}

func valueOrDefault(value *string, fallback string) string {
	if value == nil || strings.TrimSpace(*value) == "" {
		return fallback
	}
	return *value
}

func ptrString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func ptrFloat(value *float64) float64 {
	if value == nil {
		return 0
	}
	return *value
}

func formatDate(value *time.Time) string {
	if value == nil {
		return ""
	}
	return value.Format("02 Jan 2006")
}
