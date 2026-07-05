package ports

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

type HolidayRepo interface {
	CreateHoliday(ctx context.Context, item *domain.Holiday, actorID *uuid.UUID) (*domain.Holiday, error)
	ListHolidays(ctx context.Context, tenantID uuid.UUID) ([]*domain.Holiday, error)
	ListHolidaysByDateRange(ctx context.Context, tenantID uuid.UUID, startDate time.Time, endDate time.Time) ([]*domain.Holiday, error)
	ListHolidaysByFinancialYear(ctx context.Context, tenantID uuid.UUID, fyID uuid.UUID) ([]*domain.Holiday, error)
	ListHolidaysByBranch(ctx context.Context, tenantID uuid.UUID, branchID uuid.UUID) ([]*domain.Holiday, error)
	ListUpcomingHolidays(ctx context.Context, tenantID uuid.UUID, limit int32) ([]*domain.Holiday, error)
	GetHoliday(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.Holiday, error)
	UpdateHoliday(ctx context.Context, item *domain.Holiday, actorID *uuid.UUID) (*domain.Holiday, error)
	DeleteHoliday(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
}

type HolidayCommand struct {
	ID         uuid.UUID  `json:"id,omitempty"`
	TenantID   uuid.UUID  `json:"tenant_id"`
	BranchID   *uuid.UUID `json:"branch_id,omitempty"`
	FYID       *uuid.UUID `json:"fy_id,omitempty"`
	Name       string     `json:"name"`
	Date       string     `json:"date"`
	IsOptional bool       `json:"is_optional"`
	ActorID    *uuid.UUID `json:"-"`
}

type ListHolidaysCommand struct {
	TenantID  uuid.UUID  `json:"tenant_id"`
	BranchID  *uuid.UUID `json:"branch_id,omitempty"`
	FYID      *uuid.UUID `json:"fy_id,omitempty"`
	StartDate string     `json:"start_date,omitempty"`
	EndDate   string     `json:"end_date,omitempty"`
	Upcoming  bool       `json:"upcoming,omitempty"`
	Limit     int32      `json:"limit,omitempty"`
}
