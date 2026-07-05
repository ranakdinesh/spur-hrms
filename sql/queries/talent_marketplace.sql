-- name: CreateTalentMarketplaceOpportunity :one
INSERT INTO hrms.talent_marketplace_opportunities (
    tenant_id, project_id, engagement_id, source_requirement_id, job_posting_id,
    title, description, opportunity_type, status, visibility, priority, seats,
    location_mode, min_allocation_percent, duration_label, start_date, due_date,
    candidate_fallback_enabled, candidate_fallback_status, metadata, created_by, updated_by
) VALUES (
    $1, $2, $3, $4, $5,
    $6, $7, $8, $9, $10, $11, $12,
    $13, $14, $15, $16, $17,
    $18, $19, $20, $21, $21
)
RETURNING *;

-- name: UpdateTalentMarketplaceOpportunity :one
UPDATE hrms.talent_marketplace_opportunities
SET project_id = $3,
    engagement_id = $4,
    source_requirement_id = $5,
    job_posting_id = $6,
    title = $7,
    description = $8,
    opportunity_type = $9,
    status = $10,
    visibility = $11,
    priority = $12,
    seats = $13,
    location_mode = $14,
    min_allocation_percent = $15,
    duration_label = $16,
    start_date = $17,
    due_date = $18,
    candidate_fallback_enabled = $19,
    candidate_fallback_status = $20,
    metadata = $21,
    updated_by = $22
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: UpdateTalentMarketplaceOpportunityFallback :one
UPDATE hrms.talent_marketplace_opportunities
SET candidate_fallback_enabled = TRUE,
    candidate_fallback_status = $3,
    updated_by = $4
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: GetTalentMarketplaceOpportunity :one
SELECT * FROM hrms.talent_marketplace_opportunities
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: ListTalentMarketplaceOpportunities :many
SELECT
    tmo.*,
    p.name AS project_name,
    p.project_code,
    e.title AS engagement_title,
    e.engagement_code,
    jp.title AS job_posting_title,
    jp.code AS job_posting_code,
    COUNT(DISTINCT tma.id)::int AS application_count,
    COUNT(DISTINCT tma.id) FILTER (WHERE tma.status IN ('recommended', 'invited'))::int AS recommended_count,
    COUNT(DISTINCT tma.id) FILTER (WHERE tma.status IN ('accepted', 'assigned'))::int AS selected_count
FROM hrms.talent_marketplace_opportunities tmo
LEFT JOIN hrms.projects p ON p.tenant_id = tmo.tenant_id AND p.id = tmo.project_id AND NOT p.inactive
LEFT JOIN hrms.engagements e ON e.tenant_id = tmo.tenant_id AND e.id = tmo.engagement_id AND NOT e.inactive
LEFT JOIN hrms.job_postings jp ON jp.tenant_id = tmo.tenant_id AND jp.id = tmo.job_posting_id AND NOT jp.inactive
LEFT JOIN hrms.talent_marketplace_applications tma ON tma.tenant_id = tmo.tenant_id AND tma.opportunity_id = tmo.id AND NOT tma.inactive
WHERE tmo.tenant_id = $1
  AND (sqlc.narg('project_id')::uuid IS NULL OR tmo.project_id = sqlc.narg('project_id')::uuid)
  AND (sqlc.narg('engagement_id')::uuid IS NULL OR tmo.engagement_id = sqlc.narg('engagement_id')::uuid)
  AND (sqlc.narg('status')::text IS NULL OR tmo.status = sqlc.narg('status')::text)
  AND (sqlc.narg('opportunity_type')::text IS NULL OR tmo.opportunity_type = sqlc.narg('opportunity_type')::text)
  AND (
      sqlc.narg('search')::text IS NULL
      OR tmo.title ILIKE '%' || sqlc.narg('search')::text || '%'
      OR p.name ILIKE '%' || sqlc.narg('search')::text || '%'
      OR e.title ILIKE '%' || sqlc.narg('search')::text || '%'
      OR jp.title ILIKE '%' || sqlc.narg('search')::text || '%'
  )
  AND NOT tmo.inactive
GROUP BY tmo.id, p.name, p.project_code, e.title, e.engagement_code, jp.title, jp.code
ORDER BY
    CASE tmo.priority WHEN 'critical' THEN 1 WHEN 'high' THEN 2 WHEN 'normal' THEN 3 ELSE 4 END,
    CASE tmo.status WHEN 'open' THEN 1 WHEN 'draft' THEN 2 WHEN 'paused' THEN 3 ELSE 4 END,
    tmo.updated_at DESC;

-- name: SoftDeleteTalentMarketplaceOpportunity :exec
UPDATE hrms.talent_marketplace_opportunities
SET inactive = TRUE,
    updated_by = $3
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: CreateTalentMarketplaceApplication :one
INSERT INTO hrms.talent_marketplace_applications (
    tenant_id, opportunity_id, worker_profile_id, status, match_score,
    match_reasons, worker_note, manager_note, created_by, updated_by
) VALUES (
    $1, $2, $3, $4, $5,
    $6, $7, $8, $9, $9
)
RETURNING *;

-- name: UpdateTalentMarketplaceApplicationStatus :one
UPDATE hrms.talent_marketplace_applications
SET status = $3,
    worker_note = COALESCE($4, worker_note),
    manager_note = COALESCE($5, manager_note),
    decided_at = CASE WHEN $3 IN ('accepted', 'declined', 'withdrawn', 'rejected', 'assigned') THEN now() ELSE decided_at END,
    decided_by = CASE WHEN $3 IN ('accepted', 'declined', 'withdrawn', 'rejected', 'assigned') THEN $6 ELSE decided_by END,
    updated_by = $6
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: GetTalentMarketplaceApplication :one
SELECT * FROM hrms.talent_marketplace_applications
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: ListTalentMarketplaceApplications :many
SELECT
    tma.*,
    tmo.title AS opportunity_title,
    tmo.status AS opportunity_status,
    wp.display_name AS worker_display_name,
    wp.worker_code,
    p.name AS project_name,
    e.title AS engagement_title
FROM hrms.talent_marketplace_applications tma
JOIN hrms.talent_marketplace_opportunities tmo ON tmo.tenant_id = tma.tenant_id AND tmo.id = tma.opportunity_id AND NOT tmo.inactive
JOIN hrms.worker_profiles wp ON wp.tenant_id = tma.tenant_id AND wp.id = tma.worker_profile_id AND NOT wp.inactive
LEFT JOIN hrms.projects p ON p.tenant_id = tmo.tenant_id AND p.id = tmo.project_id AND NOT p.inactive
LEFT JOIN hrms.engagements e ON e.tenant_id = tmo.tenant_id AND e.id = tmo.engagement_id AND NOT e.inactive
WHERE tma.tenant_id = $1
  AND (sqlc.narg('opportunity_id')::uuid IS NULL OR tma.opportunity_id = sqlc.narg('opportunity_id')::uuid)
  AND (sqlc.narg('worker_profile_id')::uuid IS NULL OR tma.worker_profile_id = sqlc.narg('worker_profile_id')::uuid)
  AND (sqlc.narg('status')::text IS NULL OR tma.status = sqlc.narg('status')::text)
  AND NOT tma.inactive
ORDER BY tma.updated_at DESC;

-- name: ListTalentMarketplaceRecommendations :many
WITH opp AS (
    SELECT *
    FROM hrms.talent_marketplace_opportunities
    WHERE tenant_id = $1 AND id = $2 AND NOT inactive
),
req AS (
    SELECT
        psr.*,
        CASE psr.required_proficiency
            WHEN 'beginner' THEN 1
            WHEN 'intermediate' THEN 2
            WHEN 'advanced' THEN 3
            WHEN 'expert' THEN 4
            ELSE 1
        END AS required_rank
    FROM hrms.project_skill_requirements psr
    JOIN opp ON opp.tenant_id = psr.tenant_id
    WHERE NOT psr.inactive
      AND (
          (opp.source_requirement_id IS NOT NULL AND psr.id = opp.source_requirement_id)
          OR (opp.source_requirement_id IS NULL AND opp.engagement_id IS NOT NULL AND psr.engagement_id = opp.engagement_id)
          OR (opp.source_requirement_id IS NULL AND opp.engagement_id IS NULL AND opp.project_id IS NOT NULL AND psr.project_id = opp.project_id)
      )
),
worker_base AS (
    SELECT wp.id, wp.display_name, wp.worker_code, wp.profile_status
    FROM hrms.worker_profiles wp
    JOIN opp ON opp.tenant_id = wp.tenant_id
    WHERE wp.profile_status IN ('active', 'draft', 'paused')
      AND NOT wp.inactive
),
scored AS (
    SELECT
        wb.id AS worker_profile_id,
        wb.display_name AS worker_display_name,
        wb.worker_code,
        COUNT(DISTINCT req.id)::int AS required_skill_count,
        COUNT(DISTINCT req.id) FILTER (
            WHERE EXISTS (
                SELECT 1
                FROM hrms.worker_skills ws
                WHERE ws.tenant_id = $1
                  AND ws.worker_profile_id = wb.id
                  AND ws.skill_id = req.skill_id
                  AND ws.verification_status IN ('manager_endorsed', 'hr_verified')
                  AND NOT ws.inactive
                  AND (req.min_years_experience IS NULL OR COALESCE(ws.years_experience, 0) >= req.min_years_experience)
                  AND (CASE ws.proficiency WHEN 'beginner' THEN 1 WHEN 'intermediate' THEN 2 WHEN 'advanced' THEN 3 WHEN 'expert' THEN 4 ELSE 0 END) >= req.required_rank
            )
        )::int AS matched_skill_count
    FROM worker_base wb
    LEFT JOIN req ON TRUE
    GROUP BY wb.id, wb.display_name, wb.worker_code
)
SELECT
    scored.worker_profile_id,
    scored.worker_display_name,
    scored.worker_code,
    scored.required_skill_count,
    scored.matched_skill_count,
    GREATEST(scored.required_skill_count - scored.matched_skill_count, 0)::int AS missing_skill_count,
    (CASE WHEN scored.required_skill_count = 0 THEN 0 ELSE ROUND((scored.matched_skill_count::numeric / scored.required_skill_count::numeric) * 100, 2) END)::numeric(5,2) AS match_score,
    jsonb_build_object(
        'matched_skill_count', scored.matched_skill_count,
        'required_skill_count', scored.required_skill_count,
        'missing_skill_count', GREATEST(scored.required_skill_count - scored.matched_skill_count, 0),
        'basis', 'verified_skills'
    ) AS match_reasons,
    existing.id AS application_id,
    existing.status AS application_status
FROM scored
LEFT JOIN hrms.talent_marketplace_applications existing
    ON existing.tenant_id = $1
   AND existing.opportunity_id = $2
   AND existing.worker_profile_id = scored.worker_profile_id
   AND NOT existing.inactive
ORDER BY match_score DESC, scored.worker_display_name ASC;

-- name: CreateTalentMarketplaceEvent :one
INSERT INTO hrms.talent_marketplace_events (
    tenant_id, opportunity_id, application_id, actor_user_id, action,
    from_status, to_status, notes, metadata
) VALUES (
    $1, $2, $3, $4, $5,
    $6, $7, $8, $9
)
RETURNING *;

-- name: ListTalentMarketplaceEvents :many
SELECT
    tme.*,
    tmo.title AS opportunity_title,
    tma.worker_profile_id,
    wp.display_name AS worker_display_name
FROM hrms.talent_marketplace_events tme
LEFT JOIN hrms.talent_marketplace_opportunities tmo ON tmo.tenant_id = tme.tenant_id AND tmo.id = tme.opportunity_id
LEFT JOIN hrms.talent_marketplace_applications tma ON tma.tenant_id = tme.tenant_id AND tma.id = tme.application_id
LEFT JOIN hrms.worker_profiles wp ON wp.tenant_id = tma.tenant_id AND wp.id = tma.worker_profile_id
WHERE tme.tenant_id = $1
  AND (sqlc.narg('opportunity_id')::uuid IS NULL OR tme.opportunity_id = sqlc.narg('opportunity_id')::uuid)
  AND (sqlc.narg('application_id')::uuid IS NULL OR tme.application_id = sqlc.narg('application_id')::uuid)
  AND NOT tme.inactive
ORDER BY tme.created_at DESC;
