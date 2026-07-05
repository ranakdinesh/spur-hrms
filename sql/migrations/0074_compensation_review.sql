CREATE TABLE IF NOT EXISTS hrms.compensation_pay_bands (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    code TEXT NOT NULL,
    name TEXT NOT NULL,
    job_family TEXT,
    level_code TEXT,
    location_label TEXT,
    currency_code TEXT NOT NULL DEFAULT 'INR',
    min_pay NUMERIC(14,2) NOT NULL,
    midpoint_pay NUMERIC(14,2) NOT NULL,
    max_pay NUMERIC(14,2) NOT NULL,
    effective_from DATE,
    effective_to DATE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    notes TEXT,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    CONSTRAINT compensation_pay_bands_range_chk CHECK (min_pay >= 0 AND midpoint_pay >= min_pay AND max_pay >= midpoint_pay),
    CONSTRAINT compensation_pay_bands_dates_chk CHECK (effective_to IS NULL OR effective_from IS NULL OR effective_to >= effective_from)
);

CREATE UNIQUE INDEX IF NOT EXISTS compensation_pay_bands_tenant_code_uq
    ON hrms.compensation_pay_bands (tenant_id, lower(code))
    WHERE inactive = FALSE;

CREATE INDEX IF NOT EXISTS compensation_pay_bands_tenant_filters_idx
    ON hrms.compensation_pay_bands (tenant_id, is_active, currency_code, level_code)
    WHERE inactive = FALSE;

CREATE TABLE IF NOT EXISTS hrms.compensation_cycles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    code TEXT NOT NULL,
    name TEXT NOT NULL,
    fiscal_year_id UUID REFERENCES hrms.financial_years(id),
    status TEXT NOT NULL DEFAULT 'draft',
    cycle_type TEXT NOT NULL DEFAULT 'annual',
    starts_on DATE,
    ends_on DATE,
    effective_date DATE,
    currency_code TEXT NOT NULL DEFAULT 'INR',
    budget_amount NUMERIC(14,2) NOT NULL DEFAULT 0,
    planning_guidance TEXT,
    approval_policy TEXT,
    finalized_at TIMESTAMPTZ,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    CONSTRAINT compensation_cycles_status_chk CHECK (status IN ('draft','open','submitted','approved','finalized','cancelled')),
    CONSTRAINT compensation_cycles_type_chk CHECK (cycle_type IN ('annual','mid_year','promotion','market_adjustment','equity','bonus','custom')),
    CONSTRAINT compensation_cycles_budget_chk CHECK (budget_amount >= 0),
    CONSTRAINT compensation_cycles_dates_chk CHECK (ends_on IS NULL OR starts_on IS NULL OR ends_on >= starts_on)
);

CREATE UNIQUE INDEX IF NOT EXISTS compensation_cycles_tenant_code_uq
    ON hrms.compensation_cycles (tenant_id, lower(code))
    WHERE inactive = FALSE;

CREATE INDEX IF NOT EXISTS compensation_cycles_tenant_status_idx
    ON hrms.compensation_cycles (tenant_id, status, effective_date DESC)
    WHERE inactive = FALSE;

CREATE TABLE IF NOT EXISTS hrms.compensation_budget_pools (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    cycle_id UUID NOT NULL REFERENCES hrms.compensation_cycles(id),
    name TEXT NOT NULL,
    pool_type TEXT NOT NULL DEFAULT 'merit',
    owner_user_id UUID,
    department_id UUID REFERENCES hrms.departments(id),
    branch_id UUID REFERENCES hrms.branches(id),
    budget_amount NUMERIC(14,2) NOT NULL DEFAULT 0,
    allocated_amount NUMERIC(14,2) NOT NULL DEFAULT 0,
    notes TEXT,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    CONSTRAINT compensation_budget_pools_type_chk CHECK (pool_type IN ('merit','promotion','equity','retention','market_adjustment','bonus','custom')),
    CONSTRAINT compensation_budget_pools_budget_chk CHECK (budget_amount >= 0 AND allocated_amount >= 0)
);

CREATE INDEX IF NOT EXISTS compensation_budget_pools_cycle_idx
    ON hrms.compensation_budget_pools (tenant_id, cycle_id, pool_type)
    WHERE inactive = FALSE;

CREATE TABLE IF NOT EXISTS hrms.compensation_recommendations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    cycle_id UUID NOT NULL REFERENCES hrms.compensation_cycles(id),
    worker_profile_id UUID NOT NULL REFERENCES hrms.worker_profiles(id),
    pay_band_id UUID REFERENCES hrms.compensation_pay_bands(id),
    budget_pool_id UUID REFERENCES hrms.compensation_budget_pools(id),
    current_salary NUMERIC(14,2) NOT NULL DEFAULT 0,
    current_compa_ratio NUMERIC(8,4) NOT NULL DEFAULT 0,
    recommended_salary NUMERIC(14,2) NOT NULL DEFAULT 0,
    recommended_increment_amount NUMERIC(14,2) NOT NULL DEFAULT 0,
    recommended_increment_percent NUMERIC(8,4) NOT NULL DEFAULT 0,
    promotion_recommended BOOLEAN NOT NULL DEFAULT FALSE,
    recommended_designation_id UUID REFERENCES hrms.designations(id),
    reason TEXT,
    performance_rating TEXT,
    equity_flag BOOLEAN NOT NULL DEFAULT FALSE,
    equity_notes TEXT,
    status TEXT NOT NULL DEFAULT 'draft',
    effective_date DATE,
    payroll_handoff_at TIMESTAMPTZ,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    CONSTRAINT compensation_recommendations_status_chk CHECK (status IN ('draft','submitted','approved','rejected','finalized','handed_to_payroll')),
    CONSTRAINT compensation_recommendations_amount_chk CHECK (current_salary >= 0 AND recommended_salary >= 0 AND recommended_increment_amount >= 0),
    CONSTRAINT compensation_recommendations_percent_chk CHECK (recommended_increment_percent >= 0)
);

CREATE UNIQUE INDEX IF NOT EXISTS compensation_recommendations_cycle_worker_uq
    ON hrms.compensation_recommendations (tenant_id, cycle_id, worker_profile_id)
    WHERE inactive = FALSE;

CREATE INDEX IF NOT EXISTS compensation_recommendations_cycle_status_idx
    ON hrms.compensation_recommendations (tenant_id, cycle_id, status)
    WHERE inactive = FALSE;

CREATE TABLE IF NOT EXISTS hrms.compensation_equity_checks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    cycle_id UUID NOT NULL REFERENCES hrms.compensation_cycles(id),
    worker_profile_id UUID NOT NULL REFERENCES hrms.worker_profiles(id),
    pay_band_id UUID REFERENCES hrms.compensation_pay_bands(id),
    check_type TEXT NOT NULL DEFAULT 'band_position',
    severity TEXT NOT NULL DEFAULT 'medium',
    current_salary NUMERIC(14,2) NOT NULL DEFAULT 0,
    band_min NUMERIC(14,2) NOT NULL DEFAULT 0,
    band_midpoint NUMERIC(14,2) NOT NULL DEFAULT 0,
    band_max NUMERIC(14,2) NOT NULL DEFAULT 0,
    variance_percent NUMERIC(8,4) NOT NULL DEFAULT 0,
    finding TEXT NOT NULL,
    recommendation TEXT,
    status TEXT NOT NULL DEFAULT 'open',
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    CONSTRAINT compensation_equity_checks_type_chk CHECK (check_type IN ('below_band','above_band','low_compa_ratio','high_compa_ratio','manual','band_position')),
    CONSTRAINT compensation_equity_checks_severity_chk CHECK (severity IN ('low','medium','high','critical')),
    CONSTRAINT compensation_equity_checks_status_chk CHECK (status IN ('open','acknowledged','resolved','waived'))
);

CREATE INDEX IF NOT EXISTS compensation_equity_checks_cycle_status_idx
    ON hrms.compensation_equity_checks (tenant_id, cycle_id, status, severity)
    WHERE inactive = FALSE;

CREATE TABLE IF NOT EXISTS hrms.compensation_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    cycle_id UUID REFERENCES hrms.compensation_cycles(id),
    source_type TEXT NOT NULL,
    source_id UUID,
    action TEXT NOT NULL,
    from_status TEXT,
    to_status TEXT,
    remarks TEXT,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID
);

CREATE INDEX IF NOT EXISTS compensation_events_source_idx
    ON hrms.compensation_events (tenant_id, source_type, source_id, created_at DESC);

DROP TRIGGER IF EXISTS trg_compensation_pay_bands_updated_at ON hrms.compensation_pay_bands;
CREATE TRIGGER trg_compensation_pay_bands_updated_at BEFORE UPDATE ON hrms.compensation_pay_bands FOR EACH ROW EXECUTE FUNCTION hrms.set_updated_at();

DROP TRIGGER IF EXISTS trg_compensation_cycles_updated_at ON hrms.compensation_cycles;
CREATE TRIGGER trg_compensation_cycles_updated_at BEFORE UPDATE ON hrms.compensation_cycles FOR EACH ROW EXECUTE FUNCTION hrms.set_updated_at();

DROP TRIGGER IF EXISTS trg_compensation_budget_pools_updated_at ON hrms.compensation_budget_pools;
CREATE TRIGGER trg_compensation_budget_pools_updated_at BEFORE UPDATE ON hrms.compensation_budget_pools FOR EACH ROW EXECUTE FUNCTION hrms.set_updated_at();

DROP TRIGGER IF EXISTS trg_compensation_recommendations_updated_at ON hrms.compensation_recommendations;
CREATE TRIGGER trg_compensation_recommendations_updated_at BEFORE UPDATE ON hrms.compensation_recommendations FOR EACH ROW EXECUTE FUNCTION hrms.set_updated_at();

DROP TRIGGER IF EXISTS trg_compensation_equity_checks_updated_at ON hrms.compensation_equity_checks;
CREATE TRIGGER trg_compensation_equity_checks_updated_at BEFORE UPDATE ON hrms.compensation_equity_checks FOR EACH ROW EXECUTE FUNCTION hrms.set_updated_at();

ALTER TABLE hrms.compensation_pay_bands ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.compensation_cycles ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.compensation_budget_pools ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.compensation_recommendations ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.compensation_equity_checks ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.compensation_events ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS compensation_pay_bands_tenant_isolation ON hrms.compensation_pay_bands;
CREATE POLICY compensation_pay_bands_tenant_isolation ON hrms.compensation_pay_bands
USING (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true')
WITH CHECK (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true');

DROP POLICY IF EXISTS compensation_cycles_tenant_isolation ON hrms.compensation_cycles;
CREATE POLICY compensation_cycles_tenant_isolation ON hrms.compensation_cycles
USING (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true')
WITH CHECK (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true');

DROP POLICY IF EXISTS compensation_budget_pools_tenant_isolation ON hrms.compensation_budget_pools;
CREATE POLICY compensation_budget_pools_tenant_isolation ON hrms.compensation_budget_pools
USING (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true')
WITH CHECK (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true');

DROP POLICY IF EXISTS compensation_recommendations_tenant_isolation ON hrms.compensation_recommendations;
CREATE POLICY compensation_recommendations_tenant_isolation ON hrms.compensation_recommendations
USING (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true')
WITH CHECK (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true');

DROP POLICY IF EXISTS compensation_equity_checks_tenant_isolation ON hrms.compensation_equity_checks;
CREATE POLICY compensation_equity_checks_tenant_isolation ON hrms.compensation_equity_checks
USING (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true')
WITH CHECK (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true');

DROP POLICY IF EXISTS compensation_events_tenant_isolation ON hrms.compensation_events;
CREATE POLICY compensation_events_tenant_isolation ON hrms.compensation_events
USING (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true')
WITH CHECK (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true');
