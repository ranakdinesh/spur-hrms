-- Hrms sqlc queries
-- Run: sqlc generate  after adding queries

-- name: CreateHrms :one
INSERT INTO hrmss (id, tenant_id)
VALUES ($1, $2)
RETURNING *;

-- name: GetHrmsByID :one
SELECT * FROM hrmss
WHERE id = $1 AND tenant_id = $2;

-- name: ListHrmssByTenant :many
SELECT * FROM hrmss
WHERE tenant_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: DeleteHrms :exec
DELETE FROM hrmss
WHERE id = $1 AND tenant_id = $2;
