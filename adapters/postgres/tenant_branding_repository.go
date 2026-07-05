package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func (s *Store) GetTenantBranding(ctx context.Context, tenantID uuid.UUID) (*domain.TenantBranding, error) {
	row, err := s.getQueries(ctx).GetTenantBrandingByTenantID(ctx, tenantID)
	if err != nil {
		return nil, s.logDBError(ctx, "get tenant branding", err, tenantIDField(tenantID))
	}
	return mapTenantBrandingByTenantID(row), nil
}

func (s *Store) UpsertTenantBranding(ctx context.Context, branding *domain.TenantBranding, actorID *uuid.UUID, activationCode string) (*domain.TenantBranding, error) {
	row, err := s.getQueries(ctx).UpsertTenantBranding(ctx, sqlc.UpsertTenantBrandingParams{
		TenantID:             branding.TenantID,
		LogoPath:             textFromPtr(branding.LogoPath),
		FaviconPath:          textFromPtr(branding.FaviconPath),
		Layout:               branding.Layout,
		ColorMode:            branding.ColorMode,
		SidebarSize:          branding.SidebarSize,
		LayoutWidth:          branding.LayoutWidth,
		CardLayout:           branding.CardLayout,
		ThemeColor:           branding.ThemeColor,
		PrimaryColor:         branding.PrimaryColor,
		SecondaryColor:       branding.SecondaryColor,
		TertiaryColor:        branding.TertiaryColor,
		TopbarColor:          branding.TopbarColor,
		SidebarColor:         branding.SidebarColor,
		TopbarBackground:     branding.TopbarBackground,
		SidebarBackground:    branding.SidebarBackground,
		FontFamily:           branding.FontFamily,
		Preloader:            branding.Preloader,
		CreatedBy:            uuidFromPtr(actorID),
		Subdomain:            branding.Subdomain,
		MobileActivationCode: activationCode,
		DisplayName:          textFromPtr(branding.DisplayName),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "upsert tenant branding", err, tenantIDField(branding.TenantID))
	}
	return mapTenantBrandingUpsert(row), nil
}

func (s *Store) ResolveTenantBrandingBySubdomain(ctx context.Context, subdomain string) (*domain.TenantBranding, error) {
	row, err := s.getQueries(ctx).ResolveTenantBrandingBySubdomain(ctx, subdomain)
	if err != nil {
		return nil, s.logDBError(ctx, "resolve tenant branding by subdomain", err, stringField("subdomain", subdomain))
	}
	return mapTenantBrandingBySubdomain(row), nil
}
