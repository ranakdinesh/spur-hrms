-- name: UpsertPrivacyConsent :one
INSERT INTO hrms.privacy_consents (
    tenant_id, employee_user_id, worker_profile_id, consent_key, consent_area, status,
    lawful_basis, channel, source, purpose, granted_at, revoked_at, expires_at,
    evidence, metadata, created_by, updated_by
) VALUES (
    $1, sqlc.narg('employee_user_id')::uuid, sqlc.narg('worker_profile_id')::uuid, $2, $3, $4,
    $5, $6, $7, $8, sqlc.narg('granted_at')::timestamptz, sqlc.narg('revoked_at')::timestamptz,
    sqlc.narg('expires_at')::timestamptz, $9, $10, sqlc.narg('actor_id')::uuid, sqlc.narg('actor_id')::uuid
)
ON CONFLICT (tenant_id, consent_key) WHERE NOT inactive DO UPDATE SET
    employee_user_id = EXCLUDED.employee_user_id,
    worker_profile_id = EXCLUDED.worker_profile_id,
    consent_area = EXCLUDED.consent_area,
    status = EXCLUDED.status,
    lawful_basis = EXCLUDED.lawful_basis,
    channel = EXCLUDED.channel,
    source = EXCLUDED.source,
    purpose = EXCLUDED.purpose,
    granted_at = EXCLUDED.granted_at,
    revoked_at = EXCLUDED.revoked_at,
    expires_at = EXCLUDED.expires_at,
    evidence = EXCLUDED.evidence,
    metadata = EXCLUDED.metadata,
    inactive = FALSE,
    updated_at = NOW(),
    updated_by = EXCLUDED.updated_by
RETURNING *;

-- name: ListPrivacyConsents :many
SELECT *
FROM hrms.privacy_consents
WHERE tenant_id = $1
  AND NOT inactive
  AND (sqlc.narg('status')::text IS NULL OR status = sqlc.narg('status')::text)
  AND (sqlc.narg('consent_area')::text IS NULL OR consent_area = sqlc.narg('consent_area')::text)
ORDER BY updated_at DESC
LIMIT $2 OFFSET $3;

-- name: CreateDataErasureRequest :one
INSERT INTO hrms.data_erasure_requests (
    tenant_id, request_key, subject_user_id, worker_profile_id, request_type, status,
    priority, requested_by, reason, scope, retained_reason, due_at, audit_summary,
    created_by, updated_by
) VALUES (
    $1, $2, sqlc.narg('subject_user_id')::uuid, sqlc.narg('worker_profile_id')::uuid, $3, $4,
    $5, sqlc.narg('requested_by')::uuid, $6, $7, sqlc.narg('retained_reason')::text,
    sqlc.narg('due_at')::timestamptz, $8, sqlc.narg('actor_id')::uuid, sqlc.narg('actor_id')::uuid
)
RETURNING *;

-- name: ListDataErasureRequests :many
SELECT *
FROM hrms.data_erasure_requests
WHERE tenant_id = $1
  AND NOT inactive
  AND (sqlc.narg('status')::text IS NULL OR status = sqlc.narg('status')::text)
  AND (sqlc.narg('priority')::text IS NULL OR priority = sqlc.narg('priority')::text)
ORDER BY
  CASE priority WHEN 'urgent' THEN 1 WHEN 'high' THEN 2 WHEN 'normal' THEN 3 ELSE 4 END,
  COALESCE(due_at, created_at) ASC
LIMIT $2 OFFSET $3;

-- name: UpdateDataErasureRequestStatus :one
UPDATE hrms.data_erasure_requests
SET status = $3,
    retained_reason = COALESCE(sqlc.narg('retained_reason')::text, retained_reason),
    audit_summary = COALESCE(sqlc.narg('audit_summary')::jsonb, audit_summary),
    reviewed_by = COALESCE(sqlc.narg('actor_id')::uuid, reviewed_by),
    reviewed_at = CASE WHEN $3 IN ('approved', 'rejected', 'blocked_legal_hold', 'completed', 'cancelled') THEN NOW() ELSE reviewed_at END,
    completed_at = CASE WHEN $3 = 'completed' THEN NOW() ELSE completed_at END,
    updated_at = NOW(),
    updated_by = sqlc.narg('actor_id')::uuid
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: UpsertEcosystemIntegrationHook :one
INSERT INTO hrms.ecosystem_integration_hooks (
    tenant_id, hook_key, provider, channel, direction, status, display_name,
    endpoint_url, event_types, secret_ref, consent_required, mobile_safe,
    last_checked_at, last_error, config, created_by, updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7,
    sqlc.narg('endpoint_url')::text, $8, sqlc.narg('secret_ref')::text, $9, $10,
    sqlc.narg('last_checked_at')::timestamptz, sqlc.narg('last_error')::text, $11,
    sqlc.narg('actor_id')::uuid, sqlc.narg('actor_id')::uuid
)
ON CONFLICT (tenant_id, hook_key) WHERE NOT inactive DO UPDATE SET
    provider = EXCLUDED.provider,
    channel = EXCLUDED.channel,
    direction = EXCLUDED.direction,
    status = EXCLUDED.status,
    display_name = EXCLUDED.display_name,
    endpoint_url = EXCLUDED.endpoint_url,
    event_types = EXCLUDED.event_types,
    secret_ref = EXCLUDED.secret_ref,
    consent_required = EXCLUDED.consent_required,
    mobile_safe = EXCLUDED.mobile_safe,
    last_checked_at = EXCLUDED.last_checked_at,
    last_error = EXCLUDED.last_error,
    config = EXCLUDED.config,
    inactive = FALSE,
    updated_at = NOW(),
    updated_by = EXCLUDED.updated_by
RETURNING *;

-- name: ListEcosystemIntegrationHooks :many
SELECT *
FROM hrms.ecosystem_integration_hooks
WHERE tenant_id = $1
  AND NOT inactive
  AND (sqlc.narg('channel')::text IS NULL OR channel = sqlc.narg('channel')::text)
  AND (sqlc.narg('status')::text IS NULL OR status = sqlc.narg('status')::text)
ORDER BY channel, display_name
LIMIT $2 OFFSET $3;

-- name: UpsertMobileAPIConstraint :one
INSERT INTO hrms.mobile_api_constraints (
    tenant_id, constraint_key, workflow, min_android_version, min_ios_version,
    offline_supported, low_bandwidth_mode, requires_location, requires_device_binding,
    max_payload_kb, status, notes, config, created_by, updated_by
) VALUES (
    $1, $2, $3, sqlc.narg('min_android_version')::text, sqlc.narg('min_ios_version')::text,
    $4, $5, $6, $7, $8, $9, sqlc.narg('notes')::text, $10,
    sqlc.narg('actor_id')::uuid, sqlc.narg('actor_id')::uuid
)
ON CONFLICT (tenant_id, constraint_key) WHERE NOT inactive DO UPDATE SET
    workflow = EXCLUDED.workflow,
    min_android_version = EXCLUDED.min_android_version,
    min_ios_version = EXCLUDED.min_ios_version,
    offline_supported = EXCLUDED.offline_supported,
    low_bandwidth_mode = EXCLUDED.low_bandwidth_mode,
    requires_location = EXCLUDED.requires_location,
    requires_device_binding = EXCLUDED.requires_device_binding,
    max_payload_kb = EXCLUDED.max_payload_kb,
    status = EXCLUDED.status,
    notes = EXCLUDED.notes,
    config = EXCLUDED.config,
    inactive = FALSE,
    updated_at = NOW(),
    updated_by = EXCLUDED.updated_by
RETURNING *;

-- name: ListMobileAPIConstraints :many
SELECT *
FROM hrms.mobile_api_constraints
WHERE tenant_id = $1
  AND NOT inactive
  AND (sqlc.narg('workflow')::text IS NULL OR workflow = sqlc.narg('workflow')::text)
  AND (sqlc.narg('status')::text IS NULL OR status = sqlc.narg('status')::text)
ORDER BY workflow, constraint_key
LIMIT $2 OFFSET $3;
