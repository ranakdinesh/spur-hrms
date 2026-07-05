-- name: CreateEmployeeCredentialEvent :one
INSERT INTO hrms.employee_credential_events (
    tenant_id, employee_id, user_id, event_type, delivery_channel, delivery_target, status, failure_reason, created_by
) VALUES (
    @tenant_id, @employee_id, @user_id, @event_type, @delivery_channel, @delivery_target, @status, @failure_reason, @created_by
)
RETURNING *;

-- name: ListEmployeeCredentialEvents :many
SELECT *
FROM hrms.employee_credential_events
WHERE tenant_id = @tenant_id AND employee_id = @employee_id
ORDER BY created_at DESC
LIMIT sqlc.arg('limit_rows');
