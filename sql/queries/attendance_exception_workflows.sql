-- name: CreateAttendanceExceptionWorkflow :one
INSERT INTO hrms.attendance_exception_workflows (
    tenant_id, code, name, description, branch_id, department_id, request_type,
    route_mode, max_requests_per_month, escalation_hours, escalation_route_mode,
    block_payroll_when_pending, auto_approve, is_active, created_by, updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7,
    $8, $9, $10, $11,
    $12, $13, $14, $15, $15
)
RETURNING *;

-- name: ListAttendanceExceptionWorkflows :many
SELECT * FROM hrms.attendance_exception_workflows
WHERE tenant_id = $1 AND NOT inactive
ORDER BY is_active DESC, request_type ASC, name ASC;

-- name: GetAttendanceExceptionWorkflow :one
SELECT * FROM hrms.attendance_exception_workflows
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: ResolveAttendanceExceptionWorkflow :one
SELECT * FROM hrms.attendance_exception_workflows
WHERE tenant_id = $1
  AND request_type = $2
  AND is_active
  AND NOT inactive
  AND (
    (department_id = $3)
    OR (department_id IS NULL AND branch_id = $4)
    OR (department_id IS NULL AND branch_id IS NULL)
  )
ORDER BY
  CASE
    WHEN department_id = $3 THEN 1
    WHEN branch_id = $4 THEN 2
    ELSE 3
  END,
  updated_at DESC
LIMIT 1;

-- name: UpdateAttendanceExceptionWorkflow :one
UPDATE hrms.attendance_exception_workflows
SET
    code = $3,
    name = $4,
    description = $5,
    branch_id = $6,
    department_id = $7,
    request_type = $8,
    route_mode = $9,
    max_requests_per_month = $10,
    escalation_hours = $11,
    escalation_route_mode = $12,
    block_payroll_when_pending = $13,
    auto_approve = $14,
    is_active = $15,
    updated_by = $16,
    updated_at = NOW()
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: SoftDeleteAttendanceExceptionWorkflow :exec
UPDATE hrms.attendance_exception_workflows
SET inactive = TRUE, updated_by = $3, updated_at = NOW()
WHERE tenant_id = $1 AND id = $2;

-- name: SetAttendanceRequestWorkflow :one
UPDATE hrms.attendance_requests
SET workflow_id = $3,
    route_mode = $4,
    escalation_due_at = $5,
    payroll_blocking = $6,
    updated_by = $7,
    updated_at = NOW()
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: CreateAttendanceExceptionEvent :one
INSERT INTO hrms.attendance_exception_events (
    tenant_id, attendance_request_id, workflow_id, action, from_status,
    to_status, routed_to, remarks, metadata, created_by, updated_by
) VALUES (
    $1, $2, $3, $4, $5,
    $6, $7, $8, $9, $10, $10
)
RETURNING *;

-- name: ListAttendanceExceptionEvents :many
SELECT * FROM hrms.attendance_exception_events
WHERE tenant_id = $1 AND attendance_request_id = $2 AND NOT inactive
ORDER BY created_at DESC;

-- name: ListPayrollBlockingAttendanceRequests :many
SELECT * FROM hrms.attendance_requests
WHERE tenant_id = $1
  AND date >= $2
  AND date <= $3
  AND payroll_blocking
  AND status = 'pending'
  AND NOT inactive
ORDER BY date ASC, created_at ASC;
