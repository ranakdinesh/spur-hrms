package domain

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidTenantID             = errors.New("tenant_id is required")
	ErrInvalidSubdomain            = errors.New("subdomain must contain only lowercase letters, numbers, and hyphens")
	ErrReservedSubdomain           = errors.New("subdomain is reserved")
	ErrConfusingSubdomain          = errors.New("subdomain must be a short tenant code and cannot include business suffixes like infotech, pvt, ltd, services, or solutions")
	ErrConflictingSubdomain        = errors.New("subdomain conflicts with an existing tenant")
	ErrInvalidMobileActivationCode = errors.New("mobile activation code must contain only uppercase letters, numbers, and hyphens")
	ErrInvalidSettingKey           = errors.New("setting key must start with a letter and contain only lowercase letters, numbers, dots, underscores, and hyphens")
)

var (
	subdomainPattern            = regexp.MustCompile(`^[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?$`)
	mobileActivationCodePattern = regexp.MustCompile(`^[A-Z0-9-]{4,32}$`)
	settingKeyPattern           = regexp.MustCompile(`^[a-z][a-z0-9_.-]*$`)
)

type TenantProfile struct {
	TenantID             uuid.UUID `json:"tenant_id"`
	Subdomain            string    `json:"subdomain"`
	MobileActivationCode string    `json:"mobile_activation_code"`
	DisplayName          *string   `json:"display_name,omitempty"`
	LogoObjectKey        *string   `json:"logo_object_key,omitempty"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

type TenantSetting struct {
	TenantID  uuid.UUID      `json:"tenant_id"`
	Key       string         `json:"key"`
	Value     map[string]any `json:"value"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

func NewTenantProfile(tenantID uuid.UUID, subdomain, mobileActivationCode string, displayName, logoObjectKey *string) (*TenantProfile, error) {
	if tenantID == uuid.Nil {
		return nil, ErrInvalidTenantID
	}
	subdomain = NormalizeSubdomain(subdomain)
	if !subdomainPattern.MatchString(subdomain) {
		return nil, ErrInvalidSubdomain
	}
	if IsReservedTenantSubdomain(subdomain) || TenantSubdomainCollisionKey(subdomain) == "" {
		return nil, ErrReservedSubdomain
	}
	if HasTenantSubdomainBusinessSuffix(subdomain) {
		return nil, ErrConfusingSubdomain
	}
	mobileActivationCode = NormalizeMobileActivationCode(mobileActivationCode)
	if !mobileActivationCodePattern.MatchString(mobileActivationCode) {
		return nil, ErrInvalidMobileActivationCode
	}

	now := time.Now().UTC()
	return &TenantProfile{
		TenantID:             tenantID,
		Subdomain:            subdomain,
		MobileActivationCode: mobileActivationCode,
		DisplayName:          cleanOptional(displayName),
		LogoObjectKey:        cleanOptional(logoObjectKey),
		CreatedAt:            now,
		UpdatedAt:            now,
	}, nil
}

func NormalizeSubdomain(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func IsReservedTenantSubdomain(value string) bool {
	_, reserved := reservedTenantSubdomains[NormalizeSubdomain(value)]
	return reserved
}

func TenantSubdomainCollisionKey(value string) string {
	value = NormalizeSubdomain(value)
	if value == "" {
		return ""
	}
	parts := strings.FieldsFunc(value, func(r rune) bool {
		return r == '-' || r == '_' || r == '.' || r == ' ' || r == '\t' || r == '\n'
	})
	cleaned := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			cleaned = append(cleaned, part)
		}
	}
	for len(cleaned) > 0 {
		if _, ignored := ignoredTenantSubdomainSuffixes[cleaned[len(cleaned)-1]]; !ignored {
			break
		}
		cleaned = cleaned[:len(cleaned)-1]
	}
	return strings.Join(cleaned, "")
}

func HasTenantSubdomainBusinessSuffix(value string) bool {
	value = NormalizeSubdomain(value)
	parts := strings.FieldsFunc(value, func(r rune) bool {
		return r == '-' || r == '_' || r == '.' || r == ' ' || r == '\t' || r == '\n'
	})
	if len(parts) < 2 {
		return false
	}
	last := strings.TrimSpace(parts[len(parts)-1])
	_, ignored := ignoredTenantSubdomainSuffixes[last]
	return ignored
}

var reservedTenantSubdomains = map[string]struct{}{
	"admin":    {},
	"api":      {},
	"app":      {},
	"auth":     {},
	"billing":  {},
	"dev":      {},
	"files":    {},
	"identity": {},
	"mail":     {},
	"staging":  {},
	"support":  {},
	"www":      {},
}

var ignoredTenantSubdomainSuffixes = map[string]struct{}{
	"co":           {},
	"company":      {},
	"consultants":  {},
	"consulting":   {},
	"corp":         {},
	"corporation":  {},
	"digital":      {},
	"enterprise":   {},
	"enterprises":  {},
	"global":       {},
	"india":        {},
	"inc":          {},
	"infotech":     {},
	"it":           {},
	"limited":      {},
	"llp":          {},
	"ltd":          {},
	"private":      {},
	"pvt":          {},
	"service":      {},
	"services":     {},
	"solution":     {},
	"solutions":    {},
	"software":     {},
	"system":       {},
	"systems":      {},
	"tech":         {},
	"technologies": {},
	"technology":   {},
}

func NormalizeMobileActivationCode(value string) string {
	return strings.ToUpper(strings.TrimSpace(value))
}

func ValidateSettingKey(key string) (string, error) {
	key = strings.TrimSpace(key)
	if !settingKeyPattern.MatchString(key) {
		return "", ErrInvalidSettingKey
	}
	return key, nil
}

func cleanOptional(value *string) *string {
	if value == nil {
		return nil
	}
	clean := strings.TrimSpace(*value)
	if clean == "" {
		return nil
	}
	return &clean
}

func cleanUUIDOptional(value *uuid.UUID) *uuid.UUID {
	if value == nil || *value == uuid.Nil {
		return nil
	}
	clean := *value
	return &clean
}
