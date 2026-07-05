CREATE TABLE IF NOT EXISTS hrms.worker_types (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    code VARCHAR(80) NOT NULL,
    name VARCHAR(160) NOT NULL,
    classification_group VARCHAR(40) NOT NULL DEFAULT 'employee',
    description TEXT,
    attendance_mode VARCHAR(40) NOT NULL DEFAULT 'checkin_checkout',
    pay_mode VARCHAR(40) NOT NULL DEFAULT 'monthly_salary',
    tds_section VARCHAR(20) NOT NULL DEFAULT 'none',
    pf_applicable BOOLEAN NOT NULL DEFAULT FALSE,
    esic_applicable BOOLEAN NOT NULL DEFAULT FALSE,
    pt_applicable BOOLEAN NOT NULL DEFAULT FALSE,
    lwf_applicable BOOLEAN NOT NULL DEFAULT FALSE,
    clra_applicable BOOLEAN NOT NULL DEFAULT FALSE,
    leave_applicable BOOLEAN NOT NULL DEFAULT FALSE,
    overtime_applicable BOOLEAN NOT NULL DEFAULT FALSE,
    requires_agreement BOOLEAN NOT NULL DEFAULT FALSE,
    requires_invoice BOOLEAN NOT NULL DEFAULT FALSE,
    requires_attendance BOOLEAN NOT NULL DEFAULT TRUE,
    statutory_defaults JSONB NOT NULL DEFAULT '{}'::jsonb,
    compliance_notes TEXT,
    is_system_default BOOLEAN NOT NULL DEFAULT FALSE,
    sort_order INT NOT NULL DEFAULT 100,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT worker_types_code_check CHECK (code IN (
        'permanent_fulltime',
        'permanent_parttime',
        'fixed_term_contract',
        'project_based',
        'freelancer_gig',
        'intern',
        'consultant_retainer',
        'agency_staff'
    )),
    CONSTRAINT worker_types_classification_group_check CHECK (classification_group IN ('employee','contractor','trainee','agency')),
    CONSTRAINT worker_types_attendance_mode_check CHECK (attendance_mode IN ('checkin_checkout','hours_logged','milestone_only','none')),
    CONSTRAINT worker_types_pay_mode_check CHECK (pay_mode IN ('monthly_salary','hourly','project_milestone','invoice','retainer','stipend')),
    CONSTRAINT worker_types_tds_section_check CHECK (tds_section IN ('192','194C','194J','194I','none')),
    CONSTRAINT worker_types_defaults_object_check CHECK (jsonb_typeof(statutory_defaults) = 'object')
);

CREATE UNIQUE INDEX IF NOT EXISTS worker_types_tenant_code_idx
    ON hrms.worker_types (tenant_id, code)
    WHERE NOT inactive;
CREATE INDEX IF NOT EXISTS worker_types_tenant_group_idx
    ON hrms.worker_types (tenant_id, classification_group, sort_order)
    WHERE NOT inactive;
CREATE INDEX IF NOT EXISTS worker_types_tenant_pay_attendance_idx
    ON hrms.worker_types (tenant_id, pay_mode, attendance_mode)
    WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.worker_classification_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    worker_type_id UUID NOT NULL REFERENCES hrms.worker_types(id) ON DELETE CASCADE,
    rule_name VARCHAR(180) NOT NULL,
    rule_type VARCHAR(40) NOT NULL DEFAULT 'manual_guidance',
    priority INT NOT NULL DEFAULT 100,
    conditions JSONB NOT NULL DEFAULT '{}'::jsonb,
    outcome JSONB NOT NULL DEFAULT '{}'::jsonb,
    notes TEXT,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT worker_classification_rules_type_check CHECK (rule_type IN ('manual_guidance','compliance','payroll','attendance')),
    CONSTRAINT worker_classification_rules_priority_check CHECK (priority BETWEEN 1 AND 9999),
    CONSTRAINT worker_classification_rules_conditions_object_check CHECK (jsonb_typeof(conditions) = 'object'),
    CONSTRAINT worker_classification_rules_outcome_object_check CHECK (jsonb_typeof(outcome) = 'object')
);

CREATE INDEX IF NOT EXISTS worker_classification_rules_type_idx
    ON hrms.worker_classification_rules (tenant_id, worker_type_id, priority)
    WHERE NOT inactive;
CREATE UNIQUE INDEX IF NOT EXISTS worker_classification_rules_name_idx
    ON hrms.worker_classification_rules (tenant_id, worker_type_id, lower(rule_name))
    WHERE NOT inactive;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'worker_types_updated_at') THEN
        CREATE TRIGGER worker_types_updated_at
        BEFORE UPDATE ON hrms.worker_types
        FOR EACH ROW EXECUTE FUNCTION hrms.update_updated_at();
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'worker_classification_rules_updated_at') THEN
        CREATE TRIGGER worker_classification_rules_updated_at
        BEFORE UPDATE ON hrms.worker_classification_rules
        FOR EACH ROW EXECUTE FUNCTION hrms.update_updated_at();
    END IF;
END $$;

ALTER TABLE hrms.worker_types ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.worker_classification_rules ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS worker_types_tenant_isolation ON hrms.worker_types;
CREATE POLICY worker_types_tenant_isolation ON hrms.worker_types
USING (tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::uuid OR current_setting('app.is_super_admin', true) = 'true')
WITH CHECK (tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::uuid OR current_setting('app.is_super_admin', true) = 'true');

DROP POLICY IF EXISTS worker_classification_rules_tenant_isolation ON hrms.worker_classification_rules;
CREATE POLICY worker_classification_rules_tenant_isolation ON hrms.worker_classification_rules
USING (tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::uuid OR current_setting('app.is_super_admin', true) = 'true')
WITH CHECK (tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::uuid OR current_setting('app.is_super_admin', true) = 'true');
