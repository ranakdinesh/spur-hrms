package email

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net"
	"net/http"
	"net/mail"
	"net/smtp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
	"github.com/ranakdinesh/spur-hrms/internal/logging"
	"github.com/rs/zerolog"
)

type MultiProviderSender struct {
	client *http.Client
	log    *zerolog.Logger
}

func NewMultiProviderSender(log *zerolog.Logger) *MultiProviderSender {
	return &MultiProviderSender{client: &http.Client{Timeout: 15 * time.Second}, log: logging.Component(log, "email_delivery")}
}

func (s *MultiProviderSender) SendEmail(ctx context.Context, settings *domain.EmailProviderSettings, message ports.EmailMessage) (*ports.EmailDeliveryResult, error) {
	if settings == nil {
		return nil, fmt.Errorf("hrms: email provider settings are required")
	}
	if !settings.IsEnabled {
		return nil, fmt.Errorf("hrms: email provider is disabled")
	}
	if strings.TrimSpace(message.To) == "" {
		return nil, domain.ErrInvalidEmailRecipient
	}
	switch settings.Provider {
	case domain.EmailProviderLocal:
		return s.sendLocal(settings, message), nil
	case domain.EmailProviderSMTP:
		return s.sendSMTP(ctx, settings, message)
	case domain.EmailProviderSendGrid:
		return s.sendSendGrid(ctx, settings, message)
	default:
		return nil, domain.ErrInvalidEmailProvider
	}
}

func (s *MultiProviderSender) sendLocal(settings *domain.EmailProviderSettings, message ports.EmailMessage) *ports.EmailDeliveryResult {
	ref := uuid.NewString()
	if s.log != nil {
		s.log.Info().Str("provider", settings.Provider).Str("to", message.To).Str("subject", message.Subject).Str("external_reference_id", ref).Msg("hrms local email delivery")
	}
	return &ports.EmailDeliveryResult{Provider: settings.Provider, Status: domain.NotifStatusSent, ExternalReference: ref, MessageID: ref, EventStatus: "accepted"}
}

func (s *MultiProviderSender) sendSMTP(ctx context.Context, settings *domain.EmailProviderSettings, message ports.EmailMessage) (*ports.EmailDeliveryResult, error) {
	if settings.SMTPHost == nil || settings.SMTPPort == nil {
		return nil, domain.ErrInvalidEmailSMTPHost
	}
	host := strings.TrimSpace(*settings.SMTPHost)
	port := *settings.SMTPPort
	addr := fmt.Sprintf("%s:%d", host, port)
	from := mail.Address{Name: valueOrEmpty(settings.FromName), Address: settings.FromEmail}
	to := strings.TrimSpace(message.To)
	raw, err := buildMIMEMessage(from, settings.ReplyToEmail, to, message.Subject, message.TextBody, message.HTMLBody)
	if err != nil {
		return nil, err
	}
	dialer := &net.Dialer{Timeout: 15 * time.Second}
	var conn net.Conn
	if settings.SMTPEncryption == domain.EmailEncryptionTLS {
		conn, err = tls.DialWithDialer(dialer, "tcp", addr, &tls.Config{ServerName: host, MinVersion: tls.VersionTLS12})
	} else {
		conn, err = dialer.DialContext(ctx, "tcp", addr)
	}
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return nil, err
	}
	defer client.Quit()
	if settings.SMTPEncryption == domain.EmailEncryptionStartTLS {
		if ok, _ := client.Extension("STARTTLS"); ok {
			if err := client.StartTLS(&tls.Config{ServerName: host, MinVersion: tls.VersionTLS12}); err != nil {
				return nil, err
			}
		}
	}
	if settings.SMTPUsername != nil && strings.TrimSpace(*settings.SMTPUsername) != "" {
		auth := smtp.PlainAuth("", strings.TrimSpace(*settings.SMTPUsername), valueOrEmpty(settings.SMTPPassword), host)
		if err := client.Auth(auth); err != nil {
			return nil, err
		}
	}
	if err := client.Mail(settings.FromEmail); err != nil {
		return nil, err
	}
	if err := client.Rcpt(to); err != nil {
		return nil, err
	}
	writer, err := client.Data()
	if err != nil {
		return nil, err
	}
	if _, err := writer.Write(raw); err != nil {
		_ = writer.Close()
		return nil, err
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}
	ref := uuid.NewString()
	return &ports.EmailDeliveryResult{Provider: settings.Provider, Status: domain.NotifStatusSent, ExternalReference: ref, MessageID: ref, EventStatus: "accepted"}, nil
}

func (s *MultiProviderSender) sendSendGrid(ctx context.Context, settings *domain.EmailProviderSettings, message ports.EmailMessage) (*ports.EmailDeliveryResult, error) {
	if settings.SendGridAPIKey == nil || strings.TrimSpace(*settings.SendGridAPIKey) == "" {
		return nil, domain.ErrInvalidEmailAPIKey
	}
	body := map[string]any{
		"personalizations": []map[string]any{{"to": []map[string]string{{"email": message.To}}}},
		"from":             map[string]string{"email": settings.FromEmail, "name": valueOrEmpty(settings.FromName)},
		"subject":          message.Subject,
		"content":          []map[string]string{{"type": "text/plain", "value": message.TextBody}},
	}
	if settings.ReplyToEmail != nil {
		body["reply_to"] = map[string]string{"email": strings.TrimSpace(*settings.ReplyToEmail)}
	}
	if message.HTMLBody != nil && strings.TrimSpace(*message.HTMLBody) != "" {
		body["content"] = append(body["content"].([]map[string]string), map[string]string{"type": "text/html", "value": *message.HTMLBody})
	}
	if message.IdempotencyKey != nil {
		body["custom_args"] = map[string]string{"idempotency_key": *message.IdempotencyKey}
	}
	if settings.SendGridSandboxMode {
		body["mail_settings"] = map[string]any{"sandbox_mode": map[string]bool{"enable": true}}
	}
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.sendgrid.com/v3/mail/send", bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+strings.TrimSpace(*settings.SendGridAPIKey))
	req.Header.Set("Content-Type", "application/json")
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("sendgrid delivery failed: status=%d body=%s", resp.StatusCode, strings.TrimSpace(string(bodyBytes)))
	}
	msgID := resp.Header.Get("X-Message-Id")
	if msgID == "" {
		msgID = uuid.NewString()
	}
	return &ports.EmailDeliveryResult{Provider: settings.Provider, Status: domain.NotifStatusSent, ExternalReference: msgID, MessageID: msgID, EventStatus: "accepted"}, nil
}

func buildMIMEMessage(from mail.Address, replyTo *string, to string, subject string, text string, html *string) ([]byte, error) {
	var buf bytes.Buffer
	writeHeader := func(key, value string) {
		if strings.TrimSpace(value) != "" {
			buf.WriteString(key)
			buf.WriteString(": ")
			buf.WriteString(value)
			buf.WriteString("\r\n")
		}
	}
	writeHeader("From", from.String())
	writeHeader("To", to)
	writeHeader("Subject", mime.QEncoding.Encode("utf-8", subject))
	writeHeader("MIME-Version", "1.0")
	if replyTo != nil {
		writeHeader("Reply-To", strings.TrimSpace(*replyTo))
	}
	if html == nil || strings.TrimSpace(*html) == "" {
		writeHeader("Content-Type", `text/plain; charset="utf-8"`)
		writeHeader("Content-Transfer-Encoding", "base64")
		buf.WriteString("\r\n")
		buf.WriteString(base64.StdEncoding.EncodeToString([]byte(text)))
		return buf.Bytes(), nil
	}
	boundary := "hrms-" + uuid.NewString()
	writeHeader("Content-Type", `multipart/alternative; boundary="`+boundary+`"`)
	buf.WriteString("\r\n--" + boundary + "\r\n")
	buf.WriteString("Content-Type: text/plain; charset=\"utf-8\"\r\nContent-Transfer-Encoding: base64\r\n\r\n")
	buf.WriteString(base64.StdEncoding.EncodeToString([]byte(text)))
	buf.WriteString("\r\n--" + boundary + "\r\n")
	buf.WriteString("Content-Type: text/html; charset=\"utf-8\"\r\nContent-Transfer-Encoding: base64\r\n\r\n")
	buf.WriteString(base64.StdEncoding.EncodeToString([]byte(*html)))
	buf.WriteString("\r\n--" + boundary + "--\r\n")
	return buf.Bytes(), nil
}

func valueOrEmpty(value *string) string {
	if value == nil {
		return ""
	}
	return strings.TrimSpace(*value)
}
