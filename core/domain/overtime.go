package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	OvertimeStatusPending  = "pending"
	OvertimeStatusApproved = "approved"
	OvertimeStatusRejected = "rejected"
	OvertimeStatusCanceled = "canceled"

	OvertimeCalculationFixedRate        = "fixed_rate"
	OvertimeCalculationMultiplier       = "multiplier"
	OvertimeCalculationPayrollComponent = "payroll_component"
	OvertimeCalculationManual           = "manual"

	OvertimePayrollExportNotReady      = "not_ready"
	OvertimePayrollExportReady         = "ready"
	OvertimePayrollExportExported      = "exported"
	OvertimePayrollExportNotApplicable = "not_applicable"
)

var (
	ErrInvalidOvertimeRequest    = errors.New("overtime request is invalid")
	ErrInvalidOvertimeStatus     = errors.New("overtime status is invalid")
	ErrOvertimeRequestNotFound   = errors.New("overtime request not found")
	ErrOvertimeRequestNotPending = errors.New("overtime request is not pending")
)

type OvertimeRequest struct {
	ID                   uuid.UUID      `json:"id"`
	TenantID             uuid.UUID      `json:"tenant_id"`
	UserID               uuid.UUID      `json:"user_id"`
	WorkDate             time.Time      `json:"work_date"`
	RequestedMinutes     int32          `json:"requested_minutes"`
	ApprovedMinutes      *int32         `json:"approved_minutes,omitempty"`
	Reason               *string        `json:"reason,omitempty"`
	Status               string         `json:"status"`
	ReviewedBy           *uuid.UUID     `json:"reviewed_by,omitempty"`
	ReviewedAt           *time.Time     `json:"reviewed_at,omitempty"`
	ReviewRemarks        *string        `json:"review_remarks,omitempty"`
	CalculationType      string         `json:"calculation_type"`
	RateMultiplier       float64        `json:"rate_multiplier"`
	PayrollComponentCode *string        `json:"payroll_component_code,omitempty"`
	PayrollExportStatus  string         `json:"payroll_export_status"`
	PayrollExportedAt    *time.Time     `json:"payroll_exported_at,omitempty"`
	PayrollExportedBy    *uuid.UUID     `json:"payroll_exported_by,omitempty"`
	SourceAttendanceID   *uuid.UUID     `json:"source_attendance_id,omitempty"`
	SourceSegmentID      *uuid.UUID     `json:"source_segment_id,omitempty"`
	Metadata             map[string]any `json:"metadata"`
	Inactive             bool           `json:"inactive"`
	CreatedAt            time.Time      `json:"created_at"`
	CreatedBy            *uuid.UUID     `json:"created_by,omitempty"`
	UpdatedAt            time.Time      `json:"updated_at"`
	UpdatedBy            *uuid.UUID     `json:"updated_by,omitempty"`
}

func NormalizeOvertimeStatus(value string) string {
	status := strings.ToLower(strings.TrimSpace(value))
	if status == "" {
		return OvertimeStatusPending
	}
	switch status {
	case OvertimeStatusPending, OvertimeStatusApproved, OvertimeStatusRejected, OvertimeStatusCanceled:
		return status
	default:
		return ""
	}
}

func NormalizeOvertimeCalculationType(value string) string {
	kind := strings.ToLower(strings.TrimSpace(value))
	if kind == "" {
		return OvertimeCalculationMultiplier
	}
	switch kind {
	case OvertimeCalculationFixedRate, OvertimeCalculationMultiplier, OvertimeCalculationPayrollComponent, OvertimeCalculationManual:
		return kind
	default:
		return ""
	}
}

func NormalizeOvertimePayrollExportStatus(value string) string {
	status := strings.ToLower(strings.TrimSpace(value))
	if status == "" {
		return OvertimePayrollExportNotReady
	}
	switch status {
	case OvertimePayrollExportNotReady, OvertimePayrollExportReady, OvertimePayrollExportExported, OvertimePayrollExportNotApplicable:
		return status
	default:
		return ""
	}
}

func ValidateOvertimeRequest(item *OvertimeRequest) error {
	if item == nil {
		return ErrInvalidOvertimeRequest
	}
	if item.TenantID == uuid.Nil {
		return ErrInvalidTenantID
	}
	if item.UserID == uuid.Nil {
		return ErrInvalidEmployeeUserID
	}
	if item.WorkDate.IsZero() || item.RequestedMinutes <= 0 {
		return ErrInvalidOvertimeRequest
	}
	item.Status = NormalizeOvertimeStatus(item.Status)
	if item.Status == "" {
		return ErrInvalidOvertimeStatus
	}
	item.CalculationType = NormalizeOvertimeCalculationType(item.CalculationType)
	if item.CalculationType == "" {
		return ErrInvalidOvertimeRequest
	}
	item.PayrollExportStatus = NormalizeOvertimePayrollExportStatus(item.PayrollExportStatus)
	if item.PayrollExportStatus == "" {
		return ErrInvalidOvertimeRequest
	}
	if item.RateMultiplier < 0 {
		return ErrInvalidOvertimeRequest
	}
	if item.ApprovedMinutes != nil && *item.ApprovedMinutes < 0 {
		return ErrInvalidOvertimeRequest
	}
	if item.Metadata == nil {
		item.Metadata = map[string]any{}
	}
	return nil
}
