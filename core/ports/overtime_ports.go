package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

type OvertimeRequestRepo interface {
	CreateOvertimeRequest(ctx context.Context, item *domain.OvertimeRequest, actorID *uuid.UUID) (*domain.OvertimeRequest, error)
	GetOvertimeRequest(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.OvertimeRequest, error)
	ListOvertimeRequestsByUser(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID) ([]*domain.OvertimeRequest, error)
	ListOvertimeRequestsByStatus(ctx context.Context, tenantID uuid.UUID, status string) ([]*domain.OvertimeRequest, error)
	ListOvertimeRequestsByPayrollExportStatus(ctx context.Context, tenantID uuid.UUID, status string) ([]*domain.OvertimeRequest, error)
	ReviewOvertimeRequest(ctx context.Context, item *domain.OvertimeRequest, actorID *uuid.UUID) (*domain.OvertimeRequest, error)
}

type OvertimeRequestCommand struct {
	TenantID             uuid.UUID      `json:"tenant_id"`
	UserID               uuid.UUID      `json:"user_id"`
	WorkDate             string         `json:"work_date"`
	RequestedMinutes     int32          `json:"requested_minutes"`
	Reason               *string        `json:"reason,omitempty"`
	CalculationType      string         `json:"calculation_type"`
	RateMultiplier       float64        `json:"rate_multiplier"`
	PayrollComponentCode *string        `json:"payroll_component_code,omitempty"`
	SourceAttendanceID   *uuid.UUID     `json:"source_attendance_id,omitempty"`
	SourceSegmentID      *uuid.UUID     `json:"source_segment_id,omitempty"`
	Metadata             map[string]any `json:"metadata"`
	ActorID              *uuid.UUID     `json:"-"`
}

type OvertimeReviewCommand struct {
	TenantID             uuid.UUID      `json:"tenant_id"`
	RequestID            uuid.UUID      `json:"request_id"`
	Status               string         `json:"status"`
	ApprovedMinutes      *int32         `json:"approved_minutes,omitempty"`
	Remarks              *string        `json:"remarks,omitempty"`
	CalculationType      string         `json:"calculation_type"`
	RateMultiplier       float64        `json:"rate_multiplier"`
	PayrollComponentCode *string        `json:"payroll_component_code,omitempty"`
	PayrollExportStatus  string         `json:"payroll_export_status"`
	Metadata             map[string]any `json:"metadata"`
	ActorID              *uuid.UUID     `json:"-"`
}
