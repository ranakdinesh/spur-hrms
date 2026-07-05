package domain

import (
	"errors"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	LeaveAccrualFixedYearly      = "fixed_yearly"
	LeaveAccrualMonthlyFixed     = "monthly_fixed"
	LeaveAccrualProbationMonthly = "probation_monthly"
	LeaveAccrualWorkedDays       = "worked_days"
	LeaveAccrualWorkedDayRange   = "worked_day_range"
	LeaveAccrualWorkedPercentage = "worked_percentage"
	LeaveAccrualTenureSlab       = "tenure_slab"
	LeaveAccrualPTOBank          = "pto_bank"
	LeaveAccrualCompOff          = "comp_off"
	LeaveAccrualManualAdjustment = "manual_adjustment"

	LeaveAccrualFrequencyInstant  = "instant"
	LeaveAccrualFrequencyDaily    = "daily"
	LeaveAccrualFrequencyWeekly   = "weekly"
	LeaveAccrualFrequencyBiweekly = "biweekly"
	LeaveAccrualFrequencyMonthly  = "monthly"
	LeaveAccrualFrequencyYearly   = "yearly"
	LeaveAccrualFrequencyManual   = "manual"

	LeaveProbationAny       = "any"
	LeaveProbationOnly      = "probation"
	LeaveProbationConfirmed = "confirmed"
)

var (
	ErrInvalidLeaveTemplateID        = errors.New("leave template id is required")
	ErrInvalidLeaveTemplateName      = errors.New("leave template name is required")
	ErrInvalidLeaveTemplateCode      = errors.New("leave template code is required")
	ErrInvalidLeaveTemplateRuleID    = errors.New("leave template rule id is required")
	ErrInvalidLeaveAccrualMethod     = errors.New("leave accrual method is invalid")
	ErrInvalidLeaveAccrualFrequency  = errors.New("leave accrual frequency is invalid")
	ErrInvalidLeaveProbationStatus   = errors.New("leave probation status is invalid")
	ErrInvalidLeaveAccrualConfig     = errors.New("leave accrual config is invalid")
	ErrLeaveTemplateNotFound         = errors.New("leave template not found")
	ErrLeaveTemplateRuleNotFound     = errors.New("leave template rule not found")
	ErrLeavePolicyAssignmentNotFound = errors.New("leave policy assignment not found")
)

type LeavePolicyTemplate struct {
	ID            uuid.UUID  `json:"id"`
	TenantID      *uuid.UUID `json:"tenant_id,omitempty"`
	Name          string     `json:"name"`
	Code          string     `json:"code"`
	Description   *string    `json:"description,omitempty"`
	IsSystem      bool       `json:"is_system"`
	EffectiveFrom *time.Time `json:"effective_from,omitempty"`
	EffectiveTo   *time.Time `json:"effective_to,omitempty"`
	Inactive      bool       `json:"inactive"`
	CreatedAt     time.Time  `json:"created_at"`
	CreatedBy     *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt     time.Time  `json:"updated_at"`
	UpdatedBy     *uuid.UUID `json:"updated_by,omitempty"`
}

type LeavePolicyTemplateRule struct {
	ID                        uuid.UUID      `json:"id"`
	TenantID                  uuid.UUID      `json:"tenant_id"`
	TemplateID                uuid.UUID      `json:"template_id"`
	LeaveTypeID               uuid.UUID      `json:"leave_type_id"`
	FYID                      *uuid.UUID     `json:"fy_id,omitempty"`
	EmploymentTypeID          *uuid.UUID     `json:"employment_type_id,omitempty"`
	DepartmentID              *uuid.UUID     `json:"department_id,omitempty"`
	DesignationID             *uuid.UUID     `json:"designation_id,omitempty"`
	ProbationStatus           *string        `json:"probation_status,omitempty"`
	AccrualMethod             string         `json:"accrual_method"`
	AccrualFrequency          string         `json:"accrual_frequency"`
	CreditDays                float64        `json:"credit_days"`
	CreditHours               float64        `json:"credit_hours"`
	AnnualEntitlement         float64        `json:"annual_entitlement"`
	MinWorkedDays             int32          `json:"min_worked_days"`
	MaxBalance                *float64       `json:"max_balance,omitempty"`
	CarryForwardEnabled       bool           `json:"carry_forward_enabled"`
	MaxCarryForward           float64        `json:"max_carry_forward"`
	CarryForwardExpiryMonths  int32          `json:"carry_forward_expiry_months"`
	EncashmentEnabled         bool           `json:"encashment_enabled"`
	EncashmentLimit           float64        `json:"encashment_limit"`
	EncashmentPayablePercent  float64        `json:"encashment_payable_percent"`
	NegativeBalanceAllowed    bool           `json:"negative_balance_allowed"`
	MaxNegativeBalance        float64        `json:"max_negative_balance"`
	SandwichApplicable        bool           `json:"sandwich_applicable"`
	IncludeHolidays           bool           `json:"include_holidays"`
	IncludeWeekoffs           bool           `json:"include_weekoffs"`
	RequiresDocumentAfterDays *float64       `json:"requires_document_after_days,omitempty"`
	MinRequestDays            float64        `json:"min_request_days"`
	MaxRequestDays            *float64       `json:"max_request_days,omitempty"`
	MaxRequestsPerYear        int32          `json:"max_requests_per_year"`
	AccrualDay                int32          `json:"accrual_day"`
	LapseUnutilized           bool           `json:"lapse_unutilized"`
	AllowHalfDay              bool           `json:"allow_half_day"`
	RequiresApproval          bool           `json:"requires_approval"`
	CalculationConfig         map[string]any `json:"calculation_config"`
	Priority                  int32          `json:"priority"`
	Inactive                  bool           `json:"inactive"`
	CreatedAt                 time.Time      `json:"created_at"`
	CreatedBy                 *uuid.UUID     `json:"created_by,omitempty"`
	UpdatedAt                 time.Time      `json:"updated_at"`
	UpdatedBy                 *uuid.UUID     `json:"updated_by,omitempty"`
}

type EmployeeLeavePolicyAssignment struct {
	ID            uuid.UUID  `json:"id"`
	TenantID      uuid.UUID  `json:"tenant_id"`
	UserID        uuid.UUID  `json:"user_id"`
	TemplateID    uuid.UUID  `json:"template_id"`
	FYID          *uuid.UUID `json:"fy_id,omitempty"`
	EffectiveFrom time.Time  `json:"effective_from"`
	EffectiveTo   *time.Time `json:"effective_to,omitempty"`
	Inactive      bool       `json:"inactive"`
	CreatedAt     time.Time  `json:"created_at"`
	CreatedBy     *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt     time.Time  `json:"updated_at"`
	UpdatedBy     *uuid.UUID `json:"updated_by,omitempty"`
}

type LeaveAccrualInput struct {
	JoiningDate        time.Time
	ConfirmationDate   *time.Time
	AsOfDate           time.Time
	PeriodStart        time.Time
	PeriodEnd          time.Time
	IsProbation        bool
	PayableDays        int32
	WorkedDays         int32
	OvertimeHours      float64
	HolidayWorkedHours float64
}

type LeaveAccrualResult struct {
	Days       float64 `json:"days"`
	SourceType string  `json:"source_type"`
	Remarks    string  `json:"remarks"`
}

func ValidateLeaveTemplate(template *LeavePolicyTemplate) error {
	if template == nil {
		return ErrInvalidLeaveTemplateID
	}
	if template.Name == "" {
		return ErrInvalidLeaveTemplateName
	}
	if template.Code == "" {
		return ErrInvalidLeaveTemplateCode
	}
	if !template.IsSystem && (template.TenantID == nil || *template.TenantID == uuid.Nil) {
		return ErrInvalidTenantID
	}
	return nil
}

func ValidateLeaveTemplateRule(rule *LeavePolicyTemplateRule) error {
	if rule == nil {
		return ErrInvalidLeaveTemplateRuleID
	}
	if rule.TenantID == uuid.Nil {
		return ErrInvalidTenantID
	}
	if rule.TemplateID == uuid.Nil {
		return ErrInvalidLeaveTemplateID
	}
	if rule.LeaveTypeID == uuid.Nil {
		return ErrInvalidLeavePolicyType
	}
	if !validAccrualMethod(rule.AccrualMethod) {
		return ErrInvalidLeaveAccrualMethod
	}
	if rule.AccrualFrequency == "" {
		rule.AccrualFrequency = LeaveAccrualFrequencyMonthly
	}
	if !validAccrualFrequency(rule.AccrualFrequency) {
		return ErrInvalidLeaveAccrualFrequency
	}
	if rule.ProbationStatus != nil && !validProbationStatus(*rule.ProbationStatus) {
		return ErrInvalidLeaveProbationStatus
	}
	if rule.CreditDays < 0 || rule.CreditHours < 0 || rule.AnnualEntitlement < 0 || rule.MinWorkedDays < 0 || rule.MaxCarryForward < 0 || rule.MaxNegativeBalance < 0 || rule.MinRequestDays < 0 || rule.MaxRequestsPerYear < 0 || rule.EncashmentLimit < 0 || rule.EncashmentPayablePercent < 0 || rule.EncashmentPayablePercent > 100 {
		return ErrInvalidLeaveAccrualConfig
	}
	if rule.MaxRequestDays != nil && *rule.MaxRequestDays < rule.MinRequestDays {
		return ErrInvalidLeaveAccrualConfig
	}
	if rule.AccrualDay == 0 {
		rule.AccrualDay = 1
	}
	if rule.AccrualDay < 1 || rule.AccrualDay > 31 {
		return ErrInvalidLeaveAccrualConfig
	}
	if rule.EncashmentPayablePercent == 0 && rule.EncashmentEnabled {
		rule.EncashmentPayablePercent = 100
	}
	if rule.Priority == 0 {
		rule.Priority = 100
	}
	if rule.CalculationConfig == nil {
		rule.CalculationConfig = map[string]any{}
	}
	return nil
}

func CalculateLeaveAccrual(rule *LeavePolicyTemplateRule, input LeaveAccrualInput) (LeaveAccrualResult, error) {
	if err := ValidateLeaveTemplateRule(rule); err != nil {
		return LeaveAccrualResult{}, err
	}
	if rule.ProbationStatus != nil {
		if *rule.ProbationStatus == LeaveProbationOnly && !input.IsProbation {
			return LeaveAccrualResult{SourceType: "monthly_accrual", Remarks: "rule applies to probation employees only"}, nil
		}
		if *rule.ProbationStatus == LeaveProbationConfirmed && input.IsProbation {
			return LeaveAccrualResult{SourceType: "monthly_accrual", Remarks: "rule applies to confirmed employees only"}, nil
		}
	}
	days := rule.CreditDays
	sourceType := "monthly_accrual"
	switch rule.AccrualMethod {
	case LeaveAccrualFixedYearly, LeaveAccrualPTOBank:
		sourceType = "yearly_accrual"
		days = proratedYearlyDays(rule, input, days)
	case LeaveAccrualMonthlyFixed, LeaveAccrualProbationMonthly:
		sourceType = "monthly_accrual"
	case LeaveAccrualWorkedDays:
		if input.WorkedDays < rule.MinWorkedDays {
			days = 0
		} else if rule.MinWorkedDays > 0 {
			days = math.Floor(float64(input.WorkedDays)/float64(rule.MinWorkedDays)) * rule.CreditDays
		}
	case LeaveAccrualWorkedPercentage:
		if input.PayableDays <= 0 {
			days = 0
		} else {
			percentage := (float64(input.WorkedDays) / float64(input.PayableDays)) * 100
			days = workedPercentageDays(rule.CalculationConfig, percentage, rule.CreditDays)
		}
	case LeaveAccrualWorkedDayRange:
		days = workedRangeDays(rule.CalculationConfig, input.WorkedDays, rule.CreditDays)
	case LeaveAccrualTenureSlab:
		days = tenureSlabDays(rule.CalculationConfig, input.JoiningDate, input.AsOfDate, rule.CreditDays)
	case LeaveAccrualCompOff:
		sourceType = "comp_off"
		multiplier := floatFromConfig(rule.CalculationConfig, "multiplier", 1)
		hoursPerDay := floatFromConfig(rule.CalculationConfig, "hours_per_day", 8)
		hours := input.HolidayWorkedHours
		if hours == 0 {
			hours = input.OvertimeHours
		}
		if hoursPerDay <= 0 {
			return LeaveAccrualResult{}, ErrInvalidLeaveAccrualConfig
		}
		days = (hours * multiplier) / hoursPerDay
	case LeaveAccrualManualAdjustment:
		sourceType = "manual_adjustment"
	}
	if rule.MaxBalance != nil && *rule.MaxBalance >= 0 && days > *rule.MaxBalance {
		days = *rule.MaxBalance
	}
	return LeaveAccrualResult{Days: math.Round(days*100) / 100, SourceType: sourceType, Remarks: rule.AccrualMethod}, nil
}

func validAccrualMethod(value string) bool {
	switch value {
	case LeaveAccrualFixedYearly, LeaveAccrualMonthlyFixed, LeaveAccrualProbationMonthly, LeaveAccrualWorkedDays, LeaveAccrualWorkedDayRange, LeaveAccrualWorkedPercentage, LeaveAccrualTenureSlab, LeaveAccrualPTOBank, LeaveAccrualCompOff, LeaveAccrualManualAdjustment:
		return true
	default:
		return false
	}
}

func validAccrualFrequency(value string) bool {
	switch value {
	case LeaveAccrualFrequencyInstant, LeaveAccrualFrequencyDaily, LeaveAccrualFrequencyWeekly, LeaveAccrualFrequencyBiweekly, LeaveAccrualFrequencyMonthly, LeaveAccrualFrequencyYearly, LeaveAccrualFrequencyManual:
		return true
	default:
		return false
	}
}

func validProbationStatus(value string) bool {
	switch value {
	case LeaveProbationAny, LeaveProbationOnly, LeaveProbationConfirmed:
		return true
	default:
		return false
	}
}

func floatFromConfig(config map[string]any, key string, fallback float64) float64 {
	if config == nil {
		return fallback
	}
	switch value := config[key].(type) {
	case float64:
		return value
	case int:
		return float64(value)
	case int32:
		return float64(value)
	case int64:
		return float64(value)
	default:
		return fallback
	}
}

func boolFromConfig(config map[string]any, key string, fallback bool) bool {
	if config == nil {
		return fallback
	}
	switch value := config[key].(type) {
	case bool:
		return value
	case string:
		switch strings.ToLower(strings.TrimSpace(value)) {
		case "true", "yes", "1", "on":
			return true
		case "false", "no", "0", "off":
			return false
		default:
			return fallback
		}
	default:
		return fallback
	}
}

func stringFromConfig(config map[string]any, key string, fallback string) string {
	if config == nil {
		return fallback
	}
	value, ok := config[key].(string)
	if !ok {
		return fallback
	}
	value = strings.TrimSpace(value)
	if value == "" {
		return fallback
	}
	return value
}

func proratedYearlyDays(rule *LeavePolicyTemplateRule, input LeaveAccrualInput, fallback float64) float64 {
	if !boolFromConfig(rule.CalculationConfig, "prorate", false) {
		return fallback
	}
	periodStart := dateOnly(input.PeriodStart)
	periodEnd := dateOnly(input.PeriodEnd)
	if periodStart.IsZero() || periodEnd.IsZero() || periodEnd.Before(periodStart) {
		return fallback
	}
	fromDate := periodStart
	switch strings.ToLower(stringFromConfig(rule.CalculationConfig, "prorate_basis", "joining_date")) {
	case "confirmation_date", "confirmed_date":
		if input.ConfirmationDate != nil && !input.ConfirmationDate.IsZero() {
			fromDate = dateOnly(*input.ConfirmationDate)
		} else if input.JoiningDate.After(periodStart) {
			fromDate = dateOnly(input.JoiningDate)
		}
	case "joining_date":
		if input.JoiningDate.After(periodStart) {
			fromDate = dateOnly(input.JoiningDate)
		}
	case "cycle_start", "period_start":
		fromDate = periodStart
	default:
		if input.JoiningDate.After(periodStart) {
			fromDate = dateOnly(input.JoiningDate)
		}
	}
	if fromDate.Before(periodStart) {
		fromDate = periodStart
	}
	if fromDate.After(periodEnd) {
		return 0
	}
	totalDays := periodEnd.Sub(periodStart).Hours()/24 + 1
	eligibleDays := periodEnd.Sub(fromDate).Hours()/24 + 1
	if totalDays <= 0 || eligibleDays <= 0 {
		return 0
	}
	return roundLeaveDays((fallback*eligibleDays)/totalDays, stringFromConfig(rule.CalculationConfig, "rounding", "none"))
}

func roundLeaveDays(days float64, mode string) float64 {
	switch strings.ToLower(strings.TrimSpace(mode)) {
	case "nearest_half":
		return math.Round(days*2) / 2
	case "floor_half":
		return math.Floor(days*2) / 2
	case "ceil_half":
		return math.Ceil(days*2) / 2
	case "nearest_day", "nearest":
		return math.Round(days)
	case "floor_day", "floor":
		return math.Floor(days)
	case "ceil_day", "ceil":
		return math.Ceil(days)
	default:
		return days
	}
}

func workedPercentageDays(config map[string]any, percentage float64, fallback float64) float64 {
	return rangeConfigDays(config, percentage, fallback)
}

func workedRangeDays(config map[string]any, workedDays int32, fallback float64) float64 {
	return rangeConfigDays(config, float64(workedDays), fallback)
}

func rangeConfigDays(config map[string]any, value float64, fallback float64) float64 {
	rawRanges, ok := config["ranges"].([]any)
	if !ok {
		return fallback
	}
	for _, raw := range rawRanges {
		rangeMap, ok := raw.(map[string]any)
		if !ok {
			continue
		}
		minValue := floatFromConfig(rangeMap, "min", 0)
		maxValue := floatFromConfig(rangeMap, "max", math.MaxFloat64)
		if value >= minValue && value <= maxValue {
			return floatFromConfig(rangeMap, "days", fallback)
		}
	}
	return 0
}

func tenureSlabDays(config map[string]any, joiningDate time.Time, asOfDate time.Time, fallback float64) float64 {
	if joiningDate.IsZero() || asOfDate.IsZero() || asOfDate.Before(joiningDate) {
		return fallback
	}
	years := asOfDate.Sub(joiningDate).Hours() / 24 / 365
	return rangeConfigDays(config, years, fallback)
}
