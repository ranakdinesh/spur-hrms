CREATE TABLE IF NOT EXISTS hrms.privacy_consents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    employee_user_id UUID,
    worker_profile_id UUID REFERENCES hrms.worker_profiles(id) ON DELETE SET NULL,
    consent_key TEXT NOT NULL,
    consent_area TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'granted',
    lawful_basis TEXT NOT NULL DEFAULT 'consent',
    channel TEXT NOT NULL DEFAULT 'web',
    source TEXT NOT NULL DEFAULT 'hrms',
    purpose TEXT NOT NULL,
    granted_at TIMESTAMPTZ,
    revoked_at TIMESTAMPTZ,
    expires_at TIMESTAMPTZ,
    evidence JSONB NOT NULL DEFAULT '{}'::jsonb,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by UUID,
    CHECK (status IN ('draft', 'granted', 'revoked', 'expired')),
    CHECK (lawful_basis IN ('consent', 'contract', 'legal_obligation', 'legitimate_interest', 'vital_interest')),
    CHECK (jsonb_typeof(evidence) = 'object'),
    CHECK (jsonb_typeof(metadata) = 'object')
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_hrms_privacy_consents_key
    ON hrms.privacy_consents (tenant_id, consent_key)
    WHERE NOT inactive;

CREATE INDEX IF NOT EXISTS idx_hrms_privacy_consents_status
    ON hrms.privacy_consents (tenant_id, status, consent_area)
    WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.data_erasure_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    request_key TEXT NOT NULL,
    subject_user_id UUID,
    worker_profile_id UUID REFERENCES hrms.worker_profiles(id) ON DELETE SET NULL,
    request_type TEXT NOT NULL DEFAULT 'erasure',
    status TEXT NOT NULL DEFAULT 'intake',
    priority TEXT NOT NULL DEFAULT 'normal',
    requested_by UUID,
    reason TEXT NOT NULL,
    scope JSONB NOT NULL DEFAULT '{}'::jsonb,
    retained_reason TEXT,
    due_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    reviewed_by UUID,
    reviewed_at TIMESTAMPTZ,
    audit_summary JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by UUID,
    CHECK (request_type IN ('erasure', 'export', 'correction', 'restriction')),
    CHECK (status IN ('intake', 'validating', 'blocked_legal_hold', 'approved', 'processing', 'completed', 'rejected', 'cancelled')),
    CHECK (priority IN ('low', 'normal', 'high', 'urgent')),
    CHECK (jsonb_typeof(scope) = 'object'),
    CHECK (jsonb_typeof(audit_summary) = 'object')
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_hrms_data_erasure_requests_key
    ON hrms.data_erasure_requests (tenant_id, request_key)
    WHERE NOT inactive;

CREATE INDEX IF NOT EXISTS idx_hrms_data_erasure_requests_status
    ON hrms.data_erasure_requests (tenant_id, status, due_at)
    WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.ecosystem_integration_hooks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    hook_key TEXT NOT NULL,
    provider TEXT NOT NULL,
    channel TEXT NOT NULL,
    direction TEXT NOT NULL DEFAULT 'outbound',
    status TEXT NOT NULL DEFAULT 'draft',
    display_name TEXT NOT NULL,
    endpoint_url TEXT,
    event_types TEXT[] NOT NULL DEFAULT ARRAY[]::TEXT[],
    secret_ref TEXT,
    consent_required BOOLEAN NOT NULL DEFAULT TRUE,
    mobile_safe BOOLEAN NOT NULL DEFAULT FALSE,
    last_checked_at TIMESTAMPTZ,
    last_error TEXT,
    config JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by UUID,
    CHECK (channel IN ('whatsapp', 'slack', 'email', 'git', 'webhook', 'mobile', 'api')),
    CHECK (direction IN ('inbound', 'outbound', 'bidirectional')),
    CHECK (status IN ('draft', 'active', 'paused', 'failed', 'archived')),
    CHECK (jsonb_typeof(config) = 'object')
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_hrms_ecosystem_integration_hooks_key
    ON hrms.ecosystem_integration_hooks (tenant_id, hook_key)
    WHERE NOT inactive;

CREATE INDEX IF NOT EXISTS idx_hrms_ecosystem_integration_hooks_status
    ON hrms.ecosystem_integration_hooks (tenant_id, channel, status)
    WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.mobile_api_constraints (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    constraint_key TEXT NOT NULL,
    workflow TEXT NOT NULL,
    min_android_version TEXT,
    min_ios_version TEXT,
    offline_supported BOOLEAN NOT NULL DEFAULT FALSE,
    low_bandwidth_mode BOOLEAN NOT NULL DEFAULT TRUE,
    requires_location BOOLEAN NOT NULL DEFAULT FALSE,
    requires_device_binding BOOLEAN NOT NULL DEFAULT FALSE,
    max_payload_kb INT NOT NULL DEFAULT 256,
    status TEXT NOT NULL DEFAULT 'active',
    notes TEXT,
    config JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by UUID,
    CHECK (status IN ('active', 'paused', 'archived')),
    CHECK (max_payload_kb > 0),
    CHECK (jsonb_typeof(config) = 'object')
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_hrms_mobile_api_constraints_key
    ON hrms.mobile_api_constraints (tenant_id, constraint_key)
    WHERE NOT inactive;

CREATE INDEX IF NOT EXISTS idx_hrms_mobile_api_constraints_workflow
    ON hrms.mobile_api_constraints (tenant_id, workflow, status)
    WHERE NOT inactive;

DROP TRIGGER IF EXISTS set_privacy_consents_updated_at ON hrms.privacy_consents;
CREATE TRIGGER set_privacy_consents_updated_at
    BEFORE UPDATE ON hrms.privacy_consents
    FOR EACH ROW EXECUTE FUNCTION hrms.set_updated_at();

DROP TRIGGER IF EXISTS set_data_erasure_requests_updated_at ON hrms.data_erasure_requests;
CREATE TRIGGER set_data_erasure_requests_updated_at
    BEFORE UPDATE ON hrms.data_erasure_requests
    FOR EACH ROW EXECUTE FUNCTION hrms.set_updated_at();

DROP TRIGGER IF EXISTS set_ecosystem_integration_hooks_updated_at ON hrms.ecosystem_integration_hooks;
CREATE TRIGGER set_ecosystem_integration_hooks_updated_at
    BEFORE UPDATE ON hrms.ecosystem_integration_hooks
    FOR EACH ROW EXECUTE FUNCTION hrms.set_updated_at();

DROP TRIGGER IF EXISTS set_mobile_api_constraints_updated_at ON hrms.mobile_api_constraints;
CREATE TRIGGER set_mobile_api_constraints_updated_at
    BEFORE UPDATE ON hrms.mobile_api_constraints
    FOR EACH ROW EXECUTE FUNCTION hrms.set_updated_at();

ALTER TABLE hrms.privacy_consents ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.data_erasure_requests ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.ecosystem_integration_hooks ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.mobile_api_constraints ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS privacy_consents_tenant_isolation ON hrms.privacy_consents;
CREATE POLICY privacy_consents_tenant_isolation ON hrms.privacy_consents
    USING (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true))
    WITH CHECK (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true));

DROP POLICY IF EXISTS data_erasure_requests_tenant_isolation ON hrms.data_erasure_requests;
CREATE POLICY data_erasure_requests_tenant_isolation ON hrms.data_erasure_requests
    USING (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true))
    WITH CHECK (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true));

DROP POLICY IF EXISTS ecosystem_integration_hooks_tenant_isolation ON hrms.ecosystem_integration_hooks;
CREATE POLICY ecosystem_integration_hooks_tenant_isolation ON hrms.ecosystem_integration_hooks
    USING (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true))
    WITH CHECK (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true));

DROP POLICY IF EXISTS mobile_api_constraints_tenant_isolation ON hrms.mobile_api_constraints;
CREATE POLICY mobile_api_constraints_tenant_isolation ON hrms.mobile_api_constraints
    USING (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true))
    WITH CHECK (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true));
