package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	SalaryTemplateTypeCTC       = "ctc"
	SalaryTemplateTypeGross     = "gross"
	SalaryTemplateTypeNet       = "net"
	SalaryTemplateTypeAllowance = "allowance"
	SalaryTemplateTypeDeduction = "deduction"

	SalaryTemplateAppliesAll          = "all"
	SalaryTemplateAppliesGrade        = "grade"
	SalaryTemplateAppliesDepartment   = "department"
	SalaryTemplateAppliesDesignation  = "designation"
	SalaryTemplateAppliesEmployeeType = "employee_type"
	SalaryTemplateAppliesCustom       = "custom"

	SalaryItemEmployerContribution = "employer_contribution"
	SalaryItemReimbursement        = "reimbursement"

	SalaryCalculationFixed      = "fixed"
	SalaryCalculationPercentage = "percentage"
	SalaryCalculationFormula    = "formula"
	SalaryCalculationManual     = "manual"

	SalaryCalculationBaseCTC     = "ctc"
	SalaryCalculationBaseGross   = "gross"
	SalaryCalculationBaseBasic   = "basic"
	SalaryCalculationBaseTaxable = "taxable"
	SalaryCalculationBaseNet     = "net"
	SalaryCalculationBaseCustom  = "custom"

	SalaryContributionEmployee = "employee"
	SalaryContributionEmployer = "employer"
	SalaryContributionNone     = "none"
)

var (
	ErrInvalidSalaryTemplate     = errors.New("invalid salary template")
	ErrInvalidSalaryTemplateID   = errors.New("salary_template_id is required")
	ErrInvalidSalaryTemplateItem = errors.New("invalid salary template item")
	ErrInvalidSalaryItemID       = errors.New("salary_template_item_id is required")
)

type SalaryTemplate struct {
	ID            uuid.UUID             `json:"id"`
	TenantID      uuid.UUID             `json:"tenant_id"`
	FYID          uuid.UUID             `json:"fy_id"`
	Code          string                `json:"code"`
	Name          string                `json:"name"`
	Description   *string               `json:"description,omitempty"`
	TemplateType  string                `json:"template_type"`
	AppliesTo     string                `json:"applies_to"`
	CurrencyCode  string                `json:"currency_code"`
	EffectiveFrom *time.Time            `json:"effective_from,omitempty"`
	EffectiveTo   *time.Time            `json:"effective_to,omitempty"`
	Notes         *string               `json:"notes,omitempty"`
	IsActive      bool                  `json:"is_active"`
	Inactive      bool                  `json:"inactive"`
	CreatedAt     time.Time             `json:"created_at"`
	CreatedBy     *uuid.UUID            `json:"created_by,omitempty"`
	UpdatedAt     time.Time             `json:"updated_at"`
	UpdatedBy     *uuid.UUID            `json:"updated_by,omitempty"`
	Items         []*SalaryTemplateItem `json:"items,omitempty"`
}

type SalaryTemplateItem struct {
	ID               uuid.UUID  `json:"id"`
	TenantID         uuid.UUID  `json:"tenant_id"`
	TemplateID       uuid.UUID  `json:"template_id"`
	ItemType         string     `json:"item_type"`
	Code             string     `json:"code"`
	Name             string     `json:"name"`
	Percentage       *float64   `json:"percentage,omitempty"`
	Amount           *float64   `json:"amount,omitempty"`
	CalculationMode  string     `json:"calculation_mode"`
	CalculationBase  string     `json:"calculation_base"`
	Formula          *string    `json:"formula,omitempty"`
	ContributionSide string     `json:"contribution_side"`
	IsTaxExempt      bool       `json:"is_tax_exempt"`
	IsStatutory      bool       `json:"is_statutory"`
	IsVariable       bool       `json:"is_variable"`
	AffectsGross     bool       `json:"affects_gross"`
	AffectsNet       bool       `json:"affects_net"`
	CapAmount        *float64   `json:"cap_amount,omitempty"`
	MinAmount        *float64   `json:"min_amount,omitempty"`
	MaxAmount        *float64   `json:"max_amount,omitempty"`
	SortOrder        int32      `json:"sort_order"`
	Inactive         bool       `json:"inactive"`
	CreatedAt        time.Time  `json:"created_at"`
	CreatedBy        *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt        time.Time  `json:"updated_at"`
	UpdatedBy        *uuid.UUID `json:"updated_by,omitempty"`
}

type SalaryTemplateInput struct {
	TenantID      uuid.UUID
	FYID          uuid.UUID
	Code          string
	Name          string
	Description   *string
	TemplateType  string
	AppliesTo     string
	CurrencyCode  string
	EffectiveFrom *time.Time
	EffectiveTo   *time.Time
	Notes         *string
}

type SalaryTemplateItemInput struct {
	TenantID         uuid.UUID
	TemplateID       uuid.UUID
	ItemType         string
	Code             string
	Name             string
	Percentage       *float64
	Amount           *float64
	CalculationMode  string
	CalculationBase  string
	Formula          *string
	ContributionSide string
	IsTaxExempt      bool
	IsStatutory      bool
	IsVariable       bool
	AffectsGross     bool
	AffectsNet       bool
	CapAmount        *float64
	MinAmount        *float64
	MaxAmount        *float64
	SortOrder        int32
}

func NewSalaryTemplate(input SalaryTemplateInput) (*SalaryTemplate, error) {
	if input.TenantID == uuid.Nil {
		return nil, ErrInvalidTenantID
	}
	if input.FYID == uuid.Nil {
		return nil, ErrInvalidFinancialYearID
	}
	name := strings.TrimSpace(input.Name)
	code := normalizeCode(input.Code)
	if name == "" {
		return nil, ErrInvalidSalaryTemplate
	}
	if code == "" {
		code = normalizeCode(name)
	}
	templateType := firstClean(input.TemplateType, SalaryTemplateTypeCTC)
	appliesTo := firstClean(input.AppliesTo, SalaryTemplateAppliesAll)
	currency := strings.ToUpper(firstClean(input.CurrencyCode, "INR"))
	if input.EffectiveFrom != nil && input.EffectiveTo != nil && dateOnly(*input.EffectiveTo).Before(dateOnly(*input.EffectiveFrom)) {
		return nil, ErrInvalidSalaryTemplate
	}
	item := &SalaryTemplate{TenantID: input.TenantID, FYID: input.FYID, Code: code, Name: name, Description: cleanOptional(input.Description), TemplateType: templateType, AppliesTo: appliesTo, CurrencyCode: currency, EffectiveFrom: salaryTemplateDatePtr(input.EffectiveFrom), EffectiveTo: salaryTemplateDatePtr(input.EffectiveTo), Notes: cleanOptional(input.Notes), CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()}
	if !validSalaryTemplate(item) {
		return nil, ErrInvalidSalaryTemplate
	}
	return item, nil
}

func NewSalaryTemplateItem(input SalaryTemplateItemInput) (*SalaryTemplateItem, error) {
	if input.TenantID == uuid.Nil {
		return nil, ErrInvalidTenantID
	}
	if input.TemplateID == uuid.Nil {
		return nil, ErrInvalidSalaryTemplateID
	}
	name := strings.TrimSpace(input.Name)
	code := normalizeCode(input.Code)
	if name == "" || code == "" {
		return nil, ErrInvalidSalaryTemplateItem
	}
	mode := firstClean(input.CalculationMode, SalaryCalculationFixed)
	base := firstClean(input.CalculationBase, SalaryCalculationBaseGross)
	side := firstClean(input.ContributionSide, SalaryContributionEmployee)
	itemType := firstClean(input.ItemType, SalaryItemEarning)
	item := &SalaryTemplateItem{TenantID: input.TenantID, TemplateID: input.TemplateID, ItemType: itemType, Code: code, Name: name, Percentage: positiveFloatPtr(input.Percentage), Amount: positiveFloatPtr(input.Amount), CalculationMode: mode, CalculationBase: base, Formula: cleanOptional(input.Formula), ContributionSide: side, IsTaxExempt: input.IsTaxExempt, IsStatutory: input.IsStatutory, IsVariable: input.IsVariable, AffectsGross: input.AffectsGross, AffectsNet: input.AffectsNet, CapAmount: positiveFloatPtr(input.CapAmount), MinAmount: positiveFloatPtr(input.MinAmount), MaxAmount: positiveFloatPtr(input.MaxAmount), SortOrder: input.SortOrder, CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()}
	if mode == SalaryCalculationFixed && item.Amount == nil {
		return nil, ErrInvalidSalaryTemplateItem
	}
	if mode == SalaryCalculationPercentage && item.Percentage == nil {
		return nil, ErrInvalidSalaryTemplateItem
	}
	if mode == SalaryCalculationFormula && item.Formula == nil {
		return nil, ErrInvalidSalaryTemplateItem
	}
	if item.MinAmount != nil && item.MaxAmount != nil && *item.MaxAmount < *item.MinAmount {
		return nil, ErrInvalidSalaryTemplateItem
	}
	if !validSalaryTemplateItem(item) {
		return nil, ErrInvalidSalaryTemplateItem
	}
	return item, nil
}

func validSalaryTemplate(item *SalaryTemplate) bool {
	return item != nil && item.Code != "" && item.Name != "" && containsString([]string{SalaryTemplateTypeCTC, SalaryTemplateTypeGross, SalaryTemplateTypeNet, SalaryTemplateTypeAllowance, SalaryTemplateTypeDeduction}, item.TemplateType) && containsString([]string{SalaryTemplateAppliesAll, SalaryTemplateAppliesGrade, SalaryTemplateAppliesDepartment, SalaryTemplateAppliesDesignation, SalaryTemplateAppliesEmployeeType, SalaryTemplateAppliesCustom}, item.AppliesTo) && len(item.CurrencyCode) == 3
}

func validSalaryTemplateItem(item *SalaryTemplateItem) bool {
	return item != nil && item.Code != "" && item.Name != "" && containsString([]string{SalaryItemEarning, SalaryItemDeduction, SalaryItemEmployerContribution, SalaryItemReimbursement}, item.ItemType) && containsString([]string{SalaryCalculationFixed, SalaryCalculationPercentage, SalaryCalculationFormula, SalaryCalculationManual}, item.CalculationMode) && containsString([]string{SalaryCalculationBaseCTC, SalaryCalculationBaseGross, SalaryCalculationBaseBasic, SalaryCalculationBaseTaxable, SalaryCalculationBaseNet, SalaryCalculationBaseCustom}, item.CalculationBase) && containsString([]string{SalaryContributionEmployee, SalaryContributionEmployer, SalaryContributionNone}, item.ContributionSide)
}

func normalizeCode(value string) string {
	value = strings.TrimSpace(strings.ToLower(value))
	if value == "" {
		return ""
	}
	var builder strings.Builder
	lastUnderscore := false
	for _, r := range value {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			builder.WriteRune(r)
			lastUnderscore = false
			continue
		}
		if !lastUnderscore {
			builder.WriteByte('_')
			lastUnderscore = true
		}
	}
	return strings.Trim(builder.String(), "_")
}

func salaryTemplateDatePtr(value *time.Time) *time.Time {
	if value == nil || value.IsZero() {
		return nil
	}
	clean := dateOnly(*value)
	return &clean
}

func positiveFloatPtr(value *float64) *float64 {
	if value == nil || *value < 0 {
		return nil
	}
	clean := *value
	return &clean
}
