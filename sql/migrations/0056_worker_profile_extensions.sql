CREATE TABLE IF NOT EXISTS hrms.worker_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    worker_type_id UUID NOT NULL REFERENCES hrms.worker_types(id) ON DELETE RESTRICT,
    employee_id UUID REFERENCES hrms.employees(id) ON DELETE SET NULL,
    employee_user_id UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    worker_code VARCHAR(80),
    display_name VARCHAR(180) NOT NULL,
    legal_name VARCHAR(180),
    email VARCHAR(255),
    mobile VARCHAR(30),
    profile_status VARCHAR(40) NOT NULL DEFAULT 'draft',
    start_date DATE,
    end_date DATE,
    branch_id UUID REFERENCES hrms.branches(id) ON DELETE SET NULL,
    department_id UUID REFERENCES hrms.departments(id) ON DELETE SET NULL,
    reporting_manager_id UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    work_location_label VARCHAR(180),
    source_partner VARCHAR(180),
    external_reference VARCHAR(120),
    compliance_status VARCHAR(40) NOT NULL DEFAULT 'pending',
    payroll_status VARCHAR(40) NOT NULL DEFAULT 'not_applicable',
    notes TEXT,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT worker_profiles_status_check CHECK (profile_status IN ('draft','active','paused','ended','blacklisted')),
    CONSTRAINT worker_profiles_compliance_status_check CHECK (compliance_status IN ('pending','ready','review_required','blocked')),
    CONSTRAINT worker_profiles_payroll_status_check CHECK (payroll_status IN ('not_applicable','pending','ready','blocked')),
    CONSTRAINT worker_profiles_dates_check CHECK (end_date IS NULL OR start_date IS NULL OR end_date >= start_date),
    CONSTRAINT worker_profiles_metadata_object_check CHECK (jsonb_typeof(metadata) = 'object'),
    CONSTRAINT worker_profiles_employee_link_check CHECK (
        (employee_id IS NULL AND employee_user_id IS NULL)
        OR (employee_id IS NOT NULL AND employee_user_id IS NOT NULL)
    )
);

CREATE UNIQUE INDEX IF NOT EXISTS worker_profiles_tenant_worker_code_idx
    ON hrms.worker_profiles (tenant_id, lower(worker_code))
    WHERE worker_code IS NOT NULL AND NOT inactive;
CREATE UNIQUE INDEX IF NOT EXISTS worker_profiles_tenant_employee_idx
    ON hrms.worker_profiles (tenant_id, employee_id)
    WHERE employee_id IS NOT NULL AND NOT inactive;
CREATE INDEX IF NOT EXISTS worker_profiles_tenant_status_idx
    ON hrms.worker_profiles (tenant_id, profile_status, updated_at DESC)
    WHERE NOT inactive;
CREATE INDEX IF NOT EXISTS worker_profiles_tenant_type_idx
    ON hrms.worker_profiles (tenant_id, worker_type_id, profile_status)
    WHERE NOT inactive;
CREATE INDEX IF NOT EXISTS worker_profiles_tenant_org_idx
    ON hrms.worker_profiles (tenant_id, branch_id, department_id)
    WHERE NOT inactive;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'worker_profiles_updated_at') THEN
        CREATE TRIGGER worker_profiles_updated_at
        BEFORE UPDATE ON hrms.worker_profiles
        FOR EACH ROW EXECUTE FUNCTION hrms.update_updated_at();
    END IF;
END $$;

ALTER TABLE hrms.worker_profiles ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS worker_profiles_tenant_isolation ON hrms.worker_profiles;
CREATE POLICY worker_profiles_tenant_isolation ON hrms.worker_profiles
USING (tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::uuid OR current_setting('app.is_super_admin', true) = 'true')
WITH CHECK (tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::uuid OR current_setting('app.is_super_admin', true) = 'true');
