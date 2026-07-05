package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidHolidayID   = errors.New("holiday_id is required")
	ErrInvalidHolidayName = errors.New("holiday name is required")
	ErrInvalidHolidayDate = errors.New("holiday date is required")
	ErrHolidayNotFound    = errors.New("holiday not found")
)

type Holiday struct {
	ID         uuid.UUID  `json:"id"`
	TenantID   uuid.UUID  `json:"tenant_id"`
	BranchID   *uuid.UUID `json:"branch_id,omitempty"`
	FYID       *uuid.UUID `json:"fy_id,omitempty"`
	Name       string     `json:"name"`
	Date       time.Time  `json:"date"`
	IsOptional bool       `json:"is_optional"`
	Inactive   bool       `json:"inactive"`
	CreatedAt  time.Time  `json:"created_at"`
	CreatedBy  *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt  time.Time  `json:"updated_at"`
	UpdatedBy  *uuid.UUID `json:"updated_by,omitempty"`
}

type HolidayInput struct {
	TenantID   uuid.UUID
	BranchID   *uuid.UUID
	FYID       *uuid.UUID
	Name       string
	Date       time.Time
	IsOptional bool
}

func NewHoliday(input HolidayInput) (*Holiday, error) {
	if input.TenantID == uuid.Nil {
		return nil, ErrInvalidTenantID
	}
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, ErrInvalidHolidayName
	}
	if input.Date.IsZero() {
		return nil, ErrInvalidHolidayDate
	}
	now := time.Now().UTC()
	return &Holiday{
		TenantID:   input.TenantID,
		BranchID:   cleanUUIDOptional(input.BranchID),
		FYID:       cleanUUIDOptional(input.FYID),
		Name:       name,
		Date:       dateOnly(input.Date),
		IsOptional: input.IsOptional,
		CreatedAt:  now,
		UpdatedAt:  now,
	}, nil
}
