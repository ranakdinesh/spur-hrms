-- name: CreateProject :one
INSERT INTO hrms.projects (
    tenant_id,
    project_code,
    name,
    description,
    status,
    department_id,
    branch_id,
    project_manager_id,
    start_date,
    due_date,
    completed_at,
    budget_amount,
    currency_code,
    billing_type,
    client_label,
    priority,
    notes,
    metadata,
    created_by,
    updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
    $11, $12, $13, $14, $15, $16, $17, $18, $19, $19
)
RETURNING *;

-- name: UpdateProject :one
UPDATE hrms.projects
SET project_code = $3,
    name = $4,
    description = $5,
    status = $6,
    department_id = $7,
    branch_id = $8,
    project_manager_id = $9,
    start_date = $10,
    due_date = $11,
    completed_at = $12,
    budget_amount = $13,
    currency_code = $14,
    billing_type = $15,
    client_label = $16,
    priority = $17,
    notes = $18,
    metadata = $19,
    updated_by = $20
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: UpdateProjectStatus :one
UPDATE hrms.projects
SET status = $3,
    completed_at = $4,
    updated_by = $5
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: GetProject :one
SELECT * FROM hrms.projects
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: ListProjects :many
SELECT
    p.id,
    p.tenant_id,
    p.project_code,
    p.name,
    p.description,
    p.status,
    p.department_id,
    d.name AS department_name,
    p.branch_id,
    b.branch_name AS branch_name,
    p.project_manager_id,
    p.start_date,
    p.due_date,
    p.completed_at,
    p.budget_amount,
    p.currency_code,
    p.billing_type,
    p.client_label,
    p.priority,
    p.notes,
    p.metadata,
    p.inactive,
    p.created_at,
    p.created_by,
    p.updated_at,
    p.updated_by,
    COUNT(pm.id)::int AS milestone_count,
    COUNT(pm.id) FILTER (WHERE pm.status = 'submitted')::int AS submitted_milestone_count,
    COUNT(pm.id) FILTER (WHERE pm.status = 'accepted')::int AS accepted_milestone_count,
    COUNT(pm.id) FILTER (WHERE pm.status = 'rejected')::int AS rejected_milestone_count,
    COALESCE(SUM(pm.amount) FILTER (WHERE NOT pm.inactive), 0)::numeric AS milestone_amount,
    COALESCE(SUM(pm.amount) FILTER (WHERE pm.status = 'accepted' AND NOT pm.inactive), 0)::numeric AS accepted_amount
FROM hrms.projects p
LEFT JOIN hrms.departments d ON d.tenant_id = p.tenant_id AND d.id = p.department_id AND NOT d.inactive
LEFT JOIN hrms.branches b ON b.tenant_id = p.tenant_id AND b.id = p.branch_id AND NOT b.inactive
LEFT JOIN hrms.project_milestones pm ON pm.tenant_id = p.tenant_id AND pm.project_id = p.id AND NOT pm.inactive
WHERE p.tenant_id = $1
  AND (sqlc.narg('status')::text IS NULL OR p.status = sqlc.narg('status')::text)
  AND (sqlc.narg('department_id')::uuid IS NULL OR p.department_id = sqlc.narg('department_id')::uuid)
  AND (sqlc.narg('branch_id')::uuid IS NULL OR p.branch_id = sqlc.narg('branch_id')::uuid)
  AND (sqlc.narg('project_manager_id')::uuid IS NULL OR p.project_manager_id = sqlc.narg('project_manager_id')::uuid)
  AND (
      sqlc.narg('search')::text IS NULL
      OR p.name ILIKE '%' || sqlc.narg('search')::text || '%'
      OR p.project_code ILIKE '%' || sqlc.narg('search')::text || '%'
      OR p.client_label ILIKE '%' || sqlc.narg('search')::text || '%'
  )
  AND NOT p.inactive
GROUP BY p.id, d.name, b.branch_name
ORDER BY
    CASE p.status WHEN 'active' THEN 0 WHEN 'draft' THEN 1 WHEN 'on_hold' THEN 2 WHEN 'completed' THEN 3 ELSE 4 END,
    p.due_date NULLS LAST,
    p.updated_at DESC;

-- name: SoftDeleteProject :exec
UPDATE hrms.projects
SET inactive = TRUE,
    status = 'cancelled',
    updated_by = $3
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: CreateProjectMilestone :one
INSERT INTO hrms.project_milestones (
    tenant_id,
    project_id,
    engagement_id,
    milestone_code,
    title,
    description,
    acceptance_criteria,
    due_date,
    status,
    amount,
    currency_code,
    payment_trigger,
    submitted_at,
    submitted_by,
    accepted_at,
    accepted_by,
    rejected_at,
    rejected_by,
    review_comment,
    notes,
    metadata,
    created_by,
    updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
    $11, $12, $13, $14, $15, $16, $17, $18, $19, $20,
    $21, $22, $22
)
RETURNING *;

-- name: UpdateProjectMilestone :one
UPDATE hrms.project_milestones
SET project_id = $3,
    engagement_id = $4,
    milestone_code = $5,
    title = $6,
    description = $7,
    acceptance_criteria = $8,
    due_date = $9,
    status = $10,
    amount = $11,
    currency_code = $12,
    payment_trigger = $13,
    submitted_at = $14,
    submitted_by = $15,
    accepted_at = $16,
    accepted_by = $17,
    rejected_at = $18,
    rejected_by = $19,
    review_comment = $20,
    notes = $21,
    metadata = $22,
    updated_by = $23
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: UpdateProjectMilestoneStatus :one
UPDATE hrms.project_milestones
SET status = $3,
    submitted_at = $4,
    submitted_by = $5,
    accepted_at = $6,
    accepted_by = $7,
    rejected_at = $8,
    rejected_by = $9,
    review_comment = $10,
    updated_by = $11
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: GetProjectMilestone :one
SELECT * FROM hrms.project_milestones
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: ListProjectMilestones :many
SELECT
    pm.id,
    pm.tenant_id,
    pm.project_id,
    p.name AS project_name,
    p.project_code,
    p.project_manager_id,
    p.department_id,
    d.name AS department_name,
    pm.engagement_id,
    e.title AS engagement_title,
    e.engagement_code,
    e.worker_profile_id,
    wp.display_name AS worker_display_name,
    wp.worker_code,
    pm.milestone_code,
    pm.title,
    pm.description,
    pm.acceptance_criteria,
    pm.due_date,
    pm.status,
    pm.amount,
    pm.currency_code,
    pm.payment_trigger,
    pm.submitted_at,
    pm.submitted_by,
    pm.accepted_at,
    pm.accepted_by,
    pm.rejected_at,
    pm.rejected_by,
    pm.review_comment,
    pm.notes,
    pm.metadata,
    pm.inactive,
    pm.created_at,
    pm.created_by,
    pm.updated_at,
    pm.updated_by
FROM hrms.project_milestones pm
JOIN hrms.projects p ON p.tenant_id = pm.tenant_id AND p.id = pm.project_id AND NOT p.inactive
LEFT JOIN hrms.departments d ON d.tenant_id = p.tenant_id AND d.id = p.department_id AND NOT d.inactive
LEFT JOIN hrms.engagements e ON e.tenant_id = pm.tenant_id AND e.id = pm.engagement_id AND NOT e.inactive
LEFT JOIN hrms.worker_profiles wp ON wp.tenant_id = pm.tenant_id AND wp.id = e.worker_profile_id AND NOT wp.inactive
WHERE pm.tenant_id = $1
  AND (sqlc.narg('project_id')::uuid IS NULL OR pm.project_id = sqlc.narg('project_id')::uuid)
  AND (sqlc.narg('engagement_id')::uuid IS NULL OR pm.engagement_id = sqlc.narg('engagement_id')::uuid)
  AND (sqlc.narg('status')::text IS NULL OR pm.status = sqlc.narg('status')::text)
  AND (
      sqlc.narg('search')::text IS NULL
      OR pm.title ILIKE '%' || sqlc.narg('search')::text || '%'
      OR pm.milestone_code ILIKE '%' || sqlc.narg('search')::text || '%'
      OR p.name ILIKE '%' || sqlc.narg('search')::text || '%'
      OR p.project_code ILIKE '%' || sqlc.narg('search')::text || '%'
      OR e.title ILIKE '%' || sqlc.narg('search')::text || '%'
      OR wp.display_name ILIKE '%' || sqlc.narg('search')::text || '%'
      OR wp.worker_code ILIKE '%' || sqlc.narg('search')::text || '%'
  )
  AND NOT pm.inactive
ORDER BY
    CASE pm.status WHEN 'submitted' THEN 0 WHEN 'open' THEN 1 WHEN 'draft' THEN 2 WHEN 'rejected' THEN 3 WHEN 'accepted' THEN 4 ELSE 5 END,
    pm.due_date NULLS LAST,
    pm.updated_at DESC;

-- name: SoftDeleteProjectMilestone :exec
UPDATE hrms.project_milestones
SET inactive = TRUE,
    status = 'cancelled',
    updated_by = $3
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: CreateProjectMilestoneEvent :one
INSERT INTO hrms.project_milestone_events (
    tenant_id,
    project_id,
    milestone_id,
    event_type,
    from_status,
    to_status,
    comment,
    actor_id,
    metadata
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING *;

-- name: ListProjectMilestoneEvents :many
SELECT * FROM hrms.project_milestone_events
WHERE tenant_id = $1 AND milestone_id = $2
ORDER BY created_at DESC;
