package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	SubscriptionStatusTrialing  = "trialing"
	SubscriptionStatusActive    = "active"
	SubscriptionStatusPastDue   = "past_due"
	SubscriptionStatusCancelled = "cancelled"
	SubscriptionStatusExpired   = "expired"
)

var (
	ErrInvalidSubscriptionID             = errors.New("subscription_id is required")
	ErrInvalidSubscriptionStatus         = errors.New("subscription status is invalid")
	ErrInvalidSubscriptionPeriod         = errors.New("subscription end date must be on or after start date")
	ErrInvalidSubscriptionEmployee       = errors.New("subscription employee limit cannot be negative")
	ErrTenantSubscriptionRequired        = errors.New("active tenant subscription is required before creating employees")
	ErrSubscriptionEmployeeLimitExceeded = errors.New("subscription employee limit reached")
)

type TenantSubscription struct {
	ID           uuid.UUID  `json:"id"`
	TenantID     uuid.UUID  `json:"tenant_id"`
	PlanID       *uuid.UUID `json:"plan_id,omitempty"`
	StartDate    *time.Time `json:"start_date,omitempty"`
	EndDate      *time.Time `json:"end_date,omitempty"`
	Status       string     `json:"status"`
	MaxEmployees int32      `json:"max_employees"`
	Inactive     bool       `json:"inactive"`
	CreatedAt    time.Time  `json:"created_at"`
	CreatedBy    *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt    time.Time  `json:"updated_at"`
	UpdatedBy    *uuid.UUID `json:"updated_by,omitempty"`
}

type TenantSubscriptionInput struct {
	TenantID     uuid.UUID
	PlanID       *uuid.UUID
	StartDate    *time.Time
	EndDate      *time.Time
	Status       string
	MaxEmployees int32
}

func NewTenantSubscription(input TenantSubscriptionInput) (*TenantSubscription, error) {
	if input.TenantID == uuid.Nil {
		return nil, ErrInvalidTenantID
	}
	status := normalizeSubscriptionStatus(input.Status)
	if !IsValidSubscriptionStatus(status) {
		return nil, ErrInvalidSubscriptionStatus
	}
	if input.MaxEmployees < 0 {
		return nil, ErrInvalidSubscriptionEmployee
	}
	startDate := datePtr(input.StartDate)
	endDate := datePtr(input.EndDate)
	if startDate != nil && endDate != nil && endDate.Before(*startDate) {
		return nil, ErrInvalidSubscriptionPeriod
	}
	now := time.Now().UTC()
	return &TenantSubscription{
		TenantID:     input.TenantID,
		PlanID:       cleanUUIDPtr(input.PlanID),
		StartDate:    startDate,
		EndDate:      endDate,
		Status:       status,
		MaxEmployees: input.MaxEmployees,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

func IsValidSubscriptionStatus(status string) bool {
	switch status {
	case SubscriptionStatusTrialing, SubscriptionStatusActive, SubscriptionStatusPastDue, SubscriptionStatusCancelled, SubscriptionStatusExpired:
		return true
	default:
		return false
	}
}

func normalizeSubscriptionStatus(status string) string {
	status = strings.TrimSpace(strings.ToLower(status))
	if status == "" {
		return SubscriptionStatusActive
	}
	return status
}

func datePtr(value *time.Time) *time.Time {
	if value == nil || value.IsZero() {
		return nil
	}
	date := dateOnly(*value)
	return &date
}

func cleanUUIDPtr(value *uuid.UUID) *uuid.UUID {
	if value == nil || *value == uuid.Nil {
		return nil
	}
	id := *value
	return &id
}
