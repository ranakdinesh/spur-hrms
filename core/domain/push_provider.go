package domain

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	PushProviderLocal = "local"
	PushProviderFCM   = "fcm"
)

var (
	ErrPushProviderSettingsNotFound = errors.New("push provider settings not found")
	ErrInvalidPushProvider          = errors.New("push provider is invalid")
	ErrInvalidPushProjectID         = errors.New("push provider project_id is required")
	ErrInvalidPushClientEmail       = errors.New("push provider client_email is required")
	ErrInvalidPushPrivateKey        = errors.New("push provider private_key is required")
	ErrPushProviderDisabled         = errors.New("push provider is disabled")
	ErrInvalidPushToken             = errors.New("push device token is required")
)

type PushProviderSettings struct {
	ID                 uuid.UUID       `json:"id"`
	TenantID           uuid.UUID       `json:"tenant_id"`
	Provider           string          `json:"provider"`
	IsEnabled          bool            `json:"is_enabled"`
	ProjectID          *string         `json:"project_id,omitempty"`
	ClientEmail        *string         `json:"client_email,omitempty"`
	PrivateKey         *string         `json:"-"`
	PrivateKeyID       *string         `json:"private_key_id,omitempty"`
	AuthURI            *string         `json:"auth_uri,omitempty"`
	TokenURI           *string         `json:"token_uri,omitempty"`
	AndroidEnabled     bool            `json:"android_enabled"`
	IOSEnabled         bool            `json:"ios_enabled"`
	WebEnabled         bool            `json:"web_enabled"`
	DefaultClickAction *string         `json:"default_click_action,omitempty"`
	DefaultImageURL    *string         `json:"default_image_url,omitempty"`
	TTLSeconds         int32           `json:"ttl_seconds"`
	CollapseKey        *string         `json:"collapse_key,omitempty"`
	LastTestAt         *time.Time      `json:"last_test_at,omitempty"`
	LastTestStatus     *string         `json:"last_test_status,omitempty"`
	LastTestMessage    *string         `json:"last_test_message,omitempty"`
	Metadata           json.RawMessage `json:"metadata,omitempty"`
	HasPrivateKey      bool            `json:"has_private_key"`
	ReadinessHints     []string        `json:"readiness_hints,omitempty"`
	Inactive           bool            `json:"inactive"`
	CreatedAt          time.Time       `json:"created_at"`
	CreatedBy          *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt          time.Time       `json:"updated_at"`
	UpdatedBy          *uuid.UUID      `json:"updated_by,omitempty"`
}

type PushProviderSettingsInput struct {
	TenantID           uuid.UUID
	Provider           string
	IsEnabled          bool
	ProjectID          *string
	ClientEmail        *string
	PrivateKey         *string
	PrivateKeyID       *string
	AuthURI            *string
	TokenURI           *string
	AndroidEnabled     bool
	IOSEnabled         bool
	WebEnabled         bool
	DefaultClickAction *string
	DefaultImageURL    *string
	TTLSeconds         int32
	CollapseKey        *string
	Metadata           json.RawMessage
}

func NewPushProviderSettings(input PushProviderSettingsInput) (*PushProviderSettings, error) {
	if input.TenantID == uuid.Nil {
		return nil, ErrInvalidTenantID
	}
	provider, err := ValidatePushProvider(input.Provider)
	if err != nil {
		return nil, err
	}
	projectID := cleanOptional(input.ProjectID)
	clientEmail := cleanOptional(input.ClientEmail)
	privateKey := cleanOptional(input.PrivateKey)
	if provider == PushProviderFCM && input.IsEnabled {
		if projectID == nil {
			return nil, ErrInvalidPushProjectID
		}
		if clientEmail == nil {
			return nil, ErrInvalidPushClientEmail
		}
		if privateKey == nil {
			return nil, ErrInvalidPushPrivateKey
		}
	}
	ttl := input.TTLSeconds
	if ttl <= 0 {
		ttl = 3600
	}
	if ttl > 2419200 {
		ttl = 2419200
	}
	metadata := input.Metadata
	if len(metadata) == 0 {
		metadata = json.RawMessage(`{}`)
	}
	now := time.Now().UTC()
	return &PushProviderSettings{
		TenantID: input.TenantID, Provider: provider, IsEnabled: input.IsEnabled, ProjectID: projectID, ClientEmail: clientEmail, PrivateKey: privateKey, PrivateKeyID: cleanOptional(input.PrivateKeyID), AuthURI: cleanOptional(input.AuthURI), TokenURI: cleanOptional(input.TokenURI), AndroidEnabled: input.AndroidEnabled, IOSEnabled: input.IOSEnabled, WebEnabled: input.WebEnabled, DefaultClickAction: cleanOptional(input.DefaultClickAction), DefaultImageURL: cleanOptional(input.DefaultImageURL), TTLSeconds: ttl, CollapseKey: cleanOptional(input.CollapseKey), Metadata: metadata,
		CreatedAt: now, UpdatedAt: now,
	}, nil
}

func ValidatePushProvider(value string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "", PushProviderLocal:
		return PushProviderLocal, nil
	case PushProviderFCM:
		return PushProviderFCM, nil
	default:
		return "", ErrInvalidPushProvider
	}
}

func (s *PushProviderSettings) RedactSecrets() *PushProviderSettings {
	if s == nil {
		return nil
	}
	s.HasPrivateKey = s.PrivateKey != nil && strings.TrimSpace(*s.PrivateKey) != ""
	s.PrivateKey = nil
	s.ReadinessHints = []string{"Flutter apps should register refreshed FCM tokens with the HRMS device-token API.", "The backend owns targeting, preferences, audit logs, and retries.", "FCM delivers Android pushes and can bridge iOS APNs when the Firebase iOS app is configured.", "Keep a local provider for development and CI without external delivery."}
	return s
}
