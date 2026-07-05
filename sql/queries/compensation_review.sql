-- name: CreateCompensationPayBand :one
INSERT INTO hrms.compensation_pay_bands (
    tenant_id, code, name, job_family, level_code, location_label, currency_code,
    min_pay, midpoint_pay, max_pay, effective_from, effective_to, is_active, notes, metadata,
    created_by, updated_by
)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$16)
RETURNING *;

-- name: UpdateCompensationPayBand :one
UPDATE hrms.compensation_pay_bands
SET code = $3,
    name = $4,
    job_family = $5,
    level_code = $6,
    location_label = $7,
    currency_code = $8,
    min_pay = $9,
    midpoint_pay = $10,
    max_pay = $11,
    effective_from = $12,
    effective_to = $13,
    is_active = $14,
    notes = $15,
    metadata = $16,
    updated_by = $17,
    updated_at = now()
WHERE tenant_id = $1 AND id = $2 AND inactive = FALSE
RETURNING *;

-- name: GetCompensationPayBand :one
SELECT * FROM hrms.compensation_pay_bands
WHERE tenant_id = $1 AND id = $2 AND inactive = FALSE;

-- name: ListCompensationPayBands :many
SELECT * FROM hrms.compensation_pay_bands
WHERE tenant_id = $1
  AND inactive = FALSE
  AND (sqlc.narg('is_active')::boolean IS NULL OR is_active = sqlc.narg('is_active')::boolean)
  AND (sqlc.narg('currency_code')::text IS NULL OR currency_code = sqlc.narg('currency_code')::text)
  AND (sqlc.narg('search')::text IS NULL OR code ILIKE '%' || sqlc.narg('search')::text || '%' OR name ILIKE '%' || sqlc.narg('search')::text || '%' OR COALESCE(job_family, '') ILIKE '%' || sqlc.narg('search')::text || '%')
ORDER BY is_active DESC, level_code NULLS LAST, name ASC
LIMIT $2 OFFSET $3;

-- name: SoftDeleteCompensationPayBand :exec
UPDATE hrms.compensation_pay_bands
SET inactive = TRUE, is_active = FALSE, updated_by = $3, updated_at = now()
WHERE tenant_id = $1 AND id = $2;

-- name: CreateCompensationCycle :one
INSERT INTO hrms.compensation_cycles (
    tenant_id, code, name, fiscal_year_id, status, cycle_type, starts_on, ends_on,
    effective_date, currency_code, budget_amount, planning_guidance, approval_policy, metadata,
    created_by, updated_by
)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$15)
RETURNING *;

-- name: UpdateCompensationCycle :one
UPDATE hrms.compensation_cycles
SET code = $3,
    name = $4,
    fiscal_year_id = $5,
    cycle_type = $6,
    starts_on = $7,
    ends_on = $8,
    effective_date = $9,
    currency_code = $10,
    budget_amount = $11,
    planning_guidance = $12,
    approval_policy = $13,
    metadata = $14,
    updated_by = $15,
    updated_at = now()
WHERE tenant_id = $1 AND id = $2 AND inactive = FALSE
RETURNING *;

-- name: UpdateCompensationCycleStatus :one
UPDATE hrms.compensation_cycles
SET status = $3,
    finalized_at = CASE WHEN $3 = 'finalized' THEN now() ELSE finalized_at END,
    updated_by = $4,
    updated_at = now()
WHERE tenant_id = $1 AND id = $2 AND inactive = FALSE
RETURNING *;

-- name: GetCompensationCycle :one
SELECT * FROM hrms.compensation_cycles
WHERE tenant_id = $1 AND id = $2 AND inactive = FALSE;

-- name: ListCompensationCycles :many
SELECT * FROM hrms.compensation_cycles
WHERE tenant_id = $1
  AND inactive = FALSE
  AND (sqlc.narg('status')::text IS NULL OR status = sqlc.narg('status')::text)
  AND (sqlc.narg('search')::text IS NULL OR code ILIKE '%' || sqlc.narg('search')::text || '%' OR name ILIKE '%' || sqlc.narg('search')::text || '%')
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CreateCompensationBudgetPool :one
INSERT INTO hrms.compensation_budget_pools (
    tenant_id, cycle_id, name, pool_type, owner_user_id, department_id, branch_id,
    budget_amount, allocated_amount, notes, created_by, updated_by
)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$11)
RETURNING *;

-- name: UpdateCompensationBudgetPool :one
UPDATE hrms.compensation_budget_pools
SET name = $4,
    pool_type = $5,
    owner_user_id = $6,
    department_id = $7,
    branch_id = $8,
    budget_amount = $9,
    allocated_amount = $10,
    notes = $11,
    updated_by = $12,
    updated_at = now()
WHERE tenant_id = $1 AND id = $2 AND cycle_id = $3 AND inactive = FALSE
RETURNING *;

-- name: ListCompensationBudgetPools :many
SELECT bp.*,
       COALESCE(SUM(cr.recommended_increment_amount) FILTER (WHERE cr.inactive = FALSE AND cr.status IN ('submitted','approved','finalized','handed_to_payroll')), 0)::numeric AS committed_amount,
       COALESCE(COUNT(cr.id) FILTER (WHERE cr.inactive = FALSE), 0)::bigint AS recommendation_count
FROM hrms.compensation_budget_pools bp
LEFT JOIN hrms.compensation_recommendations cr ON cr.tenant_id = bp.tenant_id AND cr.budget_pool_id = bp.id
WHERE bp.tenant_id = $1 AND bp.cycle_id = $2 AND bp.inactive = FALSE
GROUP BY bp.id
ORDER BY bp.pool_type ASC, bp.name ASC;

-- name: SoftDeleteCompensationBudgetPool :exec
UPDATE hrms.compensation_budget_pools
SET inactive = TRUE, updated_by = $4, updated_at = now()
WHERE tenant_id = $1 AND cycle_id = $2 AND id = $3;

-- name: CreateCompensationRecommendation :one
INSERT INTO hrms.compensation_recommendations (
    tenant_id, cycle_id, worker_profile_id, pay_band_id, budget_pool_id, current_salary,
    current_compa_ratio, recommended_salary, recommended_increment_amount,
    recommended_increment_percent, promotion_recommended, recommended_designation_id,
    reason, performance_rating, equity_flag, equity_notes, status, effective_date, metadata,
    created_by, updated_by
)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$20)
RETURNING *;

-- name: UpdateCompensationRecommendation :one
UPDATE hrms.compensation_recommendations
SET pay_band_id = $4,
    budget_pool_id = $5,
    current_salary = $6,
    current_compa_ratio = $7,
    recommended_salary = $8,
    recommended_increment_amount = $9,
    recommended_increment_percent = $10,
    promotion_recommended = $11,
    recommended_designation_id = $12,
    reason = $13,
    performance_rating = $14,
    equity_flag = $15,
    equity_notes = $16,
    effective_date = $17,
    metadata = $18,
    updated_by = $19,
    updated_at = now()
WHERE tenant_id = $1 AND id = $2 AND cycle_id = $3 AND inactive = FALSE
RETURNING *;

-- name: UpdateCompensationRecommendationStatus :one
UPDATE hrms.compensation_recommendations
SET status = $3,
    payroll_handoff_at = CASE WHEN $3 = 'handed_to_payroll' THEN now() ELSE payroll_handoff_at END,
    updated_by = $4,
    updated_at = now()
WHERE tenant_id = $1 AND id = $2 AND inactive = FALSE
RETURNING *;

-- name: GetCompensationRecommendation :one
SELECT cr.*, wp.display_name AS worker_display_name, wp.worker_code, pb.code AS pay_band_code, pb.name AS pay_band_name, bp.name AS budget_pool_name
FROM hrms.compensation_recommendations cr
JOIN hrms.worker_profiles wp ON wp.tenant_id = cr.tenant_id AND wp.id = cr.worker_profile_id
LEFT JOIN hrms.compensation_pay_bands pb ON pb.tenant_id = cr.tenant_id AND pb.id = cr.pay_band_id
LEFT JOIN hrms.compensation_budget_pools bp ON bp.tenant_id = cr.tenant_id AND bp.id = cr.budget_pool_id
WHERE cr.tenant_id = $1 AND cr.id = $2 AND cr.inactive = FALSE;

-- name: ListCompensationRecommendations :many
SELECT cr.*, wp.display_name AS worker_display_name, wp.worker_code, pb.code AS pay_band_code, pb.name AS pay_band_name, bp.name AS budget_pool_name
FROM hrms.compensation_recommendations cr
JOIN hrms.worker_profiles wp ON wp.tenant_id = cr.tenant_id AND wp.id = cr.worker_profile_id
LEFT JOIN hrms.compensation_pay_bands pb ON pb.tenant_id = cr.tenant_id AND pb.id = cr.pay_band_id
LEFT JOIN hrms.compensation_budget_pools bp ON bp.tenant_id = cr.tenant_id AND bp.id = cr.budget_pool_id
WHERE cr.tenant_id = $1
  AND cr.inactive = FALSE
  AND (sqlc.narg('cycle_id')::uuid IS NULL OR cr.cycle_id = sqlc.narg('cycle_id')::uuid)
  AND (sqlc.narg('status')::text IS NULL OR cr.status = sqlc.narg('status')::text)
  AND (sqlc.narg('search')::text IS NULL OR wp.display_name ILIKE '%' || sqlc.narg('search')::text || '%' OR COALESCE(wp.worker_code, '') ILIKE '%' || sqlc.narg('search')::text || '%')
ORDER BY cr.updated_at DESC
LIMIT $2 OFFSET $3;

-- name: SoftDeleteCompensationRecommendation :exec
UPDATE hrms.compensation_recommendations
SET inactive = TRUE, updated_by = $3, updated_at = now()
WHERE tenant_id = $1 AND id = $2;

-- name: GenerateCompensationEquityChecks :many
WITH inserted AS (
    INSERT INTO hrms.compensation_equity_checks (
        tenant_id, cycle_id, worker_profile_id, pay_band_id, check_type, severity,
        current_salary, band_min, band_midpoint, band_max, variance_percent,
        finding, recommendation, created_by, updated_by
    )
    SELECT cr.tenant_id,
           cr.cycle_id,
           cr.worker_profile_id,
           cr.pay_band_id,
           CASE
               WHEN cr.current_salary < pb.min_pay THEN 'below_band'
               WHEN cr.current_salary > pb.max_pay THEN 'above_band'
               WHEN pb.midpoint_pay > 0 AND cr.current_salary / pb.midpoint_pay < 0.80 THEN 'low_compa_ratio'
               WHEN pb.midpoint_pay > 0 AND cr.current_salary / pb.midpoint_pay > 1.20 THEN 'high_compa_ratio'
               ELSE 'band_position'
           END,
           CASE
               WHEN cr.current_salary < pb.min_pay * 0.90 OR cr.current_salary > pb.max_pay * 1.10 THEN 'critical'
               WHEN cr.current_salary < pb.min_pay OR cr.current_salary > pb.max_pay THEN 'high'
               ELSE 'medium'
           END,
           cr.current_salary,
           pb.min_pay,
           pb.midpoint_pay,
           pb.max_pay,
           CASE WHEN pb.midpoint_pay > 0 THEN round(((cr.current_salary - pb.midpoint_pay) / pb.midpoint_pay) * 100, 4) ELSE 0 END,
           CASE
               WHEN cr.current_salary < pb.min_pay THEN 'Current salary is below band minimum.'
               WHEN cr.current_salary > pb.max_pay THEN 'Current salary is above band maximum.'
               WHEN pb.midpoint_pay > 0 AND cr.current_salary / pb.midpoint_pay < 0.80 THEN 'Compa-ratio is materially below midpoint.'
               WHEN pb.midpoint_pay > 0 AND cr.current_salary / pb.midpoint_pay > 1.20 THEN 'Compa-ratio is materially above midpoint.'
               ELSE 'Current salary is inside the band but should be reviewed for equity.'
           END,
           CASE
               WHEN cr.current_salary < pb.min_pay THEN 'Review correction adjustment to bring pay within band.'
               WHEN cr.current_salary > pb.max_pay THEN 'Review role leveling, promotion history, or red-circle handling.'
               ELSE 'Review peer equity and performance context before approval.'
           END,
           $2,
           $2
    FROM hrms.compensation_recommendations cr
    JOIN hrms.compensation_pay_bands pb ON pb.tenant_id = cr.tenant_id AND pb.id = cr.pay_band_id
    WHERE cr.tenant_id = $1
      AND cr.cycle_id = $3
      AND cr.inactive = FALSE
      AND pb.inactive = FALSE
      AND (
          cr.current_salary < pb.min_pay
          OR cr.current_salary > pb.max_pay
          OR (pb.midpoint_pay > 0 AND (cr.current_salary / pb.midpoint_pay < 0.80 OR cr.current_salary / pb.midpoint_pay > 1.20))
      )
    RETURNING *
)
SELECT * FROM inserted
ORDER BY severity DESC, created_at DESC;

-- name: ListCompensationEquityChecks :many
SELECT ec.*, wp.display_name AS worker_display_name, wp.worker_code, pb.code AS pay_band_code, pb.name AS pay_band_name
FROM hrms.compensation_equity_checks ec
JOIN hrms.worker_profiles wp ON wp.tenant_id = ec.tenant_id AND wp.id = ec.worker_profile_id
LEFT JOIN hrms.compensation_pay_bands pb ON pb.tenant_id = ec.tenant_id AND pb.id = ec.pay_band_id
WHERE ec.tenant_id = $1
  AND ec.inactive = FALSE
  AND (sqlc.narg('cycle_id')::uuid IS NULL OR ec.cycle_id = sqlc.narg('cycle_id')::uuid)
  AND (sqlc.narg('status')::text IS NULL OR ec.status = sqlc.narg('status')::text)
ORDER BY ec.created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateCompensationEquityCheckStatus :one
UPDATE hrms.compensation_equity_checks
SET status = $3,
    updated_by = $4,
    updated_at = now()
WHERE tenant_id = $1 AND id = $2 AND inactive = FALSE
RETURNING *;

-- name: CreateCompensationEvent :one
INSERT INTO hrms.compensation_events (
    tenant_id, cycle_id, source_type, source_id, action, from_status, to_status, remarks, metadata, created_by
)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
RETURNING *;

-- name: ListCompensationEvents :many
SELECT * FROM hrms.compensation_events
WHERE tenant_id = $1
  AND (sqlc.narg('cycle_id')::uuid IS NULL OR cycle_id = sqlc.narg('cycle_id')::uuid)
  AND (sqlc.narg('source_type')::text IS NULL OR source_type = sqlc.narg('source_type')::text)
  AND (sqlc.narg('source_id')::uuid IS NULL OR source_id = sqlc.narg('source_id')::uuid)
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetCompensationSummary :many
SELECT 'cycles'::text AS metric, COUNT(*)::bigint AS metric_count, COALESCE(SUM(budget_amount), 0)::numeric AS metric_amount
FROM hrms.compensation_cycles cc
WHERE cc.tenant_id = $1 AND cc.inactive = FALSE
UNION ALL
SELECT 'open_cycles', COUNT(*)::bigint, COALESCE(SUM(budget_amount), 0)::numeric
FROM hrms.compensation_cycles cc
WHERE cc.tenant_id = $1 AND cc.inactive = FALSE AND cc.status IN ('draft','open','submitted')
UNION ALL
SELECT 'recommendations', COUNT(*)::bigint, COALESCE(SUM(recommended_increment_amount), 0)::numeric
FROM hrms.compensation_recommendations cr
WHERE cr.tenant_id = $1 AND cr.inactive = FALSE
UNION ALL
SELECT 'approved_recommendations', COUNT(*)::bigint, COALESCE(SUM(recommended_increment_amount), 0)::numeric
FROM hrms.compensation_recommendations cr
WHERE cr.tenant_id = $1 AND cr.inactive = FALSE AND cr.status IN ('approved','finalized','handed_to_payroll')
UNION ALL
SELECT 'equity_flags', COUNT(*)::bigint, 0::numeric
FROM hrms.compensation_equity_checks ec
WHERE ec.tenant_id = $1 AND ec.inactive = FALSE AND ec.status = 'open';
