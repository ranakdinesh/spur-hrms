CREATE TABLE IF NOT EXISTS hrms.learning_courses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    code TEXT NOT NULL,
    title TEXT NOT NULL,
    description TEXT,
    course_type TEXT NOT NULL DEFAULT 'functional',
    delivery_mode TEXT NOT NULL DEFAULT 'self_paced',
    provider TEXT,
    duration_minutes INTEGER NOT NULL DEFAULT 0,
    skill_id UUID REFERENCES hrms.skills(id),
    compliance_rule_id UUID REFERENCES hrms.compliance_rules(id),
    mandatory BOOLEAN NOT NULL DEFAULT FALSE,
    ai_readiness BOOLEAN NOT NULL DEFAULT FALSE,
    certificate_required BOOLEAN NOT NULL DEFAULT FALSE,
    budget_amount NUMERIC(14,2) NOT NULL DEFAULT 0,
    currency_code TEXT NOT NULL DEFAULT 'INR',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    CONSTRAINT learning_courses_type_chk CHECK (course_type IN ('technical','functional','compliance','leadership','behavioral','ai_readiness','custom')),
    CONSTRAINT learning_courses_delivery_chk CHECK (delivery_mode IN ('self_paced','classroom','virtual','blended','external')),
    CONSTRAINT learning_courses_duration_chk CHECK (duration_minutes >= 0),
    CONSTRAINT learning_courses_budget_chk CHECK (budget_amount >= 0)
);

CREATE UNIQUE INDEX IF NOT EXISTS learning_courses_tenant_code_uq
    ON hrms.learning_courses (tenant_id, lower(code))
    WHERE inactive = FALSE;

CREATE INDEX IF NOT EXISTS learning_courses_tenant_filters_idx
    ON hrms.learning_courses (tenant_id, course_type, mandatory, ai_readiness, is_active)
    WHERE inactive = FALSE;

CREATE TABLE IF NOT EXISTS hrms.learning_paths (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    code TEXT NOT NULL,
    title TEXT NOT NULL,
    description TEXT,
    path_type TEXT NOT NULL DEFAULT 'upskilling',
    target_role TEXT,
    skill_id UUID REFERENCES hrms.skills(id),
    ai_readiness BOOLEAN NOT NULL DEFAULT FALSE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    CONSTRAINT learning_paths_type_chk CHECK (path_type IN ('onboarding','compliance','upskilling','leadership','ai_readiness','custom'))
);

CREATE UNIQUE INDEX IF NOT EXISTS learning_paths_tenant_code_uq
    ON hrms.learning_paths (tenant_id, lower(code))
    WHERE inactive = FALSE;

CREATE TABLE IF NOT EXISTS hrms.learning_path_courses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    path_id UUID NOT NULL REFERENCES hrms.learning_paths(id),
    course_id UUID NOT NULL REFERENCES hrms.learning_courses(id),
    sort_order INTEGER NOT NULL DEFAULT 100,
    required BOOLEAN NOT NULL DEFAULT TRUE,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID
);

CREATE UNIQUE INDEX IF NOT EXISTS learning_path_courses_path_course_uq
    ON hrms.learning_path_courses (tenant_id, path_id, course_id)
    WHERE inactive = FALSE;

CREATE TABLE IF NOT EXISTS hrms.learning_enrollments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    course_id UUID NOT NULL REFERENCES hrms.learning_courses(id),
    path_id UUID REFERENCES hrms.learning_paths(id),
    worker_profile_id UUID NOT NULL REFERENCES hrms.worker_profiles(id),
    assignment_source TEXT NOT NULL DEFAULT 'hr',
    status TEXT NOT NULL DEFAULT 'assigned',
    nominated_by UUID,
    assigned_by UUID,
    due_date DATE,
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    score NUMERIC(6,2),
    certificate_url TEXT,
    certificate_file_name TEXT,
    certificate_content_type TEXT,
    notes TEXT,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    CONSTRAINT learning_enrollments_source_chk CHECK (assignment_source IN ('self','manager','hr','compliance','skill_gap','ai','manual')),
    CONSTRAINT learning_enrollments_status_chk CHECK (status IN ('nominated','assigned','approved','in_progress','completed','overdue','waived','cancelled')),
    CONSTRAINT learning_enrollments_score_chk CHECK (score IS NULL OR score >= 0)
);

CREATE UNIQUE INDEX IF NOT EXISTS learning_enrollments_worker_course_uq
    ON hrms.learning_enrollments (tenant_id, worker_profile_id, course_id)
    WHERE inactive = FALSE AND status NOT IN ('cancelled','waived');

CREATE INDEX IF NOT EXISTS learning_enrollments_tenant_status_idx
    ON hrms.learning_enrollments (tenant_id, status, due_date)
    WHERE inactive = FALSE;

CREATE TABLE IF NOT EXISTS hrms.learning_recommendations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    worker_profile_id UUID REFERENCES hrms.worker_profiles(id),
    skill_id UUID REFERENCES hrms.skills(id),
    course_id UUID REFERENCES hrms.learning_courses(id),
    path_id UUID REFERENCES hrms.learning_paths(id),
    source_type TEXT NOT NULL DEFAULT 'manual',
    reason TEXT NOT NULL,
    priority TEXT NOT NULL DEFAULT 'medium',
    confidence_score NUMERIC(5,2) NOT NULL DEFAULT 0,
    status TEXT NOT NULL DEFAULT 'open',
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    CONSTRAINT learning_recommendations_source_chk CHECK (source_type IN ('skill_gap','compliance','performance','ai','manager','manual')),
    CONSTRAINT learning_recommendations_priority_chk CHECK (priority IN ('low','medium','high','urgent')),
    CONSTRAINT learning_recommendations_status_chk CHECK (status IN ('open','accepted','assigned','dismissed','completed'))
);

CREATE INDEX IF NOT EXISTS learning_recommendations_tenant_status_idx
    ON hrms.learning_recommendations (tenant_id, status, priority)
    WHERE inactive = FALSE;

DROP TRIGGER IF EXISTS trg_learning_courses_updated_at ON hrms.learning_courses;
CREATE TRIGGER trg_learning_courses_updated_at
BEFORE UPDATE ON hrms.learning_courses
FOR EACH ROW EXECUTE FUNCTION hrms.set_updated_at();

DROP TRIGGER IF EXISTS trg_learning_paths_updated_at ON hrms.learning_paths;
CREATE TRIGGER trg_learning_paths_updated_at
BEFORE UPDATE ON hrms.learning_paths
FOR EACH ROW EXECUTE FUNCTION hrms.set_updated_at();

DROP TRIGGER IF EXISTS trg_learning_path_courses_updated_at ON hrms.learning_path_courses;
CREATE TRIGGER trg_learning_path_courses_updated_at
BEFORE UPDATE ON hrms.learning_path_courses
FOR EACH ROW EXECUTE FUNCTION hrms.set_updated_at();

DROP TRIGGER IF EXISTS trg_learning_enrollments_updated_at ON hrms.learning_enrollments;
CREATE TRIGGER trg_learning_enrollments_updated_at
BEFORE UPDATE ON hrms.learning_enrollments
FOR EACH ROW EXECUTE FUNCTION hrms.set_updated_at();

DROP TRIGGER IF EXISTS trg_learning_recommendations_updated_at ON hrms.learning_recommendations;
CREATE TRIGGER trg_learning_recommendations_updated_at
BEFORE UPDATE ON hrms.learning_recommendations
FOR EACH ROW EXECUTE FUNCTION hrms.set_updated_at();

ALTER TABLE hrms.learning_courses ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.learning_paths ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.learning_path_courses ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.learning_enrollments ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.learning_recommendations ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS learning_courses_tenant_isolation ON hrms.learning_courses;
CREATE POLICY learning_courses_tenant_isolation ON hrms.learning_courses
USING (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true')
WITH CHECK (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true');

DROP POLICY IF EXISTS learning_paths_tenant_isolation ON hrms.learning_paths;
CREATE POLICY learning_paths_tenant_isolation ON hrms.learning_paths
USING (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true')
WITH CHECK (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true');

DROP POLICY IF EXISTS learning_path_courses_tenant_isolation ON hrms.learning_path_courses;
CREATE POLICY learning_path_courses_tenant_isolation ON hrms.learning_path_courses
USING (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true')
WITH CHECK (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true');

DROP POLICY IF EXISTS learning_enrollments_tenant_isolation ON hrms.learning_enrollments;
CREATE POLICY learning_enrollments_tenant_isolation ON hrms.learning_enrollments
USING (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true')
WITH CHECK (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true');

DROP POLICY IF EXISTS learning_recommendations_tenant_isolation ON hrms.learning_recommendations;
CREATE POLICY learning_recommendations_tenant_isolation ON hrms.learning_recommendations
USING (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true')
WITH CHECK (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true');
