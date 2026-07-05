package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

type CommunicationProviderRepo interface {
	GetCommunicationProviderSettings(ctx context.Context, tenantID uuid.UUID) (*domain.CommunicationProviderSettings, error)
	UpsertCommunicationProviderSettings(ctx context.Context, item *domain.CommunicationProviderSettings, actorID *uuid.UUID) (*domain.CommunicationProviderSettings, error)
	UpdateCommunicationProviderTestResult(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, channel string, status string, message *string, actorID *uuid.UUID) (*domain.CommunicationProviderSettings, error)
	DeleteCommunicationProviderSettings(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
}

type CommunicationDeliverySender interface {
	SendSMS(ctx context.Context, settings *domain.CommunicationProviderSettings, message CommunicationMessage) (*CommunicationDeliveryResult, error)
	SendWhatsApp(ctx context.Context, settings *domain.CommunicationProviderSettings, message CommunicationMessage) (*CommunicationDeliveryResult, error)
}

type CommunicationMessage struct {
	TenantID       uuid.UUID
	To             string
	Body           string
	TemplateID     *string
	TemplateName   *string
	Variables      map[string]string
	IdempotencyKey *string
	ReferenceID    *uuid.UUID
}

type CommunicationDeliveryResult struct {
	Provider          string
	Status            string
	ExternalReference string
	MessageID         string
	EventStatus       string
}

type CommunicationProviderSettingsCommand struct {
	TenantID             uuid.UUID  `json:"tenant_id"`
	SMSProvider          string     `json:"sms_provider"`
	SMSEnabled           bool       `json:"sms_enabled"`
	SMSSenderID          *string    `json:"sms_sender_id,omitempty"`
	SMSAuthKey           *string    `json:"sms_auth_key,omitempty"`
	SMSTemplateID        *string    `json:"sms_template_id,omitempty"`
	SMSRoute             *string    `json:"sms_route,omitempty"`
	SMSCountryCode       *string    `json:"sms_country_code,omitempty"`
	SMSBaseURL           *string    `json:"sms_base_url,omitempty"`
	WhatsAppProvider     string     `json:"whatsapp_provider"`
	WhatsAppEnabled      bool       `json:"whatsapp_enabled"`
	WhatsAppAuthKey      *string    `json:"whatsapp_auth_key,omitempty"`
	WhatsAppAppName      *string    `json:"whatsapp_app_name,omitempty"`
	WhatsAppSourceNumber *string    `json:"whatsapp_source_number,omitempty"`
	WhatsAppTemplateID   *string    `json:"whatsapp_template_id,omitempty"`
	WhatsAppTemplateName *string    `json:"whatsapp_template_name,omitempty"`
	WhatsAppNamespace    *string    `json:"whatsapp_namespace,omitempty"`
	WhatsAppBaseURL      *string    `json:"whatsapp_base_url,omitempty"`
	WebhookSigningSecret *string    `json:"webhook_signing_secret,omitempty"`
	ActorID              *uuid.UUID `json:"-"`
}

type CommunicationProviderTestCommand struct {
	TenantID uuid.UUID  `json:"tenant_id"`
	Channel  string     `json:"channel"`
	ToPhone  string     `json:"to_phone"`
	Message  *string    `json:"message,omitempty"`
	ActorID  *uuid.UUID `json:"-"`
}
