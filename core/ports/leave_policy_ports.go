package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

type LeavePolicyRepo interface {
	CreateLeavePolicy(ctx context.Context, item *domain.LeavePolicy, actorID *uuid.UUID) (*domain.LeavePolicy, error)
	ListLeavePolicies(ctx context.Context, tenantID uuid.UUID) ([]*domain.LeavePolicy, error)
	GetLeavePolicy(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.LeavePolicy, error)
	GetLeavePolicyByTypeAndFY(ctx context.Context, tenantID uuid.UUID, leaveTypeID uuid.UUID, fyID uuid.UUID) (*domain.LeavePolicy, error)
	UpdateLeavePolicy(ctx context.Context, item *domain.LeavePolicy, actorID *uuid.UUID) (*domain.LeavePolicy, error)
	DeleteLeavePolicy(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
}

type LeavePolicyCommand struct {
	ID                   uuid.UUID  `json:"id,omitempty"`
	TenantID             uuid.UUID  `json:"tenant_id"`
	LeaveTypeID          uuid.UUID  `json:"leave_type_id"`
	FYID                 uuid.UUID  `json:"fy_id"`
	TotalDays            float64    `json:"total_days"`
	AllocationType       string     `json:"allocation_type"`
	Jan                  int32      `json:"jan"`
	Feb                  int32      `json:"feb"`
	Mar                  int32      `json:"mar"`
	Apr                  int32      `json:"apr"`
	May                  int32      `json:"may"`
	Jun                  int32      `json:"jun"`
	Jul                  int32      `json:"jul"`
	Aug                  int32      `json:"aug"`
	Sep                  int32      `json:"sep"`
	Oct                  int32      `json:"oct"`
	Nov                  int32      `json:"nov"`
	Dec                  int32      `json:"dec"`
	IsSandwichApplicable bool       `json:"is_sandwich_applicable"`
	ActorID              *uuid.UUID `json:"-"`
}

type MonthlyLeaveAllocationCommand struct {
	TenantID uuid.UUID `json:"tenant_id"`
	FYID     uuid.UUID `json:"fy_id"`
	Month    int32     `json:"month"`
}
