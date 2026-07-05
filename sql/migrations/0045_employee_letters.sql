CREATE TABLE IF NOT EXISTS hrms.employee_letter_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    letter_type VARCHAR(40) NOT NULL,
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
    CONSTRAINT employee_letter_templates_type_check CHECK (letter_type IN ('appointment', 'experience', 'relieving')),
    CONSTRAINT employee_letter_templates_name_check CHECK (BTRIM(name) <> ''),
    CONSTRAINT employee_letter_templates_body_check CHECK (BTRIM(body_html) <> '')
);

CREATE UNIQUE INDEX IF NOT EXISTS employee_letter_templates_default_idx
    ON hrms.employee_letter_templates(tenant_id, letter_type)
    WHERE is_default AND is_active AND NOT inactive;

CREATE INDEX IF NOT EXISTS employee_letter_templates_tenant_idx
    ON hrms.employee_letter_templates(tenant_id, letter_type, name)
    WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.employee_letters (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    employee_id UUID NOT NULL REFERENCES hrms.employees(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    template_id UUID REFERENCES hrms.employee_letter_templates(id) ON DELETE SET NULL,
    document_type_id UUID REFERENCES hrms.document_types(id) ON DELETE SET NULL,
    employee_document_id UUID REFERENCES hrms.employee_documents(id) ON DELETE SET NULL,
    letter_type VARCHAR(40) NOT NULL,
    subject VARCHAR(255),
    rendered_html TEXT,
    status VARCHAR(40) NOT NULL DEFAULT 'Generated',
    issue_date DATE,
    effective_date DATE,
    end_date DATE,
    pdf_path TEXT,
    version INTEGER NOT NULL DEFAULT 1,
    is_latest BOOLEAN NOT NULL DEFAULT TRUE,
    approval_requested_at TIMESTAMPTZ,
    approved_at TIMESTAMPTZ,
    approved_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    sent_at TIMESTAMPTZ,
    revoked_at TIMESTAMPTZ,
    signature_token TEXT,
    signature_requested_at TIMESTAMPTZ,
    signature_completed_at TIMESTAMPTZ,
    signer_name VARCHAR(255),
    signer_email VARCHAR(255),
    signer_ip INET,
    signer_user_agent TEXT,
    signature_hash TEXT,
    audit_certificate_url TEXT,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT employee_letters_type_check CHECK (letter_type IN ('appointment', 'experience', 'relieving')),
    CONSTRAINT employee_letters_status_check CHECK (status IN ('Generated', 'Approved', 'Sent', 'Signed', 'Revoked'))
);

CREATE INDEX IF NOT EXISTS employee_letters_employee_idx
    ON hrms.employee_letters(tenant_id, employee_id, letter_type, created_at DESC)
    WHERE NOT inactive;

CREATE UNIQUE INDEX IF NOT EXISTS employee_letters_latest_idx
    ON hrms.employee_letters(tenant_id, employee_id, letter_type)
    WHERE is_latest AND NOT inactive;

CREATE UNIQUE INDEX IF NOT EXISTS employee_letters_signature_token_idx
    ON hrms.employee_letters(signature_token)
    WHERE signature_token IS NOT NULL AND NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.employee_letter_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    employee_letter_id UUID NOT NULL REFERENCES hrms.employee_letters(id) ON DELETE CASCADE,
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

CREATE INDEX IF NOT EXISTS employee_letter_events_letter_idx
    ON hrms.employee_letter_events(tenant_id, employee_letter_id, created_at DESC)
    WHERE NOT inactive;

ALTER TABLE hrms.employee_letter_templates ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.employee_letters ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.employee_letter_events ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS tenant_isolation_employee_letter_templates ON hrms.employee_letter_templates;
CREATE POLICY tenant_isolation_employee_letter_templates ON hrms.employee_letter_templates
USING (
    current_setting('app.is_super_admin', true) = 'true'
    OR tenant_id::text = current_setting('app.tenant_id', true)
)
WITH CHECK (
    current_setting('app.is_super_admin', true) = 'true'
    OR tenant_id::text = current_setting('app.tenant_id', true)
);

DROP POLICY IF EXISTS tenant_isolation_employee_letters ON hrms.employee_letters;
CREATE POLICY tenant_isolation_employee_letters ON hrms.employee_letters
USING (
    current_setting('app.is_super_admin', true) = 'true'
    OR tenant_id::text = current_setting('app.tenant_id', true)
)
WITH CHECK (
    current_setting('app.is_super_admin', true) = 'true'
    OR tenant_id::text = current_setting('app.tenant_id', true)
);

DROP POLICY IF EXISTS tenant_isolation_employee_letter_events ON hrms.employee_letter_events;
CREATE POLICY tenant_isolation_employee_letter_events ON hrms.employee_letter_events
USING (
    current_setting('app.is_super_admin', true) = 'true'
    OR tenant_id::text = current_setting('app.tenant_id', true)
)
WITH CHECK (
    current_setting('app.is_super_admin', true) = 'true'
    OR tenant_id::text = current_setting('app.tenant_id', true)
);
