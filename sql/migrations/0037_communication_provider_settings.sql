CREATE TABLE IF NOT EXISTS hrms.communication_provider_settings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    sms_provider TEXT NOT NULL DEFAULT 'local',
    sms_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    sms_sender_id TEXT NULL,
    sms_auth_key TEXT NULL,
    sms_template_id TEXT NULL,
    sms_route TEXT NULL,
    sms_country_code TEXT NULL,
    sms_base_url TEXT NULL,
    whatsapp_provider TEXT NOT NULL DEFAULT 'local',
    whatsapp_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    whatsapp_auth_key TEXT NULL,
    whatsapp_app_name TEXT NULL,
    whatsapp_source_number TEXT NULL,
    whatsapp_template_id TEXT NULL,
    whatsapp_template_name TEXT NULL,
    whatsapp_namespace TEXT NULL,
    whatsapp_base_url TEXT NULL,
    webhook_signing_secret TEXT NULL,
    sms_last_test_at TIMESTAMPTZ NULL,
    sms_last_test_status TEXT NULL,
    sms_last_test_message TEXT NULL,
    whatsapp_last_test_at TIMESTAMPTZ NULL,
    whatsapp_last_test_status TEXT NULL,
    whatsapp_last_test_message TEXT NULL,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by UUID NULL,
    CONSTRAINT communication_provider_settings_sms_provider_check CHECK (sms_provider IN ('local','msg91')),
    CONSTRAINT communication_provider_settings_whatsapp_provider_check CHECK (whatsapp_provider IN ('local','msg91','gupshup')),
    CONSTRAINT communication_provider_settings_sms_test_status_check CHECK (sms_last_test_status IS NULL OR sms_last_test_status IN ('Pending','Sent','Failed','Suppressed')),
    CONSTRAINT communication_provider_settings_whatsapp_test_status_check CHECK (whatsapp_last_test_status IS NULL OR whatsapp_last_test_status IN ('Pending','Sent','Failed','Suppressed'))
);

CREATE UNIQUE INDEX IF NOT EXISTS communication_provider_settings_tenant_active_idx
    ON hrms.communication_provider_settings(tenant_id)
    WHERE inactive = FALSE;

CREATE INDEX IF NOT EXISTS communication_provider_settings_provider_idx
    ON hrms.communication_provider_settings(tenant_id, sms_provider, whatsapp_provider)
    WHERE inactive = FALSE;

ALTER TABLE hrms.communication_provider_settings ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS tenant_isolation_communication_provider_settings ON hrms.communication_provider_settings;
CREATE POLICY tenant_isolation_communication_provider_settings ON hrms.communication_provider_settings
    USING (tenant_id = current_setting('app.tenant_id', true)::uuid OR current_setting('app.is_super_admin', true) = 'true')
    WITH CHECK (tenant_id = current_setting('app.tenant_id', true)::uuid OR current_setting('app.is_super_admin', true) = 'true');
