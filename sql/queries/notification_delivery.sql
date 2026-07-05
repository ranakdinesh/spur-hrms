-- name: ListNotificationMasters :many
SELECT * FROM hrms.notification_masters
WHERE tenant_id = $1 AND NOT inactive
ORDER BY code ASC;

-- name: CreateNotificationMaster :one
INSERT INTO hrms.notification_masters (
    tenant_id, code, name, description, is_in_app_enabled, is_email_enabled, is_push_enabled,
    email_subject_template, email_text_template, email_html_template, created_by
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING *;

-- name: GetNotificationMaster :one
SELECT * FROM hrms.notification_masters
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: GetNotificationMasterByCode :one
SELECT * FROM hrms.notification_masters
WHERE tenant_id = $1 AND code = $2 AND NOT inactive;

-- name: UpdateNotificationMaster :one
UPDATE hrms.notification_masters
SET code = $3,
    name = $4,
    description = $5,
    is_in_app_enabled = $6,
    is_email_enabled = $7,
    is_push_enabled = $8,
    email_subject_template = $9,
    email_text_template = $10,
    email_html_template = $11,
    updated_at = NOW(),
    updated_by = $12
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: SoftDeleteNotificationMaster :exec
UPDATE hrms.notification_masters
SET inactive = TRUE, updated_at = NOW(), updated_by = $3
WHERE tenant_id = $1 AND id = $2;

-- name: ListNotificationPreferencesByUser :many
SELECT * FROM hrms.notification_preferences
WHERE tenant_id = $1 AND user_id = $2 AND NOT inactive
ORDER BY created_at DESC;

-- name: UpsertNotificationPreference :one
INSERT INTO hrms.notification_preferences (
    tenant_id, user_id, notification_master_id, is_in_app_enabled, is_email_enabled, is_push_enabled,
    digest_frequency, quiet_hours_start, quiet_hours_end, timezone, created_by
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
ON CONFLICT (tenant_id, user_id, notification_master_id) WHERE NOT inactive
DO UPDATE SET
    is_in_app_enabled = EXCLUDED.is_in_app_enabled,
    is_email_enabled = EXCLUDED.is_email_enabled,
    is_push_enabled = EXCLUDED.is_push_enabled,
    digest_frequency = EXCLUDED.digest_frequency,
    quiet_hours_start = EXCLUDED.quiet_hours_start,
    quiet_hours_end = EXCLUDED.quiet_hours_end,
    timezone = EXCLUDED.timezone,
    updated_at = NOW(),
    updated_by = $11
RETURNING *;

-- name: GetNotificationPreference :one
SELECT * FROM hrms.notification_preferences
WHERE tenant_id = $1 AND user_id = $2 AND notification_master_id = $3 AND NOT inactive;

-- name: SoftDeleteNotificationPreference :exec
UPDATE hrms.notification_preferences
SET inactive = TRUE, updated_at = NOW(), updated_by = $3
WHERE tenant_id = $1 AND id = $2;

-- name: ListNotificationInboxByUser :many
SELECT
    i.id,
    i.tenant_id,
    i.user_id,
    i.notification_master_id,
    m.code AS notification_code,
    m.name AS notification_name,
    i.title,
    i.message,
    i.reference_table,
    i.reference_id,
    i.is_read,
    i.read_date,
    i.inactive,
    i.created_at,
    i.created_by,
    i.updated_at,
    i.updated_by
FROM hrms.notification_inbox i
JOIN hrms.notification_masters m ON m.tenant_id = i.tenant_id AND m.id = i.notification_master_id
WHERE i.tenant_id = $1
  AND i.user_id = $2
  AND (sqlc.narg('is_read')::boolean IS NULL OR i.is_read = sqlc.narg('is_read')::boolean)
  AND (sqlc.narg('notification_master_id')::uuid IS NULL OR i.notification_master_id = sqlc.narg('notification_master_id')::uuid)
  AND (sqlc.narg('notification_code')::text IS NULL OR m.code = sqlc.narg('notification_code')::text)
  AND (
      sqlc.narg('search')::text IS NULL
      OR i.title ILIKE '%' || sqlc.narg('search')::text || '%'
      OR i.message ILIKE '%' || sqlc.narg('search')::text || '%'
  )
  AND NOT i.inactive
ORDER BY i.is_read ASC, i.created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: CountNotificationInboxByUser :one
SELECT COUNT(*) FROM hrms.notification_inbox i
JOIN hrms.notification_masters m ON m.tenant_id = i.tenant_id AND m.id = i.notification_master_id
WHERE i.tenant_id = $1
  AND i.user_id = $2
  AND (sqlc.narg('is_read')::boolean IS NULL OR i.is_read = sqlc.narg('is_read')::boolean)
  AND (sqlc.narg('notification_master_id')::uuid IS NULL OR i.notification_master_id = sqlc.narg('notification_master_id')::uuid)
  AND (sqlc.narg('notification_code')::text IS NULL OR m.code = sqlc.narg('notification_code')::text)
  AND (
      sqlc.narg('search')::text IS NULL
      OR i.title ILIKE '%' || sqlc.narg('search')::text || '%'
      OR i.message ILIKE '%' || sqlc.narg('search')::text || '%'
  )
  AND NOT i.inactive;

-- name: CreateNotificationInboxItem :one
INSERT INTO hrms.notification_inbox (
    tenant_id, user_id, notification_master_id, title, message, reference_table, reference_id, created_by
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: CountUnreadNotificationsByUser :one
SELECT COUNT(*) FROM hrms.notification_inbox
WHERE tenant_id = $1 AND user_id = $2 AND NOT is_read AND NOT inactive;

-- name: GetNotificationInboxItem :one
SELECT * FROM hrms.notification_inbox
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: MarkNotificationInboxItemRead :exec
UPDATE hrms.notification_inbox
SET is_read = TRUE, read_date = NOW(), updated_at = NOW(), updated_by = $3
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: MarkNotificationInboxItemUnread :exec
UPDATE hrms.notification_inbox
SET is_read = FALSE, read_date = NULL, updated_at = NOW(), updated_by = $3
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: MarkNotificationInboxAllRead :exec
UPDATE hrms.notification_inbox
SET is_read = TRUE,
    read_date = COALESCE(read_date, NOW()),
    updated_at = NOW(),
    updated_by = $3
WHERE tenant_id = $1 AND user_id = $2 AND NOT is_read AND NOT inactive;

-- name: SoftDeleteNotificationInboxItem :exec
UPDATE hrms.notification_inbox
SET inactive = TRUE, updated_at = NOW(), updated_by = $3
WHERE tenant_id = $1 AND id = $2;

-- name: ListNotificationLogs :many
SELECT * FROM hrms.notification_logs
WHERE tenant_id = $1
  AND (sqlc.narg('user_id')::uuid IS NULL OR user_id = sqlc.narg('user_id')::uuid)
  AND (sqlc.narg('status')::text IS NULL OR status = sqlc.narg('status')::text)
  AND (sqlc.narg('channel')::text IS NULL OR channel = sqlc.narg('channel')::text)
  AND (sqlc.narg('bulk_id')::uuid IS NULL OR bulk_id = sqlc.narg('bulk_id')::uuid)
  AND NOT inactive
ORDER BY created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: ListNotificationLogsByUser :many
SELECT * FROM hrms.notification_logs
WHERE tenant_id = $1 AND user_id = $2 AND NOT inactive
ORDER BY created_at DESC
LIMIT $3 OFFSET $4;

-- name: CreateNotificationLog :one
INSERT INTO hrms.notification_logs (
    tenant_id, notification_master_id, user_id, channel, target_address, subject, body, status,
    sent_date, error_message, external_reference_id, idempotency_key, bulk_id, created_by
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8,
    CASE WHEN $8 = 'Sent' THEN NOW() ELSE NULL END,
    $9, $10, $11, $12, $13
)
ON CONFLICT (tenant_id, idempotency_key) WHERE idempotency_key IS NOT NULL AND NOT inactive
DO UPDATE SET updated_at = NOW()
RETURNING *;

-- name: GetNotificationLog :one
SELECT * FROM hrms.notification_logs
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: SoftDeleteNotificationLog :exec
UPDATE hrms.notification_logs
SET inactive = TRUE, updated_by = $3
WHERE tenant_id = $1 AND id = $2;

-- name: ListDeviceTokensByUser :many
SELECT * FROM hrms.device_tokens
WHERE tenant_id = $1 AND user_id = $2 AND NOT inactive
ORDER BY updated_at DESC;

-- name: UpsertDeviceToken :one
INSERT INTO hrms.device_tokens (
    tenant_id, user_id, device_token, device_type, device_id, created_by
)
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (tenant_id, user_id, device_id) WHERE device_id IS NOT NULL AND NOT inactive
DO UPDATE SET
    device_token = EXCLUDED.device_token,
    device_type = EXCLUDED.device_type,
    inactive = FALSE,
    updated_at = NOW(),
    updated_by = $6
RETURNING *;

-- name: DeactivateRotatedDeviceTokens :exec
UPDATE hrms.device_tokens
SET inactive = TRUE,
    updated_at = NOW(),
    updated_by = $4
WHERE hrms.device_tokens.tenant_id = $1
  AND hrms.device_tokens.user_id = $2
  AND hrms.device_tokens.device_token <> $3
  AND NOT inactive
  AND (
      hrms.device_tokens.device_id IS NULL
      OR hrms.device_tokens.device_id NOT IN (
          SELECT current_tokens.device_id
          FROM hrms.device_tokens AS current_tokens
          WHERE current_tokens.tenant_id = $1
            AND current_tokens.user_id = $2
            AND current_tokens.device_token = $3
            AND current_tokens.device_id IS NOT NULL
            AND NOT current_tokens.inactive
      )
  );

-- name: ListDeviceTokensByDeviceID :many
SELECT * FROM hrms.device_tokens
WHERE tenant_id = $1 AND device_id = $2 AND NOT inactive
ORDER BY updated_at DESC;

-- name: GetDeviceToken :one
SELECT * FROM hrms.device_tokens
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: SoftDeleteDeviceToken :exec
UPDATE hrms.device_tokens
SET inactive = TRUE, updated_at = NOW(), updated_by = $3
WHERE tenant_id = $1 AND id = $2;
