-- name: CreateSkillCategory :one
INSERT INTO hrms.skill_categories (
    tenant_id, parent_id, code, name, description, source_scope, sort_order,
    metadata, created_by, updated_by
) VALUES (
    $1, $2, $3, $4, $5, 'tenant', $6,
    $7, $8, $8
)
RETURNING *;

-- name: UpdateSkillCategory :one
UPDATE hrms.skill_categories
SET parent_id = $3,
    code = $4,
    name = $5,
    description = $6,
    sort_order = $7,
    metadata = $8,
    updated_by = $9
WHERE tenant_id = $1 AND id = $2 AND source_scope = 'tenant' AND NOT inactive
RETURNING *;

-- name: GetSkillCategory :one
SELECT * FROM hrms.skill_categories
WHERE (tenant_id = $1 OR tenant_id IS NULL) AND id = $2 AND NOT inactive;

-- name: ListSkillCategories :many
SELECT
    sc.*,
    parent.name AS parent_name
FROM hrms.skill_categories sc
LEFT JOIN hrms.skill_categories parent ON parent.id = sc.parent_id AND NOT parent.inactive
WHERE (sc.tenant_id = $1 OR sc.tenant_id IS NULL)
  AND (sqlc.narg('source_scope')::text IS NULL OR sc.source_scope = sqlc.narg('source_scope')::text)
  AND (sqlc.narg('parent_id')::uuid IS NULL OR sc.parent_id = sqlc.narg('parent_id')::uuid)
  AND (
      sqlc.narg('search')::text IS NULL
      OR sc.code ILIKE '%' || sqlc.narg('search')::text || '%'
      OR sc.name ILIKE '%' || sqlc.narg('search')::text || '%'
      OR sc.description ILIKE '%' || sqlc.narg('search')::text || '%'
  )
  AND NOT sc.inactive
ORDER BY sc.sort_order ASC, sc.name ASC;

-- name: SoftDeleteSkillCategory :exec
UPDATE hrms.skill_categories
SET inactive = TRUE,
    updated_by = $3
WHERE tenant_id = $1 AND id = $2 AND source_scope = 'tenant' AND NOT inactive;

-- name: CreateSkill :one
INSERT INTO hrms.skills (
    tenant_id, category_id, code, name, description, skill_type, source_scope,
    certificate_required, assessment_required, is_active, metadata, created_by, updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6, 'tenant',
    $7, $8, $9, $10, $11, $11
)
RETURNING *;

-- name: UpdateSkill :one
UPDATE hrms.skills
SET category_id = $3,
    code = $4,
    name = $5,
    description = $6,
    skill_type = $7,
    certificate_required = $8,
    assessment_required = $9,
    is_active = $10,
    metadata = $11,
    updated_by = $12
WHERE tenant_id = $1 AND id = $2 AND source_scope = 'tenant' AND NOT inactive
RETURNING *;

-- name: GetSkill :one
SELECT * FROM hrms.skills
WHERE (tenant_id = $1 OR tenant_id IS NULL) AND id = $2 AND NOT inactive;

-- name: GetSkillByCode :one
SELECT * FROM hrms.skills
WHERE (tenant_id = $1 OR tenant_id IS NULL) AND lower(code) = lower($2) AND NOT inactive
ORDER BY CASE WHEN tenant_id = $1 THEN 0 ELSE 1 END
LIMIT 1;

-- name: ListSkills :many
SELECT
    s.*,
    sc.name AS category_name,
    sc.code AS category_code
FROM hrms.skills s
LEFT JOIN hrms.skill_categories sc ON sc.id = s.category_id AND NOT sc.inactive
WHERE (s.tenant_id = $1 OR s.tenant_id IS NULL)
  AND (sqlc.narg('category_id')::uuid IS NULL OR s.category_id = sqlc.narg('category_id')::uuid)
  AND (sqlc.narg('skill_type')::text IS NULL OR s.skill_type = sqlc.narg('skill_type')::text)
  AND (sqlc.narg('source_scope')::text IS NULL OR s.source_scope = sqlc.narg('source_scope')::text)
  AND (sqlc.narg('is_active')::boolean IS NULL OR s.is_active = sqlc.narg('is_active')::boolean)
  AND (
      sqlc.narg('search')::text IS NULL
      OR s.code ILIKE '%' || sqlc.narg('search')::text || '%'
      OR s.name ILIKE '%' || sqlc.narg('search')::text || '%'
      OR s.description ILIKE '%' || sqlc.narg('search')::text || '%'
      OR sc.name ILIKE '%' || sqlc.narg('search')::text || '%'
  )
  AND NOT s.inactive
ORDER BY s.source_scope ASC, sc.sort_order ASC NULLS LAST, s.name ASC;

-- name: SoftDeleteSkill :exec
UPDATE hrms.skills
SET inactive = TRUE,
    updated_by = $3
WHERE tenant_id = $1 AND id = $2 AND source_scope = 'tenant' AND NOT inactive;

-- name: CreateWorkerSkill :one
INSERT INTO hrms.worker_skills (
    tenant_id, worker_profile_id, skill_id, skill_name_snapshot, proficiency,
    years_experience, last_used_on, verification_status, certificate_url,
    certificate_expires_on, assessment_score, assessed_on, notes, metadata,
    created_by, updated_by
) VALUES (
    $1, $2, $3, $4, $5,
    $6, $7, $8, $9,
    $10, $11, $12, $13, $14,
    $15, $15
)
RETURNING *;

-- name: UpdateWorkerSkill :one
UPDATE hrms.worker_skills
SET skill_id = $3,
    skill_name_snapshot = $4,
    proficiency = $5,
    years_experience = $6,
    last_used_on = $7,
    verification_status = $8,
    certificate_url = $9,
    certificate_expires_on = $10,
    assessment_score = $11,
    assessed_on = $12,
    notes = $13,
    metadata = $14,
    updated_by = $15
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: GetWorkerSkill :one
SELECT * FROM hrms.worker_skills
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: ListWorkerSkills :many
SELECT
    ws.*,
    wp.display_name AS worker_display_name,
    wp.worker_code,
    s.code AS skill_code,
    s.name AS skill_name,
    s.skill_type,
    s.source_scope AS skill_source_scope,
    s.certificate_required,
    s.assessment_required,
    sc.name AS category_name,
    sc.code AS category_code
FROM hrms.worker_skills ws
JOIN hrms.worker_profiles wp ON wp.tenant_id = ws.tenant_id AND wp.id = ws.worker_profile_id AND NOT wp.inactive
JOIN hrms.skills s ON s.id = ws.skill_id AND NOT s.inactive
LEFT JOIN hrms.skill_categories sc ON sc.id = s.category_id AND NOT sc.inactive
WHERE ws.tenant_id = $1
  AND (sqlc.narg('worker_profile_id')::uuid IS NULL OR ws.worker_profile_id = sqlc.narg('worker_profile_id')::uuid)
  AND (sqlc.narg('skill_id')::uuid IS NULL OR ws.skill_id = sqlc.narg('skill_id')::uuid)
  AND (sqlc.narg('category_id')::uuid IS NULL OR s.category_id = sqlc.narg('category_id')::uuid)
  AND (sqlc.narg('proficiency')::text IS NULL OR ws.proficiency = sqlc.narg('proficiency')::text)
  AND (sqlc.narg('verification_status')::text IS NULL OR ws.verification_status = sqlc.narg('verification_status')::text)
  AND (sqlc.narg('certificate_expiring_before')::date IS NULL OR ws.certificate_expires_on <= sqlc.narg('certificate_expiring_before')::date)
  AND (
      sqlc.narg('search')::text IS NULL
      OR wp.display_name ILIKE '%' || sqlc.narg('search')::text || '%'
      OR wp.worker_code ILIKE '%' || sqlc.narg('search')::text || '%'
      OR s.code ILIKE '%' || sqlc.narg('search')::text || '%'
      OR s.name ILIKE '%' || sqlc.narg('search')::text || '%'
      OR ws.skill_name_snapshot ILIKE '%' || sqlc.narg('search')::text || '%'
  )
  AND NOT ws.inactive
ORDER BY
    CASE ws.verification_status
        WHEN 'hr_verified' THEN 1
        WHEN 'manager_endorsed' THEN 2
        WHEN 'self_declared' THEN 3
        WHEN 'expired' THEN 4
        ELSE 5
    END,
    wp.display_name ASC,
    s.name ASC;

-- name: UpdateWorkerSkillVerification :one
UPDATE hrms.worker_skills
SET verification_status = $3,
    endorsed_by = CASE WHEN $3 = 'manager_endorsed' THEN $4 ELSE endorsed_by END,
    endorsed_at = CASE WHEN $3 = 'manager_endorsed' THEN now() ELSE endorsed_at END,
    verified_by = CASE WHEN $3 IN ('hr_verified', 'rejected', 'expired') THEN $4 ELSE verified_by END,
    verified_at = CASE WHEN $3 IN ('hr_verified', 'rejected', 'expired') THEN now() ELSE verified_at END,
    notes = COALESCE($5, notes),
    updated_by = $4
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: SoftDeleteWorkerSkill :exec
UPDATE hrms.worker_skills
SET inactive = TRUE,
    updated_by = $3
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: CreateWorkerSkillAssessment :one
INSERT INTO hrms.worker_skill_assessments (
    tenant_id, worker_skill_id, assessment_type, result_status, score, max_score,
    assessed_by, assessed_on, evidence_url, notes, metadata, created_by
) VALUES (
    $1, $2, $3, $4, $5, $6,
    $7, $8, $9, $10, $11, $12
)
RETURNING *;

-- name: ListWorkerSkillAssessments :many
SELECT * FROM hrms.worker_skill_assessments
WHERE tenant_id = $1
  AND (sqlc.narg('worker_skill_id')::uuid IS NULL OR worker_skill_id = sqlc.narg('worker_skill_id')::uuid)
ORDER BY assessed_on DESC, created_at DESC;

-- name: GetSkillsSummary :many
SELECT
    COALESCE(ws.verification_status, 'catalog') AS status,
    COUNT(DISTINCT ws.id)::int AS worker_skill_count,
    COUNT(DISTINCT ws.worker_profile_id)::int AS worker_count,
    COUNT(DISTINCT s.id)::int AS skill_count,
    COUNT(DISTINCT ws.id) FILTER (
        WHERE ws.certificate_expires_on IS NOT NULL
          AND ws.certificate_expires_on <= CURRENT_DATE + INTERVAL '30 days'
    )::int AS expiring_certificate_count
FROM hrms.skills s
LEFT JOIN hrms.worker_skills ws ON ws.skill_id = s.id AND ws.tenant_id = $1 AND NOT ws.inactive
WHERE (s.tenant_id = $1 OR s.tenant_id IS NULL) AND NOT s.inactive
GROUP BY ROLLUP(ws.verification_status)
ORDER BY status ASC;
