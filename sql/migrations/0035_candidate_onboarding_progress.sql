ALTER TABLE hrms.candidate_onboardings
    ADD COLUMN IF NOT EXISTS started_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS completed_at TIMESTAMPTZ;

ALTER TABLE hrms.candidate_onboarding_tasks
    ADD COLUMN IF NOT EXISTS due_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS started_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS completed_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS reviewed_by UUID REFERENCES auth.users(id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS candidate_onboardings_status_idx
    ON hrms.candidate_onboardings(tenant_id, onboarding_status)
    WHERE NOT inactive;

CREATE INDEX IF NOT EXISTS candidate_onboarding_tasks_due_idx
    ON hrms.candidate_onboarding_tasks(tenant_id, due_at, status)
    WHERE NOT inactive;

CREATE UNIQUE INDEX IF NOT EXISTS candidate_onboarding_tasks_unique_task_idx
    ON hrms.candidate_onboarding_tasks(tenant_id, candidate_onboarding_id, onboarding_task_id)
    WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.candidate_onboarding_events (
    id                          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id                   UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    candidate_onboarding_id     UUID NOT NULL REFERENCES hrms.candidate_onboardings(id) ON DELETE CASCADE,
    candidate_onboarding_task_id UUID REFERENCES hrms.candidate_onboarding_tasks(id) ON DELETE SET NULL,
    action                      VARCHAR(80) NOT NULL,
    from_status                 VARCHAR(50),
    to_status                   VARCHAR(50),
    remarks                     TEXT,
    metadata                    JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive                    BOOLEAN NOT NULL DEFAULT FALSE,
    created_at                  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by                  UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at                  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by                  UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT candidate_onboarding_events_action_check CHECK (length(trim(action)) > 0)
);

CREATE INDEX IF NOT EXISTS candidate_onboarding_events_onboarding_idx
    ON hrms.candidate_onboarding_events(tenant_id, candidate_onboarding_id, created_at DESC)
    WHERE NOT inactive;

ALTER TABLE hrms.candidate_onboarding_events ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS tenant_isolation_candidate_onboarding_events ON hrms.candidate_onboarding_events;
CREATE POLICY tenant_isolation_candidate_onboarding_events ON hrms.candidate_onboarding_events
    USING (
        current_setting('app.is_super_admin', true) = 'true'
        OR tenant_id::text = current_setting('app.tenant_id', true)
    )
    WITH CHECK (
        current_setting('app.is_super_admin', true) = 'true'
        OR tenant_id::text = current_setting('app.tenant_id', true)
    );
