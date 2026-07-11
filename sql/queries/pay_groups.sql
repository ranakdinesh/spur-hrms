-- name: CreatePayGroup :one
INSERT INTO hrms.pay_groups (
    tenant_id, code, name, description, grouping_type, branch_id, department_id,
    employment_type_id, reporting_tag, rules, is_active, created_by, updated_by
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $12)
RETURNING *;

-- name: UpdatePayGroup :one
UPDATE hrms.pay_groups
SET code = $3,
    name = $4,
    description = $5,
    grouping_type = $6,
    branch_id = $7,
    department_id = $8,
    employment_type_id = $9,
    reporting_tag = $10,
    rules = $11,
    is_active = $12,
    updated_by = $13,
    updated_at = NOW()
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: GetPayGroup :one
SELECT * FROM hrms.pay_groups
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: ListPayGroups :many
SELECT * FROM hrms.pay_groups
WHERE tenant_id = $1 AND NOT inactive
ORDER BY is_active DESC, name ASC;

-- name: SoftDeletePayGroup :exec
UPDATE hrms.pay_groups
SET inactive = TRUE, is_active = FALSE, updated_by = $3, updated_at = NOW()
WHERE tenant_id = $1 AND id = $2;

-- name: UpsertPayGroupMember :one
INSERT INTO hrms.pay_group_members (
    tenant_id, pay_group_id, user_id, membership_type, effective_from, effective_to, created_by, updated_by
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $7)
ON CONFLICT (tenant_id, pay_group_id, user_id, membership_type) WHERE NOT inactive
DO UPDATE SET
    effective_from = EXCLUDED.effective_from,
    effective_to = EXCLUDED.effective_to,
    updated_by = EXCLUDED.updated_by,
    updated_at = NOW()
RETURNING *;

-- name: ListPayGroupMembers :many
SELECT * FROM hrms.pay_group_members
WHERE tenant_id = $1 AND pay_group_id = $2 AND NOT inactive
ORDER BY membership_type ASC, created_at DESC;

-- name: SoftDeletePayGroupMember :exec
UPDATE hrms.pay_group_members
SET inactive = TRUE, updated_by = $3, updated_at = NOW()
WHERE tenant_id = $1 AND id = $2;

-- name: ListPayGroupEmployees :many
WITH group_row AS (
    SELECT * FROM hrms.pay_groups
    WHERE hrms.pay_groups.tenant_id = $1 AND hrms.pay_groups.id = $2 AND NOT hrms.pay_groups.inactive AND hrms.pay_groups.is_active
)
SELECT
    e.id AS employee_id,
    e.user_id,
    e.employee_code,
    e.firstname,
    e.lastname,
    e.branch_id,
    b.branch_name,
    e.department_id,
    d.name AS department_name,
    e.employment_type_id,
    et.name AS employment_type_name,
    CASE
        WHEN EXISTS (
            SELECT 1 FROM hrms.pay_group_members pgm
            WHERE pgm.tenant_id = e.tenant_id AND pgm.pay_group_id = (SELECT id FROM group_row)
              AND pgm.user_id = e.user_id AND pgm.membership_type = 'manual_include'
              AND NOT pgm.inactive
        ) THEN 'manual_include'
        ELSE 'rule'
    END AS match_source
FROM hrms.employees e
JOIN group_row pg ON pg.tenant_id = e.tenant_id
LEFT JOIN hrms.branches b ON b.tenant_id = e.tenant_id AND b.id = e.branch_id AND NOT b.inactive
LEFT JOIN hrms.departments d ON d.tenant_id = e.tenant_id AND d.id = e.department_id AND NOT d.inactive
LEFT JOIN hrms.employment_types et ON et.tenant_id = e.tenant_id AND et.id = e.employment_type_id AND NOT et.inactive
WHERE e.tenant_id = $1
  AND NOT e.inactive
  AND (
      pg.grouping_type = 'all'
      OR (pg.grouping_type IN ('branch','mixed') AND pg.branch_id IS NOT NULL AND e.branch_id = pg.branch_id)
      OR (pg.grouping_type IN ('department','mixed') AND pg.department_id IS NOT NULL AND e.department_id = pg.department_id)
      OR (pg.grouping_type IN ('employment_type','mixed') AND pg.employment_type_id IS NOT NULL AND e.employment_type_id = pg.employment_type_id)
      OR (pg.grouping_type IN ('reporting_tag','mixed') AND pg.reporting_tag IS NOT NULL AND lower(coalesce(e.grade, '')) = lower(pg.reporting_tag))
      OR EXISTS (
          SELECT 1 FROM hrms.pay_group_members pgm
          WHERE pgm.tenant_id = e.tenant_id AND pgm.pay_group_id = pg.id
            AND pgm.user_id = e.user_id AND pgm.membership_type = 'manual_include'
            AND NOT pgm.inactive
      )
  )
  AND NOT EXISTS (
      SELECT 1 FROM hrms.pay_group_members pgm
      WHERE pgm.tenant_id = e.tenant_id AND pgm.pay_group_id = pg.id
        AND pgm.user_id = e.user_id AND pgm.membership_type = 'manual_exclude'
        AND NOT pgm.inactive
  )
ORDER BY e.employee_code ASC NULLS LAST, e.firstname ASC;

-- name: CreatePayRun :one
INSERT INTO hrms.pay_runs (
    tenant_id, pay_group_id, fy_id, month, year, status, employee_count,
    ready_count, blocked_count, generated_count, readiness, notes, created_by, updated_by
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $13)
RETURNING *;

-- name: GetPayRun :one
SELECT * FROM hrms.pay_runs
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: ListPayRuns :many
SELECT * FROM hrms.pay_runs
WHERE tenant_id = $1
  AND (sqlc.narg('pay_group_id')::uuid IS NULL OR pay_group_id = sqlc.narg('pay_group_id')::uuid)
  AND (sqlc.narg('month')::int IS NULL OR month = sqlc.narg('month')::int)
  AND (sqlc.narg('year')::int IS NULL OR year = sqlc.narg('year')::int)
  AND NOT inactive
ORDER BY year DESC, month DESC, created_at DESC;

-- name: UpdatePayRunStatus :one
UPDATE hrms.pay_runs
SET status = $3,
    employee_count = $4,
    ready_count = $5,
    blocked_count = $6,
    generated_count = $7,
    readiness = $8,
    notes = $9,
    attendance_frozen_at = CASE WHEN $3 = 'frozen' AND attendance_frozen_at IS NULL THEN NOW() ELSE attendance_frozen_at END,
    lop_frozen_at = CASE WHEN $3 = 'frozen' AND lop_frozen_at IS NULL THEN NOW() ELSE lop_frozen_at END,
    adjustments_frozen_at = CASE WHEN $3 = 'frozen' AND adjustments_frozen_at IS NULL THEN NOW() ELSE adjustments_frozen_at END,
    generated_at = CASE WHEN $3 = 'generated' THEN NOW() ELSE generated_at END,
    locked_at = CASE WHEN $3 = 'locked' THEN NOW() ELSE locked_at END,
    locked_by = CASE WHEN $3 = 'locked' THEN $10 ELSE locked_by END,
    unlocked_at = CASE WHEN $3 = 'unlocked' THEN NOW() ELSE unlocked_at END,
    unlocked_by = CASE WHEN $3 = 'unlocked' THEN $10 ELSE unlocked_by END,
    updated_by = $10,
    updated_at = NOW()
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: SoftDeletePayRun :exec
UPDATE hrms.pay_runs
SET inactive = TRUE, updated_by = $3, updated_at = NOW()
WHERE tenant_id = $1 AND id = $2 AND status IN ('draft','blocked','failed','unlocked');

-- name: UpsertPayRunEmployee :one
INSERT INTO hrms.pay_run_employees (
    tenant_id, pay_run_id, user_id, readiness_status, blocker_reason, salary_slip_id, generated_at, created_by, updated_by
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $8)
ON CONFLICT (tenant_id, pay_run_id, user_id) WHERE NOT inactive
DO UPDATE SET
    readiness_status = EXCLUDED.readiness_status,
    blocker_reason = EXCLUDED.blocker_reason,
    salary_slip_id = EXCLUDED.salary_slip_id,
    generated_at = EXCLUDED.generated_at,
    updated_by = EXCLUDED.updated_by,
    updated_at = NOW()
RETURNING *;

-- name: ListPayRunEmployees :many
SELECT
    pre.*,
    e.employee_code,
    e.firstname,
    e.lastname,
    b.branch_name,
    d.name AS department_name,
    COALESCE(SUM(prc.amount) FILTER (WHERE prc.component_type = 'earning'), 0)::numeric AS earnings_amount,
    COALESCE(SUM(prc.amount) FILTER (WHERE prc.component_type = 'deduction'), 0)::numeric AS deductions_amount,
    COALESCE(SUM(prc.amount) FILTER (WHERE prc.component_type = 'employer_contribution'), 0)::numeric AS employer_cost_amount,
    COALESCE(SUM(prc.amount) FILTER (WHERE prc.component_type = 'earning'), 0)::numeric AS gross_amount,
    (
        COALESCE(SUM(prc.amount) FILTER (WHERE prc.component_type = 'earning'), 0)
        - COALESCE(SUM(prc.amount) FILTER (WHERE prc.component_type = 'deduction'), 0)
    )::numeric AS net_amount
FROM hrms.pay_run_employees pre
JOIN hrms.employees e ON e.tenant_id = pre.tenant_id AND e.user_id = pre.user_id AND NOT e.inactive
LEFT JOIN hrms.branches b ON b.tenant_id = e.tenant_id AND b.id = e.branch_id AND NOT b.inactive
LEFT JOIN hrms.departments d ON d.tenant_id = e.tenant_id AND d.id = e.department_id AND NOT d.inactive
LEFT JOIN hrms.pay_run_components prc ON prc.tenant_id = pre.tenant_id AND prc.pay_run_id = pre.pay_run_id AND prc.user_id = pre.user_id AND NOT prc.inactive
WHERE pre.tenant_id = $1 AND pre.pay_run_id = $2 AND NOT pre.inactive
GROUP BY pre.id, e.employee_code, e.firstname, e.lastname, b.branch_name, d.name
ORDER BY pre.readiness_status ASC, e.employee_code ASC NULLS LAST, e.firstname ASC;

-- name: DeletePayRunLedger :exec
UPDATE hrms.pay_run_inputs
SET inactive = TRUE, updated_by = $3, updated_at = NOW()
WHERE tenant_id = $1 AND pay_run_id = $2 AND NOT inactive;

-- name: DeletePayRunComponentLedger :exec
UPDATE hrms.pay_run_components
SET inactive = TRUE, updated_by = $3, updated_at = NOW()
WHERE tenant_id = $1 AND pay_run_id = $2 AND NOT inactive;

-- name: CreatePayRunInput :one
INSERT INTO hrms.pay_run_inputs (
    tenant_id, pay_run_id, user_id, input_type, source_type, source_id,
    description, quantity, amount, metadata, created_by, updated_by
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $11)
RETURNING *;

-- name: CreatePayRunComponent :one
INSERT INTO hrms.pay_run_components (
    tenant_id, pay_run_id, user_id, component_type, code, name, amount,
    source_input_id, salary_template_id, taxable, statutory, employer_cost,
    sort_order, metadata, created_by, updated_by
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $15)
RETURNING *;

-- name: ListPayRunInputs :many
SELECT
    pri.*,
    e.employee_code,
    e.firstname,
    e.lastname
FROM hrms.pay_run_inputs pri
JOIN hrms.employees e ON e.tenant_id = pri.tenant_id AND e.user_id = pri.user_id AND NOT e.inactive
WHERE pri.tenant_id = $1 AND pri.pay_run_id = $2 AND NOT pri.inactive
ORDER BY e.employee_code ASC NULLS LAST, e.firstname ASC, pri.input_type ASC, pri.created_at ASC;

-- name: ListPayRunComponents :many
SELECT
    prc.*,
    e.employee_code,
    e.firstname,
    e.lastname
FROM hrms.pay_run_components prc
JOIN hrms.employees e ON e.tenant_id = prc.tenant_id AND e.user_id = prc.user_id AND NOT e.inactive
WHERE prc.tenant_id = $1 AND prc.pay_run_id = $2 AND NOT prc.inactive
ORDER BY e.employee_code ASC NULLS LAST, e.firstname ASC, prc.sort_order ASC, prc.created_at ASC;

-- name: GetPayRunLedgerSummary :one
SELECT
    pr.id AS pay_run_id,
    pr.employee_count::int AS employee_count,
    COUNT(DISTINCT prc.user_id)::int AS draft_employee_count,
    COALESCE(SUM(prc.amount) FILTER (WHERE prc.component_type = 'earning'), 0)::numeric AS gross_amount,
    COALESCE(SUM(prc.amount) FILTER (WHERE prc.component_type = 'earning'), 0)::numeric AS total_earnings,
    COALESCE(SUM(prc.amount) FILTER (WHERE prc.component_type = 'deduction'), 0)::numeric AS total_deductions,
    (
        COALESCE(SUM(prc.amount) FILTER (WHERE prc.component_type = 'earning'), 0)
        - COALESCE(SUM(prc.amount) FILTER (WHERE prc.component_type = 'deduction'), 0)
    )::numeric AS net_amount,
    COALESCE(SUM(prc.amount) FILTER (WHERE prc.component_type = 'employer_contribution'), 0)::numeric AS employer_cost_amount,
    (SELECT COUNT(*)::int FROM hrms.pay_run_inputs pri WHERE pri.tenant_id = pr.tenant_id AND pri.pay_run_id = pr.id AND NOT pri.inactive) AS input_count,
    COUNT(prc.id)::int AS component_count
FROM hrms.pay_runs pr
LEFT JOIN hrms.pay_run_components prc ON prc.tenant_id = pr.tenant_id AND prc.pay_run_id = pr.id AND NOT prc.inactive
WHERE pr.tenant_id = $1 AND pr.id = $2 AND NOT pr.inactive
GROUP BY pr.id, pr.employee_count;

-- name: CreatePayRunEvent :one
INSERT INTO hrms.pay_run_events (
    tenant_id, pay_run_id, action, from_status, to_status, remarks, metadata, created_by, updated_by
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $8)
RETURNING *;

-- name: ListPayRunEvents :many
SELECT * FROM hrms.pay_run_events
WHERE tenant_id = $1 AND pay_run_id = $2 AND NOT inactive
ORDER BY created_at DESC;
