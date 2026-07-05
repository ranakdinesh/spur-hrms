CREATE TABLE IF NOT EXISTS hrms.push_provider_settings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    provider TEXT NOT NULL DEFAULT 'local',
    is_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    project_id TEXT NULL,
    client_email TEXT NULL,
    private_key TEXT NULL,
    private_key_id TEXT NULL,
    auth_uri TEXT NULL,
    token_uri TEXT NULL,
    android_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    ios_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    web_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    default_click_action TEXT NULL,
    default_image_url TEXT NULL,
    ttl_seconds INTEGER NOT NULL DEFAULT 3600,
    collapse_key TEXT NULL,
    last_test_at TIMESTAMPTZ NULL,
    last_test_status TEXT NULL,
    last_test_message TEXT NULL,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by UUID NULL,
    CONSTRAINT push_provider_settings_provider_check CHECK (provider IN ('local','fcm')),
    CONSTRAINT push_provider_settings_ttl_check CHECK (ttl_seconds > 0 AND ttl_seconds <= 2419200),
    CONSTRAINT push_provider_settings_test_status_check CHECK (last_test_status IS NULL OR last_test_status IN ('Pending','Sent','Failed','Suppressed'))
);

CREATE UNIQUE INDEX IF NOT EXISTS push_provider_settings_tenant_active_idx
    ON hrms.push_provider_settings(tenant_id)
    WHERE inactive = FALSE;

ALTER TABLE hrms.push_provider_settings ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS tenant_isolation_push_provider_settings ON hrms.push_provider_settings;
CREATE POLICY tenant_isolation_push_provider_settings ON hrms.push_provider_settings
    USING (tenant_id = current_setting('app.tenant_id', true)::uuid OR current_setting('app.is_super_admin', true) = 'true')
    WITH CHECK (tenant_id = current_setting('app.tenant_id', true)::uuid OR current_setting('app.is_super_admin', true) = 'true');
