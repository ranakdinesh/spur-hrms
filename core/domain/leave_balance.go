package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

const (
	LeaveLedgerDebit  = "debit"
	LeaveLedgerCredit = "credit"

	LeaveLedgerSourceOpeningBalance   = "opening_balance"
	LeaveLedgerSourceMonthlyAccrual   = "monthly_accrual"
	LeaveLedgerSourceYearlyAccrual    = "yearly_accrual"
	LeaveLedgerSourceLeaveApply       = "leave_apply"
	LeaveLedgerSourceLeaveApprove     = "leave_approve"
	LeaveLedgerSourceLeaveReject      = "leave_reject"
	LeaveLedgerSourceLeaveCancel      = "leave_cancel"
	LeaveLedgerSourceCompOff          = "comp_off"
	LeaveLedgerSourceCarryForward     = "carry_forward"
	LeaveLedgerSourceEncashment       = "encashment"
	LeaveLedgerSourceManualAdjustment = "manual_adjustment"
)

var (
	ErrInvalidLeaveBalanceID    = errors.New("leave balance id is required")
	ErrInvalidLeaveBalanceUser  = errors.New("leave balance user_id is required")
	ErrInvalidLeaveBalanceDays  = errors.New("leave balance days cannot be negative")
	ErrInvalidLeaveLedgerType   = errors.New("leave ledger transaction type must be credit or debit")
	ErrInvalidLeaveLedgerSource = errors.New("leave ledger source type is invalid")
	ErrLeaveBalanceNotFound     = errors.New("leave balance not found")
	ErrLeaveLedgerEntryNotFound = errors.New("leave ledger entry not found")
)

type LeaveBalance struct {
	ID          uuid.UUID  `json:"id"`
	TenantID    uuid.UUID  `json:"tenant_id"`
	UserID      uuid.UUID  `json:"user_id"`
	LeaveTypeID uuid.UUID  `json:"leave_type_id"`
	FYID        uuid.UUID  `json:"fy_id"`
	TotalDays   float64    `json:"total_days"`
	UsedDays    float64    `json:"used_days"`
	PendingDays float64    `json:"pending_days"`
	BalanceDays float64    `json:"balance_days"`
	Inactive    bool       `json:"inactive"`
	CreatedAt   time.Time  `json:"created_at"`
	CreatedBy   *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt   time.Time  `json:"updated_at"`
	UpdatedBy   *uuid.UUID `json:"updated_by,omitempty"`
}

type LeaveLedgerEntry struct {
	ID              uuid.UUID      `json:"id"`
	TenantID        uuid.UUID      `json:"tenant_id"`
	UserID          uuid.UUID      `json:"user_id"`
	LeaveTypeID     uuid.UUID      `json:"leave_type_id"`
	FYID            uuid.UUID      `json:"fy_id"`
	LeaveID         *uuid.UUID     `json:"leave_id,omitempty"`
	TransactionType string         `json:"transaction_type"`
	Days            float64        `json:"days"`
	Remarks         *string        `json:"remarks,omitempty"`
	SourceType      string         `json:"source_type"`
	SourceID        *uuid.UUID     `json:"source_id,omitempty"`
	BalanceBefore   *float64       `json:"balance_before,omitempty"`
	BalanceAfter    *float64       `json:"balance_after,omitempty"`
	PendingBefore   *float64       `json:"pending_before,omitempty"`
	PendingAfter    *float64       `json:"pending_after,omitempty"`
	UsedBefore      *float64       `json:"used_before,omitempty"`
	UsedAfter       *float64       `json:"used_after,omitempty"`
	Metadata        map[string]any `json:"metadata"`
	Inactive        bool           `json:"inactive"`
	CreatedAt       time.Time      `json:"created_at"`
	CreatedBy       *uuid.UUID     `json:"created_by,omitempty"`
}

func ValidateLeaveBalance(balance *LeaveBalance) error {
	if balance == nil {
		return ErrInvalidLeaveBalanceID
	}
	if balance.TenantID == uuid.Nil {
		return ErrInvalidTenantID
	}
	if balance.UserID == uuid.Nil {
		return ErrInvalidLeaveBalanceUser
	}
	if balance.LeaveTypeID == uuid.Nil {
		return ErrInvalidLeavePolicyType
	}
	if balance.FYID == uuid.Nil {
		return ErrInvalidLeavePolicyFY
	}
	if balance.TotalDays < 0 || balance.UsedDays < 0 || balance.PendingDays < 0 {
		return ErrInvalidLeaveBalanceDays
	}
	return nil
}

func ValidateLeaveLedgerEntry(entry *LeaveLedgerEntry) error {
	if entry == nil {
		return ErrLeaveLedgerEntryNotFound
	}
	if entry.TenantID == uuid.Nil {
		return ErrInvalidTenantID
	}
	if entry.UserID == uuid.Nil {
		return ErrInvalidLeaveBalanceUser
	}
	if entry.LeaveTypeID == uuid.Nil {
		return ErrInvalidLeavePolicyType
	}
	if entry.FYID == uuid.Nil {
		return ErrInvalidLeavePolicyFY
	}
	if entry.TransactionType != LeaveLedgerDebit && entry.TransactionType != LeaveLedgerCredit {
		return ErrInvalidLeaveLedgerType
	}
	if !validLeaveLedgerSource(entry.SourceType) {
		return ErrInvalidLeaveLedgerSource
	}
	if entry.Days < 0 {
		return ErrInvalidLeaveBalanceDays
	}
	if entry.Metadata == nil {
		entry.Metadata = map[string]any{}
	}
	return nil
}

func validLeaveLedgerSource(value string) bool {
	switch value {
	case LeaveLedgerSourceOpeningBalance, LeaveLedgerSourceMonthlyAccrual, LeaveLedgerSourceYearlyAccrual, LeaveLedgerSourceLeaveApply, LeaveLedgerSourceLeaveApprove, LeaveLedgerSourceLeaveReject, LeaveLedgerSourceLeaveCancel, LeaveLedgerSourceCompOff, LeaveLedgerSourceCarryForward, LeaveLedgerSourceEncashment, LeaveLedgerSourceManualAdjustment:
		return true
	default:
		return false
	}
}
