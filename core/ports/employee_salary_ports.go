package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

type EmployeeSalaryRepo interface {
	CreateEmployeeSalary(ctx context.Context, item *domain.EmployeeSalary, actorID *uuid.UUID) (*domain.EmployeeSalary, error)
	UpdateEmployeeSalary(ctx context.Context, item *domain.EmployeeSalary, actorID *uuid.UUID) (*domain.EmployeeSalary, error)
	ListEmployeeSalariesByUser(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID) ([]*domain.EmployeeSalary, error)
	GetEmployeeSalary(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.EmployeeSalary, error)
	DeleteEmployeeSalary(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	DeleteEmployeeSalariesByUserFY(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID, fyID uuid.UUID, actorID *uuid.UUID) error
	CreateEmployeeSalaryStructure(ctx context.Context, item *domain.EmployeeSalaryStructure, actorID *uuid.UUID) (*domain.EmployeeSalaryStructure, error)
	ListEmployeeSalaryStructures(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID, fyID uuid.UUID) ([]*domain.EmployeeSalaryStructure, error)
	DeleteEmployeeSalaryStructuresByUserFY(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID, fyID uuid.UUID, actorID *uuid.UUID) error
}

type EmployeeSalaryCommand struct {
	ID            uuid.UUID  `json:"id,omitempty"`
	TenantID      uuid.UUID  `json:"tenant_id"`
	UserID        uuid.UUID  `json:"user_id"`
	FYID          uuid.UUID  `json:"fy_id"`
	TemplateID    uuid.UUID  `json:"template_id"`
	GrossSalary   float64    `json:"gross_salary"`
	EffectiveFrom *string    `json:"effective_from,omitempty"`
	ActorID       *uuid.UUID `json:"-"`
}

type EmployeeSalaryCalculationCommand struct {
	TenantID    uuid.UUID  `json:"tenant_id"`
	UserID      uuid.UUID  `json:"user_id"`
	FYID        uuid.UUID  `json:"fy_id"`
	SalaryID    *uuid.UUID `json:"salary_id,omitempty"`
	Month       int        `json:"month"`
	Year        int        `json:"year"`
	PresentDays *int       `json:"present_days,omitempty"`
	AbsentDays  *int       `json:"absent_days,omitempty"`
	TotalDays   *int       `json:"total_days,omitempty"`
	IsSpecial   bool       `json:"is_special"`
}

type EmployeeSalaryCalculation struct {
	TenantID     uuid.UUID           `json:"tenant_id"`
	UserID       uuid.UUID           `json:"user_id"`
	FYID         uuid.UUID           `json:"fy_id"`
	SalaryID     uuid.UUID           `json:"salary_id"`
	Month        int                 `json:"month"`
	Year         int                 `json:"year"`
	PresentDays  int                 `json:"present_days"`
	AbsentDays   int                 `json:"absent_days"`
	TotalDays    int                 `json:"total_days"`
	LWPDays      int                 `json:"lwp_days"`
	IsSpecial    bool                `json:"is_special"`
	GrossSalary  float64             `json:"gross_salary"`
	SalaryResult domain.SalaryResult `json:"salary_result"`
}
