CREATE TABLE IF NOT EXISTS hrms.performance_checkins (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    worker_profile_id UUID NOT NULL REFERENCES hrms.worker_profiles(id) ON DELETE CASCADE,
    reviewer_worker_profile_id UUID REFERENCES hrms.worker_profiles(id) ON DELETE SET NULL,
    cycle_id UUID REFERENCES hrms.okr_cycles(id) ON DELETE SET NULL,
    checkin_date DATE NOT NULL DEFAULT CURRENT_DATE,
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    mood TEXT NOT NULL DEFAULT 'neutral',
    status TEXT NOT NULL DEFAULT 'draft',
    visibility TEXT NOT NULL DEFAULT 'worker_manager_hr',
    highlights TEXT,
    blockers TEXT,
    next_plan TEXT,
    employee_comment TEXT,
    manager_comment TEXT,
    score NUMERIC(5,2),
    calibration_bucket TEXT,
    reviewed_at TIMESTAMPTZ,
    reviewed_by UUID,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    CONSTRAINT performance_checkins_period_chk CHECK (period_end >= period_start),
    CONSTRAINT performance_checkins_mood_chk CHECK (mood IN ('great', 'good', 'neutral', 'low', 'stressed')),
    CONSTRAINT performance_checkins_status_chk CHECK (status IN ('draft', 'submitted', 'reviewed', 'closed')),
    CONSTRAINT performance_checkins_visibility_chk CHECK (visibility IN ('worker_manager_hr', 'manager_hr', 'hr_only')),
    CONSTRAINT performance_checkins_score_chk CHECK (score IS NULL OR (score >= 0 AND score <= 5)),
    CONSTRAINT performance_checkins_calibration_chk CHECK (calibration_bucket IS NULL OR calibration_bucket IN ('high', 'solid', 'watch', 'improve')),
    CONSTRAINT performance_checkins_metadata_object_chk CHECK (jsonb_typeof(metadata) = 'object')
);

CREATE UNIQUE INDEX IF NOT EXISTS performance_checkins_worker_period_uq
    ON hrms.performance_checkins (tenant_id, worker_profile_id, period_start, period_end)
    WHERE inactive = FALSE;

CREATE INDEX IF NOT EXISTS performance_checkins_tenant_status_idx
    ON hrms.performance_checkins (tenant_id, status, checkin_date DESC)
    WHERE inactive = FALSE;

CREATE INDEX IF NOT EXISTS performance_checkins_worker_idx
    ON hrms.performance_checkins (tenant_id, worker_profile_id, checkin_date DESC)
    WHERE inactive = FALSE;

CREATE TABLE IF NOT EXISTS hrms.feedback_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    subject_worker_profile_id UUID NOT NULL REFERENCES hrms.worker_profiles(id) ON DELETE CASCADE,
    requester_worker_profile_id UUID REFERENCES hrms.worker_profiles(id) ON DELETE SET NULL,
    objective_id UUID REFERENCES hrms.objectives(id) ON DELETE SET NULL,
    relationship TEXT NOT NULL DEFAULT 'peer',
    feedback_type TEXT NOT NULL DEFAULT '360',
    status TEXT NOT NULL DEFAULT 'requested',
    is_anonymous BOOLEAN NOT NULL DEFAULT FALSE,
    visibility TEXT NOT NULL DEFAULT 'subject_manager_hr',
    due_date DATE,
    prompt TEXT,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    CONSTRAINT feedback_requests_relationship_chk CHECK (relationship IN ('manager', 'peer', 'direct_report', 'self', 'hr', 'client')),
    CONSTRAINT feedback_requests_type_chk CHECK (feedback_type IN ('360', 'project', 'general', 'okr', 'manager_review')),
    CONSTRAINT feedback_requests_status_chk CHECK (status IN ('requested', 'submitted', 'declined', 'expired', 'cancelled')),
    CONSTRAINT feedback_requests_visibility_chk CHECK (visibility IN ('subject_manager_hr', 'manager_hr', 'hr_only', 'subject_only')),
    CONSTRAINT feedback_requests_metadata_object_chk CHECK (jsonb_typeof(metadata) = 'object')
);

CREATE INDEX IF NOT EXISTS feedback_requests_subject_idx
    ON hrms.feedback_requests (tenant_id, subject_worker_profile_id, status, due_date)
    WHERE inactive = FALSE;

CREATE INDEX IF NOT EXISTS feedback_requests_requester_idx
    ON hrms.feedback_requests (tenant_id, requester_worker_profile_id, created_at DESC)
    WHERE requester_worker_profile_id IS NOT NULL AND inactive = FALSE;

CREATE TABLE IF NOT EXISTS hrms.feedback_responses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    request_id UUID NOT NULL REFERENCES hrms.feedback_requests(id) ON DELETE CASCADE,
    respondent_worker_profile_id UUID REFERENCES hrms.worker_profiles(id) ON DELETE SET NULL,
    rating NUMERIC(5,2),
    strengths TEXT,
    improvements TEXT,
    comments TEXT,
    submitted_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID,
    CONSTRAINT feedback_responses_rating_chk CHECK (rating IS NULL OR (rating >= 0 AND rating <= 5)),
    CONSTRAINT feedback_responses_metadata_object_chk CHECK (jsonb_typeof(metadata) = 'object')
);

CREATE UNIQUE INDEX IF NOT EXISTS feedback_responses_request_uq
    ON hrms.feedback_responses (tenant_id, request_id)
    WHERE inactive = FALSE;

CREATE INDEX IF NOT EXISTS feedback_responses_respondent_idx
    ON hrms.feedback_responses (tenant_id, respondent_worker_profile_id, submitted_at DESC)
    WHERE respondent_worker_profile_id IS NOT NULL AND inactive = FALSE;

CREATE TABLE IF NOT EXISTS hrms.performance_timeline_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    worker_profile_id UUID NOT NULL REFERENCES hrms.worker_profiles(id) ON DELETE CASCADE,
    event_type TEXT NOT NULL,
    checkin_id UUID REFERENCES hrms.performance_checkins(id) ON DELETE SET NULL,
    feedback_request_id UUID REFERENCES hrms.feedback_requests(id) ON DELETE SET NULL,
    feedback_response_id UUID REFERENCES hrms.feedback_responses(id) ON DELETE SET NULL,
    objective_id UUID REFERENCES hrms.objectives(id) ON DELETE SET NULL,
    actor_worker_profile_id UUID REFERENCES hrms.worker_profiles(id) ON DELETE SET NULL,
    title TEXT NOT NULL,
    notes TEXT,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID,
    CONSTRAINT performance_timeline_event_type_chk CHECK (event_type IN ('checkin_created', 'checkin_submitted', 'checkin_reviewed', 'feedback_requested', 'feedback_submitted', 'calibration_note', 'objective_update')),
    CONSTRAINT performance_timeline_metadata_object_chk CHECK (jsonb_typeof(metadata) = 'object')
);

CREATE INDEX IF NOT EXISTS performance_timeline_worker_idx
    ON hrms.performance_timeline_events (tenant_id, worker_profile_id, created_at DESC)
    WHERE inactive = FALSE;

DROP TRIGGER IF EXISTS trg_performance_checkins_updated_at ON hrms.performance_checkins;
CREATE TRIGGER trg_performance_checkins_updated_at
BEFORE UPDATE ON hrms.performance_checkins
FOR EACH ROW EXECUTE FUNCTION hrms.set_updated_at();

DROP TRIGGER IF EXISTS trg_feedback_requests_updated_at ON hrms.feedback_requests;
CREATE TRIGGER trg_feedback_requests_updated_at
BEFORE UPDATE ON hrms.feedback_requests
FOR EACH ROW EXECUTE FUNCTION hrms.set_updated_at();

ALTER TABLE hrms.performance_checkins ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.feedback_requests ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.feedback_responses ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.performance_timeline_events ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS performance_checkins_tenant_isolation ON hrms.performance_checkins;
CREATE POLICY performance_checkins_tenant_isolation ON hrms.performance_checkins
USING (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true')
WITH CHECK (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true');

DROP POLICY IF EXISTS feedback_requests_tenant_isolation ON hrms.feedback_requests;
CREATE POLICY feedback_requests_tenant_isolation ON hrms.feedback_requests
USING (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true')
WITH CHECK (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true');

DROP POLICY IF EXISTS feedback_responses_tenant_isolation ON hrms.feedback_responses;
CREATE POLICY feedback_responses_tenant_isolation ON hrms.feedback_responses
USING (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true')
WITH CHECK (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true');

DROP POLICY IF EXISTS performance_timeline_tenant_isolation ON hrms.performance_timeline_events;
CREATE POLICY performance_timeline_tenant_isolation ON hrms.performance_timeline_events
USING (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true')
WITH CHECK (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true');
