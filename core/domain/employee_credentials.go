package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

const (
	CredentialEventResendCredentials      = "resend_credentials"
	CredentialEventResetTemporaryPassword = "reset_temporary_password"
	CredentialDeliveryEmail               = "email"
	CredentialDeliverySent                = "sent"
	CredentialDeliveryFailed              = "failed"
)

var (
	ErrEmployeeCredentialPortMissing = errors.New("employee credential identity port is not configured")
	ErrInvalidCredentialEventType    = errors.New("invalid credential event type")
	ErrInvalidTemporaryPassword      = errors.New("temporary password must be at least 8 characters")
	ErrEmployeeCredentialTarget      = errors.New("employee email is required for credential delivery")
)

type EmployeeCredentialEvent struct {
	ID              uuid.UUID  `json:"id"`
	TenantID        uuid.UUID  `json:"tenant_id"`
	EmployeeID      uuid.UUID  `json:"employee_id"`
	UserID          uuid.UUID  `json:"user_id"`
	EventType       string     `json:"event_type"`
	DeliveryChannel string     `json:"delivery_channel"`
	DeliveryTarget  string     `json:"delivery_target"`
	Status          string     `json:"status"`
	FailureReason   *string    `json:"failure_reason,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	CreatedBy       *uuid.UUID `json:"created_by,omitempty"`
}
