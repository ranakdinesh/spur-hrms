-- HRMS core schema baseline.
-- Identity owns auth.tenants/auth.users. HRMS references those tables only.

CREATE SCHEMA IF NOT EXISTS hrms;

-- ============================================================
-- TENANT / COMPANY SETUP
-- ============================================================

CREATE TABLE IF NOT EXISTS hrms.departments (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id      UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    name           VARCHAR(255) NOT NULL,
    description    TEXT,
    inactive       BOOLEAN NOT NULL DEFAULT FALSE,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by     UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by     UUID REFERENCES auth.users(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS departments_tenant_idx ON hrms.departments(tenant_id);
CREATE UNIQUE INDEX IF NOT EXISTS departments_tenant_name_idx ON hrms.departments(tenant_id, lower(name)) WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.branches (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id      UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    branch_name    VARCHAR(255) NOT NULL,
    address        TEXT,
    city           VARCHAR(100),
    state          VARCHAR(100),
    country        VARCHAR(100),
    pincode        VARCHAR(20),
    phone          VARCHAR(30),
    inactive       BOOLEAN NOT NULL DEFAULT FALSE,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by     UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by     UUID REFERENCES auth.users(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS branches_tenant_idx ON hrms.branches(tenant_id);
CREATE UNIQUE INDEX IF NOT EXISTS branches_tenant_name_idx ON hrms.branches(tenant_id, lower(branch_name)) WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.designations (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id      UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    name           VARCHAR(255) NOT NULL,
    description    TEXT,
    inactive       BOOLEAN NOT NULL DEFAULT FALSE,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by     UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by     UUID REFERENCES auth.users(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS designations_tenant_idx ON hrms.designations(tenant_id);
CREATE UNIQUE INDEX IF NOT EXISTS designations_tenant_name_idx ON hrms.designations(tenant_id, lower(name)) WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.employment_types (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    name        VARCHAR(100) NOT NULL,
    inactive    BOOLEAN NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by  UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by  UUID REFERENCES auth.users(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS employment_types_tenant_idx ON hrms.employment_types(tenant_id);
CREATE UNIQUE INDEX IF NOT EXISTS employment_types_tenant_name_idx ON hrms.employment_types(tenant_id, lower(name)) WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.marital_statuses (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    name        VARCHAR(100) NOT NULL,
    inactive    BOOLEAN NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by  UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by  UUID REFERENCES auth.users(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS marital_statuses_tenant_idx ON hrms.marital_statuses(tenant_id);
CREATE UNIQUE INDEX IF NOT EXISTS marital_statuses_tenant_name_idx ON hrms.marital_statuses(tenant_id, lower(name)) WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.financial_years (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id     UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    name          VARCHAR(100),
    start_date    DATE NOT NULL,
    end_date      DATE NOT NULL,
    is_active     BOOLEAN NOT NULL DEFAULT FALSE,
    inactive      BOOLEAN NOT NULL DEFAULT FALSE,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by    UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by    UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT financial_years_date_order CHECK (start_date <= end_date)
);
CREATE INDEX IF NOT EXISTS financial_years_tenant_idx ON hrms.financial_years(tenant_id);
CREATE UNIQUE INDEX IF NOT EXISTS financial_years_active_idx ON hrms.financial_years(tenant_id) WHERE is_active AND NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.working_hours (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    branch_id       UUID REFERENCES hrms.branches(id) ON DELETE CASCADE,
    user_id         UUID REFERENCES auth.users(id) ON DELETE CASCADE,
    day_of_week     VARCHAR(20) NOT NULL,
    is_working_day  BOOLEAN NOT NULL DEFAULT TRUE,
    start_time      TIME NOT NULL,
    end_time        TIME NOT NULL,
    break_minutes   INT NOT NULL DEFAULT 0,
    inactive        BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by      UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by      UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT working_hours_day_check CHECK (day_of_week IN ('Monday','Tuesday','Wednesday','Thursday','Friday','Saturday','Sunday')),
    CONSTRAINT working_hours_break_check CHECK (break_minutes >= 0)
);
CREATE INDEX IF NOT EXISTS working_hours_tenant_idx ON hrms.working_hours(tenant_id);
CREATE INDEX IF NOT EXISTS working_hours_scope_idx ON hrms.working_hours(tenant_id, branch_id, user_id) WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.holidays (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    branch_id    UUID REFERENCES hrms.branches(id) ON DELETE CASCADE,
    fy_id        UUID REFERENCES hrms.financial_years(id) ON DELETE SET NULL,
    name         VARCHAR(255) NOT NULL,
    date         DATE NOT NULL,
    is_optional  BOOLEAN NOT NULL DEFAULT FALSE,
    inactive     BOOLEAN NOT NULL DEFAULT FALSE,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by   UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by   UUID REFERENCES auth.users(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS holidays_tenant_idx ON hrms.holidays(tenant_id);
CREATE INDEX IF NOT EXISTS holidays_date_idx ON hrms.holidays(date);

CREATE TABLE IF NOT EXISTS hrms.pay_cycles (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id         UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    cycle_type        VARCHAR(50) NOT NULL DEFAULT 'monthly',
    pay_day           INT,
    start_day         INT,
    end_day           INT,
    inactive          BOOLEAN NOT NULL DEFAULT FALSE,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by        UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by        UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT pay_cycles_day_check CHECK (
        (pay_day IS NULL OR pay_day BETWEEN 1 AND 31)
        AND (start_day IS NULL OR start_day BETWEEN 1 AND 31)
        AND (end_day IS NULL OR end_day BETWEEN 1 AND 31)
    )
);
CREATE UNIQUE INDEX IF NOT EXISTS pay_cycles_tenant_active_idx ON hrms.pay_cycles(tenant_id) WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.policy_types (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    name        VARCHAR(255) NOT NULL,
    inactive    BOOLEAN NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by  UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by  UUID REFERENCES auth.users(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS policy_types_tenant_idx ON hrms.policy_types(tenant_id);
CREATE UNIQUE INDEX IF NOT EXISTS policy_types_tenant_name_idx ON hrms.policy_types(tenant_id, lower(name)) WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.company_policies (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    policy_type_id  UUID REFERENCES hrms.policy_types(id) ON DELETE SET NULL,
    title           VARCHAR(255) NOT NULL,
    file_path       TEXT,
    description     TEXT,
    inactive        BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by      UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by      UUID REFERENCES auth.users(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS company_policies_tenant_idx ON hrms.company_policies(tenant_id);

CREATE TABLE IF NOT EXISTS hrms.tenant_subscriptions (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id        UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    plan_id          UUID,
    start_date       DATE,
    end_date         DATE,
    status           VARCHAR(50),
    max_employees    INT NOT NULL DEFAULT 0,
    inactive         BOOLEAN NOT NULL DEFAULT FALSE,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by       UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by       UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT tenant_subscriptions_employee_limit_check CHECK (max_employees >= 0)
);
CREATE INDEX IF NOT EXISTS tenant_subscriptions_tenant_idx ON hrms.tenant_subscriptions(tenant_id);

-- ============================================================
-- EMPLOYEE
-- ============================================================

CREATE TABLE IF NOT EXISTS hrms.employees (
    id                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id            UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    user_id              UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    employee_code        VARCHAR(50),
    firstname            VARCHAR(100) NOT NULL,
    lastname             VARCHAR(100),
    email                VARCHAR(255),
    mobile               VARCHAR(20),
    dob                  DATE,
    gender               VARCHAR(20),
    marital_status       VARCHAR(30),
    blood_group          VARCHAR(10),
    profile_photo_path   TEXT,
    address              TEXT,
    city                 VARCHAR(100),
    state                VARCHAR(100),
    country              VARCHAR(100),
    pincode              VARCHAR(20),
    emergency_contact    VARCHAR(255),
    joining_date         DATE,
    resignation_date     DATE,
    department_id        UUID REFERENCES hrms.departments(id) ON DELETE SET NULL,
    branch_id            UUID REFERENCES hrms.branches(id) ON DELETE SET NULL,
    designation_id       UUID REFERENCES hrms.designations(id) ON DELETE SET NULL,
    reporting_manager_id UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    employment_type_id   UUID REFERENCES hrms.employment_types(id) ON DELETE SET NULL,
    role                 VARCHAR(50),
    grade                VARCHAR(50),
    experience_year      INT NOT NULL DEFAULT 0,
    experience_month     INT NOT NULL DEFAULT 0,
    inactive             BOOLEAN NOT NULL DEFAULT FALSE,
    created_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by           UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by           UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT employees_experience_check CHECK (experience_year >= 0 AND experience_month BETWEEN 0 AND 11),
    CONSTRAINT employees_resignation_after_joining CHECK (resignation_date IS NULL OR joining_date IS NULL OR resignation_date >= joining_date)
);
CREATE INDEX IF NOT EXISTS employees_tenant_idx ON hrms.employees(tenant_id);
CREATE INDEX IF NOT EXISTS employees_user_idx ON hrms.employees(user_id);
CREATE UNIQUE INDEX IF NOT EXISTS employees_code_tenant_idx ON hrms.employees(tenant_id, employee_code) WHERE employee_code IS NOT NULL AND NOT inactive;
CREATE UNIQUE INDEX IF NOT EXISTS employees_user_tenant_idx ON hrms.employees(tenant_id, user_id) WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.employee_statutory (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    user_id         UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    pf_no           VARCHAR(100),
    uan_no          VARCHAR(100),
    esic_no         VARCHAR(100),
    pan             VARCHAR(20),
    aadhaar         VARCHAR(20),
    pt_applicable   BOOLEAN NOT NULL DEFAULT FALSE,
    pf_applicable   BOOLEAN NOT NULL DEFAULT FALSE,
    esic_applicable BOOLEAN NOT NULL DEFAULT FALSE,
    inactive        BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by      UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by      UUID REFERENCES auth.users(id) ON DELETE SET NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS employee_statutory_user_tenant_idx ON hrms.employee_statutory(tenant_id, user_id) WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.employee_banks (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    user_id         UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    bank_name       VARCHAR(255),
    account_number  VARCHAR(50),
    ifsc_code       VARCHAR(20),
    account_type    VARCHAR(50),
    branch_name     VARCHAR(255),
    is_primary      BOOLEAN NOT NULL DEFAULT TRUE,
    inactive        BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by      UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by      UUID REFERENCES auth.users(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS employee_banks_tenant_user_idx ON hrms.employee_banks(tenant_id, user_id);
CREATE UNIQUE INDEX IF NOT EXISTS employee_banks_primary_idx ON hrms.employee_banks(tenant_id, user_id) WHERE is_primary AND NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.document_types (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    name        VARCHAR(255) NOT NULL,
    inactive    BOOLEAN NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by  UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by  UUID REFERENCES auth.users(id) ON DELETE SET NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS document_types_tenant_name_idx ON hrms.document_types(tenant_id, lower(name)) WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.employee_documents (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id        UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    user_id          UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    document_type_id UUID REFERENCES hrms.document_types(id) ON DELETE SET NULL,
    title            VARCHAR(255),
    file_path        TEXT,
    inactive         BOOLEAN NOT NULL DEFAULT FALSE,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by       UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by       UUID REFERENCES auth.users(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS employee_documents_tenant_user_idx ON hrms.employee_documents(tenant_id, user_id);

-- ============================================================
-- LEAVE
-- ============================================================

CREATE TABLE IF NOT EXISTS hrms.leave_types (
    id                     UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id              UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    name                   VARCHAR(100) NOT NULL,
    shortcode              VARCHAR(20),
    description            TEXT,
    is_paid                BOOLEAN NOT NULL DEFAULT TRUE,
    is_carry_forward       BOOLEAN NOT NULL DEFAULT FALSE,
    max_carry_forward      INT NOT NULL DEFAULT 0,
    is_consecutive_limit   BOOLEAN NOT NULL DEFAULT FALSE,
    consecutive_days_limit INT NOT NULL DEFAULT 0,
    is_system              BOOLEAN NOT NULL DEFAULT FALSE,
    inactive               BOOLEAN NOT NULL DEFAULT FALSE,
    created_at             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by             UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by             UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT leave_types_limits_check CHECK (max_carry_forward >= 0 AND consecutive_days_limit >= 0)
);
CREATE INDEX IF NOT EXISTS leave_types_tenant_idx ON hrms.leave_types(tenant_id);
CREATE UNIQUE INDEX IF NOT EXISTS leave_types_shortcode_idx ON hrms.leave_types(tenant_id, lower(shortcode)) WHERE shortcode IS NOT NULL AND NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.leave_policies (
    id                       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id                UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    leave_type_id            UUID NOT NULL REFERENCES hrms.leave_types(id) ON DELETE CASCADE,
    fy_id                    UUID NOT NULL REFERENCES hrms.financial_years(id) ON DELETE CASCADE,
    total_days               NUMERIC(5,1) NOT NULL DEFAULT 0,
    allocation_type          VARCHAR(20) NOT NULL DEFAULT 'fixed',
    jan INT NOT NULL DEFAULT 0, feb INT NOT NULL DEFAULT 0, mar INT NOT NULL DEFAULT 0,
    apr INT NOT NULL DEFAULT 0, may INT NOT NULL DEFAULT 0, jun INT NOT NULL DEFAULT 0,
    jul INT NOT NULL DEFAULT 0, aug INT NOT NULL DEFAULT 0, sep INT NOT NULL DEFAULT 0,
    oct INT NOT NULL DEFAULT 0, nov INT NOT NULL DEFAULT 0, dec INT NOT NULL DEFAULT 0,
    is_sandwich_applicable  BOOLEAN NOT NULL DEFAULT FALSE,
    inactive                 BOOLEAN NOT NULL DEFAULT FALSE,
    created_at               TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by               UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at               TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by               UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT leave_policies_allocation_check CHECK (allocation_type IN ('fixed','monthly')),
    CONSTRAINT leave_policies_days_check CHECK (total_days >= 0)
);
CREATE UNIQUE INDEX IF NOT EXISTS leave_policy_unique_idx ON hrms.leave_policies(tenant_id, leave_type_id, fy_id) WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.leave_balances (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    user_id         UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    leave_type_id   UUID NOT NULL REFERENCES hrms.leave_types(id) ON DELETE CASCADE,
    fy_id           UUID NOT NULL REFERENCES hrms.financial_years(id) ON DELETE CASCADE,
    total_days      NUMERIC(5,1) NOT NULL DEFAULT 0,
    used_days       NUMERIC(5,1) NOT NULL DEFAULT 0,
    pending_days    NUMERIC(5,1) NOT NULL DEFAULT 0,
    balance_days    NUMERIC(5,1) GENERATED ALWAYS AS (total_days - used_days - pending_days) STORED,
    inactive        BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by      UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by      UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT leave_balances_days_check CHECK (total_days >= 0 AND used_days >= 0 AND pending_days >= 0)
);
CREATE UNIQUE INDEX IF NOT EXISTS leave_balance_unique_idx ON hrms.leave_balances(tenant_id, user_id, leave_type_id, fy_id) WHERE NOT inactive;
CREATE INDEX IF NOT EXISTS leave_balances_user_idx ON hrms.leave_balances(tenant_id, user_id);

CREATE TABLE IF NOT EXISTS hrms.leaves (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id           UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    user_id             UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    leave_type_id       UUID NOT NULL REFERENCES hrms.leave_types(id) ON DELETE RESTRICT,
    fy_id               UUID NOT NULL REFERENCES hrms.financial_years(id) ON DELETE RESTRICT,
    start_date          DATE NOT NULL,
    end_date            DATE NOT NULL,
    start_day_type      VARCHAR(20) NOT NULL DEFAULT 'fullday',
    end_day_type        VARCHAR(20) NOT NULL DEFAULT 'fullday',
    days                NUMERIC(5,1),
    reason              TEXT,
    status              VARCHAR(20) NOT NULL DEFAULT 'pending',
    applied_date        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    from_leave_type     UUID REFERENCES hrms.leave_types(id) ON DELETE SET NULL,
    to_leave_type       UUID REFERENCES hrms.leave_types(id) ON DELETE SET NULL,
    is_sandwich         BOOLEAN NOT NULL DEFAULT FALSE,
    inactive            BOOLEAN NOT NULL DEFAULT FALSE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by          UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by          UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT leaves_date_order CHECK (start_date <= end_date),
    CONSTRAINT leaves_day_type_check CHECK (start_day_type IN ('fullday','firsthalf','secondhalf') AND end_day_type IN ('fullday','firsthalf','secondhalf')),
    CONSTRAINT leaves_status_check CHECK (status IN ('pending','approved','rejected','canceled')),
    CONSTRAINT leaves_days_check CHECK (days IS NULL OR days >= 0)
);
CREATE INDEX IF NOT EXISTS leaves_tenant_user_idx ON hrms.leaves(tenant_id, user_id);
CREATE INDEX IF NOT EXISTS leaves_dates_idx ON hrms.leaves(start_date, end_date);
CREATE INDEX IF NOT EXISTS leaves_status_idx ON hrms.leaves(status);

CREATE TABLE IF NOT EXISTS hrms.leave_approvals (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    leave_id        UUID NOT NULL REFERENCES hrms.leaves(id) ON DELETE CASCADE,
    approver_id     UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    status          VARCHAR(20) NOT NULL DEFAULT 'pending',
    remarks         TEXT,
    action_date     TIMESTAMPTZ,
    inactive        BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by      UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by      UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT leave_approvals_status_check CHECK (status IN ('pending','approved','rejected','canceled'))
);
CREATE INDEX IF NOT EXISTS leave_approvals_leave_idx ON hrms.leave_approvals(leave_id);
CREATE INDEX IF NOT EXISTS leave_approvals_approver_idx ON hrms.leave_approvals(approver_id);

CREATE TABLE IF NOT EXISTS hrms.leave_ledger (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id        UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    user_id          UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    leave_type_id    UUID NOT NULL REFERENCES hrms.leave_types(id) ON DELETE RESTRICT,
    fy_id            UUID NOT NULL REFERENCES hrms.financial_years(id) ON DELETE RESTRICT,
    leave_id         UUID REFERENCES hrms.leaves(id) ON DELETE SET NULL,
    transaction_type VARCHAR(10) NOT NULL,
    days             NUMERIC(5,1) NOT NULL,
    remarks          TEXT,
    inactive         BOOLEAN NOT NULL DEFAULT FALSE,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by       UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT leave_ledger_transaction_type_check CHECK (transaction_type IN ('debit','credit')),
    CONSTRAINT leave_ledger_days_check CHECK (days >= 0)
);
CREATE INDEX IF NOT EXISTS leave_ledger_user_idx ON hrms.leave_ledger(tenant_id, user_id);

-- ============================================================
-- ATTENDANCE
-- ============================================================

CREATE TABLE IF NOT EXISTS hrms.attendances (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    user_id      UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    date         DATE NOT NULL,
    time         TIMESTAMPTZ,
    type         VARCHAR(20),
    status       VARCHAR(30),
    source       VARCHAR(50),
    latitude     NUMERIC(10,7),
    longitude    NUMERIC(10,7),
    work_mode    VARCHAR(50),
    remarks      TEXT,
    inactive     BOOLEAN NOT NULL DEFAULT FALSE,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by   UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by   UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT attendances_type_check CHECK (type IS NULL OR type IN ('checkin','checkout')),
    CONSTRAINT attendances_status_check CHECK (status IS NULL OR status IN ('present','leave','absent','holiday','weekoff'))
);
CREATE INDEX IF NOT EXISTS attendances_tenant_user_date_idx ON hrms.attendances(tenant_id, user_id, date);

CREATE TABLE IF NOT EXISTS hrms.attendance_requests (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    user_id         UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    date            DATE NOT NULL,
    requested_type  VARCHAR(20),
    reason          TEXT,
    status          VARCHAR(20) NOT NULL DEFAULT 'pending',
    reviewed_by     UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    reviewed_at     TIMESTAMPTZ,
    remarks         TEXT,
    inactive        BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by      UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by      UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT attendance_requests_type_check CHECK (requested_type IS NULL OR requested_type IN ('halfday','absent')),
    CONSTRAINT attendance_requests_status_check CHECK (status IN ('pending','approved','rejected','canceled'))
);
CREATE INDEX IF NOT EXISTS attendance_requests_user_idx ON hrms.attendance_requests(tenant_id, user_id);

CREATE TABLE IF NOT EXISTS hrms.device_logs (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    user_id     UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    device_id   VARCHAR(255),
    device_type VARCHAR(50),
    ip_address  VARCHAR(45),
    action      VARCHAR(50),
    inactive    BOOLEAN NOT NULL DEFAULT FALSE,
    logged_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by  UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by  UUID REFERENCES auth.users(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS device_logs_tenant_user_idx ON hrms.device_logs(tenant_id, user_id);

-- ============================================================
-- SALARY
-- ============================================================

CREATE TABLE IF NOT EXISTS hrms.salary_templates (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    fy_id       UUID NOT NULL REFERENCES hrms.financial_years(id) ON DELETE CASCADE,
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    is_active   BOOLEAN NOT NULL DEFAULT FALSE,
    inactive    BOOLEAN NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by  UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by  UUID REFERENCES auth.users(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS salary_templates_tenant_fy_idx ON hrms.salary_templates(tenant_id, fy_id);
CREATE UNIQUE INDEX IF NOT EXISTS salary_templates_active_idx ON hrms.salary_templates(tenant_id, fy_id) WHERE is_active AND NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.salary_template_items (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    template_id     UUID NOT NULL REFERENCES hrms.salary_templates(id) ON DELETE CASCADE,
    item_type       VARCHAR(20) NOT NULL,
    code            VARCHAR(50) NOT NULL,
    name            VARCHAR(255) NOT NULL,
    percentage      NUMERIC(7,2),
    amount          NUMERIC(12,2),
    is_tax_exempt   BOOLEAN NOT NULL DEFAULT FALSE,
    sort_order      INT NOT NULL DEFAULT 0,
    inactive        BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by      UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by      UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT salary_template_items_type_check CHECK (item_type IN ('earning','deduction')),
    CONSTRAINT salary_template_items_amount_source_check CHECK (percentage IS NOT NULL OR amount IS NOT NULL)
);
CREATE INDEX IF NOT EXISTS salary_template_items_template_idx ON hrms.salary_template_items(template_id);

CREATE TABLE IF NOT EXISTS hrms.employee_salaries (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    user_id         UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    fy_id           UUID NOT NULL REFERENCES hrms.financial_years(id) ON DELETE CASCADE,
    template_id     UUID NOT NULL REFERENCES hrms.salary_templates(id) ON DELETE RESTRICT,
    gross_salary    NUMERIC(12,2) NOT NULL,
    effective_from  DATE,
    inactive        BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by      UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by      UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT employee_salaries_gross_check CHECK (gross_salary >= 0)
);
CREATE INDEX IF NOT EXISTS employee_salaries_tenant_user_idx ON hrms.employee_salaries(tenant_id, user_id);

CREATE TABLE IF NOT EXISTS hrms.employee_salary_structures (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    user_id         UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    template_id     UUID NOT NULL REFERENCES hrms.salary_templates(id) ON DELETE RESTRICT,
    fy_id           UUID NOT NULL REFERENCES hrms.financial_years(id) ON DELETE CASCADE,
    item_type       VARCHAR(20) NOT NULL,
    code            VARCHAR(50) NOT NULL,
    name            VARCHAR(255) NOT NULL,
    amount          NUMERIC(12,2) NOT NULL DEFAULT 0,
    sort_order      INT NOT NULL DEFAULT 0,
    inactive        BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by      UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by      UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT employee_salary_structures_type_check CHECK (item_type IN ('earning','deduction')),
    CONSTRAINT employee_salary_structures_amount_check CHECK (amount >= 0)
);
CREATE INDEX IF NOT EXISTS employee_salary_structures_user_idx ON hrms.employee_salary_structures(tenant_id, user_id, fy_id);

CREATE TABLE IF NOT EXISTS hrms.salary_slips (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id           UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    user_id             UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    fy_id               UUID NOT NULL REFERENCES hrms.financial_years(id) ON DELETE RESTRICT,
    template_id         UUID NOT NULL REFERENCES hrms.salary_templates(id) ON DELETE RESTRICT,
    month               INT NOT NULL,
    year                INT NOT NULL,
    gross_salary        NUMERIC(12,2),
    total_earnings      NUMERIC(12,2),
    total_deductions    NUMERIC(12,2),
    absent_deduction    NUMERIC(12,2) NOT NULL DEFAULT 0,
    net_salary          NUMERIC(12,2),
    absent_days         INT NOT NULL DEFAULT 0,
    present_days        INT NOT NULL DEFAULT 0,
    total_days          INT NOT NULL DEFAULT 0,
    lwp_days            NUMERIC(5,1) NOT NULL DEFAULT 0,
    no_of_ph_weo        INT NOT NULL DEFAULT 0,
    is_special          BOOLEAN NOT NULL DEFAULT FALSE,
    is_regenerated      BOOLEAN NOT NULL DEFAULT FALSE,
    pdf_path            TEXT,
    inactive            BOOLEAN NOT NULL DEFAULT FALSE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by          UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by          UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT salary_slips_month_check CHECK (month BETWEEN 1 AND 12),
    CONSTRAINT salary_slips_day_counts_check CHECK (absent_days >= 0 AND present_days >= 0 AND total_days >= 0 AND lwp_days >= 0 AND no_of_ph_weo >= 0)
);
CREATE UNIQUE INDEX IF NOT EXISTS salary_slips_unique_idx ON hrms.salary_slips(tenant_id, user_id, month, year) WHERE NOT inactive;
CREATE INDEX IF NOT EXISTS salary_slips_tenant_user_idx ON hrms.salary_slips(tenant_id, user_id);

CREATE TABLE IF NOT EXISTS hrms.salary_slip_items (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    slip_id     UUID NOT NULL REFERENCES hrms.salary_slips(id) ON DELETE CASCADE,
    item_type   VARCHAR(20) NOT NULL,
    code        VARCHAR(50) NOT NULL,
    name        VARCHAR(255) NOT NULL,
    amount      NUMERIC(12,2) NOT NULL DEFAULT 0,
    sort_order  INT NOT NULL DEFAULT 0,
    inactive    BOOLEAN NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by  UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by  UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT salary_slip_items_type_check CHECK (item_type IN ('earning','deduction')),
    CONSTRAINT salary_slip_items_amount_check CHECK (amount >= 0)
);
CREATE INDEX IF NOT EXISTS salary_slip_items_slip_idx ON hrms.salary_slip_items(slip_id);

CREATE TABLE IF NOT EXISTS hrms.salary_slip_leaves (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    slip_id         UUID NOT NULL REFERENCES hrms.salary_slips(id) ON DELETE CASCADE,
    leave_type_id   UUID NOT NULL REFERENCES hrms.leave_types(id) ON DELETE RESTRICT,
    leave_type_name VARCHAR(100),
    total_days      NUMERIC(5,1) NOT NULL DEFAULT 0,
    used_days       NUMERIC(5,1) NOT NULL DEFAULT 0,
    balance_days    NUMERIC(5,1) NOT NULL DEFAULT 0,
    inactive        BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by      UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by      UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT salary_slip_leaves_days_check CHECK (total_days >= 0 AND used_days >= 0)
);
CREATE INDEX IF NOT EXISTS salary_slip_leaves_slip_idx ON hrms.salary_slip_leaves(slip_id);

-- ============================================================
-- CELEBRATIONS
-- ============================================================

CREATE TABLE IF NOT EXISTS hrms.celebration_types (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id           UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    name                VARCHAR(100) NOT NULL,
    is_yearly           BOOLEAN NOT NULL DEFAULT TRUE,
    is_user_celebration BOOLEAN NOT NULL DEFAULT TRUE,
    inactive            BOOLEAN NOT NULL DEFAULT FALSE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by          UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by          UUID REFERENCES auth.users(id) ON DELETE SET NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS celebration_types_tenant_name_idx ON hrms.celebration_types(tenant_id, lower(name)) WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.celebrations (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id           UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    branch_id           UUID REFERENCES hrms.branches(id) ON DELETE SET NULL,
    user_id             UUID REFERENCES auth.users(id) ON DELETE CASCADE,
    celebration_type_id UUID NOT NULL REFERENCES hrms.celebration_types(id) ON DELETE RESTRICT,
    celebration_date    DATE,
    custom_title        VARCHAR(255),
    description         TEXT,
    inactive            BOOLEAN NOT NULL DEFAULT FALSE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by          UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by          UUID REFERENCES auth.users(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS celebrations_tenant_idx ON hrms.celebrations(tenant_id);
CREATE UNIQUE INDEX IF NOT EXISTS celebrations_user_type_idx ON hrms.celebrations(tenant_id, user_id, celebration_type_id) WHERE user_id IS NOT NULL AND NOT inactive;

-- ============================================================
-- NOTIFICATIONS
-- ============================================================

CREATE TABLE IF NOT EXISTS hrms.notification_masters (
    id                    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id             UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    code                  VARCHAR(100) NOT NULL,
    name                  VARCHAR(255),
    description           TEXT,
    is_in_app_enabled     BOOLEAN NOT NULL DEFAULT TRUE,
    is_email_enabled      BOOLEAN NOT NULL DEFAULT FALSE,
    is_push_enabled       BOOLEAN NOT NULL DEFAULT TRUE,
    inactive              BOOLEAN NOT NULL DEFAULT FALSE,
    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by            UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by            UUID REFERENCES auth.users(id) ON DELETE SET NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS notification_masters_code_idx ON hrms.notification_masters(tenant_id, code) WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.notification_preferences (
    id                     UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id              UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    user_id                UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    notification_master_id UUID NOT NULL REFERENCES hrms.notification_masters(id) ON DELETE CASCADE,
    is_in_app_enabled      BOOLEAN NOT NULL DEFAULT TRUE,
    is_email_enabled       BOOLEAN NOT NULL DEFAULT FALSE,
    inactive               BOOLEAN NOT NULL DEFAULT FALSE,
    created_at             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by             UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by             UUID REFERENCES auth.users(id) ON DELETE SET NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS notification_prefs_user_idx ON hrms.notification_preferences(tenant_id, user_id, notification_master_id) WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.notification_inbox (
    id                     UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id              UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    user_id                UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    notification_master_id UUID NOT NULL REFERENCES hrms.notification_masters(id) ON DELETE RESTRICT,
    title                  TEXT NOT NULL,
    message                TEXT NOT NULL,
    reference_table        VARCHAR(100),
    reference_id           UUID,
    is_read                BOOLEAN NOT NULL DEFAULT FALSE,
    read_date              TIMESTAMPTZ,
    inactive               BOOLEAN NOT NULL DEFAULT FALSE,
    created_at             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by             UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by             UUID REFERENCES auth.users(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS notification_inbox_user_idx ON hrms.notification_inbox(tenant_id, user_id, is_read);

CREATE TABLE IF NOT EXISTS hrms.notification_logs (
    id                     UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id              UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    notification_master_id UUID REFERENCES hrms.notification_masters(id) ON DELETE SET NULL,
    user_id                UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    channel                VARCHAR(20) NOT NULL,
    target_address         TEXT NOT NULL,
    subject                TEXT,
    body                   TEXT,
    status                 VARCHAR(20) NOT NULL DEFAULT 'Pending',
    sent_date              TIMESTAMPTZ,
    error_message          TEXT,
    external_reference_id  TEXT,
    inactive               BOOLEAN NOT NULL DEFAULT FALSE,
    created_at             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by             UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by             UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT notification_logs_channel_check CHECK (channel IN ('Push','Email')),
    CONSTRAINT notification_logs_status_check CHECK (status IN ('Pending','Sent','Failed'))
);
CREATE INDEX IF NOT EXISTS notification_logs_tenant_user_idx ON hrms.notification_logs(tenant_id, user_id);

CREATE TABLE IF NOT EXISTS hrms.device_tokens (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id     UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    user_id       UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    device_token  TEXT NOT NULL,
    device_type   VARCHAR(50),
    device_id     VARCHAR(255),
    inactive      BOOLEAN NOT NULL DEFAULT FALSE,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by    UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by    UUID REFERENCES auth.users(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS device_tokens_user_idx ON hrms.device_tokens(tenant_id, user_id) WHERE NOT inactive;
CREATE INDEX IF NOT EXISTS device_tokens_device_idx ON hrms.device_tokens(tenant_id, device_id) WHERE device_id IS NOT NULL AND NOT inactive;

-- ============================================================
-- ONBOARDING
-- ============================================================

CREATE TABLE IF NOT EXISTS hrms.job_positions (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id           UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    code                VARCHAR(50),
    title               VARCHAR(255) NOT NULL,
    level               VARCHAR(100),
    category            VARCHAR(100),
    description         TEXT,
    department_id       UUID REFERENCES hrms.departments(id) ON DELETE SET NULL,
    employment_type_id  UUID REFERENCES hrms.employment_types(id) ON DELETE SET NULL,
    work_mode           VARCHAR(50),
    total_position      INT NOT NULL DEFAULT 1,
    budgeted_cost       NUMERIC(12,2),
    inactive            BOOLEAN NOT NULL DEFAULT FALSE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by          UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by          UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT job_positions_total_position_check CHECK (total_position >= 0)
);
CREATE INDEX IF NOT EXISTS job_positions_tenant_idx ON hrms.job_positions(tenant_id);
CREATE UNIQUE INDEX IF NOT EXISTS job_positions_code_idx ON hrms.job_positions(tenant_id, code) WHERE code IS NOT NULL AND NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.job_position_locations (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    job_position_id UUID NOT NULL REFERENCES hrms.job_positions(id) ON DELETE CASCADE,
    location        VARCHAR(255),
    city            VARCHAR(100),
    state           VARCHAR(100),
    country         VARCHAR(100),
    is_remote       BOOLEAN NOT NULL DEFAULT FALSE,
    inactive        BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by      UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by      UUID REFERENCES auth.users(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS job_position_locations_position_idx ON hrms.job_position_locations(job_position_id);

CREATE TABLE IF NOT EXISTS hrms.job_requisitions (
    id                    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id             UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    job_position_id       UUID NOT NULL REFERENCES hrms.job_positions(id) ON DELETE RESTRICT,
    code                  VARCHAR(50),
    title                 VARCHAR(255) NOT NULL,
    level                 VARCHAR(100),
    category              VARCHAR(100),
    department_id         UUID REFERENCES hrms.departments(id) ON DELETE SET NULL,
    employment_type_id    UUID REFERENCES hrms.employment_types(id) ON DELETE SET NULL,
    description           TEXT,
    work_mode             VARCHAR(50),
    total_openings        INT NOT NULL DEFAULT 1,
    reason_for_hire       TEXT,
    min_salary            NUMERIC(12,2),
    max_salary            NUMERIC(12,2),
    currency              VARCHAR(10) NOT NULL DEFAULT 'INR',
    target_hire_date      DATE,
    expected_closure_date DATE,
    requested_by          UUID NOT NULL REFERENCES auth.users(id) ON DELETE RESTRICT,
    requested_date        DATE,
    is_approved           BOOLEAN NOT NULL DEFAULT FALSE,
    approved_by           UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    approved_date         TIMESTAMPTZ,
    priority              VARCHAR(50),
    status                VARCHAR(50) NOT NULL DEFAULT 'Draft',
    notes                 TEXT,
    inactive              BOOLEAN NOT NULL DEFAULT FALSE,
    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by            UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by            UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT job_requisitions_status_check CHECK (status IN ('Draft','Pending','Approved','Rejected','Closed')),
    CONSTRAINT job_requisitions_openings_check CHECK (total_openings > 0)
);
CREATE INDEX IF NOT EXISTS job_requisitions_tenant_idx ON hrms.job_requisitions(tenant_id);
CREATE UNIQUE INDEX IF NOT EXISTS job_requisitions_code_idx ON hrms.job_requisitions(tenant_id, code) WHERE code IS NOT NULL AND NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.job_postings (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id           UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    job_requisition_id  UUID REFERENCES hrms.job_requisitions(id) ON DELETE SET NULL,
    code                VARCHAR(50),
    title               VARCHAR(255),
    job_summary         TEXT,
    description         TEXT,
    job_category        VARCHAR(100),
    department_id       UUID REFERENCES hrms.departments(id) ON DELETE SET NULL,
    industry            VARCHAR(100),
    employment_type_id  UUID REFERENCES hrms.employment_types(id) ON DELETE SET NULL,
    work_mode           VARCHAR(50),
    role_type           VARCHAR(100),
    min_experience      NUMERIC(5,1),
    max_experience      NUMERIC(5,1),
    min_salary          NUMERIC(12,2),
    max_salary          NUMERIC(12,2),
    salary_currency     VARCHAR(10),
    salary_period       VARCHAR(30),
    is_salary_visible   BOOLEAN NOT NULL DEFAULT FALSE,
    job_status          VARCHAR(50),
    publish_date        DATE,
    expiry_date         DATE,
    is_published        BOOLEAN NOT NULL DEFAULT FALSE,
    inactive            BOOLEAN NOT NULL DEFAULT FALSE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by          UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by          UUID REFERENCES auth.users(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS job_postings_tenant_idx ON hrms.job_postings(tenant_id);
CREATE UNIQUE INDEX IF NOT EXISTS job_postings_code_idx ON hrms.job_postings(tenant_id, code) WHERE code IS NOT NULL AND NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.candidates (
    id                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id            UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    firstname            VARCHAR(100),
    lastname             VARCHAR(100),
    email                VARCHAR(255),
    phone                VARCHAR(30),
    dob                  DATE,
    gender               VARCHAR(20),
    total_experience     NUMERIC(5,1),
    current_company      VARCHAR(255),
    current_designation  VARCHAR(255),
    current_salary       NUMERIC(12,2),
    expected_salary      NUMERIC(12,2),
    notice_period        INT,
    current_location     VARCHAR(255),
    preferred_location   VARCHAR(255),
    source               VARCHAR(100),
    resume_url           TEXT,
    inactive             BOOLEAN NOT NULL DEFAULT FALSE,
    created_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by           UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by           UUID REFERENCES auth.users(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS candidates_tenant_idx ON hrms.candidates(tenant_id);
CREATE INDEX IF NOT EXISTS candidates_email_idx ON hrms.candidates(tenant_id, lower(email)) WHERE email IS NOT NULL AND NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.candidate_applications (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    candidate_id    UUID REFERENCES hrms.candidates(id) ON DELETE SET NULL,
    job_posting_id  UUID REFERENCES hrms.job_postings(id) ON DELETE SET NULL,
    resume_url      TEXT,
    cover_letter    TEXT,
    current_ctc     NUMERIC(12,2),
    expected_ctc    NUMERIC(12,2),
    notice_period   INT,
    referred_by     VARCHAR(128),
    source          VARCHAR(50),
    status          VARCHAR(50) NOT NULL DEFAULT 'New',
    comments        TEXT,
    inactive        BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by      UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by      UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT candidate_applications_status_check CHECK (status IN ('New','Screening','Interview','Offered','Hired','Rejected','Withdrawn'))
);
CREATE INDEX IF NOT EXISTS candidate_applications_tenant_idx ON hrms.candidate_applications(tenant_id);
CREATE INDEX IF NOT EXISTS candidate_applications_candidate_idx ON hrms.candidate_applications(candidate_id);

CREATE TABLE IF NOT EXISTS hrms.candidate_education (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    candidate_id    UUID NOT NULL REFERENCES hrms.candidates(id) ON DELETE CASCADE,
    institution     VARCHAR(255),
    degree          VARCHAR(255),
    field_of_study  VARCHAR(255),
    start_date      DATE,
    end_date        DATE,
    grade           VARCHAR(50),
    inactive        BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by      UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by      UUID REFERENCES auth.users(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS candidate_education_candidate_idx ON hrms.candidate_education(candidate_id);

CREATE TABLE IF NOT EXISTS hrms.candidate_experience (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    candidate_id    UUID NOT NULL REFERENCES hrms.candidates(id) ON DELETE CASCADE,
    company         VARCHAR(255),
    designation     VARCHAR(255),
    start_date      DATE,
    end_date        DATE,
    responsibilities TEXT,
    inactive        BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by      UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by      UUID REFERENCES auth.users(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS candidate_experience_candidate_idx ON hrms.candidate_experience(candidate_id);

CREATE TABLE IF NOT EXISTS hrms.interview_rounds (
    id                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id            UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    application_id       UUID NOT NULL REFERENCES hrms.candidate_applications(id) ON DELETE CASCADE,
    round_name           VARCHAR(100),
    round_number         INT,
    scheduled_date       TIMESTAMPTZ,
    duration_minutes     INT,
    interviewer_user_id  UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    mode                 VARCHAR(50),
    meeting_link         TEXT,
    location             VARCHAR(255),
    status               VARCHAR(50),
    remarks              TEXT,
    inactive             BOOLEAN NOT NULL DEFAULT FALSE,
    created_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by           UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by           UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT interview_rounds_duration_check CHECK (duration_minutes IS NULL OR duration_minutes > 0)
);
CREATE INDEX IF NOT EXISTS interview_rounds_application_idx ON hrms.interview_rounds(application_id);

CREATE TABLE IF NOT EXISTS hrms.offer_letters (
    id                         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id                  UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    application_id             UUID NOT NULL REFERENCES hrms.candidate_applications(id) ON DELETE CASCADE,
    candidate_id               UUID REFERENCES hrms.candidates(id) ON DELETE SET NULL,
    offered_ctc                NUMERIC(12,2),
    currency                   VARCHAR(10) NOT NULL DEFAULT 'INR',
    salary_breakdown           JSONB,
    joining_date               DATE,
    valid_until_date           DATE,
    status                     VARCHAR(50) NOT NULL DEFAULT 'Generated',
    offer_letter_url           TEXT,
    candidate_reaction_date    TIMESTAMPTZ,
    candidate_rejection_reason TEXT,
    version                    INT NOT NULL DEFAULT 1,
    is_latest                  BOOLEAN NOT NULL DEFAULT FALSE,
    inactive                   BOOLEAN NOT NULL DEFAULT FALSE,
    created_at                 TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by                 UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at                 TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by                 UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT offer_letters_status_check CHECK (status IN ('Generated','Sent','Accepted','Declined','Revoked')),
    CONSTRAINT offer_letters_version_check CHECK (version > 0)
);
CREATE INDEX IF NOT EXISTS offer_letters_application_idx ON hrms.offer_letters(application_id);
CREATE UNIQUE INDEX IF NOT EXISTS offer_letters_latest_idx ON hrms.offer_letters(application_id) WHERE is_latest AND NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.onboarding_workflows (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    inactive    BOOLEAN NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by  UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by  UUID REFERENCES auth.users(id) ON DELETE SET NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS onboarding_workflows_tenant_name_idx ON hrms.onboarding_workflows(tenant_id, lower(name)) WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.onboarding_tasks (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    workflow_id     UUID NOT NULL REFERENCES hrms.onboarding_workflows(id) ON DELETE CASCADE,
    title           VARCHAR(255) NOT NULL,
    description     TEXT,
    due_days        INT NOT NULL DEFAULT 0,
    is_required     BOOLEAN NOT NULL DEFAULT TRUE,
    sort_order      INT NOT NULL DEFAULT 0,
    inactive        BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by      UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by      UUID REFERENCES auth.users(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS onboarding_tasks_workflow_idx ON hrms.onboarding_tasks(workflow_id);

CREATE TABLE IF NOT EXISTS hrms.candidate_onboardings (
    id                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id            UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    candidate_id         UUID NOT NULL REFERENCES hrms.candidates(id) ON DELETE CASCADE,
    workflow_id          UUID NOT NULL REFERENCES hrms.onboarding_workflows(id) ON DELETE RESTRICT,
    onboarding_status    VARCHAR(50) NOT NULL DEFAULT 'Pending',
    progress_percentage  INT NOT NULL DEFAULT 0,
    inactive             BOOLEAN NOT NULL DEFAULT FALSE,
    created_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by           UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by           UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT candidate_onboardings_status_check CHECK (onboarding_status IN ('Pending','InProgress','Completed')),
    CONSTRAINT candidate_onboardings_progress_check CHECK (progress_percentage BETWEEN 0 AND 100)
);
CREATE UNIQUE INDEX IF NOT EXISTS candidate_onboardings_candidate_idx ON hrms.candidate_onboardings(tenant_id, candidate_id) WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.candidate_onboarding_tasks (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id               UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    candidate_onboarding_id UUID NOT NULL REFERENCES hrms.candidate_onboardings(id) ON DELETE CASCADE,
    onboarding_task_id      UUID NOT NULL REFERENCES hrms.onboarding_tasks(id) ON DELETE RESTRICT,
    status                  VARCHAR(50) NOT NULL DEFAULT 'Pending',
    completed_at            TIMESTAMPTZ,
    remarks                 TEXT,
    inactive                BOOLEAN NOT NULL DEFAULT FALSE,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by              UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by              UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT candidate_onboarding_tasks_status_check CHECK (status IN ('Pending','InProgress','Completed'))
);
CREATE INDEX IF NOT EXISTS candidate_onboarding_tasks_onboarding_idx ON hrms.candidate_onboarding_tasks(candidate_onboarding_id);

-- ============================================================
-- AUTH SUPPORT OWNED BY HRMS FLOWS
-- ============================================================

CREATE TABLE IF NOT EXISTS hrms.user_otps (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    user_id     UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    otp         VARCHAR(10) NOT NULL,
    otp_for     VARCHAR(50) NOT NULL,
    mobile      VARCHAR(20),
    is_used     BOOLEAN NOT NULL DEFAULT FALSE,
    expires_at  TIMESTAMPTZ NOT NULL,
    inactive    BOOLEAN NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by  UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by  UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT user_otps_for_check CHECK (otp_for IN ('passwordreset','login'))
);
CREATE INDEX IF NOT EXISTS user_otps_user_idx ON hrms.user_otps(tenant_id, user_id);
CREATE INDEX IF NOT EXISTS user_otps_expires_idx ON hrms.user_otps(expires_at) WHERE NOT is_used AND NOT inactive;

-- ============================================================
-- UPDATED_AT TRIGGERS AND RLS
-- ============================================================

DO $$
DECLARE
    table_name TEXT;
    updated_tables TEXT[] := ARRAY[
        'departments','branches','designations','employment_types','marital_statuses',
        'financial_years','working_hours','holidays','pay_cycles','policy_types','company_policies',
        'tenant_subscriptions','employees','employee_statutory','employee_banks','document_types',
        'employee_documents','leave_types','leave_policies','leave_balances','leaves','leave_approvals',
        'attendances','attendance_requests','device_logs','salary_templates','salary_template_items',
        'employee_salaries','employee_salary_structures','salary_slips','salary_slip_items',
        'salary_slip_leaves','celebration_types','celebrations','notification_masters',
        'notification_preferences','notification_inbox','notification_logs','device_tokens','job_positions',
        'job_position_locations','job_requisitions','job_postings','candidates','candidate_applications',
        'candidate_education','candidate_experience','interview_rounds','offer_letters','onboarding_workflows',
        'onboarding_tasks','candidate_onboardings','candidate_onboarding_tasks','user_otps'
    ];
BEGIN
    FOREACH table_name IN ARRAY updated_tables LOOP
        EXECUTE format('DROP TRIGGER IF EXISTS %I ON hrms.%I', table_name || '_updated_at', table_name);
        EXECUTE format(
            'CREATE TRIGGER %I BEFORE UPDATE ON hrms.%I FOR EACH ROW EXECUTE FUNCTION hrms.update_updated_at()',
            table_name || '_updated_at',
            table_name
        );
    END LOOP;
END $$;

DO $$
DECLARE
    table_name TEXT;
    tenant_tables TEXT[] := ARRAY[
        'tenant_profiles','tenant_settings','departments','branches',
        'designations','employment_types','marital_statuses','financial_years','working_hours','holidays',
        'pay_cycles','policy_types','company_policies','tenant_subscriptions','employees','employee_statutory',
        'employee_banks','document_types','employee_documents','leave_types','leave_policies','leave_balances',
        'leaves','leave_approvals','leave_ledger','attendances','attendance_requests','device_logs',
        'salary_templates','salary_template_items','employee_salaries','employee_salary_structures',
        'salary_slips','salary_slip_items','salary_slip_leaves','celebration_types','celebrations',
        'notification_masters','notification_preferences','notification_inbox','notification_logs','device_tokens',
        'job_positions','job_position_locations','job_requisitions','job_postings','candidates',
        'candidate_applications','candidate_education','candidate_experience','interview_rounds','offer_letters',
        'onboarding_workflows','onboarding_tasks','candidate_onboardings','candidate_onboarding_tasks','user_otps'
    ];
BEGIN
    FOREACH table_name IN ARRAY tenant_tables LOOP
        EXECUTE format('ALTER TABLE hrms.%I ENABLE ROW LEVEL SECURITY', table_name);
        EXECUTE format('DROP POLICY IF EXISTS %I ON hrms.%I', table_name || '_tenant_isolation', table_name);
        EXECUTE format(
            'CREATE POLICY %I ON hrms.%I USING (tenant_id = NULLIF(current_setting(''app.tenant_id'', true), '''')::uuid OR current_setting(''app.is_super_admin'', true) = ''true'') WITH CHECK (tenant_id = NULLIF(current_setting(''app.tenant_id'', true), '''')::uuid OR current_setting(''app.is_super_admin'', true) = ''true'')',
            table_name || '_tenant_isolation',
            table_name
        );
    END LOOP;
END $$;
