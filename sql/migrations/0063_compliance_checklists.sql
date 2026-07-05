CREATE TABLE IF NOT EXISTS hrms.compliance_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    code TEXT NOT NULL,
    title TEXT NOT NULL,
    description TEXT,
    category TEXT NOT NULL,
    scope TEXT NOT NULL DEFAULT 'worker',
    severity TEXT NOT NULL DEFAULT 'medium',
    classification_group TEXT,
    worker_type_id UUID,
    engagement_type TEXT,
    branch_id UUID,
    department_id UUID,
    country_code TEXT NOT NULL DEFAULT 'IN',
    state_code TEXT,
    trigger_event TEXT NOT NULL DEFAULT 'onboarding',
    default_due_days INTEGER NOT NULL DEFAULT 0,
    recurring_days INTEGER,
    requires_evidence BOOLEAN NOT NULL DEFAULT TRUE,
    evidence_label TEXT,
    auto_detect_key TEXT,
    blocks_payroll BOOLEAN NOT NULL DEFAULT FALSE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    effective_from DATE,
    effective_to DATE,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    CONSTRAINT compliance_rules_scope_chk CHECK (scope IN ('worker', 'engagement', 'worker_or_engagement')),
    CONSTRAINT compliance_rules_category_chk CHECK (category IN ('clra', 'fixed_term', 'gig_worker', 'tds', 'pf', 'esic', 'pt', 'lwf', 'document', 'safety', 'contract', 'custom')),
    CONSTRAINT compliance_rules_severity_chk CHECK (severity IN ('low', 'medium', 'high', 'critical')),
    CONSTRAINT compliance_rules_due_chk CHECK (default_due_days >= 0 AND (recurring_days IS NULL OR recurring_days > 0)),
    CONSTRAINT compliance_rules_effective_chk CHECK (effective_to IS NULL OR effective_from IS NULL OR effective_to >= effective_from),
    CONSTRAINT compliance_rules_worker_type_fk FOREIGN KEY (worker_type_id) REFERENCES hrms.worker_types(id),
    CONSTRAINT compliance_rules_branch_fk FOREIGN KEY (branch_id) REFERENCES hrms.branches(id),
    CONSTRAINT compliance_rules_department_fk FOREIGN KEY (department_id) REFERENCES hrms.departments(id)
);

CREATE UNIQUE INDEX IF NOT EXISTS compliance_rules_code_active_uq
    ON hrms.compliance_rules (tenant_id, lower(code))
    WHERE inactive = FALSE;

CREATE INDEX IF NOT EXISTS compliance_rules_tenant_active_idx
    ON hrms.compliance_rules (tenant_id, is_active, category, scope)
    WHERE inactive = FALSE;

CREATE TABLE IF NOT EXISTS hrms.compliance_checklist_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    rule_id UUID NOT NULL REFERENCES hrms.compliance_rules(id),
    worker_profile_id UUID REFERENCES hrms.worker_profiles(id),
    engagement_id UUID REFERENCES hrms.engagements(id),
    status TEXT NOT NULL DEFAULT 'pending',
    due_date DATE,
    completed_at TIMESTAMPTZ,
    reviewed_at TIMESTAMPTZ,
    reviewed_by UUID,
    evidence_path TEXT,
    evidence_file_name TEXT,
    evidence_content_type TEXT,
    evidence_uploaded_at TIMESTAMPTZ,
    evidence_uploaded_by UUID,
    waiver_reason TEXT,
    waiver_until DATE,
    waived_at TIMESTAMPTZ,
    waived_by UUID,
    detected_value TEXT,
    notes TEXT,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    CONSTRAINT compliance_checklist_status_chk CHECK (status IN ('pending', 'in_review', 'compliant', 'non_compliant', 'waived', 'expired', 'not_applicable')),
    CONSTRAINT compliance_checklist_target_chk CHECK (worker_profile_id IS NOT NULL OR engagement_id IS NOT NULL),
    CONSTRAINT compliance_checklist_waiver_chk CHECK ((status <> 'waived') OR waiver_reason IS NOT NULL)
);

CREATE UNIQUE INDEX IF NOT EXISTS compliance_checklist_rule_worker_active_uq
    ON hrms.compliance_checklist_items (tenant_id, rule_id, worker_profile_id)
    WHERE inactive = FALSE AND worker_profile_id IS NOT NULL AND engagement_id IS NULL;

CREATE UNIQUE INDEX IF NOT EXISTS compliance_checklist_rule_engagement_active_uq
    ON hrms.compliance_checklist_items (tenant_id, rule_id, engagement_id)
    WHERE inactive = FALSE AND engagement_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS compliance_checklist_tenant_status_idx
    ON hrms.compliance_checklist_items (tenant_id, status, due_date)
    WHERE inactive = FALSE;

CREATE INDEX IF NOT EXISTS compliance_checklist_worker_idx
    ON hrms.compliance_checklist_items (tenant_id, worker_profile_id)
    WHERE inactive = FALSE;

CREATE INDEX IF NOT EXISTS compliance_checklist_engagement_idx
    ON hrms.compliance_checklist_items (tenant_id, engagement_id)
    WHERE inactive = FALSE;

CREATE TABLE IF NOT EXISTS hrms.compliance_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    checklist_item_id UUID REFERENCES hrms.compliance_checklist_items(id),
    rule_id UUID REFERENCES hrms.compliance_rules(id),
    event_type TEXT NOT NULL,
    from_status TEXT,
    to_status TEXT,
    comment TEXT,
    actor_id UUID,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS compliance_events_item_idx
    ON hrms.compliance_events (tenant_id, checklist_item_id, created_at DESC);

CREATE INDEX IF NOT EXISTS compliance_events_rule_idx
    ON hrms.compliance_events (tenant_id, rule_id, created_at DESC);

DROP TRIGGER IF EXISTS trg_compliance_rules_updated_at ON hrms.compliance_rules;
CREATE TRIGGER trg_compliance_rules_updated_at
BEFORE UPDATE ON hrms.compliance_rules
FOR EACH ROW EXECUTE FUNCTION hrms.set_updated_at();

DROP TRIGGER IF EXISTS trg_compliance_checklist_items_updated_at ON hrms.compliance_checklist_items;
CREATE TRIGGER trg_compliance_checklist_items_updated_at
BEFORE UPDATE ON hrms.compliance_checklist_items
FOR EACH ROW EXECUTE FUNCTION hrms.set_updated_at();

ALTER TABLE hrms.compliance_rules ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.compliance_checklist_items ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.compliance_events ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS compliance_rules_tenant_isolation ON hrms.compliance_rules;
CREATE POLICY compliance_rules_tenant_isolation ON hrms.compliance_rules
USING (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true')
WITH CHECK (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true');

DROP POLICY IF EXISTS compliance_checklist_items_tenant_isolation ON hrms.compliance_checklist_items;
CREATE POLICY compliance_checklist_items_tenant_isolation ON hrms.compliance_checklist_items
USING (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true')
WITH CHECK (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true');

DROP POLICY IF EXISTS compliance_events_tenant_isolation ON hrms.compliance_events;
CREATE POLICY compliance_events_tenant_isolation ON hrms.compliance_events
USING (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true')
WITH CHECK (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true');
