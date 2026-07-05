package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) CreateAgreementTemplate(ctx context.Context, cmd ports.AgreementTemplateCommand) (*domain.AgreementTemplate, error) {
	item, err := domain.NewAgreementTemplate(agreementTemplateInput(cmd))
	if err != nil {
		s.logError("validate agreement template create", err, serviceTenantIDField(cmd.TenantID), serviceStringField("agreement_type", cmd.AgreementType))
		return nil, err
	}
	result, err := s.agreements.CreateAgreementTemplate(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create agreement template", err, serviceTenantIDField(cmd.TenantID), serviceStringField("agreement_type", item.AgreementType))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) ListAgreementTemplates(ctx context.Context, tenantID uuid.UUID, agreementType *string) ([]*domain.AgreementTemplate, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate agreement template list tenant", err)
		return nil, err
	}
	cleanType := cleanStringPtr(agreementType)
	if cleanType != nil {
		normalized, err := domain.ValidateAgreementType(*cleanType)
		if err != nil {
			s.logError("validate agreement template list type", err, serviceTenantIDField(tenantID), serviceStringField("agreement_type", *cleanType))
			return nil, err
		}
		cleanType = &normalized
	}
	items, err := s.agreements.ListAgreementTemplates(ctx, tenantID, cleanType)
	if err != nil {
		s.logError("list agreement templates", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	return items, nil
}

func (s *TenantService) GetAgreementTemplate(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.AgreementTemplate, error) {
	if tenantID == uuid.Nil || id == uuid.Nil {
		err := domain.ErrInvalidAgreementTemplate
		s.logError("validate agreement template get", err)
		return nil, err
	}
	item, err := s.agreements.GetAgreementTemplate(ctx, tenantID, id)
	if err != nil {
		s.logError("get agreement template", err, serviceTenantIDField(tenantID), serviceStringField("agreement_template_id", id.String()))
		return nil, err
	}
	return item, nil
}

func (s *TenantService) UpdateAgreementTemplate(ctx context.Context, cmd ports.AgreementTemplateCommand) (*domain.AgreementTemplate, error) {
	if cmd.ID == uuid.Nil {
		err := domain.ErrInvalidAgreementTemplate
		s.logError("validate agreement template update id", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	if _, err := s.GetAgreementTemplate(ctx, cmd.TenantID, cmd.ID); err != nil {
		return nil, err
	}
	item, err := domain.NewAgreementTemplate(agreementTemplateInput(cmd))
	if err != nil {
		s.logError("validate agreement template update", err, serviceTenantIDField(cmd.TenantID), serviceStringField("agreement_template_id", cmd.ID.String()))
		return nil, err
	}
	item.ID = cmd.ID
	result, err := s.agreements.UpdateAgreementTemplate(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("update agreement template", err, serviceTenantIDField(cmd.TenantID), serviceStringField("agreement_template_id", cmd.ID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) DeleteAgreementTemplate(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if _, err := s.GetAgreementTemplate(ctx, tenantID, id); err != nil {
		return err
	}
	if err := s.agreements.DeleteAgreementTemplate(ctx, tenantID, id, actorID); err != nil {
		s.logError("delete agreement template", err, serviceTenantIDField(tenantID), serviceStringField("agreement_template_id", id.String()))
		return err
	}
	return nil
}

func (s *TenantService) GenerateAgreement(ctx context.Context, cmd ports.AgreementCommand) (*domain.Agreement, error) {
	agreementType, err := domain.ValidateAgreementType(cmd.AgreementType)
	if err != nil {
		s.logError("validate agreement type", err, serviceTenantIDField(cmd.TenantID), serviceStringField("agreement_type", cmd.AgreementType))
		return nil, err
	}
	cmd.AgreementType = agreementType
	if cmd.WorkerProfileID != nil {
		if _, err := s.GetWorkerProfile(ctx, cmd.TenantID, *cmd.WorkerProfileID); err != nil {
			return nil, err
		}
	}
	if cmd.EngagementID != nil {
		if _, err := s.GetEngagement(ctx, cmd.TenantID, *cmd.EngagementID); err != nil {
			return nil, err
		}
	}
	if cmd.ProjectID != nil {
		if _, err := s.GetProject(ctx, cmd.TenantID, *cmd.ProjectID); err != nil {
			return nil, err
		}
	}
	if cmd.TemplateID == nil {
		if tmpl, err := s.agreements.GetDefaultAgreementTemplate(ctx, cmd.TenantID, agreementType); err == nil && tmpl != nil {
			cmd.TemplateID = &tmpl.ID
		}
	}
	item, err := domain.NewAgreement(agreementInput(cmd))
	if err != nil {
		s.logError("validate agreement create", err, serviceTenantIDField(cmd.TenantID), serviceStringField("agreement_type", agreementType))
		return nil, err
	}
	if item.TemplateID != nil {
		tmpl, err := s.GetAgreementTemplate(ctx, cmd.TenantID, *item.TemplateID)
		if err != nil {
			return nil, err
		}
		s.renderAgreement(ctx, item, tmpl)
	}
	if cmd.SignatureRequired {
		token, err := domain.NewAgreementSignatureToken()
		if err != nil {
			s.logError("generate agreement signature token", err, serviceTenantIDField(cmd.TenantID))
			return nil, err
		}
		item.SignatureToken = &token
	}
	var result *domain.Agreement
	err = s.system.RunAsSystem(ctx, func(txCtx context.Context) error {
		saved, err := s.agreements.CreateAgreement(txCtx, item, cmd.ActorID)
		if err != nil {
			return err
		}
		result = saved
		_, _ = s.agreements.CreateAgreementEvent(txCtx, &domain.AgreementEvent{TenantID: saved.TenantID, AgreementID: saved.ID, ToStatus: saved.Status, Action: "generated", Remarks: saved.Subject, Metadata: map[string]any{"version": saved.Version, "agreement_type": saved.AgreementType}}, cmd.ActorID)
		return nil
	})
	if err != nil {
		s.logError("create agreement", err, serviceTenantIDField(cmd.TenantID), serviceStringField("agreement_type", agreementType))
		return nil, err
	}
	if _, err := s.ensureAgreementPDF(ctx, result, cmd.ActorID); err != nil {
		s.logError("store generated agreement pdf", err, serviceTenantIDField(cmd.TenantID), serviceStringField("agreement_id", result.ID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) ListAgreements(ctx context.Context, filter domain.AgreementFilter) ([]*domain.Agreement, error) {
	if filter.TenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate agreement list tenant", err)
		return nil, err
	}
	items, err := s.agreements.ListAgreements(ctx, filter)
	if err != nil {
		s.logError("list agreements", err, serviceTenantIDField(filter.TenantID))
		return nil, err
	}
	return items, nil
}

func (s *TenantService) GetAgreement(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.Agreement, error) {
	if tenantID == uuid.Nil || id == uuid.Nil {
		err := domain.ErrInvalidAgreementID
		s.logError("validate agreement get", err)
		return nil, err
	}
	item, err := s.agreements.GetAgreement(ctx, tenantID, id)
	if err != nil {
		s.logError("get agreement", err, serviceTenantIDField(tenantID), serviceStringField("agreement_id", id.String()))
		return nil, err
	}
	return item, nil
}

func (s *TenantService) GetAgreementBySignatureToken(ctx context.Context, token string) (*domain.Agreement, error) {
	if strings.TrimSpace(token) == "" {
		return nil, domain.ErrAgreementTokenMissing
	}
	return s.agreements.GetAgreementBySignatureToken(ctx, token)
}

func (s *TenantService) UpdateAgreementStatus(ctx context.Context, cmd ports.AgreementStatusCommand) (*domain.Agreement, error) {
	existing, err := s.GetAgreement(ctx, cmd.TenantID, cmd.ID)
	if err != nil {
		return nil, err
	}
	status, err := domain.ValidateAgreementStatus(cmd.Status)
	if err != nil {
		s.logError("validate agreement status", err, serviceTenantIDField(cmd.TenantID), serviceStringField("agreement_id", cmd.ID.String()), serviceStringField("status", cmd.Status))
		return nil, err
	}
	if existing.Status == domain.AgreementStatusRevoked && status != domain.AgreementStatusRevoked {
		err := domain.ErrAgreementCannotBeChanged
		s.logError("validate revoked agreement status", err, serviceTenantIDField(cmd.TenantID), serviceStringField("agreement_id", cmd.ID.String()))
		return nil, err
	}
	if existing.Status == domain.AgreementStatusSigned && status != domain.AgreementStatusSigned {
		err := domain.ErrAgreementCannotBeChanged
		s.logError("validate signed agreement status", err, serviceTenantIDField(cmd.TenantID), serviceStringField("agreement_id", cmd.ID.String()))
		return nil, err
	}
	result, err := s.agreements.UpdateAgreementStatus(ctx, cmd.TenantID, cmd.ID, status, cmd.ActorID)
	if err != nil {
		s.logError("update agreement status", err, serviceTenantIDField(cmd.TenantID), serviceStringField("agreement_id", cmd.ID.String()), serviceStringField("status", status))
		return nil, err
	}
	_, _ = s.agreements.CreateAgreementEvent(ctx, &domain.AgreementEvent{TenantID: result.TenantID, AgreementID: result.ID, FromStatus: &existing.Status, ToStatus: result.Status, Action: strings.ToLower(result.Status), Remarks: cmd.Remarks, Metadata: map[string]any{"version": result.Version, "agreement_type": result.AgreementType}}, cmd.ActorID)
	return result, nil
}

func (s *TenantService) SignAgreement(ctx context.Context, cmd ports.AgreementSignatureCommand) (*domain.Agreement, error) {
	if strings.TrimSpace(cmd.Token) == "" || strings.TrimSpace(cmd.SignerName) == "" || strings.TrimSpace(cmd.SignerEmail) == "" {
		err := domain.ErrAgreementSignature
		s.logError("validate agreement signature", err)
		return nil, err
	}
	existing, err := s.agreements.GetAgreementBySignatureToken(ctx, cmd.Token)
	if err != nil {
		s.logError("get agreement by signature token", err)
		return nil, err
	}
	hash := agreementSignatureHash(existing, cmd)
	result, err := s.agreements.SignAgreement(ctx, cmd.Token, cmd.SignerName, cmd.SignerEmail, cmd.IPAddress, cmd.UserAgent, hash)
	if err != nil {
		s.logError("sign agreement", err, serviceTenantIDField(existing.TenantID), serviceStringField("agreement_id", existing.ID.String()))
		return nil, err
	}
	_, _ = s.agreements.CreateAgreementEvent(ctx, &domain.AgreementEvent{TenantID: result.TenantID, AgreementID: result.ID, FromStatus: &existing.Status, ToStatus: result.Status, Action: "signed", ActorEmail: &cmd.SignerEmail, IPAddress: cmd.IPAddress, UserAgent: cmd.UserAgent, Metadata: map[string]any{"signature_hash": hash}}, nil)
	if _, err := s.ensureAgreementPDF(ctx, result, nil); err != nil {
		s.logError("refresh signed agreement pdf", err, serviceTenantIDField(result.TenantID), serviceStringField("agreement_id", result.ID.String()))
	}
	return result, nil
}

func (s *TenantService) DeleteAgreement(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if _, err := s.GetAgreement(ctx, tenantID, id); err != nil {
		return err
	}
	if err := s.agreements.DeleteAgreement(ctx, tenantID, id, actorID); err != nil {
		s.logError("delete agreement", err, serviceTenantIDField(tenantID), serviceStringField("agreement_id", id.String()))
		return err
	}
	return nil
}

func (s *TenantService) ListAgreementEvents(ctx context.Context, tenantID uuid.UUID, agreementID uuid.UUID) ([]*domain.AgreementEvent, error) {
	if _, err := s.GetAgreement(ctx, tenantID, agreementID); err != nil {
		return nil, err
	}
	return s.agreements.ListAgreementEvents(ctx, tenantID, agreementID)
}

func (s *TenantService) RenderAgreementPDF(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) ([]byte, string, error) {
	agreement, err := s.GetAgreement(ctx, tenantID, id)
	if err != nil {
		return nil, "", err
	}
	data, err := s.ensureAgreementPDF(ctx, agreement, actorID)
	if err != nil {
		return nil, "", err
	}
	_, _ = s.agreements.CreateAgreementEvent(ctx, &domain.AgreementEvent{TenantID: agreement.TenantID, AgreementID: agreement.ID, FromStatus: &agreement.Status, ToStatus: agreement.Status, Action: "downloaded", Metadata: map[string]any{"version": agreement.Version, "agreement_type": agreement.AgreementType}}, actorID)
	return data, agreementPDFName(agreement), nil
}

func (s *TenantService) ensureAgreementPDF(ctx context.Context, agreement *domain.Agreement, actorID *uuid.UUID) ([]byte, error) {
	if s.agreementPDF == nil || s.agreementStorage == nil {
		return nil, domain.ErrAgreementPDFMissing
	}
	data, err := s.agreementPDF.RenderAgreementPDF(ctx, ports.AgreementDocument{Agreement: agreement})
	if err != nil {
		return nil, err
	}
	path, err := s.agreementStorage.StoreAgreementPDF(ctx, ports.StoreAgreementPDFInput{TenantID: agreement.TenantID, AgreementID: agreement.ID, AgreementType: agreement.AgreementType, FileName: agreementPDFName(agreement), ContentType: "application/pdf", Content: data})
	if err != nil {
		return nil, err
	}
	updated, err := s.agreements.UpdateAgreementPDF(ctx, agreement.TenantID, agreement.ID, &path, actorID)
	if err != nil {
		return nil, err
	}
	*agreement = *updated
	return data, nil
}

func (s *TenantService) renderAgreement(ctx context.Context, item *domain.Agreement, tmpl *domain.AgreementTemplate) {
	subject := valueOrDefault(tmpl.Subject, strings.Title(strings.ReplaceAll(item.AgreementType, "_", " "))+" Agreement")
	body := tmpl.BodyHTML
	replacements := map[string]string{
		"agreement_title": item.Title,
		"agreement_type":  strings.Title(strings.ReplaceAll(item.AgreementType, "_", " ")),
		"issue_date":      formatDate(item.IssueDate),
		"effective_date":  formatDate(item.EffectiveDate),
		"end_date":        formatDate(item.EndDate),
	}
	if item.WorkerProfileID != nil {
		if worker, err := s.GetWorkerProfile(ctx, item.TenantID, *item.WorkerProfileID); err == nil && worker != nil {
			replacements["worker_name"] = worker.DisplayName
			replacements["worker_code"] = ptrString(worker.WorkerCode)
			replacements["worker_email"] = ptrString(worker.Email)
		}
	}
	if item.EngagementID != nil {
		if engagement, err := s.GetEngagement(ctx, item.TenantID, *item.EngagementID); err == nil && engagement != nil {
			replacements["engagement_title"] = engagement.Title
			replacements["engagement_code"] = ptrString(engagement.EngagementCode)
			replacements["cost_center"] = ptrString(engagement.CostCenter)
		}
	}
	if item.ProjectID != nil {
		if project, err := s.GetProject(ctx, item.TenantID, *item.ProjectID); err == nil && project != nil {
			replacements["project_name"] = project.Name
			replacements["project_code"] = ptrString(project.ProjectCode)
		}
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

func agreementTemplateInput(cmd ports.AgreementTemplateCommand) domain.AgreementTemplateInput {
	return domain.AgreementTemplateInput{TenantID: cmd.TenantID, AgreementType: cmd.AgreementType, Name: cmd.Name, Description: cmd.Description, Subject: cmd.Subject, BodyHTML: cmd.BodyHTML, FooterHTML: cmd.FooterHTML, Locale: cmd.Locale, IsDefault: cmd.IsDefault, IsActive: cmd.IsActive, Metadata: cmd.Metadata}
}

func agreementInput(cmd ports.AgreementCommand) domain.AgreementInput {
	return domain.AgreementInput{TenantID: cmd.TenantID, AgreementType: cmd.AgreementType, Title: cmd.Title, TemplateID: cmd.TemplateID, WorkerProfileID: cmd.WorkerProfileID, EngagementID: cmd.EngagementID, ProjectID: cmd.ProjectID, Subject: cmd.Subject, RenderedHTML: cmd.RenderedHTML, IssueDate: cmd.IssueDate, EffectiveDate: cmd.EffectiveDate, EndDate: cmd.EndDate, SignatureRequired: cmd.SignatureRequired, SignerEmail: cmd.SignerEmail, Metadata: cmd.Metadata}
}

func agreementSignatureHash(agreement *domain.Agreement, cmd ports.AgreementSignatureCommand) string {
	sum := sha256.Sum256([]byte(strings.Join([]string{agreement.ID.String(), ptrString(agreement.SignatureToken), cmd.SignerName, cmd.SignerEmail, ptrString(cmd.IPAddress), ptrString(cmd.UserAgent)}, "|")))
	return hex.EncodeToString(sum[:])
}

func agreementPDFName(agreement *domain.Agreement) string {
	return fmt.Sprintf("%s-agreement-%s-v%d.pdf", agreement.AgreementType, agreement.ID.String(), agreement.Version)
}
