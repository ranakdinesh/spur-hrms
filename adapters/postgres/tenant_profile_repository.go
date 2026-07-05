package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func (s *Store) UpsertTenantProfile(ctx context.Context, profile *domain.TenantProfile) (*domain.TenantProfile, error) {
	row, err := s.getQueries(ctx).UpsertTenantProfile(ctx, sqlc.UpsertTenantProfileParams{
		TenantID:             profile.TenantID,
		Subdomain:            profile.Subdomain,
		MobileActivationCode: profile.MobileActivationCode,
		DisplayName:          textFromPtr(profile.DisplayName),
		LogoObjectKey:        textFromPtr(profile.LogoObjectKey),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "upsert tenant profile", err, tenantIDField(profile.TenantID), stringField("subdomain", profile.Subdomain))
	}
	return mapTenantProfile(row), nil
}

func (s *Store) GetTenantProfile(ctx context.Context, tenantID uuid.UUID) (*domain.TenantProfile, error) {
	row, err := s.getQueries(ctx).GetTenantProfile(ctx, tenantID)
	if err != nil {
		return nil, s.logDBError(ctx, "get tenant profile", err, tenantIDField(tenantID))
	}
	return mapTenantProfile(row), nil
}

func (s *Store) GetTenantProfileBySubdomain(ctx context.Context, subdomain string) (*domain.TenantProfile, error) {
	row, err := s.getQueries(ctx).GetTenantProfileBySubdomain(ctx, subdomain)
	if err != nil {
		return nil, s.logDBError(ctx, "get tenant profile by subdomain", err, stringField("subdomain", subdomain))
	}
	return mapTenantProfile(row), nil
}

func (s *Store) GetTenantProfileByActivationCode(ctx context.Context, code string) (*domain.TenantProfile, error) {
	row, err := s.getQueries(ctx).GetTenantProfileByActivationCode(ctx, code)
	if err != nil {
		return nil, s.logDBError(ctx, "get tenant profile by activation code", err, stringField("activation_code", code))
	}
	return mapTenantProfile(row), nil
}
