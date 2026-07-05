CREATE TABLE IF NOT EXISTS hrms.hr_case_categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    code VARCHAR(80) NOT NULL,
    name VARCHAR(160) NOT NULL,
    description TEXT,
    confidentiality_default VARCHAR(30) NOT NULL DEFAULT 'normal',
    default_owner_role VARCHAR(80),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT hr_case_categories_confidentiality_chk CHECK (confidentiality_default IN ('normal','restricted','sensitive','grievance')),
    CONSTRAINT hr_case_categories_code_chk CHECK (length(trim(code)) > 0),
    CONSTRAINT hr_case_categories_name_chk CHECK (length(trim(name)) > 0)
);

CREATE UNIQUE INDEX IF NOT EXISTS hr_case_categories_tenant_code_idx
    ON hrms.hr_case_categories(tenant_id, lower(code)) WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.hr_case_sla_policies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    category_id UUID REFERENCES hrms.hr_case_categories(id) ON DELETE SET NULL,
    priority VARCHAR(20) NOT NULL DEFAULT 'normal',
    response_hours INT NOT NULL DEFAULT 8,
    resolution_hours INT NOT NULL DEFAULT 48,
    escalation_hours INT NOT NULL DEFAULT 24,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT hr_case_sla_priority_chk CHECK (priority IN ('low','normal','high','urgent')),
    CONSTRAINT hr_case_sla_hours_chk CHECK (response_hours >= 0 AND resolution_hours > 0 AND escalation_hours > 0)
);

CREATE UNIQUE INDEX IF NOT EXISTS hr_case_sla_policy_unique_idx
    ON hrms.hr_case_sla_policies(tenant_id, COALESCE(category_id, '00000000-0000-0000-0000-000000000000'::uuid), priority)
    WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.hr_cases (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    case_number VARCHAR(40) NOT NULL,
    category_id UUID REFERENCES hrms.hr_case_categories(id) ON DELETE SET NULL,
    case_type VARCHAR(80) NOT NULL DEFAULT 'general',
    title VARCHAR(220) NOT NULL,
    description TEXT NOT NULL,
    confidentiality_level VARCHAR(30) NOT NULL DEFAULT 'normal',
    requester_user_id UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    subject_employee_user_id UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    owner_user_id UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    owner_role VARCHAR(80),
    status VARCHAR(30) NOT NULL DEFAULT 'new',
    priority VARCHAR(20) NOT NULL DEFAULT 'normal',
    source_channel VARCHAR(40) NOT NULL DEFAULT 'web',
    first_response_due_at TIMESTAMPTZ,
    first_responded_at TIMESTAMPTZ,
    due_at TIMESTAMPTZ,
    resolved_at TIMESTAMPTZ,
    closed_at TIMESTAMPTZ,
    escalated_at TIMESTAMPTZ,
    escalation_level INT NOT NULL DEFAULT 0,
    last_activity_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    resolution_summary TEXT,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT hr_cases_case_number_chk CHECK (length(trim(case_number)) > 0),
    CONSTRAINT hr_cases_title_chk CHECK (length(trim(title)) > 0),
    CONSTRAINT hr_cases_description_chk CHECK (length(trim(description)) > 0),
    CONSTRAINT hr_cases_confidentiality_chk CHECK (confidentiality_level IN ('normal','restricted','sensitive','grievance')),
    CONSTRAINT hr_cases_status_chk CHECK (status IN ('new','open','in_progress','waiting_on_employee','waiting_on_hr','escalated','resolved','closed','cancelled')),
    CONSTRAINT hr_cases_priority_chk CHECK (priority IN ('low','normal','high','urgent'))
);

CREATE UNIQUE INDEX IF NOT EXISTS hr_cases_tenant_number_idx
    ON hrms.hr_cases(tenant_id, case_number) WHERE NOT inactive;
CREATE INDEX IF NOT EXISTS hr_cases_tenant_status_idx
    ON hrms.hr_cases(tenant_id, status, priority, due_at) WHERE NOT inactive;
CREATE INDEX IF NOT EXISTS hr_cases_requester_idx
    ON hrms.hr_cases(tenant_id, requester_user_id, status) WHERE requester_user_id IS NOT NULL AND NOT inactive;
CREATE INDEX IF NOT EXISTS hr_cases_subject_idx
    ON hrms.hr_cases(tenant_id, subject_employee_user_id, status) WHERE subject_employee_user_id IS NOT NULL AND NOT inactive;
CREATE INDEX IF NOT EXISTS hr_cases_owner_idx
    ON hrms.hr_cases(tenant_id, owner_user_id, owner_role, status) WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.hr_case_comments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    case_id UUID NOT NULL REFERENCES hrms.hr_cases(id) ON DELETE CASCADE,
    author_user_id UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    visibility VARCHAR(20) NOT NULL DEFAULT 'public',
    body TEXT NOT NULL,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT hr_case_comments_visibility_chk CHECK (visibility IN ('public','internal')),
    CONSTRAINT hr_case_comments_body_chk CHECK (length(trim(body)) > 0)
);

CREATE INDEX IF NOT EXISTS hr_case_comments_case_idx
    ON hrms.hr_case_comments(tenant_id, case_id, created_at DESC) WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.hr_case_attachments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    case_id UUID NOT NULL REFERENCES hrms.hr_cases(id) ON DELETE CASCADE,
    comment_id UUID REFERENCES hrms.hr_case_comments(id) ON DELETE SET NULL,
    file_name VARCHAR(255) NOT NULL,
    content_type VARCHAR(120) NOT NULL,
    object_key TEXT NOT NULL,
    visibility VARCHAR(20) NOT NULL DEFAULT 'public',
    uploaded_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT hr_case_attachments_visibility_chk CHECK (visibility IN ('public','internal')),
    CONSTRAINT hr_case_attachments_file_chk CHECK (length(trim(file_name)) > 0 AND length(trim(object_key)) > 0)
);

CREATE INDEX IF NOT EXISTS hr_case_attachments_case_idx
    ON hrms.hr_case_attachments(tenant_id, case_id, created_at DESC) WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.hr_case_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    case_id UUID NOT NULL REFERENCES hrms.hr_cases(id) ON DELETE CASCADE,
    event_type VARCHAR(80) NOT NULL,
    from_status VARCHAR(30),
    to_status VARCHAR(30),
    actor_user_id UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    comment TEXT,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT hr_case_events_type_chk CHECK (length(trim(event_type)) > 0)
);

CREATE INDEX IF NOT EXISTS hr_case_events_case_idx
    ON hrms.hr_case_events(tenant_id, case_id, created_at DESC) WHERE NOT inactive;

ALTER TABLE hrms.hr_case_categories ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.hr_case_sla_policies ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.hr_cases ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.hr_case_comments ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.hr_case_attachments ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.hr_case_events ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS tenant_isolation_hr_case_categories ON hrms.hr_case_categories;
CREATE POLICY tenant_isolation_hr_case_categories ON hrms.hr_case_categories
    USING (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true))
    WITH CHECK (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true));

DROP POLICY IF EXISTS tenant_isolation_hr_case_sla_policies ON hrms.hr_case_sla_policies;
CREATE POLICY tenant_isolation_hr_case_sla_policies ON hrms.hr_case_sla_policies
    USING (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true))
    WITH CHECK (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true));

DROP POLICY IF EXISTS tenant_isolation_hr_cases ON hrms.hr_cases;
CREATE POLICY tenant_isolation_hr_cases ON hrms.hr_cases
    USING (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true))
    WITH CHECK (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true));

DROP POLICY IF EXISTS tenant_isolation_hr_case_comments ON hrms.hr_case_comments;
CREATE POLICY tenant_isolation_hr_case_comments ON hrms.hr_case_comments
    USING (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true))
    WITH CHECK (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true));

DROP POLICY IF EXISTS tenant_isolation_hr_case_attachments ON hrms.hr_case_attachments;
CREATE POLICY tenant_isolation_hr_case_attachments ON hrms.hr_case_attachments
    USING (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true))
    WITH CHECK (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true));

DROP POLICY IF EXISTS tenant_isolation_hr_case_events ON hrms.hr_case_events;
CREATE POLICY tenant_isolation_hr_case_events ON hrms.hr_case_events
    USING (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true))
    WITH CHECK (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true));
