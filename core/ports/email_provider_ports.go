package ports

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

type EmailProviderRepo interface {
	GetEmailProviderSettings(ctx context.Context, tenantID uuid.UUID) (*domain.EmailProviderSettings, error)
	UpsertEmailProviderSettings(ctx context.Context, item *domain.EmailProviderSettings, actorID *uuid.UUID) (*domain.EmailProviderSettings, error)
	UpdateEmailProviderTestResult(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, message *string, actorID *uuid.UUID) (*domain.EmailProviderSettings, error)
	DeleteEmailProviderSettings(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	UpdateNotificationLogDelivery(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, update NotificationDeliveryUpdate, actorID *uuid.UUID) (*domain.NotificationLog, error)
}

type EmailDeliverySender interface {
	SendEmail(ctx context.Context, settings *domain.EmailProviderSettings, message EmailMessage) (*EmailDeliveryResult, error)
}

type EmailMessage struct {
	TenantID       uuid.UUID
	To             string
	Subject        string
	TextBody       string
	HTMLBody       *string
	IdempotencyKey *string
	ReferenceID    *uuid.UUID
}

type EmailDeliveryResult struct {
	Provider          string
	Status            string
	ExternalReference string
	MessageID         string
	EventStatus       string
}

type NotificationDeliveryUpdate struct {
	Status              string
	ErrorMessage        *string
	ExternalReferenceID *string
	Provider            *string
	ProviderMessageID   *string
	ProviderEventStatus *string
	ProviderEventAt     *time.Time
}

type EmailProviderSettingsCommand struct {
	TenantID             uuid.UUID  `json:"tenant_id"`
	Provider             string     `json:"provider"`
	IsEnabled            bool       `json:"is_enabled"`
	FromName             *string    `json:"from_name,omitempty"`
	FromEmail            string     `json:"from_email"`
	ReplyToEmail         *string    `json:"reply_to_email,omitempty"`
	SMTPHost             *string    `json:"smtp_host,omitempty"`
	SMTPPort             *int32     `json:"smtp_port,omitempty"`
	SMTPUsername         *string    `json:"smtp_username,omitempty"`
	SMTPPassword         *string    `json:"smtp_password,omitempty"`
	SMTPEncryption       string     `json:"smtp_encryption"`
	SendGridAPIKey       *string    `json:"sendgrid_api_key,omitempty"`
	SendGridSandboxMode  bool       `json:"sendgrid_sandbox_mode"`
	WebhookSigningSecret *string    `json:"webhook_signing_secret,omitempty"`
	SPFStatus            *string    `json:"spf_status,omitempty"`
	DKIMStatus           *string    `json:"dkim_status,omitempty"`
	DMARCStatus          *string    `json:"dmarc_status,omitempty"`
	ActorID              *uuid.UUID `json:"-"`
}

type EmailProviderTestCommand struct {
	TenantID uuid.UUID  `json:"tenant_id"`
	ToEmail  string     `json:"to_email"`
	Subject  *string    `json:"subject,omitempty"`
	Message  *string    `json:"message,omitempty"`
	ActorID  *uuid.UUID `json:"-"`
}
