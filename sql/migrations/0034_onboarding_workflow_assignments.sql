ALTER TABLE hrms.onboarding_workflows
    ADD COLUMN IF NOT EXISTS is_default BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS is_active BOOLEAN NOT NULL DEFAULT TRUE;

CREATE UNIQUE INDEX IF NOT EXISTS onboarding_workflows_default_idx
    ON hrms.onboarding_workflows(tenant_id)
    WHERE is_default AND is_active AND NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.onboarding_workflow_assignments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    workflow_id UUID NOT NULL REFERENCES hrms.onboarding_workflows(id) ON DELETE CASCADE,
    name VARCHAR(160) NOT NULL,
    job_posting_id UUID REFERENCES hrms.job_postings(id) ON DELETE SET NULL,
    job_position_id UUID REFERENCES hrms.job_positions(id) ON DELETE SET NULL,
    department_id UUID REFERENCES hrms.departments(id) ON DELETE SET NULL,
    employment_type_id UUID REFERENCES hrms.employment_types(id) ON DELETE SET NULL,
    priority INT NOT NULL DEFAULT 100,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT onboarding_workflow_assignments_name_check CHECK (BTRIM(name) <> '')
);

CREATE INDEX IF NOT EXISTS onboarding_workflow_assignments_tenant_idx
    ON hrms.onboarding_workflow_assignments(tenant_id, priority, name)
    WHERE NOT inactive;

CREATE INDEX IF NOT EXISTS onboarding_workflow_assignments_workflow_idx
    ON hrms.onboarding_workflow_assignments(tenant_id, workflow_id)
    WHERE NOT inactive;

ALTER TABLE hrms.onboarding_workflow_assignments ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS tenant_isolation_onboarding_workflow_assignments ON hrms.onboarding_workflow_assignments;
CREATE POLICY tenant_isolation_onboarding_workflow_assignments ON hrms.onboarding_workflow_assignments
USING (
    current_setting('app.is_super_admin', true) = 'true'
    OR tenant_id::text = current_setting('app.tenant_id', true)
)
WITH CHECK (
    current_setting('app.is_super_admin', true) = 'true'
    OR tenant_id::text = current_setting('app.tenant_id', true)
);
