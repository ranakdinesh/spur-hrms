CREATE TABLE IF NOT EXISTS hrms.offer_letter_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    name VARCHAR(160) NOT NULL,
    description TEXT,
    subject VARCHAR(255),
    body_html TEXT NOT NULL,
    footer_html TEXT,
    locale VARCHAR(20) NOT NULL DEFAULT 'en-IN',
    is_default BOOLEAN NOT NULL DEFAULT FALSE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT offer_letter_templates_name_check CHECK (BTRIM(name) <> ''),
    CONSTRAINT offer_letter_templates_body_check CHECK (BTRIM(body_html) <> '')
);

CREATE UNIQUE INDEX IF NOT EXISTS offer_letter_templates_default_idx
    ON hrms.offer_letter_templates(tenant_id)
    WHERE is_default AND is_active AND NOT inactive;

CREATE INDEX IF NOT EXISTS offer_letter_templates_tenant_idx
    ON hrms.offer_letter_templates(tenant_id, name)
    WHERE NOT inactive;

ALTER TABLE hrms.offer_letters
    ADD COLUMN IF NOT EXISTS template_id UUID REFERENCES hrms.offer_letter_templates(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS subject VARCHAR(255),
    ADD COLUMN IF NOT EXISTS rendered_html TEXT,
    ADD COLUMN IF NOT EXISTS sent_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS revoked_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS signature_token TEXT,
    ADD COLUMN IF NOT EXISTS signature_requested_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS signature_completed_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS signer_name VARCHAR(255),
    ADD COLUMN IF NOT EXISTS signer_email VARCHAR(255),
    ADD COLUMN IF NOT EXISTS signer_ip INET,
    ADD COLUMN IF NOT EXISTS signer_user_agent TEXT,
    ADD COLUMN IF NOT EXISTS signature_hash TEXT,
    ADD COLUMN IF NOT EXISTS audit_certificate_url TEXT;

CREATE UNIQUE INDEX IF NOT EXISTS offer_letters_signature_token_idx
    ON hrms.offer_letters(signature_token)
    WHERE signature_token IS NOT NULL AND NOT inactive;

CREATE INDEX IF NOT EXISTS offer_letters_status_idx
    ON hrms.offer_letters(tenant_id, status)
    WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.offer_letter_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    offer_letter_id UUID NOT NULL REFERENCES hrms.offer_letters(id) ON DELETE CASCADE,
    from_status VARCHAR(50),
    to_status VARCHAR(50) NOT NULL,
    action VARCHAR(80) NOT NULL,
    remarks TEXT,
    actor_email VARCHAR(255),
    ip_address INET,
    user_agent TEXT,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by UUID REFERENCES auth.users(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS offer_letter_events_offer_idx
    ON hrms.offer_letter_events(tenant_id, offer_letter_id, created_at DESC)
    WHERE NOT inactive;

ALTER TABLE hrms.offer_letter_templates ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.offer_letter_events ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS tenant_isolation_offer_letter_templates ON hrms.offer_letter_templates;
CREATE POLICY tenant_isolation_offer_letter_templates ON hrms.offer_letter_templates
USING (
    current_setting('app.is_super_admin', true) = 'true'
    OR tenant_id::text = current_setting('app.tenant_id', true)
)
WITH CHECK (
    current_setting('app.is_super_admin', true) = 'true'
    OR tenant_id::text = current_setting('app.tenant_id', true)
);

DROP POLICY IF EXISTS tenant_isolation_offer_letter_events ON hrms.offer_letter_events;
CREATE POLICY tenant_isolation_offer_letter_events ON hrms.offer_letter_events
USING (
    current_setting('app.is_super_admin', true) = 'true'
    OR tenant_id::text = current_setting('app.tenant_id', true)
)
WITH CHECK (
    current_setting('app.is_super_admin', true) = 'true'
    OR tenant_id::text = current_setting('app.tenant_id', true)
);
