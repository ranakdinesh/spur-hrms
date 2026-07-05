package domain

import (
	"errors"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	BillingCycleMonthly   = "monthly"
	BillingCycleQuarterly = "quarterly"
	BillingCycleYearly    = "yearly"
	BillingCycleOneTime   = "one_time"
	BillingCycleCustom    = "custom"

	SubscriptionPlanPriceBasisPerEmployee = "per_employee"
	SubscriptionPlanPriceBasisPackage     = "package_plus_overage"
	SubscriptionPlanPriceBasisFlat        = "flat"
	SubscriptionPlanPriceBasisCustomQuote = "custom_quote"

	SubscriptionPlanVisibilityPublic   = "public"
	SubscriptionPlanVisibilityInternal = "internal"
)

var (
	ErrInvalidSubscriptionPlanID       = errors.New("subscription plan id is required")
	ErrInvalidSubscriptionPlanCode     = errors.New("subscription plan code is required")
	ErrInvalidSubscriptionPlanName     = errors.New("subscription plan name is required")
	ErrInvalidSubscriptionPlanPrice    = errors.New("subscription plan price cannot be negative")
	ErrInvalidSubscriptionPriceBasis   = errors.New("subscription plan price basis is invalid")
	ErrInvalidSubscriptionPlanBilling  = errors.New("subscription plan billing amounts cannot be negative")
	ErrInvalidSubscriptionPlanCurrency = errors.New("subscription plan currency must be a 3-letter code")
	ErrInvalidSubscriptionBillingCycle = errors.New("subscription billing cycle is invalid")
	ErrInvalidSubscriptionPlanLimit    = errors.New("subscription employee limit cannot be negative")
	ErrInvalidSubscriptionTrialDays    = errors.New("subscription trial days cannot be negative")
	ErrInvalidSubscriptionVisibility   = errors.New("subscription plan visibility is invalid")
	ErrInactiveSubscriptionPlan        = errors.New("subscription plan is inactive")
)

type SubscriptionPlan struct {
	ID                uuid.UUID  `json:"id"`
	Code              string     `json:"code"`
	Name              string     `json:"name"`
	Description       *string    `json:"description,omitempty"`
	PriceAmount       float64    `json:"price_amount"`
	PriceBasis        string     `json:"price_basis"`
	MinimumAmount     float64    `json:"minimum_amount"`
	IncludedEmployees int32      `json:"included_employees"`
	OverageAmount     float64    `json:"overage_amount"`
	CurrencyCode      string     `json:"currency_code"`
	BillingCycle      string     `json:"billing_cycle"`
	EmployeeLimit     int32      `json:"employee_limit"`
	TrialDays         int32      `json:"trial_days"`
	Visibility        string     `json:"visibility"`
	IsActive          bool       `json:"is_active"`
	Inactive          bool       `json:"inactive"`
	CreatedAt         time.Time  `json:"created_at"`
	CreatedBy         *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt         time.Time  `json:"updated_at"`
	UpdatedBy         *uuid.UUID `json:"updated_by,omitempty"`
}

type SubscriptionPlanInput struct {
	ID                uuid.UUID
	Code              string
	Name              string
	Description       *string
	PriceAmount       float64
	PriceBasis        string
	MinimumAmount     float64
	IncludedEmployees int32
	OverageAmount     float64
	CurrencyCode      string
	BillingCycle      string
	EmployeeLimit     int32
	TrialDays         int32
	Visibility        string
	IsActive          bool
}

func NewSubscriptionPlan(input SubscriptionPlanInput) (*SubscriptionPlan, error) {
	code := strings.TrimSpace(input.Code)
	name := strings.TrimSpace(input.Name)
	if code == "" {
		return nil, ErrInvalidSubscriptionPlanCode
	}
	if name == "" {
		return nil, ErrInvalidSubscriptionPlanName
	}
	if input.PriceAmount < 0 || math.IsNaN(input.PriceAmount) || math.IsInf(input.PriceAmount, 0) {
		return nil, ErrInvalidSubscriptionPlanPrice
	}
	priceBasis := normalizeSubscriptionPlanPriceBasis(input.PriceBasis)
	if !IsValidSubscriptionPlanPriceBasis(priceBasis) {
		return nil, ErrInvalidSubscriptionPriceBasis
	}
	if input.MinimumAmount < 0 || input.OverageAmount < 0 || input.IncludedEmployees < 0 || math.IsNaN(input.MinimumAmount) || math.IsNaN(input.OverageAmount) || math.IsInf(input.MinimumAmount, 0) || math.IsInf(input.OverageAmount, 0) {
		return nil, ErrInvalidSubscriptionPlanBilling
	}
	currency := strings.ToUpper(strings.TrimSpace(input.CurrencyCode))
	if currency == "" {
		currency = "INR"
	}
	if len(currency) != 3 {
		return nil, ErrInvalidSubscriptionPlanCurrency
	}
	cycle := normalizeBillingCycle(input.BillingCycle)
	if !IsValidBillingCycle(cycle) {
		return nil, ErrInvalidSubscriptionBillingCycle
	}
	if input.EmployeeLimit < 0 {
		return nil, ErrInvalidSubscriptionPlanLimit
	}
	if input.TrialDays < 0 {
		return nil, ErrInvalidSubscriptionTrialDays
	}
	visibility := normalizeSubscriptionPlanVisibility(input.Visibility)
	if !IsValidSubscriptionPlanVisibility(visibility) {
		return nil, ErrInvalidSubscriptionVisibility
	}
	now := time.Now().UTC()
	return &SubscriptionPlan{
		ID:                input.ID,
		Code:              code,
		Name:              name,
		Description:       cleanStringPtr(input.Description),
		PriceAmount:       roundMoney(input.PriceAmount),
		PriceBasis:        priceBasis,
		MinimumAmount:     roundMoney(input.MinimumAmount),
		IncludedEmployees: input.IncludedEmployees,
		OverageAmount:     roundMoney(input.OverageAmount),
		CurrencyCode:      currency,
		BillingCycle:      cycle,
		EmployeeLimit:     input.EmployeeLimit,
		TrialDays:         input.TrialDays,
		Visibility:        visibility,
		IsActive:          input.IsActive,
		CreatedAt:         now,
		UpdatedAt:         now,
	}, nil
}

func IsValidSubscriptionPlanPriceBasis(value string) bool {
	switch value {
	case SubscriptionPlanPriceBasisPerEmployee, SubscriptionPlanPriceBasisPackage, SubscriptionPlanPriceBasisFlat, SubscriptionPlanPriceBasisCustomQuote:
		return true
	default:
		return false
	}
}

func (p *SubscriptionPlan) MonthlyAmountForEmployees(employeeCount int32) float64 {
	if p == nil || employeeCount < 0 {
		return 0
	}
	switch p.PriceBasis {
	case SubscriptionPlanPriceBasisPackage:
		overageEmployees := employeeCount - p.IncludedEmployees
		if overageEmployees < 0 {
			overageEmployees = 0
		}
		return roundMoney(p.PriceAmount + float64(overageEmployees)*p.OverageAmount)
	case SubscriptionPlanPriceBasisPerEmployee:
		total := float64(employeeCount) * p.PriceAmount
		if total < p.MinimumAmount {
			total = p.MinimumAmount
		}
		return roundMoney(total)
	case SubscriptionPlanPriceBasisFlat:
		return roundMoney(p.PriceAmount)
	default:
		return 0
	}
}

func IsValidSubscriptionPlanVisibility(value string) bool {
	switch value {
	case SubscriptionPlanVisibilityPublic, SubscriptionPlanVisibilityInternal:
		return true
	default:
		return false
	}
}

func IsValidBillingCycle(cycle string) bool {
	switch cycle {
	case BillingCycleMonthly, BillingCycleQuarterly, BillingCycleYearly, BillingCycleOneTime, BillingCycleCustom:
		return true
	default:
		return false
	}
}

func normalizeSubscriptionPlanPriceBasis(value string) string {
	value = strings.TrimSpace(strings.ToLower(value))
	if value == "" {
		return SubscriptionPlanPriceBasisPerEmployee
	}
	return value
}

func normalizeSubscriptionPlanVisibility(value string) string {
	value = strings.TrimSpace(strings.ToLower(value))
	if value == "" {
		return SubscriptionPlanVisibilityPublic
	}
	return value
}

func normalizeBillingCycle(cycle string) string {
	cycle = strings.TrimSpace(strings.ToLower(cycle))
	if cycle == "" {
		return BillingCycleMonthly
	}
	return cycle
}
