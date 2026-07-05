-- name: CreateLearningCourse :one
INSERT INTO hrms.learning_courses (
    tenant_id, code, title, description, course_type, delivery_mode, provider,
    duration_minutes, skill_id, compliance_rule_id, mandatory, ai_readiness,
    certificate_required, budget_amount, currency_code, is_active, metadata,
    created_by, updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7,
    $8, $9, $10, $11, $12,
    $13, $14, $15, $16, $17,
    $18, $18
)
RETURNING *;

-- name: UpdateLearningCourse :one
UPDATE hrms.learning_courses
SET code = $3,
    title = $4,
    description = $5,
    course_type = $6,
    delivery_mode = $7,
    provider = $8,
    duration_minutes = $9,
    skill_id = $10,
    compliance_rule_id = $11,
    mandatory = $12,
    ai_readiness = $13,
    certificate_required = $14,
    budget_amount = $15,
    currency_code = $16,
    is_active = $17,
    metadata = $18,
    updated_by = $19
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: ListLearningCourses :many
SELECT
    c.*,
    s.name AS skill_name,
    s.code AS skill_code,
    cr.title AS compliance_rule_title
FROM hrms.learning_courses c
LEFT JOIN hrms.skills s ON s.id = c.skill_id AND NOT s.inactive
LEFT JOIN hrms.compliance_rules cr ON cr.tenant_id = c.tenant_id AND cr.id = c.compliance_rule_id AND NOT cr.inactive
WHERE c.tenant_id = $1
  AND (sqlc.narg('course_type')::text IS NULL OR c.course_type = sqlc.narg('course_type')::text)
  AND (sqlc.narg('skill_id')::uuid IS NULL OR c.skill_id = sqlc.narg('skill_id')::uuid)
  AND (sqlc.narg('mandatory')::boolean IS NULL OR c.mandatory = sqlc.narg('mandatory')::boolean)
  AND (sqlc.narg('ai_readiness')::boolean IS NULL OR c.ai_readiness = sqlc.narg('ai_readiness')::boolean)
  AND (sqlc.narg('is_active')::boolean IS NULL OR c.is_active = sqlc.narg('is_active')::boolean)
  AND (
      sqlc.narg('search')::text IS NULL
      OR c.code ILIKE '%' || sqlc.narg('search')::text || '%'
      OR c.title ILIKE '%' || sqlc.narg('search')::text || '%'
      OR c.description ILIKE '%' || sqlc.narg('search')::text || '%'
      OR s.name ILIKE '%' || sqlc.narg('search')::text || '%'
  )
  AND NOT c.inactive
ORDER BY c.mandatory DESC, c.ai_readiness DESC, c.title ASC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: CountLearningCourses :one
SELECT COUNT(*)
FROM hrms.learning_courses c
LEFT JOIN hrms.skills s ON s.id = c.skill_id AND NOT s.inactive
WHERE c.tenant_id = $1
  AND (sqlc.narg('course_type')::text IS NULL OR c.course_type = sqlc.narg('course_type')::text)
  AND (sqlc.narg('skill_id')::uuid IS NULL OR c.skill_id = sqlc.narg('skill_id')::uuid)
  AND (sqlc.narg('mandatory')::boolean IS NULL OR c.mandatory = sqlc.narg('mandatory')::boolean)
  AND (sqlc.narg('ai_readiness')::boolean IS NULL OR c.ai_readiness = sqlc.narg('ai_readiness')::boolean)
  AND (sqlc.narg('is_active')::boolean IS NULL OR c.is_active = sqlc.narg('is_active')::boolean)
  AND (
      sqlc.narg('search')::text IS NULL
      OR c.code ILIKE '%' || sqlc.narg('search')::text || '%'
      OR c.title ILIKE '%' || sqlc.narg('search')::text || '%'
      OR c.description ILIKE '%' || sqlc.narg('search')::text || '%'
      OR s.name ILIKE '%' || sqlc.narg('search')::text || '%'
  )
  AND NOT c.inactive;

-- name: GetLearningCourse :one
SELECT * FROM hrms.learning_courses
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: SoftDeleteLearningCourse :exec
UPDATE hrms.learning_courses
SET inactive = TRUE,
    is_active = FALSE,
    updated_by = $3
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: CreateLearningPath :one
INSERT INTO hrms.learning_paths (
    tenant_id, code, title, description, path_type, target_role, skill_id,
    ai_readiness, is_active, metadata, created_by, updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7,
    $8, $9, $10, $11, $11
)
RETURNING *;

-- name: UpdateLearningPath :one
UPDATE hrms.learning_paths
SET code = $3,
    title = $4,
    description = $5,
    path_type = $6,
    target_role = $7,
    skill_id = $8,
    ai_readiness = $9,
    is_active = $10,
    metadata = $11,
    updated_by = $12
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: ListLearningPaths :many
SELECT
    p.*,
    s.name AS skill_name,
    COUNT(pc.id)::int AS course_count,
    COALESCE(SUM(c.duration_minutes) FILTER (WHERE NOT c.inactive), 0)::int AS total_minutes
FROM hrms.learning_paths p
LEFT JOIN hrms.skills s ON s.id = p.skill_id AND NOT s.inactive
LEFT JOIN hrms.learning_path_courses pc ON pc.tenant_id = p.tenant_id AND pc.path_id = p.id AND NOT pc.inactive
LEFT JOIN hrms.learning_courses c ON c.tenant_id = p.tenant_id AND c.id = pc.course_id AND NOT c.inactive
WHERE p.tenant_id = $1
  AND (sqlc.narg('path_type')::text IS NULL OR p.path_type = sqlc.narg('path_type')::text)
  AND (sqlc.narg('skill_id')::uuid IS NULL OR p.skill_id = sqlc.narg('skill_id')::uuid)
  AND (sqlc.narg('ai_readiness')::boolean IS NULL OR p.ai_readiness = sqlc.narg('ai_readiness')::boolean)
  AND (sqlc.narg('is_active')::boolean IS NULL OR p.is_active = sqlc.narg('is_active')::boolean)
  AND (
      sqlc.narg('search')::text IS NULL
      OR p.code ILIKE '%' || sqlc.narg('search')::text || '%'
      OR p.title ILIKE '%' || sqlc.narg('search')::text || '%'
      OR p.description ILIKE '%' || sqlc.narg('search')::text || '%'
  )
  AND NOT p.inactive
GROUP BY p.id, s.name
ORDER BY p.ai_readiness DESC, p.title ASC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: GetLearningPath :one
SELECT * FROM hrms.learning_paths
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: SoftDeleteLearningPath :exec
UPDATE hrms.learning_paths
SET inactive = TRUE,
    is_active = FALSE,
    updated_by = $3
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: UpsertLearningPathCourse :one
INSERT INTO hrms.learning_path_courses (
    tenant_id, path_id, course_id, sort_order, required, created_by, updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $6
)
ON CONFLICT (tenant_id, path_id, course_id) WHERE inactive = FALSE
DO UPDATE SET sort_order = EXCLUDED.sort_order,
              required = EXCLUDED.required,
              updated_by = EXCLUDED.updated_by,
              updated_at = NOW()
RETURNING *;

-- name: ListLearningPathCourses :many
SELECT
    pc.*,
    c.code AS course_code,
    c.title AS course_title,
    c.course_type,
    c.delivery_mode,
    c.duration_minutes
FROM hrms.learning_path_courses pc
JOIN hrms.learning_courses c ON c.tenant_id = pc.tenant_id AND c.id = pc.course_id AND NOT c.inactive
WHERE pc.tenant_id = $1
  AND pc.path_id = $2
  AND NOT pc.inactive
ORDER BY pc.sort_order ASC, c.title ASC;

-- name: SoftDeleteLearningPathCourse :exec
UPDATE hrms.learning_path_courses
SET inactive = TRUE,
    updated_by = $3
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: CreateLearningEnrollment :one
INSERT INTO hrms.learning_enrollments (
    tenant_id, course_id, path_id, worker_profile_id, assignment_source, status,
    nominated_by, assigned_by, due_date, notes, metadata, created_by, updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6,
    $7, $8, $9, $10, $11, $12, $12
)
RETURNING *;

-- name: UpdateLearningEnrollmentStatus :one
UPDATE hrms.learning_enrollments
SET status = $3,
    started_at = CASE WHEN $3 = 'in_progress' AND started_at IS NULL THEN NOW() ELSE started_at END,
    completed_at = CASE WHEN $3 = 'completed' THEN NOW() ELSE completed_at END,
    score = $4,
    certificate_url = COALESCE($5, certificate_url),
    certificate_file_name = COALESCE($6, certificate_file_name),
    certificate_content_type = COALESCE($7, certificate_content_type),
    notes = $8,
    updated_by = $9
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: ListLearningEnrollments :many
SELECT
    e.*,
    c.title AS course_title,
    c.code AS course_code,
    c.course_type,
    c.ai_readiness,
    c.mandatory,
    p.title AS path_title,
    wp.display_name AS worker_display_name,
    wp.worker_code
FROM hrms.learning_enrollments e
JOIN hrms.learning_courses c ON c.tenant_id = e.tenant_id AND c.id = e.course_id AND NOT c.inactive
LEFT JOIN hrms.learning_paths p ON p.tenant_id = e.tenant_id AND p.id = e.path_id AND NOT p.inactive
JOIN hrms.worker_profiles wp ON wp.tenant_id = e.tenant_id AND wp.id = e.worker_profile_id AND NOT wp.inactive
WHERE e.tenant_id = $1
  AND (sqlc.narg('worker_profile_id')::uuid IS NULL OR e.worker_profile_id = sqlc.narg('worker_profile_id')::uuid)
  AND (sqlc.narg('course_id')::uuid IS NULL OR e.course_id = sqlc.narg('course_id')::uuid)
  AND (sqlc.narg('status')::text IS NULL OR e.status = sqlc.narg('status')::text)
  AND (sqlc.narg('assignment_source')::text IS NULL OR e.assignment_source = sqlc.narg('assignment_source')::text)
  AND (
      sqlc.narg('search')::text IS NULL
      OR c.title ILIKE '%' || sqlc.narg('search')::text || '%'
      OR wp.display_name ILIKE '%' || sqlc.narg('search')::text || '%'
      OR wp.worker_code ILIKE '%' || sqlc.narg('search')::text || '%'
  )
  AND NOT e.inactive
ORDER BY CASE e.status WHEN 'overdue' THEN 1 WHEN 'assigned' THEN 2 WHEN 'nominated' THEN 3 WHEN 'in_progress' THEN 4 ELSE 5 END,
         e.due_date ASC NULLS LAST,
         e.updated_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: GetLearningEnrollment :one
SELECT * FROM hrms.learning_enrollments
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: SoftDeleteLearningEnrollment :exec
UPDATE hrms.learning_enrollments
SET inactive = TRUE,
    status = 'cancelled',
    updated_by = $3
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: CreateLearningRecommendation :one
INSERT INTO hrms.learning_recommendations (
    tenant_id, worker_profile_id, skill_id, course_id, path_id, source_type,
    reason, priority, confidence_score, status, metadata, created_by, updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6,
    $7, $8, $9, $10, $11, $12, $12
)
RETURNING *;

-- name: UpdateLearningRecommendationStatus :one
UPDATE hrms.learning_recommendations
SET status = $3,
    updated_by = $4
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: ListLearningRecommendations :many
SELECT
    r.*,
    wp.display_name AS worker_display_name,
    wp.worker_code,
    s.name AS skill_name,
    c.title AS course_title,
    p.title AS path_title
FROM hrms.learning_recommendations r
LEFT JOIN hrms.worker_profiles wp ON wp.tenant_id = r.tenant_id AND wp.id = r.worker_profile_id AND NOT wp.inactive
LEFT JOIN hrms.skills s ON s.id = r.skill_id AND NOT s.inactive
LEFT JOIN hrms.learning_courses c ON c.tenant_id = r.tenant_id AND c.id = r.course_id AND NOT c.inactive
LEFT JOIN hrms.learning_paths p ON p.tenant_id = r.tenant_id AND p.id = r.path_id AND NOT p.inactive
WHERE r.tenant_id = $1
  AND (sqlc.narg('worker_profile_id')::uuid IS NULL OR r.worker_profile_id = sqlc.narg('worker_profile_id')::uuid)
  AND (sqlc.narg('skill_id')::uuid IS NULL OR r.skill_id = sqlc.narg('skill_id')::uuid)
  AND (sqlc.narg('source_type')::text IS NULL OR r.source_type = sqlc.narg('source_type')::text)
  AND (sqlc.narg('status')::text IS NULL OR r.status = sqlc.narg('status')::text)
  AND (
      sqlc.narg('search')::text IS NULL
      OR r.reason ILIKE '%' || sqlc.narg('search')::text || '%'
      OR wp.display_name ILIKE '%' || sqlc.narg('search')::text || '%'
      OR s.name ILIKE '%' || sqlc.narg('search')::text || '%'
      OR c.title ILIKE '%' || sqlc.narg('search')::text || '%'
  )
  AND NOT r.inactive
ORDER BY CASE r.priority WHEN 'urgent' THEN 1 WHEN 'high' THEN 2 WHEN 'medium' THEN 3 ELSE 4 END,
         r.created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: GenerateSkillGapLearningRecommendations :many
WITH missing AS (
    SELECT DISTINCT
        psr.tenant_id,
        psr.skill_id,
        s.name AS skill_name,
        c.id AS course_id,
        p.id AS path_id,
        psr.importance
    FROM hrms.project_skill_requirements psr
    JOIN hrms.skills s ON s.id = psr.skill_id AND NOT s.inactive
    LEFT JOIN hrms.learning_courses c ON c.tenant_id = psr.tenant_id AND c.skill_id = psr.skill_id AND c.is_active AND NOT c.inactive
    LEFT JOIN hrms.learning_paths p ON p.tenant_id = psr.tenant_id AND p.skill_id = psr.skill_id AND p.is_active AND NOT p.inactive
    WHERE psr.tenant_id = $1
      AND NOT psr.inactive
      AND (c.id IS NOT NULL OR p.id IS NOT NULL)
      AND NOT EXISTS (
          SELECT 1
          FROM hrms.worker_skills ws
          WHERE ws.tenant_id = psr.tenant_id
            AND ws.skill_id = psr.skill_id
            AND ws.verification_status IN ('manager_endorsed','hr_verified')
            AND NOT ws.inactive
      )
)
INSERT INTO hrms.learning_recommendations (
    tenant_id, skill_id, course_id, path_id, source_type, reason, priority,
    confidence_score, status, metadata, created_by, updated_by
)
SELECT
    m.tenant_id,
    m.skill_id,
    m.course_id,
    m.path_id,
    'skill_gap',
    ('Skill gap detected for ' || m.skill_name || '. Assign matching learning to build coverage.')::text,
    CASE m.importance WHEN 'critical' THEN 'urgent' WHEN 'high' THEN 'high' ELSE 'medium' END,
    82,
    'open',
    jsonb_build_object('importance', m.importance, 'generated', true),
    $2,
    $2
FROM missing m
WHERE NOT EXISTS (
    SELECT 1
    FROM hrms.learning_recommendations existing
    WHERE existing.tenant_id = m.tenant_id
      AND existing.skill_id = m.skill_id
      AND COALESCE(existing.course_id, '00000000-0000-0000-0000-000000000000'::uuid) = COALESCE(m.course_id, '00000000-0000-0000-0000-000000000000'::uuid)
      AND COALESCE(existing.path_id, '00000000-0000-0000-0000-000000000000'::uuid) = COALESCE(m.path_id, '00000000-0000-0000-0000-000000000000'::uuid)
      AND existing.source_type = 'skill_gap'
      AND existing.status IN ('open','accepted','assigned')
      AND NOT existing.inactive
)
RETURNING *;

-- name: GetLearningSummary :many
SELECT 'courses'::text AS metric, COUNT(*)::int AS metric_count
FROM hrms.learning_courses lc
WHERE lc.tenant_id = $1 AND NOT lc.inactive AND lc.is_active
UNION ALL
SELECT 'mandatory_courses'::text, COUNT(*)::int
FROM hrms.learning_courses lc
WHERE lc.tenant_id = $1 AND NOT lc.inactive AND lc.is_active AND lc.mandatory
UNION ALL
SELECT 'ai_readiness_courses'::text, COUNT(*)::int
FROM hrms.learning_courses lc
WHERE lc.tenant_id = $1 AND NOT lc.inactive AND lc.is_active AND lc.ai_readiness
UNION ALL
SELECT 'assigned'::text, COUNT(*)::int
FROM hrms.learning_enrollments le
WHERE le.tenant_id = $1 AND NOT le.inactive AND le.status IN ('assigned','nominated','approved','in_progress')
UNION ALL
SELECT 'completed'::text, COUNT(*)::int
FROM hrms.learning_enrollments le
WHERE le.tenant_id = $1 AND NOT le.inactive AND le.status = 'completed'
UNION ALL
SELECT 'overdue'::text, COUNT(*)::int
FROM hrms.learning_enrollments le
WHERE le.tenant_id = $1 AND NOT le.inactive AND le.due_date IS NOT NULL AND le.due_date < CURRENT_DATE AND le.status NOT IN ('completed','waived','cancelled')
UNION ALL
SELECT 'recommendations'::text, COUNT(*)::int
FROM hrms.learning_recommendations lr
WHERE lr.tenant_id = $1 AND NOT lr.inactive AND lr.status = 'open';
