package domain

import (
	"errors"
	"net"
	"net/url"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

var (
	ErrInvalidBrandingColor  = errors.New("branding colors must be hex values like #588368")
	ErrInvalidBrandingOption = errors.New("branding option is not supported")
	ErrInvalidBrandingHost   = errors.New("branding host is required")
)

var hexColorPattern = regexp.MustCompile(`^#[0-9A-Fa-f]{6}$`)

type TenantBranding struct {
	TenantID          uuid.UUID `json:"tenant_id"`
	Subdomain         string    `json:"subdomain"`
	DisplayName       *string   `json:"display_name,omitempty"`
	LogoPath          *string   `json:"logo_path,omitempty"`
	FaviconPath       *string   `json:"favicon_path,omitempty"`
	Layout            string    `json:"layout"`
	ColorMode         string    `json:"color_mode"`
	SidebarSize       string    `json:"sidebar_size"`
	LayoutWidth       string    `json:"layout_width"`
	CardLayout        string    `json:"card_layout"`
	ThemeColor        string    `json:"theme_color"`
	PrimaryColor      string    `json:"primary_color"`
	SecondaryColor    string    `json:"secondary_color"`
	TertiaryColor     string    `json:"tertiary_color"`
	TopbarColor       string    `json:"topbar_color"`
	SidebarColor      string    `json:"sidebar_color"`
	TopbarBackground  string    `json:"topbar_background"`
	SidebarBackground string    `json:"sidebar_background"`
	FontFamily        string    `json:"font_family"`
	Preloader         bool      `json:"preloader"`
}

type TenantBrandingInput struct {
	TenantID          uuid.UUID
	Subdomain         string
	DisplayName       *string
	LogoPath          *string
	FaviconPath       *string
	Layout            string
	ColorMode         string
	SidebarSize       string
	LayoutWidth       string
	CardLayout        string
	ThemeColor        string
	PrimaryColor      string
	SecondaryColor    string
	TertiaryColor     string
	TopbarColor       string
	SidebarColor      string
	TopbarBackground  string
	SidebarBackground string
	FontFamily        string
	Preloader         bool
}

func NewTenantBranding(input TenantBrandingInput) (*TenantBranding, error) {
	if input.TenantID == uuid.Nil {
		return nil, ErrInvalidTenantID
	}
	subdomain := NormalizeSubdomain(input.Subdomain)
	if !subdomainPattern.MatchString(subdomain) {
		return nil, ErrInvalidSubdomain
	}
	if IsReservedTenantSubdomain(subdomain) || TenantSubdomainCollisionKey(subdomain) == "" {
		return nil, ErrReservedSubdomain
	}
	if HasTenantSubdomainBusinessSuffix(subdomain) {
		return nil, ErrConfusingSubdomain
	}
	branding := &TenantBranding{
		TenantID:          input.TenantID,
		Subdomain:         subdomain,
		DisplayName:       cleanOptional(input.DisplayName),
		LogoPath:          cleanOptional(input.LogoPath),
		FaviconPath:       cleanOptional(input.FaviconPath),
		Layout:            normalizeBrandingOption(input.Layout, "vertical"),
		ColorMode:         normalizeBrandingOption(input.ColorMode, "light"),
		SidebarSize:       normalizeBrandingOption(input.SidebarSize, "default"),
		LayoutWidth:       normalizeBrandingOption(input.LayoutWidth, "fluid"),
		CardLayout:        normalizeBrandingOption(input.CardLayout, "bordered"),
		ThemeColor:        normalizeBrandingOption(input.ThemeColor, "#588368"),
		PrimaryColor:      normalizeBrandingOption(input.PrimaryColor, "#588368"),
		SecondaryColor:    normalizeBrandingOption(input.SecondaryColor, "#2f6f7d"),
		TertiaryColor:     normalizeBrandingOption(input.TertiaryColor, "#e87839"),
		TopbarColor:       normalizeBrandingOption(input.TopbarColor, "#ffffff"),
		SidebarColor:      normalizeBrandingOption(input.SidebarColor, "#111827"),
		TopbarBackground:  normalizeBrandingOption(input.TopbarBackground, "none"),
		SidebarBackground: normalizeBrandingOption(input.SidebarBackground, "solid"),
		FontFamily:        normalizeBrandingOption(input.FontFamily, "Inter, sans-serif"),
		Preloader:         input.Preloader,
	}
	if !validBrandingOptions(branding) {
		return nil, ErrInvalidBrandingOption
	}
	if !validBrandingColors(branding) {
		return nil, ErrInvalidBrandingColor
	}
	return branding, nil
}

func BrandingSubdomainFromHost(value string) (string, error) {
	host := strings.TrimSpace(strings.ToLower(value))
	if host == "" {
		return "", ErrInvalidBrandingHost
	}
	if strings.Contains(host, "://") {
		parsed, err := url.Parse(host)
		if err != nil {
			return "", err
		}
		host = parsed.Host
	}
	if parsedHost, _, err := net.SplitHostPort(host); err == nil {
		host = parsedHost
	}
	host = strings.Trim(host, ".")
	if host == "" || net.ParseIP(host) != nil {
		return "", ErrInvalidBrandingHost
	}
	parts := strings.Split(host, ".")
	if len(parts) == 0 || parts[0] == "" || parts[0] == "www" || parts[0] == "localhost" {
		return "", ErrInvalidBrandingHost
	}
	subdomain := NormalizeSubdomain(parts[0])
	if !subdomainPattern.MatchString(subdomain) {
		return "", ErrInvalidSubdomain
	}
	return subdomain, nil
}

func normalizeBrandingOption(value string, fallback string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return fallback
	}
	return value
}

func validBrandingColors(branding *TenantBranding) bool {
	return hexColorPattern.MatchString(branding.ThemeColor) &&
		hexColorPattern.MatchString(branding.PrimaryColor) &&
		hexColorPattern.MatchString(branding.SecondaryColor) &&
		hexColorPattern.MatchString(branding.TertiaryColor) &&
		hexColorPattern.MatchString(branding.TopbarColor) &&
		hexColorPattern.MatchString(branding.SidebarColor)
}

func validBrandingOptions(branding *TenantBranding) bool {
	return inSet(branding.Layout, "vertical", "horizontal", "detached", "two-column") &&
		inSet(branding.ColorMode, "light", "dark", "system") &&
		inSet(branding.SidebarSize, "default", "compact", "condensed", "icon") &&
		inSet(branding.LayoutWidth, "fluid", "boxed") &&
		inSet(branding.CardLayout, "bordered", "shadow", "plain")
}

func inSet(value string, allowed ...string) bool {
	for _, item := range allowed {
		if value == item {
			return true
		}
	}
	return false
}
