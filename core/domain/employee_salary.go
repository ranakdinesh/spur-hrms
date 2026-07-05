package domain

import (
	"errors"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidEmployeeSalary          = errors.New("invalid employee salary")
	ErrInvalidEmployeeSalaryID        = errors.New("employee_salary_id is required")
	ErrInvalidEmployeeSalaryStructure = errors.New("invalid employee salary structure")
)

type EmployeeSalary struct {
	ID            uuid.UUID                  `json:"id"`
	TenantID      uuid.UUID                  `json:"tenant_id"`
	UserID        uuid.UUID                  `json:"user_id"`
	FYID          uuid.UUID                  `json:"fy_id"`
	TemplateID    uuid.UUID                  `json:"template_id"`
	GrossSalary   float64                    `json:"gross_salary"`
	EffectiveFrom *time.Time                 `json:"effective_from,omitempty"`
	Inactive      bool                       `json:"inactive"`
	CreatedAt     time.Time                  `json:"created_at"`
	CreatedBy     *uuid.UUID                 `json:"created_by,omitempty"`
	UpdatedAt     time.Time                  `json:"updated_at"`
	UpdatedBy     *uuid.UUID                 `json:"updated_by,omitempty"`
	Template      *SalaryTemplate            `json:"template,omitempty"`
	Structures    []*EmployeeSalaryStructure `json:"structures,omitempty"`
}

type EmployeeSalaryStructure struct {
	ID         uuid.UUID  `json:"id"`
	TenantID   uuid.UUID  `json:"tenant_id"`
	UserID     uuid.UUID  `json:"user_id"`
	TemplateID uuid.UUID  `json:"template_id"`
	FYID       uuid.UUID  `json:"fy_id"`
	ItemType   string     `json:"item_type"`
	Code       string     `json:"code"`
	Name       string     `json:"name"`
	Amount     float64    `json:"amount"`
	SortOrder  int32      `json:"sort_order"`
	Inactive   bool       `json:"inactive"`
	CreatedAt  time.Time  `json:"created_at"`
	CreatedBy  *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt  time.Time  `json:"updated_at"`
	UpdatedBy  *uuid.UUID `json:"updated_by,omitempty"`
}

type EmployeeSalaryInput struct {
	TenantID      uuid.UUID
	UserID        uuid.UUID
	FYID          uuid.UUID
	TemplateID    uuid.UUID
	GrossSalary   float64
	EffectiveFrom *time.Time
}

func NewEmployeeSalary(input EmployeeSalaryInput) (*EmployeeSalary, error) {
	if input.TenantID == uuid.Nil {
		return nil, ErrInvalidTenantID
	}
	if input.UserID == uuid.Nil {
		return nil, ErrInvalidEmployeeUserID
	}
	if input.FYID == uuid.Nil {
		return nil, ErrInvalidFinancialYearID
	}
	if input.TemplateID == uuid.Nil {
		return nil, ErrInvalidSalaryTemplateID
	}
	if input.GrossSalary < 0 || math.IsNaN(input.GrossSalary) || math.IsInf(input.GrossSalary, 0) {
		return nil, ErrInvalidEmployeeSalary
	}
	item := &EmployeeSalary{TenantID: input.TenantID, UserID: input.UserID, FYID: input.FYID, TemplateID: input.TemplateID, GrossSalary: roundMoney(input.GrossSalary), EffectiveFrom: salaryTemplateDatePtr(input.EffectiveFrom), CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()}
	return item, nil
}

func NewEmployeeSalaryStructure(tenantID uuid.UUID, userID uuid.UUID, fyID uuid.UUID, templateID uuid.UUID, templateItem *SalaryTemplateItem, amount float64) (*EmployeeSalaryStructure, error) {
	if tenantID == uuid.Nil || userID == uuid.Nil || fyID == uuid.Nil || templateID == uuid.Nil || templateItem == nil {
		return nil, ErrInvalidEmployeeSalaryStructure
	}
	itemType := strings.TrimSpace(strings.ToLower(templateItem.ItemType))
	if itemType == SalaryItemEmployerContribution || itemType == SalaryItemReimbursement {
		itemType = SalaryItemEarning
	}
	if itemType != SalaryItemEarning && itemType != SalaryItemDeduction {
		return nil, ErrInvalidEmployeeSalaryStructure
	}
	if templateItem.Code == "" || templateItem.Name == "" || amount < 0 || math.IsNaN(amount) || math.IsInf(amount, 0) {
		return nil, ErrInvalidEmployeeSalaryStructure
	}
	return &EmployeeSalaryStructure{TenantID: tenantID, UserID: userID, FYID: fyID, TemplateID: templateID, ItemType: itemType, Code: templateItem.Code, Name: templateItem.Name, Amount: roundMoney(amount), SortOrder: templateItem.SortOrder, CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()}, nil
}

func roundMoney(value float64) float64 {
	return math.Round(value*100) / 100
}
