-- name: GetCommunicationProviderSettings :one
SELECT * FROM hrms.communication_provider_settings
WHERE tenant_id = $1 AND NOT inactive
LIMIT 1;

-- name: UpsertCommunicationProviderSettings :one
INSERT INTO hrms.communication_provider_settings (
    tenant_id,
    sms_provider, sms_enabled, sms_sender_id, sms_auth_key, sms_template_id, sms_route, sms_country_code, sms_base_url,
    whatsapp_provider, whatsapp_enabled, whatsapp_auth_key, whatsapp_app_name, whatsapp_source_number, whatsapp_template_id, whatsapp_template_name, whatsapp_namespace, whatsapp_base_url,
    webhook_signing_secret, metadata, created_by
)
VALUES (
    $1,
    $2, $3, $4, $5, $6, $7, $8, $9,
    $10, $11, $12, $13, $14, $15, $16, $17, $18,
    $19, COALESCE($20, '{}'::jsonb), $21
)
ON CONFLICT (tenant_id) WHERE NOT inactive
DO UPDATE SET
    sms_provider = EXCLUDED.sms_provider,
    sms_enabled = EXCLUDED.sms_enabled,
    sms_sender_id = EXCLUDED.sms_sender_id,
    sms_auth_key = COALESCE(EXCLUDED.sms_auth_key, hrms.communication_provider_settings.sms_auth_key),
    sms_template_id = EXCLUDED.sms_template_id,
    sms_route = EXCLUDED.sms_route,
    sms_country_code = EXCLUDED.sms_country_code,
    sms_base_url = EXCLUDED.sms_base_url,
    whatsapp_provider = EXCLUDED.whatsapp_provider,
    whatsapp_enabled = EXCLUDED.whatsapp_enabled,
    whatsapp_auth_key = COALESCE(EXCLUDED.whatsapp_auth_key, hrms.communication_provider_settings.whatsapp_auth_key),
    whatsapp_app_name = EXCLUDED.whatsapp_app_name,
    whatsapp_source_number = EXCLUDED.whatsapp_source_number,
    whatsapp_template_id = EXCLUDED.whatsapp_template_id,
    whatsapp_template_name = EXCLUDED.whatsapp_template_name,
    whatsapp_namespace = EXCLUDED.whatsapp_namespace,
    whatsapp_base_url = EXCLUDED.whatsapp_base_url,
    webhook_signing_secret = COALESCE(EXCLUDED.webhook_signing_secret, hrms.communication_provider_settings.webhook_signing_secret),
    metadata = EXCLUDED.metadata,
    updated_at = NOW(),
    updated_by = $21
RETURNING *;

-- name: UpdateCommunicationProviderTestResult :one
UPDATE hrms.communication_provider_settings
SET sms_last_test_at = CASE WHEN $3 = 'sms' THEN NOW() ELSE sms_last_test_at END,
    sms_last_test_status = CASE WHEN $3 = 'sms' THEN $4 ELSE sms_last_test_status END,
    sms_last_test_message = CASE WHEN $3 = 'sms' THEN $5 ELSE sms_last_test_message END,
    whatsapp_last_test_at = CASE WHEN $3 = 'whatsapp' THEN NOW() ELSE whatsapp_last_test_at END,
    whatsapp_last_test_status = CASE WHEN $3 = 'whatsapp' THEN $4 ELSE whatsapp_last_test_status END,
    whatsapp_last_test_message = CASE WHEN $3 = 'whatsapp' THEN $5 ELSE whatsapp_last_test_message END,
    updated_at = NOW(),
    updated_by = $6
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: SoftDeleteCommunicationProviderSettings :exec
UPDATE hrms.communication_provider_settings
SET inactive = TRUE,
    sms_enabled = FALSE,
    whatsapp_enabled = FALSE,
    updated_at = NOW(),
    updated_by = $3
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;
