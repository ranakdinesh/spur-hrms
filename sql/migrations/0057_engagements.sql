CREATE TABLE IF NOT EXISTS hrms.engagements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    worker_profile_id UUID NOT NULL REFERENCES hrms.worker_profiles(id) ON DELETE RESTRICT,
    engagement_code VARCHAR(80),
    title VARCHAR(180) NOT NULL,
    description TEXT,
    engagement_type VARCHAR(40) NOT NULL,
    status VARCHAR(40) NOT NULL DEFAULT 'draft',
    start_date DATE NOT NULL,
    end_date DATE,
    hours_budget NUMERIC(12,2),
    rate_amount NUMERIC(14,2),
    currency_code CHAR(3) NOT NULL DEFAULT 'INR',
    rate_unit VARCHAR(40) NOT NULL DEFAULT 'none',
    branch_id UUID REFERENCES hrms.branches(id) ON DELETE SET NULL,
    department_id UUID REFERENCES hrms.departments(id) ON DELETE SET NULL,
    reporting_manager_id UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    project_label VARCHAR(180),
    project_code VARCHAR(80),
    cost_center VARCHAR(80),
    renewal_due_date DATE,
    renewal_status VARCHAR(40) NOT NULL DEFAULT 'not_required',
    termination_reason TEXT,
    terminated_at TIMESTAMPTZ,
    notes TEXT,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT engagements_type_check CHECK (engagement_type IN ('employee_assignment','fixed_term','project','hourly','retainer','stipend','agency','consulting')),
    CONSTRAINT engagements_status_check CHECK (status IN ('draft','active','paused','completed','terminated','cancelled')),
    CONSTRAINT engagements_dates_check CHECK (end_date IS NULL OR end_date >= start_date),
    CONSTRAINT engagements_hours_budget_check CHECK (hours_budget IS NULL OR hours_budget >= 0),
    CONSTRAINT engagements_rate_amount_check CHECK (rate_amount IS NULL OR rate_amount >= 0),
    CONSTRAINT engagements_currency_code_check CHECK (currency_code ~ '^[A-Z]{3}$'),
    CONSTRAINT engagements_rate_unit_check CHECK (rate_unit IN ('none','hour','day','month','milestone','retainer','stipend')),
    CONSTRAINT engagements_renewal_status_check CHECK (renewal_status IN ('not_required','pending','renewed','not_renewed')),
    CONSTRAINT engagements_metadata_object_check CHECK (jsonb_typeof(metadata) = 'object')
);

CREATE UNIQUE INDEX IF NOT EXISTS engagements_tenant_code_idx
    ON hrms.engagements (tenant_id, lower(engagement_code))
    WHERE engagement_code IS NOT NULL AND NOT inactive;
CREATE INDEX IF NOT EXISTS engagements_tenant_worker_idx
    ON hrms.engagements (tenant_id, worker_profile_id, status)
    WHERE NOT inactive;
CREATE INDEX IF NOT EXISTS engagements_tenant_status_idx
    ON hrms.engagements (tenant_id, status, start_date DESC)
    WHERE NOT inactive;
CREATE INDEX IF NOT EXISTS engagements_tenant_type_idx
    ON hrms.engagements (tenant_id, engagement_type, status)
    WHERE NOT inactive;
CREATE INDEX IF NOT EXISTS engagements_tenant_org_idx
    ON hrms.engagements (tenant_id, branch_id, department_id)
    WHERE NOT inactive;
CREATE INDEX IF NOT EXISTS engagements_tenant_renewal_idx
    ON hrms.engagements (tenant_id, renewal_due_date)
    WHERE renewal_due_date IS NOT NULL AND NOT inactive;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'engagements_updated_at') THEN
        CREATE TRIGGER engagements_updated_at
        BEFORE UPDATE ON hrms.engagements
        FOR EACH ROW EXECUTE FUNCTION hrms.update_updated_at();
    END IF;
END $$;

ALTER TABLE hrms.engagements ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS engagements_tenant_isolation ON hrms.engagements;
CREATE POLICY engagements_tenant_isolation ON hrms.engagements
USING (tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::uuid OR current_setting('app.is_super_admin', true) = 'true')
WITH CHECK (tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::uuid OR current_setting('app.is_super_admin', true) = 'true');
