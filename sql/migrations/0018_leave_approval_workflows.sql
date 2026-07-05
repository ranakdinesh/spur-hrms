CREATE TABLE IF NOT EXISTS hrms.leave_approval_workflows (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    name        VARCHAR(120) NOT NULL,
    code        VARCHAR(60) NOT NULL,
    description TEXT,
    is_default  BOOLEAN NOT NULL DEFAULT FALSE,
    inactive    BOOLEAN NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by  UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by  UUID REFERENCES auth.users(id) ON DELETE SET NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS leave_approval_workflows_code_idx ON hrms.leave_approval_workflows(tenant_id, lower(code)) WHERE NOT inactive;
CREATE UNIQUE INDEX IF NOT EXISTS leave_approval_workflows_default_idx ON hrms.leave_approval_workflows(tenant_id) WHERE is_default AND NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.leave_approval_workflow_steps (
    id                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id          UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    workflow_id        UUID NOT NULL REFERENCES hrms.leave_approval_workflows(id) ON DELETE CASCADE,
    step_order         INT NOT NULL,
    name               VARCHAR(120) NOT NULL,
    approver_type      VARCHAR(40) NOT NULL,
    approver_user_id   UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    approver_role      VARCHAR(80),
    decision_rule      VARCHAR(20) NOT NULL DEFAULT 'all',
    required_approvals INT NOT NULL DEFAULT 1,
    auto_approve       BOOLEAN NOT NULL DEFAULT FALSE,
    sla_hours          INT NOT NULL DEFAULT 0,
    inactive           BOOLEAN NOT NULL DEFAULT FALSE,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by         UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by         UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT leave_approval_workflow_steps_order_check CHECK (step_order > 0),
    CONSTRAINT leave_approval_workflow_steps_type_check CHECK (approver_type IN ('reporting_manager','manager_manager','hr_user','specific_user','role','applicant')),
    CONSTRAINT leave_approval_workflow_steps_rule_check CHECK (decision_rule IN ('all','any')),
    CONSTRAINT leave_approval_workflow_steps_required_check CHECK (required_approvals > 0 AND sla_hours >= 0)
);
CREATE UNIQUE INDEX IF NOT EXISTS leave_approval_workflow_steps_order_idx ON hrms.leave_approval_workflow_steps(tenant_id, workflow_id, step_order, approver_type, COALESCE(approver_user_id, '00000000-0000-0000-0000-000000000000'::uuid), COALESCE(approver_role, '')) WHERE NOT inactive;
CREATE INDEX IF NOT EXISTS leave_approval_workflow_steps_workflow_idx ON hrms.leave_approval_workflow_steps(tenant_id, workflow_id, step_order) WHERE NOT inactive;

ALTER TABLE hrms.leave_approvals ADD COLUMN IF NOT EXISTS workflow_id UUID REFERENCES hrms.leave_approval_workflows(id) ON DELETE SET NULL;
ALTER TABLE hrms.leave_approvals ADD COLUMN IF NOT EXISTS workflow_step_id UUID REFERENCES hrms.leave_approval_workflow_steps(id) ON DELETE SET NULL;
ALTER TABLE hrms.leave_approvals ADD COLUMN IF NOT EXISTS step_order INT NOT NULL DEFAULT 1;
ALTER TABLE hrms.leave_approvals ADD COLUMN IF NOT EXISTS decision_rule VARCHAR(20) NOT NULL DEFAULT 'all';
ALTER TABLE hrms.leave_approvals ADD COLUMN IF NOT EXISTS required_approvals INT NOT NULL DEFAULT 1;
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'leave_approvals_decision_rule_check') THEN
        ALTER TABLE hrms.leave_approvals ADD CONSTRAINT leave_approvals_decision_rule_check CHECK (decision_rule IN ('all','any'));
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'leave_approvals_required_check') THEN
        ALTER TABLE hrms.leave_approvals ADD CONSTRAINT leave_approvals_required_check CHECK (required_approvals > 0 AND step_order > 0);
    END IF;
END $$;
CREATE INDEX IF NOT EXISTS leave_approvals_workflow_step_idx ON hrms.leave_approvals(tenant_id, leave_id, step_order) WHERE NOT inactive;

DO $$
DECLARE
    table_name TEXT;
    updated_tables TEXT[] := ARRAY['leave_approval_workflows','leave_approval_workflow_steps'];
BEGIN
    FOREACH table_name IN ARRAY updated_tables LOOP
        EXECUTE format('DROP TRIGGER IF EXISTS %I ON hrms.%I', table_name || '_updated_at', table_name);
        EXECUTE format('CREATE TRIGGER %I BEFORE UPDATE ON hrms.%I FOR EACH ROW EXECUTE FUNCTION hrms.update_updated_at()', table_name || '_updated_at', table_name);
        EXECUTE format('ALTER TABLE hrms.%I ENABLE ROW LEVEL SECURITY', table_name);
        EXECUTE format('DROP POLICY IF EXISTS %I ON hrms.%I', table_name || '_tenant_isolation', table_name);
        EXECUTE format(
            'CREATE POLICY %I ON hrms.%I USING (tenant_id = NULLIF(current_setting(''app.tenant_id'', true), '''')::uuid OR current_setting(''app.is_super_admin'', true) = ''true'') WITH CHECK (tenant_id = NULLIF(current_setting(''app.tenant_id'', true), '''')::uuid OR current_setting(''app.is_super_admin'', true) = ''true'')',
            table_name || '_tenant_isolation',
            table_name
        );
    END LOOP;
END $$;
