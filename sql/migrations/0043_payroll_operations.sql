CREATE TABLE IF NOT EXISTS hrms.payroll_period_locks (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id           UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    month               INTEGER NOT NULL,
    year                INTEGER NOT NULL,
    status              VARCHAR(30) NOT NULL DEFAULT 'open',
    locked_at           TIMESTAMPTZ,
    locked_by           UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    unlocked_at         TIMESTAMPTZ,
    unlocked_by         UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    unlock_reason       TEXT,
    notes               TEXT,
    inactive            BOOLEAN NOT NULL DEFAULT FALSE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by          UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by          UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT payroll_period_locks_period_check CHECK (month BETWEEN 1 AND 12 AND year BETWEEN 1900 AND 9999),
    CONSTRAINT payroll_period_locks_status_check CHECK (status IN ('open','locked','unlocked'))
);

CREATE UNIQUE INDEX IF NOT EXISTS payroll_period_locks_tenant_period_idx
    ON hrms.payroll_period_locks(tenant_id, month, year)
    WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.payroll_period_lock_events (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id           UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    payroll_lock_id     UUID NOT NULL REFERENCES hrms.payroll_period_locks(id) ON DELETE CASCADE,
    action              VARCHAR(50) NOT NULL,
    from_status         VARCHAR(30),
    to_status           VARCHAR(30),
    remarks             TEXT,
    inactive            BOOLEAN NOT NULL DEFAULT FALSE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by          UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by          UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT payroll_period_lock_events_action_check CHECK (length(trim(action)) > 0)
);

CREATE INDEX IF NOT EXISTS payroll_period_lock_events_lock_idx
    ON hrms.payroll_period_lock_events(tenant_id, payroll_lock_id, created_at DESC)
    WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.payroll_statutory_rules (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id           UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    rule_type           VARCHAR(20) NOT NULL,
    name                VARCHAR(160) NOT NULL,
    state               VARCHAR(80),
    branch_id           UUID REFERENCES hrms.branches(id) ON DELETE SET NULL,
    effective_from      DATE NOT NULL,
    effective_to        DATE,
    min_gross_salary    NUMERIC(12,2),
    max_gross_salary    NUMERIC(12,2),
    employee_amount     NUMERIC(12,2) NOT NULL DEFAULT 0,
    employer_amount     NUMERIC(12,2) NOT NULL DEFAULT 0,
    frequency           VARCHAR(30) NOT NULL DEFAULT 'monthly',
    deduction_month     INTEGER,
    notes               TEXT,
    inactive            BOOLEAN NOT NULL DEFAULT FALSE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by          UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by          UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT payroll_statutory_rules_type_check CHECK (rule_type IN ('pt','lwf')),
    CONSTRAINT payroll_statutory_rules_frequency_check CHECK (frequency IN ('monthly','quarterly','half_yearly','annual','one_time')),
    CONSTRAINT payroll_statutory_rules_date_check CHECK (effective_to IS NULL OR effective_to >= effective_from),
    CONSTRAINT payroll_statutory_rules_amount_check CHECK (
        employee_amount >= 0 AND employer_amount >= 0
        AND (min_gross_salary IS NULL OR min_gross_salary >= 0)
        AND (max_gross_salary IS NULL OR max_gross_salary >= 0)
        AND (max_gross_salary IS NULL OR min_gross_salary IS NULL OR max_gross_salary >= min_gross_salary)
        AND (deduction_month IS NULL OR deduction_month BETWEEN 1 AND 12)
    )
);

CREATE INDEX IF NOT EXISTS payroll_statutory_rules_lookup_idx
    ON hrms.payroll_statutory_rules(tenant_id, rule_type, state, branch_id, effective_from)
    WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.payroll_import_batches (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id           UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    import_type         VARCHAR(40) NOT NULL DEFAULT 'salary_revision',
    month               INTEGER,
    year                INTEGER,
    fy_id               UUID REFERENCES hrms.financial_years(id) ON DELETE SET NULL,
    template_id         UUID REFERENCES hrms.salary_templates(id) ON DELETE SET NULL,
    file_name           TEXT,
    status              VARCHAR(30) NOT NULL DEFAULT 'validated',
    total_rows          INTEGER NOT NULL DEFAULT 0,
    valid_rows          INTEGER NOT NULL DEFAULT 0,
    invalid_rows        INTEGER NOT NULL DEFAULT 0,
    applied_rows        INTEGER NOT NULL DEFAULT 0,
    error_report        JSONB NOT NULL DEFAULT '[]'::jsonb,
    notes               TEXT,
    inactive            BOOLEAN NOT NULL DEFAULT FALSE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by          UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by          UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT payroll_import_batches_status_check CHECK (status IN ('validated','applied','failed','partial')),
    CONSTRAINT payroll_import_batches_type_check CHECK (import_type IN ('salary_revision','attendance_lop','variable_pay','adjustment')),
    CONSTRAINT payroll_import_batches_period_check CHECK ((month IS NULL OR month BETWEEN 1 AND 12) AND (year IS NULL OR year BETWEEN 1900 AND 9999))
);

CREATE INDEX IF NOT EXISTS payroll_import_batches_tenant_idx
    ON hrms.payroll_import_batches(tenant_id, created_at DESC)
    WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.payroll_import_rows (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id           UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    batch_id            UUID NOT NULL REFERENCES hrms.payroll_import_batches(id) ON DELETE CASCADE,
    row_number          INTEGER NOT NULL,
    employee_code       VARCHAR(80),
    employee_user_id    UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    employee_name       TEXT,
    gross_salary        NUMERIC(12,2),
    present_days        NUMERIC(8,2),
    absent_days         NUMERIC(8,2),
    lop_days            NUMERIC(8,2),
    variable_earnings   NUMERIC(12,2),
    variable_deductions NUMERIC(12,2),
    status              VARCHAR(30) NOT NULL DEFAULT 'valid',
    error_message       TEXT,
    raw_data            JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive            BOOLEAN NOT NULL DEFAULT FALSE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by          UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by          UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT payroll_import_rows_status_check CHECK (status IN ('valid','invalid','applied','skipped'))
);

CREATE INDEX IF NOT EXISTS payroll_import_rows_batch_idx
    ON hrms.payroll_import_rows(tenant_id, batch_id, row_number)
    WHERE NOT inactive;

ALTER TABLE hrms.payroll_period_locks ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.payroll_period_lock_events ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.payroll_statutory_rules ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.payroll_import_batches ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.payroll_import_rows ENABLE ROW LEVEL SECURITY;

DO $$
DECLARE
    table_name text;
BEGIN
    FOREACH table_name IN ARRAY ARRAY[
        'payroll_period_locks',
        'payroll_period_lock_events',
        'payroll_statutory_rules',
        'payroll_import_batches',
        'payroll_import_rows'
    ]
    LOOP
        EXECUTE format('DROP POLICY IF EXISTS tenant_isolation_%I ON hrms.%I', table_name, table_name);
        EXECUTE format(
            'CREATE POLICY tenant_isolation_%I ON hrms.%I USING (current_setting(''app.is_super_admin'', true) = ''true'' OR tenant_id::text = current_setting(''app.tenant_id'', true)) WITH CHECK (current_setting(''app.is_super_admin'', true) = ''true'' OR tenant_id::text = current_setting(''app.tenant_id'', true))',
            table_name,
            table_name
        );
    END LOOP;
END $$;
