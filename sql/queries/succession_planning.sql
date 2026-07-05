-- name: CreateSuccessionReviewCycle :one
INSERT INTO hrms.succession_review_cycles (
    tenant_id, code, name, status, starts_on, ends_on, confidentiality_level, notes, metadata, created_by, updated_by
)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$10)
RETURNING *;

-- name: UpdateSuccessionReviewCycle :one
UPDATE hrms.succession_review_cycles
SET code = $3,
    name = $4,
    starts_on = $5,
    ends_on = $6,
    confidentiality_level = $7,
    notes = $8,
    metadata = $9,
    updated_by = $10,
    updated_at = now()
WHERE tenant_id = $1 AND id = $2 AND inactive = FALSE
RETURNING *;

-- name: UpdateSuccessionReviewCycleStatus :one
UPDATE hrms.succession_review_cycles
SET status = $3, updated_by = $4, updated_at = now()
WHERE tenant_id = $1 AND id = $2 AND inactive = FALSE
RETURNING *;

-- name: GetSuccessionReviewCycle :one
SELECT * FROM hrms.succession_review_cycles
WHERE tenant_id = $1 AND id = $2 AND inactive = FALSE;

-- name: ListSuccessionReviewCycles :many
SELECT * FROM hrms.succession_review_cycles
WHERE tenant_id = $1
  AND inactive = FALSE
  AND (sqlc.narg('status')::text IS NULL OR status = sqlc.narg('status')::text)
  AND (sqlc.narg('search')::text IS NULL OR code ILIKE '%' || sqlc.narg('search')::text || '%' OR name ILIKE '%' || sqlc.narg('search')::text || '%')
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CreateSuccessionCriticalRole :one
INSERT INTO hrms.succession_critical_roles (
    tenant_id, cycle_id, code, title, department_id, designation_id, incumbent_worker_profile_id,
    emergency_cover_worker_profile_id, criticality, impact_level, vacancy_risk, attrition_risk,
    readiness_target, successor_required_count, role_summary, status, metadata, created_by, updated_by
)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$18)
RETURNING *;

-- name: UpdateSuccessionCriticalRole :one
UPDATE hrms.succession_critical_roles
SET cycle_id = $3,
    code = $4,
    title = $5,
    department_id = $6,
    designation_id = $7,
    incumbent_worker_profile_id = $8,
    emergency_cover_worker_profile_id = $9,
    criticality = $10,
    impact_level = $11,
    vacancy_risk = $12,
    attrition_risk = $13,
    readiness_target = $14,
    successor_required_count = $15,
    role_summary = $16,
    metadata = $17,
    updated_by = $18,
    updated_at = now()
WHERE tenant_id = $1 AND id = $2 AND inactive = FALSE
RETURNING *;

-- name: UpdateSuccessionCriticalRoleStatus :one
UPDATE hrms.succession_critical_roles
SET status = $3, updated_by = $4, updated_at = now()
WHERE tenant_id = $1 AND id = $2 AND inactive = FALSE
RETURNING *;

-- name: GetSuccessionCriticalRole :one
SELECT cr.*, d.name AS department_name, ds.name AS designation_name, inc.display_name AS incumbent_name, inc.worker_code AS incumbent_code, cover.display_name AS emergency_cover_name, cover.worker_code AS emergency_cover_code
FROM hrms.succession_critical_roles cr
LEFT JOIN hrms.departments d ON d.tenant_id = cr.tenant_id AND d.id = cr.department_id
LEFT JOIN hrms.designations ds ON ds.tenant_id = cr.tenant_id AND ds.id = cr.designation_id
LEFT JOIN hrms.worker_profiles inc ON inc.tenant_id = cr.tenant_id AND inc.id = cr.incumbent_worker_profile_id
LEFT JOIN hrms.worker_profiles cover ON cover.tenant_id = cr.tenant_id AND cover.id = cr.emergency_cover_worker_profile_id
WHERE cr.tenant_id = $1 AND cr.id = $2 AND cr.inactive = FALSE;

-- name: ListSuccessionCriticalRoles :many
SELECT cr.*, d.name AS department_name, ds.name AS designation_name, inc.display_name AS incumbent_name, inc.worker_code AS incumbent_code, cover.display_name AS emergency_cover_name, cover.worker_code AS emergency_cover_code,
       COALESCE(COUNT(sn.id) FILTER (WHERE sn.inactive = FALSE AND sn.nomination_status NOT IN ('rejected','withdrawn')), 0)::bigint AS successor_count,
       COALESCE(COUNT(sn.id) FILTER (WHERE sn.inactive = FALSE AND sn.readiness_level = 'ready_now' AND sn.nomination_status NOT IN ('rejected','withdrawn')), 0)::bigint AS ready_now_count
FROM hrms.succession_critical_roles cr
LEFT JOIN hrms.departments d ON d.tenant_id = cr.tenant_id AND d.id = cr.department_id
LEFT JOIN hrms.designations ds ON ds.tenant_id = cr.tenant_id AND ds.id = cr.designation_id
LEFT JOIN hrms.worker_profiles inc ON inc.tenant_id = cr.tenant_id AND inc.id = cr.incumbent_worker_profile_id
LEFT JOIN hrms.worker_profiles cover ON cover.tenant_id = cr.tenant_id AND cover.id = cr.emergency_cover_worker_profile_id
LEFT JOIN hrms.succession_successor_nominations sn ON sn.tenant_id = cr.tenant_id AND sn.critical_role_id = cr.id
WHERE cr.tenant_id = $1
  AND cr.inactive = FALSE
  AND (sqlc.narg('cycle_id')::uuid IS NULL OR cr.cycle_id = sqlc.narg('cycle_id')::uuid)
  AND (sqlc.narg('status')::text IS NULL OR cr.status = sqlc.narg('status')::text)
  AND (sqlc.narg('search')::text IS NULL OR cr.code ILIKE '%' || sqlc.narg('search')::text || '%' OR cr.title ILIKE '%' || sqlc.narg('search')::text || '%')
GROUP BY cr.id, d.name, ds.name, inc.display_name, inc.worker_code, cover.display_name, cover.worker_code
ORDER BY cr.criticality DESC, cr.vacancy_risk DESC, cr.updated_at DESC
LIMIT $2 OFFSET $3;

-- name: SoftDeleteSuccessionCriticalRole :exec
UPDATE hrms.succession_critical_roles
SET inactive = TRUE, updated_by = $3, updated_at = now()
WHERE tenant_id = $1 AND id = $2;

-- name: CreateSuccessionSuccessorNomination :one
INSERT INTO hrms.succession_successor_nominations (
    tenant_id, critical_role_id, worker_profile_id, nominated_by, readiness_level, readiness_months,
    potential_rating, performance_rating, retention_risk, mobility_preference, nomination_status,
    development_notes, metadata, created_by, updated_by
)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$14)
RETURNING *;

-- name: UpdateSuccessionSuccessorNomination :one
UPDATE hrms.succession_successor_nominations
SET readiness_level = $4,
    readiness_months = $5,
    potential_rating = $6,
    performance_rating = $7,
    retention_risk = $8,
    mobility_preference = $9,
    development_notes = $10,
    metadata = $11,
    updated_by = $12,
    updated_at = now()
WHERE tenant_id = $1 AND id = $2 AND critical_role_id = $3 AND inactive = FALSE
RETURNING *;

-- name: UpdateSuccessionSuccessorNominationStatus :one
UPDATE hrms.succession_successor_nominations
SET nomination_status = $3, updated_by = $4, updated_at = now()
WHERE tenant_id = $1 AND id = $2 AND inactive = FALSE
RETURNING *;

-- name: ListSuccessionSuccessorNominations :many
SELECT sn.*, wp.display_name AS worker_display_name, wp.worker_code, cr.title AS critical_role_title
FROM hrms.succession_successor_nominations sn
JOIN hrms.worker_profiles wp ON wp.tenant_id = sn.tenant_id AND wp.id = sn.worker_profile_id
JOIN hrms.succession_critical_roles cr ON cr.tenant_id = sn.tenant_id AND cr.id = sn.critical_role_id
WHERE sn.tenant_id = $1
  AND sn.inactive = FALSE
  AND (sqlc.narg('critical_role_id')::uuid IS NULL OR sn.critical_role_id = sqlc.narg('critical_role_id')::uuid)
  AND (sqlc.narg('status')::text IS NULL OR sn.nomination_status = sqlc.narg('status')::text)
  AND (sqlc.narg('search')::text IS NULL OR wp.display_name ILIKE '%' || sqlc.narg('search')::text || '%' OR COALESCE(wp.worker_code, '') ILIKE '%' || sqlc.narg('search')::text || '%' OR cr.title ILIKE '%' || sqlc.narg('search')::text || '%')
ORDER BY sn.readiness_months ASC, sn.updated_at DESC
LIMIT $2 OFFSET $3;

-- name: CreateSuccessionDevelopmentAction :one
INSERT INTO hrms.succession_development_actions (
    tenant_id, nomination_id, critical_role_id, worker_profile_id, action_type, title,
    learning_course_id, learning_path_id, owner_user_id, due_date, status, notes, metadata, created_by, updated_by
)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$14)
RETURNING *;

-- name: UpdateSuccessionDevelopmentAction :one
UPDATE hrms.succession_development_actions
SET nomination_id = $3,
    critical_role_id = $4,
    worker_profile_id = $5,
    action_type = $6,
    title = $7,
    learning_course_id = $8,
    learning_path_id = $9,
    owner_user_id = $10,
    due_date = $11,
    notes = $12,
    metadata = $13,
    updated_by = $14,
    updated_at = now()
WHERE tenant_id = $1 AND id = $2 AND inactive = FALSE
RETURNING *;

-- name: UpdateSuccessionDevelopmentActionStatus :one
UPDATE hrms.succession_development_actions
SET status = $3, updated_by = $4, updated_at = now()
WHERE tenant_id = $1 AND id = $2 AND inactive = FALSE
RETURNING *;

-- name: ListSuccessionDevelopmentActions :many
SELECT da.*, wp.display_name AS worker_display_name, wp.worker_code, cr.title AS critical_role_title, lc.title AS learning_course_title, lp.title AS learning_path_title
FROM hrms.succession_development_actions da
JOIN hrms.worker_profiles wp ON wp.tenant_id = da.tenant_id AND wp.id = da.worker_profile_id
LEFT JOIN hrms.succession_critical_roles cr ON cr.tenant_id = da.tenant_id AND cr.id = da.critical_role_id
LEFT JOIN hrms.learning_courses lc ON lc.tenant_id = da.tenant_id AND lc.id = da.learning_course_id
LEFT JOIN hrms.learning_paths lp ON lp.tenant_id = da.tenant_id AND lp.id = da.learning_path_id
WHERE da.tenant_id = $1
  AND da.inactive = FALSE
  AND (sqlc.narg('critical_role_id')::uuid IS NULL OR da.critical_role_id = sqlc.narg('critical_role_id')::uuid)
  AND (sqlc.narg('status')::text IS NULL OR da.status = sqlc.narg('status')::text)
  AND (sqlc.narg('search')::text IS NULL OR da.title ILIKE '%' || sqlc.narg('search')::text || '%' OR wp.display_name ILIKE '%' || sqlc.narg('search')::text || '%')
ORDER BY da.due_date ASC NULLS LAST, da.updated_at DESC
LIMIT $2 OFFSET $3;

-- name: CreateSuccessionEvent :one
INSERT INTO hrms.succession_events (tenant_id, source_type, source_id, action, from_status, to_status, remarks, metadata, created_by)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
RETURNING *;

-- name: ListSuccessionEvents :many
SELECT * FROM hrms.succession_events
WHERE tenant_id = $1
  AND (sqlc.narg('source_type')::text IS NULL OR source_type = sqlc.narg('source_type')::text)
  AND (sqlc.narg('source_id')::uuid IS NULL OR source_id = sqlc.narg('source_id')::uuid)
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetSuccessionSummary :many
SELECT 'critical_roles'::text AS metric, COUNT(*)::bigint AS metric_count
FROM hrms.succession_critical_roles cr
WHERE cr.tenant_id = $1 AND cr.inactive = FALSE
UNION ALL
SELECT 'at_risk_roles', COUNT(*)::bigint
FROM hrms.succession_critical_roles cr
WHERE cr.tenant_id = $1 AND cr.inactive = FALSE AND (cr.status = 'at_risk' OR cr.criticality IN ('high','critical') OR cr.vacancy_risk IN ('high','critical'))
UNION ALL
SELECT 'ready_now_successors', COUNT(*)::bigint
FROM hrms.succession_successor_nominations sn
WHERE sn.tenant_id = $1 AND sn.inactive = FALSE AND sn.readiness_level = 'ready_now' AND sn.nomination_status NOT IN ('rejected','withdrawn')
UNION ALL
SELECT 'open_development_actions', COUNT(*)::bigint
FROM hrms.succession_development_actions da
WHERE da.tenant_id = $1 AND da.inactive = FALSE AND da.status IN ('open','in_progress','overdue')
UNION ALL
SELECT 'roles_without_successor', COUNT(*)::bigint
FROM hrms.succession_critical_roles cr
WHERE cr.tenant_id = $1 AND cr.inactive = FALSE AND NOT EXISTS (
    SELECT 1 FROM hrms.succession_successor_nominations sn
    WHERE sn.tenant_id = cr.tenant_id AND sn.critical_role_id = cr.id AND sn.inactive = FALSE AND sn.nomination_status NOT IN ('rejected','withdrawn')
);
