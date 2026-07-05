CREATE OR REPLACE FUNCTION hrms.set_updated_at()
RETURNS TRIGGER LANGUAGE plpgsql AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$;

CREATE TABLE IF NOT EXISTS hrms.work_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    engagement_id UUID NOT NULL REFERENCES hrms.engagements(id) ON DELETE RESTRICT,
    worker_profile_id UUID NOT NULL REFERENCES hrms.worker_profiles(id) ON DELETE RESTRICT,
    log_date DATE NOT NULL,
    hours_worked NUMERIC(8,2) NOT NULL,
    billable_hours NUMERIC(8,2),
    work_summary TEXT,
    deliverable_reference TEXT,
    status TEXT NOT NULL DEFAULT 'draft',
    submitted_at TIMESTAMPTZ,
    submitted_by UUID,
    reviewed_at TIMESTAMPTZ,
    reviewed_by UUID,
    review_comment TEXT,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    CONSTRAINT work_logs_status_check CHECK (status IN ('draft', 'submitted', 'approved', 'rejected', 'cancelled')),
    CONSTRAINT work_logs_hours_check CHECK (hours_worked > 0 AND hours_worked <= 24),
    CONSTRAINT work_logs_billable_hours_check CHECK (billable_hours IS NULL OR (billable_hours >= 0 AND billable_hours <= hours_worked)),
    CONSTRAINT work_logs_metadata_object_check CHECK (jsonb_typeof(metadata) = 'object')
);

CREATE UNIQUE INDEX IF NOT EXISTS work_logs_tenant_engagement_worker_date_active_uq
    ON hrms.work_logs (tenant_id, engagement_id, worker_profile_id, log_date)
    WHERE inactive = FALSE;

CREATE INDEX IF NOT EXISTS work_logs_tenant_status_date_idx
    ON hrms.work_logs (tenant_id, status, log_date DESC)
    WHERE inactive = FALSE;

CREATE INDEX IF NOT EXISTS work_logs_tenant_worker_date_idx
    ON hrms.work_logs (tenant_id, worker_profile_id, log_date DESC)
    WHERE inactive = FALSE;

CREATE INDEX IF NOT EXISTS work_logs_tenant_engagement_status_idx
    ON hrms.work_logs (tenant_id, engagement_id, status)
    WHERE inactive = FALSE;

DROP TRIGGER IF EXISTS trg_work_logs_updated_at ON hrms.work_logs;
CREATE TRIGGER trg_work_logs_updated_at
BEFORE UPDATE ON hrms.work_logs
FOR EACH ROW
EXECUTE FUNCTION hrms.set_updated_at();

ALTER TABLE hrms.work_logs ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS work_logs_tenant_isolation ON hrms.work_logs;
CREATE POLICY work_logs_tenant_isolation ON hrms.work_logs
USING (
    current_setting('app.is_super_admin', true) = 'true'
    OR tenant_id::text = current_setting('app.tenant_id', true)
)
WITH CHECK (
    current_setting('app.is_super_admin', true) = 'true'
    OR tenant_id::text = current_setting('app.tenant_id', true)
);
