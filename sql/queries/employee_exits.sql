-- name: CreateEmployeeExitRequest :one
INSERT INTO hrms.employee_exit_requests (
    tenant_id, employee_id, employee_user_id, initiated_by, status, exit_type, reason,
    resignation_date, notice_start_date, last_working_date, requested_relieving_date,
    final_settlement_status, access_revocation_status, asset_clearance_status,
    handover_status, exit_interview_status, notes, created_by
)
VALUES (
    $1, $2, $3, $4, $5, $6, $7,
    $8, $9, $10, $11,
    $12, $13, $14,
    $15, $16, $17, $18
)
RETURNING *;

-- name: ListEmployeeExitRequests :many
SELECT
    er.*,
    e.firstname AS employee_firstname,
    e.lastname AS employee_lastname,
    e.employee_code AS employee_code,
    e.email AS employee_email,
    d.name AS department_name,
    b.branch_name AS branch_name,
    COALESCE(tasks.total_tasks, 0)::int AS total_tasks,
    COALESCE(tasks.completed_tasks, 0)::int AS completed_tasks,
    COALESCE(tasks.blocked_tasks, 0)::int AS blocked_tasks
FROM hrms.employee_exit_requests er
JOIN hrms.employees e ON e.tenant_id = er.tenant_id AND e.id = er.employee_id
LEFT JOIN hrms.departments d ON d.tenant_id = e.tenant_id AND d.id = e.department_id AND NOT d.inactive
LEFT JOIN hrms.branches b ON b.tenant_id = e.tenant_id AND b.id = e.branch_id AND NOT b.inactive
LEFT JOIN LATERAL (
    SELECT
        COUNT(*) AS total_tasks,
        COUNT(*) FILTER (WHERE status IN ('completed','waived')) AS completed_tasks,
        COUNT(*) FILTER (WHERE status = 'blocked') AS blocked_tasks
    FROM hrms.employee_exit_tasks t
    WHERE t.tenant_id = er.tenant_id AND t.exit_request_id = er.id AND NOT t.inactive
) tasks ON TRUE
WHERE er.tenant_id = $1
  AND (sqlc.narg('status')::text IS NULL OR er.status = sqlc.narg('status')::text)
  AND (sqlc.narg('employee_user_id')::uuid IS NULL OR er.employee_user_id = sqlc.narg('employee_user_id')::uuid)
  AND (
      sqlc.narg('search')::text IS NULL
      OR e.firstname ILIKE '%' || sqlc.narg('search')::text || '%'
      OR e.lastname ILIKE '%' || sqlc.narg('search')::text || '%'
      OR e.employee_code ILIKE '%' || sqlc.narg('search')::text || '%'
      OR e.email ILIKE '%' || sqlc.narg('search')::text || '%'
  )
  AND NOT er.inactive
ORDER BY er.created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: CountEmployeeExitRequests :one
SELECT COUNT(*)
FROM hrms.employee_exit_requests er
JOIN hrms.employees e ON e.tenant_id = er.tenant_id AND e.id = er.employee_id
WHERE er.tenant_id = $1
  AND (sqlc.narg('status')::text IS NULL OR er.status = sqlc.narg('status')::text)
  AND (sqlc.narg('employee_user_id')::uuid IS NULL OR er.employee_user_id = sqlc.narg('employee_user_id')::uuid)
  AND (
      sqlc.narg('search')::text IS NULL
      OR e.firstname ILIKE '%' || sqlc.narg('search')::text || '%'
      OR e.lastname ILIKE '%' || sqlc.narg('search')::text || '%'
      OR e.employee_code ILIKE '%' || sqlc.narg('search')::text || '%'
      OR e.email ILIKE '%' || sqlc.narg('search')::text || '%'
  )
  AND NOT er.inactive;

-- name: GetEmployeeExitRequest :one
SELECT
    er.*,
    e.firstname AS employee_firstname,
    e.lastname AS employee_lastname,
    e.employee_code AS employee_code,
    e.email AS employee_email,
    d.name AS department_name,
    b.branch_name AS branch_name
FROM hrms.employee_exit_requests er
JOIN hrms.employees e ON e.tenant_id = er.tenant_id AND e.id = er.employee_id
LEFT JOIN hrms.departments d ON d.tenant_id = e.tenant_id AND d.id = e.department_id AND NOT d.inactive
LEFT JOIN hrms.branches b ON b.tenant_id = e.tenant_id AND b.id = e.branch_id AND NOT b.inactive
WHERE er.tenant_id = $1 AND er.id = $2 AND NOT er.inactive;

-- name: GetActiveEmployeeExitRequestByUserID :one
SELECT *
FROM hrms.employee_exit_requests
WHERE tenant_id = $1
  AND employee_user_id = $2
  AND status IN ('submitted','approved')
  AND NOT inactive
LIMIT 1;

-- name: UpdateEmployeeExitRequestStatus :one
UPDATE hrms.employee_exit_requests
SET status = $3,
    approved_by = CASE WHEN $3 = 'approved' THEN $4 ELSE approved_by END,
    approved_at = CASE WHEN $3 = 'approved' THEN NOW() ELSE approved_at END,
    completed_by = CASE WHEN $3 = 'completed' THEN $4 ELSE completed_by END,
    completed_at = CASE WHEN $3 = 'completed' THEN NOW() ELSE completed_at END,
    approved_relieving_date = COALESCE($5, approved_relieving_date),
    rejection_reason = CASE WHEN $3 = 'rejected' THEN $6 ELSE rejection_reason END,
    cancel_reason = CASE WHEN $3 = 'canceled' THEN $6 ELSE cancel_reason END,
    updated_at = NOW(),
    updated_by = $4
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: CreateEmployeeExitTask :one
INSERT INTO hrms.employee_exit_tasks (
    tenant_id, exit_request_id, employee_user_id, task_key, title, description,
    owner_role, due_date, status, remarks, sort_order, created_by
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
RETURNING *;

-- name: ListEmployeeExitTasks :many
SELECT *
FROM hrms.employee_exit_tasks
WHERE tenant_id = $1 AND exit_request_id = $2 AND NOT inactive
ORDER BY sort_order ASC, created_at ASC;

-- name: GetEmployeeExitTask :one
SELECT *
FROM hrms.employee_exit_tasks
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: UpdateEmployeeExitTaskStatus :one
UPDATE hrms.employee_exit_tasks
SET status = $3,
    remarks = $4,
    completed_by = CASE WHEN $3 IN ('completed','waived') THEN $5 ELSE completed_by END,
    completed_at = CASE WHEN $3 IN ('completed','waived') THEN NOW() ELSE NULL END,
    updated_at = NOW(),
    updated_by = $5
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: CreateEmployeeExitEvent :one
INSERT INTO hrms.employee_exit_events (
    tenant_id, exit_request_id, exit_task_id, action, from_status, to_status, remarks, metadata, created_by
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: ListEmployeeExitEvents :many
SELECT *
FROM hrms.employee_exit_events
WHERE tenant_id = $1 AND exit_request_id = $2 AND NOT inactive
ORDER BY created_at DESC;
