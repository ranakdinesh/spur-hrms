-- name: ListHolidays :many
SELECT * FROM hrms.holidays
WHERE tenant_id = $1 AND NOT inactive
ORDER BY date ASC, name ASC;

-- name: CreateHoliday :one
INSERT INTO hrms.holidays (
    tenant_id,
    branch_id,
    fy_id,
    name,
    date,
    is_optional,
    created_by,
    updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $7
)
RETURNING *;

-- name: GetHoliday :one
SELECT * FROM hrms.holidays
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: ListHolidaysByDateRange :many
SELECT * FROM hrms.holidays
WHERE tenant_id = $1 AND date BETWEEN $2 AND $3 AND NOT inactive
ORDER BY date ASC, name ASC;

-- name: ListHolidaysByFinancialYear :many
SELECT * FROM hrms.holidays
WHERE tenant_id = $1 AND fy_id = $2 AND NOT inactive
ORDER BY date ASC, name ASC;

-- name: ListHolidaysByBranch :many
SELECT * FROM hrms.holidays
WHERE tenant_id = $1 AND branch_id = $2 AND NOT inactive
ORDER BY date ASC, name ASC;

-- name: ListUpcomingHolidays :many
SELECT * FROM hrms.holidays
WHERE tenant_id = $1 AND date >= CURRENT_DATE AND NOT inactive
ORDER BY date ASC, name ASC
LIMIT $2;

-- name: UpdateHoliday :one
UPDATE hrms.holidays
SET branch_id = $3,
    fy_id = $4,
    name = $5,
    date = $6,
    is_optional = $7,
    updated_by = $8,
    updated_at = NOW()
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: SoftDeleteHoliday :exec
UPDATE hrms.holidays
SET inactive = TRUE,
    updated_by = $3,
    updated_at = NOW()
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: GetPayCycle :one
SELECT * FROM hrms.pay_cycles
WHERE tenant_id = $1 AND NOT inactive;

-- name: UpsertPayCycle :one
INSERT INTO hrms.pay_cycles (
    tenant_id,
    name,
    cycle_type,
    pay_day,
    start_day,
    end_day,
    attendance_source,
    attendance_period_type,
    attendance_cutoff_day,
    payout_timing,
    payout_offset_days,
    include_weekly_offs,
    include_holidays,
    prorate_joining_exit,
    proration_basis,
    allow_arrears,
    arrears_mode,
    allow_negative_net_pay,
    overtime_component_code,
    lwp_component_code,
    rounding_mode,
    payment_mode,
    payment_file_format,
    requires_approval,
    auto_lock_after_approval,
    payroll_lock_day,
    pf_enabled,
    pf_employee_rate,
    pf_employer_rate,
    pf_wage_ceiling,
    pf_apply_ceiling,
    esi_enabled,
    esi_employee_rate,
    esi_employer_rate,
    esi_wage_ceiling,
    professional_tax_enabled,
    tds_enabled,
    country_code,
    state_code,
    notes,
    created_by,
    updated_by
)
VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
    $11, $12, $13, $14, $15, $16, $17, $18, $19, $20,
    $21, $22, $23, $24, $25, $26, $27, $28, $29, $30,
    $31, $32, $33, $34, $35, $36, $37, $38, $39, $40,
    $41, $41
)
ON CONFLICT (tenant_id) WHERE NOT inactive
DO UPDATE SET
    name = EXCLUDED.name,
    cycle_type = EXCLUDED.cycle_type,
    pay_day = EXCLUDED.pay_day,
    start_day = EXCLUDED.start_day,
    end_day = EXCLUDED.end_day,
    attendance_source = EXCLUDED.attendance_source,
    attendance_period_type = EXCLUDED.attendance_period_type,
    attendance_cutoff_day = EXCLUDED.attendance_cutoff_day,
    payout_timing = EXCLUDED.payout_timing,
    payout_offset_days = EXCLUDED.payout_offset_days,
    include_weekly_offs = EXCLUDED.include_weekly_offs,
    include_holidays = EXCLUDED.include_holidays,
    prorate_joining_exit = EXCLUDED.prorate_joining_exit,
    proration_basis = EXCLUDED.proration_basis,
    allow_arrears = EXCLUDED.allow_arrears,
    arrears_mode = EXCLUDED.arrears_mode,
    allow_negative_net_pay = EXCLUDED.allow_negative_net_pay,
    overtime_component_code = EXCLUDED.overtime_component_code,
    lwp_component_code = EXCLUDED.lwp_component_code,
    rounding_mode = EXCLUDED.rounding_mode,
    payment_mode = EXCLUDED.payment_mode,
    payment_file_format = EXCLUDED.payment_file_format,
    requires_approval = EXCLUDED.requires_approval,
    auto_lock_after_approval = EXCLUDED.auto_lock_after_approval,
    payroll_lock_day = EXCLUDED.payroll_lock_day,
    pf_enabled = EXCLUDED.pf_enabled,
    pf_employee_rate = EXCLUDED.pf_employee_rate,
    pf_employer_rate = EXCLUDED.pf_employer_rate,
    pf_wage_ceiling = EXCLUDED.pf_wage_ceiling,
    pf_apply_ceiling = EXCLUDED.pf_apply_ceiling,
    esi_enabled = EXCLUDED.esi_enabled,
    esi_employee_rate = EXCLUDED.esi_employee_rate,
    esi_employer_rate = EXCLUDED.esi_employer_rate,
    esi_wage_ceiling = EXCLUDED.esi_wage_ceiling,
    professional_tax_enabled = EXCLUDED.professional_tax_enabled,
    tds_enabled = EXCLUDED.tds_enabled,
    country_code = EXCLUDED.country_code,
    state_code = EXCLUDED.state_code,
    notes = EXCLUDED.notes,
    updated_by = EXCLUDED.updated_by,
    updated_at = NOW()
RETURNING *;

-- name: SoftDeletePayCycle :exec
UPDATE hrms.pay_cycles
SET inactive = TRUE, updated_by = $2
WHERE tenant_id = $1 AND NOT inactive;

-- name: ListPolicyTypes :many
SELECT * FROM hrms.policy_types
WHERE (tenant_id = $1 OR is_system) AND NOT inactive
ORDER BY is_system DESC, name ASC;

-- name: CreatePolicyType :one
INSERT INTO hrms.policy_types (
    tenant_id,
    name,
    is_system,
    created_by,
    updated_by
) VALUES (
    $1, $2, $3, $4, $4
)
RETURNING *;

-- name: GetPolicyType :one
SELECT * FROM hrms.policy_types
WHERE (tenant_id = $1 OR is_system) AND id = $2 AND NOT inactive;

-- name: UpdatePolicyType :one
UPDATE hrms.policy_types
SET name = $3,
    updated_by = $4,
    updated_at = NOW()
WHERE tenant_id = $1 AND id = $2 AND NOT is_system AND NOT inactive
RETURNING *;

-- name: SoftDeletePolicyType :exec
UPDATE hrms.policy_types
SET inactive = TRUE,
    updated_by = $3,
    updated_at = NOW()
WHERE tenant_id = $1 AND id = $2 AND NOT is_system AND NOT inactive;

-- name: ListCompanyPolicies :many
SELECT * FROM hrms.company_policies
WHERE tenant_id = $1 AND NOT inactive
ORDER BY created_at DESC;

-- name: CreateCompanyPolicy :one
INSERT INTO hrms.company_policies (
    tenant_id,
    policy_type_id,
    title,
    file_path,
    description,
    created_by,
    updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $6
)
RETURNING *;

-- name: GetCompanyPolicy :one
SELECT * FROM hrms.company_policies
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: ListCompanyPoliciesByType :many
SELECT * FROM hrms.company_policies
WHERE tenant_id = $1 AND policy_type_id = $2 AND NOT inactive
ORDER BY created_at DESC;

-- name: UpdateCompanyPolicy :one
UPDATE hrms.company_policies
SET policy_type_id = $3,
    title = $4,
    file_path = $5,
    description = $6,
    updated_by = $7,
    updated_at = NOW()
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: SoftDeleteCompanyPolicy :exec
UPDATE hrms.company_policies
SET inactive = TRUE,
    updated_by = $3,
    updated_at = NOW()
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: ListTenantSubscriptions :many
SELECT * FROM hrms.tenant_subscriptions
WHERE tenant_id = $1 AND NOT inactive
ORDER BY created_at DESC;

-- name: GetTenantSubscription :one
SELECT * FROM hrms.tenant_subscriptions
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: GetCurrentTenantSubscription :one
SELECT * FROM hrms.tenant_subscriptions
WHERE tenant_id = $1 AND status IN ('trialing','active','past_due') AND NOT inactive
ORDER BY created_at DESC
LIMIT 1;

-- name: CreateTenantSubscription :one
INSERT INTO hrms.tenant_subscriptions (
    tenant_id,
    plan_id,
    start_date,
    end_date,
    status,
    max_employees,
    created_by,
    updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $7
)
RETURNING *;

-- name: UpdateTenantSubscription :one
UPDATE hrms.tenant_subscriptions
SET plan_id = $3,
    start_date = $4,
    end_date = $5,
    status = $6,
    max_employees = $7,
    updated_by = $8,
    updated_at = NOW()
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: SoftDeleteTenantSubscription :exec
UPDATE hrms.tenant_subscriptions
SET inactive = TRUE,
    updated_by = $3,
    updated_at = NOW()
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;
