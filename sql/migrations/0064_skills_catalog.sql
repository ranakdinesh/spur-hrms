CREATE TABLE IF NOT EXISTS hrms.skill_categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID,
    parent_id UUID REFERENCES hrms.skill_categories(id),
    code TEXT NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    source_scope TEXT NOT NULL DEFAULT 'tenant',
    sort_order INTEGER NOT NULL DEFAULT 100,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    CONSTRAINT skill_categories_scope_chk CHECK (source_scope IN ('global', 'tenant')),
    CONSTRAINT skill_categories_scope_tenant_chk CHECK (
        (source_scope = 'global' AND tenant_id IS NULL)
        OR (source_scope = 'tenant' AND tenant_id IS NOT NULL)
    )
);

CREATE UNIQUE INDEX IF NOT EXISTS skill_categories_global_code_uq
    ON hrms.skill_categories (lower(code))
    WHERE tenant_id IS NULL AND inactive = FALSE;

CREATE UNIQUE INDEX IF NOT EXISTS skill_categories_tenant_code_uq
    ON hrms.skill_categories (tenant_id, lower(code))
    WHERE tenant_id IS NOT NULL AND inactive = FALSE;

CREATE INDEX IF NOT EXISTS skill_categories_tenant_parent_idx
    ON hrms.skill_categories (tenant_id, parent_id, sort_order)
    WHERE inactive = FALSE;

CREATE TABLE IF NOT EXISTS hrms.skills (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID,
    category_id UUID REFERENCES hrms.skill_categories(id),
    code TEXT NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    skill_type TEXT NOT NULL DEFAULT 'technical',
    source_scope TEXT NOT NULL DEFAULT 'tenant',
    certificate_required BOOLEAN NOT NULL DEFAULT FALSE,
    assessment_required BOOLEAN NOT NULL DEFAULT FALSE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    CONSTRAINT skills_type_chk CHECK (skill_type IN ('technical', 'functional', 'behavioral', 'compliance', 'tool', 'language', 'domain', 'custom')),
    CONSTRAINT skills_scope_chk CHECK (source_scope IN ('global', 'tenant')),
    CONSTRAINT skills_scope_tenant_chk CHECK (
        (source_scope = 'global' AND tenant_id IS NULL)
        OR (source_scope = 'tenant' AND tenant_id IS NOT NULL)
    )
);

CREATE UNIQUE INDEX IF NOT EXISTS skills_global_code_uq
    ON hrms.skills (lower(code))
    WHERE tenant_id IS NULL AND inactive = FALSE;

CREATE UNIQUE INDEX IF NOT EXISTS skills_tenant_code_uq
    ON hrms.skills (tenant_id, lower(code))
    WHERE tenant_id IS NOT NULL AND inactive = FALSE;

CREATE INDEX IF NOT EXISTS skills_tenant_category_idx
    ON hrms.skills (tenant_id, category_id, skill_type, is_active)
    WHERE inactive = FALSE;

CREATE INDEX IF NOT EXISTS skills_global_category_idx
    ON hrms.skills (category_id, skill_type, is_active)
    WHERE tenant_id IS NULL AND inactive = FALSE;

CREATE TABLE IF NOT EXISTS hrms.worker_skills (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    worker_profile_id UUID NOT NULL REFERENCES hrms.worker_profiles(id),
    skill_id UUID NOT NULL REFERENCES hrms.skills(id),
    skill_name_snapshot TEXT NOT NULL,
    proficiency TEXT NOT NULL DEFAULT 'beginner',
    years_experience NUMERIC(4,1),
    last_used_on DATE,
    verification_status TEXT NOT NULL DEFAULT 'self_declared',
    certificate_url TEXT,
    certificate_expires_on DATE,
    assessment_score NUMERIC(5,2),
    assessed_on DATE,
    endorsed_by UUID,
    endorsed_at TIMESTAMPTZ,
    verified_by UUID,
    verified_at TIMESTAMPTZ,
    notes TEXT,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    CONSTRAINT worker_skills_proficiency_chk CHECK (proficiency IN ('beginner', 'intermediate', 'advanced', 'expert')),
    CONSTRAINT worker_skills_verification_chk CHECK (verification_status IN ('self_declared', 'manager_endorsed', 'hr_verified', 'expired', 'rejected')),
    CONSTRAINT worker_skills_years_chk CHECK (years_experience IS NULL OR years_experience >= 0),
    CONSTRAINT worker_skills_score_chk CHECK (assessment_score IS NULL OR (assessment_score >= 0 AND assessment_score <= 100))
);

CREATE UNIQUE INDEX IF NOT EXISTS worker_skills_worker_skill_uq
    ON hrms.worker_skills (tenant_id, worker_profile_id, skill_id)
    WHERE inactive = FALSE;

CREATE INDEX IF NOT EXISTS worker_skills_tenant_status_idx
    ON hrms.worker_skills (tenant_id, verification_status, proficiency)
    WHERE inactive = FALSE;

CREATE INDEX IF NOT EXISTS worker_skills_skill_idx
    ON hrms.worker_skills (tenant_id, skill_id)
    WHERE inactive = FALSE;

CREATE INDEX IF NOT EXISTS worker_skills_certificate_expiry_idx
    ON hrms.worker_skills (tenant_id, certificate_expires_on)
    WHERE inactive = FALSE AND certificate_expires_on IS NOT NULL;

CREATE TABLE IF NOT EXISTS hrms.worker_skill_assessments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    worker_skill_id UUID NOT NULL REFERENCES hrms.worker_skills(id),
    assessment_type TEXT NOT NULL DEFAULT 'manager',
    result_status TEXT NOT NULL DEFAULT 'submitted',
    score NUMERIC(5,2),
    max_score NUMERIC(5,2),
    assessed_by UUID,
    assessed_on DATE NOT NULL DEFAULT CURRENT_DATE,
    evidence_url TEXT,
    notes TEXT,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID,
    CONSTRAINT worker_skill_assessments_type_chk CHECK (assessment_type IN ('self', 'manager', 'hr', 'external')),
    CONSTRAINT worker_skill_assessments_status_chk CHECK (result_status IN ('submitted', 'observed', 'passed', 'failed')),
    CONSTRAINT worker_skill_assessments_score_chk CHECK (
        score IS NULL OR score >= 0
    ),
    CONSTRAINT worker_skill_assessments_max_score_chk CHECK (
        max_score IS NULL OR max_score > 0
    )
);

CREATE INDEX IF NOT EXISTS worker_skill_assessments_skill_idx
    ON hrms.worker_skill_assessments (tenant_id, worker_skill_id, assessed_on DESC);

DROP TRIGGER IF EXISTS trg_skill_categories_updated_at ON hrms.skill_categories;
CREATE TRIGGER trg_skill_categories_updated_at
BEFORE UPDATE ON hrms.skill_categories
FOR EACH ROW EXECUTE FUNCTION hrms.set_updated_at();

DROP TRIGGER IF EXISTS trg_skills_updated_at ON hrms.skills;
CREATE TRIGGER trg_skills_updated_at
BEFORE UPDATE ON hrms.skills
FOR EACH ROW EXECUTE FUNCTION hrms.set_updated_at();

DROP TRIGGER IF EXISTS trg_worker_skills_updated_at ON hrms.worker_skills;
CREATE TRIGGER trg_worker_skills_updated_at
BEFORE UPDATE ON hrms.worker_skills
FOR EACH ROW EXECUTE FUNCTION hrms.set_updated_at();

ALTER TABLE hrms.skill_categories ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.skills ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.worker_skills ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.worker_skill_assessments ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS skill_categories_tenant_isolation ON hrms.skill_categories;
CREATE POLICY skill_categories_tenant_isolation ON hrms.skill_categories
USING (tenant_id IS NULL OR tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true')
WITH CHECK (
    (tenant_id IS NULL AND current_setting('app.is_super_admin', true) = 'true')
    OR tenant_id::text = current_setting('app.tenant_id', true)
    OR current_setting('app.is_super_admin', true) = 'true'
);

DROP POLICY IF EXISTS skills_tenant_isolation ON hrms.skills;
CREATE POLICY skills_tenant_isolation ON hrms.skills
USING (tenant_id IS NULL OR tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true')
WITH CHECK (
    (tenant_id IS NULL AND current_setting('app.is_super_admin', true) = 'true')
    OR tenant_id::text = current_setting('app.tenant_id', true)
    OR current_setting('app.is_super_admin', true) = 'true'
);

DROP POLICY IF EXISTS worker_skills_tenant_isolation ON hrms.worker_skills;
CREATE POLICY worker_skills_tenant_isolation ON hrms.worker_skills
USING (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true')
WITH CHECK (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true');

DROP POLICY IF EXISTS worker_skill_assessments_tenant_isolation ON hrms.worker_skill_assessments;
CREATE POLICY worker_skill_assessments_tenant_isolation ON hrms.worker_skill_assessments
USING (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true')
WITH CHECK (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true');

WITH seeded_categories AS (
    INSERT INTO hrms.skill_categories (code, name, description, source_scope, sort_order)
    VALUES
        ('technical', 'Technical Skills', 'Engineering, IT, analytics, and product delivery skills.', 'global', 10),
        ('functional', 'Functional Skills', 'Business-function skills used by HR, finance, sales, operations, and support teams.', 'global', 20),
        ('behavioral', 'Behavioral Skills', 'Collaboration, leadership, communication, and workplace behavior skills.', 'global', 30),
        ('compliance', 'Compliance Skills', 'Mandatory statutory, safety, and regulated-work capabilities.', 'global', 40),
        ('tools', 'Tools and Platforms', 'Software, machinery, and platform-specific capabilities.', 'global', 50),
        ('languages', 'Languages', 'Spoken and written language capabilities.', 'global', 60)
    ON CONFLICT DO NOTHING
    RETURNING id, code
)
INSERT INTO hrms.skills (category_id, code, name, description, skill_type, source_scope, certificate_required, assessment_required)
VALUES
    ((SELECT id FROM hrms.skill_categories WHERE tenant_id IS NULL AND code = 'technical' AND NOT inactive), 'excel-advanced', 'Advanced Excel', 'Advanced spreadsheet analysis, formulas, pivots, and operational reporting.', 'tool', 'global', FALSE, FALSE),
    ((SELECT id FROM hrms.skill_categories WHERE tenant_id IS NULL AND code = 'technical' AND NOT inactive), 'sql', 'SQL', 'Structured query writing, data extraction, joins, aggregations, and reporting.', 'technical', 'global', FALSE, FALSE),
    ((SELECT id FROM hrms.skill_categories WHERE tenant_id IS NULL AND code = 'technical' AND NOT inactive), 'python', 'Python', 'Python programming for automation, data, integrations, or application development.', 'technical', 'global', FALSE, FALSE),
    ((SELECT id FROM hrms.skill_categories WHERE tenant_id IS NULL AND code = 'functional' AND NOT inactive), 'payroll-operations', 'Payroll Operations', 'Payroll processing, controls, statutory inputs, and employee pay queries.', 'functional', 'global', FALSE, TRUE),
    ((SELECT id FROM hrms.skill_categories WHERE tenant_id IS NULL AND code = 'functional' AND NOT inactive), 'recruitment', 'Recruitment', 'Hiring workflow execution, screening, interviews, offers, and onboarding handoff.', 'functional', 'global', FALSE, FALSE),
    ((SELECT id FROM hrms.skill_categories WHERE tenant_id IS NULL AND code = 'behavioral' AND NOT inactive), 'people-management', 'People Management', 'Coaching, feedback, performance conversations, and team coordination.', 'behavioral', 'global', FALSE, TRUE),
    ((SELECT id FROM hrms.skill_categories WHERE tenant_id IS NULL AND code = 'behavioral' AND NOT inactive), 'communication', 'Communication', 'Clear written and verbal communication across teams and stakeholders.', 'behavioral', 'global', FALSE, FALSE),
    ((SELECT id FROM hrms.skill_categories WHERE tenant_id IS NULL AND code = 'compliance' AND NOT inactive), 'workplace-safety', 'Workplace Safety', 'Safety awareness, incident prevention, and site-specific safety practices.', 'compliance', 'global', TRUE, TRUE),
    ((SELECT id FROM hrms.skill_categories WHERE tenant_id IS NULL AND code = 'languages' AND NOT inactive), 'english', 'English', 'Professional English communication.', 'language', 'global', FALSE, FALSE),
    ((SELECT id FROM hrms.skill_categories WHERE tenant_id IS NULL AND code = 'languages' AND NOT inactive), 'hindi', 'Hindi', 'Professional Hindi communication.', 'language', 'global', FALSE, FALSE)
ON CONFLICT DO NOTHING;
