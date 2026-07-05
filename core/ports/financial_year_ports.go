package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

type FinancialYearRepo interface {
	CreateFinancialYear(ctx context.Context, item *domain.FinancialYear, actorID *uuid.UUID) (*domain.FinancialYear, error)
	ListFinancialYears(ctx context.Context, tenantID uuid.UUID) ([]*domain.FinancialYear, error)
	GetFinancialYear(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.FinancialYear, error)
	GetActiveFinancialYear(ctx context.Context, tenantID uuid.UUID) (*domain.FinancialYear, error)
	UpdateFinancialYear(ctx context.Context, item *domain.FinancialYear, actorID *uuid.UUID) (*domain.FinancialYear, error)
	DeleteFinancialYear(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	SetActiveFinancialYear(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) (*domain.FinancialYear, error)
	SetFinancialYearLock(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, isLocked bool, closeNote *string, actorID *uuid.UUID) (*domain.FinancialYear, error)
}

type FinancialYearCommand struct {
	ID            uuid.UUID  `json:"id,omitempty"`
	TenantID      uuid.UUID  `json:"tenant_id"`
	Name          string     `json:"name"`
	StartDate     string     `json:"start_date"`
	EndDate       string     `json:"end_date"`
	IsActive      bool       `json:"is_active"`
	PayrollYear   bool       `json:"payroll_year"`
	LeaveYear     bool       `json:"leave_year"`
	HolidayYear   bool       `json:"holiday_year"`
	ReportingYear bool       `json:"reporting_year"`
	CloseNote     *string    `json:"close_note,omitempty"`
	ActorID       *uuid.UUID `json:"-"`
}

type FinancialYearLockCommand struct {
	ID        uuid.UUID  `json:"id,omitempty"`
	TenantID  uuid.UUID  `json:"tenant_id"`
	IsLocked  bool       `json:"is_locked"`
	CloseNote *string    `json:"close_note,omitempty"`
	ActorID   *uuid.UUID `json:"-"`
}
