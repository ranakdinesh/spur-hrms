CREATE TABLE IF NOT EXISTS hrms.projects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    project_code TEXT,
    name TEXT NOT NULL,
    description TEXT,
    status TEXT NOT NULL DEFAULT 'draft',
    department_id UUID REFERENCES hrms.departments(id) ON DELETE SET NULL,
    branch_id UUID REFERENCES hrms.branches(id) ON DELETE SET NULL,
    project_manager_id UUID,
    start_date DATE,
    due_date DATE,
    completed_at TIMESTAMPTZ,
    budget_amount NUMERIC(14,2),
    currency_code CHAR(3) NOT NULL DEFAULT 'INR',
    billing_type TEXT NOT NULL DEFAULT 'none',
    client_label TEXT,
    priority TEXT NOT NULL DEFAULT 'normal',
    notes TEXT,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    CONSTRAINT projects_status_check CHECK (status IN ('draft', 'active', 'on_hold', 'completed', 'cancelled')),
    CONSTRAINT projects_billing_type_check CHECK (billing_type IN ('none', 'fixed', 'hourly', 'milestone', 'retainer')),
    CONSTRAINT projects_priority_check CHECK (priority IN ('low', 'normal', 'high', 'critical')),
    CONSTRAINT projects_budget_amount_check CHECK (budget_amount IS NULL OR budget_amount >= 0),
    CONSTRAINT projects_date_order_check CHECK (start_date IS NULL OR due_date IS NULL OR due_date >= start_date),
    CONSTRAINT projects_currency_code_check CHECK (currency_code ~ '^[A-Z]{3}$'),
    CONSTRAINT projects_metadata_object_check CHECK (jsonb_typeof(metadata) = 'object')
);

CREATE UNIQUE INDEX IF NOT EXISTS projects_tenant_code_active_uq
    ON hrms.projects (tenant_id, lower(project_code))
    WHERE project_code IS NOT NULL AND inactive = FALSE;

CREATE INDEX IF NOT EXISTS projects_tenant_status_due_idx
    ON hrms.projects (tenant_id, status, due_date)
    WHERE inactive = FALSE;

CREATE INDEX IF NOT EXISTS projects_tenant_department_idx
    ON hrms.projects (tenant_id, department_id)
    WHERE inactive = FALSE;

CREATE INDEX IF NOT EXISTS projects_tenant_manager_idx
    ON hrms.projects (tenant_id, project_manager_id)
    WHERE inactive = FALSE;

DROP TRIGGER IF EXISTS trg_projects_updated_at ON hrms.projects;
CREATE TRIGGER trg_projects_updated_at
BEFORE UPDATE ON hrms.projects
FOR EACH ROW
EXECUTE FUNCTION hrms.set_updated_at();

ALTER TABLE hrms.projects ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS projects_tenant_isolation ON hrms.projects;
CREATE POLICY projects_tenant_isolation ON hrms.projects
USING (
    current_setting('app.is_super_admin', true) = 'true'
    OR tenant_id::text = current_setting('app.tenant_id', true)
)
WITH CHECK (
    current_setting('app.is_super_admin', true) = 'true'
    OR tenant_id::text = current_setting('app.tenant_id', true)
);

CREATE TABLE IF NOT EXISTS hrms.project_milestones (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    project_id UUID NOT NULL REFERENCES hrms.projects(id) ON DELETE RESTRICT,
    engagement_id UUID REFERENCES hrms.engagements(id) ON DELETE SET NULL,
    milestone_code TEXT,
    title TEXT NOT NULL,
    description TEXT,
    acceptance_criteria TEXT,
    due_date DATE,
    status TEXT NOT NULL DEFAULT 'draft',
    amount NUMERIC(14,2),
    currency_code CHAR(3) NOT NULL DEFAULT 'INR',
    payment_trigger JSONB NOT NULL DEFAULT '{}'::jsonb,
    submitted_at TIMESTAMPTZ,
    submitted_by UUID,
    accepted_at TIMESTAMPTZ,
    accepted_by UUID,
    rejected_at TIMESTAMPTZ,
    rejected_by UUID,
    review_comment TEXT,
    notes TEXT,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    CONSTRAINT project_milestones_status_check CHECK (status IN ('draft', 'open', 'submitted', 'accepted', 'rejected', 'cancelled')),
    CONSTRAINT project_milestones_amount_check CHECK (amount IS NULL OR amount >= 0),
    CONSTRAINT project_milestones_currency_code_check CHECK (currency_code ~ '^[A-Z]{3}$'),
    CONSTRAINT project_milestones_payment_trigger_object_check CHECK (jsonb_typeof(payment_trigger) = 'object'),
    CONSTRAINT project_milestones_metadata_object_check CHECK (jsonb_typeof(metadata) = 'object')
);

CREATE UNIQUE INDEX IF NOT EXISTS project_milestones_tenant_project_code_active_uq
    ON hrms.project_milestones (tenant_id, project_id, lower(milestone_code))
    WHERE milestone_code IS NOT NULL AND inactive = FALSE;

CREATE INDEX IF NOT EXISTS project_milestones_tenant_project_status_idx
    ON hrms.project_milestones (tenant_id, project_id, status)
    WHERE inactive = FALSE;

CREATE INDEX IF NOT EXISTS project_milestones_tenant_status_due_idx
    ON hrms.project_milestones (tenant_id, status, due_date)
    WHERE inactive = FALSE;

CREATE INDEX IF NOT EXISTS project_milestones_tenant_engagement_idx
    ON hrms.project_milestones (tenant_id, engagement_id)
    WHERE inactive = FALSE;

DROP TRIGGER IF EXISTS trg_project_milestones_updated_at ON hrms.project_milestones;
CREATE TRIGGER trg_project_milestones_updated_at
BEFORE UPDATE ON hrms.project_milestones
FOR EACH ROW
EXECUTE FUNCTION hrms.set_updated_at();

ALTER TABLE hrms.project_milestones ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS project_milestones_tenant_isolation ON hrms.project_milestones;
CREATE POLICY project_milestones_tenant_isolation ON hrms.project_milestones
USING (
    current_setting('app.is_super_admin', true) = 'true'
    OR tenant_id::text = current_setting('app.tenant_id', true)
)
WITH CHECK (
    current_setting('app.is_super_admin', true) = 'true'
    OR tenant_id::text = current_setting('app.tenant_id', true)
);

CREATE TABLE IF NOT EXISTS hrms.project_milestone_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    project_id UUID NOT NULL REFERENCES hrms.projects(id) ON DELETE CASCADE,
    milestone_id UUID NOT NULL REFERENCES hrms.project_milestones(id) ON DELETE CASCADE,
    event_type TEXT NOT NULL,
    from_status TEXT,
    to_status TEXT,
    comment TEXT,
    actor_id UUID,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT project_milestone_events_metadata_object_check CHECK (jsonb_typeof(metadata) = 'object')
);

CREATE INDEX IF NOT EXISTS project_milestone_events_tenant_milestone_created_idx
    ON hrms.project_milestone_events (tenant_id, milestone_id, created_at DESC);

CREATE INDEX IF NOT EXISTS project_milestone_events_tenant_project_created_idx
    ON hrms.project_milestone_events (tenant_id, project_id, created_at DESC);

ALTER TABLE hrms.project_milestone_events ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS project_milestone_events_tenant_isolation ON hrms.project_milestone_events;
CREATE POLICY project_milestone_events_tenant_isolation ON hrms.project_milestone_events
USING (
    current_setting('app.is_super_admin', true) = 'true'
    OR tenant_id::text = current_setting('app.tenant_id', true)
)
WITH CHECK (
    current_setting('app.is_super_admin', true) = 'true'
    OR tenant_id::text = current_setting('app.tenant_id', true)
);
