package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

type LeaveRequestRepo interface {
	CreateLeave(ctx context.Context, item *domain.Leave, actorID *uuid.UUID) (*domain.Leave, error)
	CreateLeaveApproval(ctx context.Context, item *domain.LeaveApproval, actorID *uuid.UUID) (*domain.LeaveApproval, error)
	ListLeavesByUser(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID) ([]*domain.Leave, error)
	ListLeavesByFY(ctx context.Context, tenantID uuid.UUID, fyID uuid.UUID) ([]*domain.Leave, error)
	ListLeaveReportRows(ctx context.Context, filter domain.LeaveReportFilter) ([]*domain.LeaveReportRow, error)
	GetLeaveReportSummary(ctx context.Context, filter domain.LeaveReportFilter) (*domain.LeaveReportSummary, error)
	ListOverlappingLeaves(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID, startDate string, endDate string) ([]*domain.Leave, error)
	GetLeave(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.Leave, error)
}

type ApplyLeaveCommand struct {
	TenantID     uuid.UUID  `json:"tenant_id"`
	UserID       uuid.UUID  `json:"user_id"`
	LeaveTypeID  uuid.UUID  `json:"leave_type_id"`
	FYID         uuid.UUID  `json:"fy_id"`
	StartDate    string     `json:"start_date"`
	EndDate      string     `json:"end_date"`
	StartDayType string     `json:"start_day_type"`
	EndDayType   string     `json:"end_day_type"`
	Reason       *string    `json:"reason,omitempty"`
	ApproverID   *uuid.UUID `json:"approver_id,omitempty"`
	ActorID      *uuid.UUID `json:"-"`
}

type LeaveNotifier interface {
	NotifyLeaveApplied(ctx context.Context, application *domain.LeaveApplication) error
}
