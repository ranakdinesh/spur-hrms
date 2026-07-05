CREATE TABLE IF NOT EXISTS hrms.okr_cycles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    cycle_code TEXT NOT NULL,
    description TEXT,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    status TEXT NOT NULL DEFAULT 'draft',
    review_cadence TEXT NOT NULL DEFAULT 'weekly',
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    CONSTRAINT okr_cycles_status_chk CHECK (status IN ('draft', 'active', 'closed', 'archived')),
    CONSTRAINT okr_cycles_cadence_chk CHECK (review_cadence IN ('weekly', 'biweekly', 'monthly')),
    CONSTRAINT okr_cycles_dates_chk CHECK (end_date >= start_date),
    CONSTRAINT okr_cycles_metadata_object_chk CHECK (jsonb_typeof(metadata) = 'object')
);

CREATE UNIQUE INDEX IF NOT EXISTS okr_cycles_tenant_code_uq
    ON hrms.okr_cycles (tenant_id, lower(cycle_code))
    WHERE inactive = FALSE;

CREATE INDEX IF NOT EXISTS okr_cycles_tenant_status_idx
    ON hrms.okr_cycles (tenant_id, status, start_date DESC)
    WHERE inactive = FALSE;

CREATE TABLE IF NOT EXISTS hrms.objectives (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    cycle_id UUID NOT NULL REFERENCES hrms.okr_cycles(id) ON DELETE CASCADE,
    parent_objective_id UUID REFERENCES hrms.objectives(id) ON DELETE SET NULL,
    owner_type TEXT NOT NULL,
    owner_worker_profile_id UUID REFERENCES hrms.worker_profiles(id) ON DELETE SET NULL,
    owner_department_id UUID REFERENCES hrms.departments(id) ON DELETE SET NULL,
    owner_project_id UUID REFERENCES hrms.projects(id) ON DELETE SET NULL,
    title TEXT NOT NULL,
    description TEXT,
    status TEXT NOT NULL DEFAULT 'draft',
    priority TEXT NOT NULL DEFAULT 'normal',
    progress_percent NUMERIC(5,2) NOT NULL DEFAULT 0,
    weight NUMERIC(6,2) NOT NULL DEFAULT 1,
    start_date DATE,
    due_date DATE,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    CONSTRAINT objectives_owner_type_chk CHECK (owner_type IN ('company', 'department', 'project', 'worker')),
    CONSTRAINT objectives_owner_required_chk CHECK (
        (owner_type = 'company' AND owner_worker_profile_id IS NULL AND owner_department_id IS NULL AND owner_project_id IS NULL)
        OR (owner_type = 'department' AND owner_department_id IS NOT NULL AND owner_worker_profile_id IS NULL AND owner_project_id IS NULL)
        OR (owner_type = 'project' AND owner_project_id IS NOT NULL AND owner_worker_profile_id IS NULL AND owner_department_id IS NULL)
        OR (owner_type = 'worker' AND owner_worker_profile_id IS NOT NULL AND owner_department_id IS NULL AND owner_project_id IS NULL)
    ),
    CONSTRAINT objectives_status_chk CHECK (status IN ('draft', 'active', 'at_risk', 'completed', 'closed', 'cancelled')),
    CONSTRAINT objectives_priority_chk CHECK (priority IN ('low', 'normal', 'high', 'critical')),
    CONSTRAINT objectives_progress_chk CHECK (progress_percent >= 0 AND progress_percent <= 100),
    CONSTRAINT objectives_weight_chk CHECK (weight > 0),
    CONSTRAINT objectives_dates_chk CHECK (start_date IS NULL OR due_date IS NULL OR due_date >= start_date),
    CONSTRAINT objectives_metadata_object_chk CHECK (jsonb_typeof(metadata) = 'object')
);

CREATE INDEX IF NOT EXISTS objectives_tenant_cycle_idx
    ON hrms.objectives (tenant_id, cycle_id, owner_type, status)
    WHERE inactive = FALSE;

CREATE INDEX IF NOT EXISTS objectives_parent_idx
    ON hrms.objectives (tenant_id, parent_objective_id)
    WHERE parent_objective_id IS NOT NULL AND inactive = FALSE;

CREATE INDEX IF NOT EXISTS objectives_worker_owner_idx
    ON hrms.objectives (tenant_id, owner_worker_profile_id)
    WHERE owner_worker_profile_id IS NOT NULL AND inactive = FALSE;

CREATE TABLE IF NOT EXISTS hrms.key_results (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    objective_id UUID NOT NULL REFERENCES hrms.objectives(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    description TEXT,
    metric_type TEXT NOT NULL DEFAULT 'number',
    start_value NUMERIC(14,2) NOT NULL DEFAULT 0,
    target_value NUMERIC(14,2) NOT NULL DEFAULT 100,
    current_value NUMERIC(14,2) NOT NULL DEFAULT 0,
    progress_percent NUMERIC(5,2) NOT NULL DEFAULT 0,
    confidence TEXT NOT NULL DEFAULT 'medium',
    status TEXT NOT NULL DEFAULT 'not_started',
    weight NUMERIC(6,2) NOT NULL DEFAULT 1,
    unit_label TEXT,
    due_date DATE,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    CONSTRAINT key_results_metric_type_chk CHECK (metric_type IN ('number', 'percent', 'currency', 'boolean')),
    CONSTRAINT key_results_status_chk CHECK (status IN ('not_started', 'on_track', 'at_risk', 'behind', 'completed', 'closed', 'cancelled')),
    CONSTRAINT key_results_confidence_chk CHECK (confidence IN ('low', 'medium', 'high')),
    CONSTRAINT key_results_progress_chk CHECK (progress_percent >= 0 AND progress_percent <= 100),
    CONSTRAINT key_results_weight_chk CHECK (weight > 0),
    CONSTRAINT key_results_metadata_object_chk CHECK (jsonb_typeof(metadata) = 'object')
);

CREATE INDEX IF NOT EXISTS key_results_objective_idx
    ON hrms.key_results (tenant_id, objective_id, status)
    WHERE inactive = FALSE;

CREATE TABLE IF NOT EXISTS hrms.key_result_checkins (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    key_result_id UUID NOT NULL REFERENCES hrms.key_results(id) ON DELETE CASCADE,
    checkin_date DATE NOT NULL DEFAULT CURRENT_DATE,
    value NUMERIC(14,2) NOT NULL,
    progress_percent NUMERIC(5,2) NOT NULL,
    confidence TEXT NOT NULL DEFAULT 'medium',
    status TEXT NOT NULL DEFAULT 'on_track',
    note TEXT,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID,
    CONSTRAINT key_result_checkins_status_chk CHECK (status IN ('not_started', 'on_track', 'at_risk', 'behind', 'completed', 'closed', 'cancelled')),
    CONSTRAINT key_result_checkins_confidence_chk CHECK (confidence IN ('low', 'medium', 'high')),
    CONSTRAINT key_result_checkins_progress_chk CHECK (progress_percent >= 0 AND progress_percent <= 100),
    CONSTRAINT key_result_checkins_metadata_object_chk CHECK (jsonb_typeof(metadata) = 'object')
);

CREATE INDEX IF NOT EXISTS key_result_checkins_kr_idx
    ON hrms.key_result_checkins (tenant_id, key_result_id, checkin_date DESC, created_at DESC)
    WHERE inactive = FALSE;

DROP TRIGGER IF EXISTS trg_okr_cycles_updated_at ON hrms.okr_cycles;
CREATE TRIGGER trg_okr_cycles_updated_at
BEFORE UPDATE ON hrms.okr_cycles
FOR EACH ROW EXECUTE FUNCTION hrms.set_updated_at();

DROP TRIGGER IF EXISTS trg_objectives_updated_at ON hrms.objectives;
CREATE TRIGGER trg_objectives_updated_at
BEFORE UPDATE ON hrms.objectives
FOR EACH ROW EXECUTE FUNCTION hrms.set_updated_at();

DROP TRIGGER IF EXISTS trg_key_results_updated_at ON hrms.key_results;
CREATE TRIGGER trg_key_results_updated_at
BEFORE UPDATE ON hrms.key_results
FOR EACH ROW EXECUTE FUNCTION hrms.set_updated_at();

ALTER TABLE hrms.okr_cycles ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.objectives ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.key_results ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.key_result_checkins ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS okr_cycles_tenant_isolation ON hrms.okr_cycles;
CREATE POLICY okr_cycles_tenant_isolation ON hrms.okr_cycles
USING (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true')
WITH CHECK (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true');

DROP POLICY IF EXISTS objectives_tenant_isolation ON hrms.objectives;
CREATE POLICY objectives_tenant_isolation ON hrms.objectives
USING (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true')
WITH CHECK (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true');

DROP POLICY IF EXISTS key_results_tenant_isolation ON hrms.key_results;
CREATE POLICY key_results_tenant_isolation ON hrms.key_results
USING (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true')
WITH CHECK (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true');

DROP POLICY IF EXISTS key_result_checkins_tenant_isolation ON hrms.key_result_checkins;
CREATE POLICY key_result_checkins_tenant_isolation ON hrms.key_result_checkins
USING (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true')
WITH CHECK (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true');
