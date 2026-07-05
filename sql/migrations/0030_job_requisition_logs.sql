CREATE TABLE IF NOT EXISTS hrms.job_requisition_logs (
    id                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id          UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    job_requisition_id UUID NOT NULL REFERENCES hrms.job_requisitions(id) ON DELETE CASCADE,
    from_status        VARCHAR(50),
    to_status          VARCHAR(50) NOT NULL,
    action             VARCHAR(50) NOT NULL,
    remarks            TEXT,
    inactive           BOOLEAN NOT NULL DEFAULT FALSE,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by         UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by         UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT job_requisition_logs_action_check CHECK (action IN ('Created','Submitted','Approved','Rejected','Closed','Updated','Reopened')),
    CONSTRAINT job_requisition_logs_status_check CHECK (to_status IN ('Draft','Pending','Approved','Rejected','Closed'))
);

CREATE INDEX IF NOT EXISTS job_requisition_logs_requisition_idx ON hrms.job_requisition_logs(tenant_id, job_requisition_id, created_at DESC) WHERE NOT inactive;

DROP TRIGGER IF EXISTS job_requisition_logs_updated_at ON hrms.job_requisition_logs;
CREATE TRIGGER job_requisition_logs_updated_at
BEFORE UPDATE ON hrms.job_requisition_logs
FOR EACH ROW EXECUTE FUNCTION hrms.update_updated_at();

ALTER TABLE hrms.job_requisition_logs ENABLE ROW LEVEL SECURITY;
DROP POLICY IF EXISTS job_requisition_logs_tenant_isolation ON hrms.job_requisition_logs;
CREATE POLICY job_requisition_logs_tenant_isolation ON hrms.job_requisition_logs
USING (
    tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::uuid
    OR current_setting('app.is_super_admin', true) = 'true'
)
WITH CHECK (
    tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::uuid
    OR current_setting('app.is_super_admin', true) = 'true'
);
