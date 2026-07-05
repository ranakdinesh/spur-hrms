package domain

import (
	"encoding/json"
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	CommunicationChannelSMS      = "sms"
	CommunicationChannelWhatsApp = "whatsapp"

	CommunicationProviderLocal   = "local"
	CommunicationProviderMSG91   = "msg91"
	CommunicationProviderGupshup = "gupshup"
)

var (
	ErrCommunicationProviderSettingsNotFound = errors.New("communication provider settings not found")
	ErrInvalidCommunicationChannel           = errors.New("communication channel is invalid")
	ErrInvalidSMSProvider                    = errors.New("sms provider is invalid")
	ErrInvalidWhatsAppProvider               = errors.New("whatsapp provider is invalid")
	ErrInvalidCommunicationAuthKey           = errors.New("communication provider auth key is required")
	ErrInvalidCommunicationRecipient         = errors.New("communication recipient is invalid")
	ErrInvalidCommunicationMessage           = errors.New("communication message is required")
	ErrSMSProviderDisabled                   = errors.New("sms provider is disabled")
	ErrWhatsAppProviderDisabled              = errors.New("whatsapp provider is disabled")
)

var phonePattern = regexp.MustCompile(`^\+?[1-9]\d{7,14}$`)

type CommunicationProviderSettings struct {
	ID                      uuid.UUID       `json:"id"`
	TenantID                uuid.UUID       `json:"tenant_id"`
	SMSProvider             string          `json:"sms_provider"`
	SMSEnabled              bool            `json:"sms_enabled"`
	SMSSenderID             *string         `json:"sms_sender_id,omitempty"`
	SMSAuthKey              *string         `json:"-"`
	SMSTemplateID           *string         `json:"sms_template_id,omitempty"`
	SMSRoute                *string         `json:"sms_route,omitempty"`
	SMSCountryCode          *string         `json:"sms_country_code,omitempty"`
	SMSBaseURL              *string         `json:"sms_base_url,omitempty"`
	WhatsAppProvider        string          `json:"whatsapp_provider"`
	WhatsAppEnabled         bool            `json:"whatsapp_enabled"`
	WhatsAppAuthKey         *string         `json:"-"`
	WhatsAppAppName         *string         `json:"whatsapp_app_name,omitempty"`
	WhatsAppSourceNumber    *string         `json:"whatsapp_source_number,omitempty"`
	WhatsAppTemplateID      *string         `json:"whatsapp_template_id,omitempty"`
	WhatsAppTemplateName    *string         `json:"whatsapp_template_name,omitempty"`
	WhatsAppNamespace       *string         `json:"whatsapp_namespace,omitempty"`
	WhatsAppBaseURL         *string         `json:"whatsapp_base_url,omitempty"`
	WebhookSigningSecret    *string         `json:"-"`
	SMSLastTestAt           *time.Time      `json:"sms_last_test_at,omitempty"`
	SMSLastTestStatus       *string         `json:"sms_last_test_status,omitempty"`
	SMSLastTestMessage      *string         `json:"sms_last_test_message,omitempty"`
	WhatsAppLastTestAt      *time.Time      `json:"whatsapp_last_test_at,omitempty"`
	WhatsAppLastTestStatus  *string         `json:"whatsapp_last_test_status,omitempty"`
	WhatsAppLastTestMessage *string         `json:"whatsapp_last_test_message,omitempty"`
	Metadata                json.RawMessage `json:"metadata,omitempty"`
	HasSMSAuthKey           bool            `json:"has_sms_auth_key"`
	HasWhatsAppAuthKey      bool            `json:"has_whatsapp_auth_key"`
	HasWebhookSecret        bool            `json:"has_webhook_secret"`
	ReadinessHints          []string        `json:"readiness_hints,omitempty"`
	Inactive                bool            `json:"inactive"`
	CreatedAt               time.Time       `json:"created_at"`
	CreatedBy               *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt               time.Time       `json:"updated_at"`
	UpdatedBy               *uuid.UUID      `json:"updated_by,omitempty"`
}

type CommunicationProviderSettingsInput struct {
	TenantID             uuid.UUID
	SMSProvider          string
	SMSEnabled           bool
	SMSSenderID          *string
	SMSAuthKey           *string
	SMSTemplateID        *string
	SMSRoute             *string
	SMSCountryCode       *string
	SMSBaseURL           *string
	WhatsAppProvider     string
	WhatsAppEnabled      bool
	WhatsAppAuthKey      *string
	WhatsAppAppName      *string
	WhatsAppSourceNumber *string
	WhatsAppTemplateID   *string
	WhatsAppTemplateName *string
	WhatsAppNamespace    *string
	WhatsAppBaseURL      *string
	WebhookSigningSecret *string
	Metadata             json.RawMessage
}

func NewCommunicationProviderSettings(input CommunicationProviderSettingsInput) (*CommunicationProviderSettings, error) {
	if input.TenantID == uuid.Nil {
		return nil, ErrInvalidTenantID
	}
	smsProvider, err := ValidateSMSProvider(input.SMSProvider)
	if err != nil {
		return nil, err
	}
	whatsAppProvider, err := ValidateWhatsAppProvider(input.WhatsAppProvider)
	if err != nil {
		return nil, err
	}
	if input.SMSEnabled && smsProvider == CommunicationProviderMSG91 && input.SMSAuthKey == nil {
		return nil, ErrInvalidCommunicationAuthKey
	}
	if input.WhatsAppEnabled && (whatsAppProvider == CommunicationProviderMSG91 || whatsAppProvider == CommunicationProviderGupshup) && input.WhatsAppAuthKey == nil {
		return nil, ErrInvalidCommunicationAuthKey
	}
	metadata := input.Metadata
	if len(metadata) == 0 {
		metadata = json.RawMessage(`{}`)
	}
	now := time.Now().UTC()
	return &CommunicationProviderSettings{
		TenantID: input.TenantID, SMSProvider: smsProvider, SMSEnabled: input.SMSEnabled, SMSSenderID: cleanOptional(input.SMSSenderID), SMSAuthKey: cleanOptional(input.SMSAuthKey), SMSTemplateID: cleanOptional(input.SMSTemplateID), SMSRoute: cleanOptional(input.SMSRoute), SMSCountryCode: cleanOptional(input.SMSCountryCode), SMSBaseURL: cleanOptional(input.SMSBaseURL),
		WhatsAppProvider: whatsAppProvider, WhatsAppEnabled: input.WhatsAppEnabled, WhatsAppAuthKey: cleanOptional(input.WhatsAppAuthKey), WhatsAppAppName: cleanOptional(input.WhatsAppAppName), WhatsAppSourceNumber: cleanOptional(input.WhatsAppSourceNumber), WhatsAppTemplateID: cleanOptional(input.WhatsAppTemplateID), WhatsAppTemplateName: cleanOptional(input.WhatsAppTemplateName), WhatsAppNamespace: cleanOptional(input.WhatsAppNamespace), WhatsAppBaseURL: cleanOptional(input.WhatsAppBaseURL), WebhookSigningSecret: cleanOptional(input.WebhookSigningSecret),
		Metadata: metadata, CreatedAt: now, UpdatedAt: now,
	}, nil
}

func ValidateSMSProvider(value string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "", CommunicationProviderLocal:
		return CommunicationProviderLocal, nil
	case CommunicationProviderMSG91:
		return CommunicationProviderMSG91, nil
	default:
		return "", ErrInvalidSMSProvider
	}
}

func ValidateWhatsAppProvider(value string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "", CommunicationProviderLocal:
		return CommunicationProviderLocal, nil
	case CommunicationProviderMSG91:
		return CommunicationProviderMSG91, nil
	case CommunicationProviderGupshup:
		return CommunicationProviderGupshup, nil
	default:
		return "", ErrInvalidWhatsAppProvider
	}
}

func ValidateCommunicationChannel(value string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case CommunicationChannelSMS:
		return CommunicationChannelSMS, nil
	case CommunicationChannelWhatsApp:
		return CommunicationChannelWhatsApp, nil
	default:
		return "", ErrInvalidCommunicationChannel
	}
}

func ValidateCommunicationRecipient(value string) (string, error) {
	clean := strings.ReplaceAll(strings.TrimSpace(value), " ", "")
	if !phonePattern.MatchString(clean) {
		return "", ErrInvalidCommunicationRecipient
	}
	return clean, nil
}

func (s *CommunicationProviderSettings) RedactSecrets() *CommunicationProviderSettings {
	if s == nil {
		return nil
	}
	s.HasSMSAuthKey = s.SMSAuthKey != nil && strings.TrimSpace(*s.SMSAuthKey) != ""
	s.HasWhatsAppAuthKey = s.WhatsAppAuthKey != nil && strings.TrimSpace(*s.WhatsAppAuthKey) != ""
	s.HasWebhookSecret = s.WebhookSigningSecret != nil && strings.TrimSpace(*s.WebhookSigningSecret) != ""
	s.SMSAuthKey = nil
	s.WhatsAppAuthKey = nil
	s.WebhookSigningSecret = nil
	s.ReadinessHints = []string{"Keep SMS and WhatsApp disabled until tenant opt-in is complete.", "Use approved WhatsApp templates for business-initiated messages.", "Track delivery receipts and failures through provider webhooks.", "Keep OTP, HR alerts, and marketing messages separated by template and consent."}
	return s
}
