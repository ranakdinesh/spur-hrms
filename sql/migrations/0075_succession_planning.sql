CREATE TABLE IF NOT EXISTS hrms.succession_review_cycles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    code TEXT NOT NULL,
    name TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'draft',
    starts_on DATE,
    ends_on DATE,
    confidentiality_level TEXT NOT NULL DEFAULT 'hr_only',
    notes TEXT,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    CONSTRAINT succession_review_cycles_status_chk CHECK (status IN ('draft','active','review','closed','archived','cancelled')),
    CONSTRAINT succession_review_cycles_confidentiality_chk CHECK (confidentiality_level IN ('hr_only','leadership','restricted')),
    CONSTRAINT succession_review_cycles_dates_chk CHECK (ends_on IS NULL OR starts_on IS NULL OR ends_on >= starts_on)
);

CREATE UNIQUE INDEX IF NOT EXISTS succession_review_cycles_tenant_code_uq
    ON hrms.succession_review_cycles (tenant_id, lower(code))
    WHERE inactive = FALSE;

CREATE TABLE IF NOT EXISTS hrms.succession_critical_roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    cycle_id UUID REFERENCES hrms.succession_review_cycles(id),
    code TEXT NOT NULL,
    title TEXT NOT NULL,
    department_id UUID REFERENCES hrms.departments(id),
    designation_id UUID REFERENCES hrms.designations(id),
    incumbent_worker_profile_id UUID REFERENCES hrms.worker_profiles(id),
    emergency_cover_worker_profile_id UUID REFERENCES hrms.worker_profiles(id),
    criticality TEXT NOT NULL DEFAULT 'medium',
    impact_level TEXT NOT NULL DEFAULT 'medium',
    vacancy_risk TEXT NOT NULL DEFAULT 'medium',
    attrition_risk TEXT NOT NULL DEFAULT 'unknown',
    readiness_target TEXT NOT NULL DEFAULT 'ready_now',
    successor_required_count INTEGER NOT NULL DEFAULT 2,
    role_summary TEXT,
    status TEXT NOT NULL DEFAULT 'open',
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    CONSTRAINT succession_critical_roles_criticality_chk CHECK (criticality IN ('low','medium','high','critical')),
    CONSTRAINT succession_critical_roles_impact_chk CHECK (impact_level IN ('low','medium','high','critical')),
    CONSTRAINT succession_critical_roles_vacancy_risk_chk CHECK (vacancy_risk IN ('low','medium','high','critical')),
    CONSTRAINT succession_critical_roles_attrition_risk_chk CHECK (attrition_risk IN ('unknown','low','medium','high','critical')),
    CONSTRAINT succession_critical_roles_readiness_target_chk CHECK (readiness_target IN ('ready_now','ready_soon','ready_later','future_potential')),
    CONSTRAINT succession_critical_roles_status_chk CHECK (status IN ('open','covered','at_risk','closed','archived')),
    CONSTRAINT succession_critical_roles_required_count_chk CHECK (successor_required_count >= 0)
);

CREATE UNIQUE INDEX IF NOT EXISTS succession_critical_roles_tenant_code_uq
    ON hrms.succession_critical_roles (tenant_id, lower(code))
    WHERE inactive = FALSE;

CREATE INDEX IF NOT EXISTS succession_critical_roles_tenant_risk_idx
    ON hrms.succession_critical_roles (tenant_id, criticality, vacancy_risk, status)
    WHERE inactive = FALSE;

CREATE TABLE IF NOT EXISTS hrms.succession_successor_nominations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    critical_role_id UUID NOT NULL REFERENCES hrms.succession_critical_roles(id),
    worker_profile_id UUID NOT NULL REFERENCES hrms.worker_profiles(id),
    nominated_by UUID,
    readiness_level TEXT NOT NULL DEFAULT 'ready_later',
    readiness_months INTEGER NOT NULL DEFAULT 12,
    potential_rating TEXT,
    performance_rating TEXT,
    retention_risk TEXT NOT NULL DEFAULT 'medium',
    mobility_preference TEXT,
    nomination_status TEXT NOT NULL DEFAULT 'nominated',
    development_notes TEXT,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    CONSTRAINT succession_successor_readiness_chk CHECK (readiness_level IN ('ready_now','ready_soon','ready_later','future_potential')),
    CONSTRAINT succession_successor_months_chk CHECK (readiness_months >= 0),
    CONSTRAINT succession_successor_retention_chk CHECK (retention_risk IN ('low','medium','high','critical')),
    CONSTRAINT succession_successor_status_chk CHECK (nomination_status IN ('nominated','approved','development','ready','rejected','withdrawn'))
);

CREATE UNIQUE INDEX IF NOT EXISTS succession_successor_role_worker_uq
    ON hrms.succession_successor_nominations (tenant_id, critical_role_id, worker_profile_id)
    WHERE inactive = FALSE AND nomination_status NOT IN ('rejected','withdrawn');

CREATE INDEX IF NOT EXISTS succession_successor_readiness_idx
    ON hrms.succession_successor_nominations (tenant_id, readiness_level, nomination_status)
    WHERE inactive = FALSE;

CREATE TABLE IF NOT EXISTS hrms.succession_development_actions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    nomination_id UUID REFERENCES hrms.succession_successor_nominations(id),
    critical_role_id UUID REFERENCES hrms.succession_critical_roles(id),
    worker_profile_id UUID NOT NULL REFERENCES hrms.worker_profiles(id),
    action_type TEXT NOT NULL DEFAULT 'learning',
    title TEXT NOT NULL,
    learning_course_id UUID REFERENCES hrms.learning_courses(id),
    learning_path_id UUID REFERENCES hrms.learning_paths(id),
    owner_user_id UUID,
    due_date DATE,
    status TEXT NOT NULL DEFAULT 'open',
    notes TEXT,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    CONSTRAINT succession_development_action_type_chk CHECK (action_type IN ('learning','stretch_assignment','mentoring','coaching','certification','performance_goal','mobility','other')),
    CONSTRAINT succession_development_status_chk CHECK (status IN ('open','in_progress','completed','cancelled','overdue'))
);

CREATE INDEX IF NOT EXISTS succession_development_actions_worker_idx
    ON hrms.succession_development_actions (tenant_id, worker_profile_id, status, due_date)
    WHERE inactive = FALSE;

CREATE TABLE IF NOT EXISTS hrms.succession_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    source_type TEXT NOT NULL,
    source_id UUID,
    action TEXT NOT NULL,
    from_status TEXT,
    to_status TEXT,
    remarks TEXT,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID
);

CREATE INDEX IF NOT EXISTS succession_events_source_idx
    ON hrms.succession_events (tenant_id, source_type, source_id, created_at DESC);

DROP TRIGGER IF EXISTS trg_succession_review_cycles_updated_at ON hrms.succession_review_cycles;
CREATE TRIGGER trg_succession_review_cycles_updated_at BEFORE UPDATE ON hrms.succession_review_cycles FOR EACH ROW EXECUTE FUNCTION hrms.set_updated_at();

DROP TRIGGER IF EXISTS trg_succession_critical_roles_updated_at ON hrms.succession_critical_roles;
CREATE TRIGGER trg_succession_critical_roles_updated_at BEFORE UPDATE ON hrms.succession_critical_roles FOR EACH ROW EXECUTE FUNCTION hrms.set_updated_at();

DROP TRIGGER IF EXISTS trg_succession_successor_nominations_updated_at ON hrms.succession_successor_nominations;
CREATE TRIGGER trg_succession_successor_nominations_updated_at BEFORE UPDATE ON hrms.succession_successor_nominations FOR EACH ROW EXECUTE FUNCTION hrms.set_updated_at();

DROP TRIGGER IF EXISTS trg_succession_development_actions_updated_at ON hrms.succession_development_actions;
CREATE TRIGGER trg_succession_development_actions_updated_at BEFORE UPDATE ON hrms.succession_development_actions FOR EACH ROW EXECUTE FUNCTION hrms.set_updated_at();

ALTER TABLE hrms.succession_review_cycles ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.succession_critical_roles ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.succession_successor_nominations ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.succession_development_actions ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.succession_events ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS succession_review_cycles_tenant_isolation ON hrms.succession_review_cycles;
CREATE POLICY succession_review_cycles_tenant_isolation ON hrms.succession_review_cycles
USING (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true')
WITH CHECK (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true');

DROP POLICY IF EXISTS succession_critical_roles_tenant_isolation ON hrms.succession_critical_roles;
CREATE POLICY succession_critical_roles_tenant_isolation ON hrms.succession_critical_roles
USING (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true')
WITH CHECK (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true');

DROP POLICY IF EXISTS succession_successor_nominations_tenant_isolation ON hrms.succession_successor_nominations;
CREATE POLICY succession_successor_nominations_tenant_isolation ON hrms.succession_successor_nominations
USING (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true')
WITH CHECK (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true');

DROP POLICY IF EXISTS succession_development_actions_tenant_isolation ON hrms.succession_development_actions;
CREATE POLICY succession_development_actions_tenant_isolation ON hrms.succession_development_actions
USING (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true')
WITH CHECK (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true');

DROP POLICY IF EXISTS succession_events_tenant_isolation ON hrms.succession_events;
CREATE POLICY succession_events_tenant_isolation ON hrms.succession_events
USING (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true')
WITH CHECK (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true');
