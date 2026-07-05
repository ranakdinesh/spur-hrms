-- name: CreateEngagement :one
INSERT INTO hrms.engagements (
    tenant_id,
    worker_profile_id,
    engagement_code,
    title,
    description,
    engagement_type,
    status,
    start_date,
    end_date,
    hours_budget,
    rate_amount,
    currency_code,
    rate_unit,
    branch_id,
    department_id,
    reporting_manager_id,
    project_label,
    project_code,
    cost_center,
    renewal_due_date,
    renewal_status,
    termination_reason,
    terminated_at,
    notes,
    metadata,
    created_by,
    updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
    $11, $12, $13, $14, $15, $16, $17, $18, $19, $20,
    $21, $22, $23, $24, $25, $26, $26
)
RETURNING *;

-- name: UpdateEngagement :one
UPDATE hrms.engagements
SET worker_profile_id = $3,
    engagement_code = $4,
    title = $5,
    description = $6,
    engagement_type = $7,
    status = $8,
    start_date = $9,
    end_date = $10,
    hours_budget = $11,
    rate_amount = $12,
    currency_code = $13,
    rate_unit = $14,
    branch_id = $15,
    department_id = $16,
    reporting_manager_id = $17,
    project_label = $18,
    project_code = $19,
    cost_center = $20,
    renewal_due_date = $21,
    renewal_status = $22,
    termination_reason = $23,
    terminated_at = $24,
    notes = $25,
    metadata = $26,
    updated_by = $27
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: UpdateEngagementStatus :one
UPDATE hrms.engagements
SET status = $3,
    termination_reason = $4,
    terminated_at = $5,
    updated_by = $6
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: GetEngagement :one
SELECT * FROM hrms.engagements
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: ListEngagements :many
SELECT
    e.id,
    e.tenant_id,
    e.worker_profile_id,
    wp.display_name AS worker_display_name,
    wp.worker_code,
    wp.employee_id,
    wt.name AS worker_type_name,
    wt.classification_group,
    e.engagement_code,
    e.title,
    e.description,
    e.engagement_type,
    e.status,
    e.start_date,
    e.end_date,
    e.hours_budget,
    e.rate_amount,
    e.currency_code,
    e.rate_unit,
    e.branch_id,
    b.branch_name AS branch_name,
    e.department_id,
    d.name AS department_name,
    e.reporting_manager_id,
    e.project_label,
    e.project_code,
    e.cost_center,
    e.renewal_due_date,
    e.renewal_status,
    e.termination_reason,
    e.terminated_at,
    e.notes,
    e.metadata,
    e.inactive,
    e.created_at,
    e.created_by,
    e.updated_at,
    e.updated_by
FROM hrms.engagements e
JOIN hrms.worker_profiles wp ON wp.tenant_id = e.tenant_id AND wp.id = e.worker_profile_id AND NOT wp.inactive
JOIN hrms.worker_types wt ON wt.tenant_id = wp.tenant_id AND wt.id = wp.worker_type_id AND NOT wt.inactive
LEFT JOIN hrms.branches b ON b.tenant_id = e.tenant_id AND b.id = e.branch_id AND NOT b.inactive
LEFT JOIN hrms.departments d ON d.tenant_id = e.tenant_id AND d.id = e.department_id AND NOT d.inactive
WHERE e.tenant_id = $1
  AND (sqlc.narg('worker_profile_id')::uuid IS NULL OR e.worker_profile_id = sqlc.narg('worker_profile_id')::uuid)
  AND (sqlc.narg('engagement_type')::text IS NULL OR e.engagement_type = sqlc.narg('engagement_type')::text)
  AND (sqlc.narg('status')::text IS NULL OR e.status = sqlc.narg('status')::text)
  AND (sqlc.narg('department_id')::uuid IS NULL OR e.department_id = sqlc.narg('department_id')::uuid)
  AND (
      sqlc.narg('search')::text IS NULL
      OR e.title ILIKE '%' || sqlc.narg('search')::text || '%'
      OR e.engagement_code ILIKE '%' || sqlc.narg('search')::text || '%'
      OR e.project_label ILIKE '%' || sqlc.narg('search')::text || '%'
      OR e.project_code ILIKE '%' || sqlc.narg('search')::text || '%'
      OR wp.display_name ILIKE '%' || sqlc.narg('search')::text || '%'
      OR wp.worker_code ILIKE '%' || sqlc.narg('search')::text || '%'
  )
  AND NOT e.inactive
ORDER BY
    CASE e.status WHEN 'active' THEN 0 WHEN 'draft' THEN 1 WHEN 'paused' THEN 2 ELSE 3 END,
    e.start_date DESC,
    e.updated_at DESC;

-- name: SoftDeleteEngagement :exec
UPDATE hrms.engagements
SET inactive = TRUE,
    status = 'cancelled',
    updated_by = $3
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: CreateWorkLog :one
INSERT INTO hrms.work_logs (
    tenant_id,
    engagement_id,
    worker_profile_id,
    log_date,
    hours_worked,
    billable_hours,
    work_summary,
    deliverable_reference,
    status,
    submitted_at,
    submitted_by,
    reviewed_at,
    reviewed_by,
    review_comment,
    metadata,
    created_by,
    updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
    $11, $12, $13, $14, $15, $16, $16
)
RETURNING *;

-- name: UpdateWorkLog :one
UPDATE hrms.work_logs
SET log_date = $3,
    hours_worked = $4,
    billable_hours = $5,
    work_summary = $6,
    deliverable_reference = $7,
    status = $8,
    submitted_at = $9,
    submitted_by = $10,
    reviewed_at = $11,
    reviewed_by = $12,
    review_comment = $13,
    metadata = $14,
    updated_by = $15
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: UpdateWorkLogStatus :one
UPDATE hrms.work_logs
SET status = $3,
    submitted_at = $4,
    submitted_by = $5,
    reviewed_at = $6,
    reviewed_by = $7,
    review_comment = $8,
    updated_by = $9
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: GetWorkLog :one
SELECT * FROM hrms.work_logs
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: ListWorkLogs :many
SELECT
    wl.id,
    wl.tenant_id,
    wl.engagement_id,
    e.title AS engagement_title,
    e.engagement_code,
    e.project_label,
    e.project_code,
    e.cost_center,
    wl.worker_profile_id,
    wp.display_name AS worker_display_name,
    wp.worker_code,
    wp.employee_id,
    e.reporting_manager_id,
    e.department_id,
    d.name AS department_name,
    e.branch_id,
    b.branch_name AS branch_name,
    wl.log_date,
    wl.hours_worked,
    wl.billable_hours,
    wl.work_summary,
    wl.deliverable_reference,
    wl.status,
    wl.submitted_at,
    wl.submitted_by,
    wl.reviewed_at,
    wl.reviewed_by,
    wl.review_comment,
    wl.metadata,
    wl.inactive,
    wl.created_at,
    wl.created_by,
    wl.updated_at,
    wl.updated_by
FROM hrms.work_logs wl
JOIN hrms.engagements e ON e.tenant_id = wl.tenant_id AND e.id = wl.engagement_id AND NOT e.inactive
JOIN hrms.worker_profiles wp ON wp.tenant_id = wl.tenant_id AND wp.id = wl.worker_profile_id AND NOT wp.inactive
LEFT JOIN hrms.departments d ON d.tenant_id = wl.tenant_id AND d.id = e.department_id AND NOT d.inactive
LEFT JOIN hrms.branches b ON b.tenant_id = wl.tenant_id AND b.id = e.branch_id AND NOT b.inactive
WHERE wl.tenant_id = $1
  AND (sqlc.narg('engagement_id')::uuid IS NULL OR wl.engagement_id = sqlc.narg('engagement_id')::uuid)
  AND (sqlc.narg('worker_profile_id')::uuid IS NULL OR wl.worker_profile_id = sqlc.narg('worker_profile_id')::uuid)
  AND (sqlc.narg('status')::text IS NULL OR wl.status = sqlc.narg('status')::text)
  AND (sqlc.narg('date_from')::date IS NULL OR wl.log_date >= sqlc.narg('date_from')::date)
  AND (sqlc.narg('date_to')::date IS NULL OR wl.log_date <= sqlc.narg('date_to')::date)
  AND (
      sqlc.narg('search')::text IS NULL
      OR wl.work_summary ILIKE '%' || sqlc.narg('search')::text || '%'
      OR wl.deliverable_reference ILIKE '%' || sqlc.narg('search')::text || '%'
      OR e.title ILIKE '%' || sqlc.narg('search')::text || '%'
      OR e.project_label ILIKE '%' || sqlc.narg('search')::text || '%'
      OR wp.display_name ILIKE '%' || sqlc.narg('search')::text || '%'
      OR wp.worker_code ILIKE '%' || sqlc.narg('search')::text || '%'
  )
  AND NOT wl.inactive
ORDER BY
    CASE wl.status WHEN 'submitted' THEN 0 WHEN 'draft' THEN 1 WHEN 'rejected' THEN 2 WHEN 'approved' THEN 3 ELSE 4 END,
    wl.log_date DESC,
    wl.updated_at DESC;

-- name: GetWorkLogBudgetUsage :one
SELECT
    e.hours_budget,
    COALESCE(SUM(wl.hours_worked) FILTER (
        WHERE wl.status IN ('submitted', 'approved')
          AND NOT wl.inactive
          AND (sqlc.narg('exclude_work_log_id')::uuid IS NULL OR wl.id <> sqlc.narg('exclude_work_log_id')::uuid)
    ), 0)::numeric AS used_hours
FROM hrms.engagements e
LEFT JOIN hrms.work_logs wl ON wl.tenant_id = e.tenant_id AND wl.engagement_id = e.id
WHERE e.tenant_id = $1 AND e.id = $2 AND NOT e.inactive
GROUP BY e.hours_budget;

-- name: ListWorkLogRollups :many
SELECT
    wl.tenant_id,
    wl.engagement_id,
    e.title AS engagement_title,
    e.engagement_code,
    wl.worker_profile_id,
    wp.display_name AS worker_display_name,
    COUNT(*)::int AS log_count,
    COALESCE(SUM(wl.hours_worked), 0)::numeric AS total_hours,
    COALESCE(SUM(wl.billable_hours), 0)::numeric AS billable_hours,
    COALESCE(SUM(wl.hours_worked) FILTER (WHERE wl.status = 'approved'), 0)::numeric AS approved_hours,
    COALESCE(SUM(wl.hours_worked) FILTER (WHERE wl.status = 'submitted'), 0)::numeric AS submitted_hours,
    COALESCE(SUM(wl.hours_worked) FILTER (WHERE wl.status = 'rejected'), 0)::numeric AS rejected_hours,
    e.hours_budget
FROM hrms.work_logs wl
JOIN hrms.engagements e ON e.tenant_id = wl.tenant_id AND e.id = wl.engagement_id AND NOT e.inactive
JOIN hrms.worker_profiles wp ON wp.tenant_id = wl.tenant_id AND wp.id = wl.worker_profile_id AND NOT wp.inactive
WHERE wl.tenant_id = $1
  AND (sqlc.narg('engagement_id')::uuid IS NULL OR wl.engagement_id = sqlc.narg('engagement_id')::uuid)
  AND (sqlc.narg('worker_profile_id')::uuid IS NULL OR wl.worker_profile_id = sqlc.narg('worker_profile_id')::uuid)
  AND (sqlc.narg('date_from')::date IS NULL OR wl.log_date >= sqlc.narg('date_from')::date)
  AND (sqlc.narg('date_to')::date IS NULL OR wl.log_date <= sqlc.narg('date_to')::date)
  AND NOT wl.inactive
GROUP BY wl.tenant_id, wl.engagement_id, e.title, e.engagement_code, wl.worker_profile_id, wp.display_name, e.hours_budget
ORDER BY e.title, wp.display_name;

-- name: SoftDeleteWorkLog :exec
UPDATE hrms.work_logs
SET inactive = TRUE,
    status = 'cancelled',
    updated_by = $3
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;
