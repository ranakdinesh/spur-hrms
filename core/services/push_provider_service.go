package services

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) GetPushProviderSettings(ctx context.Context, tenantID uuid.UUID) (*domain.PushProviderSettings, error) {
	if tenantID == uuid.Nil {
		s.logError("get push provider settings", domain.ErrInvalidTenantID)
		return nil, domain.ErrInvalidTenantID
	}
	item, err := s.pushProviders.GetPushProviderSettings(ctx, tenantID)
	if err == nil {
		return item.RedactSecrets(), nil
	}
	if errors.Is(err, domain.ErrPushProviderSettingsNotFound) && s.defaultPushProvider != nil {
		copy := *s.defaultPushProvider
		copy.TenantID = tenantID
		return copy.RedactSecrets(), nil
	}
	s.logError("get push provider settings", err, serviceTenantIDField(tenantID))
	return nil, err
}

func (s *TenantService) UpsertPushProviderSettings(ctx context.Context, cmd ports.PushProviderSettingsCommand) (*domain.PushProviderSettings, error) {
	if existing, err := s.pushProviders.GetPushProviderSettings(ctx, cmd.TenantID); err == nil && cmd.PrivateKey == nil {
		cmd.PrivateKey = existing.PrivateKey
	} else if err != nil && !errors.Is(err, domain.ErrPushProviderSettingsNotFound) {
		s.logError("get existing push provider settings", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	item, err := domain.NewPushProviderSettings(domain.PushProviderSettingsInput{TenantID: cmd.TenantID, Provider: cmd.Provider, IsEnabled: cmd.IsEnabled, ProjectID: cmd.ProjectID, ClientEmail: cmd.ClientEmail, PrivateKey: cmd.PrivateKey, PrivateKeyID: cmd.PrivateKeyID, AuthURI: cmd.AuthURI, TokenURI: cmd.TokenURI, AndroidEnabled: cmd.AndroidEnabled, IOSEnabled: cmd.IOSEnabled, WebEnabled: cmd.WebEnabled, DefaultClickAction: cmd.DefaultClickAction, DefaultImageURL: cmd.DefaultImageURL, TTLSeconds: cmd.TTLSeconds, CollapseKey: cmd.CollapseKey})
	if err != nil {
		s.logError("validate push provider settings", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	result, err := s.pushProviders.UpsertPushProviderSettings(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("upsert push provider settings", err, serviceTenantIDField(cmd.TenantID), serviceStringField("provider", item.Provider))
		return nil, err
	}
	return result.RedactSecrets(), nil
}

func (s *TenantService) TestPushProviderSettings(ctx context.Context, cmd ports.PushProviderTestCommand) (*domain.PushProviderSettings, error) {
	settings, err := s.activePushProviderSettings(ctx, cmd.TenantID)
	if err != nil {
		s.logError("get push provider settings for test", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	status := domain.NotifStatusSent
	message := "Push provider test accepted."
	title := "HRMS push test"
	if cmd.Title != nil && strings.TrimSpace(*cmd.Title) != "" {
		title = strings.TrimSpace(*cmd.Title)
	}
	body := "This confirms the tenant push provider can send to a Flutter device token."
	if cmd.Message != nil && strings.TrimSpace(*cmd.Message) != "" {
		body = strings.TrimSpace(*cmd.Message)
	}
	if s.pushDelivery == nil {
		status = domain.NotifStatusFailed
		message = "Push delivery sender is not configured."
	} else if _, err := s.pushDelivery.SendPush(ctx, settings, ports.PushMessage{TenantID: cmd.TenantID, Token: strings.TrimSpace(cmd.Token), Title: title, Body: body}); err != nil {
		status = domain.NotifStatusFailed
		message = err.Error()
		s.logError("send push provider test", err, serviceTenantIDField(cmd.TenantID), serviceStringField("provider", settings.Provider))
	}
	updated, updateErr := s.pushProviders.UpdatePushProviderTestResult(ctx, cmd.TenantID, settings.ID, status, &message, cmd.ActorID)
	if updateErr != nil {
		s.logError("update push provider test result", updateErr, serviceTenantIDField(cmd.TenantID))
		return nil, updateErr
	}
	return updated.RedactSecrets(), nil
}

func (s *TenantService) DeletePushProviderSettings(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if tenantID == uuid.Nil || id == uuid.Nil {
		s.logError("delete push provider settings", domain.ErrInvalidTenantID)
		return domain.ErrInvalidTenantID
	}
	if err := s.pushProviders.DeletePushProviderSettings(ctx, tenantID, id, actorID); err != nil {
		s.logError("delete push provider settings", err, serviceTenantIDField(tenantID), serviceStringField("push_provider_settings_id", id.String()))
		return err
	}
	return nil
}

func (s *TenantService) activePushProviderSettings(ctx context.Context, tenantID uuid.UUID) (*domain.PushProviderSettings, error) {
	settings, err := s.pushProviders.GetPushProviderSettings(ctx, tenantID)
	if err == nil {
		return settings, nil
	}
	if errors.Is(err, domain.ErrPushProviderSettingsNotFound) && s.defaultPushProvider != nil {
		copy := *s.defaultPushProvider
		copy.TenantID = tenantID
		return &copy, nil
	}
	return nil, err
}
