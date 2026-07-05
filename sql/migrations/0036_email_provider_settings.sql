CREATE TABLE IF NOT EXISTS hrms.email_provider_settings (
    id                     UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id              UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    provider               VARCHAR(40) NOT NULL DEFAULT 'smtp',
    is_enabled             BOOLEAN NOT NULL DEFAULT FALSE,
    from_name              VARCHAR(160),
    from_email             VARCHAR(255) NOT NULL,
    reply_to_email         VARCHAR(255),
    smtp_host              VARCHAR(255),
    smtp_port              INT,
    smtp_username          VARCHAR(255),
    smtp_password          TEXT,
    smtp_encryption        VARCHAR(20) NOT NULL DEFAULT 'starttls',
    sendgrid_api_key       TEXT,
    sendgrid_sandbox_mode  BOOLEAN NOT NULL DEFAULT FALSE,
    webhook_signing_secret TEXT,
    spf_status             VARCHAR(40),
    dkim_status            VARCHAR(40),
    dmarc_status           VARCHAR(40),
    last_test_at           TIMESTAMPTZ,
    last_test_status       VARCHAR(40),
    last_test_message      TEXT,
    metadata               JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive               BOOLEAN NOT NULL DEFAULT FALSE,
    created_at             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by             UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by             UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT email_provider_settings_provider_check CHECK (provider IN ('local','smtp','sendgrid')),
    CONSTRAINT email_provider_settings_smtp_port_check CHECK (smtp_port IS NULL OR (smtp_port > 0 AND smtp_port <= 65535)),
    CONSTRAINT email_provider_settings_encryption_check CHECK (smtp_encryption IN ('none','starttls','tls')),
    CONSTRAINT email_provider_settings_test_status_check CHECK (last_test_status IS NULL OR last_test_status IN ('Pending','Sent','Failed','Suppressed'))
);

CREATE UNIQUE INDEX IF NOT EXISTS email_provider_settings_tenant_active_idx
    ON hrms.email_provider_settings(tenant_id)
    WHERE NOT inactive;

CREATE INDEX IF NOT EXISTS email_provider_settings_provider_idx
    ON hrms.email_provider_settings(tenant_id, provider)
    WHERE NOT inactive;

ALTER TABLE hrms.email_provider_settings ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS tenant_isolation_email_provider_settings ON hrms.email_provider_settings;
CREATE POLICY tenant_isolation_email_provider_settings ON hrms.email_provider_settings
    USING (
        current_setting('app.is_super_admin', true) = 'true'
        OR tenant_id::text = current_setting('app.tenant_id', true)
    )
    WITH CHECK (
        current_setting('app.is_super_admin', true) = 'true'
        OR tenant_id::text = current_setting('app.tenant_id', true)
    );

ALTER TABLE hrms.notification_logs
    ADD COLUMN IF NOT EXISTS provider VARCHAR(40),
    ADD COLUMN IF NOT EXISTS provider_message_id TEXT,
    ADD COLUMN IF NOT EXISTS provider_event_status VARCHAR(80),
    ADD COLUMN IF NOT EXISTS provider_event_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS attempt_count INT NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS last_attempt_at TIMESTAMPTZ;

CREATE INDEX IF NOT EXISTS notification_logs_provider_event_idx
    ON hrms.notification_logs(tenant_id, provider, provider_event_status)
    WHERE provider IS NOT NULL AND NOT inactive;
