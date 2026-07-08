package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	CompOffStatusPending  = "pending"
	CompOffStatusApproved = "approved"
	CompOffStatusRejected = "rejected"
	CompOffStatusCanceled = "canceled"
)

var (
	ErrInvalidCompOffRequest      = errors.New("comp-off request is invalid")
	ErrInvalidCompOffStatus       = errors.New("comp-off status is invalid")
	ErrCompOffRequestNotFound     = errors.New("comp-off request not found")
	ErrCompOffRequestNotPending   = errors.New("comp-off request is not pending")
	ErrCompOffLedgerAlreadyExists = errors.New("comp-off ledger credit already exists")
)

type CompOffRequest struct {
	ID                 uuid.UUID      `json:"id"`
	TenantID           uuid.UUID      `json:"tenant_id"`
	UserID             uuid.UUID      `json:"user_id"`
	LeaveTypeID        uuid.UUID      `json:"leave_type_id"`
	FYID               uuid.UUID      `json:"fy_id"`
	WorkDate           time.Time      `json:"work_date"`
	WorkedMinutes      int32          `json:"worked_minutes"`
	RequestedDays      float64        `json:"requested_days"`
	ApprovedDays       *float64       `json:"approved_days,omitempty"`
	ExpiryDate         *time.Time     `json:"expiry_date,omitempty"`
	Reason             *string        `json:"reason,omitempty"`
	Status             string         `json:"status"`
	ReviewedBy         *uuid.UUID     `json:"reviewed_by,omitempty"`
	ReviewedAt         *time.Time     `json:"reviewed_at,omitempty"`
	ReviewRemarks      *string        `json:"review_remarks,omitempty"`
	PayrollImpact      bool           `json:"payroll_impact"`
	SourceAttendanceID *uuid.UUID     `json:"source_attendance_id,omitempty"`
	SourceSegmentID    *uuid.UUID     `json:"source_segment_id,omitempty"`
	Metadata           map[string]any `json:"metadata"`
	Inactive           bool           `json:"inactive"`
	CreatedAt          time.Time      `json:"created_at"`
	CreatedBy          *uuid.UUID     `json:"created_by,omitempty"`
	UpdatedAt          time.Time      `json:"updated_at"`
	UpdatedBy          *uuid.UUID     `json:"updated_by,omitempty"`
}

func NormalizeCompOffStatus(value string) string {
	status := strings.ToLower(strings.TrimSpace(value))
	if status == "" {
		return CompOffStatusPending
	}
	switch status {
	case CompOffStatusPending, CompOffStatusApproved, CompOffStatusRejected, CompOffStatusCanceled:
		return status
	default:
		return ""
	}
}

func ValidateCompOffRequest(item *CompOffRequest) error {
	if item == nil {
		return ErrInvalidCompOffRequest
	}
	if item.TenantID == uuid.Nil {
		return ErrInvalidTenantID
	}
	if item.UserID == uuid.Nil {
		return ErrInvalidLeaveBalanceUser
	}
	if item.LeaveTypeID == uuid.Nil {
		return ErrInvalidLeavePolicyType
	}
	if item.FYID == uuid.Nil {
		return ErrInvalidLeavePolicyFY
	}
	if item.WorkDate.IsZero() || item.RequestedDays <= 0 || item.WorkedMinutes < 0 {
		return ErrInvalidCompOffRequest
	}
	item.Status = NormalizeCompOffStatus(item.Status)
	if item.Status == "" {
		return ErrInvalidCompOffStatus
	}
	if item.ApprovedDays != nil && *item.ApprovedDays < 0 {
		return ErrInvalidCompOffRequest
	}
	if item.Metadata == nil {
		item.Metadata = map[string]any{}
	}
	return nil
}
