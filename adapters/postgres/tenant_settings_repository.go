package postgres

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func (s *Store) UpsertTenantSetting(ctx context.Context, tenantID uuid.UUID, key string, value map[string]any) (*domain.TenantSetting, error) {
	valueBytes, err := json.Marshal(value)
	if err != nil {
		return nil, s.logDBError(ctx, "marshal tenant setting", fmt.Errorf("hrms: marshal tenant setting: %w", err), tenantIDField(tenantID), stringField("setting_key", key))
	}
	row, err := s.getQueries(ctx).UpsertTenantSetting(ctx, sqlc.UpsertTenantSettingParams{
		TenantID: tenantID,
		Key:      key,
		Value:    valueBytes,
	})
	if err != nil {
		return nil, s.logDBError(ctx, "upsert tenant setting", err, tenantIDField(tenantID), stringField("setting_key", key))
	}
	item, err := mapTenantSetting(row)
	if err != nil {
		return nil, s.logDBError(ctx, "map tenant setting", err, tenantIDField(tenantID), stringField("setting_key", key))
	}
	return item, nil
}

func (s *Store) GetTenantSetting(ctx context.Context, tenantID uuid.UUID, key string) (*domain.TenantSetting, error) {
	row, err := s.getQueries(ctx).GetTenantSetting(ctx, sqlc.GetTenantSettingParams{
		TenantID: tenantID,
		Key:      key,
	})
	if err != nil {
		return nil, s.logDBError(ctx, "get tenant setting", err, tenantIDField(tenantID), stringField("setting_key", key))
	}
	item, err := mapTenantSetting(row)
	if err != nil {
		return nil, s.logDBError(ctx, "map tenant setting", err, tenantIDField(tenantID), stringField("setting_key", key))
	}
	return item, nil
}

func (s *Store) ListTenantSettings(ctx context.Context, tenantID uuid.UUID) ([]*domain.TenantSetting, error) {
	rows, err := s.getQueries(ctx).ListTenantSettings(ctx, tenantID)
	if err != nil {
		return nil, s.logDBError(ctx, "list tenant settings", err, tenantIDField(tenantID))
	}
	items := make([]*domain.TenantSetting, 0, len(rows))
	for _, row := range rows {
		item, err := mapTenantSetting(row)
		if err != nil {
			return nil, s.logDBError(ctx, "map tenant setting", err, tenantIDField(tenantID), stringField("setting_key", row.Key))
		}
		items = append(items, item)
	}
	return items, nil
}

func (s *Store) DeleteTenantSetting(ctx context.Context, tenantID uuid.UUID, key string) error {
	if err := s.getQueries(ctx).DeleteTenantSetting(ctx, sqlc.DeleteTenantSettingParams{
		TenantID: tenantID,
		Key:      key,
	}); err != nil {
		return s.logDBError(ctx, "delete tenant setting", err, tenantIDField(tenantID), stringField("setting_key", key))
	}
	return nil
}
