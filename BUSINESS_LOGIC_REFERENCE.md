# HRMS Business Logic Reference — v2
> Complete extraction from SetikaHRMS .NET Core. Covers all SQL, business rules, Go module structure, CREATE TABLE migrations, and agent guardrails.

---

## TABLE OF CONTENTS
1. [Enums & Constants](#enums)
2. [Table Map](#table-map)
3. [CREATE TABLE Migrations](#migrations)
4. [Go Module Structure](#module-structure)
5. [Module: Auth / Tenant Provisioning](#module-auth)
6. [Module: Employee](#module-employee)
7. [Module: Leave](#module-leave)
8. [Module: Attendance](#module-attendance)
9. [Module: Salary](#module-salary)
10. [Module: PDF Salary Slip](#module-pdf)
11. [Module: Dashboard](#module-dashboard)
12. [Module: Celebrations](#module-celebrations)
13. [Module: Notifications](#module-notifications)
14. [Module: Onboarding](#module-onboarding)
15. [External Integrations](#integrations)
16. [Password Migration](#password-migration)
17. [Agent Guardrails](#guardrails)

---

## 1. ENUMS & CONSTANTS {#enums}

```go
// Roles
const (
    RoleSuperAdmin = "superadmin"
    RoleTenant     = "tenant"   // company admin
    RoleHR         = "hr"
    RoleManager    = "manager"
    RoleEmployee   = "employee"
)

// Leave status
const (
    LeaveStatusPending  = "pending"
    LeaveStatusApproved = "approved"
    LeaveStatusRejected = "rejected"
    LeaveStatusCanceled = "canceled"
)

// Leave day type
const (
    LeaveDayFullDay    = "fullday"
    LeaveDayFirstHalf  = "firsthalf"
    LeaveDaySecondHalf = "secondhalf"
)

// Leave type short names (special codes)
const LeaveTypeShortEarnLeave = "earnleave"  // triggers auto-conversion logic

// Attendance type
const (
    AttendanceCheckin  = "checkin"
    AttendanceCheckout = "checkout"
)

// Attendance status
const (
    AttendanceStatusPresent = "present"
    AttendanceStatusLeave   = "leave"
    AttendanceStatusAbsent  = "absent"
    AttendanceStatusHoliday = "holiday"
    AttendanceStatusWeekoff = "weekoff"
)

// Leave allocation type
const (
    AllocationFixed   = "fixed"
    AllocationMonthly = "monthly"
)

// Notification codes
const (
    NotifLeaveApplied        = "leaveapplied"
    NotifLeaveApproved       = "leaveapproved"
    NotifLeaveRejected       = "leaverejected"
    NotifCompanyPolicy       = "companypolicy"
    NotifGeneralNotif        = "generalnotification"
    NotifUserCelebration     = "usercelebration"
)

// Notification channels
const (
    NotifChannelPush  = "Push"
    NotifChannelEmail = "Email"
)

// Notification status
const (
    NotifStatusPending = "Pending"
    NotifStatusSent    = "Sent"
    NotifStatusFailed  = "Failed"
)

// Notification reference tables
const (
    RefTableUserLeaves      = "userleaves"
    RefTableCompanyPolicy   = "companypolicy"
    RefTableUserCelebration = "usercelebration"
)

// Transaction type
const (
    TransactionDebit  = "debit"
    TransactionCredit = "credit"
)

// Regularisation key
const (
    RangeKeyHalfday = "halfday"
    RangeKeyAbsent  = "absent"
)

// OTP for
const (
    OtpForPasswordReset = "passwordreset"
    OtpForLogin         = "login"
)

// Salary item types
const (
    SalaryItemEarning   = "earning"
    SalaryItemDeduction = "deduction"
)

// Reserved salary item codes
const (
    SalaryCodeBasic = "basic"
    SalaryCodeHRA   = "hra"
    SalaryCodeLWP   = "LWP"  // Leave Without Pay — dynamically added
)

// Dashboard quick tools
const (
    DashToolLeaveRequest   = "leaverequest"
    DashToolAttendance     = "attendancerecord"
    DashToolForgotPunch    = "forgottopunch"
    DashToolLeaveApproval  = "leaveapproval"  // manager only
    DashToolCelebration    = "celebration"
    DashToolHolidays       = "holidays"
    DashToolPolicies       = "policies"
)

// Onboarding application status
const (
    AppStatusNew        = "New"
    AppStatusScreening  = "Screening"
    AppStatusInterview  = "Interview"
    AppStatusOffered    = "Offered"
    AppStatusHired      = "Hired"
    AppStatusRejected   = "Rejected"
    AppStatusWithdrawn  = "Withdrawn"
)

// Offer letter status
const (
    OfferStatusGenerated = "Generated"
    OfferStatusSent      = "Sent"
    OfferStatusAccepted  = "Accepted"
    OfferStatusDeclined  = "Declined"
    OfferStatusRevoked   = "Revoked"
)

// Job requisition status
const (
    ReqStatusDraft     = "Draft"
    ReqStatusPending   = "Pending"
    ReqStatusApproved  = "Approved"
    ReqStatusRejected  = "Rejected"
    ReqStatusClosed    = "Closed"
)

// Onboarding status
const (
    OnboardStatusPending    = "Pending"
    OnboardStatusInProgress = "InProgress"
    OnboardStatusCompleted  = "Completed"
)
```

---

## 2. TABLE MAP {#table-map}

| .NET Entity Name | Go Table Name | Schema | Owned By Module |
|---|---|---|---|
| Tenants | tenants | public | auth |
| ApplicationUser (AspNetUsers) | users | auth | auth (spur-identity) |
| AspNetRoles | roles | auth | auth (spur-identity) |
| AspNetUserRoles | user_roles | auth | auth (spur-identity) |
| TenantBranding | tenant_brandings | public | tenant |
| UserInfo | employees | public | employee |
| UserStatutoryInfo | employee_statutory | public | employee |
| UserBankInfo | employee_banks | public | employee |
| Departments | departments | public | tenant |
| Branch | branches | public | tenant |
| Designation | designations | public | tenant |
| FinancialYear | financial_years | public | tenant |
| LeaveType | leave_types | public | leave |
| LeavePolicy | leave_policies | public | leave |
| UserLeaveBalance | leave_balances | public | leave |
| UserLeaveLedger | leave_ledger | public | leave |
| UserLeaves | leaves | public | leave |
| UserLeaveApproval | leave_approvals | public | leave |
| Holiday | holidays | public | tenant |
| WorkingHour | working_hours | public | tenant |
| PayCycleConfig | pay_cycles | public | salary |
| SalaryTemplate | salary_templates | public | salary |
| SalaryTemplateItem | salary_template_items | public | salary |
| UserSalary | employee_salaries | public | salary |
| UserSalaryStructure | employee_salary_structures | public | salary |
| UserSalarySlip | salary_slips | public | salary |
| UserSalarySlipItems | salary_slip_items | public | salary |
| UserSalarySlipLeave | salary_slip_leaves | public | salary |
| UserAttendance | attendances | public | attendance |
| UserAttendanceRequest | attendance_requests | public | attendance |
| UserDeviceLog | device_logs | public | attendance |
| CompanyPolicy | company_policies | public | tenant |
| PolicyType | policy_types | public | tenant |
| DocumentType | document_types | public | employee |
| UserDocument | employee_documents | public | employee |
| UserCelebration | celebrations | public | celebration |
| CelebrationType | celebration_types | public | celebration |
| NotificationMaster | notification_masters | public | notification |
| NotificationPreference | notification_preferences | public | notification |
| NotificationInbox | notification_inbox | public | notification |
| NotificationLog | notification_logs | public | notification |
| UserDeviceToken | device_tokens | public | notification |
| JobPosition | job_positions | public | onboarding |
| JobPositionLocation | job_position_locations | public | onboarding |
| JobPosting | job_postings | public | onboarding |
| JobRequisition | job_requisitions | public | onboarding |
| JobRequisitionLog | job_requisition_logs | public | onboarding |
| Candidate | candidates | public | onboarding |
| CandidateApplication | candidate_applications | public | onboarding |
| CandidateEducation | candidate_education | public | onboarding |
| CandidateExperience | candidate_experience | public | onboarding |
| InterviewRound | interview_rounds | public | onboarding |
| OfferLetter | offer_letters | public | onboarding |
| CandidateOnboarding | candidate_onboardings | public | onboarding |
| CandidateOnboardingTask | candidate_onboarding_tasks | public | onboarding |
| OnboardingTask | onboarding_tasks | public | onboarding |
| OnboardingWorkflow | onboarding_workflows | public | onboarding |
| UserOtp | user_otps | public | auth |
| TenantSubscription | tenant_subscriptions | public | tenant |
| EmploymentType | employment_types | public | tenant |
| MaritalStatus | marital_statuses | public | tenant |

---

## 3. CREATE TABLE MIGRATIONS {#migrations}

> All tables use `tenant_id UUID NOT NULL` + RLS. `inactive BOOLEAN NOT NULL DEFAULT FALSE` is the soft-delete flag.
> Run migrations in order. auth schema tables are owned by spur-identity — do NOT recreate them, only reference them.

```sql
-- ============================================================
-- TENANT / COMPANY SETUP
-- ============================================================

CREATE TABLE IF NOT EXISTS tenant_brandings (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id         UUID NOT NULL,
    company_name      TEXT,
    domain            TEXT UNIQUE,
    logo_path         TEXT,
    primary_color     VARCHAR(10),    -- hex e.g. "#2C435B"
    address1          TEXT,
    address2          TEXT,
    city              VARCHAR(100),
    state             VARCHAR(100),
    country           VARCHAR(100),
    pincode           VARCHAR(20),
    phone             VARCHAR(30),
    email             VARCHAR(255),
    website           TEXT,
    gstin             VARCHAR(50),
    inactive          BOOLEAN NOT NULL DEFAULT FALSE,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by        TEXT,
    updated_at        TIMESTAMPTZ,
    updated_by        TEXT
);
CREATE INDEX idx_tenant_brandings_tenant_id ON tenant_brandings(tenant_id);
CREATE INDEX idx_tenant_brandings_domain    ON tenant_brandings(domain);

CREATE TABLE IF NOT EXISTS departments (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id      UUID NOT NULL,
    name           VARCHAR(255) NOT NULL,
    description    TEXT,
    inactive       BOOLEAN NOT NULL DEFAULT FALSE,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by     TEXT,
    updated_at     TIMESTAMPTZ,
    updated_by     TEXT
);
CREATE INDEX idx_departments_tenant_id ON departments(tenant_id);

CREATE TABLE IF NOT EXISTS branches (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id      UUID NOT NULL,
    branch_name    VARCHAR(255) NOT NULL,
    address        TEXT,
    city           VARCHAR(100),
    state          VARCHAR(100),
    country        VARCHAR(100),
    pincode        VARCHAR(20),
    phone          VARCHAR(30),
    inactive       BOOLEAN NOT NULL DEFAULT FALSE,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by     TEXT,
    updated_at     TIMESTAMPTZ,
    updated_by     TEXT
);
CREATE INDEX idx_branches_tenant_id ON branches(tenant_id);

CREATE TABLE IF NOT EXISTS designations (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id      UUID NOT NULL,
    name           VARCHAR(255) NOT NULL,
    description    TEXT,
    inactive       BOOLEAN NOT NULL DEFAULT FALSE,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by     TEXT,
    updated_at     TIMESTAMPTZ,
    updated_by     TEXT
);
CREATE INDEX idx_designations_tenant_id ON designations(tenant_id);

CREATE TABLE IF NOT EXISTS employment_types (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL,
    name        VARCHAR(100) NOT NULL,
    inactive    BOOLEAN NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by  TEXT,
    updated_at  TIMESTAMPTZ,
    updated_by  TEXT
);

CREATE TABLE IF NOT EXISTS financial_years (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id     UUID NOT NULL,
    name          VARCHAR(100),
    start_date    DATE NOT NULL,
    end_date      DATE NOT NULL,
    is_active     BOOLEAN NOT NULL DEFAULT FALSE,
    inactive      BOOLEAN NOT NULL DEFAULT FALSE,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by    TEXT,
    updated_at    TIMESTAMPTZ,
    updated_by    TEXT
);
CREATE INDEX idx_financial_years_tenant_id ON financial_years(tenant_id);

CREATE TABLE IF NOT EXISTS working_hours (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL,
    branch_id       UUID,                     -- NULL = tenant-level default
    user_id         TEXT,                     -- NULL = not user-specific
    day_of_week     VARCHAR(20) NOT NULL,     -- "Monday".."Sunday"
    is_working_day  BOOLEAN NOT NULL DEFAULT TRUE,
    start_time      TIME NOT NULL,
    end_time        TIME NOT NULL,
    break_minutes   INT DEFAULT 0,
    inactive        BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by      TEXT,
    updated_at      TIMESTAMPTZ,
    updated_by      TEXT
);
CREATE INDEX idx_working_hours_tenant_id ON working_hours(tenant_id);

CREATE TABLE IF NOT EXISTS holidays (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    UUID NOT NULL,
    branch_id    UUID,            -- NULL = applies to all branches
    fy_id        UUID,            -- links to financial_years
    name         VARCHAR(255) NOT NULL,
    date         DATE NOT NULL,
    is_optional  BOOLEAN NOT NULL DEFAULT FALSE,
    inactive     BOOLEAN NOT NULL DEFAULT FALSE,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by   TEXT,
    updated_at   TIMESTAMPTZ,
    updated_by   TEXT
);
CREATE INDEX idx_holidays_tenant_id ON holidays(tenant_id);
CREATE INDEX idx_holidays_date      ON holidays(date);

CREATE TABLE IF NOT EXISTS pay_cycles (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id         UUID NOT NULL UNIQUE,
    cycle_type        VARCHAR(50),        -- "monthly"
    pay_day           INT,                -- day of month salary is paid
    start_day         INT,                -- attendance cycle start (e.g. 26)
    end_day           INT,                -- attendance cycle end (e.g. 25)
    inactive          BOOLEAN NOT NULL DEFAULT FALSE,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by        TEXT,
    updated_at        TIMESTAMPTZ,
    updated_by        TEXT
);

CREATE TABLE IF NOT EXISTS policy_types (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL,
    name        VARCHAR(255) NOT NULL,
    inactive    BOOLEAN NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by  TEXT
);

CREATE TABLE IF NOT EXISTS company_policies (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL,
    policy_type_id  UUID,
    title           VARCHAR(255) NOT NULL,
    file_path       TEXT,
    description     TEXT,
    inactive        BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by      TEXT,
    updated_at      TIMESTAMPTZ,
    updated_by      TEXT
);
CREATE INDEX idx_company_policies_tenant_id ON company_policies(tenant_id);

CREATE TABLE IF NOT EXISTS tenant_subscriptions (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id        UUID NOT NULL,
    plan_id          UUID,
    start_date       DATE,
    end_date         DATE,
    status           VARCHAR(50),
    max_employees    INT DEFAULT 0,
    inactive         BOOLEAN NOT NULL DEFAULT FALSE,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by       TEXT,
    updated_at       TIMESTAMPTZ,
    updated_by       TEXT
);

-- ============================================================
-- EMPLOYEE
-- ============================================================

CREATE TABLE IF NOT EXISTS employees (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id           UUID NOT NULL,
    user_id             TEXT NOT NULL,          -- FK -> auth.users.id
    employee_code       VARCHAR(50),
    firstname           VARCHAR(100) NOT NULL,
    lastname            VARCHAR(100),
    email               VARCHAR(255),
    mobile              VARCHAR(20),
    dob                 DATE,
    gender              VARCHAR(20),
    marital_status      VARCHAR(30),
    blood_group         VARCHAR(10),
    profile_photo_path  TEXT,
    address             TEXT,
    city                VARCHAR(100),
    state               VARCHAR(100),
    country             VARCHAR(100),
    pincode             VARCHAR(20),
    emergency_contact   VARCHAR(255),
    joining_date        DATE,
    resignation_date    DATE,
    department_id       UUID,
    branch_id           UUID,
    designation_id      UUID,
    reporting_manager_id TEXT,                 -- FK -> auth.users.id
    employment_type_id  UUID,
    role                VARCHAR(50),
    grade               VARCHAR(50),
    experience_year     INT DEFAULT 0,
    experience_month    INT DEFAULT 0,
    inactive            BOOLEAN NOT NULL DEFAULT FALSE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by          TEXT,
    updated_at          TIMESTAMPTZ,
    updated_by          TEXT
);
CREATE INDEX idx_employees_tenant_id  ON employees(tenant_id);
CREATE INDEX idx_employees_user_id    ON employees(user_id);
CREATE UNIQUE INDEX idx_employees_code_tenant ON employees(tenant_id, employee_code) WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS employee_statutory (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL,
    user_id         TEXT NOT NULL,
    pf_no           VARCHAR(100),
    uan_no          VARCHAR(100),
    esic_no         VARCHAR(100),
    pan             VARCHAR(20),
    aadhaar         VARCHAR(20),
    pt_applicable   BOOLEAN DEFAULT FALSE,
    pf_applicable   BOOLEAN DEFAULT FALSE,
    esic_applicable BOOLEAN DEFAULT FALSE,
    inactive        BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by      TEXT,
    updated_at      TIMESTAMPTZ,
    updated_by      TEXT
);
CREATE UNIQUE INDEX idx_employee_statutory_user_tenant ON employee_statutory(tenant_id, user_id) WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS employee_banks (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL,
    user_id         TEXT NOT NULL,
    bank_name       VARCHAR(255),
    account_number  VARCHAR(50),
    ifsc_code       VARCHAR(20),
    account_type    VARCHAR(50),
    branch_name     VARCHAR(255),
    is_primary      BOOLEAN DEFAULT TRUE,
    inactive        BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by      TEXT,
    updated_at      TIMESTAMPTZ,
    updated_by      TEXT
);
CREATE INDEX idx_employee_banks_tenant_user ON employee_banks(tenant_id, user_id);

CREATE TABLE IF NOT EXISTS document_types (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL,
    name        VARCHAR(255) NOT NULL,
    inactive    BOOLEAN NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS employee_documents (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id        UUID NOT NULL,
    user_id          TEXT NOT NULL,
    document_type_id UUID,
    title            VARCHAR(255),
    file_path        TEXT,
    inactive         BOOLEAN NOT NULL DEFAULT FALSE,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by       TEXT,
    updated_at       TIMESTAMPTZ,
    updated_by       TEXT
);
CREATE INDEX idx_employee_documents_tenant_user ON employee_documents(tenant_id, user_id);

-- ============================================================
-- LEAVE
-- ============================================================

CREATE TABLE IF NOT EXISTS leave_types (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL,
    name            VARCHAR(100) NOT NULL,
    shortcode       VARCHAR(20),            -- e.g. "CL", "SL", "EL"
    description     TEXT,
    is_paid         BOOLEAN DEFAULT TRUE,
    is_carry_forward BOOLEAN DEFAULT FALSE,
    max_carry_forward INT DEFAULT 0,
    is_consecutive_limit BOOLEAN DEFAULT FALSE,
    consecutive_days_limit INT DEFAULT 0,
    inactive        BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by      TEXT,
    updated_at      TIMESTAMPTZ,
    updated_by      TEXT
);
CREATE INDEX idx_leave_types_tenant_id ON leave_types(tenant_id);

CREATE TABLE IF NOT EXISTS leave_policies (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id           UUID NOT NULL,
    leave_type_id       UUID NOT NULL,
    fy_id               UUID NOT NULL,
    total_days          NUMERIC(5,1) NOT NULL DEFAULT 0,
    allocation_type     VARCHAR(20) NOT NULL DEFAULT 'fixed',  -- 'fixed' | 'monthly'
    jan INT DEFAULT 0, feb INT DEFAULT 0, mar INT DEFAULT 0,
    apr INT DEFAULT 0, may INT DEFAULT 0, jun INT DEFAULT 0,
    jul INT DEFAULT 0, aug INT DEFAULT 0, sep INT DEFAULT 0,
    oct INT DEFAULT 0, nov INT DEFAULT 0, dec INT DEFAULT 0,
    is_sandwich_applicable BOOLEAN DEFAULT FALSE,
    inactive            BOOLEAN NOT NULL DEFAULT FALSE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by          TEXT,
    updated_at          TIMESTAMPTZ,
    updated_by          TEXT
);
CREATE UNIQUE INDEX idx_leave_policy_unique ON leave_policies(tenant_id, leave_type_id, fy_id) WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS leave_balances (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL,
    user_id         TEXT NOT NULL,
    leave_type_id   UUID NOT NULL,
    fy_id           UUID NOT NULL,
    total_days      NUMERIC(5,1) NOT NULL DEFAULT 0,
    used_days       NUMERIC(5,1) NOT NULL DEFAULT 0,
    pending_days    NUMERIC(5,1) NOT NULL DEFAULT 0,
    balance_days    NUMERIC(5,1) GENERATED ALWAYS AS (total_days - used_days - pending_days) STORED,
    inactive        BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by      TEXT,
    updated_at      TIMESTAMPTZ,
    updated_by      TEXT
);
CREATE UNIQUE INDEX idx_leave_balance_unique ON leave_balances(tenant_id, user_id, leave_type_id, fy_id) WHERE NOT inactive;
CREATE INDEX idx_leave_balances_user ON leave_balances(tenant_id, user_id);

CREATE TABLE IF NOT EXISTS leaves (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id           UUID NOT NULL,
    user_id             TEXT NOT NULL,
    leave_type_id       UUID NOT NULL,
    fy_id               UUID NOT NULL,
    start_date          DATE NOT NULL,
    end_date            DATE NOT NULL,
    start_day_type      VARCHAR(20) DEFAULT 'fullday',
    end_day_type        VARCHAR(20) DEFAULT 'fullday',
    days                NUMERIC(5,1),
    reason              TEXT,
    status              VARCHAR(20) NOT NULL DEFAULT 'pending',
    applied_date        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    from_leave_type     UUID,               -- for auto-conversion: original leave type
    to_leave_type       UUID,               -- for auto-conversion: new leave type
    is_sandwich         BOOLEAN DEFAULT FALSE,
    inactive            BOOLEAN NOT NULL DEFAULT FALSE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by          TEXT,
    updated_at          TIMESTAMPTZ,
    updated_by          TEXT
);
CREATE INDEX idx_leaves_tenant_user     ON leaves(tenant_id, user_id);
CREATE INDEX idx_leaves_dates           ON leaves(start_date, end_date);
CREATE INDEX idx_leaves_status          ON leaves(status);

CREATE TABLE IF NOT EXISTS leave_approvals (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL,
    leave_id        UUID NOT NULL,
    approver_id     TEXT NOT NULL,          -- FK -> auth.users.id
    status          VARCHAR(20) NOT NULL DEFAULT 'pending',
    remarks         TEXT,
    action_date     TIMESTAMPTZ,
    inactive        BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by      TEXT,
    updated_at      TIMESTAMPTZ,
    updated_by      TEXT
);
CREATE INDEX idx_leave_approvals_leave_id    ON leave_approvals(leave_id);
CREATE INDEX idx_leave_approvals_approver_id ON leave_approvals(approver_id);

CREATE TABLE IF NOT EXISTS leave_ledger (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL,
    user_id         TEXT NOT NULL,
    leave_type_id   UUID NOT NULL,
    fy_id           UUID NOT NULL,
    leave_id        UUID,                   -- nullable, for context
    transaction_type VARCHAR(10) NOT NULL,  -- 'debit' | 'credit'
    days            NUMERIC(5,1) NOT NULL,
    remarks         TEXT,
    inactive        BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by      TEXT
);
CREATE INDEX idx_leave_ledger_user ON leave_ledger(tenant_id, user_id);

-- ============================================================
-- ATTENDANCE
-- ============================================================

CREATE TABLE IF NOT EXISTS attendances (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    UUID NOT NULL,
    user_id      TEXT NOT NULL,
    date         DATE NOT NULL,
    time         TIMESTAMPTZ,
    type         VARCHAR(20),               -- 'checkin' | 'checkout'
    status       VARCHAR(30),               -- 'present' | 'leave' | 'absent' | 'holiday' | 'weekoff'
    source       VARCHAR(50),               -- 'app' | 'web' | 'biometric'
    latitude     NUMERIC(10,7),
    longitude    NUMERIC(10,7),
    work_mode    VARCHAR(50),               -- 'wfh' | 'wfo'
    remarks      TEXT,
    inactive     BOOLEAN NOT NULL DEFAULT FALSE,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by   TEXT,
    updated_at   TIMESTAMPTZ,
    updated_by   TEXT
);
CREATE INDEX idx_attendances_tenant_user_date ON attendances(tenant_id, user_id, date);

CREATE TABLE IF NOT EXISTS attendance_requests (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL,
    user_id         TEXT NOT NULL,
    date            DATE NOT NULL,
    requested_type  VARCHAR(20),            -- 'halfday' | 'absent'
    reason          TEXT,
    status          VARCHAR(20) DEFAULT 'pending',
    reviewed_by     TEXT,
    reviewed_at     TIMESTAMPTZ,
    remarks         TEXT,
    inactive        BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by      TEXT
);
CREATE INDEX idx_attendance_requests_user ON attendance_requests(tenant_id, user_id);

CREATE TABLE IF NOT EXISTS device_logs (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL,
    user_id     TEXT NOT NULL,
    device_id   VARCHAR(255),
    device_type VARCHAR(50),
    ip_address  VARCHAR(45),
    action      VARCHAR(50),
    logged_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================
-- SALARY
-- ============================================================

CREATE TABLE IF NOT EXISTS salary_templates (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL,
    fy_id       UUID NOT NULL,
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    is_active   BOOLEAN NOT NULL DEFAULT FALSE,   -- only ONE active per FY per tenant
    inactive    BOOLEAN NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by  TEXT,
    updated_at  TIMESTAMPTZ,
    updated_by  TEXT
);
CREATE INDEX idx_salary_templates_tenant_fy ON salary_templates(tenant_id, fy_id);

CREATE TABLE IF NOT EXISTS salary_template_items (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL,
    template_id     UUID NOT NULL,
    item_type       VARCHAR(20) NOT NULL,   -- 'earning' | 'deduction'
    code            VARCHAR(50) NOT NULL,   -- reserved: 'basic', 'hra', 'LWP'
    name            VARCHAR(255) NOT NULL,
    percentage      NUMERIC(7,2),           -- % of gross (nullable)
    amount          NUMERIC(12,2),          -- fixed amount (nullable)
    is_tax_exempt   BOOLEAN DEFAULT FALSE,
    sort_order      INT DEFAULT 0,
    inactive        BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by      TEXT,
    updated_at      TIMESTAMPTZ,
    updated_by      TEXT
);
CREATE INDEX idx_salary_template_items_template ON salary_template_items(template_id);

CREATE TABLE IF NOT EXISTS employee_salaries (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL,
    user_id         TEXT NOT NULL,
    fy_id           UUID NOT NULL,
    template_id     UUID NOT NULL,
    gross_salary    NUMERIC(12,2) NOT NULL,
    effective_from  DATE,
    inactive        BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by      TEXT,
    updated_at      TIMESTAMPTZ,
    updated_by      TEXT
);
CREATE INDEX idx_employee_salaries_tenant_user ON employee_salaries(tenant_id, user_id);

CREATE TABLE IF NOT EXISTS employee_salary_structures (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL,
    user_id         TEXT NOT NULL,
    template_id     UUID NOT NULL,
    fy_id           UUID NOT NULL,
    item_type       VARCHAR(20) NOT NULL,
    code            VARCHAR(50) NOT NULL,
    name            VARCHAR(255) NOT NULL,
    amount          NUMERIC(12,2) NOT NULL DEFAULT 0,
    sort_order      INT DEFAULT 0,
    inactive        BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by      TEXT,
    updated_at      TIMESTAMPTZ,
    updated_by      TEXT
);
CREATE INDEX idx_employee_salary_structures_user ON employee_salary_structures(tenant_id, user_id, fy_id);

CREATE TABLE IF NOT EXISTS salary_slips (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id           UUID NOT NULL,
    user_id             TEXT NOT NULL,
    fy_id               UUID NOT NULL,
    template_id         UUID NOT NULL,
    month               INT NOT NULL,       -- 1-12
    year                INT NOT NULL,
    gross_salary        NUMERIC(12,2),
    total_earnings      NUMERIC(12,2),
    total_deductions    NUMERIC(12,2),
    absent_deduction    NUMERIC(12,2) DEFAULT 0,
    net_salary          NUMERIC(12,2),
    absent_days         INT DEFAULT 0,
    present_days        INT DEFAULT 0,
    total_days          INT DEFAULT 0,
    lwp_days            NUMERIC(5,1) DEFAULT 0,
    no_of_ph_weo        INT DEFAULT 0,
    is_special          BOOLEAN DEFAULT FALSE,   -- pro-rata for partial month joiners
    is_regenerated      BOOLEAN DEFAULT FALSE,
    pdf_path            TEXT,
    inactive            BOOLEAN NOT NULL DEFAULT FALSE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by          TEXT,
    updated_at          TIMESTAMPTZ,
    updated_by          TEXT
);
CREATE UNIQUE INDEX idx_salary_slips_unique ON salary_slips(tenant_id, user_id, month, year) WHERE NOT inactive;
CREATE INDEX idx_salary_slips_tenant_user   ON salary_slips(tenant_id, user_id);

CREATE TABLE IF NOT EXISTS salary_slip_items (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL,
    slip_id     UUID NOT NULL,
    item_type   VARCHAR(20) NOT NULL,
    code        VARCHAR(50) NOT NULL,
    name        VARCHAR(255) NOT NULL,
    amount      NUMERIC(12,2) NOT NULL DEFAULT 0,
    sort_order  INT DEFAULT 0,
    inactive    BOOLEAN NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_salary_slip_items_slip ON salary_slip_items(slip_id);

CREATE TABLE IF NOT EXISTS salary_slip_leaves (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL,
    slip_id         UUID NOT NULL,
    leave_type_id   UUID NOT NULL,
    leave_type_name VARCHAR(100),
    total_days      NUMERIC(5,1) DEFAULT 0,
    used_days       NUMERIC(5,1) DEFAULT 0,
    balance_days    NUMERIC(5,1) DEFAULT 0,
    inactive        BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_salary_slip_leaves_slip ON salary_slip_leaves(slip_id);

-- ============================================================
-- CELEBRATIONS
-- ============================================================

CREATE TABLE IF NOT EXISTS celebration_types (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id           UUID NOT NULL,
    name                VARCHAR(100) NOT NULL,  -- "Birthday","Work Anniversary","Marriage Anniversary"
    is_yearly           BOOLEAN NOT NULL DEFAULT TRUE,
    is_user_celebration BOOLEAN NOT NULL DEFAULT TRUE,  -- FALSE = company event
    inactive            BOOLEAN NOT NULL DEFAULT FALSE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by          TEXT
);

CREATE TABLE IF NOT EXISTS celebrations (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id           UUID NOT NULL,
    branch_id           UUID,
    user_id             TEXT,                   -- NULL for company events
    celebration_type_id UUID NOT NULL,
    celebration_date    DATE,
    custom_title        VARCHAR(255),           -- for company events
    description         TEXT,
    inactive            BOOLEAN NOT NULL DEFAULT FALSE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by          TEXT,
    updated_at          TIMESTAMPTZ,
    updated_by          TEXT
);
CREATE INDEX idx_celebrations_tenant_id ON celebrations(tenant_id);
CREATE UNIQUE INDEX idx_celebrations_user_type ON celebrations(tenant_id, user_id, celebration_type_id) WHERE NOT inactive AND user_id IS NOT NULL;

-- ============================================================
-- NOTIFICATIONS
-- ============================================================

CREATE TABLE IF NOT EXISTS notification_masters (
    id                    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id             UUID NOT NULL,
    code                  VARCHAR(100) NOT NULL,
    name                  VARCHAR(255),
    description           TEXT,
    is_in_app_enabled     BOOLEAN NOT NULL DEFAULT TRUE,
    is_email_enabled      BOOLEAN NOT NULL DEFAULT FALSE,
    is_push_enabled       BOOLEAN NOT NULL DEFAULT TRUE,
    inactive              BOOLEAN NOT NULL DEFAULT FALSE,
    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by            TEXT
);
CREATE UNIQUE INDEX idx_notification_masters_code ON notification_masters(tenant_id, code) WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS notification_preferences (
    id                     UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id              UUID NOT NULL,
    user_id                TEXT NOT NULL,
    notification_master_id UUID NOT NULL,
    is_in_app_enabled      BOOLEAN NOT NULL DEFAULT TRUE,
    is_email_enabled       BOOLEAN NOT NULL DEFAULT FALSE,
    inactive               BOOLEAN NOT NULL DEFAULT FALSE,
    created_at             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at             TIMESTAMPTZ
);
CREATE UNIQUE INDEX idx_notification_prefs_user ON notification_preferences(tenant_id, user_id, notification_master_id);

CREATE TABLE IF NOT EXISTS notification_inbox (
    id                     UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id              UUID NOT NULL,
    user_id                TEXT NOT NULL,
    notification_master_id UUID NOT NULL,
    title                  TEXT NOT NULL,
    message                TEXT NOT NULL,
    reference_table        VARCHAR(100),
    reference_id           UUID,
    is_read                BOOLEAN NOT NULL DEFAULT FALSE,
    read_date              TIMESTAMPTZ,
    inactive               BOOLEAN NOT NULL DEFAULT FALSE,
    created_at             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by             TEXT,
    updated_at             TIMESTAMPTZ,
    updated_by             TEXT
);
CREATE INDEX idx_notification_inbox_user ON notification_inbox(tenant_id, user_id, is_read);

CREATE TABLE IF NOT EXISTS notification_logs (
    id                     UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id              UUID NOT NULL,
    notification_master_id UUID,
    user_id                TEXT,
    channel                VARCHAR(20) NOT NULL,   -- 'Push' | 'Email'
    target_address         TEXT NOT NULL,
    subject                TEXT,
    body                   TEXT,
    status                 VARCHAR(20) NOT NULL DEFAULT 'Pending',
    sent_date              TIMESTAMPTZ,
    error_message          TEXT,
    external_reference_id  TEXT,
    created_at             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by             TEXT
);

CREATE TABLE IF NOT EXISTS device_tokens (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id     UUID NOT NULL,
    user_id       TEXT NOT NULL,
    device_token  TEXT NOT NULL,
    device_type   VARCHAR(50),   -- 'android' | 'ios' | 'web'
    device_id     VARCHAR(255),
    inactive      BOOLEAN NOT NULL DEFAULT FALSE,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by    TEXT,
    updated_at    TIMESTAMPTZ,
    updated_by    TEXT
);
CREATE INDEX idx_device_tokens_user ON device_tokens(tenant_id, user_id) WHERE NOT inactive;

-- ============================================================
-- ONBOARDING
-- ============================================================

CREATE TABLE IF NOT EXISTS job_positions (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id           UUID NOT NULL,
    code                VARCHAR(50),
    title               VARCHAR(255) NOT NULL,
    level               VARCHAR(100),
    category            VARCHAR(100),
    description         TEXT,
    department_id       UUID,
    employment_type_id  UUID,
    work_mode           VARCHAR(50),       -- 'Remote' | 'OnSite' | 'Hybrid'
    total_position      INT DEFAULT 1,
    budgeted_cost       NUMERIC(12,2),
    inactive            BOOLEAN NOT NULL DEFAULT FALSE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by          TEXT,
    updated_at          TIMESTAMPTZ,
    updated_by          TEXT
);

CREATE TABLE IF NOT EXISTS job_position_locations (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL,
    job_position_id UUID NOT NULL,
    location        VARCHAR(255),
    city            VARCHAR(100),
    state           VARCHAR(100),
    country         VARCHAR(100),
    is_remote       BOOLEAN DEFAULT FALSE,
    inactive        BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS job_requisitions (
    id                    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id             UUID NOT NULL,
    job_position_id       UUID NOT NULL,
    code                  VARCHAR(50),
    title                 VARCHAR(255) NOT NULL,
    level                 VARCHAR(100),
    category              VARCHAR(100),
    department_id         UUID,
    employment_type_id    UUID,
    description           TEXT,
    work_mode             VARCHAR(50),
    total_openings        INT NOT NULL DEFAULT 1,
    reason_for_hire       TEXT,
    min_salary            NUMERIC(12,2),
    max_salary            NUMERIC(12,2),
    currency              VARCHAR(10) DEFAULT 'INR',
    target_hire_date      DATE,
    expected_closure_date DATE,
    requested_by          TEXT NOT NULL,
    requested_date        DATE,
    is_approved           BOOLEAN DEFAULT FALSE,
    approved_by           TEXT,
    approved_date         TIMESTAMPTZ,
    priority              VARCHAR(50),    -- 'Low' | 'Medium' | 'High' | 'Critical'
    status                VARCHAR(50) NOT NULL DEFAULT 'Draft',
    notes                 TEXT,
    inactive              BOOLEAN NOT NULL DEFAULT FALSE,
    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by            TEXT,
    updated_at            TIMESTAMPTZ,
    updated_by            TEXT
);

CREATE TABLE IF NOT EXISTS job_postings (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id           UUID NOT NULL,
    job_requisition_id  UUID,
    code                VARCHAR(50),
    title               VARCHAR(255),
    job_summary         TEXT,
    description         TEXT,
    job_category        VARCHAR(100),
    department_id       UUID,
    industry            VARCHAR(100),
    employment_type_id  UUID,
    work_mode           VARCHAR(50),
    role_type           VARCHAR(100),
    min_experience      NUMERIC(5,1),
    max_experience      NUMERIC(5,1),
    min_salary          NUMERIC(12,2),
    max_salary          NUMERIC(12,2),
    salary_currency     VARCHAR(10),
    salary_period       VARCHAR(30),
    is_salary_visible   BOOLEAN DEFAULT FALSE,
    job_status          VARCHAR(50),
    publish_date        DATE,
    expiry_date         DATE,
    is_published        BOOLEAN DEFAULT FALSE,
    inactive            BOOLEAN NOT NULL DEFAULT FALSE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by          TEXT,
    updated_at          TIMESTAMPTZ,
    updated_by          TEXT
);

CREATE TABLE IF NOT EXISTS candidates (
    id                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id            UUID NOT NULL,
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
    notice_period        INT,               -- days
    current_location     VARCHAR(255),
    preferred_location   VARCHAR(255),
    source               VARCHAR(100),
    resume_url           TEXT,
    inactive             BOOLEAN NOT NULL DEFAULT FALSE,
    created_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by           TEXT,
    updated_at           TIMESTAMPTZ,
    updated_by           TEXT
);
CREATE INDEX idx_candidates_tenant ON candidates(tenant_id);

CREATE TABLE IF NOT EXISTS candidate_applications (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL,
    candidate_id    UUID,
    job_posting_id  UUID,
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
    created_by      TEXT,
    updated_at      TIMESTAMPTZ,
    updated_by      TEXT
);
CREATE INDEX idx_candidate_applications_tenant ON candidate_applications(tenant_id);

CREATE TABLE IF NOT EXISTS interview_rounds (
    id                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id            UUID NOT NULL,
    application_id       UUID NOT NULL,
    round_name           VARCHAR(100),
    round_number         INT,
    scheduled_date       TIMESTAMPTZ,
    duration_minutes     INT,
    interviewer_user_id  TEXT,
    mode                 VARCHAR(50),       -- 'Online' | 'InPerson' | 'Phone'
    meeting_link         TEXT,
    location             VARCHAR(255),
    status               VARCHAR(50),       -- 'Scheduled' | 'Completed' | 'Cancelled'
    remarks              TEXT,
    inactive             BOOLEAN NOT NULL DEFAULT FALSE,
    created_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by           TEXT,
    updated_at           TIMESTAMPTZ,
    updated_by           TEXT
);

CREATE TABLE IF NOT EXISTS offer_letters (
    id                        UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id                 UUID NOT NULL,
    application_id            UUID NOT NULL,
    candidate_id              UUID,
    offered_ctc               NUMERIC(12,2),
    currency                  VARCHAR(10) DEFAULT 'INR',
    salary_breakdown          JSONB,
    joining_date              DATE,
    valid_until_date          DATE,
    status                    VARCHAR(50) DEFAULT 'Generated',
    offer_letter_url          TEXT,
    candidate_reaction_date   TIMESTAMPTZ,
    candidate_rejection_reason TEXT,
    version                   INT DEFAULT 1,
    is_latest                 BOOLEAN DEFAULT FALSE,
    inactive                  BOOLEAN NOT NULL DEFAULT FALSE,
    created_at                TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by                TEXT,
    updated_at                TIMESTAMPTZ,
    updated_by                TEXT
);

CREATE TABLE IF NOT EXISTS onboarding_workflows (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL,
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    inactive    BOOLEAN NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by  TEXT
);

CREATE TABLE IF NOT EXISTS onboarding_tasks (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL,
    workflow_id     UUID NOT NULL,
    title           VARCHAR(255) NOT NULL,
    description     TEXT,
    due_days        INT DEFAULT 0,          -- days after joining
    is_required     BOOLEAN DEFAULT TRUE,
    sort_order      INT DEFAULT 0,
    inactive        BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by      TEXT
);

CREATE TABLE IF NOT EXISTS candidate_onboardings (
    id                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id            UUID NOT NULL,
    candidate_id         UUID NOT NULL,
    workflow_id          UUID NOT NULL,
    onboarding_status    VARCHAR(50) DEFAULT 'Pending',
    progress_percentage  INT DEFAULT 0,
    inactive             BOOLEAN NOT NULL DEFAULT FALSE,
    created_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by           TEXT,
    updated_at           TIMESTAMPTZ,
    updated_by           TEXT
);

CREATE TABLE IF NOT EXISTS candidate_onboarding_tasks (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id               UUID NOT NULL,
    candidate_onboarding_id UUID NOT NULL,
    onboarding_task_id      UUID NOT NULL,
    status                  VARCHAR(50) DEFAULT 'Pending',
    completed_at            TIMESTAMPTZ,
    remarks                 TEXT,
    inactive                BOOLEAN NOT NULL DEFAULT FALSE,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ
);

-- ============================================================
-- AUTH SUPPORT
-- ============================================================

CREATE TABLE IF NOT EXISTS user_otps (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id       TEXT NOT NULL,
    otp           VARCHAR(10) NOT NULL,
    otp_for       VARCHAR(50) NOT NULL,    -- 'passwordreset' | 'login'
    mobile        VARCHAR(20),
    is_used       BOOLEAN DEFAULT FALSE,
    expires_at    TIMESTAMPTZ NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_user_otps_user ON user_otps(user_id);
```

---

## 4. GO MODULE STRUCTURE {#module-structure}

```
hrms/
├── cmd/
│   └── server/
│       └── main.go              -- wire all modules, start server
├── internal/
│   └── app/
│       └── app.go               -- Spur App, RegisterModule() calls
├── modules/
│   ├── tenant/                  -- company setup, branches, departments, FY, working hours
│   │   ├── module.go
│   │   ├── core/
│   │   │   ├── domain/          -- TenantBranding, Department, Branch, Designation structs
│   │   │   ├── ports/           -- ITenantService interface
│   │   │   └── services/        -- TenantService implementation
│   │   ├── adapters/
│   │   │   ├── postgres/        -- sqlc generated queries
│   │   │   └── http/            -- HTTP handlers + routes
│   │   └── sql/
│   │       ├── migrations/
│   │       └── queries/
│   │
│   ├── employee/
│   │   ├── module.go
│   │   ├── core/
│   │   ├── adapters/
│   │   └── sql/
│   │
│   ├── leave/
│   │   ├── module.go
│   │   ├── core/
│   │   ├── adapters/
│   │   └── sql/
│   │
│   ├── attendance/
│   │   ├── module.go
│   │   ├── core/
│   │   ├── adapters/
│   │   └── sql/
│   │
│   ├── salary/
│   │   ├── module.go
│   │   ├── core/
│   │   │   └── services/        -- SalaryService, PdfService
│   │   ├── adapters/
│   │   └── sql/
│   │
│   ├── dashboard/
│   │   ├── module.go
│   │   ├── core/
│   │   └── adapters/
│   │
│   ├── celebration/
│   │   ├── module.go
│   │   ├── core/
│   │   └── adapters/
│   │
│   ├── notification/
│   │   ├── module.go
│   │   ├── core/
│   │   └── adapters/
│   │
│   └── onboarding/
│       ├── module.go
│       ├── core/
│       ├── adapters/
│       └── sql/
│
├── pkg/
│   ├── pbkdf2/                  -- legacy password verifier
│   ├── pdf/                     -- salary slip PDF (gofpdf)
│   ├── storage/                 -- AWS S3
│   ├── email/                   -- SendGrid
│   ├── sms/                     -- MSG91
│   ├── push/                    -- Firebase FCM
│   └── numberwords/             -- amount to words
│
└── sqlc.yaml                    -- sqlc config per module
```

### sqlc.yaml per module:
```yaml
version: "2"
sql:
  - engine: "postgresql"
    queries: "modules/leave/sql/queries/"
    schema: "modules/leave/sql/migrations/"
    gen:
      go:
        package: "leavedb"
        out: "modules/leave/adapters/postgres"
        sql_package: "pgx/v5"
        emit_interface: true
        emit_json_tags: true
```

---

## 5. MODULE: AUTH / TENANT PROVISIONING {#module-auth}

> spur-identity handles users/roles/JWT/OAuth2. The tenant module handles company provisioning only.

### CreateTenant Flow

```
1. Validate email & mobile not already registered (global unique in auth.users)
2. Generate slug: lowercase(companyname) → replace spaces → strip special chars
3. Loop slug uniqueness: slug, slug-1, slug-2...
4. Generate 8-char code: random uppercase+digits until unique
5. INSERT INTO tenants
6. Create Identity user via spur-identity → assign role "tenant"
7. Seed tenant defaults (HrmsDefaults.json + TenantDefaults.json)
8. Send registration email via SendGrid
```

```sql
-- Check email uniqueness globally
SELECT COUNT(*) FROM auth.users WHERE email = $1 AND NOT inactive;

-- Insert tenant
INSERT INTO public.tenants (id, name, slug, code, email, mobile, status, created_at)
VALUES ($1, $2, $3, $4, $5, $6, 'active', NOW());

-- Insert tenant branding
INSERT INTO tenant_brandings (id, tenant_id, company_name, domain, primary_color, ...)
VALUES ($1, $2, $3, $4, $5, ...);
```

### CreateEmployee Flow

```
1. Check employee_code unique per tenant
2. Check email globally unique in auth.users
3. Check mobile globally unique
4. Call spur-identity CreateUser (email, mobile, password)
5. Assign role (employee/manager/hr)
6. INSERT INTO employees
7. IF statutory → INSERT INTO employee_statutory
8. IF bank → INSERT INTO employee_banks
```

### Password Migration: Legacy PBKDF2

```go
// auth.users schema additions (run once)
// ALTER TABLE auth.users ADD COLUMN legacy_password_hash TEXT;
// ALTER TABLE auth.users ADD COLUMN password_algo VARCHAR(20) DEFAULT 'bcrypt';

// ASP.NET Core Identity v3 PBKDF2 format:
// byte[0]    = 0x01 (format marker)
// byte[1..4] = PRF (HMACSHA256 = 1, big-endian)
// byte[5..8] = iteration count (10000, big-endian)
// byte[9..12]= salt length (16)
// byte[13..28] = 16-byte salt
// byte[29..60] = 32-byte subkey

func VerifyPBKDF2(storedBase64, password string) bool {
    data, _ := base64.StdEncoding.DecodeString(storedBase64)
    iterCount := int(binary.BigEndian.Uint32(data[5:9]))
    salt := data[13:29]
    storedSubkey := data[29:]
    computed := pbkdf2.Key([]byte(password), salt, iterCount, 32, sha256.New)
    return subtle.ConstantTimeCompare(computed, storedSubkey) == 1
}

// On successful PBKDF2 login → rehash to bcrypt cost 14, clear legacy hash
```

---

## 6. MODULE: EMPLOYEE {#module-employee}

### GetEmployeeList

```sql
SELECT
    e.id, e.user_id, e.employee_code,
    e.firstname, e.lastname, e.email, e.mobile,
    e.joining_date, e.role, e.grade,
    d.name AS department_name, b.branch_name, des.name AS designation_name,
    e.profile_photo_path, e.inactive
FROM employees e
LEFT JOIN departments d ON d.id = e.department_id AND NOT d.inactive
LEFT JOIN branches b ON b.id = e.branch_id AND NOT b.inactive
LEFT JOIN designations des ON des.id = e.designation_id AND NOT des.inactive
WHERE e.tenant_id = $1 AND NOT e.inactive
ORDER BY e.firstname;
```

### GetEmployeeById

```sql
-- Employee + joins
SELECT e.*, d.name as dept_name, b.branch_name, des.name as designation_name
FROM employees e
LEFT JOIN departments d ON d.id = e.department_id
LEFT JOIN branches b ON b.id = e.branch_id
LEFT JOIN designations des ON des.id = e.designation_id
WHERE e.tenant_id = $1 AND e.user_id = $2 AND NOT e.inactive;

-- Statutory
SELECT * FROM employee_statutory WHERE tenant_id=$1 AND user_id=$2 AND NOT inactive LIMIT 1;

-- Bank
SELECT * FROM employee_banks WHERE tenant_id=$1 AND user_id=$2 AND NOT inactive ORDER BY is_primary DESC;
```

### UpdateEmployee

```sql
-- UPSERT bank
INSERT INTO employee_banks (id, tenant_id, user_id, bank_name, account_number, ifsc_code, ...)
VALUES ($1, $2, $3, $4, $5, $6, ...)
ON CONFLICT (tenant_id, user_id) WHERE NOT inactive
DO UPDATE SET bank_name=$4, account_number=$5, ifsc_code=$6, updated_at=NOW(), updated_by=$7;

-- UPSERT statutory
INSERT INTO employee_statutory (id, tenant_id, user_id, pf_no, uan_no, esic_no, pan, ...)
VALUES (...)
ON CONFLICT (tenant_id, user_id) WHERE NOT inactive
DO UPDATE SET pf_no=EXCLUDED.pf_no, uan_no=EXCLUDED.uan_no, updated_at=NOW();
```

### DeactivateEmployee

```sql
UPDATE employees SET inactive=TRUE, updated_at=NOW(), updated_by=$3 WHERE tenant_id=$1 AND user_id=$2;
UPDATE auth.users SET inactive=TRUE WHERE id=$1;   -- via spur-identity
```

---

## 7. MODULE: LEAVE {#module-leave}

### ApplyLeave Flow (most complex)

```
1. Parse/normalize dates: multi-day half-days normalized
2. Sandwich check → may add gap days to leave count
3. Validate dates within active FY
4. Check reporting manager exists; if NULL → notify HR, still allow
5. Load leave balance
6. Overlap detection with DayPart bitmask (Morning=1, Afternoon=2, Full=3)
7. Single-day rule: firsthalf+secondhalf = full (OK); duplicate halves = overlap error
8. Calculate total_days (Full=1.0, Half=0.5, Sandwich gap days added)
9. Check balance: balance_days - new_days >= 0 (else reject)
10. INSERT INTO leaves (status='pending')
11. INSERT INTO leave_approvals (approver=reporting_manager)
12. UPDATE leave_balances SET pending_days += new_days
13. Notify manager: NotifLeaveApplied
```

### IsSandwich

```go
func IsSandwich(startDate, endDate time.Time, holidays []time.Time, weekoffs []time.Weekday) (bool, float64) {
    if startDate.Equal(endDate) { return false, 0 }
    gapDays := float64(0)
    allNonWorking := true
    for d := startDate.AddDate(0,0,1); d.Before(endDate); d = d.AddDate(0,0,1) {
        isHoliday := slices.Contains(holidays, d)
        isWeekoff := slices.Contains(weekoffs, d.Weekday())
        if !isHoliday && !isWeekoff { allNonWorking = false; break }
        gapDays++
    }
    if allNonWorking && gapDays > 0 { return true, gapDays }
    return false, 0
}
```

### ApproveLeave Flow

```
1. CL consecutive limit: IF leave_type.is_consecutive_limit:
   count existing approved CL days in FY; IF count + new > limit → reject

2. EarnLeave auto-conversion (shortcode == "earnleave"):
   a. Scan BACKWARD from start_date for adjacent approved leaves of different type
   b. Scan FORWARD from end_date for adjacent approved leaves of different type
   c. IF adjacent found:
      - Validate combined balance covers total days
      - Refund old type: used_days -= old.days; INSERT credit ledger
      - Reassign old leave_type_id to EL
      - Deduct EL: used_days += combined; INSERT debit ledger

3. INSERT leave_approvals (status='approved')
4. UPDATE leave_balances: used_days += days, pending_days -= days
5. INSERT leave_ledger (transaction_type='debit')
6. UPDATE leaves SET status='approved'
7. Notify employee: NotifLeaveApproved
```

### RejectLeave

```sql
UPDATE leave_approvals SET status='rejected', remarks=$3, action_date=NOW() WHERE leave_id=$1 AND approver_id=$2;
UPDATE leave_balances SET pending_days = pending_days - (SELECT days FROM leaves WHERE id=$1)
  WHERE tenant_id=$1 AND user_id=(SELECT user_id FROM leaves WHERE id=$1)
    AND leave_type_id=(SELECT leave_type_id FROM leaves WHERE id=$1)
    AND fy_id=(SELECT fy_id FROM leaves WHERE id=$1);
UPDATE leaves SET status='rejected' WHERE id=$1;
```

### GetLeaveList (Manager)

```sql
SELECT l.id, l.user_id, l.start_date, l.end_date, l.days, l.status, l.applied_date,
       lt.name AS leave_type_name, lt.shortcode,
       e.firstname, e.lastname, e.employee_code, e.profile_photo_path
FROM leaves l
JOIN leave_types lt ON lt.id = l.leave_type_id AND NOT lt.inactive
JOIN employees e ON e.user_id = l.user_id AND NOT e.inactive
WHERE l.tenant_id=$1 AND NOT l.inactive AND l.fy_id=$2
  AND ($3::text IS NULL OR e.reporting_manager_id = $3)
ORDER BY l.applied_date DESC;
```

### DistributeInto12Months

```go
func DistributeInto12Months(totalDays float64) [12]int {
    base := int(totalDays) / 12
    remainder := int(totalDays) % 12
    var months [12]int
    for i := range months { months[i] = base }
    for i := 0; i < remainder; i++ { months[i]++ }
    return months
}
```

---

## 8. MODULE: ATTENDANCE {#module-attendance}

### Attendance Status Priority

```
1. Before joining date → skip (not applicable)
2. Approved leave exists for date → "leave"
3. Attendance record has status already set → use it
4. Check-in record exists → "present"
5. Holiday in holidays table → "holiday"
6. Week-off per working_hours → "weekoff"
7. Date is past → "absent"
8. Date is today, no check-in → "" (empty)
```

### GetWorkingHour Hierarchy

```go
// 1. User-specific hours
// 2. Branch-level hours (employee's branch_id)
// 3. Tenant-level hours (no branch, no user)
// 4. Default: Mon-Fri 9-6, Sat-Sun off
```

### Working Hours SQL

```sql
-- Day order for display: Mon=1, Tue=2, Wed=3, Thu=4, Fri=5, Sat=6, Sun=7
SELECT *, CASE day_of_week
    WHEN 'Monday' THEN 1 WHEN 'Tuesday' THEN 2 WHEN 'Wednesday' THEN 3
    WHEN 'Thursday' THEN 4 WHEN 'Friday' THEN 5 WHEN 'Saturday' THEN 6
    WHEN 'Sunday' THEN 7 END AS sort_order
FROM working_hours
WHERE tenant_id=$1 AND branch_id IS NULL AND user_id IS NULL AND NOT inactive
ORDER BY sort_order;
```

### CopyWorkingHoursToBranch

```sql
UPDATE working_hours SET inactive=TRUE WHERE tenant_id=$1 AND branch_id=$2;
INSERT INTO working_hours (id, tenant_id, branch_id, day_of_week, is_working_day, start_time, end_time, created_at)
SELECT gen_random_uuid(), $1, $2, day_of_week, is_working_day, start_time, end_time, NOW()
FROM working_hours WHERE tenant_id=$1 AND branch_id IS NULL AND user_id IS NULL AND NOT inactive;
```

---

## 9. MODULE: SALARY {#module-salary}

### CalculateUserSalary

```go
func CalculateUserSalary(grossSalary float64, structure []SalaryItem, presentDays, absentDays, daysInMonth int, isSpecial bool) SalaryResult {
    var earnings, deductions float64
    for _, item := range structure {
        if item.ItemType == "earning" { earnings += item.Amount } else { deductions += item.Amount }
    }
    absentDeduction := float64(0)
    if absentDays > 0 {
        absentDeduction = (grossSalary / float64(daysInMonth)) * float64(absentDays)
        deductions += absentDeduction
    }
    if isSpecial && presentDays > 0 {
        earnings = (earnings / float64(daysInMonth)) * float64(presentDays)
    }
    return SalaryResult{
        TotalEarnings: earnings, TotalDeductions: deductions,
        AbsentDeduction: absentDeduction, NetSalary: earnings - deductions,
    }
}
```

### Template Activation

```sql
-- Deactivate all others first
UPDATE salary_templates SET is_active=FALSE WHERE tenant_id=$1 AND fy_id=$2;
UPDATE salary_templates SET is_active=TRUE, updated_at=NOW() WHERE id=$3;
```

### Salary Slip Generation (Upsert)

```go
// 1. Check existing slip → if exists and !isRegenerated → error
// 2. Load active FY, active template, employee salary structure
// 3. Calculate attendance counts for month/year
// 4. CalculateUserSalary()
// 5. If absentDeduction > 0 → add LWP item dynamically
// 6. Snapshot leave balances
// 7. UPSERT salary_slips
// 8. DELETE + INSERT salary_slip_items
// 9. DELETE + INSERT salary_slip_leaves
```

---

## 10. MODULE: PDF SALARY SLIP {#module-pdf}

> Port iTextSharp → use `github.com/jung-kurt/gofpdf`

### Page Setup
```
Size: A4 (595.28 x 841.89 pt)
Margins: Left=40, Right=40, Top=35, Bottom=40
Font: Helvetica
```

### Section 1: Header (3-column, widths 20%/50%/30%)
```
Col 1: Company logo (90x70 pt max)
Col 2: Company name (14pt bold) + address (9pt)
Col 3: "Payslip For the Month" (right-aligned) + "Month Year" (14pt bold)
```

### Section 2: EMPLOYEE SUMMARY (7-col grid)
```
Widths: [18, 4, 26, 15, 18, 4, 15]
Pattern: label | ":" | value | spacer | label | ":" | value
Rows:
  Employee Code / Date of Joining
  Employee Name / UAN No.
  Branch / Department
  Designation / Grade
  PF No. / ESIC No.
  PAN Card / "Basic Salary"
```

### Section 3: LEAVE DETAILS (rounded border)
```
8-col inner table, widths [18,7,18,7,18,7,18,7]
Row 1: Total Days / Days Present / PH+WEO / LWP Absent
Row 2+: Per leave type — name : utilized | balance
```

### Section 4: EARNINGS / DEDUCTIONS (side-by-side tables)
```
Both tables equalized to same row count
Left: EARNINGS header → items → "Gross Earnings: {total}"
Right: DEDUCTIONS header → items → "Total Deductions: {total}"
```

### Section 5: NET PAYABLE (rounded border box)
```
NET PAYABLE: ₹{net_salary N0}
Amount In Words: {NumberToWords(net)} Only
```

### Section 6: PAYMENT DETAILS (4-col, bottom-border-only cells)
```
Headers: Mode of Payment | Employee Bank | Account No | Amount
Data:    Bank Transfer   | {bank_name}   | {account}  | {net N0}
```

### Section 7: Footer
```
Horizontal rule
"This is a system-generated payslip and does not require a signature."
"Powered by Setika" (gray, right-aligned)
```

### Colors
```go
NavyBlue    = RGB{44, 67, 91}
Primary     = parse(tenant_brandings.primary_color)  // hex → RGB
LightGrayBg = RGB{245, 245, 245}
BorderColor = RGB{211, 211, 211}
```

---

## 11. MODULE: DASHBOARD {#module-dashboard}

### GetDashboardData (Employee)

```go
type DashboardDto struct {
    Profile          UserProfileDto
    Celebrations     []CelebrationDto       // top 5, sorted by DaysLeft ascending
    AttendanceSummary AttendanceSummaryDto  // last 7 days
    LeaveSummary     LeaveSummaryDto        // requested/approved/rejected/pending counts
}
```

```sql
-- Profile
SELECT e.user_id, e.firstname, e.lastname, e.profile_photo_path,
       des.name AS designation, e.city, e.state, e.experience_year, e.experience_month
FROM employees e LEFT JOIN designations des ON des.id = e.designation_id
WHERE e.tenant_id=$1 AND e.user_id=$2 AND NOT e.inactive;

-- Leave summary (current FY)
SELECT
    COUNT(*) AS requested,
    COUNT(*) FILTER (WHERE status='approved') AS approved,
    COUNT(*) FILTER (WHERE status='rejected') AS rejected
FROM leaves WHERE tenant_id=$1 AND user_id=$2 AND NOT inactive AND fy_id=$3;

-- Last 7 days attendance
SELECT date, type, status FROM attendances
WHERE tenant_id=$1 AND user_id=$2 AND date >= $3 AND NOT inactive ORDER BY time;
-- For each day: if attendance.status empty → check approved leaves
```

### GetEmployeeDashboard (Full)

**Additional data:**
```go
TodayStatus:      first checkin time + last checkout time for today
LeaveList:        pending+approved future leaves in current FY (enddate >= today)
QuickTools:       role-based action list
AttendanceSummary: last 7 days with worked hours per day
TotalWork:        "X.Xh/Y.Yh" worked vs expected
DailyUpdates:     clock-in reminder + next upcoming holiday
WorkModeList:     available work modes
AttendanceConfig: attendance settings for user
```

**Quick tools (always):** leaverequest, attendancerecord, forgottopunch, celebration, holidays, policies
**Quick tools (manager only):** + leaveapproval

**Total work calculation:**
```go
for each day in last 7:
    wh = GetWorkingHour(tenantID, userID, dayOfWeek)
    if wh.IsWorkingDay: expectedHours += wh.EndTime - wh.StartTime
    clockIn = first checkin for day
    if clockIn != nil:
        clockOut = last checkout ?? now
        totalWorked += clockOut - clockIn
result = fmt.Sprintf("%.1fh/%.1fh", totalWorked.Hours(), totalExpected.Hours())
```

---

## 12. MODULE: CELEBRATIONS {#module-celebrations}

### CreateCelebration

```sql
-- Uniqueness check
SELECT COUNT(*) FROM celebrations
WHERE tenant_id=$1 AND user_id=$2 AND celebration_type_id=$3 AND NOT inactive;

INSERT INTO celebrations (id, tenant_id, branch_id, user_id, celebration_type_id,
    celebration_date, custom_title, description, created_at, created_by)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,NOW(),$9);
```

### GetCelebrationList

```sql
SELECT c.id, c.user_id, c.celebration_date, c.custom_title,
       ct.name AS celebration_name,
       e.firstname || ' ' || e.lastname AS user_name,
       b.branch_name
FROM celebrations c
JOIN celebration_types ct ON ct.id = c.celebration_type_id AND NOT ct.inactive
JOIN employees e ON e.user_id = c.user_id AND NOT e.inactive
LEFT JOIN branches b ON b.id = e.branch_id AND NOT b.inactive
WHERE c.tenant_id=$1 AND NOT c.inactive
ORDER BY CASE WHEN c.celebration_date >= CURRENT_DATE THEN 0 ELSE 1 END, c.celebration_date;
```

### SendCelebrationMail (daily job)

```go
// Query: celebrations where today matches:
//   Birthday/Work Anniversary/Marriage Anniversary → month+day == today.month+day
//   Other types → date == today exactly
//
// Filter out inactive employees
//
// For is_user_celebration=true:
//   - Send personal email to celebrating employee (birthday/anniversary template)
//   - Send group notification email to all OTHER employees (exclude celebrant)
//   - Send Firebase push to all employees (different messages for celebrant vs others)
//
// For is_user_celebration=false (company event):
//   - Send single email to ALL employees
//   - Send Firebase push to ALL employees
```

### NextOccurrence

```go
func NextOccurrence(original, fromDate time.Time, isYearly bool) time.Time {
    if !isYearly { return original }
    maxDay := daysInMonth(fromDate.Year(), original.Month())
    day := min(original.Day(), maxDay)
    candidate := time.Date(fromDate.Year(), original.Month(), day, 0,0,0,0, time.UTC)
    if candidate.Before(fromDate) { candidate = candidate.AddDate(1,0,0) }
    return candidate
}
```

---

## 13. MODULE: NOTIFICATIONS {#module-notifications}

### SendNotification

```go
func SendNotification(tenantID, senderID, recipientID, code, title, message, refTable string, refID uuid.UUID) {
    // 1. Load notification_masters by code+tenant
    // 2. Load notification_preferences for user (opt-out check)
    // 3. Effective flags: pref?.IsInApp ?? config.IsInApp
    // 4. If sendInApp: INSERT notification_inbox
    // 5. Load device_tokens for user
    // 6. For each token: Firebase.Send() → INSERT notification_logs (Sent/Failed)
}
```

### Bulk SendNotification

```go
// For celebrations: accepts []string recipientUserIDs
// Loop and call SendNotification per user
```

### GetUserNotifications (paged)

```sql
SELECT ni.id, nm.code, ni.title, ni.message, ni.reference_table, ni.reference_id, ni.is_read, ni.created_at
FROM notification_inbox ni
JOIN notification_masters nm ON nm.id = ni.notification_master_id AND NOT nm.inactive
WHERE ni.tenant_id=$1 AND ni.user_id=$2 AND NOT ni.inactive
ORDER BY ni.created_at DESC
LIMIT $3 OFFSET $4;
-- Enrich leave notifications with leave_types + employees data
```

### GetUserInbox (policy/celebration/general only)

```sql
-- Same as above but filter nm.code IN ('companypolicy','generalnotification','usercelebration')
-- Mark fetched items as read immediately after fetch
```

### SaveDeviceToken (upsert by device_id)

```go
// If same device_id + same token exists → no-op
// Else: deactivate all tokens for that device_id, insert new
```

---

## 14. MODULE: ONBOARDING {#module-onboarding}

### Entity Relationships

```
job_positions (1) → (N) job_position_locations
job_positions (1) → (N) job_requisitions
job_requisitions (1) → (1) job_postings
job_postings (1) → (N) candidate_applications
candidates (1) → (N) candidate_applications
candidate_applications (1) → (N) interview_rounds
candidate_applications (1) → (N) offer_letters
candidates (1) → (1) candidate_onboardings
onboarding_workflows (1) → (N) onboarding_tasks
candidate_onboardings (1) → (N) candidate_onboarding_tasks
```

### Candidate Pipeline

```
Draft Requisition → Approve → Create Job Posting → Publish
→ Candidate applies (status='New') → Screening → Interview rounds
→ Selected → Create Offer Letter (status='Generated') → Send → Accept/Decline
→ If Accepted: application.status='Hired', create candidate_onboarding
→ Assign workflow tasks → Complete tasks → progress 0..100% → status='Completed'
```

### Create Offer Letter

```sql
UPDATE offer_letters SET is_latest=FALSE WHERE application_id=$1 AND NOT inactive;
INSERT INTO offer_letters (id, tenant_id, application_id, candidate_id, offered_ctc,
    salary_breakdown, joining_date, valid_until_date, status='Generated', is_latest=TRUE, ...)
VALUES (...);
UPDATE candidate_applications SET status='Offered' WHERE id=$1;
```

### Start Onboarding

```sql
INSERT INTO candidate_onboardings (id, tenant_id, candidate_id, workflow_id, onboarding_status='Pending', progress_percentage=0, ...);
INSERT INTO candidate_onboarding_tasks (id, tenant_id, candidate_onboarding_id, onboarding_task_id, status='Pending', ...)
SELECT gen_random_uuid(), $1, $2, ot.id, 'Pending', NOW()
FROM onboarding_tasks ot WHERE ot.workflow_id=$3 AND NOT ot.inactive;
```

### Update Progress

```go
total := countTasks(onboardingID)
completed := countCompletedTasks(onboardingID)
progress := (completed * 100) / total
status := "InProgress"
if progress == 100 { status = "Completed" }
UPDATE candidate_onboardings SET progress_percentage=$1, onboarding_status=$2 WHERE id=$3
```

---

## 15. EXTERNAL INTEGRATIONS {#integrations}

### SendGrid
```go
import "github.com/sendgrid/sendgrid-go"
// Env: SENDGRID_API_KEY, SENDGRID_FROM_EMAIL
// HTML emails only
```

### MSG91 (SMS OTP)
```
POST https://api.msg91.com/api/v5/otp
Header: authkey: MSG91_AUTH_KEY
Body: { "template_id": "...", "mobile": "91{phone}", "otp": "{otp}" }
Verify: POST https://api.msg91.com/api/v5/otp/verify
```

### Gupshup (WhatsApp)
```
Template ID: bfb8b1fd-2f54-4c4d-966c-d01db1547891
POST https://api.gupshup.io/sm/api/v1/template/msg
Header: apikey: GUPSHUP_API_KEY
```

### AWS S3
```
Key structure:
  Profile photo:  /employees/{tenantID}/{userID}/profile.jpg
  Salary slip:    /salary/{tenantID}/{userID}/{year}/{month}/slip.pdf
  Company policy: /policies/{tenantID}/{policyID}/{filename}
```

### Firebase FCM
```go
import "firebase.google.com/go/v4/messaging"
// Initialize from service account JSON (env: FIREBASE_CREDENTIALS_JSON)
// Returns message ID on success
```

---

## 16. PASSWORD MIGRATION {#password-migration}

```sql
-- Step 1: Add columns
ALTER TABLE auth.users ADD COLUMN IF NOT EXISTS legacy_password_hash TEXT;
ALTER TABLE auth.users ADD COLUMN IF NOT EXISTS password_algo VARCHAR(20) NOT NULL DEFAULT 'bcrypt';

-- Step 2: Mark existing .NET users
UPDATE auth.users
SET legacy_password_hash = password_hash,
    password_algo = 'pbkdf2',
    password_hash = ''
WHERE password_hash LIKE 'AQ%';  -- Base64(0x01...) starts with "AQ"
```

```go
// Login flow:
func Login(email, password string) (*User, error) {
    user := findByEmail(email)
    valid, needsRehash := VerifyPassword(user, password)
    if !valid { return nil, ErrInvalidCredentials }
    if needsRehash {
        hash, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
        db.Exec("UPDATE auth.users SET password_hash=$1, password_algo='bcrypt', legacy_password_hash=NULL WHERE id=$2", hash, user.ID)
    }
    return user, nil
}
```

---

## 17. AGENT GUARDRAILS {#guardrails}

### 🔴 NEVER DO

1. **Never skip `tenant_id` in queries.** Every SELECT/INSERT/UPDATE/DELETE must filter by `tenant_id = $1`. Cross-tenant leakage is a critical security bug.

2. **Never write raw SQL in service layer.** All queries go through sqlc-generated functions. No `db.Exec("SELECT ...")` in services.

3. **Never activate two salary templates simultaneously.** Always deactivate all others for `(tenant_id, fy_id)` before activating one.

4. **Never allow duplicate active leave balances.** Use UPSERT, not blind INSERT. Unique index on `(tenant_id, user_id, leave_type_id, fy_id)`.

5. **Never deduct balance on apply — only on approve.** Apply → increment `pending_days`. Approve → move to `used_days`. Reject → decrement `pending_days`.

6. **Never regenerate salary slip without `regenerate=true` flag.** Return error if slip exists and flag is false.

7. **Never store plaintext passwords.** bcrypt cost 14 minimum.

8. **Never hard-delete records.** Soft-delete only (`inactive=TRUE`).

9. **Never modify auth schema tables directly.** Use spur-identity APIs.

10. **Never return soft-deleted records.** Always add `AND NOT inactive` to queries.

---

### 🟡 ALWAYS DO

1. **Always set `created_by`/`updated_by`** from JWT context user ID.

2. **Always validate active financial year** before leave, salary, attendance ops.

3. **Always use `TenantIsolation()` middleware** — tenant_id from JWT only, never from request body.

4. **Always check overlap before inserting leave** using DayPart bitmask (both pending + approved).

5. **Always snapshot leave balances** in `salary_slip_leaves` when generating a slip.

6. **Always use UTC for timestamps.** IST conversion is frontend responsibility.

7. **Always wrap leave approval in a DB transaction** (approval + balance update + ledger + notification).

8. **Always check celebrating employee uniqueness** (one celebration per user per type per tenant).

9. **Always check reporting manager exists** before applying leave; notify HR if missing.

10. **Always validate FY dates** — leave start+end must both be within same active FY.

---

### 🔵 EDGE CASES

| Edge Case | Rule |
|---|---|
| Leave crossing FY boundary | Reject — both dates must be in same active FY |
| Half-day overlap same date | Morning+Afternoon = 1 full day (OK); two Mornings = error |
| Salary for partial month joiner | `isSpecial=true` → pro-rate earnings = (earnings/daysInMonth)*presentDays |
| LWP deduction = 0 | Do NOT add LWP line item to salary_slip_items |
| Feb 29 birthday in non-leap year | Use min(29, daysInMonth(year, Feb)) → project to Feb 28 |
| EarnLeave auto-conversion | Only when shortcode == "earnleave"; check adjacency both directions |
| CL consecutive limit | Count ALL approved CL days in current FY before approving |
| Device token rotation | Deactivate ALL old tokens for same device_id before new insert |
| Notification master missing | Silent skip (no crash); log the miss |
| Attendance today with no checkin | Return `""` not `"absent"` — absent only for past dates |
| Celebration email to celebrant | They get personal email; remove them from group notification list |
| Onboarding progress 100% | Set `onboarding_status = "Completed"` |
| Working hour fallback | User → Branch → Tenant → Default (Mon-Fri 9-6) |
| PBKDF2 header v3 format | Byte layout: [0]=0x01, [1-4]=PRF, [5-8]=iterations, [9-12]=saltLen, [13-28]=salt, [29+]=subkey |

---

### 🧪 TESTING RULES

1. Table-driven tests for: `CalculateLeaveDays`, `IsSandwich`, `CalculateUserSalary`, `DistributeInto12Months`, `NextOccurrence`, `VerifyPBKDF2`.

2. HTTP handler tests use sqlc mock querier interface (set `emit_interface: true` in sqlc.yaml).

3. Salary tests must cover: absentDays=0, absentDays>0, isSpecial=true/false, LWP added/not added.

4. Leave overlap tests must cover: full+full same day, firsthalf+secondhalf same day, two firsthalfs (error), adjacent days (no overlap), date range overlap.

5. Integration tests for attendance status must test all 6 priority levels.

6. PBKDF2 tests use real ASP.NET Identity test vectors — never mock.

---

### 🏗️ MODULE WIRING ORDER (app.go)

```go
// Wire in this order (dependency chain):
// 1. spur-identity (external — auth, users, roles, JWT)
// 2. tenant      (no custom module deps)
// 3. employee    (depends on tenant: dept/branch/designation)
// 4. leave       (depends on employee, tenant)
// 5. attendance  (depends on employee, tenant)
// 6. salary      (depends on employee, leave, attendance)
// 7. notification (depends on employee)
// 8. celebration  (depends on employee, notification)
// 9. dashboard    (depends on employee, leave, attendance, celebration)
// 10. onboarding  (depends on tenant, employee)
```

---

### 🔔 SCHEDULED JOBS

| Job | Cron | Description |
|---|---|---|
| SendCelebrationEmails | `0 8 * * *` (IST) | Today's birthdays/anniversaries: emails + push |
| LeaveBalanceAllocation | `0 0 1 * *` | Monthly allocation type: credit this month's days |
| HolidayNotification | `0 8 * * *` | Notify employees of holidays in next 3 days |
| SalarySlipReminder | Configurable | Remind HR to generate monthly slips |

---

### 📋 SEED DATA (run on new tenant creation)

```json
{
  "notification_masters": [
    {"code": "leaveapplied",        "is_in_app": true,  "is_push": true,  "is_email": false},
    {"code": "leaveapproved",       "is_in_app": true,  "is_push": true,  "is_email": false},
    {"code": "leaverejected",       "is_in_app": true,  "is_push": true,  "is_email": false},
    {"code": "companypolicy",       "is_in_app": true,  "is_push": false, "is_email": false},
    {"code": "generalnotification", "is_in_app": true,  "is_push": true,  "is_email": false},
    {"code": "usercelebration",     "is_in_app": true,  "is_push": true,  "is_email": false}
  ],
  "celebration_types": [
    {"name": "Birthday",             "is_yearly": true,  "is_user_celebration": true},
    {"name": "Work Anniversary",     "is_yearly": true,  "is_user_celebration": true},
    {"name": "Marriage Anniversary", "is_yearly": true,  "is_user_celebration": true},
    {"name": "Company Event",        "is_yearly": false, "is_user_celebration": false}
  ],
  "financial_year": {
    "name": "FY 2024-25",
    "start_date": "2024-04-01",
    "end_date": "2025-03-31",
    "is_active": true
  }
}
```

---

*End of BUSINESS_LOGIC_REFERENCE.md v2*
*Extracted from SetikaHRMS .NET Core | 10 modules | Complete SQL + Go structure + guardrails*
*Last updated: 2026-06-11*
