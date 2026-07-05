package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

type TenantRegistrationEmailSender interface {
	SendTenantRegistration(ctx context.Context, message TenantRegistrationEmail) error
}

type TenantRegistrationEmail struct {
	TenantID             uuid.UUID `json:"tenant_id"`
	CompanyName          string    `json:"company_name"`
	Subdomain            string    `json:"subdomain"`
	MobileActivationCode string    `json:"mobile_activation_code"`
	AdminEmail           string    `json:"admin_email"`
	AdminName            string    `json:"admin_name"`
	TenantURL            string    `json:"tenant_url"`
}

type ProvisionTenantCommand struct {
	TenantID             uuid.UUID  `json:"tenant_id"`
	CompanyName          string     `json:"company_name"`
	Subdomain            string     `json:"subdomain"`
	MobileActivationCode string     `json:"mobile_activation_code"`
	AdminEmail           *string    `json:"admin_email,omitempty"`
	AdminName            *string    `json:"admin_name,omitempty"`
	TenantURL            *string    `json:"tenant_url,omitempty"`
	ActorID              *uuid.UUID `json:"-"`
}

type TenantProvisioningStatus struct {
	TenantID                 uuid.UUID                  `json:"tenant_id"`
	ProfileReady             bool                       `json:"profile_ready"`
	Subdomain                string                     `json:"subdomain,omitempty"`
	MobileActivationCode     string                     `json:"mobile_activation_code,omitempty"`
	WorkingHoursCount        int                        `json:"working_hours_count"`
	FinancialYearsCount      int                        `json:"financial_years_count"`
	ActiveFinancialYearReady bool                       `json:"active_financial_year_ready"`
	EmploymentTypesCount     int                        `json:"employment_types_count"`
	MaritalStatusesCount     int                        `json:"marital_statuses_count"`
	LevelCodesCount          int                        `json:"level_codes_count"`
	SeniorityRanksCount      int                        `json:"seniority_ranks_count"`
	NotificationMastersCount int                        `json:"notification_masters_count"`
	PayCycleReady            bool                       `json:"pay_cycle_ready"`
	CurrentSubscription      *domain.TenantSubscription `json:"current_subscription,omitempty"`
}
