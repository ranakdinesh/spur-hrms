package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

type LeaveBalanceRepo interface {
	UpsertLeaveBalance(ctx context.Context, item *domain.LeaveBalance, actorID *uuid.UUID) (*domain.LeaveBalance, error)
	GetLeaveBalance(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID, leaveTypeID uuid.UUID, fyID uuid.UUID) (*domain.LeaveBalance, error)
	ListLeaveBalancesByUser(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID) ([]*domain.LeaveBalance, error)
	ListLeaveBalancesByTenantFY(ctx context.Context, tenantID uuid.UUID, fyID uuid.UUID) ([]*domain.LeaveBalance, error)
	AddLeaveBalanceCredit(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID, leaveTypeID uuid.UUID, fyID uuid.UUID, days float64, actorID *uuid.UUID) (*domain.LeaveBalance, error)
	UpdateLeaveBalancePending(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID, leaveTypeID uuid.UUID, fyID uuid.UUID, daysDelta float64, actorID *uuid.UUID) (*domain.LeaveBalance, error)
	ReverseLeaveBalancePending(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID, leaveTypeID uuid.UUID, fyID uuid.UUID, days float64, actorID *uuid.UUID) (*domain.LeaveBalance, error)
	MoveLeaveBalancePendingToUsed(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID, leaveTypeID uuid.UUID, fyID uuid.UUID, days float64, actorID *uuid.UUID) (*domain.LeaveBalance, error)
	CreateLeaveLedgerEntry(ctx context.Context, item *domain.LeaveLedgerEntry, actorID *uuid.UUID) (*domain.LeaveLedgerEntry, error)
	ListLeaveLedgerByUser(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID) ([]*domain.LeaveLedgerEntry, error)
	ListLeaveLedgerByLeave(ctx context.Context, tenantID uuid.UUID, leaveID uuid.UUID) ([]*domain.LeaveLedgerEntry, error)
	GetLeaveLedgerBySource(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID, leaveTypeID uuid.UUID, fyID uuid.UUID, sourceType string, sourceID uuid.UUID) (*domain.LeaveLedgerEntry, error)
}

type LeaveBalanceCommand struct {
	TenantID    uuid.UUID  `json:"tenant_id"`
	UserID      uuid.UUID  `json:"user_id"`
	LeaveTypeID uuid.UUID  `json:"leave_type_id"`
	FYID        uuid.UUID  `json:"fy_id"`
	TotalDays   float64    `json:"total_days"`
	UsedDays    float64    `json:"used_days"`
	PendingDays float64    `json:"pending_days"`
	ActorID     *uuid.UUID `json:"-"`
}

type LeaveBalanceAdjustmentCommand struct {
	TenantID        uuid.UUID      `json:"tenant_id"`
	UserID          uuid.UUID      `json:"user_id"`
	LeaveTypeID     uuid.UUID      `json:"leave_type_id"`
	FYID            uuid.UUID      `json:"fy_id"`
	Days            float64        `json:"days"`
	TransactionType string         `json:"transaction_type"`
	SourceType      string         `json:"source_type"`
	SourceID        *uuid.UUID     `json:"source_id,omitempty"`
	LeaveID         *uuid.UUID     `json:"leave_id,omitempty"`
	Remarks         *string        `json:"remarks,omitempty"`
	Metadata        map[string]any `json:"metadata"`
	ActorID         *uuid.UUID     `json:"-"`
}

type RunLeaveAccrualCommand struct {
	TenantID uuid.UUID  `json:"tenant_id"`
	FYID     uuid.UUID  `json:"fy_id"`
	Month    int32      `json:"month"`
	ActorID  *uuid.UUID `json:"-"`
}
