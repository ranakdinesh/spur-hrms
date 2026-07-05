package communication

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
	"github.com/ranakdinesh/spur-hrms/internal/logging"
	"github.com/rs/zerolog"
)

const (
	defaultMSG91SMSURL      = "https://control.msg91.com/api/v5/flow/"
	defaultMSG91WhatsAppURL = "https://control.msg91.com/api/v5/whatsapp/whatsapp-outbound-message/bulk/"
	defaultGupshupURL       = "https://api.gupshup.io/sm/api/v1/msg"
)

type MultiProviderSender struct {
	client *http.Client
	log    *zerolog.Logger
}

func NewMultiProviderSender(log *zerolog.Logger) *MultiProviderSender {
	return &MultiProviderSender{client: &http.Client{Timeout: 15 * time.Second}, log: logging.Component(log, "communication_delivery")}
}

func (s *MultiProviderSender) SendSMS(ctx context.Context, settings *domain.CommunicationProviderSettings, message ports.CommunicationMessage) (*ports.CommunicationDeliveryResult, error) {
	if settings == nil {
		return nil, fmt.Errorf("hrms: communication provider settings are required")
	}
	if !settings.SMSEnabled {
		return nil, domain.ErrSMSProviderDisabled
	}
	if strings.TrimSpace(message.Body) == "" {
		return nil, domain.ErrInvalidCommunicationMessage
	}
	switch settings.SMSProvider {
	case domain.CommunicationProviderLocal:
		return s.sendLocal(domain.CommunicationChannelSMS, settings.SMSProvider, message), nil
	case domain.CommunicationProviderMSG91:
		return s.sendMSG91SMS(ctx, settings, message)
	default:
		return nil, domain.ErrInvalidSMSProvider
	}
}

func (s *MultiProviderSender) SendWhatsApp(ctx context.Context, settings *domain.CommunicationProviderSettings, message ports.CommunicationMessage) (*ports.CommunicationDeliveryResult, error) {
	if settings == nil {
		return nil, fmt.Errorf("hrms: communication provider settings are required")
	}
	if !settings.WhatsAppEnabled {
		return nil, domain.ErrWhatsAppProviderDisabled
	}
	if strings.TrimSpace(message.Body) == "" {
		return nil, domain.ErrInvalidCommunicationMessage
	}
	switch settings.WhatsAppProvider {
	case domain.CommunicationProviderLocal:
		return s.sendLocal(domain.CommunicationChannelWhatsApp, settings.WhatsAppProvider, message), nil
	case domain.CommunicationProviderMSG91:
		return s.sendMSG91WhatsApp(ctx, settings, message)
	case domain.CommunicationProviderGupshup:
		return s.sendGupshupWhatsApp(ctx, settings, message)
	default:
		return nil, domain.ErrInvalidWhatsAppProvider
	}
}

func (s *MultiProviderSender) sendLocal(channel string, provider string, message ports.CommunicationMessage) *ports.CommunicationDeliveryResult {
	ref := uuid.NewString()
	if s.log != nil {
		s.log.Info().Str("channel", channel).Str("provider", provider).Str("to", message.To).Str("external_reference_id", ref).Msg("hrms local communication delivery")
	}
	return &ports.CommunicationDeliveryResult{Provider: provider, Status: domain.NotifStatusSent, ExternalReference: ref, MessageID: ref, EventStatus: "accepted"}
}

func (s *MultiProviderSender) sendMSG91SMS(ctx context.Context, settings *domain.CommunicationProviderSettings, message ports.CommunicationMessage) (*ports.CommunicationDeliveryResult, error) {
	if settings.SMSAuthKey == nil || strings.TrimSpace(*settings.SMSAuthKey) == "" {
		return nil, domain.ErrInvalidCommunicationAuthKey
	}
	payload := map[string]any{"recipients": []map[string]string{{"mobiles": message.To}}}
	if settings.SMSTemplateID != nil {
		payload["template_id"] = strings.TrimSpace(*settings.SMSTemplateID)
	}
	if settings.SMSSenderID != nil {
		payload["sender"] = strings.TrimSpace(*settings.SMSSenderID)
	}
	if settings.SMSRoute != nil {
		payload["route"] = strings.TrimSpace(*settings.SMSRoute)
	}
	if settings.SMSCountryCode != nil {
		payload["country"] = strings.TrimSpace(*settings.SMSCountryCode)
	}
	payload["message"] = message.Body
	return s.postJSON(ctx, firstNonEmpty(settings.SMSBaseURL, defaultMSG91SMSURL), map[string]string{"authkey": strings.TrimSpace(*settings.SMSAuthKey), "Content-Type": "application/json"}, payload)
}

func (s *MultiProviderSender) sendMSG91WhatsApp(ctx context.Context, settings *domain.CommunicationProviderSettings, message ports.CommunicationMessage) (*ports.CommunicationDeliveryResult, error) {
	if settings.WhatsAppAuthKey == nil || strings.TrimSpace(*settings.WhatsAppAuthKey) == "" {
		return nil, domain.ErrInvalidCommunicationAuthKey
	}
	templateID := valueOrEmpty(settings.WhatsAppTemplateID)
	if templateID == "" && message.TemplateID != nil {
		templateID = strings.TrimSpace(*message.TemplateID)
	}
	payload := map[string]any{"integrated_number": valueOrEmpty(settings.WhatsAppSourceNumber), "content_type": "template", "payload": map[string]any{"to": message.To, "type": "template", "template": map[string]any{"name": firstNonEmptyString(valueOrEmpty(settings.WhatsAppTemplateName), templateID), "language": map[string]string{"code": "en"}, "components": []any{}}}}
	if message.Body != "" {
		payload["preview_message"] = message.Body
	}
	return s.postJSON(ctx, firstNonEmpty(settings.WhatsAppBaseURL, defaultMSG91WhatsAppURL), map[string]string{"authkey": strings.TrimSpace(*settings.WhatsAppAuthKey), "Content-Type": "application/json"}, payload)
}

func (s *MultiProviderSender) sendGupshupWhatsApp(ctx context.Context, settings *domain.CommunicationProviderSettings, message ports.CommunicationMessage) (*ports.CommunicationDeliveryResult, error) {
	if settings.WhatsAppAuthKey == nil || strings.TrimSpace(*settings.WhatsAppAuthKey) == "" {
		return nil, domain.ErrInvalidCommunicationAuthKey
	}
	form := url.Values{}
	form.Set("channel", "whatsapp")
	form.Set("source", valueOrEmpty(settings.WhatsAppSourceNumber))
	form.Set("destination", message.To)
	form.Set("src.name", valueOrEmpty(settings.WhatsAppAppName))
	form.Set("message", message.Body)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, firstNonEmpty(settings.WhatsAppBaseURL, defaultGupshupURL), strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("apikey", strings.TrimSpace(*settings.WhatsAppAuthKey))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return s.doRequest(req)
}

func (s *MultiProviderSender) postJSON(ctx context.Context, endpoint string, headers map[string]string, payload map[string]any) (*ports.CommunicationDeliveryResult, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	return s.doRequest(req)
}

func (s *MultiProviderSender) doRequest(req *http.Request) (*ports.CommunicationDeliveryResult, error) {
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("communication delivery failed: status=%d body=%s", resp.StatusCode, strings.TrimSpace(string(bodyBytes)))
	}
	ref := resp.Header.Get("X-Request-Id")
	if ref == "" {
		ref = uuid.NewString()
	}
	return &ports.CommunicationDeliveryResult{Status: domain.NotifStatusSent, ExternalReference: ref, MessageID: ref, EventStatus: "accepted"}, nil
}

func firstNonEmpty(value *string, fallback string) string {
	if value == nil || strings.TrimSpace(*value) == "" {
		return fallback
	}
	return strings.TrimSpace(*value)
}

func firstNonEmptyString(value string, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return strings.TrimSpace(value)
}

func valueOrEmpty(value *string) string {
	if value == nil {
		return ""
	}
	return strings.TrimSpace(*value)
}
