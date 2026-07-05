ALTER TABLE hrms.insights
    ADD COLUMN IF NOT EXISTS visibility_scope TEXT NOT NULL DEFAULT 'hr',
    ADD COLUMN IF NOT EXISTS human_review_required BOOLEAN NOT NULL DEFAULT TRUE,
    ADD COLUMN IF NOT EXISTS action_state TEXT NOT NULL DEFAULT 'proposed',
    ADD COLUMN IF NOT EXISTS employee_safe_summary TEXT;

CREATE TABLE IF NOT EXISTS hrms.ai_signal_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    signal_key TEXT NOT NULL,
    signal_type TEXT NOT NULL,
    source_module TEXT NOT NULL,
    source_event TEXT NOT NULL,
    severity TEXT NOT NULL DEFAULT 'medium',
    processing_status TEXT NOT NULL DEFAULT 'new',
    entity_type TEXT,
    entity_id UUID,
    employee_user_id UUID,
    visibility_scope TEXT NOT NULL DEFAULT 'hr',
    idempotency_key TEXT,
    correlation_id TEXT,
    payload JSONB NOT NULL DEFAULT '{}'::jsonb,
    explainability JSONB NOT NULL DEFAULT '{}'::jsonb,
    occurred_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    processed_at TIMESTAMPTZ,
    error_message TEXT,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by UUID,
    CHECK (severity IN ('low', 'medium', 'high', 'critical')),
    CHECK (processing_status IN ('new', 'queued', 'processed', 'ignored', 'failed')),
    CHECK (visibility_scope IN ('employee', 'manager_aggregate', 'hr', 'admin'))
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_hrms_ai_signal_logs_key
    ON hrms.ai_signal_logs (tenant_id, signal_key)
    WHERE NOT inactive;

CREATE UNIQUE INDEX IF NOT EXISTS uq_hrms_ai_signal_logs_idempotency
    ON hrms.ai_signal_logs (tenant_id, idempotency_key)
    WHERE idempotency_key IS NOT NULL AND NOT inactive;

CREATE INDEX IF NOT EXISTS idx_hrms_ai_signal_logs_tenant_status
    ON hrms.ai_signal_logs (tenant_id, processing_status, severity, occurred_at DESC)
    WHERE NOT inactive;

CREATE INDEX IF NOT EXISTS idx_hrms_ai_signal_logs_entity
    ON hrms.ai_signal_logs (tenant_id, entity_type, entity_id)
    WHERE entity_id IS NOT NULL AND NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.ai_agent_action_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    action_key TEXT NOT NULL,
    agent_key TEXT NOT NULL,
    agent_name TEXT NOT NULL,
    action_type TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'proposed',
    severity TEXT NOT NULL DEFAULT 'medium',
    title TEXT NOT NULL,
    summary TEXT NOT NULL,
    insight_id UUID REFERENCES hrms.insights(id),
    signal_id UUID REFERENCES hrms.ai_signal_logs(id),
    entity_type TEXT,
    entity_id UUID,
    employee_user_id UUID,
    visibility_scope TEXT NOT NULL DEFAULT 'hr',
    proposed_action JSONB NOT NULL DEFAULT '{}'::jsonb,
    input_snapshot JSONB NOT NULL DEFAULT '{}'::jsonb,
    output_snapshot JSONB NOT NULL DEFAULT '{}'::jsonb,
    explainability JSONB NOT NULL DEFAULT '{}'::jsonb,
    confidence_score NUMERIC(5,2) NOT NULL DEFAULT 0,
    model_version TEXT,
    sidecar_run_id TEXT,
    requires_human_review BOOLEAN NOT NULL DEFAULT TRUE,
    reviewed_by UUID,
    reviewed_at TIMESTAMPTZ,
    executed_at TIMESTAMPTZ,
    failed_at TIMESTAMPTZ,
    failure_message TEXT,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by UUID,
    CHECK (status IN ('proposed', 'queued', 'reviewing', 'approved', 'rejected', 'executed', 'failed', 'overridden', 'cancelled')),
    CHECK (severity IN ('low', 'medium', 'high', 'critical')),
    CHECK (visibility_scope IN ('employee', 'manager_aggregate', 'hr', 'admin'))
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_hrms_ai_agent_action_logs_key
    ON hrms.ai_agent_action_logs (tenant_id, action_key)
    WHERE NOT inactive;

CREATE INDEX IF NOT EXISTS idx_hrms_ai_agent_action_logs_tenant_status
    ON hrms.ai_agent_action_logs (tenant_id, status, severity, created_at DESC)
    WHERE NOT inactive;

CREATE INDEX IF NOT EXISTS idx_hrms_ai_agent_action_logs_insight
    ON hrms.ai_agent_action_logs (tenant_id, insight_id)
    WHERE insight_id IS NOT NULL AND NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.ai_human_overrides (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    insight_id UUID REFERENCES hrms.insights(id),
    action_id UUID REFERENCES hrms.ai_agent_action_logs(id),
    override_type TEXT NOT NULL,
    original_status TEXT,
    override_status TEXT NOT NULL,
    reason TEXT NOT NULL,
    decision TEXT NOT NULL,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by UUID,
    CHECK (decision IN ('accepted', 'rejected', 'replaced', 'manual_action')),
    CHECK (override_status IN ('overridden', 'dismissed', 'resolved', 'rejected', 'cancelled'))
);

CREATE INDEX IF NOT EXISTS idx_hrms_ai_human_overrides_tenant
    ON hrms.ai_human_overrides (tenant_id, created_at DESC)
    WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.ai_event_outbox (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    event_key TEXT NOT NULL,
    event_type TEXT NOT NULL,
    target_bus TEXT NOT NULL DEFAULT 'redis_stream',
    status TEXT NOT NULL DEFAULT 'pending',
    payload JSONB NOT NULL DEFAULT '{}'::jsonb,
    correlation_id TEXT,
    attempts INT NOT NULL DEFAULT 0,
    published_at TIMESTAMPTZ,
    last_error TEXT,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by UUID,
    CHECK (status IN ('pending', 'published', 'failed', 'skipped'))
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_hrms_ai_event_outbox_key
    ON hrms.ai_event_outbox (tenant_id, event_key)
    WHERE NOT inactive;

CREATE INDEX IF NOT EXISTS idx_hrms_ai_event_outbox_status
    ON hrms.ai_event_outbox (tenant_id, status, created_at DESC)
    WHERE NOT inactive;

DROP TRIGGER IF EXISTS set_ai_signal_logs_updated_at ON hrms.ai_signal_logs;
CREATE TRIGGER set_ai_signal_logs_updated_at
    BEFORE UPDATE ON hrms.ai_signal_logs
    FOR EACH ROW EXECUTE FUNCTION hrms.set_updated_at();

DROP TRIGGER IF EXISTS set_ai_agent_action_logs_updated_at ON hrms.ai_agent_action_logs;
CREATE TRIGGER set_ai_agent_action_logs_updated_at
    BEFORE UPDATE ON hrms.ai_agent_action_logs
    FOR EACH ROW EXECUTE FUNCTION hrms.set_updated_at();

DROP TRIGGER IF EXISTS set_ai_human_overrides_updated_at ON hrms.ai_human_overrides;
CREATE TRIGGER set_ai_human_overrides_updated_at
    BEFORE UPDATE ON hrms.ai_human_overrides
    FOR EACH ROW EXECUTE FUNCTION hrms.set_updated_at();

DROP TRIGGER IF EXISTS set_ai_event_outbox_updated_at ON hrms.ai_event_outbox;
CREATE TRIGGER set_ai_event_outbox_updated_at
    BEFORE UPDATE ON hrms.ai_event_outbox
    FOR EACH ROW EXECUTE FUNCTION hrms.set_updated_at();

ALTER TABLE hrms.ai_signal_logs ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.ai_agent_action_logs ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.ai_human_overrides ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.ai_event_outbox ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS ai_signal_logs_tenant_isolation ON hrms.ai_signal_logs;
CREATE POLICY ai_signal_logs_tenant_isolation ON hrms.ai_signal_logs
    USING (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true))
    WITH CHECK (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true));

DROP POLICY IF EXISTS ai_agent_action_logs_tenant_isolation ON hrms.ai_agent_action_logs;
CREATE POLICY ai_agent_action_logs_tenant_isolation ON hrms.ai_agent_action_logs
    USING (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true))
    WITH CHECK (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true));

DROP POLICY IF EXISTS ai_human_overrides_tenant_isolation ON hrms.ai_human_overrides;
CREATE POLICY ai_human_overrides_tenant_isolation ON hrms.ai_human_overrides
    USING (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true))
    WITH CHECK (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true));

DROP POLICY IF EXISTS ai_event_outbox_tenant_isolation ON hrms.ai_event_outbox;
CREATE POLICY ai_event_outbox_tenant_isolation ON hrms.ai_event_outbox
    USING (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true))
    WITH CHECK (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true));
