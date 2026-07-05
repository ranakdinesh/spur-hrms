-- name: CreateFinancialYear :one
INSERT INTO hrms.financial_years (
    tenant_id,
    name,
    start_date,
    end_date,
    is_active,
    payroll_year,
    leave_year,
    holiday_year,
    reporting_year,
    close_note,
    created_by,
    updated_by
) VALUES (
    $1, $2, $3, $4, FALSE, $5, $6, $7, $8, $9, $10, $10
)
RETURNING *;

-- name: ListFinancialYears :many
SELECT * FROM hrms.financial_years
WHERE tenant_id = $1 AND NOT inactive
ORDER BY start_date DESC, end_date DESC;

-- name: GetFinancialYear :one
SELECT * FROM hrms.financial_years
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: GetActiveFinancialYear :one
SELECT * FROM hrms.financial_years
WHERE tenant_id = $1 AND is_active AND NOT inactive;

-- name: UpdateFinancialYear :one
UPDATE hrms.financial_years
SET name = $3,
    start_date = $4,
    end_date = $5,
    payroll_year = $6,
    leave_year = $7,
    holiday_year = $8,
    reporting_year = $9,
    close_note = $10,
    updated_by = $11,
    updated_at = NOW()
WHERE tenant_id = $1 AND id = $2 AND NOT inactive AND NOT is_locked
RETURNING *;

-- name: ClearActiveFinancialYears :exec
UPDATE hrms.financial_years
SET is_active = FALSE,
    updated_by = $2,
    updated_at = NOW()
WHERE tenant_id = $1 AND NOT inactive;

-- name: MarkFinancialYearActive :one
UPDATE hrms.financial_years
SET is_active = TRUE,
    updated_by = $3,
    updated_at = NOW()
WHERE tenant_id = $1 AND id = $2 AND NOT inactive AND NOT is_locked
RETURNING *;

-- name: SetFinancialYearLock :one
UPDATE hrms.financial_years
SET is_locked = $3,
    locked_at = CASE WHEN $3 THEN NOW() ELSE NULL END,
    locked_by = CASE WHEN $3 THEN $4 ELSE NULL END,
    close_note = $5,
    updated_by = $4,
    updated_at = NOW()
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: SoftDeleteFinancialYear :exec
UPDATE hrms.financial_years
SET inactive = TRUE,
    is_active = FALSE,
    updated_by = $3,
    updated_at = NOW()
WHERE tenant_id = $1 AND id = $2 AND NOT inactive AND NOT is_locked;
