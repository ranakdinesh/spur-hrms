CREATE TABLE IF NOT EXISTS hrms.employee_exit_requests (
    id                          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id                   UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    employee_id                 UUID NOT NULL REFERENCES hrms.employees(id) ON DELETE CASCADE,
    employee_user_id            UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    initiated_by                UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    approved_by                 UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    approved_at                 TIMESTAMPTZ,
    completed_by                UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    completed_at                TIMESTAMPTZ,
    status                      VARCHAR(40) NOT NULL DEFAULT 'submitted',
    exit_type                   VARCHAR(40) NOT NULL DEFAULT 'resignation',
    reason                      TEXT,
    resignation_date            DATE,
    notice_start_date           DATE,
    last_working_date           DATE NOT NULL,
    requested_relieving_date    DATE,
    approved_relieving_date     DATE,
    final_settlement_status     VARCHAR(40) NOT NULL DEFAULT 'pending',
    access_revocation_status    VARCHAR(40) NOT NULL DEFAULT 'pending',
    asset_clearance_status      VARCHAR(40) NOT NULL DEFAULT 'pending',
    handover_status             VARCHAR(40) NOT NULL DEFAULT 'pending',
    exit_interview_status       VARCHAR(40) NOT NULL DEFAULT 'pending',
    notes                       TEXT,
    rejection_reason            TEXT,
    cancel_reason               TEXT,
    inactive                    BOOLEAN NOT NULL DEFAULT FALSE,
    created_at                  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by                  UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at                  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by                  UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT employee_exit_requests_status_check CHECK (status IN ('submitted','approved','rejected','completed','canceled')),
    CONSTRAINT employee_exit_requests_type_check CHECK (exit_type IN ('resignation','termination','retirement','contract_end','absconding','other')),
    CONSTRAINT employee_exit_requests_summary_check CHECK (
        final_settlement_status IN ('pending','in_progress','ready','paid','on_hold')
        AND access_revocation_status IN ('pending','scheduled','revoked','not_required')
        AND asset_clearance_status IN ('pending','partial','cleared','not_required')
        AND handover_status IN ('pending','in_progress','completed','not_required')
        AND exit_interview_status IN ('pending','scheduled','completed','skipped')
    ),
    CONSTRAINT employee_exit_requests_dates_check CHECK (
        notice_start_date IS NULL OR last_working_date >= notice_start_date
    )
);

CREATE UNIQUE INDEX IF NOT EXISTS employee_exit_requests_one_active_idx
    ON hrms.employee_exit_requests(tenant_id, employee_user_id)
    WHERE NOT inactive AND status IN ('submitted','approved');

CREATE INDEX IF NOT EXISTS employee_exit_requests_tenant_status_idx
    ON hrms.employee_exit_requests(tenant_id, status, last_working_date)
    WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.employee_exit_tasks (
    id                          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id                   UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    exit_request_id             UUID NOT NULL REFERENCES hrms.employee_exit_requests(id) ON DELETE CASCADE,
    employee_user_id            UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    task_key                    VARCHAR(80) NOT NULL,
    title                       VARCHAR(160) NOT NULL,
    description                 TEXT,
    owner_role                  VARCHAR(80),
    owner_user_id               UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    due_date                    DATE,
    status                      VARCHAR(40) NOT NULL DEFAULT 'pending',
    completed_by                UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    completed_at                TIMESTAMPTZ,
    remarks                     TEXT,
    sort_order                  INTEGER NOT NULL DEFAULT 0,
    inactive                    BOOLEAN NOT NULL DEFAULT FALSE,
    created_at                  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by                  UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at                  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by                  UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT employee_exit_tasks_status_check CHECK (status IN ('pending','in_progress','completed','waived','blocked')),
    CONSTRAINT employee_exit_tasks_title_check CHECK (length(trim(title)) > 0)
);

CREATE UNIQUE INDEX IF NOT EXISTS employee_exit_tasks_unique_key_idx
    ON hrms.employee_exit_tasks(tenant_id, exit_request_id, task_key)
    WHERE NOT inactive;

CREATE INDEX IF NOT EXISTS employee_exit_tasks_request_idx
    ON hrms.employee_exit_tasks(tenant_id, exit_request_id, sort_order)
    WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.employee_exit_events (
    id                          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id                   UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    exit_request_id             UUID NOT NULL REFERENCES hrms.employee_exit_requests(id) ON DELETE CASCADE,
    exit_task_id                UUID REFERENCES hrms.employee_exit_tasks(id) ON DELETE SET NULL,
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
    CONSTRAINT employee_exit_events_action_check CHECK (length(trim(action)) > 0)
);

CREATE INDEX IF NOT EXISTS employee_exit_events_request_idx
    ON hrms.employee_exit_events(tenant_id, exit_request_id, created_at DESC)
    WHERE NOT inactive;

ALTER TABLE hrms.employee_exit_requests ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.employee_exit_tasks ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.employee_exit_events ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS tenant_isolation_employee_exit_requests ON hrms.employee_exit_requests;
CREATE POLICY tenant_isolation_employee_exit_requests ON hrms.employee_exit_requests
    USING (
        current_setting('app.is_super_admin', true) = 'true'
        OR tenant_id::text = current_setting('app.tenant_id', true)
    )
    WITH CHECK (
        current_setting('app.is_super_admin', true) = 'true'
        OR tenant_id::text = current_setting('app.tenant_id', true)
    );

DROP POLICY IF EXISTS tenant_isolation_employee_exit_tasks ON hrms.employee_exit_tasks;
CREATE POLICY tenant_isolation_employee_exit_tasks ON hrms.employee_exit_tasks
    USING (
        current_setting('app.is_super_admin', true) = 'true'
        OR tenant_id::text = current_setting('app.tenant_id', true)
    )
    WITH CHECK (
        current_setting('app.is_super_admin', true) = 'true'
        OR tenant_id::text = current_setting('app.tenant_id', true)
    );

DROP POLICY IF EXISTS tenant_isolation_employee_exit_events ON hrms.employee_exit_events;
CREATE POLICY tenant_isolation_employee_exit_events ON hrms.employee_exit_events
    USING (
        current_setting('app.is_super_admin', true) = 'true'
        OR tenant_id::text = current_setting('app.tenant_id', true)
    )
    WITH CHECK (
        current_setting('app.is_super_admin', true) = 'true'
        OR tenant_id::text = current_setting('app.tenant_id', true)
    );
