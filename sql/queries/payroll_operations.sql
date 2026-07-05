-- name: UpsertPayrollPeriodLock :one
INSERT INTO hrms.payroll_period_locks (
    tenant_id, month, year, status, locked_at, locked_by, unlocked_at, unlocked_by, unlock_reason, notes, created_by, updated_by
)
VALUES (
    $1, $2, $3, $4,
    CASE WHEN $4 = 'locked' THEN NOW() ELSE NULL END,
    CASE WHEN $4 = 'locked' THEN $5 ELSE NULL END,
    CASE WHEN $4 = 'unlocked' THEN NOW() ELSE NULL END,
    CASE WHEN $4 = 'unlocked' THEN $5 ELSE NULL END,
    $6, $7, $5, $5
)
ON CONFLICT (tenant_id, month, year) WHERE NOT inactive
DO UPDATE SET
    status = EXCLUDED.status,
    locked_at = CASE WHEN EXCLUDED.status = 'locked' THEN NOW() ELSE hrms.payroll_period_locks.locked_at END,
    locked_by = CASE WHEN EXCLUDED.status = 'locked' THEN EXCLUDED.updated_by ELSE hrms.payroll_period_locks.locked_by END,
    unlocked_at = CASE WHEN EXCLUDED.status = 'unlocked' THEN NOW() ELSE hrms.payroll_period_locks.unlocked_at END,
    unlocked_by = CASE WHEN EXCLUDED.status = 'unlocked' THEN EXCLUDED.updated_by ELSE hrms.payroll_period_locks.unlocked_by END,
    unlock_reason = EXCLUDED.unlock_reason,
    notes = EXCLUDED.notes,
    updated_by = EXCLUDED.updated_by,
    updated_at = NOW()
RETURNING *;

-- name: GetPayrollPeriodLock :one
SELECT * FROM hrms.payroll_period_locks
WHERE tenant_id = $1 AND month = $2 AND year = $3 AND NOT inactive;

-- name: ListPayrollPeriodLocks :many
SELECT * FROM hrms.payroll_period_locks
WHERE tenant_id = $1 AND NOT inactive
ORDER BY year DESC, month DESC;

-- name: CreatePayrollPeriodLockEvent :one
INSERT INTO hrms.payroll_period_lock_events (
    tenant_id, payroll_lock_id, action, from_status, to_status, remarks, created_by, updated_by
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $7)
RETURNING *;

-- name: ListPayrollPeriodLockEvents :many
SELECT * FROM hrms.payroll_period_lock_events
WHERE tenant_id = $1 AND payroll_lock_id = $2 AND NOT inactive
ORDER BY created_at DESC;

-- name: ListPayrollStatutoryRules :many
SELECT * FROM hrms.payroll_statutory_rules
WHERE tenant_id = $1
  AND (sqlc.narg('rule_type')::text IS NULL OR rule_type = sqlc.narg('rule_type')::text)
  AND NOT inactive
ORDER BY rule_type ASC, state ASC NULLS LAST, effective_from DESC;

-- name: GetPayrollStatutoryRule :one
SELECT * FROM hrms.payroll_statutory_rules
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: CreatePayrollStatutoryRule :one
INSERT INTO hrms.payroll_statutory_rules (
    tenant_id, rule_type, name, state, branch_id, effective_from, effective_to,
    min_gross_salary, max_gross_salary, employee_amount, employer_amount,
    frequency, deduction_month, notes, created_by, updated_by
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $15)
RETURNING *;

-- name: UpdatePayrollStatutoryRule :one
UPDATE hrms.payroll_statutory_rules
SET rule_type = $3,
    name = $4,
    state = $5,
    branch_id = $6,
    effective_from = $7,
    effective_to = $8,
    min_gross_salary = $9,
    max_gross_salary = $10,
    employee_amount = $11,
    employer_amount = $12,
    frequency = $13,
    deduction_month = $14,
    notes = $15,
    updated_by = $16,
    updated_at = NOW()
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: SoftDeletePayrollStatutoryRule :exec
UPDATE hrms.payroll_statutory_rules
SET inactive = TRUE, updated_by = $3, updated_at = NOW()
WHERE tenant_id = $1 AND id = $2;

-- name: ResolvePayrollStatutoryRule :one
SELECT * FROM hrms.payroll_statutory_rules
WHERE tenant_id = $1
  AND rule_type = $2
  AND (branch_id IS NULL OR branch_id = sqlc.narg('branch_id')::uuid)
  AND (state IS NULL OR LOWER(state) = LOWER(sqlc.narg('state')::text))
  AND effective_from <= $3
  AND (effective_to IS NULL OR effective_to >= $3)
  AND (min_gross_salary IS NULL OR min_gross_salary <= $4)
  AND (max_gross_salary IS NULL OR max_gross_salary >= $4)
  AND (deduction_month IS NULL OR deduction_month = $5)
  AND NOT inactive
ORDER BY
  CASE WHEN branch_id = sqlc.narg('branch_id')::uuid THEN 0 ELSE 1 END,
  CASE WHEN LOWER(state) = LOWER(sqlc.narg('state')::text) THEN 0 ELSE 1 END,
  effective_from DESC
LIMIT 1;

-- name: CreatePayrollImportBatch :one
INSERT INTO hrms.payroll_import_batches (
    tenant_id, import_type, month, year, fy_id, template_id, file_name, status,
    total_rows, valid_rows, invalid_rows, applied_rows, error_report, notes, created_by, updated_by
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $15)
RETURNING *;

-- name: ListPayrollImportBatches :many
SELECT * FROM hrms.payroll_import_batches
WHERE tenant_id = $1 AND NOT inactive
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetPayrollImportBatch :one
SELECT * FROM hrms.payroll_import_batches
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: CreatePayrollImportRow :one
INSERT INTO hrms.payroll_import_rows (
    tenant_id, batch_id, row_number, employee_code, employee_user_id, employee_name,
    gross_salary, present_days, absent_days, lop_days, variable_earnings,
    variable_deductions, status, error_message, raw_data, created_by, updated_by
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $16)
RETURNING *;

-- name: ListPayrollImportRows :many
SELECT * FROM hrms.payroll_import_rows
WHERE tenant_id = $1 AND batch_id = $2 AND NOT inactive
ORDER BY row_number ASC;

-- name: ListConsolidatedSalarySheet :many
SELECT
    ss.id,
    ss.tenant_id,
    ss.user_id,
    e.employee_code,
    e.firstname,
    e.lastname,
    e.email,
    b.branch_name,
    d.name AS department_name,
    ss.month,
    ss.year,
    ss.gross_salary,
    ss.total_earnings,
    ss.total_deductions,
    ss.absent_deduction,
    ss.net_salary,
    ss.present_days,
    ss.absent_days,
    ss.lwp_days,
    ss.pdf_path,
    ss.created_at
FROM hrms.salary_slips ss
JOIN hrms.employees e ON e.tenant_id = ss.tenant_id AND e.user_id = ss.user_id
LEFT JOIN hrms.branches b ON b.tenant_id = e.tenant_id AND b.id = e.branch_id AND NOT b.inactive
LEFT JOIN hrms.departments d ON d.tenant_id = e.tenant_id AND d.id = e.department_id AND NOT d.inactive
WHERE ss.tenant_id = $1 AND ss.month = $2 AND ss.year = $3 AND NOT ss.inactive
ORDER BY e.employee_code ASC NULLS LAST, e.firstname ASC;

-- name: ListPayrollReconciliationRows :many
SELECT
    e.id AS employee_id,
    e.user_id,
    e.employee_code,
    e.firstname,
    e.lastname,
    e.email,
    b.branch_name,
    d.name AS department_name,
    ss.id AS salary_slip_id,
    ss.present_days,
    ss.absent_days,
    ss.lwp_days,
    ss.net_salary,
    CASE WHEN ss.id IS NULL THEN 'missing_payslip'
         WHEN ss.absent_days > 0 AND ss.absent_deduction <= 0 THEN 'lop_without_deduction'
         ELSE 'ok'
    END AS reconciliation_status
FROM hrms.employees e
LEFT JOIN hrms.branches b ON b.tenant_id = e.tenant_id AND b.id = e.branch_id AND NOT b.inactive
LEFT JOIN hrms.departments d ON d.tenant_id = e.tenant_id AND d.id = e.department_id AND NOT d.inactive
LEFT JOIN hrms.salary_slips ss ON ss.tenant_id = e.tenant_id AND ss.user_id = e.user_id AND ss.month = $2 AND ss.year = $3 AND NOT ss.inactive
WHERE e.tenant_id = $1 AND NOT e.inactive
ORDER BY reconciliation_status DESC, e.employee_code ASC NULLS LAST, e.firstname ASC;
