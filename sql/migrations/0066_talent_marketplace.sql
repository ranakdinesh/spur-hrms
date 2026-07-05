CREATE TABLE IF NOT EXISTS hrms.talent_marketplace_opportunities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    project_id UUID REFERENCES hrms.projects(id) ON DELETE SET NULL,
    engagement_id UUID REFERENCES hrms.engagements(id) ON DELETE SET NULL,
    source_requirement_id UUID REFERENCES hrms.project_skill_requirements(id) ON DELETE SET NULL,
    job_posting_id UUID REFERENCES hrms.job_postings(id) ON DELETE SET NULL,
    title TEXT NOT NULL,
    description TEXT,
    opportunity_type TEXT NOT NULL DEFAULT 'project_assignment',
    status TEXT NOT NULL DEFAULT 'draft',
    visibility TEXT NOT NULL DEFAULT 'all_workers',
    priority TEXT NOT NULL DEFAULT 'normal',
    seats INTEGER NOT NULL DEFAULT 1,
    location_mode TEXT NOT NULL DEFAULT 'flexible',
    min_allocation_percent INTEGER,
    duration_label TEXT,
    start_date DATE,
    due_date DATE,
    candidate_fallback_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    candidate_fallback_status TEXT NOT NULL DEFAULT 'not_needed',
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    CONSTRAINT talent_marketplace_opportunities_type_chk CHECK (opportunity_type IN ('project_assignment', 'gig', 'role', 'mentorship', 'stretch', 'backfill')),
    CONSTRAINT talent_marketplace_opportunities_status_chk CHECK (status IN ('draft', 'open', 'paused', 'filled', 'closed', 'cancelled')),
    CONSTRAINT talent_marketplace_opportunities_visibility_chk CHECK (visibility IN ('all_workers', 'invited', 'manager_nomination')),
    CONSTRAINT talent_marketplace_opportunities_priority_chk CHECK (priority IN ('low', 'normal', 'high', 'critical')),
    CONSTRAINT talent_marketplace_opportunities_location_mode_chk CHECK (location_mode IN ('onsite', 'remote', 'hybrid', 'flexible')),
    CONSTRAINT talent_marketplace_opportunities_seats_chk CHECK (seats > 0),
    CONSTRAINT talent_marketplace_opportunities_allocation_chk CHECK (min_allocation_percent IS NULL OR (min_allocation_percent >= 1 AND min_allocation_percent <= 100)),
    CONSTRAINT talent_marketplace_opportunities_dates_chk CHECK (start_date IS NULL OR due_date IS NULL OR due_date >= start_date),
    CONSTRAINT talent_marketplace_opportunities_fallback_chk CHECK (candidate_fallback_status IN ('not_needed', 'monitoring', 'recommended', 'opened')),
    CONSTRAINT talent_marketplace_opportunities_metadata_object_chk CHECK (jsonb_typeof(metadata) = 'object')
);

CREATE INDEX IF NOT EXISTS talent_marketplace_opportunities_tenant_status_idx
    ON hrms.talent_marketplace_opportunities (tenant_id, status, priority)
    WHERE inactive = FALSE;

CREATE INDEX IF NOT EXISTS talent_marketplace_opportunities_project_idx
    ON hrms.talent_marketplace_opportunities (tenant_id, project_id)
    WHERE project_id IS NOT NULL AND inactive = FALSE;

CREATE INDEX IF NOT EXISTS talent_marketplace_opportunities_engagement_idx
    ON hrms.talent_marketplace_opportunities (tenant_id, engagement_id)
    WHERE engagement_id IS NOT NULL AND inactive = FALSE;

CREATE TABLE IF NOT EXISTS hrms.talent_marketplace_applications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    opportunity_id UUID NOT NULL REFERENCES hrms.talent_marketplace_opportunities(id) ON DELETE CASCADE,
    worker_profile_id UUID NOT NULL REFERENCES hrms.worker_profiles(id) ON DELETE RESTRICT,
    status TEXT NOT NULL DEFAULT 'applied',
    match_score NUMERIC(5,2),
    match_reasons JSONB NOT NULL DEFAULT '{}'::jsonb,
    worker_note TEXT,
    manager_note TEXT,
    decided_at TIMESTAMPTZ,
    decided_by UUID,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    CONSTRAINT talent_marketplace_applications_status_chk CHECK (status IN ('recommended', 'invited', 'interested', 'applied', 'accepted', 'declined', 'withdrawn', 'rejected', 'assigned')),
    CONSTRAINT talent_marketplace_applications_score_chk CHECK (match_score IS NULL OR (match_score >= 0 AND match_score <= 100)),
    CONSTRAINT talent_marketplace_applications_reasons_object_chk CHECK (jsonb_typeof(match_reasons) = 'object')
);

CREATE UNIQUE INDEX IF NOT EXISTS talent_marketplace_applications_worker_uq
    ON hrms.talent_marketplace_applications (tenant_id, opportunity_id, worker_profile_id)
    WHERE inactive = FALSE;

CREATE INDEX IF NOT EXISTS talent_marketplace_applications_tenant_status_idx
    ON hrms.talent_marketplace_applications (tenant_id, status, updated_at DESC)
    WHERE inactive = FALSE;

CREATE INDEX IF NOT EXISTS talent_marketplace_applications_worker_idx
    ON hrms.talent_marketplace_applications (tenant_id, worker_profile_id, status)
    WHERE inactive = FALSE;

CREATE TABLE IF NOT EXISTS hrms.talent_marketplace_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    opportunity_id UUID REFERENCES hrms.talent_marketplace_opportunities(id) ON DELETE CASCADE,
    application_id UUID REFERENCES hrms.talent_marketplace_applications(id) ON DELETE CASCADE,
    actor_user_id UUID,
    action TEXT NOT NULL,
    from_status TEXT,
    to_status TEXT,
    notes TEXT,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT talent_marketplace_events_metadata_object_chk CHECK (jsonb_typeof(metadata) = 'object')
);

CREATE INDEX IF NOT EXISTS talent_marketplace_events_opportunity_idx
    ON hrms.talent_marketplace_events (tenant_id, opportunity_id, created_at DESC)
    WHERE inactive = FALSE;

CREATE INDEX IF NOT EXISTS talent_marketplace_events_application_idx
    ON hrms.talent_marketplace_events (tenant_id, application_id, created_at DESC)
    WHERE inactive = FALSE;

DROP TRIGGER IF EXISTS trg_talent_marketplace_opportunities_updated_at ON hrms.talent_marketplace_opportunities;
CREATE TRIGGER trg_talent_marketplace_opportunities_updated_at
BEFORE UPDATE ON hrms.talent_marketplace_opportunities
FOR EACH ROW EXECUTE FUNCTION hrms.set_updated_at();

DROP TRIGGER IF EXISTS trg_talent_marketplace_applications_updated_at ON hrms.talent_marketplace_applications;
CREATE TRIGGER trg_talent_marketplace_applications_updated_at
BEFORE UPDATE ON hrms.talent_marketplace_applications
FOR EACH ROW EXECUTE FUNCTION hrms.set_updated_at();

ALTER TABLE hrms.talent_marketplace_opportunities ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.talent_marketplace_applications ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.talent_marketplace_events ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS talent_marketplace_opportunities_tenant_isolation ON hrms.talent_marketplace_opportunities;
CREATE POLICY talent_marketplace_opportunities_tenant_isolation ON hrms.talent_marketplace_opportunities
USING (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true')
WITH CHECK (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true');

DROP POLICY IF EXISTS talent_marketplace_applications_tenant_isolation ON hrms.talent_marketplace_applications;
CREATE POLICY talent_marketplace_applications_tenant_isolation ON hrms.talent_marketplace_applications
USING (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true')
WITH CHECK (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true');

DROP POLICY IF EXISTS talent_marketplace_events_tenant_isolation ON hrms.talent_marketplace_events;
CREATE POLICY talent_marketplace_events_tenant_isolation ON hrms.talent_marketplace_events
USING (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true')
WITH CHECK (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true');
