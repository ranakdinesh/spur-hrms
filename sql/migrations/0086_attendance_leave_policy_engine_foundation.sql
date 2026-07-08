CREATE TABLE IF NOT EXISTS hrms.policy_sets (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    policy_kind     VARCHAR(20) NOT NULL,
    code            VARCHAR(80) NOT NULL,
    name            VARCHAR(160) NOT NULL,
    description     TEXT,
    config          JSONB NOT NULL DEFAULT '{}'::jsonb,
    is_default      BOOLEAN NOT NULL DEFAULT FALSE,
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    effective_from  DATE,
    effective_to    DATE,
    inactive        BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by      UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by      UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT policy_sets_kind_check CHECK (policy_kind IN ('attendance','leave')),
    CONSTRAINT policy_sets_effective_order_check CHECK (effective_to IS NULL OR effective_from IS NULL OR effective_from <= effective_to)
);

CREATE INDEX IF NOT EXISTS policy_sets_tenant_kind_idx ON hrms.policy_sets(tenant_id, policy_kind);
CREATE UNIQUE INDEX IF NOT EXISTS policy_sets_code_idx ON hrms.policy_sets(tenant_id, policy_kind, lower(code)) WHERE NOT inactive;
CREATE UNIQUE INDEX IF NOT EXISTS policy_sets_default_idx ON hrms.policy_sets(tenant_id, policy_kind) WHERE is_default AND is_active AND NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.policy_assignments (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    policy_set_id   UUID NOT NULL REFERENCES hrms.policy_sets(id) ON DELETE CASCADE,
    policy_kind     VARCHAR(20) NOT NULL,
    scope_type      VARCHAR(40) NOT NULL,
    scope_id        UUID,
    role_code       VARCHAR(80),
    priority        INT NOT NULL DEFAULT 0,
    effective_from  DATE,
    effective_to    DATE,
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    inactive        BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by      UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by      UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT policy_assignments_kind_check CHECK (policy_kind IN ('attendance','leave')),
    CONSTRAINT policy_assignments_scope_check CHECK (scope_type IN ('tenant','branch','department','designation','workforce_type','role_group','employee')),
    CONSTRAINT policy_assignments_scope_id_check CHECK (
        (scope_type IN ('tenant','role_group') AND scope_id IS NULL)
        OR (scope_type NOT IN ('tenant','role_group') AND scope_id IS NOT NULL)
    ),
    CONSTRAINT policy_assignments_role_code_check CHECK (
        (scope_type = 'role_group' AND role_code IS NOT NULL AND btrim(role_code) <> '')
        OR (scope_type <> 'role_group' AND role_code IS NULL)
    ),
    CONSTRAINT policy_assignments_effective_order_check CHECK (effective_to IS NULL OR effective_from IS NULL OR effective_from <= effective_to)
);

CREATE INDEX IF NOT EXISTS policy_assignments_tenant_kind_idx ON hrms.policy_assignments(tenant_id, policy_kind);
CREATE INDEX IF NOT EXISTS policy_assignments_policy_idx ON hrms.policy_assignments(policy_set_id);
CREATE INDEX IF NOT EXISTS policy_assignments_scope_idx ON hrms.policy_assignments(tenant_id, policy_kind, scope_type, scope_id);
CREATE UNIQUE INDEX IF NOT EXISTS policy_assignments_unique_scope_idx
    ON hrms.policy_assignments(tenant_id, policy_kind, scope_type, COALESCE(scope_id, '00000000-0000-0000-0000-000000000000'::uuid), COALESCE(lower(role_code), ''), COALESCE(effective_from, '1900-01-01'::date))
    WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.leave_policy_rules (
    id                                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id                           UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    policy_set_id                       UUID NOT NULL REFERENCES hrms.policy_sets(id) ON DELETE CASCADE,
    leave_type_id                       UUID NOT NULL REFERENCES hrms.leave_types(id) ON DELETE CASCADE,
    grant_mode                          VARCHAR(40) NOT NULL DEFAULT 'annual_calendar',
    accrual_frequency                   VARCHAR(40),
    entitlement_days                    NUMERIC(6,2) NOT NULL DEFAULT 0,
    accrual_amount_per_period           NUMERIC(6,2) NOT NULL DEFAULT 0,
    prorate_joiners                     BOOLEAN NOT NULL DEFAULT TRUE,
    probation_handling                  VARCHAR(40) NOT NULL DEFAULT 'eligible',
    rounding_rule                       VARCHAR(40) NOT NULL DEFAULT 'nearest_half',
    max_balance_cap                     NUMERIC(6,2),
    carry_forward_cap                   NUMERIC(6,2),
    encashment_eligible                 BOOLEAN NOT NULL DEFAULT FALSE,
    negative_balance_allowed            BOOLEAN NOT NULL DEFAULT FALSE,
    insufficient_balance_action         VARCHAR(40) NOT NULL DEFAULT 'block',
    expiry_days                         INT,
    allow_half_day                      BOOLEAN NOT NULL DEFAULT TRUE,
    attachment_required_after_days      NUMERIC(6,2),
    approval_workflow                   JSONB NOT NULL DEFAULT '{}'::jsonb,
    sandwich_enabled                    BOOLEAN NOT NULL DEFAULT FALSE,
    sandwich_include_weekly_off         BOOLEAN NOT NULL DEFAULT FALSE,
    sandwich_include_public_holiday     BOOLEAN NOT NULL DEFAULT FALSE,
    sandwich_same_leave_type_only       BOOLEAN NOT NULL DEFAULT TRUE,
    sandwich_across_leave_types         BOOLEAN NOT NULL DEFAULT FALSE,
    notice_required_after_days          NUMERIC(6,2),
    notice_days                         INT NOT NULL DEFAULT 0,
    payroll_impact                      VARCHAR(40) NOT NULL DEFAULT 'none',
    rule_config                         JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive                            BOOLEAN NOT NULL DEFAULT FALSE,
    created_at                          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by                          UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at                          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by                          UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT leave_policy_rules_grant_mode_check CHECK (grant_mode IN ('annual_calendar','annual_financial','anniversary_year','monthly_accrual','payroll_cycle_accrual','days_worked_accrual','attendance_based_accrual','probation_completion','manual_grant','comp_off_earned','no_balance')),
    CONSTRAINT leave_policy_rules_insufficient_action_check CHECK (insufficient_balance_action IN ('block','convert_to_lop','allow_negative','require_override')),
    CONSTRAINT leave_policy_rules_payroll_impact_check CHECK (payroll_impact IN ('none','paid','lop','payroll_adjustment')),
    CONSTRAINT leave_policy_rules_numbers_check CHECK (
        entitlement_days >= 0
        AND accrual_amount_per_period >= 0
        AND (max_balance_cap IS NULL OR max_balance_cap >= 0)
        AND (carry_forward_cap IS NULL OR carry_forward_cap >= 0)
        AND (expiry_days IS NULL OR expiry_days >= 0)
        AND (attachment_required_after_days IS NULL OR attachment_required_after_days >= 0)
        AND (notice_required_after_days IS NULL OR notice_required_after_days >= 0)
        AND notice_days >= 0
    )
);

CREATE INDEX IF NOT EXISTS leave_policy_rules_tenant_policy_idx ON hrms.leave_policy_rules(tenant_id, policy_set_id);
CREATE UNIQUE INDEX IF NOT EXISTS leave_policy_rules_leave_type_idx ON hrms.leave_policy_rules(tenant_id, policy_set_id, leave_type_id) WHERE NOT inactive;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'policy_sets_updated_at') THEN
        CREATE TRIGGER policy_sets_updated_at BEFORE UPDATE ON hrms.policy_sets FOR EACH ROW EXECUTE FUNCTION hrms.update_updated_at();
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'policy_assignments_updated_at') THEN
        CREATE TRIGGER policy_assignments_updated_at BEFORE UPDATE ON hrms.policy_assignments FOR EACH ROW EXECUTE FUNCTION hrms.update_updated_at();
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'leave_policy_rules_updated_at') THEN
        CREATE TRIGGER leave_policy_rules_updated_at BEFORE UPDATE ON hrms.leave_policy_rules FOR EACH ROW EXECUTE FUNCTION hrms.update_updated_at();
    END IF;
END $$;

ALTER TABLE hrms.policy_sets ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.policy_assignments ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.leave_policy_rules ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS policy_sets_tenant_isolation ON hrms.policy_sets;
CREATE POLICY policy_sets_tenant_isolation ON hrms.policy_sets
USING (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true')
WITH CHECK (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true');

DROP POLICY IF EXISTS policy_assignments_tenant_isolation ON hrms.policy_assignments;
CREATE POLICY policy_assignments_tenant_isolation ON hrms.policy_assignments
USING (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true')
WITH CHECK (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true');

DROP POLICY IF EXISTS leave_policy_rules_tenant_isolation ON hrms.leave_policy_rules;
CREATE POLICY leave_policy_rules_tenant_isolation ON hrms.leave_policy_rules
USING (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true')
WITH CHECK (tenant_id::text = current_setting('app.tenant_id', true) OR current_setting('app.is_super_admin', true) = 'true');

INSERT INTO hrms.leave_types (
    tenant_id, name, shortcode, description, is_paid, is_carry_forward,
    max_carry_forward, is_consecutive_limit, consecutive_days_limit,
    is_enabled, is_system
)
SELECT t.id, defaults.name, defaults.shortcode, defaults.description, defaults.is_paid,
       defaults.is_carry_forward, defaults.max_carry_forward, false, 0, true, true
FROM auth.tenants t
CROSS JOIN (
    VALUES
        ('Sick Leave', 'SL', 'Default India sick leave bucket.', true, false, 0),
        ('Casual Leave', 'CL', 'Default India casual leave bucket.', true, false, 0),
        ('Earned Leave', 'EL', 'Default India earned/paid/privilege leave bucket.', true, true, 30),
        ('Comp Off', 'CO', 'Compensatory off credited from approved holiday/weekoff/extra work.', true, false, 0),
        ('Loss of Pay', 'LOP', 'Unpaid leave bucket used when balance is unavailable or policy requires LOP.', false, false, 0)
) AS defaults(name, shortcode, description, is_paid, is_carry_forward, max_carry_forward)
WHERE t.kind <> 'ops'
ON CONFLICT (tenant_id, lower(shortcode)) WHERE shortcode IS NOT NULL AND NOT inactive
DO UPDATE SET
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    is_system = TRUE,
    is_enabled = hrms.leave_types.is_enabled,
    updated_at = NOW();

INSERT INTO hrms.policy_sets (
    tenant_id, policy_kind, code, name, description, config,
    is_default, is_active, effective_from
)
SELECT
    t.id,
    'attendance',
    'default_attendance',
    'Default Attendance Policy',
    'Default tenant attendance policy for fixed office or hybrid work.',
    jsonb_build_object(
        'schedule_type', 'fixed_shift',
        'standard_work_minutes', 480,
        'min_half_day_minutes', 240,
        'min_full_day_minutes', 420,
        'grace_late_minutes', 10,
        'grace_early_minutes', 10,
        'allow_flexible_hours', false,
        'strict_in_out', true,
        'location_required', false,
        'geofence_required', false,
        'missed_punch_regularization_allowed', true,
        'regularization_deadline_days', 3,
        'allowed_work_modes', jsonb_build_array('office', 'remote', 'field', 'hybrid', 'client_site', 'project_site'),
        'allowed_punch_sources', jsonb_build_array('web', 'mobile', 'biometric', 'kiosk', 'integration'),
        'holiday_weekoff_work', jsonb_build_object(
            'default_benefit', 'approval_required',
            'allow_comp_off', true,
            'allow_overtime_payment', false,
            'minimum_minutes', 240,
            'comp_off_expiry_days', 90,
            'payroll_impact', false
        ),
        'overtime', jsonb_build_object(
            'enabled', false,
            'daily_threshold_minutes', 540,
            'pre_approval_required', true,
            'calculation_type', 'none',
            'payroll_impact', false
        )
    ),
    true,
    true,
    CURRENT_DATE
FROM auth.tenants t
WHERE t.kind <> 'ops'
  AND NOT EXISTS (
      SELECT 1 FROM hrms.policy_sets ps
      WHERE ps.tenant_id = t.id
        AND ps.policy_kind = 'attendance'
        AND ps.is_default
        AND ps.is_active
        AND NOT ps.inactive
  );

INSERT INTO hrms.policy_sets (
    tenant_id, policy_kind, code, name, description, config,
    is_default, is_active, effective_from
)
SELECT
    t.id,
    'leave',
    'default_leave',
    'Default India Leave Policy',
    'Default India leave policy with Sick, Casual, Earned, Comp Off, and LOP buckets.',
    jsonb_build_object(
        'region', 'IN',
        'leave_year', 'calendar',
        'approval_workflow', jsonb_build_object('mode', 'manager_then_hr_optional'),
        'sandwich_default', jsonb_build_object(
            'enabled', false,
            'include_weekly_off', false,
            'include_public_holiday', false
        ),
        'long_leave_notice_default', jsonb_build_object(
            'duration_days_greater_than', 10,
            'notice_days', 7,
            'action', 'warn'
        ),
        'staffing_threshold_default', jsonb_build_object(
            'enabled', false,
            'action', 'warn'
        )
    ),
    true,
    true,
    CURRENT_DATE
FROM auth.tenants t
WHERE t.kind <> 'ops'
  AND NOT EXISTS (
      SELECT 1 FROM hrms.policy_sets ps
      WHERE ps.tenant_id = t.id
        AND ps.policy_kind = 'leave'
        AND ps.is_default
        AND ps.is_active
        AND NOT ps.inactive
  );

INSERT INTO hrms.policy_assignments (
    tenant_id, policy_set_id, policy_kind, scope_type, priority,
    effective_from, is_active
)
SELECT ps.tenant_id, ps.id, ps.policy_kind, 'tenant', 0, CURRENT_DATE, true
FROM hrms.policy_sets ps
WHERE ps.policy_kind IN ('attendance', 'leave')
  AND ps.is_default
  AND ps.is_active
  AND NOT ps.inactive
  AND NOT EXISTS (
      SELECT 1 FROM hrms.policy_assignments pa
      WHERE pa.tenant_id = ps.tenant_id
        AND pa.policy_kind = ps.policy_kind
        AND pa.scope_type = 'tenant'
        AND pa.is_active
        AND NOT pa.inactive
  );

INSERT INTO hrms.leave_policy_rules (
    tenant_id, policy_set_id, leave_type_id, grant_mode, accrual_frequency,
    entitlement_days, accrual_amount_per_period, prorate_joiners, probation_handling,
    rounding_rule, max_balance_cap, carry_forward_cap, encashment_eligible,
    negative_balance_allowed, insufficient_balance_action, expiry_days, allow_half_day,
    attachment_required_after_days, approval_workflow, sandwich_enabled,
    sandwich_include_weekly_off, sandwich_include_public_holiday,
    sandwich_same_leave_type_only, sandwich_across_leave_types,
    notice_required_after_days, notice_days, payroll_impact, rule_config
)
SELECT
    ps.tenant_id,
    ps.id,
    lt.id,
    CASE lt.shortcode
        WHEN 'EL' THEN 'monthly_accrual'
        WHEN 'CO' THEN 'comp_off_earned'
        WHEN 'LOP' THEN 'no_balance'
        ELSE 'annual_calendar'
    END,
    CASE lt.shortcode
        WHEN 'EL' THEN 'monthly'
        ELSE NULL
    END,
    CASE lt.shortcode
        WHEN 'SL' THEN 7
        WHEN 'CL' THEN 7
        WHEN 'EL' THEN 0
        ELSE 0
    END,
    CASE lt.shortcode
        WHEN 'EL' THEN 1.5
        ELSE 0
    END,
    true,
    CASE lt.shortcode
        WHEN 'LOP' THEN 'always_eligible'
        ELSE 'eligible'
    END,
    'nearest_half',
    CASE lt.shortcode
        WHEN 'EL' THEN 30
        ELSE NULL
    END,
    CASE lt.shortcode
        WHEN 'EL' THEN 30
        ELSE 0
    END,
    CASE lt.shortcode
        WHEN 'EL' THEN true
        ELSE false
    END,
    CASE lt.shortcode
        WHEN 'LOP' THEN true
        ELSE false
    END,
    CASE lt.shortcode
        WHEN 'LOP' THEN 'allow_negative'
        ELSE 'block'
    END,
    CASE lt.shortcode
        WHEN 'CO' THEN 90
        ELSE NULL
    END,
    true,
    CASE lt.shortcode
        WHEN 'SL' THEN 2
        ELSE NULL
    END,
    '{}'::jsonb,
    false,
    false,
    false,
    true,
    false,
    CASE lt.shortcode
        WHEN 'EL' THEN 3
        ELSE NULL
    END,
    CASE lt.shortcode
        WHEN 'EL' THEN 5
        ELSE 0
    END,
    CASE lt.shortcode
        WHEN 'LOP' THEN 'lop'
        ELSE 'paid'
    END,
    jsonb_build_object('default_india_seed', true, 'leave_shortcode', lt.shortcode)
FROM hrms.policy_sets ps
JOIN hrms.leave_types lt ON lt.tenant_id = ps.tenant_id
WHERE ps.policy_kind = 'leave'
  AND ps.is_default
  AND ps.is_active
  AND NOT ps.inactive
  AND lt.shortcode IN ('SL', 'CL', 'EL', 'CO', 'LOP')
  AND NOT lt.inactive
ON CONFLICT (tenant_id, policy_set_id, leave_type_id) WHERE NOT inactive
DO UPDATE SET
    grant_mode = EXCLUDED.grant_mode,
    accrual_frequency = EXCLUDED.accrual_frequency,
    entitlement_days = EXCLUDED.entitlement_days,
    accrual_amount_per_period = EXCLUDED.accrual_amount_per_period,
    prorate_joiners = EXCLUDED.prorate_joiners,
    probation_handling = EXCLUDED.probation_handling,
    rounding_rule = EXCLUDED.rounding_rule,
    max_balance_cap = EXCLUDED.max_balance_cap,
    carry_forward_cap = EXCLUDED.carry_forward_cap,
    encashment_eligible = EXCLUDED.encashment_eligible,
    negative_balance_allowed = EXCLUDED.negative_balance_allowed,
    insufficient_balance_action = EXCLUDED.insufficient_balance_action,
    expiry_days = EXCLUDED.expiry_days,
    allow_half_day = EXCLUDED.allow_half_day,
    attachment_required_after_days = EXCLUDED.attachment_required_after_days,
    approval_workflow = EXCLUDED.approval_workflow,
    sandwich_enabled = EXCLUDED.sandwich_enabled,
    sandwich_include_weekly_off = EXCLUDED.sandwich_include_weekly_off,
    sandwich_include_public_holiday = EXCLUDED.sandwich_include_public_holiday,
    sandwich_same_leave_type_only = EXCLUDED.sandwich_same_leave_type_only,
    sandwich_across_leave_types = EXCLUDED.sandwich_across_leave_types,
    notice_required_after_days = EXCLUDED.notice_required_after_days,
    notice_days = EXCLUDED.notice_days,
    payroll_impact = EXCLUDED.payroll_impact,
    rule_config = EXCLUDED.rule_config,
    updated_at = NOW();
