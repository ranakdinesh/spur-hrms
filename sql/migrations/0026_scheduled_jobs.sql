CREATE TABLE IF NOT EXISTS hrms.job_locks (
    tenant_id       UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    job_key         VARCHAR(120) NOT NULL,
    locked_until    TIMESTAMPTZ NOT NULL,
    owner_id        VARCHAR(120) NOT NULL,
    last_acquired_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (tenant_id, job_key)
);

CREATE TABLE IF NOT EXISTS hrms.job_runs (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    job_key         VARCHAR(120) NOT NULL,
    run_date        DATE NOT NULL,
    status          VARCHAR(30) NOT NULL DEFAULT 'running',
    owner_id        VARCHAR(120),
    started_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    finished_at     TIMESTAMPTZ,
    processed_count INTEGER NOT NULL DEFAULT 0,
    success_count   INTEGER NOT NULL DEFAULT 0,
    failed_count    INTEGER NOT NULL DEFAULT 0,
    skipped_count   INTEGER NOT NULL DEFAULT 0,
    error_message   TEXT,
    metadata        JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive        BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by      UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by      UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT job_runs_status_check CHECK (status IN ('running','succeeded','failed','skipped'))
);

CREATE UNIQUE INDEX IF NOT EXISTS job_runs_tenant_job_date_idx ON hrms.job_runs(tenant_id, job_key, run_date) WHERE NOT inactive;
CREATE INDEX IF NOT EXISTS job_runs_tenant_job_idx ON hrms.job_runs(tenant_id, job_key, started_at DESC) WHERE NOT inactive;

ALTER TABLE hrms.job_locks ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.job_runs ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS job_locks_tenant_isolation ON hrms.job_locks;
CREATE POLICY job_locks_tenant_isolation ON hrms.job_locks
USING (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true')
WITH CHECK (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true');

DROP POLICY IF EXISTS job_runs_tenant_isolation ON hrms.job_runs;
CREATE POLICY job_runs_tenant_isolation ON hrms.job_runs
USING (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true')
WITH CHECK (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true');
