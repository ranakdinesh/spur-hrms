CREATE TABLE IF NOT EXISTS hrms.shift_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    code TEXT NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    break_minutes INTEGER NOT NULL DEFAULT 0,
    paid_minutes INTEGER NOT NULL DEFAULT 0,
    work_mode TEXT NOT NULL DEFAULT 'office',
    location_type TEXT NOT NULL DEFAULT 'office',
    attendance_policy_id UUID REFERENCES hrms.attendance_policies(id) ON DELETE SET NULL,
    attendance_location_id UUID REFERENCES hrms.attendance_locations(id) ON DELETE SET NULL,
    allow_overtime BOOLEAN NOT NULL DEFAULT FALSE,
    payroll_code TEXT,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    CONSTRAINT shift_templates_time_order_check CHECK (end_time > start_time),
    CONSTRAINT shift_templates_break_check CHECK (break_minutes >= 0),
    CONSTRAINT shift_templates_paid_minutes_check CHECK (paid_minutes >= 0),
    CONSTRAINT shift_templates_work_mode_check CHECK (work_mode IN ('office','remote','field','hybrid')),
    CONSTRAINT shift_templates_location_type_check CHECK (location_type IN ('office','remote','field','client_site','hybrid','branch','warehouse','project_site','other'))
);

CREATE UNIQUE INDEX IF NOT EXISTS shift_templates_tenant_code_idx
    ON hrms.shift_templates(tenant_id, lower(code))
    WHERE NOT inactive;
CREATE INDEX IF NOT EXISTS shift_templates_tenant_active_idx
    ON hrms.shift_templates(tenant_id, is_active)
    WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.staffing_requirements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    name TEXT NOT NULL,
    requirement_date DATE,
    start_date DATE,
    end_date DATE,
    day_of_week INTEGER,
    branch_id UUID REFERENCES hrms.branches(id) ON DELETE SET NULL,
    department_id UUID REFERENCES hrms.departments(id) ON DELETE SET NULL,
    attendance_location_id UUID REFERENCES hrms.attendance_locations(id) ON DELETE SET NULL,
    role_label TEXT,
    team_label TEXT,
    shift_template_id UUID REFERENCES hrms.shift_templates(id) ON DELETE SET NULL,
    required_count INTEGER NOT NULL DEFAULT 1,
    min_count INTEGER NOT NULL DEFAULT 0,
    max_count INTEGER,
    priority TEXT NOT NULL DEFAULT 'medium',
    status TEXT NOT NULL DEFAULT 'active',
    payroll_blocking BOOLEAN NOT NULL DEFAULT FALSE,
    notes TEXT,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    CONSTRAINT staffing_requirements_day_check CHECK (day_of_week IS NULL OR day_of_week BETWEEN 0 AND 6),
    CONSTRAINT staffing_requirements_count_check CHECK (required_count >= 0 AND min_count >= 0 AND (max_count IS NULL OR max_count >= min_count)),
    CONSTRAINT staffing_requirements_priority_check CHECK (priority IN ('low','medium','high','critical')),
    CONSTRAINT staffing_requirements_status_check CHECK (status IN ('active','paused','archived'))
);

CREATE INDEX IF NOT EXISTS staffing_requirements_tenant_range_idx
    ON hrms.staffing_requirements(tenant_id, requirement_date, start_date, end_date)
    WHERE NOT inactive;
CREATE INDEX IF NOT EXISTS staffing_requirements_tenant_scope_idx
    ON hrms.staffing_requirements(tenant_id, branch_id, department_id, attendance_location_id)
    WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.shift_schedule_assignments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    schedule_date DATE NOT NULL,
    worker_profile_id UUID REFERENCES hrms.worker_profiles(id) ON DELETE SET NULL,
    employee_user_id UUID,
    shift_template_id UUID REFERENCES hrms.shift_templates(id) ON DELETE SET NULL,
    attendance_policy_id UUID REFERENCES hrms.attendance_policies(id) ON DELETE SET NULL,
    attendance_location_id UUID REFERENCES hrms.attendance_locations(id) ON DELETE SET NULL,
    attendance_roster_id UUID REFERENCES hrms.attendance_rosters(id) ON DELETE SET NULL,
    branch_id UUID REFERENCES hrms.branches(id) ON DELETE SET NULL,
    department_id UUID REFERENCES hrms.departments(id) ON DELETE SET NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    break_minutes INTEGER NOT NULL DEFAULT 0,
    work_mode TEXT NOT NULL DEFAULT 'office',
    location_type TEXT NOT NULL DEFAULT 'office',
    role_label TEXT,
    team_label TEXT,
    status TEXT NOT NULL DEFAULT 'draft',
    source TEXT NOT NULL DEFAULT 'manual',
    overtime_planned_minutes INTEGER NOT NULL DEFAULT 0,
    has_conflict BOOLEAN NOT NULL DEFAULT FALSE,
    conflict_reason TEXT,
    payroll_blocking BOOLEAN NOT NULL DEFAULT FALSE,
    notes TEXT,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    CONSTRAINT shift_schedule_assignments_time_order_check CHECK (end_time > start_time),
    CONSTRAINT shift_schedule_assignments_worker_check CHECK (worker_profile_id IS NOT NULL OR employee_user_id IS NOT NULL),
    CONSTRAINT shift_schedule_assignments_break_check CHECK (break_minutes >= 0),
    CONSTRAINT shift_schedule_assignments_overtime_check CHECK (overtime_planned_minutes >= 0),
    CONSTRAINT shift_schedule_assignments_work_mode_check CHECK (work_mode IN ('office','remote','field','hybrid')),
    CONSTRAINT shift_schedule_assignments_location_type_check CHECK (location_type IN ('office','remote','field','client_site','hybrid','branch','warehouse','project_site','other')),
    CONSTRAINT shift_schedule_assignments_status_check CHECK (status IN ('draft','published','locked','cancelled','completed')),
    CONSTRAINT shift_schedule_assignments_source_check CHECK (source IN ('manual','template','import','swap','system'))
);

CREATE INDEX IF NOT EXISTS shift_schedule_assignments_tenant_date_idx
    ON hrms.shift_schedule_assignments(tenant_id, schedule_date, status)
    WHERE NOT inactive;
CREATE INDEX IF NOT EXISTS shift_schedule_assignments_worker_date_idx
    ON hrms.shift_schedule_assignments(tenant_id, worker_profile_id, schedule_date)
    WHERE NOT inactive;
CREATE INDEX IF NOT EXISTS shift_schedule_assignments_user_date_idx
    ON hrms.shift_schedule_assignments(tenant_id, employee_user_id, schedule_date)
    WHERE NOT inactive;
CREATE INDEX IF NOT EXISTS shift_schedule_assignments_scope_idx
    ON hrms.shift_schedule_assignments(tenant_id, branch_id, department_id, attendance_location_id)
    WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.shift_swap_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    requester_assignment_id UUID NOT NULL REFERENCES hrms.shift_schedule_assignments(id) ON DELETE CASCADE,
    requester_worker_profile_id UUID REFERENCES hrms.worker_profiles(id) ON DELETE SET NULL,
    requester_user_id UUID,
    target_worker_profile_id UUID REFERENCES hrms.worker_profiles(id) ON DELETE SET NULL,
    target_user_id UUID,
    offered_assignment_id UUID REFERENCES hrms.shift_schedule_assignments(id) ON DELETE SET NULL,
    requested_date DATE,
    requested_shift_template_id UUID REFERENCES hrms.shift_templates(id) ON DELETE SET NULL,
    reason TEXT,
    status TEXT NOT NULL DEFAULT 'pending',
    reviewed_by UUID,
    reviewed_at TIMESTAMPTZ,
    review_remarks TEXT,
    payroll_blocking BOOLEAN NOT NULL DEFAULT TRUE,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    CONSTRAINT shift_swap_requests_target_check CHECK (target_worker_profile_id IS NOT NULL OR target_user_id IS NOT NULL OR offered_assignment_id IS NOT NULL OR requested_shift_template_id IS NOT NULL),
    CONSTRAINT shift_swap_requests_status_check CHECK (status IN ('pending','approved','rejected','cancelled'))
);

CREATE INDEX IF NOT EXISTS shift_swap_requests_tenant_status_idx
    ON hrms.shift_swap_requests(tenant_id, status, created_at DESC)
    WHERE NOT inactive;
CREATE INDEX IF NOT EXISTS shift_swap_requests_requester_idx
    ON hrms.shift_swap_requests(tenant_id, requester_worker_profile_id, requester_user_id)
    WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.shift_schedule_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    source_type TEXT NOT NULL,
    source_id UUID NOT NULL,
    action TEXT NOT NULL,
    from_status TEXT,
    to_status TEXT,
    actor_user_id UUID,
    remarks TEXT,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    CONSTRAINT shift_schedule_events_source_type_check CHECK (source_type IN ('template','requirement','assignment','swap_request'))
);

CREATE INDEX IF NOT EXISTS shift_schedule_events_source_idx
    ON hrms.shift_schedule_events(tenant_id, source_type, source_id, created_at DESC)
    WHERE NOT inactive;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'shift_templates_updated_at') THEN
        CREATE TRIGGER shift_templates_updated_at BEFORE UPDATE ON hrms.shift_templates FOR EACH ROW EXECUTE FUNCTION hrms.set_updated_at();
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'staffing_requirements_updated_at') THEN
        CREATE TRIGGER staffing_requirements_updated_at BEFORE UPDATE ON hrms.staffing_requirements FOR EACH ROW EXECUTE FUNCTION hrms.set_updated_at();
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'shift_schedule_assignments_updated_at') THEN
        CREATE TRIGGER shift_schedule_assignments_updated_at BEFORE UPDATE ON hrms.shift_schedule_assignments FOR EACH ROW EXECUTE FUNCTION hrms.set_updated_at();
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'shift_swap_requests_updated_at') THEN
        CREATE TRIGGER shift_swap_requests_updated_at BEFORE UPDATE ON hrms.shift_swap_requests FOR EACH ROW EXECUTE FUNCTION hrms.set_updated_at();
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'shift_schedule_events_updated_at') THEN
        CREATE TRIGGER shift_schedule_events_updated_at BEFORE UPDATE ON hrms.shift_schedule_events FOR EACH ROW EXECUTE FUNCTION hrms.set_updated_at();
    END IF;
END $$;

ALTER TABLE hrms.shift_templates ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.staffing_requirements ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.shift_schedule_assignments ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.shift_swap_requests ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.shift_schedule_events ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS shift_templates_tenant_isolation ON hrms.shift_templates;
CREATE POLICY shift_templates_tenant_isolation ON hrms.shift_templates
    USING (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true))
    WITH CHECK (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true));

DROP POLICY IF EXISTS staffing_requirements_tenant_isolation ON hrms.staffing_requirements;
CREATE POLICY staffing_requirements_tenant_isolation ON hrms.staffing_requirements
    USING (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true))
    WITH CHECK (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true));

DROP POLICY IF EXISTS shift_schedule_assignments_tenant_isolation ON hrms.shift_schedule_assignments;
CREATE POLICY shift_schedule_assignments_tenant_isolation ON hrms.shift_schedule_assignments
    USING (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true))
    WITH CHECK (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true));

DROP POLICY IF EXISTS shift_swap_requests_tenant_isolation ON hrms.shift_swap_requests;
CREATE POLICY shift_swap_requests_tenant_isolation ON hrms.shift_swap_requests
    USING (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true))
    WITH CHECK (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true));

DROP POLICY IF EXISTS shift_schedule_events_tenant_isolation ON hrms.shift_schedule_events;
CREATE POLICY shift_schedule_events_tenant_isolation ON hrms.shift_schedule_events
    USING (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true))
    WITH CHECK (current_setting('app.is_super_admin', true) = 'true' OR tenant_id::text = current_setting('app.tenant_id', true));
