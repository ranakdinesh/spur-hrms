package domain

import (
	"encoding/json"
	"errors"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidPayrollPeriod         = errors.New("payroll period is invalid")
	ErrPayrollPeriodLocked          = errors.New("payroll period is locked")
	ErrInvalidPayrollStatutoryRule  = errors.New("payroll statutory rule is invalid")
	ErrPayrollStatutoryRuleNotFound = errors.New("payroll statutory rule not found")
	ErrPayrollImportNotFound        = errors.New("payroll import not found")
)

const (
	PayrollLockStatusOpen     = "open"
	PayrollLockStatusLocked   = "locked"
	PayrollLockStatusUnlocked = "unlocked"

	PayrollRulePT  = "pt"
	PayrollRuleLWF = "lwf"

	PayrollImportSalaryRevision = "salary_revision"
	PayrollImportAttendanceLOP  = "attendance_lop"
	PayrollImportVariablePay    = "variable_pay"
	PayrollImportAdjustment     = "adjustment"

	PayrollImportStatusValidated = "validated"
	PayrollImportStatusApplied   = "applied"
	PayrollImportStatusFailed    = "failed"
	PayrollImportStatusPartial   = "partial"

	PayrollImportRowValid   = "valid"
	PayrollImportRowInvalid = "invalid"
	PayrollImportRowApplied = "applied"
	PayrollImportRowSkipped = "skipped"
)

type PayrollPeriodLock struct {
	ID           uuid.UUID                 `json:"id"`
	TenantID     uuid.UUID                 `json:"tenant_id"`
	Month        int32                     `json:"month"`
	Year         int32                     `json:"year"`
	Status       string                    `json:"status"`
	LockedAt     *time.Time                `json:"locked_at,omitempty"`
	LockedBy     *uuid.UUID                `json:"locked_by,omitempty"`
	UnlockedAt   *time.Time                `json:"unlocked_at,omitempty"`
	UnlockedBy   *uuid.UUID                `json:"unlocked_by,omitempty"`
	UnlockReason *string                   `json:"unlock_reason,omitempty"`
	Notes        *string                   `json:"notes,omitempty"`
	Inactive     bool                      `json:"inactive"`
	CreatedAt    time.Time                 `json:"created_at"`
	CreatedBy    *uuid.UUID                `json:"created_by,omitempty"`
	UpdatedAt    time.Time                 `json:"updated_at"`
	UpdatedBy    *uuid.UUID                `json:"updated_by,omitempty"`
	Events       []*PayrollPeriodLockEvent `json:"events,omitempty"`
}

type PayrollPeriodLockEvent struct {
	ID            uuid.UUID  `json:"id"`
	TenantID      uuid.UUID  `json:"tenant_id"`
	PayrollLockID uuid.UUID  `json:"payroll_lock_id"`
	Action        string     `json:"action"`
	FromStatus    *string    `json:"from_status,omitempty"`
	ToStatus      *string    `json:"to_status,omitempty"`
	Remarks       *string    `json:"remarks,omitempty"`
	Inactive      bool       `json:"inactive"`
	CreatedAt     time.Time  `json:"created_at"`
	CreatedBy     *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt     time.Time  `json:"updated_at"`
	UpdatedBy     *uuid.UUID `json:"updated_by,omitempty"`
}

type PayrollStatutoryRule struct {
	ID             uuid.UUID  `json:"id"`
	TenantID       uuid.UUID  `json:"tenant_id"`
	RuleType       string     `json:"rule_type"`
	Name           string     `json:"name"`
	State          *string    `json:"state,omitempty"`
	BranchID       *uuid.UUID `json:"branch_id,omitempty"`
	EffectiveFrom  time.Time  `json:"effective_from"`
	EffectiveTo    *time.Time `json:"effective_to,omitempty"`
	MinGrossSalary *float64   `json:"min_gross_salary,omitempty"`
	MaxGrossSalary *float64   `json:"max_gross_salary,omitempty"`
	EmployeeAmount float64    `json:"employee_amount"`
	EmployerAmount float64    `json:"employer_amount"`
	Frequency      string     `json:"frequency"`
	DeductionMonth *int32     `json:"deduction_month,omitempty"`
	Notes          *string    `json:"notes,omitempty"`
	Inactive       bool       `json:"inactive"`
	CreatedAt      time.Time  `json:"created_at"`
	CreatedBy      *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt      time.Time  `json:"updated_at"`
	UpdatedBy      *uuid.UUID `json:"updated_by,omitempty"`
}

type PayrollImportBatch struct {
	ID          uuid.UUID           `json:"id"`
	TenantID    uuid.UUID           `json:"tenant_id"`
	ImportType  string              `json:"import_type"`
	Month       *int32              `json:"month,omitempty"`
	Year        *int32              `json:"year,omitempty"`
	FYID        *uuid.UUID          `json:"fy_id,omitempty"`
	TemplateID  *uuid.UUID          `json:"template_id,omitempty"`
	FileName    *string             `json:"file_name,omitempty"`
	Status      string              `json:"status"`
	TotalRows   int32               `json:"total_rows"`
	ValidRows   int32               `json:"valid_rows"`
	InvalidRows int32               `json:"invalid_rows"`
	AppliedRows int32               `json:"applied_rows"`
	ErrorReport json.RawMessage     `json:"error_report,omitempty"`
	Notes       *string             `json:"notes,omitempty"`
	Inactive    bool                `json:"inactive"`
	CreatedAt   time.Time           `json:"created_at"`
	CreatedBy   *uuid.UUID          `json:"created_by,omitempty"`
	UpdatedAt   time.Time           `json:"updated_at"`
	UpdatedBy   *uuid.UUID          `json:"updated_by,omitempty"`
	Rows        []*PayrollImportRow `json:"rows,omitempty"`
}

type PayrollImportRow struct {
	ID                 uuid.UUID       `json:"id"`
	TenantID           uuid.UUID       `json:"tenant_id"`
	BatchID            uuid.UUID       `json:"batch_id"`
	RowNumber          int32           `json:"row_number"`
	EmployeeCode       *string         `json:"employee_code,omitempty"`
	EmployeeUserID     *uuid.UUID      `json:"employee_user_id,omitempty"`
	EmployeeName       *string         `json:"employee_name,omitempty"`
	GrossSalary        *float64        `json:"gross_salary,omitempty"`
	PresentDays        *float64        `json:"present_days,omitempty"`
	AbsentDays         *float64        `json:"absent_days,omitempty"`
	LOPDays            *float64        `json:"lop_days,omitempty"`
	VariableEarnings   *float64        `json:"variable_earnings,omitempty"`
	VariableDeductions *float64        `json:"variable_deductions,omitempty"`
	Status             string          `json:"status"`
	ErrorMessage       *string         `json:"error_message,omitempty"`
	RawData            json.RawMessage `json:"raw_data,omitempty"`
	Inactive           bool            `json:"inactive"`
	CreatedAt          time.Time       `json:"created_at"`
	CreatedBy          *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt          time.Time       `json:"updated_at"`
	UpdatedBy          *uuid.UUID      `json:"updated_by,omitempty"`
}

type ConsolidatedSalarySheetRow struct {
	SalarySlipID    uuid.UUID `json:"salary_slip_id"`
	TenantID        uuid.UUID `json:"tenant_id"`
	UserID          uuid.UUID `json:"user_id"`
	EmployeeCode    *string   `json:"employee_code,omitempty"`
	Firstname       string    `json:"firstname"`
	Lastname        *string   `json:"lastname,omitempty"`
	Email           *string   `json:"email,omitempty"`
	BranchName      *string   `json:"branch_name,omitempty"`
	DepartmentName  *string   `json:"department_name,omitempty"`
	Month           int32     `json:"month"`
	Year            int32     `json:"year"`
	GrossSalary     float64   `json:"gross_salary"`
	TotalEarnings   float64   `json:"total_earnings"`
	TotalDeductions float64   `json:"total_deductions"`
	AbsentDeduction float64   `json:"absent_deduction"`
	NetSalary       float64   `json:"net_salary"`
	PresentDays     int32     `json:"present_days"`
	AbsentDays      int32     `json:"absent_days"`
	LWPDays         float64   `json:"lwp_days"`
	PDFPath         *string   `json:"pdf_path,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}

type PayrollReconciliationRow struct {
	EmployeeID           uuid.UUID  `json:"employee_id"`
	UserID               uuid.UUID  `json:"user_id"`
	EmployeeCode         *string    `json:"employee_code,omitempty"`
	Firstname            string     `json:"firstname"`
	Lastname             *string    `json:"lastname,omitempty"`
	Email                *string    `json:"email,omitempty"`
	BranchName           *string    `json:"branch_name,omitempty"`
	DepartmentName       *string    `json:"department_name,omitempty"`
	SalarySlipID         *uuid.UUID `json:"salary_slip_id,omitempty"`
	PresentDays          *int32     `json:"present_days,omitempty"`
	AbsentDays           *int32     `json:"absent_days,omitempty"`
	LWPDays              *float64   `json:"lwp_days,omitempty"`
	NetSalary            *float64   `json:"net_salary,omitempty"`
	ReconciliationStatus string     `json:"reconciliation_status"`
}

func ValidatePayrollPeriod(month int32, year int32) error {
	if month < 1 || month > 12 || year < 1900 || year > 9999 {
		return ErrInvalidPayrollPeriod
	}
	return nil
}

func ValidatePayrollLockStatus(value string) (string, error) {
	status := strings.TrimSpace(value)
	switch status {
	case PayrollLockStatusOpen, PayrollLockStatusLocked, PayrollLockStatusUnlocked:
		return status, nil
	default:
		return "", ErrInvalidPayrollPeriod
	}
}

func NewPayrollStatutoryRule(item PayrollStatutoryRule) (*PayrollStatutoryRule, error) {
	if item.TenantID == uuid.Nil || strings.TrimSpace(item.Name) == "" {
		return nil, ErrInvalidPayrollStatutoryRule
	}
	ruleType := strings.TrimSpace(strings.ToLower(item.RuleType))
	if ruleType != PayrollRulePT && ruleType != PayrollRuleLWF {
		return nil, ErrInvalidPayrollStatutoryRule
	}
	if item.EffectiveFrom.IsZero() || (item.EffectiveTo != nil && item.EffectiveTo.Before(item.EffectiveFrom)) {
		return nil, ErrInvalidPayrollStatutoryRule
	}
	if item.EmployeeAmount < 0 || item.EmployerAmount < 0 || math.IsNaN(item.EmployeeAmount) || math.IsNaN(item.EmployerAmount) {
		return nil, ErrInvalidPayrollStatutoryRule
	}
	if item.DeductionMonth != nil && (*item.DeductionMonth < 1 || *item.DeductionMonth > 12) {
		return nil, ErrInvalidPayrollStatutoryRule
	}
	item.RuleType = ruleType
	item.Name = strings.TrimSpace(item.Name)
	if strings.TrimSpace(item.Frequency) == "" {
		item.Frequency = "monthly"
	}
	now := time.Now().UTC()
	item.CreatedAt = now
	item.UpdatedAt = now
	return &item, nil
}
