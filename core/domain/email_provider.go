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
	EmailProviderLocal    = "local"
	EmailProviderSMTP     = "smtp"
	EmailProviderSendGrid = "sendgrid"

	EmailEncryptionNone     = "none"
	EmailEncryptionStartTLS = "starttls"
	EmailEncryptionTLS      = "tls"
)

var (
	ErrEmailProviderSettingsNotFound    = errors.New("email provider settings not found")
	ErrInvalidEmailProvider             = errors.New("email provider is invalid")
	ErrInvalidEmailFromAddress          = errors.New("email from address is invalid")
	ErrInvalidEmailSMTPHost             = errors.New("smtp host is required")
	ErrInvalidEmailSMTPPort             = errors.New("smtp port is invalid")
	ErrInvalidEmailEncryption           = errors.New("smtp encryption is invalid")
	ErrInvalidEmailAPIKey               = errors.New("email provider api key is required")
	ErrInvalidEmailRecipient            = errors.New("email recipient is invalid")
	ErrGlobalEmailProviderReadOnly      = errors.New("global email provider settings are controlled by environment configuration")
	ErrGlobalEmailProviderNotConfigured = errors.New("global email provider is not configured")
)

var emailAddressPattern = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)

type EmailProviderSettings struct {
	ID                   uuid.UUID       `json:"id"`
	TenantID             uuid.UUID       `json:"tenant_id"`
	Provider             string          `json:"provider"`
	IsEnabled            bool            `json:"is_enabled"`
	FromName             *string         `json:"from_name,omitempty"`
	FromEmail            string          `json:"from_email"`
	ReplyToEmail         *string         `json:"reply_to_email,omitempty"`
	SMTPHost             *string         `json:"smtp_host,omitempty"`
	SMTPPort             *int32          `json:"smtp_port,omitempty"`
	SMTPUsername         *string         `json:"smtp_username,omitempty"`
	SMTPPassword         *string         `json:"-"`
	SMTPEncryption       string          `json:"smtp_encryption"`
	SendGridAPIKey       *string         `json:"-"`
	SendGridSandboxMode  bool            `json:"sendgrid_sandbox_mode"`
	WebhookSigningSecret *string         `json:"-"`
	SPFStatus            *string         `json:"spf_status,omitempty"`
	DKIMStatus           *string         `json:"dkim_status,omitempty"`
	DMARCStatus          *string         `json:"dmarc_status,omitempty"`
	LastTestAt           *time.Time      `json:"last_test_at,omitempty"`
	LastTestStatus       *string         `json:"last_test_status,omitempty"`
	LastTestMessage      *string         `json:"last_test_message,omitempty"`
	Metadata             json.RawMessage `json:"metadata,omitempty"`
	HasSMTPPassword      bool            `json:"has_smtp_password"`
	HasSendGridAPIKey    bool            `json:"has_sendgrid_api_key"`
	HasWebhookSecret     bool            `json:"has_webhook_secret"`
	DeliverabilityHints  []string        `json:"deliverability_hints,omitempty"`
	Inactive             bool            `json:"inactive"`
	CreatedAt            time.Time       `json:"created_at"`
	CreatedBy            *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt            time.Time       `json:"updated_at"`
	UpdatedBy            *uuid.UUID      `json:"updated_by,omitempty"`
}

type EmailProviderSettingsInput struct {
	TenantID             uuid.UUID
	Provider             string
	IsEnabled            bool
	FromName             *string
	FromEmail            string
	ReplyToEmail         *string
	SMTPHost             *string
	SMTPPort             *int32
	SMTPUsername         *string
	SMTPPassword         *string
	SMTPEncryption       string
	SendGridAPIKey       *string
	SendGridSandboxMode  bool
	WebhookSigningSecret *string
	SPFStatus            *string
	DKIMStatus           *string
	DMARCStatus          *string
	Metadata             json.RawMessage
}

func NewEmailProviderSettings(input EmailProviderSettingsInput) (*EmailProviderSettings, error) {
	if input.TenantID == uuid.Nil {
		return nil, ErrInvalidTenantID
	}
	provider, err := ValidateEmailProvider(input.Provider)
	if err != nil {
		return nil, err
	}
	fromEmail := strings.ToLower(strings.TrimSpace(input.FromEmail))
	if !emailAddressPattern.MatchString(fromEmail) {
		return nil, ErrInvalidEmailFromAddress
	}
	replyTo := cleanOptional(input.ReplyToEmail)
	if replyTo != nil {
		value := strings.ToLower(strings.TrimSpace(*replyTo))
		if !emailAddressPattern.MatchString(value) {
			return nil, ErrInvalidEmailFromAddress
		}
		replyTo = &value
	}
	encryption, err := ValidateEmailEncryption(input.SMTPEncryption)
	if err != nil {
		return nil, err
	}
	if provider == EmailProviderSMTP && input.IsEnabled {
		if input.SMTPHost == nil || strings.TrimSpace(*input.SMTPHost) == "" {
			return nil, ErrInvalidEmailSMTPHost
		}
		if input.SMTPPort == nil || *input.SMTPPort <= 0 || *input.SMTPPort > 65535 {
			return nil, ErrInvalidEmailSMTPPort
		}
	}
	if provider == EmailProviderSendGrid && input.IsEnabled && (input.SendGridAPIKey == nil || strings.TrimSpace(*input.SendGridAPIKey) == "") {
		return nil, ErrInvalidEmailAPIKey
	}
	metadata := input.Metadata
	if len(metadata) == 0 {
		metadata = json.RawMessage(`{}`)
	}
	now := time.Now().UTC()
	return &EmailProviderSettings{
		TenantID: input.TenantID, Provider: provider, IsEnabled: input.IsEnabled, FromName: cleanOptional(input.FromName), FromEmail: fromEmail, ReplyToEmail: replyTo,
		SMTPHost: cleanOptional(input.SMTPHost), SMTPPort: input.SMTPPort, SMTPUsername: cleanOptional(input.SMTPUsername), SMTPPassword: cleanOptional(input.SMTPPassword), SMTPEncryption: encryption,
		SendGridAPIKey: cleanOptional(input.SendGridAPIKey), SendGridSandboxMode: input.SendGridSandboxMode, WebhookSigningSecret: cleanOptional(input.WebhookSigningSecret),
		SPFStatus: cleanOptional(input.SPFStatus), DKIMStatus: cleanOptional(input.DKIMStatus), DMARCStatus: cleanOptional(input.DMARCStatus), Metadata: metadata,
		CreatedAt: now, UpdatedAt: now,
	}, nil
}

func ValidateEmailProvider(value string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "", EmailProviderSMTP:
		return EmailProviderSMTP, nil
	case EmailProviderLocal:
		return EmailProviderLocal, nil
	case EmailProviderSendGrid:
		return EmailProviderSendGrid, nil
	default:
		return "", ErrInvalidEmailProvider
	}
}

func ValidateEmailEncryption(value string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "", EmailEncryptionStartTLS:
		return EmailEncryptionStartTLS, nil
	case EmailEncryptionNone:
		return EmailEncryptionNone, nil
	case EmailEncryptionTLS:
		return EmailEncryptionTLS, nil
	default:
		return "", ErrInvalidEmailEncryption
	}
}

func (s *EmailProviderSettings) RedactSecrets() *EmailProviderSettings {
	if s == nil {
		return nil
	}
	s.HasSMTPPassword = s.SMTPPassword != nil && strings.TrimSpace(*s.SMTPPassword) != ""
	s.HasSendGridAPIKey = s.SendGridAPIKey != nil && strings.TrimSpace(*s.SendGridAPIKey) != ""
	s.HasWebhookSecret = s.WebhookSigningSecret != nil && strings.TrimSpace(*s.WebhookSigningSecret) != ""
	s.SMTPPassword = nil
	s.SendGridAPIKey = nil
	s.WebhookSigningSecret = nil
	s.DeliverabilityHints = []string{"Verify SPF for the sending domain.", "Enable DKIM signing with the selected provider.", "Publish a DMARC policy before high-volume sends.", "Use provider webhooks for bounces and complaints."}
	return s
}
