-- HRMS-032 attendance policy, roster, and regularisation workflow expansion.

CREATE TABLE IF NOT EXISTS hrms.attendance_policies (
    id                              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id                       UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    name                            VARCHAR(255) NOT NULL,
    code                            VARCHAR(80) NOT NULL,
    description                     TEXT,
    branch_id                       UUID REFERENCES hrms.branches(id) ON DELETE CASCADE,
    department_id                   UUID REFERENCES hrms.departments(id) ON DELETE CASCADE,
    user_id                         UUID REFERENCES auth.users(id) ON DELETE CASCADE,
    schedule_type                   VARCHAR(20) NOT NULL DEFAULT 'fixed',
    is_default                      BOOLEAN NOT NULL DEFAULT FALSE,
    standard_work_minutes           INT NOT NULL DEFAULT 480,
    min_half_day_minutes            INT NOT NULL DEFAULT 240,
    min_full_day_minutes            INT NOT NULL DEFAULT 420,
    grace_late_minutes              INT NOT NULL DEFAULT 0,
    grace_early_minutes             INT NOT NULL DEFAULT 0,
    half_day_late_after_minutes     INT,
    absent_late_after_minutes       INT,
    half_day_early_before_minutes   INT,
    absent_early_before_minutes     INT,
    allow_flexi_hours               BOOLEAN NOT NULL DEFAULT FALSE,
    core_start_time                 TIME,
    core_end_time                   TIME,
    allow_wfh                       BOOLEAN NOT NULL DEFAULT FALSE,
    wfh_days_per_week               INT NOT NULL DEFAULT 0,
    allow_permanent_remote          BOOLEAN NOT NULL DEFAULT FALSE,
    require_geo                     BOOLEAN NOT NULL DEFAULT FALSE,
    require_device                  BOOLEAN NOT NULL DEFAULT FALSE,
    regularization_window_days      INT NOT NULL DEFAULT 7,
    max_regularizations_per_month   INT NOT NULL DEFAULT 3,
    approval_mode                   VARCHAR(20) NOT NULL DEFAULT 'manager',
    effective_from                  DATE,
    effective_to                    DATE,
    inactive                        BOOLEAN NOT NULL DEFAULT FALSE,
    created_at                      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by                      UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at                      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by                      UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT attendance_policies_schedule_type_check CHECK (schedule_type IN ('fixed','flexi','daily_roster','weekly_roster','monthly_roster')),
    CONSTRAINT attendance_policies_approval_mode_check CHECK (approval_mode IN ('manager','hr','manager_hr','auto')),
    CONSTRAINT attendance_policies_minutes_check CHECK (
        standard_work_minutes >= 0 AND min_half_day_minutes >= 0 AND min_full_day_minutes >= 0
        AND grace_late_minutes >= 0 AND grace_early_minutes >= 0 AND wfh_days_per_week >= 0
        AND regularization_window_days >= 0 AND max_regularizations_per_month >= 0
    ),
    CONSTRAINT attendance_policies_scope_check CHECK (
        (CASE WHEN branch_id IS NULL THEN 0 ELSE 1 END)
        + (CASE WHEN department_id IS NULL THEN 0 ELSE 1 END)
        + (CASE WHEN user_id IS NULL THEN 0 ELSE 1 END) <= 1
    ),
    CONSTRAINT attendance_policies_effective_order_check CHECK (effective_to IS NULL OR effective_from IS NULL OR effective_from <= effective_to)
);
CREATE INDEX IF NOT EXISTS attendance_policies_tenant_idx ON hrms.attendance_policies(tenant_id);
CREATE UNIQUE INDEX IF NOT EXISTS attendance_policies_code_idx ON hrms.attendance_policies(tenant_id, lower(code)) WHERE NOT inactive;
CREATE UNIQUE INDEX IF NOT EXISTS attendance_policies_default_idx ON hrms.attendance_policies(tenant_id) WHERE is_default AND NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.attendance_rosters (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    user_id         UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    policy_id       UUID REFERENCES hrms.attendance_policies(id) ON DELETE SET NULL,
    date            DATE NOT NULL,
    start_time      TIME,
    end_time        TIME,
    break_minutes   INT NOT NULL DEFAULT 0,
    work_mode       VARCHAR(50) NOT NULL DEFAULT 'office',
    location_type   VARCHAR(50) NOT NULL DEFAULT 'office',
    remarks         TEXT,
    inactive        BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by      UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by      UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT attendance_rosters_time_order_check CHECK (end_time IS NULL OR start_time IS NULL OR end_time > start_time),
    CONSTRAINT attendance_rosters_break_check CHECK (break_minutes >= 0),
    CONSTRAINT attendance_rosters_work_mode_check CHECK (work_mode IN ('office','remote','field','hybrid')),
    CONSTRAINT attendance_rosters_location_type_check CHECK (location_type IN ('office','remote','field','client_site','hybrid'))
);
CREATE INDEX IF NOT EXISTS attendance_rosters_tenant_date_idx ON hrms.attendance_rosters(tenant_id, date) WHERE NOT inactive;
CREATE UNIQUE INDEX IF NOT EXISTS attendance_rosters_user_date_idx ON hrms.attendance_rosters(tenant_id, user_id, date) WHERE NOT inactive;

ALTER TABLE hrms.attendance_requests
    ADD COLUMN IF NOT EXISTS request_type VARCHAR(40) NOT NULL DEFAULT 'regularization',
    ADD COLUMN IF NOT EXISTS requested_checkin_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS requested_checkout_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS requested_work_mode VARCHAR(50),
    ADD COLUMN IF NOT EXISTS policy_id UUID REFERENCES hrms.attendance_policies(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS roster_id UUID REFERENCES hrms.attendance_rosters(id) ON DELETE SET NULL;

ALTER TABLE hrms.attendance_requests DROP CONSTRAINT IF EXISTS attendance_requests_type_check;
ALTER TABLE hrms.attendance_requests DROP CONSTRAINT IF EXISTS attendance_requests_request_type_check;
ALTER TABLE hrms.attendance_requests ADD CONSTRAINT attendance_requests_request_type_check CHECK (request_type IN ('regularization','missed_punch','late_exemption','early_exit_exemption','wfh','remote_work','halfday','absent','overtime'));
ALTER TABLE hrms.attendance_requests ADD CONSTRAINT attendance_requests_type_check CHECK (requested_type IS NULL OR requested_type IN ('checkin','checkout','missed_punch','halfday','absent','present','wfh','remote_work','overtime'));
ALTER TABLE hrms.attendance_requests DROP CONSTRAINT IF EXISTS attendance_requests_work_mode_check;
ALTER TABLE hrms.attendance_requests ADD CONSTRAINT attendance_requests_work_mode_check CHECK (requested_work_mode IS NULL OR requested_work_mode IN ('office','remote','field','hybrid'));
CREATE INDEX IF NOT EXISTS attendance_requests_tenant_status_idx ON hrms.attendance_requests(tenant_id, status) WHERE NOT inactive;
CREATE INDEX IF NOT EXISTS attendance_requests_tenant_date_idx ON hrms.attendance_requests(tenant_id, date) WHERE NOT inactive;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'attendance_policies_updated_at') THEN
        CREATE TRIGGER attendance_policies_updated_at BEFORE UPDATE ON hrms.attendance_policies FOR EACH ROW EXECUTE FUNCTION hrms.update_updated_at();
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'attendance_rosters_updated_at') THEN
        CREATE TRIGGER attendance_rosters_updated_at BEFORE UPDATE ON hrms.attendance_rosters FOR EACH ROW EXECUTE FUNCTION hrms.update_updated_at();
    END IF;
END $$;

ALTER TABLE hrms.attendance_policies ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.attendance_rosters ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS attendance_policies_tenant_isolation ON hrms.attendance_policies;
CREATE POLICY attendance_policies_tenant_isolation ON hrms.attendance_policies
USING (tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::uuid OR current_setting('app.is_super_admin', true) = 'true')
WITH CHECK (tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::uuid OR current_setting('app.is_super_admin', true) = 'true');

DROP POLICY IF EXISTS attendance_rosters_tenant_isolation ON hrms.attendance_rosters;
CREATE POLICY attendance_rosters_tenant_isolation ON hrms.attendance_rosters
USING (tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::uuid OR current_setting('app.is_super_admin', true) = 'true')
WITH CHECK (tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::uuid OR current_setting('app.is_super_admin', true) = 'true');
