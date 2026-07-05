package push

import (
	"bytes"
	"context"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
	"github.com/ranakdinesh/spur-hrms/internal/logging"
	"github.com/rs/zerolog"
)

const firebaseMessagingScope = "https://www.googleapis.com/auth/firebase.messaging"

type MultiProviderSender struct {
	client *http.Client
	log    *zerolog.Logger
}

func NewMultiProviderSender(log *zerolog.Logger) *MultiProviderSender {
	return &MultiProviderSender{client: &http.Client{Timeout: 15 * time.Second}, log: logging.Component(log, "push_delivery")}
}

func (s *MultiProviderSender) SendPush(ctx context.Context, settings *domain.PushProviderSettings, message ports.PushMessage) (*ports.PushDeliveryResult, error) {
	if settings == nil {
		return nil, domain.ErrPushProviderSettingsNotFound
	}
	if !settings.IsEnabled {
		return nil, domain.ErrPushProviderDisabled
	}
	if strings.TrimSpace(message.Token) == "" {
		return nil, domain.ErrInvalidPushToken
	}
	switch settings.Provider {
	case domain.PushProviderLocal:
		return s.sendLocal(settings, message), nil
	case domain.PushProviderFCM:
		return s.sendFCM(ctx, settings, message)
	default:
		return nil, domain.ErrInvalidPushProvider
	}
}

func (s *MultiProviderSender) sendLocal(settings *domain.PushProviderSettings, message ports.PushMessage) *ports.PushDeliveryResult {
	ref := uuid.NewString()
	if s.log != nil {
		s.log.Info().Str("provider", settings.Provider).Str("token", message.Token).Str("title", message.Title).Str("external_reference_id", ref).Msg("hrms local push delivery")
	}
	return &ports.PushDeliveryResult{Provider: settings.Provider, Status: domain.NotifStatusSent, ExternalReference: ref, MessageID: ref, EventStatus: "accepted"}
}

func (s *MultiProviderSender) sendFCM(ctx context.Context, settings *domain.PushProviderSettings, message ports.PushMessage) (*ports.PushDeliveryResult, error) {
	token, err := s.accessToken(ctx, settings)
	if err != nil {
		return nil, err
	}
	projectID := valueOrEmpty(settings.ProjectID)
	payload := map[string]any{
		"message": map[string]any{
			"token":        strings.TrimSpace(message.Token),
			"notification": map[string]string{"title": message.Title, "body": message.Body},
			"data":         message.Data,
			"android": map[string]any{
				"ttl": fmt.Sprintf("%ds", settings.TTLSeconds),
			},
			"apns": map[string]any{
				"payload": map[string]any{"aps": map[string]any{"sound": "default"}},
			},
		},
	}
	msg := payload["message"].(map[string]any)
	if settings.DefaultImageURL != nil && strings.TrimSpace(*settings.DefaultImageURL) != "" {
		msg["notification"].(map[string]string)["image"] = strings.TrimSpace(*settings.DefaultImageURL)
	}
	if settings.CollapseKey != nil && strings.TrimSpace(*settings.CollapseKey) != "" {
		msg["android"].(map[string]any)["collapse_key"] = strings.TrimSpace(*settings.CollapseKey)
	}
	if settings.DefaultClickAction != nil && strings.TrimSpace(*settings.DefaultClickAction) != "" {
		msg["webpush"] = map[string]any{"fcm_options": map[string]string{"link": strings.TrimSpace(*settings.DefaultClickAction)}}
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://fcm.googleapis.com/v1/projects/"+projectID+"/messages:send", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("fcm delivery failed: status=%d body=%s", resp.StatusCode, strings.TrimSpace(string(bodyBytes)))
	}
	var parsed struct {
		Name string `json:"name"`
	}
	_ = json.Unmarshal(bodyBytes, &parsed)
	if parsed.Name == "" {
		parsed.Name = uuid.NewString()
	}
	return &ports.PushDeliveryResult{Provider: settings.Provider, Status: domain.NotifStatusSent, ExternalReference: parsed.Name, MessageID: parsed.Name, EventStatus: "accepted"}, nil
}

func (s *MultiProviderSender) accessToken(ctx context.Context, settings *domain.PushProviderSettings) (string, error) {
	if settings.ProjectID == nil || strings.TrimSpace(*settings.ProjectID) == "" {
		return "", domain.ErrInvalidPushProjectID
	}
	if settings.ClientEmail == nil || strings.TrimSpace(*settings.ClientEmail) == "" {
		return "", domain.ErrInvalidPushClientEmail
	}
	if settings.PrivateKey == nil || strings.TrimSpace(*settings.PrivateKey) == "" {
		return "", domain.ErrInvalidPushPrivateKey
	}
	privateKey, err := parsePrivateKey(*settings.PrivateKey)
	if err != nil {
		return "", err
	}
	tokenURI := valueOrDefault(settings.TokenURI, "https://oauth2.googleapis.com/token")
	now := time.Now()
	claims := jwt.MapClaims{"iss": strings.TrimSpace(*settings.ClientEmail), "scope": firebaseMessagingScope, "aud": tokenURI, "iat": now.Unix(), "exp": now.Add(time.Hour).Unix()}
	j := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	if settings.PrivateKeyID != nil && strings.TrimSpace(*settings.PrivateKeyID) != "" {
		j.Header["kid"] = strings.TrimSpace(*settings.PrivateKeyID)
	}
	assertion, err := j.SignedString(privateKey)
	if err != nil {
		return "", err
	}
	form := "grant_type=urn%3Aietf%3Aparams%3Aoauth%3Agrant-type%3Ajwt-bearer&assertion=" + assertion
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURI, strings.NewReader(form))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("fcm oauth token failed: status=%d body=%s", resp.StatusCode, strings.TrimSpace(string(bodyBytes)))
	}
	var token struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.Unmarshal(bodyBytes, &token); err != nil {
		return "", err
	}
	if token.AccessToken == "" {
		return "", fmt.Errorf("fcm oauth token response missing access_token")
	}
	return token.AccessToken, nil
}

func parsePrivateKey(value string) (*rsa.PrivateKey, error) {
	clean := strings.ReplaceAll(strings.TrimSpace(value), `\n`, "\n")
	return jwt.ParseRSAPrivateKeyFromPEM([]byte(clean))
}

func valueOrEmpty(value *string) string {
	return valueOrDefault(value, "")
}

func valueOrDefault(value *string, fallback string) string {
	if value == nil || strings.TrimSpace(*value) == "" {
		return fallback
	}
	return strings.TrimSpace(*value)
}
