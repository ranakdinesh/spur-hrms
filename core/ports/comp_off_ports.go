package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

type CompOffRequestRepo interface {
	CreateCompOffRequest(ctx context.Context, item *domain.CompOffRequest, actorID *uuid.UUID) (*domain.CompOffRequest, error)
	GetCompOffRequest(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.CompOffRequest, error)
	ListCompOffRequestsByUser(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID) ([]*domain.CompOffRequest, error)
	ListCompOffRequestsByStatus(ctx context.Context, tenantID uuid.UUID, status string) ([]*domain.CompOffRequest, error)
	ReviewCompOffRequest(ctx context.Context, item *domain.CompOffRequest, actorID *uuid.UUID) (*domain.CompOffRequest, error)
}

type CompOffRequestCommand struct {
	TenantID           uuid.UUID      `json:"tenant_id"`
	UserID             uuid.UUID      `json:"user_id"`
	LeaveTypeID        uuid.UUID      `json:"leave_type_id,omitempty"`
	FYID               uuid.UUID      `json:"fy_id,omitempty"`
	WorkDate           string         `json:"work_date"`
	WorkedMinutes      int32          `json:"worked_minutes"`
	RequestedDays      float64        `json:"requested_days"`
	ExpiryDate         *string        `json:"expiry_date,omitempty"`
	Reason             *string        `json:"reason,omitempty"`
	PayrollImpact      bool           `json:"payroll_impact"`
	SourceAttendanceID *uuid.UUID     `json:"source_attendance_id,omitempty"`
	SourceSegmentID    *uuid.UUID     `json:"source_segment_id,omitempty"`
	Metadata           map[string]any `json:"metadata"`
	ActorID            *uuid.UUID     `json:"-"`
}

type CompOffReviewCommand struct {
	TenantID      uuid.UUID      `json:"tenant_id"`
	RequestID     uuid.UUID      `json:"request_id"`
	Status        string         `json:"status"`
	ApprovedDays  *float64       `json:"approved_days,omitempty"`
	ExpiryDate    *string        `json:"expiry_date,omitempty"`
	PayrollImpact bool           `json:"payroll_impact"`
	Remarks       *string        `json:"remarks,omitempty"`
	Metadata      map[string]any `json:"metadata"`
	ActorID       *uuid.UUID     `json:"-"`
}
