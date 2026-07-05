-- Advanced attendance exception workflows and payroll blockers.

CREATE TABLE IF NOT EXISTS hrms.attendance_exception_workflows (
    id                          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id                   UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    code                        VARCHAR(100) NOT NULL,
    name                        VARCHAR(180) NOT NULL,
    description                 TEXT,
    branch_id                   UUID REFERENCES hrms.branches(id) ON DELETE CASCADE,
    department_id               UUID REFERENCES hrms.departments(id) ON DELETE CASCADE,
    request_type                VARCHAR(60) NOT NULL,
    route_mode                  VARCHAR(30) NOT NULL DEFAULT 'manager',
    max_requests_per_month      INT NOT NULL DEFAULT 0,
    escalation_hours            INT NOT NULL DEFAULT 0,
    escalation_route_mode       VARCHAR(30),
    block_payroll_when_pending  BOOLEAN NOT NULL DEFAULT TRUE,
    auto_approve                BOOLEAN NOT NULL DEFAULT FALSE,
    is_active                   BOOLEAN NOT NULL DEFAULT TRUE,
    inactive                    BOOLEAN NOT NULL DEFAULT FALSE,
    created_at                  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by                  UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at                  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by                  UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT attendance_exception_workflows_code_check CHECK (length(trim(code)) > 0),
    CONSTRAINT attendance_exception_workflows_name_check CHECK (length(trim(name)) > 0),
    CONSTRAINT attendance_exception_workflows_scope_check CHECK (
        (CASE WHEN branch_id IS NULL THEN 0 ELSE 1 END)
        + (CASE WHEN department_id IS NULL THEN 0 ELSE 1 END) <= 1
    ),
    CONSTRAINT attendance_exception_workflows_request_type_check CHECK (request_type IN ('regularization','missed_punch','late_exemption','early_exit_exemption','wfh','remote_work','halfday','absent','overtime')),
    CONSTRAINT attendance_exception_workflows_route_mode_check CHECK (route_mode IN ('manager','hr','manager_hr','auto')),
    CONSTRAINT attendance_exception_workflows_escalation_route_check CHECK (escalation_route_mode IS NULL OR escalation_route_mode IN ('manager','hr','manager_hr','auto')),
    CONSTRAINT attendance_exception_workflows_limits_check CHECK (max_requests_per_month >= 0 AND escalation_hours >= 0)
);

CREATE UNIQUE INDEX IF NOT EXISTS attendance_exception_workflows_code_idx
    ON hrms.attendance_exception_workflows(tenant_id, lower(code))
    WHERE NOT inactive;

CREATE INDEX IF NOT EXISTS attendance_exception_workflows_resolve_idx
    ON hrms.attendance_exception_workflows(tenant_id, request_type, branch_id, department_id, is_active)
    WHERE NOT inactive;

ALTER TABLE hrms.attendance_requests
    ADD COLUMN IF NOT EXISTS workflow_id UUID REFERENCES hrms.attendance_exception_workflows(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS route_mode VARCHAR(30),
    ADD COLUMN IF NOT EXISTS escalation_due_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS escalated_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS payroll_blocking BOOLEAN NOT NULL DEFAULT FALSE;

ALTER TABLE hrms.attendance_requests DROP CONSTRAINT IF EXISTS attendance_requests_route_mode_check;
ALTER TABLE hrms.attendance_requests ADD CONSTRAINT attendance_requests_route_mode_check CHECK (route_mode IS NULL OR route_mode IN ('manager','hr','manager_hr','auto'));

CREATE INDEX IF NOT EXISTS attendance_requests_workflow_idx
    ON hrms.attendance_requests(tenant_id, workflow_id, status)
    WHERE workflow_id IS NOT NULL AND NOT inactive;

CREATE INDEX IF NOT EXISTS attendance_requests_payroll_blocking_idx
    ON hrms.attendance_requests(tenant_id, date, status)
    WHERE payroll_blocking AND NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.attendance_exception_events (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id               UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    attendance_request_id   UUID NOT NULL REFERENCES hrms.attendance_requests(id) ON DELETE CASCADE,
    workflow_id             UUID REFERENCES hrms.attendance_exception_workflows(id) ON DELETE SET NULL,
    action                  VARCHAR(60) NOT NULL,
    from_status             VARCHAR(30),
    to_status               VARCHAR(30),
    routed_to               VARCHAR(30),
    remarks                 TEXT,
    metadata                JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive                BOOLEAN NOT NULL DEFAULT FALSE,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by              UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by              UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT attendance_exception_events_action_check CHECK (length(trim(action)) > 0),
    CONSTRAINT attendance_exception_events_metadata_object_check CHECK (jsonb_typeof(metadata) = 'object')
);

CREATE INDEX IF NOT EXISTS attendance_exception_events_request_idx
    ON hrms.attendance_exception_events(tenant_id, attendance_request_id, created_at DESC)
    WHERE NOT inactive;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'attendance_exception_workflows_updated_at') THEN
        CREATE TRIGGER attendance_exception_workflows_updated_at BEFORE UPDATE ON hrms.attendance_exception_workflows FOR EACH ROW EXECUTE FUNCTION hrms.update_updated_at();
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'attendance_exception_events_updated_at') THEN
        CREATE TRIGGER attendance_exception_events_updated_at BEFORE UPDATE ON hrms.attendance_exception_events FOR EACH ROW EXECUTE FUNCTION hrms.update_updated_at();
    END IF;
END $$;

ALTER TABLE hrms.attendance_exception_workflows ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.attendance_exception_events ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS attendance_exception_workflows_tenant_isolation ON hrms.attendance_exception_workflows;
CREATE POLICY attendance_exception_workflows_tenant_isolation ON hrms.attendance_exception_workflows
USING (tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::uuid OR current_setting('app.is_super_admin', true) = 'true')
WITH CHECK (tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::uuid OR current_setting('app.is_super_admin', true) = 'true');

DROP POLICY IF EXISTS attendance_exception_events_tenant_isolation ON hrms.attendance_exception_events;
CREATE POLICY attendance_exception_events_tenant_isolation ON hrms.attendance_exception_events
USING (tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::uuid OR current_setting('app.is_super_admin', true) = 'true')
WITH CHECK (tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::uuid OR current_setting('app.is_super_admin', true) = 'true');
