package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

type LeaveTypeRepo interface {
	CreateLeaveType(ctx context.Context, item *domain.LeaveType, actorID *uuid.UUID) (*domain.LeaveType, error)
	UpsertSystemLeaveType(ctx context.Context, item *domain.LeaveType, actorID *uuid.UUID) (*domain.LeaveType, error)
	ListLeaveTypes(ctx context.Context, tenantID uuid.UUID) ([]*domain.LeaveType, error)
	GetLeaveType(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.LeaveType, error)
	GetLeaveTypeByShortcode(ctx context.Context, tenantID uuid.UUID, shortcode string) (*domain.LeaveType, error)
	UpdateLeaveType(ctx context.Context, item *domain.LeaveType, actorID *uuid.UUID) (*domain.LeaveType, error)
	DeleteLeaveType(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
}

type LeaveTypeCommand struct {
	ID                   uuid.UUID  `json:"id,omitempty"`
	TenantID             uuid.UUID  `json:"tenant_id"`
	Name                 string     `json:"name"`
	Shortcode            *string    `json:"shortcode,omitempty"`
	Description          *string    `json:"description,omitempty"`
	IsPaid               bool       `json:"is_paid"`
	IsCarryForward       bool       `json:"is_carry_forward"`
	MaxCarryForward      int32      `json:"max_carry_forward"`
	IsConsecutiveLimit   bool       `json:"is_consecutive_limit"`
	ConsecutiveDaysLimit int32      `json:"consecutive_days_limit"`
	IsEnabled            bool       `json:"is_enabled"`
	ActorID              *uuid.UUID `json:"-"`
}
