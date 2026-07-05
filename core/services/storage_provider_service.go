package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) GetStorageProviderSettings(ctx context.Context, tenantID uuid.UUID) (*domain.StorageProviderSettings, error) {
	if tenantID == uuid.Nil {
		s.logError("get storage provider settings", domain.ErrInvalidTenantID)
		return nil, domain.ErrInvalidTenantID
	}
	item, err := s.storageProviders.GetStorageProviderSettings(ctx, tenantID)
	if err == nil {
		return item.RedactSecrets(), nil
	}
	if errors.Is(err, domain.ErrStorageProviderSettingsNotFound) {
		if s.defaultStorageProvider != nil {
			copy := *s.defaultStorageProvider
			copy.TenantID = tenantID
			return copy.RedactSecrets(), nil
		}
		return nil, err
	}
	s.logError("get storage provider settings", err, serviceTenantIDField(tenantID))
	return nil, err
}

func (s *TenantService) UpsertStorageProviderSettings(ctx context.Context, cmd ports.StorageProviderSettingsCommand) (*domain.StorageProviderSettings, error) {
	if existing, err := s.storageProviders.GetStorageProviderSettings(ctx, cmd.TenantID); err == nil {
		if cmd.AccessKeyID == nil {
			cmd.AccessKeyID = existing.AccessKeyID
		}
		if cmd.SecretAccessKey == nil {
			cmd.SecretAccessKey = existing.SecretAccessKey
		}
	} else if !errors.Is(err, domain.ErrStorageProviderSettingsNotFound) {
		s.logError("get existing storage provider settings", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	item, err := domain.NewStorageProviderSettings(domain.StorageProviderSettingsInput{TenantID: cmd.TenantID, Provider: cmd.Provider, IsEnabled: cmd.IsEnabled, Bucket: cmd.Bucket, Region: cmd.Region, Endpoint: cmd.Endpoint, AccessKeyID: cmd.AccessKeyID, SecretAccessKey: cmd.SecretAccessKey, UseSSL: cmd.UseSSL, ForcePathStyle: cmd.ForcePathStyle, ObjectPrefix: cmd.ObjectPrefix, PublicBaseURL: cmd.PublicBaseURL, MaxFileSizeBytes: cmd.MaxFileSizeBytes, AllowedContentTypes: cmd.AllowedContentTypes})
	if err != nil {
		s.logError("validate storage provider settings", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	result, err := s.storageProviders.UpsertStorageProviderSettings(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("upsert storage provider settings", err, serviceTenantIDField(cmd.TenantID), serviceStringField("provider", item.Provider))
		return nil, err
	}
	return result.RedactSecrets(), nil
}

func (s *TenantService) TestStorageProviderSettings(ctx context.Context, cmd ports.StorageProviderTestCommand) (*domain.StorageProviderSettings, error) {
	settings, err := s.storageProviders.GetStorageProviderSettings(ctx, cmd.TenantID)
	if err != nil {
		s.logError("get storage provider settings for test", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	status := domain.NotifStatusSent
	message := "Storage provider is reachable."
	if s.objectStorage == nil {
		status = domain.NotifStatusFailed
		message = "Object storage adapter is not configured."
	} else if err := s.objectStorage.TestStorage(ctx, settings); err != nil {
		status = domain.NotifStatusFailed
		message = err.Error()
		s.logError("test storage provider", err, serviceTenantIDField(cmd.TenantID), serviceStringField("provider", settings.Provider))
	}
	updated, updateErr := s.storageProviders.UpdateStorageProviderTestResult(ctx, cmd.TenantID, settings.ID, status, &message, cmd.ActorID)
	if updateErr != nil {
		s.logError("update storage provider test result", updateErr, serviceTenantIDField(cmd.TenantID))
		return nil, updateErr
	}
	return updated.RedactSecrets(), nil
}

func (s *TenantService) DeleteStorageProviderSettings(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if tenantID == uuid.Nil || id == uuid.Nil {
		s.logError("delete storage provider settings", domain.ErrInvalidTenantID)
		return domain.ErrInvalidTenantID
	}
	if err := s.storageProviders.DeleteStorageProviderSettings(ctx, tenantID, id, actorID); err != nil {
		s.logError("delete storage provider settings", err, serviceTenantIDField(tenantID), serviceStringField("storage_provider_settings_id", id.String()))
		return err
	}
	return nil
}
