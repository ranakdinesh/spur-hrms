package services

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) GetEmailProviderSettings(ctx context.Context, tenantID uuid.UUID) (*domain.EmailProviderSettings, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate email provider get tenant", err)
		return nil, err
	}
	if s.defaultEmailProvider != nil {
		copy := *s.defaultEmailProvider
		copy.TenantID = tenantID
		return copy.RedactSecrets(), nil
	}
	if s.globalEmailProviderOnly {
		err := domain.ErrGlobalEmailProviderNotConfigured
		s.logError("get global email provider settings", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	item, err := s.emailProviders.GetEmailProviderSettings(ctx, tenantID)
	if err != nil {
		s.logError("get email provider settings", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	return item.RedactSecrets(), nil
}

func (s *TenantService) UpsertEmailProviderSettings(ctx context.Context, cmd ports.EmailProviderSettingsCommand) (*domain.EmailProviderSettings, error) {
	if s.globalEmailProviderOnly || s.defaultEmailProvider != nil {
		s.logError("reject tenant email provider settings update", domain.ErrGlobalEmailProviderReadOnly, serviceTenantIDField(cmd.TenantID))
		return nil, domain.ErrGlobalEmailProviderReadOnly
	}
	item, err := domain.NewEmailProviderSettings(domain.EmailProviderSettingsInput{
		TenantID: cmd.TenantID, Provider: cmd.Provider, IsEnabled: cmd.IsEnabled, FromName: cmd.FromName, FromEmail: cmd.FromEmail, ReplyToEmail: cmd.ReplyToEmail,
		SMTPHost: cmd.SMTPHost, SMTPPort: cmd.SMTPPort, SMTPUsername: cmd.SMTPUsername, SMTPPassword: cmd.SMTPPassword, SMTPEncryption: cmd.SMTPEncryption,
		SendGridAPIKey: cmd.SendGridAPIKey, SendGridSandboxMode: cmd.SendGridSandboxMode, WebhookSigningSecret: cmd.WebhookSigningSecret,
		SPFStatus: cmd.SPFStatus, DKIMStatus: cmd.DKIMStatus, DMARCStatus: cmd.DMARCStatus,
	})
	if err != nil {
		s.logError("validate email provider settings", err, serviceTenantIDField(cmd.TenantID), serviceStringField("provider", cmd.Provider))
		return nil, err
	}
	result, err := s.emailProviders.UpsertEmailProviderSettings(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("upsert email provider settings", err, serviceTenantIDField(cmd.TenantID), serviceStringField("provider", item.Provider))
		return nil, err
	}
	return result.RedactSecrets(), nil
}

func (s *TenantService) TestEmailProviderSettings(ctx context.Context, cmd ports.EmailProviderTestCommand) (*domain.EmailProviderSettings, error) {
	if cmd.TenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate email provider test tenant", err)
		return nil, err
	}
	toEmail := strings.ToLower(strings.TrimSpace(cmd.ToEmail))
	if toEmail == "" {
		err := domain.ErrInvalidEmailRecipient
		s.logError("validate email provider test recipient", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	settings, err := s.activeEmailProviderSettings(ctx, cmd.TenantID)
	if err != nil {
		s.logError("get email provider settings for test", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	subject := "HRMS email provider test"
	if cmd.Subject != nil && strings.TrimSpace(*cmd.Subject) != "" {
		subject = strings.TrimSpace(*cmd.Subject)
	}
	message := "This is a test email from Spur HRMS."
	if cmd.Message != nil && strings.TrimSpace(*cmd.Message) != "" {
		message = strings.TrimSpace(*cmd.Message)
	}
	status := domain.NotifStatusSent
	statusMessage := "Test email accepted by provider."
	if s.emailDelivery == nil {
		status = domain.NotifStatusFailed
		statusMessage = "Email delivery sender is not configured."
	} else if _, err := s.emailDelivery.SendEmail(ctx, settings, ports.EmailMessage{TenantID: cmd.TenantID, To: toEmail, Subject: subject, TextBody: message}); err != nil {
		status = domain.NotifStatusFailed
		statusMessage = err.Error()
		s.logError("send email provider test", err, serviceTenantIDField(cmd.TenantID), serviceStringField("provider", settings.Provider))
	}
	if s.globalEmailProviderOnly || s.defaultEmailProvider != nil {
		copy := *settings
		now := time.Now().UTC()
		copy.LastTestAt = &now
		copy.LastTestStatus = &status
		copy.LastTestMessage = &statusMessage
		return copy.RedactSecrets(), nil
	}
	updated, err := s.emailProviders.UpdateEmailProviderTestResult(ctx, cmd.TenantID, settings.ID, status, &statusMessage, cmd.ActorID)
	if err != nil {
		s.logError("update email provider test result", err, serviceTenantIDField(cmd.TenantID), serviceStringField("email_provider_settings_id", settings.ID.String()))
		return nil, err
	}
	return updated.RedactSecrets(), nil
}

func (s *TenantService) DeleteEmailProviderSettings(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if s.globalEmailProviderOnly || s.defaultEmailProvider != nil {
		s.logError("reject tenant email provider settings delete", domain.ErrGlobalEmailProviderReadOnly, serviceTenantIDField(tenantID))
		return domain.ErrGlobalEmailProviderReadOnly
	}
	if tenantID == uuid.Nil || id == uuid.Nil {
		err := domain.ErrEmailProviderSettingsNotFound
		s.logError("validate delete email provider settings", err)
		return err
	}
	if err := s.emailProviders.DeleteEmailProviderSettings(ctx, tenantID, id, actorID); err != nil {
		s.logError("delete email provider settings", err, serviceTenantIDField(tenantID), serviceStringField("email_provider_settings_id", id.String()))
		return err
	}
	return nil
}

func (s *TenantService) activeEmailProviderSettings(ctx context.Context, tenantID uuid.UUID) (*domain.EmailProviderSettings, error) {
	if s.defaultEmailProvider != nil {
		copy := *s.defaultEmailProvider
		copy.TenantID = tenantID
		return &copy, nil
	}
	if s.globalEmailProviderOnly {
		return nil, domain.ErrGlobalEmailProviderNotConfigured
	}
	item, err := s.emailProviders.GetEmailProviderSettings(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	return item, nil
}
