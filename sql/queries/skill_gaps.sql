-- name: CreateProjectSkillRequirement :one
INSERT INTO hrms.project_skill_requirements (
    tenant_id, project_id, engagement_id, skill_id, required_proficiency,
    min_years_experience, required_count, importance, requirement_source,
    notes, metadata, created_by, updated_by
) VALUES (
    $1, $2, $3, $4, $5,
    $6, $7, $8, $9,
    $10, $11, $12, $12
)
RETURNING *;

-- name: UpdateProjectSkillRequirement :one
UPDATE hrms.project_skill_requirements
SET project_id = $3,
    engagement_id = $4,
    skill_id = $5,
    required_proficiency = $6,
    min_years_experience = $7,
    required_count = $8,
    importance = $9,
    requirement_source = $10,
    notes = $11,
    metadata = $12,
    updated_by = $13
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: GetProjectSkillRequirement :one
SELECT * FROM hrms.project_skill_requirements
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: ListProjectSkillRequirements :many
SELECT
    psr.*,
    p.name AS project_name,
    p.project_code,
    e.title AS engagement_title,
    e.engagement_code,
    e.worker_profile_id,
    wp.display_name AS worker_display_name,
    wp.worker_code,
    s.name AS skill_name,
    s.code AS skill_code,
    s.skill_type,
    sc.name AS category_name
FROM hrms.project_skill_requirements psr
LEFT JOIN hrms.projects p ON p.tenant_id = psr.tenant_id AND p.id = psr.project_id AND NOT p.inactive
LEFT JOIN hrms.engagements e ON e.tenant_id = psr.tenant_id AND e.id = psr.engagement_id AND NOT e.inactive
LEFT JOIN hrms.worker_profiles wp ON wp.tenant_id = psr.tenant_id AND wp.id = e.worker_profile_id AND NOT wp.inactive
JOIN hrms.skills s ON s.id = psr.skill_id AND NOT s.inactive
LEFT JOIN hrms.skill_categories sc ON sc.id = s.category_id AND NOT sc.inactive
WHERE psr.tenant_id = $1
  AND (sqlc.narg('project_id')::uuid IS NULL OR psr.project_id = sqlc.narg('project_id')::uuid)
  AND (sqlc.narg('engagement_id')::uuid IS NULL OR psr.engagement_id = sqlc.narg('engagement_id')::uuid)
  AND (sqlc.narg('skill_id')::uuid IS NULL OR psr.skill_id = sqlc.narg('skill_id')::uuid)
  AND (sqlc.narg('importance')::text IS NULL OR psr.importance = sqlc.narg('importance')::text)
  AND (
      sqlc.narg('search')::text IS NULL
      OR p.name ILIKE '%' || sqlc.narg('search')::text || '%'
      OR p.project_code ILIKE '%' || sqlc.narg('search')::text || '%'
      OR e.title ILIKE '%' || sqlc.narg('search')::text || '%'
      OR e.engagement_code ILIKE '%' || sqlc.narg('search')::text || '%'
      OR s.name ILIKE '%' || sqlc.narg('search')::text || '%'
      OR s.code ILIKE '%' || sqlc.narg('search')::text || '%'
  )
  AND NOT psr.inactive
ORDER BY
    CASE psr.importance WHEN 'critical' THEN 1 WHEN 'required' THEN 2 ELSE 3 END,
    p.name ASC NULLS LAST,
    e.title ASC NULLS LAST,
    s.name ASC;

-- name: SoftDeleteProjectSkillRequirement :exec
UPDATE hrms.project_skill_requirements
SET inactive = TRUE,
    updated_by = $3
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: ListProjectSkillGapRows :many
WITH req AS (
    SELECT
        psr.*,
        p.name AS project_name,
        p.project_code,
        e.title AS engagement_title,
        e.engagement_code,
        s.name AS skill_name,
        s.code AS skill_code,
        s.skill_type,
        CASE psr.required_proficiency
            WHEN 'beginner' THEN 1
            WHEN 'intermediate' THEN 2
            WHEN 'advanced' THEN 3
            WHEN 'expert' THEN 4
            ELSE 1
        END AS required_rank
    FROM hrms.project_skill_requirements psr
    LEFT JOIN hrms.projects p ON p.tenant_id = psr.tenant_id AND p.id = psr.project_id AND NOT p.inactive
    LEFT JOIN hrms.engagements e ON e.tenant_id = psr.tenant_id AND e.id = psr.engagement_id AND NOT e.inactive
    JOIN hrms.skills s ON s.id = psr.skill_id AND NOT s.inactive
    WHERE psr.tenant_id = $1
      AND (sqlc.narg('project_id')::uuid IS NULL OR psr.project_id = sqlc.narg('project_id')::uuid)
      AND (sqlc.narg('engagement_id')::uuid IS NULL OR psr.engagement_id = sqlc.narg('engagement_id')::uuid)
      AND (sqlc.narg('importance')::text IS NULL OR psr.importance = sqlc.narg('importance')::text)
      AND NOT psr.inactive
)
SELECT
    req.id AS requirement_id,
    req.tenant_id,
    req.project_id,
    req.project_name,
    req.project_code,
    req.engagement_id,
    req.engagement_title,
    req.engagement_code,
    req.skill_id,
    req.skill_name,
    req.skill_code,
    req.skill_type,
    req.required_proficiency,
    req.min_years_experience,
    req.required_count,
    req.importance,
    assigned.assigned_match_count::int AS assigned_match_count,
    available.available_match_count::int AS available_match_count,
    GREATEST(req.required_count - assigned.assigned_match_count, 0)::int AS gap_count,
    LEAST(100, ROUND((assigned.assigned_match_count::numeric / req.required_count::numeric) * 100))::int AS match_percent,
    ((assigned.assigned_match_count = 1 AND req.required_count > 1) OR (available.available_match_count = 1 AND req.importance IN ('required', 'critical'))) AS single_person_dependency,
    CASE
        WHEN assigned.assigned_match_count >= req.required_count THEN 'covered'
        WHEN available.available_match_count >= (req.required_count - assigned.assigned_match_count) THEN 'train_or_assign'
        WHEN req.importance = 'critical' THEN 'hire_or_contract'
        ELSE 'train'
    END AS suggested_action
FROM req
CROSS JOIN LATERAL (
    SELECT COUNT(DISTINCT ws.worker_profile_id) AS assigned_match_count
    FROM hrms.worker_skills ws
    WHERE ws.tenant_id = req.tenant_id
      AND ws.skill_id = req.skill_id
      AND ws.verification_status IN ('manager_endorsed', 'hr_verified')
      AND NOT ws.inactive
      AND (req.min_years_experience IS NULL OR COALESCE(ws.years_experience, 0) >= req.min_years_experience)
      AND (
          CASE ws.proficiency
              WHEN 'beginner' THEN 1
              WHEN 'intermediate' THEN 2
              WHEN 'advanced' THEN 3
              WHEN 'expert' THEN 4
              ELSE 0
          END
      ) >= req.required_rank
      AND EXISTS (
          SELECT 1
          FROM hrms.engagements eg
          LEFT JOIN hrms.project_milestones pm ON pm.tenant_id = eg.tenant_id AND pm.engagement_id = eg.id AND NOT pm.inactive
          LEFT JOIN hrms.projects p2 ON p2.tenant_id = eg.tenant_id AND p2.id = req.project_id AND NOT p2.inactive
          WHERE eg.tenant_id = req.tenant_id
            AND eg.worker_profile_id = ws.worker_profile_id
            AND eg.status IN ('draft', 'active', 'paused')
            AND NOT eg.inactive
            AND (
                (req.engagement_id IS NOT NULL AND eg.id = req.engagement_id)
                OR (
                    req.engagement_id IS NULL
                    AND req.project_id IS NOT NULL
                    AND (
                        pm.project_id = req.project_id
                        OR (p2.project_code IS NOT NULL AND eg.project_code = p2.project_code)
                        OR (eg.project_label IS NOT NULL AND eg.project_label = p2.name)
                    )
                )
            )
      )
) assigned
CROSS JOIN LATERAL (
    SELECT COUNT(DISTINCT ws.worker_profile_id) AS available_match_count
    FROM hrms.worker_skills ws
    JOIN hrms.worker_profiles wp ON wp.tenant_id = ws.tenant_id AND wp.id = ws.worker_profile_id AND wp.profile_status IN ('active', 'draft', 'paused') AND NOT wp.inactive
    WHERE ws.tenant_id = req.tenant_id
      AND ws.skill_id = req.skill_id
      AND ws.verification_status IN ('manager_endorsed', 'hr_verified')
      AND NOT ws.inactive
      AND (req.min_years_experience IS NULL OR COALESCE(ws.years_experience, 0) >= req.min_years_experience)
      AND (
          CASE ws.proficiency
              WHEN 'beginner' THEN 1
              WHEN 'intermediate' THEN 2
              WHEN 'advanced' THEN 3
              WHEN 'expert' THEN 4
              ELSE 0
          END
      ) >= req.required_rank
) available
ORDER BY
    CASE req.importance WHEN 'critical' THEN 1 WHEN 'required' THEN 2 ELSE 3 END,
    match_percent ASC,
    req.skill_name ASC;

-- name: ListSkillGapSummary :many
WITH req AS (
    SELECT
        psr.*,
        p.name AS project_name,
        p.project_code,
        CASE psr.required_proficiency
            WHEN 'beginner' THEN 1
            WHEN 'intermediate' THEN 2
            WHEN 'advanced' THEN 3
            WHEN 'expert' THEN 4
            ELSE 1
        END AS required_rank
    FROM hrms.project_skill_requirements psr
    LEFT JOIN hrms.projects p ON p.tenant_id = psr.tenant_id AND p.id = psr.project_id AND NOT p.inactive
    WHERE psr.tenant_id = $1
      AND (sqlc.narg('project_id')::uuid IS NULL OR psr.project_id = sqlc.narg('project_id')::uuid)
      AND NOT psr.inactive
),
rows AS (
    SELECT
        req.project_id,
        req.project_name,
        req.project_code,
        req.importance,
        req.required_count,
        assigned.assigned_match_count::int AS assigned_match_count,
        GREATEST(req.required_count - assigned.assigned_match_count, 0)::int AS gap_count,
        LEAST(100, ROUND((assigned.assigned_match_count::numeric / req.required_count::numeric) * 100))::int AS match_percent,
        ((assigned.assigned_match_count = 1 AND req.required_count > 1) OR (available.available_match_count = 1 AND req.importance IN ('required', 'critical'))) AS single_person_dependency
    FROM req
    CROSS JOIN LATERAL (
        SELECT COUNT(DISTINCT ws.worker_profile_id) AS assigned_match_count
        FROM hrms.worker_skills ws
        WHERE ws.tenant_id = req.tenant_id
          AND ws.skill_id = req.skill_id
          AND ws.verification_status IN ('manager_endorsed', 'hr_verified')
          AND NOT ws.inactive
          AND (req.min_years_experience IS NULL OR COALESCE(ws.years_experience, 0) >= req.min_years_experience)
          AND (
              CASE ws.proficiency WHEN 'beginner' THEN 1 WHEN 'intermediate' THEN 2 WHEN 'advanced' THEN 3 WHEN 'expert' THEN 4 ELSE 0 END
          ) >= req.required_rank
          AND EXISTS (
              SELECT 1
              FROM hrms.engagements eg
              LEFT JOIN hrms.project_milestones pm ON pm.tenant_id = eg.tenant_id AND pm.engagement_id = eg.id AND NOT pm.inactive
              LEFT JOIN hrms.projects p2 ON p2.tenant_id = eg.tenant_id AND p2.id = req.project_id AND NOT p2.inactive
              WHERE eg.tenant_id = req.tenant_id
                AND eg.worker_profile_id = ws.worker_profile_id
                AND eg.status IN ('draft', 'active', 'paused')
                AND NOT eg.inactive
                AND (
                    (req.engagement_id IS NOT NULL AND eg.id = req.engagement_id)
                    OR (
                        req.engagement_id IS NULL
                        AND req.project_id IS NOT NULL
                        AND (
                            pm.project_id = req.project_id
                            OR (p2.project_code IS NOT NULL AND eg.project_code = p2.project_code)
                            OR (eg.project_label IS NOT NULL AND eg.project_label = p2.name)
                        )
                    )
                )
          )
    ) assigned
    CROSS JOIN LATERAL (
        SELECT COUNT(DISTINCT ws.worker_profile_id) AS available_match_count
        FROM hrms.worker_skills ws
        JOIN hrms.worker_profiles wp ON wp.tenant_id = ws.tenant_id AND wp.id = ws.worker_profile_id AND wp.profile_status IN ('active', 'draft', 'paused') AND NOT wp.inactive
        WHERE ws.tenant_id = req.tenant_id
          AND ws.skill_id = req.skill_id
          AND ws.verification_status IN ('manager_endorsed', 'hr_verified')
          AND NOT ws.inactive
          AND (req.min_years_experience IS NULL OR COALESCE(ws.years_experience, 0) >= req.min_years_experience)
          AND (CASE ws.proficiency WHEN 'beginner' THEN 1 WHEN 'intermediate' THEN 2 WHEN 'advanced' THEN 3 WHEN 'expert' THEN 4 ELSE 0 END) >= req.required_rank
    ) available
)
SELECT
    project_id,
    COALESCE(project_name, 'Unassigned project') AS project_name,
    project_code,
    COUNT(*)::int AS requirement_count,
    COUNT(*) FILTER (WHERE gap_count > 0)::int AS missing_skill_count,
    COUNT(*) FILTER (WHERE gap_count > 0 AND importance IN ('required', 'critical'))::int AS mandatory_gap_count,
    COALESCE(ROUND(AVG(match_percent))::int, 0)::int AS average_match_percent,
    COUNT(*) FILTER (WHERE single_person_dependency)::int AS single_person_dependency_count
FROM rows
GROUP BY project_id, project_name, project_code
ORDER BY mandatory_gap_count DESC, average_match_percent ASC, project_name ASC;

-- name: ListSinglePersonSkillDependencies :many
WITH gap AS (
    SELECT
        psr.id AS requirement_id,
        psr.tenant_id,
        psr.project_id,
        p.name AS project_name,
        psr.engagement_id,
        e.title AS engagement_title,
        s.id AS skill_id,
        s.name AS skill_name,
        psr.importance,
        psr.required_count,
        psr.min_years_experience,
        psr.required_proficiency,
        CASE psr.required_proficiency WHEN 'beginner' THEN 1 WHEN 'intermediate' THEN 2 WHEN 'advanced' THEN 3 WHEN 'expert' THEN 4 ELSE 1 END AS required_rank
    FROM hrms.project_skill_requirements psr
    LEFT JOIN hrms.projects p ON p.tenant_id = psr.tenant_id AND p.id = psr.project_id AND NOT p.inactive
    LEFT JOIN hrms.engagements e ON e.tenant_id = psr.tenant_id AND e.id = psr.engagement_id AND NOT e.inactive
    JOIN hrms.skills s ON s.id = psr.skill_id AND NOT s.inactive
    WHERE psr.tenant_id = $1
      AND (sqlc.narg('project_id')::uuid IS NULL OR psr.project_id = sqlc.narg('project_id')::uuid)
      AND NOT psr.inactive
)
SELECT
    gap.requirement_id,
    gap.project_id,
    gap.project_name,
    gap.engagement_id,
    gap.engagement_title,
    gap.skill_id,
    gap.skill_name,
    gap.importance,
    only_worker.worker_profile_id,
    only_worker.worker_display_name,
    only_worker.worker_code,
    only_worker.proficiency,
    only_worker.years_experience
FROM gap
CROSS JOIN LATERAL (
    SELECT
        ws.worker_profile_id,
        wp.display_name AS worker_display_name,
        wp.worker_code,
        ws.proficiency,
        ws.years_experience,
        COUNT(*) OVER () AS match_count
    FROM hrms.worker_skills ws
    JOIN hrms.worker_profiles wp ON wp.tenant_id = ws.tenant_id AND wp.id = ws.worker_profile_id AND wp.profile_status IN ('active', 'draft', 'paused') AND NOT wp.inactive
    WHERE ws.tenant_id = gap.tenant_id
      AND ws.skill_id = gap.skill_id
      AND ws.verification_status IN ('manager_endorsed', 'hr_verified')
      AND NOT ws.inactive
      AND (gap.min_years_experience IS NULL OR COALESCE(ws.years_experience, 0) >= gap.min_years_experience)
      AND (CASE ws.proficiency WHEN 'beginner' THEN 1 WHEN 'intermediate' THEN 2 WHEN 'advanced' THEN 3 WHEN 'expert' THEN 4 ELSE 0 END) >= gap.required_rank
) only_worker
WHERE only_worker.match_count = 1
ORDER BY
    CASE gap.importance WHEN 'critical' THEN 1 WHEN 'required' THEN 2 ELSE 3 END,
    gap.project_name ASC NULLS LAST,
    gap.skill_name ASC;
