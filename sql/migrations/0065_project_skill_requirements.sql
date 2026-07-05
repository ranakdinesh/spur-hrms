CREATE TABLE IF NOT EXISTS hrms.project_skill_requirements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    project_id UUID REFERENCES hrms.projects(id),
    engagement_id UUID REFERENCES hrms.engagements(id),
    skill_id UUID NOT NULL REFERENCES hrms.skills(id),
    required_proficiency TEXT NOT NULL DEFAULT 'intermediate',
    min_years_experience NUMERIC(4,1),
    required_count INTEGER NOT NULL DEFAULT 1,
    importance TEXT NOT NULL DEFAULT 'required',
    requirement_source TEXT NOT NULL DEFAULT 'project',
    notes TEXT,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    CONSTRAINT project_skill_requirements_target_chk CHECK (project_id IS NOT NULL OR engagement_id IS NOT NULL),
    CONSTRAINT project_skill_requirements_proficiency_chk CHECK (required_proficiency IN ('beginner', 'intermediate', 'advanced', 'expert')),
    CONSTRAINT project_skill_requirements_importance_chk CHECK (importance IN ('nice_to_have', 'required', 'critical')),
    CONSTRAINT project_skill_requirements_source_chk CHECK (requirement_source IN ('project', 'engagement', 'role', 'client', 'compliance')),
    CONSTRAINT project_skill_requirements_count_chk CHECK (required_count > 0),
    CONSTRAINT project_skill_requirements_years_chk CHECK (min_years_experience IS NULL OR min_years_experience >= 0)
);

CREATE UNIQUE INDEX IF NOT EXISTS project_skill_requirements_project_skill_uq
    ON hrms.project_skill_requirements (tenant_id, project_id, skill_id)
    WHERE inactive = FALSE AND project_id IS NOT NULL AND engagement_id IS NULL;

CREATE UNIQUE INDEX IF NOT EXISTS project_skill_requirements_engagement_skill_uq
    ON hrms.project_skill_requirements (tenant_id, engagement_id, skill_id)
    WHERE inactive = FALSE AND engagement_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS project_skill_requirements_tenant_project_idx
    ON hrms.project_skill_requirements (tenant_id, project_id, importance)
    WHERE inactive = FALSE;

CREATE INDEX IF NOT EXISTS project_skill_requirements_tenant_engagement_idx
    ON hrms.project_skill_requirements (tenant_id, engagement_id, importance)
    WHERE inactive = FALSE;

CREATE INDEX IF NOT EXISTS project_skill_requirements_skill_idx
    ON hrms.project_skill_requirements (tenant_id, skill_id)
    WHERE inactive = FALSE;

DROP TRIGGER IF EXISTS trg_project_skill_requirements_updated_at ON hrms.project_skill_requirements;
CREATE TRIGGER trg_project_skill_requirements_updated_at
BEFORE UPDATE ON hrms.project_skill_requirements
FOR EACH ROW EXECUTE FUNCTION hrms.set_updated_at();

ALTER TABLE hrms.project_skill_requirements ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS project_skill_requirements_tenant_isolation ON hrms.project_skill_requirements;
CREATE POLICY project_skill_requirements_tenant_isolation ON hrms.project_skill_requirements
USING (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true')
WITH CHECK (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true');
