CREATE TABLE IF NOT EXISTS hrms.pulse_surveys (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    description TEXT,
    survey_type TEXT NOT NULL DEFAULT 'pulse',
    status TEXT NOT NULL DEFAULT 'draft',
    audience_scope TEXT NOT NULL DEFAULT 'all',
    department_id UUID REFERENCES hrms.departments(id) ON DELETE SET NULL,
    start_date DATE NOT NULL,
    end_date DATE,
    frequency TEXT NOT NULL DEFAULT 'one_time',
    anonymity_threshold INT NOT NULL DEFAULT 5,
    consent_required BOOLEAN NOT NULL DEFAULT TRUE,
    manager_aggregate_only BOOLEAN NOT NULL DEFAULT TRUE,
    critical_alerts_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    CONSTRAINT pulse_surveys_type_chk CHECK (survey_type IN ('pulse', 'wellbeing', 'engagement')),
    CONSTRAINT pulse_surveys_status_chk CHECK (status IN ('draft', 'active', 'closed', 'archived')),
    CONSTRAINT pulse_surveys_audience_chk CHECK (audience_scope IN ('all', 'department', 'custom')),
    CONSTRAINT pulse_surveys_frequency_chk CHECK (frequency IN ('one_time', 'weekly', 'biweekly', 'monthly')),
    CONSTRAINT pulse_surveys_threshold_chk CHECK (anonymity_threshold >= 3),
    CONSTRAINT pulse_surveys_dates_chk CHECK (end_date IS NULL OR end_date >= start_date),
    CONSTRAINT pulse_surveys_metadata_object_chk CHECK (jsonb_typeof(metadata) = 'object')
);

CREATE INDEX IF NOT EXISTS pulse_surveys_tenant_status_idx
    ON hrms.pulse_surveys (tenant_id, status, start_date DESC)
    WHERE inactive = FALSE;

CREATE TABLE IF NOT EXISTS hrms.pulse_survey_questions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    survey_id UUID NOT NULL REFERENCES hrms.pulse_surveys(id) ON DELETE CASCADE,
    question_text TEXT NOT NULL,
    question_type TEXT NOT NULL DEFAULT 'scale_1_5',
    category TEXT NOT NULL DEFAULT 'general',
    is_required BOOLEAN NOT NULL DEFAULT TRUE,
    sort_order INT NOT NULL DEFAULT 0,
    options JSONB NOT NULL DEFAULT '[]'::jsonb,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    CONSTRAINT pulse_questions_type_chk CHECK (question_type IN ('scale_1_5', 'text', 'boolean', 'single_choice')),
    CONSTRAINT pulse_questions_category_chk CHECK (category IN ('mood', 'workload', 'stress', 'belonging', 'manager_support', 'safety', 'general')),
    CONSTRAINT pulse_questions_options_array_chk CHECK (jsonb_typeof(options) = 'array')
);

CREATE INDEX IF NOT EXISTS pulse_questions_survey_idx
    ON hrms.pulse_survey_questions (tenant_id, survey_id, sort_order)
    WHERE inactive = FALSE;

CREATE TABLE IF NOT EXISTS hrms.pulse_responses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    survey_id UUID NOT NULL REFERENCES hrms.pulse_surveys(id) ON DELETE CASCADE,
    question_id UUID NOT NULL REFERENCES hrms.pulse_survey_questions(id) ON DELETE CASCADE,
    worker_profile_id UUID REFERENCES hrms.worker_profiles(id) ON DELETE SET NULL,
    response_date DATE NOT NULL DEFAULT CURRENT_DATE,
    score NUMERIC(5,2),
    text_response TEXT,
    boolean_response BOOLEAN,
    option_value TEXT,
    consent_given BOOLEAN NOT NULL DEFAULT FALSE,
    is_anonymous BOOLEAN NOT NULL DEFAULT TRUE,
    risk_level TEXT NOT NULL DEFAULT 'none',
    critical_alert BOOLEAN NOT NULL DEFAULT FALSE,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID,
    CONSTRAINT pulse_responses_score_chk CHECK (score IS NULL OR (score >= 0 AND score <= 5)),
    CONSTRAINT pulse_responses_risk_chk CHECK (risk_level IN ('none', 'low', 'medium', 'high', 'critical')),
    CONSTRAINT pulse_responses_metadata_object_chk CHECK (jsonb_typeof(metadata) = 'object')
);

CREATE INDEX IF NOT EXISTS pulse_responses_survey_idx
    ON hrms.pulse_responses (tenant_id, survey_id, response_date DESC)
    WHERE inactive = FALSE;

CREATE INDEX IF NOT EXISTS pulse_responses_worker_idx
    ON hrms.pulse_responses (tenant_id, worker_profile_id, response_date DESC)
    WHERE worker_profile_id IS NOT NULL AND inactive = FALSE;

CREATE TABLE IF NOT EXISTS hrms.wellbeing_scores (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    worker_profile_id UUID NOT NULL REFERENCES hrms.worker_profiles(id) ON DELETE CASCADE,
    score_date DATE NOT NULL DEFAULT CURRENT_DATE,
    source_survey_id UUID REFERENCES hrms.pulse_surveys(id) ON DELETE SET NULL,
    wellbeing_score NUMERIC(5,2) NOT NULL,
    mood_score NUMERIC(5,2),
    stress_score NUMERIC(5,2),
    workload_score NUMERIC(5,2),
    risk_level TEXT NOT NULL DEFAULT 'none',
    consent_scope TEXT NOT NULL DEFAULT 'aggregate',
    notes TEXT,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    CONSTRAINT wellbeing_scores_score_chk CHECK (wellbeing_score >= 0 AND wellbeing_score <= 100),
    CONSTRAINT wellbeing_scores_component_chk CHECK (
        (mood_score IS NULL OR (mood_score >= 0 AND mood_score <= 5))
        AND (stress_score IS NULL OR (stress_score >= 0 AND stress_score <= 5))
        AND (workload_score IS NULL OR (workload_score >= 0 AND workload_score <= 5))
    ),
    CONSTRAINT wellbeing_scores_risk_chk CHECK (risk_level IN ('none', 'low', 'medium', 'high', 'critical')),
    CONSTRAINT wellbeing_scores_consent_chk CHECK (consent_scope IN ('private', 'aggregate', 'hr_alert')),
    CONSTRAINT wellbeing_scores_metadata_object_chk CHECK (jsonb_typeof(metadata) = 'object')
);

CREATE UNIQUE INDEX IF NOT EXISTS wellbeing_scores_worker_date_uq
    ON hrms.wellbeing_scores (tenant_id, worker_profile_id, score_date)
    WHERE inactive = FALSE;

CREATE INDEX IF NOT EXISTS wellbeing_scores_tenant_risk_idx
    ON hrms.wellbeing_scores (tenant_id, risk_level, score_date DESC)
    WHERE inactive = FALSE;

CREATE TABLE IF NOT EXISTS hrms.wellbeing_alerts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    worker_profile_id UUID REFERENCES hrms.worker_profiles(id) ON DELETE SET NULL,
    survey_id UUID REFERENCES hrms.pulse_surveys(id) ON DELETE SET NULL,
    response_id UUID REFERENCES hrms.pulse_responses(id) ON DELETE SET NULL,
    alert_type TEXT NOT NULL,
    severity TEXT NOT NULL DEFAULT 'medium',
    status TEXT NOT NULL DEFAULT 'open',
    visible_to TEXT NOT NULL DEFAULT 'hr',
    message TEXT NOT NULL,
    resolution_note TEXT,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    CONSTRAINT wellbeing_alerts_type_chk CHECK (alert_type IN ('critical_response', 'low_score', 'consent_issue')),
    CONSTRAINT wellbeing_alerts_severity_chk CHECK (severity IN ('medium', 'high', 'critical')),
    CONSTRAINT wellbeing_alerts_status_chk CHECK (status IN ('open', 'acknowledged', 'resolved', 'dismissed')),
    CONSTRAINT wellbeing_alerts_visible_chk CHECK (visible_to = 'hr'),
    CONSTRAINT wellbeing_alerts_metadata_object_chk CHECK (jsonb_typeof(metadata) = 'object')
);

CREATE INDEX IF NOT EXISTS wellbeing_alerts_tenant_status_idx
    ON hrms.wellbeing_alerts (tenant_id, status, severity, created_at DESC)
    WHERE inactive = FALSE;

DROP TRIGGER IF EXISTS trg_pulse_surveys_updated_at ON hrms.pulse_surveys;
CREATE TRIGGER trg_pulse_surveys_updated_at
BEFORE UPDATE ON hrms.pulse_surveys
FOR EACH ROW EXECUTE FUNCTION hrms.set_updated_at();

DROP TRIGGER IF EXISTS trg_pulse_questions_updated_at ON hrms.pulse_survey_questions;
CREATE TRIGGER trg_pulse_questions_updated_at
BEFORE UPDATE ON hrms.pulse_survey_questions
FOR EACH ROW EXECUTE FUNCTION hrms.set_updated_at();

DROP TRIGGER IF EXISTS trg_wellbeing_scores_updated_at ON hrms.wellbeing_scores;
CREATE TRIGGER trg_wellbeing_scores_updated_at
BEFORE UPDATE ON hrms.wellbeing_scores
FOR EACH ROW EXECUTE FUNCTION hrms.set_updated_at();

DROP TRIGGER IF EXISTS trg_wellbeing_alerts_updated_at ON hrms.wellbeing_alerts;
CREATE TRIGGER trg_wellbeing_alerts_updated_at
BEFORE UPDATE ON hrms.wellbeing_alerts
FOR EACH ROW EXECUTE FUNCTION hrms.set_updated_at();

ALTER TABLE hrms.pulse_surveys ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.pulse_survey_questions ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.pulse_responses ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.wellbeing_scores ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.wellbeing_alerts ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS pulse_surveys_tenant_isolation ON hrms.pulse_surveys;
CREATE POLICY pulse_surveys_tenant_isolation ON hrms.pulse_surveys
USING (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true')
WITH CHECK (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true');

DROP POLICY IF EXISTS pulse_questions_tenant_isolation ON hrms.pulse_survey_questions;
CREATE POLICY pulse_questions_tenant_isolation ON hrms.pulse_survey_questions
USING (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true')
WITH CHECK (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true');

DROP POLICY IF EXISTS pulse_responses_tenant_isolation ON hrms.pulse_responses;
CREATE POLICY pulse_responses_tenant_isolation ON hrms.pulse_responses
USING (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true')
WITH CHECK (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true');

DROP POLICY IF EXISTS wellbeing_scores_tenant_isolation ON hrms.wellbeing_scores;
CREATE POLICY wellbeing_scores_tenant_isolation ON hrms.wellbeing_scores
USING (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true')
WITH CHECK (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true');

DROP POLICY IF EXISTS wellbeing_alerts_tenant_isolation ON hrms.wellbeing_alerts;
CREATE POLICY wellbeing_alerts_tenant_isolation ON hrms.wellbeing_alerts
USING (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true')
WITH CHECK (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true');
