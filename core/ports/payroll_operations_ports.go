package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

type PayrollOperationsRepo interface {
	UpsertPayrollPeriodLock(ctx context.Context, item *domain.PayrollPeriodLock, actorID *uuid.UUID) (*domain.PayrollPeriodLock, error)
	GetPayrollPeriodLock(ctx context.Context, tenantID uuid.UUID, month int32, year int32) (*domain.PayrollPeriodLock, error)
	ListPayrollPeriodLocks(ctx context.Context, tenantID uuid.UUID) ([]*domain.PayrollPeriodLock, error)
	CreatePayrollPeriodLockEvent(ctx context.Context, item *domain.PayrollPeriodLockEvent, actorID *uuid.UUID) (*domain.PayrollPeriodLockEvent, error)
	ListPayrollPeriodLockEvents(ctx context.Context, tenantID uuid.UUID, lockID uuid.UUID) ([]*domain.PayrollPeriodLockEvent, error)
	ListPayrollStatutoryRules(ctx context.Context, tenantID uuid.UUID, ruleType *string) ([]*domain.PayrollStatutoryRule, error)
	GetPayrollStatutoryRule(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.PayrollStatutoryRule, error)
	CreatePayrollStatutoryRule(ctx context.Context, item *domain.PayrollStatutoryRule, actorID *uuid.UUID) (*domain.PayrollStatutoryRule, error)
	UpdatePayrollStatutoryRule(ctx context.Context, item *domain.PayrollStatutoryRule, actorID *uuid.UUID) (*domain.PayrollStatutoryRule, error)
	DeletePayrollStatutoryRule(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	ResolvePayrollStatutoryRule(ctx context.Context, tenantID uuid.UUID, ruleType string, state *string, branchID *uuid.UUID, effectiveDate string, grossSalary float64, month int32) (*domain.PayrollStatutoryRule, error)
	CreatePayrollImportBatch(ctx context.Context, item *domain.PayrollImportBatch, actorID *uuid.UUID) (*domain.PayrollImportBatch, error)
	ListPayrollImportBatches(ctx context.Context, tenantID uuid.UUID, limit int32, offset int32) ([]*domain.PayrollImportBatch, error)
	GetPayrollImportBatch(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.PayrollImportBatch, error)
	CreatePayrollImportRow(ctx context.Context, item *domain.PayrollImportRow, actorID *uuid.UUID) (*domain.PayrollImportRow, error)
	ListPayrollImportRows(ctx context.Context, tenantID uuid.UUID, batchID uuid.UUID) ([]*domain.PayrollImportRow, error)
	ListConsolidatedSalarySheet(ctx context.Context, tenantID uuid.UUID, month int32, year int32) ([]*domain.ConsolidatedSalarySheetRow, error)
	ListPayrollReconciliationRows(ctx context.Context, tenantID uuid.UUID, month int32, year int32) ([]*domain.PayrollReconciliationRow, error)
}

type PayrollPeriodLockCommand struct {
	TenantID     uuid.UUID  `json:"tenant_id"`
	Month        int32      `json:"month"`
	Year         int32      `json:"year"`
	Status       string     `json:"status"`
	UnlockReason *string    `json:"unlock_reason,omitempty"`
	Notes        *string    `json:"notes,omitempty"`
	ActorID      *uuid.UUID `json:"-"`
}

type PayrollStatutoryRuleCommand struct {
	ID             uuid.UUID  `json:"id,omitempty"`
	TenantID       uuid.UUID  `json:"tenant_id"`
	RuleType       string     `json:"rule_type"`
	Name           string     `json:"name"`
	State          *string    `json:"state,omitempty"`
	BranchID       *uuid.UUID `json:"branch_id,omitempty"`
	EffectiveFrom  string     `json:"effective_from"`
	EffectiveTo    string     `json:"effective_to,omitempty"`
	MinGrossSalary *float64   `json:"min_gross_salary,omitempty"`
	MaxGrossSalary *float64   `json:"max_gross_salary,omitempty"`
	EmployeeAmount float64    `json:"employee_amount"`
	EmployerAmount float64    `json:"employer_amount"`
	Frequency      string     `json:"frequency"`
	DeductionMonth *int32     `json:"deduction_month,omitempty"`
	Notes          *string    `json:"notes,omitempty"`
	ActorID        *uuid.UUID `json:"-"`
}

type PayrollImportRowCommand struct {
	EmployeeCode       string         `json:"employee_code"`
	GrossSalary        *float64       `json:"gross_salary,omitempty"`
	PresentDays        *float64       `json:"present_days,omitempty"`
	AbsentDays         *float64       `json:"absent_days,omitempty"`
	LOPDays            *float64       `json:"lop_days,omitempty"`
	VariableEarnings   *float64       `json:"variable_earnings,omitempty"`
	VariableDeductions *float64       `json:"variable_deductions,omitempty"`
	RawData            map[string]any `json:"raw_data,omitempty"`
}

type PayrollImportCommand struct {
	TenantID   uuid.UUID                 `json:"tenant_id"`
	ImportType string                    `json:"import_type"`
	Month      *int32                    `json:"month,omitempty"`
	Year       *int32                    `json:"year,omitempty"`
	FYID       *uuid.UUID                `json:"fy_id,omitempty"`
	TemplateID *uuid.UUID                `json:"template_id,omitempty"`
	FileName   *string                   `json:"file_name,omitempty"`
	Notes      *string                   `json:"notes,omitempty"`
	Rows       []PayrollImportRowCommand `json:"rows"`
	Apply      bool                      `json:"apply"`
	ActorID    *uuid.UUID                `json:"-"`
}
