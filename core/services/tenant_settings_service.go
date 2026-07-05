package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) UpsertTenantSetting(ctx context.Context, cmd ports.UpsertTenantSettingCmd) (*domain.TenantSetting, error) {
	key, err := domain.ValidateSettingKey(cmd.Key)
	if err != nil {
		s.logError("validate tenant setting key", err, serviceTenantIDField(cmd.TenantID), serviceStringField("setting_key", cmd.Key))
		return nil, err
	}
	result, err := s.settings.UpsertTenantSetting(ctx, cmd.TenantID, key, cmd.Value)
	if err != nil {
		s.logError("upsert tenant setting", err, serviceTenantIDField(cmd.TenantID), serviceStringField("setting_key", key))
		return nil, err
	}
	s.log.Info().Str("tenant_id", result.TenantID.String()).Str("key", result.Key).Msg("hrms: tenant setting upserted")
	return result, nil
}

func (s *TenantService) GetTenantSetting(ctx context.Context, tenantID uuid.UUID, key string) (*domain.TenantSetting, error) {
	rawKey := key
	key, err := domain.ValidateSettingKey(key)
	if err != nil {
		s.logError("validate tenant setting key", err, serviceTenantIDField(tenantID), serviceStringField("setting_key", rawKey))
		return nil, err
	}
	result, err := s.settings.GetTenantSetting(ctx, tenantID, key)
	if err != nil {
		s.logError("get tenant setting", err, serviceTenantIDField(tenantID), serviceStringField("setting_key", key))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) ListTenantSettings(ctx context.Context, tenantID uuid.UUID) ([]*domain.TenantSetting, error) {
	result, err := s.settings.ListTenantSettings(ctx, tenantID)
	if err != nil {
		s.logError("list tenant settings", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) DeleteTenantSetting(ctx context.Context, tenantID uuid.UUID, key string) error {
	rawKey := key
	key, err := domain.ValidateSettingKey(key)
	if err != nil {
		s.logError("validate tenant setting key", err, serviceTenantIDField(tenantID), serviceStringField("setting_key", rawKey))
		return err
	}
	if err := s.settings.DeleteTenantSetting(ctx, tenantID, key); err != nil {
		s.logError("delete tenant setting", err, serviceTenantIDField(tenantID), serviceStringField("setting_key", key))
		return err
	}
	return nil
}
