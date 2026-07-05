-- name: GetStorageProviderSettings :one
SELECT * FROM hrms.storage_provider_settings
WHERE tenant_id = $1 AND NOT inactive
LIMIT 1;

-- name: UpsertStorageProviderSettings :one
INSERT INTO hrms.storage_provider_settings (
    tenant_id, provider, is_enabled, bucket, region, endpoint, access_key_id, secret_access_key,
    use_ssl, force_path_style, object_prefix, public_base_url, max_file_size_bytes,
    allowed_content_types, metadata, created_by
)
VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8,
    $9, $10, $11, $12, $13,
    $14, COALESCE($15, '{}'::jsonb), $16
)
ON CONFLICT (tenant_id) WHERE NOT inactive
DO UPDATE SET
    provider = EXCLUDED.provider,
    is_enabled = EXCLUDED.is_enabled,
    bucket = EXCLUDED.bucket,
    region = EXCLUDED.region,
    endpoint = EXCLUDED.endpoint,
    access_key_id = COALESCE(EXCLUDED.access_key_id, hrms.storage_provider_settings.access_key_id),
    secret_access_key = COALESCE(EXCLUDED.secret_access_key, hrms.storage_provider_settings.secret_access_key),
    use_ssl = EXCLUDED.use_ssl,
    force_path_style = EXCLUDED.force_path_style,
    object_prefix = EXCLUDED.object_prefix,
    public_base_url = EXCLUDED.public_base_url,
    max_file_size_bytes = EXCLUDED.max_file_size_bytes,
    allowed_content_types = EXCLUDED.allowed_content_types,
    metadata = EXCLUDED.metadata,
    updated_at = NOW(),
    updated_by = $16
RETURNING *;

-- name: UpdateStorageProviderTestResult :one
UPDATE hrms.storage_provider_settings
SET last_test_at = NOW(),
    last_test_status = $3,
    last_test_message = $4,
    updated_at = NOW(),
    updated_by = $5
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: SoftDeleteStorageProviderSettings :exec
UPDATE hrms.storage_provider_settings
SET inactive = TRUE,
    is_enabled = FALSE,
    updated_at = NOW(),
    updated_by = $3
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;
