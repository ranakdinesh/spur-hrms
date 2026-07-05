CREATE TABLE IF NOT EXISTS hrms.leave_policy_templates (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id      UUID REFERENCES auth.tenants(id) ON DELETE CASCADE,
    name           VARCHAR(120) NOT NULL,
    code           VARCHAR(60) NOT NULL,
    description    TEXT,
    is_system      BOOLEAN NOT NULL DEFAULT FALSE,
    effective_from DATE,
    effective_to   DATE,
    inactive       BOOLEAN NOT NULL DEFAULT FALSE,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by     UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by     UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT leave_policy_templates_scope_check CHECK ((is_system AND tenant_id IS NULL) OR (NOT is_system AND tenant_id IS NOT NULL)),
    CONSTRAINT leave_policy_templates_effective_check CHECK (effective_to IS NULL OR effective_from IS NULL OR effective_to >= effective_from)
);
CREATE UNIQUE INDEX IF NOT EXISTS leave_policy_templates_system_code_idx ON hrms.leave_policy_templates(lower(code)) WHERE is_system AND NOT inactive;
CREATE UNIQUE INDEX IF NOT EXISTS leave_policy_templates_tenant_code_idx ON hrms.leave_policy_templates(tenant_id, lower(code)) WHERE NOT is_system AND NOT inactive;
CREATE INDEX IF NOT EXISTS leave_policy_templates_tenant_idx ON hrms.leave_policy_templates(tenant_id);

CREATE TABLE IF NOT EXISTS hrms.leave_policy_template_rules (
    id                           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id                    UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    template_id                  UUID NOT NULL REFERENCES hrms.leave_policy_templates(id) ON DELETE CASCADE,
    leave_type_id                UUID NOT NULL REFERENCES hrms.leave_types(id) ON DELETE CASCADE,
    fy_id                        UUID REFERENCES hrms.financial_years(id) ON DELETE SET NULL,
    employment_type_id           UUID REFERENCES hrms.employment_types(id) ON DELETE SET NULL,
    department_id                UUID REFERENCES hrms.departments(id) ON DELETE SET NULL,
    designation_id               UUID REFERENCES hrms.designations(id) ON DELETE SET NULL,
    probation_status             VARCHAR(30),
    accrual_method               VARCHAR(40) NOT NULL,
    accrual_frequency            VARCHAR(30) NOT NULL DEFAULT 'monthly',
    credit_days                  NUMERIC(6,2) NOT NULL DEFAULT 0,
    credit_hours                 NUMERIC(8,2) NOT NULL DEFAULT 0,
    min_worked_days              INT NOT NULL DEFAULT 0,
    max_balance                  NUMERIC(6,2),
    carry_forward_enabled        BOOLEAN NOT NULL DEFAULT FALSE,
    max_carry_forward            NUMERIC(6,2) NOT NULL DEFAULT 0,
    carry_forward_expiry_months  INT NOT NULL DEFAULT 0,
    encashment_enabled           BOOLEAN NOT NULL DEFAULT FALSE,
    negative_balance_allowed     BOOLEAN NOT NULL DEFAULT FALSE,
    max_negative_balance         NUMERIC(6,2) NOT NULL DEFAULT 0,
    sandwich_applicable          BOOLEAN NOT NULL DEFAULT FALSE,
    include_holidays             BOOLEAN NOT NULL DEFAULT FALSE,
    include_weekoffs             BOOLEAN NOT NULL DEFAULT FALSE,
    requires_document_after_days NUMERIC(6,2),
    calculation_config           JSONB NOT NULL DEFAULT '{}'::jsonb,
    priority                     INT NOT NULL DEFAULT 100,
    inactive                     BOOLEAN NOT NULL DEFAULT FALSE,
    created_at                   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by                   UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at                   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by                   UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT leave_policy_template_rules_probation_check CHECK (probation_status IS NULL OR probation_status IN ('any','probation','confirmed')),
    CONSTRAINT leave_policy_template_rules_method_check CHECK (accrual_method IN ('fixed_yearly','monthly_fixed','probation_monthly','worked_days','worked_day_range','worked_percentage','tenure_slab','pto_bank','comp_off','manual_adjustment')),
    CONSTRAINT leave_policy_template_rules_frequency_check CHECK (accrual_frequency IN ('instant','daily','weekly','biweekly','monthly','yearly','manual')),
    CONSTRAINT leave_policy_template_rules_nonnegative_check CHECK (credit_days >= 0 AND credit_hours >= 0 AND min_worked_days >= 0 AND max_carry_forward >= 0 AND carry_forward_expiry_months >= 0 AND max_negative_balance >= 0)
);
CREATE INDEX IF NOT EXISTS leave_policy_template_rules_tenant_idx ON hrms.leave_policy_template_rules(tenant_id);
CREATE INDEX IF NOT EXISTS leave_policy_template_rules_template_idx ON hrms.leave_policy_template_rules(template_id);
CREATE INDEX IF NOT EXISTS leave_policy_template_rules_leave_type_idx ON hrms.leave_policy_template_rules(tenant_id, leave_type_id);

CREATE TABLE IF NOT EXISTS hrms.employee_leave_policy_assignments (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    user_id     UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    template_id UUID NOT NULL REFERENCES hrms.leave_policy_templates(id) ON DELETE CASCADE,
    fy_id       UUID REFERENCES hrms.financial_years(id) ON DELETE SET NULL,
    effective_from DATE NOT NULL,
    effective_to   DATE,
    inactive    BOOLEAN NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by  UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by  UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT employee_leave_policy_assignments_effective_check CHECK (effective_to IS NULL OR effective_to >= effective_from)
);
CREATE UNIQUE INDEX IF NOT EXISTS employee_leave_policy_assignment_unique_idx ON hrms.employee_leave_policy_assignments(tenant_id, user_id, template_id, fy_id) WHERE NOT inactive;
CREATE INDEX IF NOT EXISTS employee_leave_policy_assignments_user_idx ON hrms.employee_leave_policy_assignments(tenant_id, user_id);

ALTER TABLE hrms.leave_ledger ADD COLUMN IF NOT EXISTS source_type VARCHAR(40) NOT NULL DEFAULT 'manual_adjustment';
ALTER TABLE hrms.leave_ledger ADD COLUMN IF NOT EXISTS source_id UUID;
ALTER TABLE hrms.leave_ledger ADD COLUMN IF NOT EXISTS balance_before NUMERIC(6,2);
ALTER TABLE hrms.leave_ledger ADD COLUMN IF NOT EXISTS balance_after NUMERIC(6,2);
ALTER TABLE hrms.leave_ledger ADD COLUMN IF NOT EXISTS pending_before NUMERIC(6,2);
ALTER TABLE hrms.leave_ledger ADD COLUMN IF NOT EXISTS pending_after NUMERIC(6,2);
ALTER TABLE hrms.leave_ledger ADD COLUMN IF NOT EXISTS used_before NUMERIC(6,2);
ALTER TABLE hrms.leave_ledger ADD COLUMN IF NOT EXISTS used_after NUMERIC(6,2);
ALTER TABLE hrms.leave_ledger ADD COLUMN IF NOT EXISTS metadata JSONB NOT NULL DEFAULT '{}'::jsonb;
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'leave_ledger_source_type_check'
    ) THEN
        ALTER TABLE hrms.leave_ledger ADD CONSTRAINT leave_ledger_source_type_check CHECK (source_type IN ('opening_balance','monthly_accrual','yearly_accrual','leave_apply','leave_approve','leave_reject','leave_cancel','comp_off','carry_forward','encashment','manual_adjustment'));
    END IF;
END $$;

CREATE TABLE IF NOT EXISTS hrms.comp_off_requests (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id           UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    user_id             UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    leave_type_id       UUID NOT NULL REFERENCES hrms.leave_types(id) ON DELETE RESTRICT,
    fy_id               UUID NOT NULL REFERENCES hrms.financial_years(id) ON DELETE RESTRICT,
    work_date           DATE NOT NULL,
    worked_hours        NUMERIC(8,2) NOT NULL DEFAULT 0,
    credit_days         NUMERIC(6,2) NOT NULL DEFAULT 0,
    multiplier          NUMERIC(5,2) NOT NULL DEFAULT 1,
    reason              TEXT,
    status              VARCHAR(20) NOT NULL DEFAULT 'pending',
    reviewed_by         UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    reviewed_at         TIMESTAMPTZ,
    remarks             TEXT,
    inactive            BOOLEAN NOT NULL DEFAULT FALSE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by          UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by          UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT comp_off_requests_status_check CHECK (status IN ('pending','approved','rejected','canceled')),
    CONSTRAINT comp_off_requests_nonnegative_check CHECK (worked_hours >= 0 AND credit_days >= 0 AND multiplier >= 0)
);
CREATE INDEX IF NOT EXISTS comp_off_requests_user_idx ON hrms.comp_off_requests(tenant_id, user_id);
CREATE INDEX IF NOT EXISTS comp_off_requests_status_idx ON hrms.comp_off_requests(tenant_id, status);

INSERT INTO hrms.leave_policy_templates (tenant_id, name, code, description, is_system)
VALUES
    (NULL, 'Probation Basic', 'PROBATION_BASIC', 'Probation model: one paid leave per month, restricted sick/earned leave until confirmation.', TRUE),
    (NULL, 'Permanent Monthly 2.0', 'PERMANENT_MONTHLY_2', 'Confirmed employee model: two paid leaves credited monthly.', TRUE),
    (NULL, 'Permanent Monthly 2.5', 'PERMANENT_MONTHLY_2_5', 'Confirmed employee model: two and a half paid leaves credited monthly.', TRUE),
    (NULL, 'Yearly Upfront', 'YEARLY_UPFRONT', 'Annual entitlement credited upfront at financial-year start or joining.', TRUE),
    (NULL, 'PTO Bank', 'PTO_BANK', 'Single combined paid-time-off bucket instead of separate leave categories.', TRUE),
    (NULL, 'Worked Day Based', 'WORKED_DAY_BASED', 'Accrual based on payable/worked days in the period.', TRUE),
    (NULL, 'Attendance Percentage Based', 'ATTENDANCE_PERCENTAGE', 'Accrual based on attendance percentage ranges.', TRUE),
    (NULL, 'Tenure Slab Based', 'TENURE_SLAB', 'Accrual changes by completed service years.', TRUE),
    (NULL, 'Comp-Off Holiday Work', 'COMP_OFF_HOLIDAY', 'Approved holiday or overtime work credits comp-off leave.', TRUE),
    (NULL, 'Manual Opening Balance', 'MANUAL_OPENING', 'Admin-managed opening balance and manual adjustments with ledger audit.', TRUE)
ON CONFLICT (lower(code)) WHERE is_system AND NOT inactive
DO UPDATE SET
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    updated_at = NOW();

DO $$
DECLARE
    table_name TEXT;
    updated_tables TEXT[] := ARRAY[
        'leave_policy_templates','leave_policy_template_rules','employee_leave_policy_assignments','comp_off_requests'
    ];
BEGIN
    FOREACH table_name IN ARRAY updated_tables LOOP
        EXECUTE format('DROP TRIGGER IF EXISTS %I ON hrms.%I', table_name || '_updated_at', table_name);
        EXECUTE format(
            'CREATE TRIGGER %I BEFORE UPDATE ON hrms.%I FOR EACH ROW EXECUTE FUNCTION hrms.update_updated_at()',
            table_name || '_updated_at',
            table_name
        );
    END LOOP;
END $$;

ALTER TABLE hrms.leave_policy_templates ENABLE ROW LEVEL SECURITY;
DROP POLICY IF EXISTS leave_policy_templates_tenant_isolation ON hrms.leave_policy_templates;
CREATE POLICY leave_policy_templates_tenant_isolation ON hrms.leave_policy_templates
USING (
    is_system
    OR tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::uuid
    OR current_setting('app.is_super_admin', true) = 'true'
)
WITH CHECK (
    is_system
    OR tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::uuid
    OR current_setting('app.is_super_admin', true) = 'true'
);

DO $$
DECLARE
    table_name TEXT;
    tenant_tables TEXT[] := ARRAY[
        'leave_policy_template_rules','employee_leave_policy_assignments','comp_off_requests'
    ];
BEGIN
    FOREACH table_name IN ARRAY tenant_tables LOOP
        EXECUTE format('ALTER TABLE hrms.%I ENABLE ROW LEVEL SECURITY', table_name);
        EXECUTE format('DROP POLICY IF EXISTS %I ON hrms.%I', table_name || '_tenant_isolation', table_name);
        EXECUTE format(
            'CREATE POLICY %I ON hrms.%I USING (tenant_id = NULLIF(current_setting(''app.tenant_id'', true), '''')::uuid OR current_setting(''app.is_super_admin'', true) = ''true'') WITH CHECK (tenant_id = NULLIF(current_setting(''app.tenant_id'', true), '''')::uuid OR current_setting(''app.is_super_admin'', true) = ''true'')',
            table_name || '_tenant_isolation',
            table_name
        );
    END LOOP;
END $$;
