-- name: GetPushProviderSettings :one
SELECT * FROM hrms.push_provider_settings
WHERE tenant_id = $1 AND NOT inactive
LIMIT 1;

-- name: UpsertPushProviderSettings :one
INSERT INTO hrms.push_provider_settings (
    tenant_id, provider, is_enabled, project_id, client_email, private_key, private_key_id,
    auth_uri, token_uri, android_enabled, ios_enabled, web_enabled, default_click_action,
    default_image_url, ttl_seconds, collapse_key, metadata, created_by
)
VALUES (
    $1, $2, $3, $4, $5, $6, $7,
    $8, $9, $10, $11, $12, $13,
    $14, $15, $16, COALESCE($17, '{}'::jsonb), $18
)
ON CONFLICT (tenant_id) WHERE NOT inactive
DO UPDATE SET
    provider = EXCLUDED.provider,
    is_enabled = EXCLUDED.is_enabled,
    project_id = EXCLUDED.project_id,
    client_email = EXCLUDED.client_email,
    private_key = COALESCE(EXCLUDED.private_key, hrms.push_provider_settings.private_key),
    private_key_id = EXCLUDED.private_key_id,
    auth_uri = EXCLUDED.auth_uri,
    token_uri = EXCLUDED.token_uri,
    android_enabled = EXCLUDED.android_enabled,
    ios_enabled = EXCLUDED.ios_enabled,
    web_enabled = EXCLUDED.web_enabled,
    default_click_action = EXCLUDED.default_click_action,
    default_image_url = EXCLUDED.default_image_url,
    ttl_seconds = EXCLUDED.ttl_seconds,
    collapse_key = EXCLUDED.collapse_key,
    metadata = EXCLUDED.metadata,
    updated_at = NOW(),
    updated_by = $18
RETURNING *;

-- name: UpdatePushProviderTestResult :one
UPDATE hrms.push_provider_settings
SET last_test_at = NOW(),
    last_test_status = $3,
    last_test_message = $4,
    updated_at = NOW(),
    updated_by = $5
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: SoftDeletePushProviderSettings :exec
UPDATE hrms.push_provider_settings
SET inactive = TRUE,
    is_enabled = FALSE,
    updated_at = NOW(),
    updated_by = $3
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;
