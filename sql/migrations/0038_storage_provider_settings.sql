CREATE TABLE IF NOT EXISTS hrms.storage_provider_settings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    provider TEXT NOT NULL DEFAULT 'minio',
    is_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    bucket TEXT NOT NULL,
    region TEXT NULL,
    endpoint TEXT NULL,
    access_key_id TEXT NULL,
    secret_access_key TEXT NULL,
    use_ssl BOOLEAN NOT NULL DEFAULT TRUE,
    force_path_style BOOLEAN NOT NULL DEFAULT TRUE,
    object_prefix TEXT NULL,
    public_base_url TEXT NULL,
    max_file_size_bytes BIGINT NOT NULL DEFAULT 10485760,
    allowed_content_types TEXT NULL,
    last_test_at TIMESTAMPTZ NULL,
    last_test_status TEXT NULL,
    last_test_message TEXT NULL,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by UUID NULL,
    CONSTRAINT storage_provider_settings_provider_check CHECK (provider IN ('s3','minio')),
    CONSTRAINT storage_provider_settings_bucket_check CHECK (length(trim(bucket)) > 0),
    CONSTRAINT storage_provider_settings_max_size_check CHECK (max_file_size_bytes > 0),
    CONSTRAINT storage_provider_settings_test_status_check CHECK (last_test_status IS NULL OR last_test_status IN ('Pending','Sent','Failed','Suppressed'))
);

CREATE UNIQUE INDEX IF NOT EXISTS storage_provider_settings_tenant_active_idx
    ON hrms.storage_provider_settings(tenant_id)
    WHERE inactive = FALSE;

CREATE INDEX IF NOT EXISTS storage_provider_settings_provider_idx
    ON hrms.storage_provider_settings(tenant_id, provider)
    WHERE inactive = FALSE;

ALTER TABLE hrms.storage_provider_settings ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS tenant_isolation_storage_provider_settings ON hrms.storage_provider_settings;
CREATE POLICY tenant_isolation_storage_provider_settings ON hrms.storage_provider_settings
    USING (tenant_id = current_setting('app.tenant_id', true)::uuid OR current_setting('app.is_super_admin', true) = 'true')
    WITH CHECK (tenant_id = current_setting('app.tenant_id', true)::uuid OR current_setting('app.is_super_admin', true) = 'true');
