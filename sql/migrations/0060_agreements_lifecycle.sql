CREATE TABLE IF NOT EXISTS hrms.agreement_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    agreement_type TEXT NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    subject TEXT,
    body_html TEXT NOT NULL,
    footer_html TEXT,
    locale TEXT NOT NULL DEFAULT 'en-IN',
    is_default BOOLEAN NOT NULL DEFAULT FALSE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    CONSTRAINT agreement_templates_type_check CHECK (agreement_type IN ('sow', 'nda', 'retainer', 'freelance_contract', 'internship_letter', 'amendment')),
    CONSTRAINT agreement_templates_name_check CHECK (btrim(name) <> ''),
    CONSTRAINT agreement_templates_body_check CHECK (btrim(body_html) <> ''),
    CONSTRAINT agreement_templates_metadata_object_check CHECK (jsonb_typeof(metadata) = 'object')
);

CREATE UNIQUE INDEX IF NOT EXISTS agreement_templates_default_idx
    ON hrms.agreement_templates (tenant_id, agreement_type)
    WHERE is_default AND is_active AND inactive = FALSE;

CREATE INDEX IF NOT EXISTS agreement_templates_tenant_type_idx
    ON hrms.agreement_templates (tenant_id, agreement_type, name)
    WHERE inactive = FALSE;

DROP TRIGGER IF EXISTS trg_agreement_templates_updated_at ON hrms.agreement_templates;
CREATE TRIGGER trg_agreement_templates_updated_at
BEFORE UPDATE ON hrms.agreement_templates
FOR EACH ROW
EXECUTE FUNCTION hrms.set_updated_at();

ALTER TABLE hrms.agreement_templates ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS agreement_templates_tenant_isolation ON hrms.agreement_templates;
CREATE POLICY agreement_templates_tenant_isolation ON hrms.agreement_templates
USING (
    current_setting('app.is_super_admin', true) = 'true'
    OR tenant_id::text = current_setting('app.tenant_id', true)
)
WITH CHECK (
    current_setting('app.is_super_admin', true) = 'true'
    OR tenant_id::text = current_setting('app.tenant_id', true)
);

CREATE TABLE IF NOT EXISTS hrms.agreements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    agreement_type TEXT NOT NULL,
    title TEXT NOT NULL,
    template_id UUID REFERENCES hrms.agreement_templates(id) ON DELETE SET NULL,
    worker_profile_id UUID REFERENCES hrms.worker_profiles(id) ON DELETE SET NULL,
    engagement_id UUID REFERENCES hrms.engagements(id) ON DELETE SET NULL,
    project_id UUID REFERENCES hrms.projects(id) ON DELETE SET NULL,
    subject TEXT,
    rendered_html TEXT,
    status TEXT NOT NULL DEFAULT 'Generated',
    issue_date DATE,
    effective_date DATE,
    end_date DATE,
    pdf_path TEXT,
    version INTEGER NOT NULL DEFAULT 1,
    is_latest BOOLEAN NOT NULL DEFAULT TRUE,
    sent_at TIMESTAMPTZ,
    revoked_at TIMESTAMPTZ,
    signature_token TEXT,
    signature_requested_at TIMESTAMPTZ,
    signature_completed_at TIMESTAMPTZ,
    signer_name TEXT,
    signer_email TEXT,
    signer_ip INET,
    signer_user_agent TEXT,
    signature_hash TEXT,
    audit_certificate_url TEXT,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    CONSTRAINT agreements_type_check CHECK (agreement_type IN ('sow', 'nda', 'retainer', 'freelance_contract', 'internship_letter', 'amendment')),
    CONSTRAINT agreements_status_check CHECK (status IN ('Generated', 'Approved', 'Sent', 'Signed', 'Revoked')),
    CONSTRAINT agreements_title_check CHECK (btrim(title) <> ''),
    CONSTRAINT agreements_date_order_check CHECK (effective_date IS NULL OR end_date IS NULL OR end_date >= effective_date),
    CONSTRAINT agreements_metadata_object_check CHECK (jsonb_typeof(metadata) = 'object')
);

CREATE INDEX IF NOT EXISTS agreements_tenant_status_idx
    ON hrms.agreements (tenant_id, status, updated_at DESC)
    WHERE inactive = FALSE;

CREATE INDEX IF NOT EXISTS agreements_tenant_worker_idx
    ON hrms.agreements (tenant_id, worker_profile_id, created_at DESC)
    WHERE inactive = FALSE;

CREATE INDEX IF NOT EXISTS agreements_tenant_engagement_idx
    ON hrms.agreements (tenant_id, engagement_id, created_at DESC)
    WHERE inactive = FALSE;

CREATE INDEX IF NOT EXISTS agreements_tenant_project_idx
    ON hrms.agreements (tenant_id, project_id, created_at DESC)
    WHERE inactive = FALSE;

CREATE UNIQUE INDEX IF NOT EXISTS agreements_signature_token_idx
    ON hrms.agreements (signature_token)
    WHERE signature_token IS NOT NULL AND inactive = FALSE;

DROP TRIGGER IF EXISTS trg_agreements_updated_at ON hrms.agreements;
CREATE TRIGGER trg_agreements_updated_at
BEFORE UPDATE ON hrms.agreements
FOR EACH ROW
EXECUTE FUNCTION hrms.set_updated_at();

ALTER TABLE hrms.agreements ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS agreements_tenant_isolation ON hrms.agreements;
CREATE POLICY agreements_tenant_isolation ON hrms.agreements
USING (
    current_setting('app.is_super_admin', true) = 'true'
    OR tenant_id::text = current_setting('app.tenant_id', true)
)
WITH CHECK (
    current_setting('app.is_super_admin', true) = 'true'
    OR tenant_id::text = current_setting('app.tenant_id', true)
);

CREATE TABLE IF NOT EXISTS hrms.agreement_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    agreement_id UUID NOT NULL REFERENCES hrms.agreements(id) ON DELETE CASCADE,
    from_status TEXT,
    to_status TEXT NOT NULL,
    action TEXT NOT NULL,
    remarks TEXT,
    actor_email TEXT,
    ip_address INET,
    user_agent TEXT,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    CONSTRAINT agreement_events_metadata_object_check CHECK (jsonb_typeof(metadata) = 'object')
);

CREATE INDEX IF NOT EXISTS agreement_events_agreement_idx
    ON hrms.agreement_events (tenant_id, agreement_id, created_at DESC)
    WHERE inactive = FALSE;

DROP TRIGGER IF EXISTS trg_agreement_events_updated_at ON hrms.agreement_events;
CREATE TRIGGER trg_agreement_events_updated_at
BEFORE UPDATE ON hrms.agreement_events
FOR EACH ROW
EXECUTE FUNCTION hrms.set_updated_at();

ALTER TABLE hrms.agreement_events ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS agreement_events_tenant_isolation ON hrms.agreement_events;
CREATE POLICY agreement_events_tenant_isolation ON hrms.agreement_events
USING (
    current_setting('app.is_super_admin', true) = 'true'
    OR tenant_id::text = current_setting('app.tenant_id', true)
)
WITH CHECK (
    current_setting('app.is_super_admin', true) = 'true'
    OR tenant_id::text = current_setting('app.tenant_id', true)
);
