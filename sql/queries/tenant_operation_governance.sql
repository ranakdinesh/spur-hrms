-- name: CreateTenantOperationRequest :one
INSERT INTO hrms.tenant_operation_requests (
    operation_number, operation_type, title, target_tenant_id, target_tenant_name, target_tenant_code,
    status, risk_level, reason, requested_by, approval_required, backup_required, backup_confirmed,
    retention_until, request_payload, validation_results, rollback_metadata, metadata, created_by, updated_by
) VALUES (
    sqlc.arg(operation_number), sqlc.arg(operation_type), sqlc.arg(title), sqlc.narg(target_tenant_id), sqlc.narg(target_tenant_name), sqlc.narg(target_tenant_code),
    sqlc.arg(status), sqlc.arg(risk_level), sqlc.arg(reason), sqlc.narg(actor_id), sqlc.arg(approval_required), sqlc.arg(backup_required), sqlc.arg(backup_confirmed),
    sqlc.narg(retention_until), sqlc.arg(request_payload), sqlc.arg(validation_results), sqlc.arg(rollback_metadata), sqlc.arg(metadata), sqlc.narg(actor_id), sqlc.narg(actor_id)
)
RETURNING *;

-- name: GetTenantOperationRequest :one
SELECT * FROM hrms.tenant_operation_requests
WHERE id = $1 AND NOT inactive;

-- name: ListTenantOperationRequests :many
SELECT * FROM hrms.tenant_operation_requests
WHERE NOT inactive
  AND (sqlc.narg(status)::text IS NULL OR status = sqlc.narg(status)::text)
  AND (sqlc.narg(operation_type)::text IS NULL OR operation_type = sqlc.narg(operation_type)::text)
  AND (sqlc.narg(risk_level)::text IS NULL OR risk_level = sqlc.narg(risk_level)::text)
  AND (sqlc.narg(target_tenant_id)::uuid IS NULL OR target_tenant_id = sqlc.narg(target_tenant_id)::uuid)
  AND (
      sqlc.narg(search)::text IS NULL
      OR operation_number ILIKE '%' || sqlc.narg(search)::text || '%'
      OR title ILIKE '%' || sqlc.narg(search)::text || '%'
      OR reason ILIKE '%' || sqlc.narg(search)::text || '%'
      OR COALESCE(target_tenant_name, '') ILIKE '%' || sqlc.narg(search)::text || '%'
      OR COALESCE(target_tenant_code, '') ILIKE '%' || sqlc.narg(search)::text || '%'
  )
ORDER BY
  CASE status WHEN 'pending_validation' THEN 0 WHEN 'pending_approval' THEN 1 WHEN 'approved' THEN 2 WHEN 'in_progress' THEN 3 ELSE 4 END,
  CASE risk_level WHEN 'critical' THEN 0 WHEN 'high' THEN 1 WHEN 'medium' THEN 2 ELSE 3 END,
  created_at DESC
LIMIT sqlc.arg(row_limit) OFFSET sqlc.arg(row_offset);

-- name: UpdateTenantOperationRequestStatus :one
UPDATE hrms.tenant_operation_requests
SET status = sqlc.arg(status),
    approved_by = COALESCE(sqlc.narg(approved_by), approved_by),
    approved_at = CASE WHEN sqlc.arg(status) = 'approved' THEN now() ELSE approved_at END,
    completed_by = COALESCE(sqlc.narg(completed_by), completed_by),
    completed_at = CASE WHEN sqlc.arg(status) IN ('completed','failed','cancelled','rejected') THEN now() ELSE completed_at END,
    backup_confirmed = COALESCE(sqlc.narg(backup_confirmed), backup_confirmed),
    validation_results = COALESCE(sqlc.narg(validation_results), validation_results),
    rollback_metadata = COALESCE(sqlc.narg(rollback_metadata), rollback_metadata),
    metadata = metadata || COALESCE(sqlc.narg(metadata), '{}'::jsonb),
    updated_by = sqlc.narg(actor_id)
WHERE id = sqlc.arg(id) AND NOT inactive
RETURNING *;

-- name: CreateTenantOperationEvent :one
INSERT INTO hrms.tenant_operation_events (
    request_id, action, from_status, to_status, actor_user_id, remarks, metadata, created_by, updated_by
) VALUES (
    sqlc.arg(request_id), sqlc.arg(action), sqlc.narg(from_status), sqlc.narg(to_status), sqlc.narg(actor_id), sqlc.narg(remarks), sqlc.arg(metadata), sqlc.narg(actor_id), sqlc.narg(actor_id)
)
RETURNING *;

-- name: ListTenantOperationEvents :many
SELECT * FROM hrms.tenant_operation_events
WHERE request_id = $1 AND NOT inactive
ORDER BY created_at DESC, id DESC;
