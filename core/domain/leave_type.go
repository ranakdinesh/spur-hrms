package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidLeaveTypeID    = errors.New("leave_type_id is required")
	ErrInvalidLeaveTypeName  = errors.New("leave type name is required")
	ErrInvalidLeaveShortcode = errors.New("leave type shortcode is required")
	ErrInvalidLeaveTypeLimit = errors.New("leave type limits cannot be negative")
	ErrLeaveTypeNotFound     = errors.New("leave type not found")
	ErrLeaveTypeDisabled     = errors.New("leave type is not enabled for this tenant")
	ErrSystemLeaveTypeLocked = errors.New("system leave types cannot be edited or deleted")
)

type LeaveType struct {
	ID                   uuid.UUID  `json:"id"`
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
	IsSystem             bool       `json:"is_system"`
	Inactive             bool       `json:"inactive"`
	CreatedAt            time.Time  `json:"created_at"`
	CreatedBy            *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt            time.Time  `json:"updated_at"`
	UpdatedBy            *uuid.UUID `json:"updated_by,omitempty"`
}

type LeaveTypeInput struct {
	TenantID             uuid.UUID
	Name                 string
	Shortcode            *string
	Description          *string
	IsPaid               bool
	IsCarryForward       bool
	MaxCarryForward      int32
	IsConsecutiveLimit   bool
	ConsecutiveDaysLimit int32
	IsEnabled            bool
	IsSystem             bool
}

func NewLeaveType(input LeaveTypeInput) (*LeaveType, error) {
	if input.TenantID == uuid.Nil {
		return nil, ErrInvalidTenantID
	}
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, ErrInvalidLeaveTypeName
	}
	shortcode := cleanShortcode(input.Shortcode)
	if shortcode == nil {
		return nil, ErrInvalidLeaveShortcode
	}
	if input.MaxCarryForward < 0 || input.ConsecutiveDaysLimit < 0 {
		return nil, ErrInvalidLeaveTypeLimit
	}
	if !input.IsCarryForward {
		input.MaxCarryForward = 0
	}
	if !input.IsConsecutiveLimit {
		input.ConsecutiveDaysLimit = 0
	}
	now := time.Now().UTC()
	return &LeaveType{
		TenantID:             input.TenantID,
		Name:                 name,
		Shortcode:            shortcode,
		Description:          cleanString(input.Description),
		IsPaid:               input.IsPaid,
		IsCarryForward:       input.IsCarryForward,
		MaxCarryForward:      input.MaxCarryForward,
		IsConsecutiveLimit:   input.IsConsecutiveLimit,
		ConsecutiveDaysLimit: input.ConsecutiveDaysLimit,
		IsEnabled:            input.IsEnabled,
		IsSystem:             input.IsSystem,
		CreatedAt:            now,
		UpdatedAt:            now,
	}, nil
}

func cleanShortcode(value *string) *string {
	if value == nil {
		return nil
	}
	clean := strings.ToUpper(strings.TrimSpace(*value))
	if clean == "" {
		return nil
	}
	return &clean
}

func cleanString(value *string) *string {
	if value == nil {
		return nil
	}
	clean := strings.TrimSpace(*value)
	if clean == "" {
		return nil
	}
	return &clean
}
