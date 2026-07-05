package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) GetCommunicationProviderSettings(ctx context.Context, tenantID uuid.UUID) (*domain.CommunicationProviderSettings, error) {
	if tenantID == uuid.Nil {
		s.logError("get communication provider settings", domain.ErrInvalidTenantID)
		return nil, domain.ErrInvalidTenantID
	}
	item, err := s.communicationProviders.GetCommunicationProviderSettings(ctx, tenantID)
	if err == nil {
		return item.RedactSecrets(), nil
	}
	if errors.Is(err, domain.ErrCommunicationProviderSettingsNotFound) {
		if s.defaultCommunicationProvider != nil {
			copy := *s.defaultCommunicationProvider
			copy.TenantID = tenantID
			return copy.RedactSecrets(), nil
		}
		item := defaultDisabledCommunicationProvider(tenantID)
		return item.RedactSecrets(), nil
	}
	s.logError("get communication provider settings", err)
	return nil, err
}

func (s *TenantService) UpsertCommunicationProviderSettings(ctx context.Context, cmd ports.CommunicationProviderSettingsCommand) (*domain.CommunicationProviderSettings, error) {
	if existing, err := s.communicationProviders.GetCommunicationProviderSettings(ctx, cmd.TenantID); err == nil {
		if cmd.SMSAuthKey == nil {
			cmd.SMSAuthKey = existing.SMSAuthKey
		}
		if cmd.WhatsAppAuthKey == nil {
			cmd.WhatsAppAuthKey = existing.WhatsAppAuthKey
		}
		if cmd.WebhookSigningSecret == nil {
			cmd.WebhookSigningSecret = existing.WebhookSigningSecret
		}
	} else if !errors.Is(err, domain.ErrCommunicationProviderSettingsNotFound) {
		s.logError("get existing communication provider settings", err)
		return nil, err
	}
	item, err := domain.NewCommunicationProviderSettings(domain.CommunicationProviderSettingsInput{
		TenantID: cmd.TenantID, SMSProvider: cmd.SMSProvider, SMSEnabled: cmd.SMSEnabled, SMSSenderID: cmd.SMSSenderID, SMSAuthKey: cmd.SMSAuthKey, SMSTemplateID: cmd.SMSTemplateID, SMSRoute: cmd.SMSRoute, SMSCountryCode: cmd.SMSCountryCode, SMSBaseURL: cmd.SMSBaseURL,
		WhatsAppProvider: cmd.WhatsAppProvider, WhatsAppEnabled: cmd.WhatsAppEnabled, WhatsAppAuthKey: cmd.WhatsAppAuthKey, WhatsAppAppName: cmd.WhatsAppAppName, WhatsAppSourceNumber: cmd.WhatsAppSourceNumber, WhatsAppTemplateID: cmd.WhatsAppTemplateID, WhatsAppTemplateName: cmd.WhatsAppTemplateName, WhatsAppNamespace: cmd.WhatsAppNamespace, WhatsAppBaseURL: cmd.WhatsAppBaseURL, WebhookSigningSecret: cmd.WebhookSigningSecret,
	})
	if err != nil {
		s.logError("validate communication provider settings", err)
		return nil, err
	}
	result, err := s.communicationProviders.UpsertCommunicationProviderSettings(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("upsert communication provider settings", err)
		return nil, err
	}
	return result.RedactSecrets(), nil
}

func (s *TenantService) TestCommunicationProviderSettings(ctx context.Context, cmd ports.CommunicationProviderTestCommand) (*domain.CommunicationProviderSettings, error) {
	channel, err := domain.ValidateCommunicationChannel(cmd.Channel)
	if err != nil {
		s.logError("validate communication provider test channel", err)
		return nil, err
	}
	to, err := domain.ValidateCommunicationRecipient(cmd.ToPhone)
	if err != nil {
		s.logError("validate communication provider test recipient", err)
		return nil, err
	}
	settings, err := s.communicationProviders.GetCommunicationProviderSettings(ctx, cmd.TenantID)
	if err != nil {
		s.logError("get communication provider settings for test", err)
		return nil, err
	}
	status := domain.NotifStatusSent
	message := "Test message accepted by provider."
	body := "HRMS provider test message."
	if cmd.Message != nil && *cmd.Message != "" {
		body = *cmd.Message
	}
	if s.communicationDelivery == nil {
		status = domain.NotifStatusFailed
		message = "Communication delivery adapter is not configured."
	} else {
		_, err = s.sendTestCommunication(ctx, settings, channel, ports.CommunicationMessage{TenantID: cmd.TenantID, To: to, Body: body})
		if err != nil {
			status = domain.NotifStatusFailed
			message = err.Error()
			s.logError("send communication provider test", err)
		}
	}
	updated, updateErr := s.communicationProviders.UpdateCommunicationProviderTestResult(ctx, cmd.TenantID, settings.ID, channel, status, &message, cmd.ActorID)
	if updateErr != nil {
		s.logError("update communication provider test result", updateErr)
		return nil, updateErr
	}
	return updated.RedactSecrets(), nil
}

func (s *TenantService) DeleteCommunicationProviderSettings(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if tenantID == uuid.Nil || id == uuid.Nil {
		s.logError("delete communication provider settings", domain.ErrInvalidTenantID)
		return domain.ErrInvalidTenantID
	}
	if err := s.communicationProviders.DeleteCommunicationProviderSettings(ctx, tenantID, id, actorID); err != nil {
		s.logError("delete communication provider settings", err)
		return err
	}
	return nil
}

func (s *TenantService) sendTestCommunication(ctx context.Context, settings *domain.CommunicationProviderSettings, channel string, message ports.CommunicationMessage) (*ports.CommunicationDeliveryResult, error) {
	switch channel {
	case domain.CommunicationChannelSMS:
		if !settings.SMSEnabled {
			if s.log != nil {
				s.log.Warn().Str("tenant_id", settings.TenantID.String()).Msg("hrms sms provider disabled")
			}
			return nil, domain.ErrSMSProviderDisabled
		}
		return s.communicationDelivery.SendSMS(ctx, settings, message)
	case domain.CommunicationChannelWhatsApp:
		if !settings.WhatsAppEnabled {
			if s.log != nil {
				s.log.Warn().Str("tenant_id", settings.TenantID.String()).Msg("hrms whatsapp provider disabled")
			}
			return nil, domain.ErrWhatsAppProviderDisabled
		}
		return s.communicationDelivery.SendWhatsApp(ctx, settings, message)
	default:
		return nil, fmt.Errorf("%w: %s", domain.ErrInvalidCommunicationChannel, channel)
	}
}

func defaultDisabledCommunicationProvider(tenantID uuid.UUID) *domain.CommunicationProviderSettings {
	item, _ := domain.NewCommunicationProviderSettings(domain.CommunicationProviderSettingsInput{TenantID: tenantID, SMSProvider: domain.CommunicationProviderLocal, WhatsAppProvider: domain.CommunicationProviderLocal})
	return item
}
