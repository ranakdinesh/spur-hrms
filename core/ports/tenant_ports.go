package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

type TenantProfileRepo interface {
	UpsertTenantProfile(ctx context.Context, profile *domain.TenantProfile) (*domain.TenantProfile, error)
	GetTenantProfile(ctx context.Context, tenantID uuid.UUID) (*domain.TenantProfile, error)
	ListTenantProfiles(ctx context.Context) ([]*domain.TenantProfile, error)
	GetTenantProfileBySubdomain(ctx context.Context, subdomain string) (*domain.TenantProfile, error)
	GetTenantProfileByActivationCode(ctx context.Context, code string) (*domain.TenantProfile, error)
}

type TenantBrandingRepo interface {
	GetTenantBranding(ctx context.Context, tenantID uuid.UUID) (*domain.TenantBranding, error)
	UpsertTenantBranding(ctx context.Context, branding *domain.TenantBranding, actorID *uuid.UUID, activationCode string) (*domain.TenantBranding, error)
	ResolveTenantBrandingBySubdomain(ctx context.Context, subdomain string) (*domain.TenantBranding, error)
}

type TenantSettingsRepo interface {
	UpsertTenantSetting(ctx context.Context, tenantID uuid.UUID, key string, value map[string]any) (*domain.TenantSetting, error)
	GetTenantSetting(ctx context.Context, tenantID uuid.UUID, key string) (*domain.TenantSetting, error)
	ListTenantSettings(ctx context.Context, tenantID uuid.UUID) ([]*domain.TenantSetting, error)
	DeleteTenantSetting(ctx context.Context, tenantID uuid.UUID, key string) error
}

type UpsertTenantProfileCmd struct {
	TenantID             uuid.UUID `json:"tenant_id"`
	Subdomain            string    `json:"subdomain"`
	MobileActivationCode string    `json:"mobile_activation_code"`
	DisplayName          *string   `json:"display_name,omitempty"`
	LogoObjectKey        *string   `json:"logo_object_key,omitempty"`
}

type UpsertTenantBrandingCmd struct {
	TenantID          uuid.UUID  `json:"tenant_id"`
	Subdomain         string     `json:"subdomain"`
	DisplayName       *string    `json:"display_name,omitempty"`
	LogoPath          *string    `json:"logo_path,omitempty"`
	FaviconPath       *string    `json:"favicon_path,omitempty"`
	Layout            string     `json:"layout"`
	ColorMode         string     `json:"color_mode"`
	SidebarSize       string     `json:"sidebar_size"`
	LayoutWidth       string     `json:"layout_width"`
	CardLayout        string     `json:"card_layout"`
	ThemeColor        string     `json:"theme_color"`
	PrimaryColor      string     `json:"primary_color"`
	SecondaryColor    string     `json:"secondary_color"`
	TertiaryColor     string     `json:"tertiary_color"`
	TopbarColor       string     `json:"topbar_color"`
	SidebarColor      string     `json:"sidebar_color"`
	TopbarBackground  string     `json:"topbar_background"`
	SidebarBackground string     `json:"sidebar_background"`
	FontFamily        string     `json:"font_family"`
	Preloader         bool       `json:"preloader"`
	ActorID           *uuid.UUID `json:"-"`
}

type UpsertTenantSettingCmd struct {
	TenantID uuid.UUID      `json:"tenant_id"`
	Key      string         `json:"key"`
	Value    map[string]any `json:"value"`
}
