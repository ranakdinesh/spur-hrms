CREATE TABLE IF NOT EXISTS hrms.pay_groups (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id           UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    code                VARCHAR(80) NOT NULL,
    name                VARCHAR(180) NOT NULL,
    description         TEXT,
    grouping_type       VARCHAR(40) NOT NULL DEFAULT 'manual',
    branch_id           UUID REFERENCES hrms.branches(id) ON DELETE SET NULL,
    department_id       UUID REFERENCES hrms.departments(id) ON DELETE SET NULL,
    employment_type_id  UUID REFERENCES hrms.employment_types(id) ON DELETE SET NULL,
    reporting_tag       VARCHAR(120),
    rules               JSONB NOT NULL DEFAULT '{}'::jsonb,
    is_active           BOOLEAN NOT NULL DEFAULT TRUE,
    inactive            BOOLEAN NOT NULL DEFAULT FALSE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by          UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by          UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT pay_groups_code_check CHECK (length(trim(code)) > 0),
    CONSTRAINT pay_groups_name_check CHECK (length(trim(name)) > 0),
    CONSTRAINT pay_groups_type_check CHECK (grouping_type IN ('all','branch','department','employment_type','reporting_tag','manual','mixed'))
);

CREATE UNIQUE INDEX IF NOT EXISTS pay_groups_tenant_code_idx
    ON hrms.pay_groups(tenant_id, lower(code))
    WHERE NOT inactive;

CREATE INDEX IF NOT EXISTS pay_groups_tenant_active_idx
    ON hrms.pay_groups(tenant_id, is_active, name)
    WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.pay_group_members (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id           UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    pay_group_id        UUID NOT NULL REFERENCES hrms.pay_groups(id) ON DELETE CASCADE,
    user_id             UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    membership_type     VARCHAR(30) NOT NULL DEFAULT 'manual_include',
    effective_from      DATE,
    effective_to        DATE,
    inactive            BOOLEAN NOT NULL DEFAULT FALSE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by          UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by          UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT pay_group_members_type_check CHECK (membership_type IN ('manual_include','manual_exclude')),
    CONSTRAINT pay_group_members_date_check CHECK (effective_to IS NULL OR effective_from IS NULL OR effective_to >= effective_from)
);

CREATE UNIQUE INDEX IF NOT EXISTS pay_group_members_unique_idx
    ON hrms.pay_group_members(tenant_id, pay_group_id, user_id, membership_type)
    WHERE NOT inactive;

CREATE INDEX IF NOT EXISTS pay_group_members_group_idx
    ON hrms.pay_group_members(tenant_id, pay_group_id, membership_type)
    WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.pay_runs (
    id                       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id                UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    pay_group_id             UUID NOT NULL REFERENCES hrms.pay_groups(id) ON DELETE RESTRICT,
    fy_id                    UUID NOT NULL REFERENCES hrms.financial_years(id) ON DELETE RESTRICT,
    month                    INTEGER NOT NULL,
    year                     INTEGER NOT NULL,
    status                   VARCHAR(40) NOT NULL DEFAULT 'draft',
    employee_count           INTEGER NOT NULL DEFAULT 0,
    ready_count              INTEGER NOT NULL DEFAULT 0,
    blocked_count            INTEGER NOT NULL DEFAULT 0,
    generated_count          INTEGER NOT NULL DEFAULT 0,
    attendance_frozen_at     TIMESTAMPTZ,
    lop_frozen_at            TIMESTAMPTZ,
    adjustments_frozen_at    TIMESTAMPTZ,
    generated_at             TIMESTAMPTZ,
    locked_at                TIMESTAMPTZ,
    locked_by                UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    unlocked_at              TIMESTAMPTZ,
    unlocked_by              UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    readiness                JSONB NOT NULL DEFAULT '{}'::jsonb,
    notes                    TEXT,
    inactive                 BOOLEAN NOT NULL DEFAULT FALSE,
    created_at               TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by               UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at               TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by               UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT pay_runs_period_check CHECK (month BETWEEN 1 AND 12 AND year BETWEEN 1900 AND 9999),
    CONSTRAINT pay_runs_status_check CHECK (status IN ('draft','readiness_ready','blocked','frozen','processing','generated','locked','unlocked','failed'))
);

CREATE UNIQUE INDEX IF NOT EXISTS pay_runs_group_period_idx
    ON hrms.pay_runs(tenant_id, pay_group_id, month, year)
    WHERE NOT inactive;

CREATE INDEX IF NOT EXISTS pay_runs_tenant_period_idx
    ON hrms.pay_runs(tenant_id, year DESC, month DESC, status)
    WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.pay_run_employees (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id           UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    pay_run_id          UUID NOT NULL REFERENCES hrms.pay_runs(id) ON DELETE CASCADE,
    user_id             UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    readiness_status    VARCHAR(40) NOT NULL DEFAULT 'pending',
    blocker_reason      TEXT,
    salary_slip_id      UUID REFERENCES hrms.salary_slips(id) ON DELETE SET NULL,
    generated_at        TIMESTAMPTZ,
    inactive            BOOLEAN NOT NULL DEFAULT FALSE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by          UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by          UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT pay_run_employees_status_check CHECK (readiness_status IN ('pending','ready','blocked','generated','skipped','failed'))
);

CREATE UNIQUE INDEX IF NOT EXISTS pay_run_employees_unique_idx
    ON hrms.pay_run_employees(tenant_id, pay_run_id, user_id)
    WHERE NOT inactive;

CREATE INDEX IF NOT EXISTS pay_run_employees_run_idx
    ON hrms.pay_run_employees(tenant_id, pay_run_id, readiness_status)
    WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.pay_run_events (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id           UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    pay_run_id          UUID NOT NULL REFERENCES hrms.pay_runs(id) ON DELETE CASCADE,
    action              VARCHAR(60) NOT NULL,
    from_status         VARCHAR(40),
    to_status           VARCHAR(40),
    remarks             TEXT,
    metadata            JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive            BOOLEAN NOT NULL DEFAULT FALSE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by          UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by          UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT pay_run_events_action_check CHECK (length(trim(action)) > 0)
);

CREATE INDEX IF NOT EXISTS pay_run_events_run_idx
    ON hrms.pay_run_events(tenant_id, pay_run_id, created_at DESC)
    WHERE NOT inactive;

ALTER TABLE hrms.pay_groups ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.pay_group_members ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.pay_runs ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.pay_run_employees ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.pay_run_events ENABLE ROW LEVEL SECURITY;

DO $$
DECLARE
    table_name text;
BEGIN
    FOREACH table_name IN ARRAY ARRAY[
        'pay_groups',
        'pay_group_members',
        'pay_runs',
        'pay_run_employees',
        'pay_run_events'
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
