package domain

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidTenantCompanyName = errors.New("company name is required")
	ErrInvalidRegistrationEmail = errors.New("registration email is invalid")
)

var (
	slugInvalidPattern = regexp.MustCompile(`[^a-z0-9]+`)
	emailLikePattern   = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)
)

type TenantProvisioning struct {
	TenantID             uuid.UUID      `json:"tenant_id"`
	CompanyName          string         `json:"company_name"`
	Subdomain            string         `json:"subdomain"`
	MobileActivationCode string         `json:"mobile_activation_code"`
	AdminEmail           *string        `json:"admin_email,omitempty"`
	AdminName            *string        `json:"admin_name,omitempty"`
	TenantURL            *string        `json:"tenant_url,omitempty"`
	Seeded               map[string]int `json:"seeded"`
	EmailSent            bool           `json:"email_sent"`
	CreatedAt            time.Time      `json:"created_at"`
}

type TenantProvisioningInput struct {
	TenantID             uuid.UUID
	CompanyName          string
	Subdomain            string
	MobileActivationCode string
	AdminEmail           *string
	AdminName            *string
	TenantURL            *string
}

func NewTenantProvisioning(input TenantProvisioningInput) (*TenantProvisioning, error) {
	if input.TenantID == uuid.Nil {
		return nil, ErrInvalidTenantID
	}
	companyName := strings.TrimSpace(input.CompanyName)
	if companyName == "" {
		return nil, ErrInvalidTenantCompanyName
	}
	subdomain := NormalizeSubdomain(input.Subdomain)
	if subdomain == "" {
		subdomain = TenantSubdomainCollisionKey(SlugFromTenantName(companyName, input.TenantID))
	}
	if subdomain == "" {
		subdomain = SlugFromTenantName(companyName, input.TenantID)
	}
	if !subdomainPattern.MatchString(subdomain) {
		return nil, ErrInvalidSubdomain
	}
	if IsReservedTenantSubdomain(subdomain) || TenantSubdomainCollisionKey(subdomain) == "" {
		return nil, ErrReservedSubdomain
	}
	if HasTenantSubdomainBusinessSuffix(subdomain) {
		return nil, ErrConfusingSubdomain
	}
	activationCode := NormalizeMobileActivationCode(input.MobileActivationCode)
	if activationCode == "" {
		activationCode = ActivationCodeFromTenantID(input.TenantID)
	}
	if !mobileActivationCodePattern.MatchString(activationCode) {
		return nil, ErrInvalidMobileActivationCode
	}
	adminEmail := cleanOptional(input.AdminEmail)
	if adminEmail != nil && !emailLikePattern.MatchString(*adminEmail) {
		return nil, ErrInvalidRegistrationEmail
	}
	return &TenantProvisioning{
		TenantID:             input.TenantID,
		CompanyName:          companyName,
		Subdomain:            subdomain,
		MobileActivationCode: activationCode,
		AdminEmail:           adminEmail,
		AdminName:            cleanOptional(input.AdminName),
		TenantURL:            cleanOptional(input.TenantURL),
		Seeded:               map[string]int{},
		CreatedAt:            time.Now().UTC(),
	}, nil
}

func SlugFromTenantName(name string, tenantID uuid.UUID) string {
	slug := strings.ToLower(strings.TrimSpace(name))
	slug = slugInvalidPattern.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	if slug == "" {
		slug = strings.ToLower(strings.ReplaceAll(tenantID.String(), "-", ""))
	}
	if len(slug) > 63 {
		slug = strings.Trim(slug[:63], "-")
	}
	if slug == "" {
		return "tenant"
	}
	return slug
}

func ActivationCodeFromTenantID(tenantID uuid.UUID) string {
	clean := strings.ToUpper(strings.ReplaceAll(tenantID.String(), "-", ""))
	if len(clean) > 8 {
		clean = clean[:8]
	}
	if clean == "" {
		clean = "TENANT"
	}
	return "HRMS-" + clean
}
