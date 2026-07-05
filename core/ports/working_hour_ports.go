package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

type WorkingHourRepo interface {
	CreateWorkingHour(ctx context.Context, item *domain.WorkingHour, actorID *uuid.UUID) (*domain.WorkingHour, error)
	ListWorkingHours(ctx context.Context, tenantID uuid.UUID) ([]*domain.WorkingHour, error)
	GetWorkingHour(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.WorkingHour, error)
	ResolveWorkingHour(ctx context.Context, tenantID uuid.UUID, dayOfWeek string, branchID *uuid.UUID, userID *uuid.UUID) (*domain.WorkingHour, error)
	UpdateWorkingHour(ctx context.Context, item *domain.WorkingHour, actorID *uuid.UUID) (*domain.WorkingHour, error)
	DeleteWorkingHour(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	CopyTenantWorkingHoursToBranch(ctx context.Context, tenantID uuid.UUID, branchID uuid.UUID, actorID *uuid.UUID) ([]*domain.WorkingHour, error)
}

type WorkingHourCommand struct {
	ID           uuid.UUID  `json:"id,omitempty"`
	TenantID     uuid.UUID  `json:"tenant_id"`
	BranchID     *uuid.UUID `json:"branch_id,omitempty"`
	UserID       *uuid.UUID `json:"user_id,omitempty"`
	DayOfWeek    string     `json:"day_of_week"`
	IsWorkingDay bool       `json:"is_working_day"`
	StartTime    string     `json:"start_time"`
	EndTime      string     `json:"end_time"`
	BreakMinutes int32      `json:"break_minutes"`
	ActorID      *uuid.UUID `json:"-"`
}

type ResolveWorkingHourCommand struct {
	TenantID  uuid.UUID  `json:"tenant_id"`
	BranchID  *uuid.UUID `json:"branch_id,omitempty"`
	UserID    *uuid.UUID `json:"user_id,omitempty"`
	DayOfWeek string     `json:"day_of_week"`
}

type CopyWorkingHoursCommand struct {
	TenantID uuid.UUID  `json:"tenant_id"`
	BranchID uuid.UUID  `json:"branch_id"`
	ActorID  *uuid.UUID `json:"-"`
}
