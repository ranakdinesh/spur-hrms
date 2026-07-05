package services

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) UpsertTenantProfile(ctx context.Context, cmd ports.UpsertTenantProfileCmd) (*domain.TenantProfile, error) {
	profile, err := domain.NewTenantProfile(cmd.TenantID, cmd.Subdomain, cmd.MobileActivationCode, cmd.DisplayName, cmd.LogoObjectKey)
	if err != nil {
		s.logError("validate tenant profile", err, serviceTenantIDField(cmd.TenantID), serviceStringField("subdomain", cmd.Subdomain))
		return nil, err
	}
	if err := s.validateTenantSubdomainAvailability(ctx, profile.TenantID, profile.Subdomain, profile.DisplayName); err != nil {
		s.logError("validate tenant subdomain availability", err, serviceTenantIDField(cmd.TenantID), serviceStringField("subdomain", profile.Subdomain))
		return nil, err
	}
	result, err := s.profiles.UpsertTenantProfile(ctx, profile)
	if err != nil {
		s.logError("upsert tenant profile", err, serviceTenantIDField(cmd.TenantID), serviceStringField("subdomain", cmd.Subdomain))
		return nil, err
	}
	s.log.Info().Str("tenant_id", result.TenantID.String()).Str("subdomain", result.Subdomain).Msg("hrms: tenant profile upserted")
	return result, nil
}

func (s *TenantService) GetTenantProfile(ctx context.Context, tenantID uuid.UUID) (*domain.TenantProfile, error) {
	result, err := s.profiles.GetTenantProfile(ctx, tenantID)
	if err != nil {
		s.logError("get tenant profile", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) ResolveTenantBySubdomain(ctx context.Context, subdomain string) (*domain.TenantProfile, error) {
	normalized := domain.NormalizeSubdomain(subdomain)
	result, err := s.profiles.GetTenantProfileBySubdomain(ctx, normalized)
	if err != nil {
		s.logError("resolve tenant by subdomain", err, serviceStringField("subdomain", normalized))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) ResolveTenantByActivationCode(ctx context.Context, code string) (*domain.TenantProfile, error) {
	normalized := domain.NormalizeMobileActivationCode(code)
	result, err := s.profiles.GetTenantProfileByActivationCode(ctx, normalized)
	if err != nil {
		s.logError("resolve tenant by activation code", err, serviceStringField("activation_code", normalized))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) GetTenantBranding(ctx context.Context, tenantID uuid.UUID) (*domain.TenantBranding, error) {
	var result *domain.TenantBranding
	if err := s.system.RunAsSystem(ctx, func(txCtx context.Context) error {
		var err error
		result, err = s.branding.GetTenantBranding(txCtx, tenantID)
		return err
	}); err != nil {
		s.logError("get tenant branding", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) UpsertTenantBranding(ctx context.Context, cmd ports.UpsertTenantBrandingCmd) (*domain.TenantBranding, error) {
	branding, err := domain.NewTenantBranding(domain.TenantBrandingInput{
		TenantID:          cmd.TenantID,
		Subdomain:         cmd.Subdomain,
		DisplayName:       cmd.DisplayName,
		LogoPath:          cmd.LogoPath,
		FaviconPath:       cmd.FaviconPath,
		Layout:            cmd.Layout,
		ColorMode:         cmd.ColorMode,
		SidebarSize:       cmd.SidebarSize,
		LayoutWidth:       cmd.LayoutWidth,
		CardLayout:        cmd.CardLayout,
		ThemeColor:        cmd.ThemeColor,
		PrimaryColor:      cmd.PrimaryColor,
		SecondaryColor:    cmd.SecondaryColor,
		TertiaryColor:     cmd.TertiaryColor,
		TopbarColor:       cmd.TopbarColor,
		SidebarColor:      cmd.SidebarColor,
		TopbarBackground:  cmd.TopbarBackground,
		SidebarBackground: cmd.SidebarBackground,
		FontFamily:        cmd.FontFamily,
		Preloader:         cmd.Preloader,
	})
	if err != nil {
		s.logError("validate tenant branding", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	if err := s.validateTenantSubdomainAvailability(ctx, branding.TenantID, branding.Subdomain, branding.DisplayName); err != nil {
		s.logError("validate tenant branding subdomain availability", err, serviceTenantIDField(cmd.TenantID), serviceStringField("subdomain", branding.Subdomain))
		return nil, err
	}
	var result *domain.TenantBranding
	if err := s.system.RunAsSystem(ctx, func(txCtx context.Context) error {
		var err error
		result, err = s.branding.UpsertTenantBranding(txCtx, branding, cmd.ActorID, tenantBrandingActivationCode(cmd.TenantID))
		return err
	}); err != nil {
		s.logError("upsert tenant branding", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	s.log.Info().Str("tenant_id", result.TenantID.String()).Str("subdomain", result.Subdomain).Msg("hrms: tenant branding upserted")
	return result, nil
}

func (s *TenantService) validateTenantSubdomainAvailability(ctx context.Context, tenantID uuid.UUID, subdomain string, displayName *string) error {
	if domain.IsReservedTenantSubdomain(subdomain) {
		return domain.ErrReservedSubdomain
	}
	candidateKey := domain.TenantSubdomainCollisionKey(subdomain)
	if candidateKey == "" {
		return domain.ErrReservedSubdomain
	}
	candidateKeys := map[string]struct{}{candidateKey: {}}
	if displayName != nil {
		if displayNameKey := domain.TenantSubdomainCollisionKey(*displayName); displayNameKey != "" {
			candidateKeys[displayNameKey] = struct{}{}
		}
	}
	existing, err := s.profiles.ListTenantProfiles(ctx)
	if err != nil {
		return err
	}
	for _, profile := range existing {
		if profile == nil || profile.TenantID == tenantID {
			continue
		}
		for _, existingKey := range tenantProfileCollisionKeys(profile) {
			if _, exists := candidateKeys[existingKey]; exists {
				return domain.ErrConflictingSubdomain
			}
		}
	}
	return nil
}

func tenantProfileCollisionKeys(profile *domain.TenantProfile) []string {
	if profile == nil {
		return nil
	}
	keys := make([]string, 0, 2)
	if key := domain.TenantSubdomainCollisionKey(profile.Subdomain); key != "" {
		keys = append(keys, key)
	}
	if profile.DisplayName != nil {
		if key := domain.TenantSubdomainCollisionKey(*profile.DisplayName); key != "" {
			keys = append(keys, key)
		}
	}
	return keys
}

func tenantBrandingActivationCode(tenantID uuid.UUID) string {
	clean := strings.ToUpper(strings.ReplaceAll(tenantID.String(), "-", ""))
	if len(clean) > 8 {
		clean = clean[:8]
	}
	return "BRAND-" + clean
}

func (s *TenantService) ResolveTenantBrandingByHost(ctx context.Context, host string) (*domain.TenantBranding, error) {
	subdomain, err := domain.BrandingSubdomainFromHost(host)
	if err != nil {
		s.logError("resolve tenant branding host", err, serviceStringField("host", host))
		return nil, err
	}
	var result *domain.TenantBranding
	if err := s.system.RunAsSystem(ctx, func(txCtx context.Context) error {
		var err error
		result, err = s.branding.ResolveTenantBrandingBySubdomain(txCtx, subdomain)
		return err
	}); err != nil {
		s.logError("resolve tenant branding", err, serviceStringField("subdomain", subdomain))
		return nil, err
	}
	return result, nil
}
