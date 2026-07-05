-- name: GetEmailProviderSettings :one
SELECT * FROM hrms.email_provider_settings
WHERE tenant_id = $1 AND NOT inactive
LIMIT 1;

-- name: UpsertEmailProviderSettings :one
INSERT INTO hrms.email_provider_settings (
    tenant_id, provider, is_enabled, from_name, from_email, reply_to_email,
    smtp_host, smtp_port, smtp_username, smtp_password, smtp_encryption,
    sendgrid_api_key, sendgrid_sandbox_mode, webhook_signing_secret,
    spf_status, dkim_status, dmarc_status, metadata, created_by
)
VALUES (
    $1, $2, $3, $4, $5, $6,
    $7, $8, $9, $10, $11,
    $12, $13, $14,
    $15, $16, $17, COALESCE($18, '{}'::jsonb), $19
)
ON CONFLICT (tenant_id) WHERE NOT inactive
DO UPDATE SET
    provider = EXCLUDED.provider,
    is_enabled = EXCLUDED.is_enabled,
    from_name = EXCLUDED.from_name,
    from_email = EXCLUDED.from_email,
    reply_to_email = EXCLUDED.reply_to_email,
    smtp_host = EXCLUDED.smtp_host,
    smtp_port = EXCLUDED.smtp_port,
    smtp_username = EXCLUDED.smtp_username,
    smtp_password = COALESCE(EXCLUDED.smtp_password, hrms.email_provider_settings.smtp_password),
    smtp_encryption = EXCLUDED.smtp_encryption,
    sendgrid_api_key = COALESCE(EXCLUDED.sendgrid_api_key, hrms.email_provider_settings.sendgrid_api_key),
    sendgrid_sandbox_mode = EXCLUDED.sendgrid_sandbox_mode,
    webhook_signing_secret = COALESCE(EXCLUDED.webhook_signing_secret, hrms.email_provider_settings.webhook_signing_secret),
    spf_status = EXCLUDED.spf_status,
    dkim_status = EXCLUDED.dkim_status,
    dmarc_status = EXCLUDED.dmarc_status,
    metadata = EXCLUDED.metadata,
    updated_at = NOW(),
    updated_by = $19
RETURNING *;

-- name: UpdateEmailProviderTestResult :one
UPDATE hrms.email_provider_settings
SET last_test_at = NOW(),
    last_test_status = $3,
    last_test_message = $4,
    updated_at = NOW(),
    updated_by = $5
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: SoftDeleteEmailProviderSettings :exec
UPDATE hrms.email_provider_settings
SET inactive = TRUE, is_enabled = FALSE, updated_at = NOW(), updated_by = $3
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: UpdateNotificationLogDelivery :one
UPDATE hrms.notification_logs
SET status = $3,
    sent_date = CASE WHEN $3 = 'Sent' THEN COALESCE(sent_date, NOW()) ELSE sent_date END,
    error_message = $4,
    external_reference_id = $5,
    provider = $6,
    provider_message_id = $7,
    provider_event_status = $8,
    provider_event_at = COALESCE($9, provider_event_at),
    attempt_count = attempt_count + 1,
    last_attempt_at = NOW(),
    updated_at = NOW(),
    updated_by = $10
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;
