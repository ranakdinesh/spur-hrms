ALTER TABLE hrms.candidate_applications
    ADD COLUMN IF NOT EXISTS applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ADD COLUMN IF NOT EXISTS status_changed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ADD COLUMN IF NOT EXISTS status_changed_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS rejection_reason TEXT,
    ADD COLUMN IF NOT EXISTS withdrawal_reason TEXT,
    ADD COLUMN IF NOT EXISTS source_detail TEXT,
    ADD COLUMN IF NOT EXISTS duplicate_of_application_id UUID REFERENCES hrms.candidate_applications(id) ON DELETE SET NULL;

UPDATE hrms.candidate_applications
SET applied_at = COALESCE(applied_at, created_at),
    status_changed_at = COALESCE(status_changed_at, updated_at, created_at)
WHERE applied_at IS NULL OR status_changed_at IS NULL;

CREATE UNIQUE INDEX IF NOT EXISTS candidate_applications_active_unique_idx
    ON hrms.candidate_applications(tenant_id, candidate_id, job_posting_id)
    WHERE candidate_id IS NOT NULL AND job_posting_id IS NOT NULL AND NOT inactive;

CREATE INDEX IF NOT EXISTS candidate_applications_status_changed_idx
    ON hrms.candidate_applications(tenant_id, status, status_changed_at)
    WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.candidate_application_events (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id      UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    application_id UUID NOT NULL REFERENCES hrms.candidate_applications(id) ON DELETE CASCADE,
    from_status    VARCHAR(50),
    to_status      VARCHAR(50) NOT NULL,
    action         VARCHAR(50) NOT NULL,
    reason         TEXT,
    remarks        TEXT,
    inactive       BOOLEAN NOT NULL DEFAULT FALSE,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by     UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by     UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT candidate_application_events_to_status_check CHECK (to_status IN ('New','Screening','Interview','Offered','Hired','Rejected','Withdrawn'))
);

CREATE INDEX IF NOT EXISTS candidate_application_events_application_idx
    ON hrms.candidate_application_events(application_id, created_at DESC)
    WHERE NOT inactive;

DO $$
BEGIN
    IF to_regclass('hrms.candidate_application_events') IS NOT NULL THEN
        ALTER TABLE hrms.candidate_application_events ENABLE ROW LEVEL SECURITY;

        IF NOT EXISTS (
            SELECT 1 FROM pg_policies
            WHERE schemaname = 'hrms'
              AND tablename = 'candidate_application_events'
              AND policyname = 'candidate_application_events_tenant_isolation'
        ) THEN
            CREATE POLICY candidate_application_events_tenant_isolation
                ON hrms.candidate_application_events
                USING (tenant_id = current_setting('app.tenant_id', true)::uuid OR current_setting('app.is_super_admin', true)::boolean)
                WITH CHECK (tenant_id = current_setting('app.tenant_id', true)::uuid OR current_setting('app.is_super_admin', true)::boolean);
        END IF;
    END IF;
END $$;
