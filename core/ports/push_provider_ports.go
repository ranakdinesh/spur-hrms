package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

type PushProviderRepo interface {
	GetPushProviderSettings(ctx context.Context, tenantID uuid.UUID) (*domain.PushProviderSettings, error)
	UpsertPushProviderSettings(ctx context.Context, item *domain.PushProviderSettings, actorID *uuid.UUID) (*domain.PushProviderSettings, error)
	UpdatePushProviderTestResult(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, message *string, actorID *uuid.UUID) (*domain.PushProviderSettings, error)
	DeletePushProviderSettings(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	UpdateNotificationLogDelivery(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, update NotificationDeliveryUpdate, actorID *uuid.UUID) (*domain.NotificationLog, error)
}

type PushDeliverySender interface {
	SendPush(ctx context.Context, settings *domain.PushProviderSettings, message PushMessage) (*PushDeliveryResult, error)
}

type PushMessage struct {
	TenantID       uuid.UUID
	Token          string
	Title          string
	Body           string
	Data           map[string]string
	IdempotencyKey *string
	ReferenceID    *uuid.UUID
}

type PushDeliveryResult struct {
	Provider          string
	Status            string
	ExternalReference string
	MessageID         string
	EventStatus       string
}

type PushProviderSettingsCommand struct {
	TenantID           uuid.UUID  `json:"tenant_id"`
	Provider           string     `json:"provider"`
	IsEnabled          bool       `json:"is_enabled"`
	ProjectID          *string    `json:"project_id,omitempty"`
	ClientEmail        *string    `json:"client_email,omitempty"`
	PrivateKey         *string    `json:"private_key,omitempty"`
	PrivateKeyID       *string    `json:"private_key_id,omitempty"`
	AuthURI            *string    `json:"auth_uri,omitempty"`
	TokenURI           *string    `json:"token_uri,omitempty"`
	AndroidEnabled     bool       `json:"android_enabled"`
	IOSEnabled         bool       `json:"ios_enabled"`
	WebEnabled         bool       `json:"web_enabled"`
	DefaultClickAction *string    `json:"default_click_action,omitempty"`
	DefaultImageURL    *string    `json:"default_image_url,omitempty"`
	TTLSeconds         int32      `json:"ttl_seconds"`
	CollapseKey        *string    `json:"collapse_key,omitempty"`
	ActorID            *uuid.UUID `json:"-"`
}

type PushProviderTestCommand struct {
	TenantID uuid.UUID  `json:"tenant_id"`
	Token    string     `json:"token"`
	Title    *string    `json:"title,omitempty"`
	Message  *string    `json:"message,omitempty"`
	ActorID  *uuid.UUID `json:"-"`
}
