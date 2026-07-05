CREATE TABLE IF NOT EXISTS hrms.insights (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    insight_key TEXT NOT NULL,
    insight_type TEXT NOT NULL,
    category TEXT NOT NULL,
    severity TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'open',
    title TEXT NOT NULL,
    summary TEXT NOT NULL,
    confidence_score NUMERIC(5,2) NOT NULL DEFAULT 0,
    score NUMERIC(8,2) NOT NULL DEFAULT 0,
    source TEXT NOT NULL DEFAULT 'deterministic',
    model_version TEXT,
    entity_type TEXT,
    entity_id UUID,
    employee_user_id UUID,
    reasons JSONB NOT NULL DEFAULT '[]'::jsonb,
    recommendations JSONB NOT NULL DEFAULT '[]'::jsonb,
    context JSONB NOT NULL DEFAULT '{}'::jsonb,
    explainability JSONB NOT NULL DEFAULT '{}'::jsonb,
    detected_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    due_at TIMESTAMPTZ,
    assigned_to UUID,
    reviewed_by UUID,
    reviewed_at TIMESTAMPTZ,
    resolved_at TIMESTAMPTZ,
    resolution_note TEXT,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by UUID,
    UNIQUE (tenant_id, insight_key)
);

CREATE INDEX IF NOT EXISTS idx_hrms_insights_tenant_status ON hrms.insights (tenant_id, status, severity, detected_at DESC) WHERE NOT inactive;
CREATE INDEX IF NOT EXISTS idx_hrms_insights_employee ON hrms.insights (tenant_id, employee_user_id) WHERE employee_user_id IS NOT NULL AND NOT inactive;
CREATE INDEX IF NOT EXISTS idx_hrms_insights_entity ON hrms.insights (tenant_id, entity_type, entity_id) WHERE entity_id IS NOT NULL AND NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.insight_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    insight_id UUID NOT NULL REFERENCES hrms.insights(id),
    action TEXT NOT NULL,
    from_status TEXT,
    to_status TEXT,
    remarks TEXT,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by UUID
);

CREATE INDEX IF NOT EXISTS idx_hrms_insight_events_insight ON hrms.insight_events (tenant_id, insight_id, created_at DESC) WHERE NOT inactive;

ALTER TABLE hrms.insights ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.insight_events ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS insights_tenant_isolation ON hrms.insights;
CREATE POLICY insights_tenant_isolation ON hrms.insights
    USING (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true))
    WITH CHECK (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true));

DROP POLICY IF EXISTS insight_events_tenant_isolation ON hrms.insight_events;
CREATE POLICY insight_events_tenant_isolation ON hrms.insight_events
    USING (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true))
    WITH CHECK (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true));
