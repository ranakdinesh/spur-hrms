CREATE TABLE IF NOT EXISTS hrms.benefit_plans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    code TEXT NOT NULL,
    name TEXT NOT NULL,
    plan_type TEXT NOT NULL,
    description TEXT,
    provider_name TEXT,
    policy_number TEXT,
    coverage_amount NUMERIC(12,2),
    employer_contribution NUMERIC(12,2) NOT NULL DEFAULT 0,
    employee_contribution NUMERIC(12,2) NOT NULL DEFAULT 0,
    currency_code TEXT NOT NULL DEFAULT 'INR',
    eligibility_rule JSONB NOT NULL DEFAULT '{}'::jsonb,
    insurance_metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    effective_from DATE,
    effective_to DATE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    CONSTRAINT benefit_plans_type_check CHECK (plan_type IN ('insurance','reimbursement','allowance','retirement','wellness','other')),
    CONSTRAINT benefit_plans_amount_check CHECK (coverage_amount IS NULL OR coverage_amount >= 0),
    CONSTRAINT benefit_plans_contribution_check CHECK (employer_contribution >= 0 AND employee_contribution >= 0)
);

CREATE UNIQUE INDEX IF NOT EXISTS benefit_plans_tenant_code_idx ON hrms.benefit_plans(tenant_id, lower(code)) WHERE NOT inactive;
CREATE INDEX IF NOT EXISTS benefit_plans_tenant_active_idx ON hrms.benefit_plans(tenant_id, is_active) WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.benefit_enrollment_windows (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    plan_id UUID NOT NULL REFERENCES hrms.benefit_plans(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    opens_on DATE NOT NULL,
    closes_on DATE NOT NULL,
    effective_from DATE,
    effective_to DATE,
    status TEXT NOT NULL DEFAULT 'open',
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    CONSTRAINT benefit_enrollment_windows_date_check CHECK (closes_on >= opens_on),
    CONSTRAINT benefit_enrollment_windows_status_check CHECK (status IN ('draft','open','closed','archived'))
);

CREATE INDEX IF NOT EXISTS benefit_enrollment_windows_tenant_plan_idx ON hrms.benefit_enrollment_windows(tenant_id, plan_id, status) WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.benefit_dependents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    employee_user_id UUID NOT NULL,
    full_name TEXT NOT NULL,
    relationship TEXT NOT NULL,
    date_of_birth DATE,
    gender TEXT,
    nominee_percentage NUMERIC(5,2),
    is_nominee BOOLEAN NOT NULL DEFAULT FALSE,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    CONSTRAINT benefit_dependents_nominee_percentage_check CHECK (nominee_percentage IS NULL OR (nominee_percentage >= 0 AND nominee_percentage <= 100))
);

CREATE INDEX IF NOT EXISTS benefit_dependents_tenant_user_idx ON hrms.benefit_dependents(tenant_id, employee_user_id) WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.benefit_enrollments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    plan_id UUID NOT NULL REFERENCES hrms.benefit_plans(id) ON DELETE RESTRICT,
    window_id UUID REFERENCES hrms.benefit_enrollment_windows(id) ON DELETE SET NULL,
    employee_user_id UUID NOT NULL,
    status TEXT NOT NULL DEFAULT 'submitted',
    coverage_level TEXT,
    selected_amount NUMERIC(12,2),
    employee_contribution NUMERIC(12,2) NOT NULL DEFAULT 0,
    employer_contribution NUMERIC(12,2) NOT NULL DEFAULT 0,
    effective_from DATE,
    effective_to DATE,
    submitted_at TIMESTAMPTZ,
    reviewed_by UUID,
    reviewed_at TIMESTAMPTZ,
    review_remarks TEXT,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    CONSTRAINT benefit_enrollments_status_check CHECK (status IN ('draft','submitted','approved','rejected','cancelled','active','ended')),
    CONSTRAINT benefit_enrollments_amount_check CHECK (selected_amount IS NULL OR selected_amount >= 0),
    CONSTRAINT benefit_enrollments_contribution_check CHECK (employee_contribution >= 0 AND employer_contribution >= 0)
);

CREATE INDEX IF NOT EXISTS benefit_enrollments_tenant_user_idx ON hrms.benefit_enrollments(tenant_id, employee_user_id, status) WHERE NOT inactive;
CREATE INDEX IF NOT EXISTS benefit_enrollments_tenant_plan_idx ON hrms.benefit_enrollments(tenant_id, plan_id, status) WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.benefit_claim_types (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    plan_id UUID REFERENCES hrms.benefit_plans(id) ON DELETE SET NULL,
    code TEXT NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    annual_limit NUMERIC(12,2),
    per_claim_limit NUMERIC(12,2),
    requires_attachment BOOLEAN NOT NULL DEFAULT TRUE,
    taxable BOOLEAN NOT NULL DEFAULT FALSE,
    payroll_component_code TEXT,
    eligibility_rule JSONB NOT NULL DEFAULT '{}'::jsonb,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    CONSTRAINT benefit_claim_types_limit_check CHECK ((annual_limit IS NULL OR annual_limit >= 0) AND (per_claim_limit IS NULL OR per_claim_limit >= 0))
);

CREATE UNIQUE INDEX IF NOT EXISTS benefit_claim_types_tenant_code_idx ON hrms.benefit_claim_types(tenant_id, lower(code)) WHERE NOT inactive;
CREATE INDEX IF NOT EXISTS benefit_claim_types_tenant_plan_idx ON hrms.benefit_claim_types(tenant_id, plan_id, is_active) WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.benefit_claims (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    claim_number TEXT NOT NULL,
    claim_type_id UUID NOT NULL REFERENCES hrms.benefit_claim_types(id) ON DELETE RESTRICT,
    plan_id UUID REFERENCES hrms.benefit_plans(id) ON DELETE SET NULL,
    employee_user_id UUID NOT NULL,
    dependent_id UUID REFERENCES hrms.benefit_dependents(id) ON DELETE SET NULL,
    expense_date DATE NOT NULL,
    submitted_at TIMESTAMPTZ,
    claim_amount NUMERIC(12,2) NOT NULL,
    approved_amount NUMERIC(12,2),
    currency_code TEXT NOT NULL DEFAULT 'INR',
    status TEXT NOT NULL DEFAULT 'draft',
    payment_status TEXT NOT NULL DEFAULT 'not_payable',
    payment_reference TEXT,
    paid_at TIMESTAMPTZ,
    reviewed_by UUID,
    reviewed_at TIMESTAMPTZ,
    review_remarks TEXT,
    payroll_export_status TEXT NOT NULL DEFAULT 'not_ready',
    payroll_exported_at TIMESTAMPTZ,
    payroll_export_reference TEXT,
    notes TEXT,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    CONSTRAINT benefit_claims_amount_check CHECK (claim_amount > 0 AND (approved_amount IS NULL OR approved_amount >= 0)),
    CONSTRAINT benefit_claims_status_check CHECK (status IN ('draft','submitted','under_review','approved','rejected','cancelled','paid')),
    CONSTRAINT benefit_claims_payment_status_check CHECK (payment_status IN ('not_payable','pending','paid','failed')),
    CONSTRAINT benefit_claims_payroll_export_status_check CHECK (payroll_export_status IN ('not_ready','ready','exported','blocked'))
);

CREATE UNIQUE INDEX IF NOT EXISTS benefit_claims_tenant_number_idx ON hrms.benefit_claims(tenant_id, claim_number) WHERE NOT inactive;
CREATE INDEX IF NOT EXISTS benefit_claims_tenant_user_idx ON hrms.benefit_claims(tenant_id, employee_user_id, status) WHERE NOT inactive;
CREATE INDEX IF NOT EXISTS benefit_claims_tenant_status_idx ON hrms.benefit_claims(tenant_id, status, payment_status, payroll_export_status) WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.benefit_claim_attachments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    claim_id UUID NOT NULL REFERENCES hrms.benefit_claims(id) ON DELETE CASCADE,
    file_name TEXT NOT NULL,
    content_type TEXT NOT NULL,
    storage_path TEXT NOT NULL,
    checksum_sha256 TEXT,
    size_bytes BIGINT NOT NULL DEFAULT 0,
    uploaded_by UUID,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID
);

CREATE INDEX IF NOT EXISTS benefit_claim_attachments_claim_idx ON hrms.benefit_claim_attachments(tenant_id, claim_id) WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.benefit_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    source_type TEXT NOT NULL,
    source_id UUID NOT NULL,
    action TEXT NOT NULL,
    from_status TEXT,
    to_status TEXT,
    actor_user_id UUID,
    remarks TEXT,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    CONSTRAINT benefit_events_source_type_check CHECK (source_type IN ('plan','window','dependent','enrollment','claim_type','claim','attachment'))
);

CREATE INDEX IF NOT EXISTS benefit_events_source_idx ON hrms.benefit_events(tenant_id, source_type, source_id, created_at DESC) WHERE NOT inactive;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'benefit_plans_updated_at') THEN
        CREATE TRIGGER benefit_plans_updated_at BEFORE UPDATE ON hrms.benefit_plans FOR EACH ROW EXECUTE FUNCTION hrms.set_updated_at();
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'benefit_enrollment_windows_updated_at') THEN
        CREATE TRIGGER benefit_enrollment_windows_updated_at BEFORE UPDATE ON hrms.benefit_enrollment_windows FOR EACH ROW EXECUTE FUNCTION hrms.set_updated_at();
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'benefit_dependents_updated_at') THEN
        CREATE TRIGGER benefit_dependents_updated_at BEFORE UPDATE ON hrms.benefit_dependents FOR EACH ROW EXECUTE FUNCTION hrms.set_updated_at();
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'benefit_enrollments_updated_at') THEN
        CREATE TRIGGER benefit_enrollments_updated_at BEFORE UPDATE ON hrms.benefit_enrollments FOR EACH ROW EXECUTE FUNCTION hrms.set_updated_at();
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'benefit_claim_types_updated_at') THEN
        CREATE TRIGGER benefit_claim_types_updated_at BEFORE UPDATE ON hrms.benefit_claim_types FOR EACH ROW EXECUTE FUNCTION hrms.set_updated_at();
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'benefit_claims_updated_at') THEN
        CREATE TRIGGER benefit_claims_updated_at BEFORE UPDATE ON hrms.benefit_claims FOR EACH ROW EXECUTE FUNCTION hrms.set_updated_at();
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'benefit_claim_attachments_updated_at') THEN
        CREATE TRIGGER benefit_claim_attachments_updated_at BEFORE UPDATE ON hrms.benefit_claim_attachments FOR EACH ROW EXECUTE FUNCTION hrms.set_updated_at();
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'benefit_events_updated_at') THEN
        CREATE TRIGGER benefit_events_updated_at BEFORE UPDATE ON hrms.benefit_events FOR EACH ROW EXECUTE FUNCTION hrms.set_updated_at();
    END IF;
END $$;

ALTER TABLE hrms.benefit_plans ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.benefit_enrollment_windows ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.benefit_dependents ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.benefit_enrollments ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.benefit_claim_types ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.benefit_claims ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.benefit_claim_attachments ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.benefit_events ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS benefit_plans_tenant_isolation ON hrms.benefit_plans;
CREATE POLICY benefit_plans_tenant_isolation ON hrms.benefit_plans USING (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true)) WITH CHECK (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true));
DROP POLICY IF EXISTS benefit_enrollment_windows_tenant_isolation ON hrms.benefit_enrollment_windows;
CREATE POLICY benefit_enrollment_windows_tenant_isolation ON hrms.benefit_enrollment_windows USING (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true)) WITH CHECK (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true));
DROP POLICY IF EXISTS benefit_dependents_tenant_isolation ON hrms.benefit_dependents;
CREATE POLICY benefit_dependents_tenant_isolation ON hrms.benefit_dependents USING (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true)) WITH CHECK (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true));
DROP POLICY IF EXISTS benefit_enrollments_tenant_isolation ON hrms.benefit_enrollments;
CREATE POLICY benefit_enrollments_tenant_isolation ON hrms.benefit_enrollments USING (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true)) WITH CHECK (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true));
DROP POLICY IF EXISTS benefit_claim_types_tenant_isolation ON hrms.benefit_claim_types;
CREATE POLICY benefit_claim_types_tenant_isolation ON hrms.benefit_claim_types USING (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true)) WITH CHECK (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true));
DROP POLICY IF EXISTS benefit_claims_tenant_isolation ON hrms.benefit_claims;
CREATE POLICY benefit_claims_tenant_isolation ON hrms.benefit_claims USING (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true)) WITH CHECK (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true));
DROP POLICY IF EXISTS benefit_claim_attachments_tenant_isolation ON hrms.benefit_claim_attachments;
CREATE POLICY benefit_claim_attachments_tenant_isolation ON hrms.benefit_claim_attachments USING (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true)) WITH CHECK (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true));
DROP POLICY IF EXISTS benefit_events_tenant_isolation ON hrms.benefit_events;
CREATE POLICY benefit_events_tenant_isolation ON hrms.benefit_events USING (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true)) WITH CHECK (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true));
