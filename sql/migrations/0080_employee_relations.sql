CREATE TABLE IF NOT EXISTS hrms.er_case_categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    code TEXT NOT NULL,
    name TEXT NOT NULL,
    case_family TEXT NOT NULL DEFAULT 'grievance',
    description TEXT,
    default_severity TEXT NOT NULL DEFAULT 'medium',
    default_owner_role TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT er_case_categories_family_check CHECK (case_family IN ('grievance','disciplinary','harassment','ethics','workplace_conflict','policy_violation','other')),
    CONSTRAINT er_case_categories_severity_check CHECK (default_severity IN ('low','medium','high','critical')),
    CONSTRAINT er_case_categories_code_check CHECK (length(trim(code)) > 0),
    CONSTRAINT er_case_categories_name_check CHECK (length(trim(name)) > 0)
);

CREATE UNIQUE INDEX IF NOT EXISTS er_case_categories_tenant_code_idx ON hrms.er_case_categories(tenant_id, lower(code)) WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.er_cases (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    case_number TEXT NOT NULL,
    source_hr_case_id UUID REFERENCES hrms.hr_cases(id) ON DELETE SET NULL,
    category_id UUID REFERENCES hrms.er_case_categories(id) ON DELETE SET NULL,
    title TEXT NOT NULL,
    intake_summary TEXT NOT NULL,
    case_family TEXT NOT NULL DEFAULT 'grievance',
    severity TEXT NOT NULL DEFAULT 'medium',
    status TEXT NOT NULL DEFAULT 'intake',
    confidentiality_level TEXT NOT NULL DEFAULT 'restricted',
    complainant_user_id UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    subject_employee_user_id UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    owner_user_id UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    owner_role TEXT,
    investigation_lead_user_id UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    legal_hold BOOLEAN NOT NULL DEFAULT FALSE,
    legal_hold_reason TEXT,
    legal_hold_at TIMESTAMPTZ,
    legal_hold_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    due_at TIMESTAMPTZ,
    closed_at TIMESTAMPTZ,
    resolution_summary TEXT,
    privacy_notes TEXT,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT er_cases_family_check CHECK (case_family IN ('grievance','disciplinary','harassment','ethics','workplace_conflict','policy_violation','other')),
    CONSTRAINT er_cases_severity_check CHECK (severity IN ('low','medium','high','critical')),
    CONSTRAINT er_cases_status_check CHECK (status IN ('intake','triage','investigation','findings','action_plan','monitoring','closed','cancelled')),
    CONSTRAINT er_cases_confidentiality_check CHECK (confidentiality_level IN ('restricted','sensitive','legal_hold')),
    CONSTRAINT er_cases_number_check CHECK (length(trim(case_number)) > 0),
    CONSTRAINT er_cases_title_check CHECK (length(trim(title)) > 0),
    CONSTRAINT er_cases_summary_check CHECK (length(trim(intake_summary)) > 0)
);

CREATE UNIQUE INDEX IF NOT EXISTS er_cases_tenant_number_idx ON hrms.er_cases(tenant_id, case_number) WHERE NOT inactive;
CREATE INDEX IF NOT EXISTS er_cases_tenant_status_idx ON hrms.er_cases(tenant_id, status, severity, due_at) WHERE NOT inactive;
CREATE INDEX IF NOT EXISTS er_cases_owner_idx ON hrms.er_cases(tenant_id, owner_user_id, owner_role, status) WHERE NOT inactive;
CREATE INDEX IF NOT EXISTS er_cases_source_hr_case_idx ON hrms.er_cases(tenant_id, source_hr_case_id) WHERE source_hr_case_id IS NOT NULL AND NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.er_case_parties (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    er_case_id UUID NOT NULL REFERENCES hrms.er_cases(id) ON DELETE CASCADE,
    party_user_id UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    party_name TEXT,
    party_role TEXT NOT NULL,
    representation_notes TEXT,
    contact_notes TEXT,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT er_case_parties_role_check CHECK (party_role IN ('complainant','respondent','witness','investigator','hr_partner','legal','manager','other')),
    CONSTRAINT er_case_parties_identity_check CHECK (party_user_id IS NOT NULL OR length(trim(COALESCE(party_name, ''))) > 0)
);

CREATE INDEX IF NOT EXISTS er_case_parties_case_idx ON hrms.er_case_parties(tenant_id, er_case_id, party_role) WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.er_allegations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    er_case_id UUID NOT NULL REFERENCES hrms.er_cases(id) ON DELETE CASCADE,
    allegation_type TEXT NOT NULL,
    incident_date DATE,
    incident_location TEXT,
    description TEXT NOT NULL,
    policy_reference TEXT,
    status TEXT NOT NULL DEFAULT 'open',
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT er_allegations_status_check CHECK (status IN ('open','substantiated','unsubstantiated','inconclusive','withdrawn')),
    CONSTRAINT er_allegations_description_check CHECK (length(trim(description)) > 0)
);

CREATE INDEX IF NOT EXISTS er_allegations_case_idx ON hrms.er_allegations(tenant_id, er_case_id, status) WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.er_investigation_steps (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    er_case_id UUID NOT NULL REFERENCES hrms.er_cases(id) ON DELETE CASCADE,
    step_type TEXT NOT NULL,
    title TEXT NOT NULL,
    description TEXT,
    owner_user_id UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    due_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    status TEXT NOT NULL DEFAULT 'pending',
    outcome_notes TEXT,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT er_steps_status_check CHECK (status IN ('pending','in_progress','completed','skipped','blocked')),
    CONSTRAINT er_steps_title_check CHECK (length(trim(title)) > 0)
);

CREATE INDEX IF NOT EXISTS er_steps_case_idx ON hrms.er_investigation_steps(tenant_id, er_case_id, status, due_at) WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.er_witness_notes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    er_case_id UUID NOT NULL REFERENCES hrms.er_cases(id) ON DELETE CASCADE,
    witness_user_id UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    witness_name TEXT,
    interview_at TIMESTAMPTZ,
    interviewer_user_id UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    statement_summary TEXT NOT NULL,
    confidentiality_level TEXT NOT NULL DEFAULT 'restricted',
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT er_witness_identity_check CHECK (witness_user_id IS NOT NULL OR length(trim(COALESCE(witness_name, ''))) > 0),
    CONSTRAINT er_witness_summary_check CHECK (length(trim(statement_summary)) > 0),
    CONSTRAINT er_witness_confidentiality_check CHECK (confidentiality_level IN ('restricted','sensitive','legal_hold'))
);

CREATE INDEX IF NOT EXISTS er_witness_notes_case_idx ON hrms.er_witness_notes(tenant_id, er_case_id, created_at DESC) WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.er_evidence_attachments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    er_case_id UUID NOT NULL REFERENCES hrms.er_cases(id) ON DELETE CASCADE,
    allegation_id UUID REFERENCES hrms.er_allegations(id) ON DELETE SET NULL,
    file_name TEXT NOT NULL,
    content_type TEXT NOT NULL,
    storage_path TEXT NOT NULL,
    checksum_sha256 TEXT,
    size_bytes BIGINT NOT NULL DEFAULT 0,
    evidence_type TEXT NOT NULL DEFAULT 'document',
    description TEXT,
    uploaded_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    legal_hold BOOLEAN NOT NULL DEFAULT FALSE,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT er_evidence_file_check CHECK (length(trim(file_name)) > 0 AND length(trim(storage_path)) > 0),
    CONSTRAINT er_evidence_size_check CHECK (size_bytes >= 0)
);

CREATE INDEX IF NOT EXISTS er_evidence_case_idx ON hrms.er_evidence_attachments(tenant_id, er_case_id, created_at DESC) WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.er_findings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    er_case_id UUID NOT NULL REFERENCES hrms.er_cases(id) ON DELETE CASCADE,
    allegation_id UUID REFERENCES hrms.er_allegations(id) ON DELETE SET NULL,
    finding TEXT NOT NULL,
    rationale TEXT NOT NULL,
    recommended_action TEXT,
    decided_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    decided_at TIMESTAMPTZ,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT er_findings_finding_check CHECK (finding IN ('substantiated','unsubstantiated','inconclusive','partially_substantiated','withdrawn')),
    CONSTRAINT er_findings_rationale_check CHECK (length(trim(rationale)) > 0)
);

CREATE INDEX IF NOT EXISTS er_findings_case_idx ON hrms.er_findings(tenant_id, er_case_id, created_at DESC) WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.er_action_plans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    er_case_id UUID NOT NULL REFERENCES hrms.er_cases(id) ON DELETE CASCADE,
    action_type TEXT NOT NULL,
    description TEXT NOT NULL,
    assigned_to_user_id UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    due_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    status TEXT NOT NULL DEFAULT 'pending',
    follow_up_notes TEXT,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT er_action_status_check CHECK (status IN ('pending','in_progress','completed','cancelled','overdue')),
    CONSTRAINT er_action_description_check CHECK (length(trim(description)) > 0)
);

CREATE INDEX IF NOT EXISTS er_action_plans_case_idx ON hrms.er_action_plans(tenant_id, er_case_id, status, due_at) WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.er_case_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    er_case_id UUID NOT NULL REFERENCES hrms.er_cases(id) ON DELETE CASCADE,
    event_type TEXT NOT NULL,
    from_status TEXT,
    to_status TEXT,
    actor_user_id UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    comment TEXT,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT er_case_events_type_check CHECK (length(trim(event_type)) > 0)
);

CREATE INDEX IF NOT EXISTS er_events_case_idx ON hrms.er_case_events(tenant_id, er_case_id, created_at DESC);

CREATE OR REPLACE FUNCTION hrms.prevent_er_event_mutation()
RETURNS trigger AS $$
BEGIN
    RAISE EXCEPTION 'employee relations audit events are immutable';
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS er_case_events_immutable_update ON hrms.er_case_events;
CREATE TRIGGER er_case_events_immutable_update
    BEFORE UPDATE OR DELETE ON hrms.er_case_events
    FOR EACH ROW EXECUTE FUNCTION hrms.prevent_er_event_mutation();

ALTER TABLE hrms.er_case_categories ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.er_cases ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.er_case_parties ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.er_allegations ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.er_investigation_steps ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.er_witness_notes ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.er_evidence_attachments ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.er_findings ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.er_action_plans ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.er_case_events ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS tenant_isolation_er_case_categories ON hrms.er_case_categories;
CREATE POLICY tenant_isolation_er_case_categories ON hrms.er_case_categories
    USING (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true))
    WITH CHECK (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true));

DROP POLICY IF EXISTS tenant_isolation_er_cases ON hrms.er_cases;
CREATE POLICY tenant_isolation_er_cases ON hrms.er_cases
    USING (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true))
    WITH CHECK (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true));

DROP POLICY IF EXISTS tenant_isolation_er_case_parties ON hrms.er_case_parties;
CREATE POLICY tenant_isolation_er_case_parties ON hrms.er_case_parties
    USING (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true))
    WITH CHECK (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true));

DROP POLICY IF EXISTS tenant_isolation_er_allegations ON hrms.er_allegations;
CREATE POLICY tenant_isolation_er_allegations ON hrms.er_allegations
    USING (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true))
    WITH CHECK (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true));

DROP POLICY IF EXISTS tenant_isolation_er_investigation_steps ON hrms.er_investigation_steps;
CREATE POLICY tenant_isolation_er_investigation_steps ON hrms.er_investigation_steps
    USING (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true))
    WITH CHECK (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true));

DROP POLICY IF EXISTS tenant_isolation_er_witness_notes ON hrms.er_witness_notes;
CREATE POLICY tenant_isolation_er_witness_notes ON hrms.er_witness_notes
    USING (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true))
    WITH CHECK (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true));

DROP POLICY IF EXISTS tenant_isolation_er_evidence_attachments ON hrms.er_evidence_attachments;
CREATE POLICY tenant_isolation_er_evidence_attachments ON hrms.er_evidence_attachments
    USING (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true))
    WITH CHECK (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true));

DROP POLICY IF EXISTS tenant_isolation_er_findings ON hrms.er_findings;
CREATE POLICY tenant_isolation_er_findings ON hrms.er_findings
    USING (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true))
    WITH CHECK (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true));

DROP POLICY IF EXISTS tenant_isolation_er_action_plans ON hrms.er_action_plans;
CREATE POLICY tenant_isolation_er_action_plans ON hrms.er_action_plans
    USING (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true))
    WITH CHECK (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true));

DROP POLICY IF EXISTS tenant_isolation_er_case_events ON hrms.er_case_events;
CREATE POLICY tenant_isolation_er_case_events ON hrms.er_case_events
    USING (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true))
    WITH CHECK (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true));
