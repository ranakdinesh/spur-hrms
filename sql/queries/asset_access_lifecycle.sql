-- name: CreateAssetItem :one
INSERT INTO hrms.asset_items (
    tenant_id, asset_code, asset_name, asset_type, category, serial_number, vendor, purchase_date,
    warranty_until, owner_user_id, custodian_worker_profile_id, status, location_label, notes, metadata,
    created_by, updated_by
)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$16)
RETURNING *;

-- name: UpdateAssetItem :one
UPDATE hrms.asset_items
SET asset_code = $3,
    asset_name = $4,
    asset_type = $5,
    category = $6,
    serial_number = $7,
    vendor = $8,
    purchase_date = $9,
    warranty_until = $10,
    owner_user_id = $11,
    custodian_worker_profile_id = $12,
    status = $13,
    location_label = $14,
    notes = $15,
    metadata = $16,
    updated_by = $17,
    updated_at = now()
WHERE tenant_id = $1 AND id = $2 AND inactive = FALSE
RETURNING *;

-- name: UpdateAssetItemStatus :one
UPDATE hrms.asset_items
SET status = $3, updated_by = $4, updated_at = now()
WHERE tenant_id = $1 AND id = $2 AND inactive = FALSE
RETURNING *;

-- name: GetAssetItem :one
SELECT * FROM hrms.asset_items
WHERE tenant_id = $1 AND id = $2 AND inactive = FALSE;

-- name: ListAssetItems :many
SELECT ai.*, wp.display_name AS custodian_name, wp.worker_code AS custodian_code
FROM hrms.asset_items ai
LEFT JOIN hrms.worker_profiles wp ON wp.tenant_id = ai.tenant_id AND wp.id = ai.custodian_worker_profile_id AND NOT wp.inactive
WHERE ai.tenant_id = $1
  AND ai.inactive = FALSE
  AND (sqlc.narg('status')::text IS NULL OR ai.status = sqlc.narg('status')::text)
  AND (sqlc.narg('category')::text IS NULL OR ai.category = sqlc.narg('category')::text)
  AND (sqlc.narg('search')::text IS NULL OR ai.asset_code ILIKE '%' || sqlc.narg('search')::text || '%' OR ai.asset_name ILIKE '%' || sqlc.narg('search')::text || '%' OR COALESCE(ai.serial_number, '') ILIKE '%' || sqlc.narg('search')::text || '%')
ORDER BY ai.updated_at DESC
LIMIT $2 OFFSET $3;

-- name: SoftDeleteAssetItem :exec
UPDATE hrms.asset_items
SET inactive = TRUE, updated_by = $3, updated_at = now()
WHERE tenant_id = $1 AND id = $2;

-- name: CreateAccessCatalogItem :one
INSERT INTO hrms.access_catalog_items (
    tenant_id, access_code, access_name, access_type, system_name, owner_user_id,
    provisioning_method, requires_approval, default_for_onboarding, default_for_exit_revocation,
    status, notes, metadata, created_by, updated_by
)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$14)
RETURNING *;

-- name: UpdateAccessCatalogItem :one
UPDATE hrms.access_catalog_items
SET access_code = $3,
    access_name = $4,
    access_type = $5,
    system_name = $6,
    owner_user_id = $7,
    provisioning_method = $8,
    requires_approval = $9,
    default_for_onboarding = $10,
    default_for_exit_revocation = $11,
    status = $12,
    notes = $13,
    metadata = $14,
    updated_by = $15,
    updated_at = now()
WHERE tenant_id = $1 AND id = $2 AND inactive = FALSE
RETURNING *;

-- name: ListAccessCatalogItems :many
SELECT * FROM hrms.access_catalog_items
WHERE tenant_id = $1
  AND inactive = FALSE
  AND (sqlc.narg('status')::text IS NULL OR status = sqlc.narg('status')::text)
  AND (sqlc.narg('access_type')::text IS NULL OR access_type = sqlc.narg('access_type')::text)
  AND (sqlc.narg('search')::text IS NULL OR access_code ILIKE '%' || sqlc.narg('search')::text || '%' OR access_name ILIKE '%' || sqlc.narg('search')::text || '%' OR COALESCE(system_name, '') ILIKE '%' || sqlc.narg('search')::text || '%')
ORDER BY access_name ASC
LIMIT $2 OFFSET $3;

-- name: GetAccessCatalogItem :one
SELECT * FROM hrms.access_catalog_items
WHERE tenant_id = $1 AND id = $2 AND inactive = FALSE;

-- name: SoftDeleteAccessCatalogItem :exec
UPDATE hrms.access_catalog_items
SET inactive = TRUE, updated_by = $3, updated_at = now()
WHERE tenant_id = $1 AND id = $2;

-- name: CreateAssetAssignment :one
INSERT INTO hrms.asset_assignments (
    tenant_id, asset_id, worker_profile_id, employee_id, candidate_onboarding_id, exit_request_id,
    requested_by, approved_by, issued_by, returned_by, approved_at, issued_on, expected_return_on,
    returned_on, issue_condition, return_condition, damage_status, recovery_amount, status, notes,
    metadata, created_by, updated_by
)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$22)
RETURNING *;

-- name: UpdateAssetAssignment :one
UPDATE hrms.asset_assignments
SET worker_profile_id = $3,
    employee_id = $4,
    candidate_onboarding_id = $5,
    exit_request_id = $6,
    requested_by = $7,
    approved_by = $8,
    issued_by = $9,
    returned_by = $10,
    approved_at = $11,
    issued_on = $12,
    expected_return_on = $13,
    returned_on = $14,
    issue_condition = $15,
    return_condition = $16,
    damage_status = $17,
    recovery_amount = $18,
    notes = $19,
    metadata = $20,
    updated_by = $21,
    updated_at = now()
WHERE tenant_id = $1 AND id = $2 AND inactive = FALSE
RETURNING *;

-- name: UpdateAssetAssignmentStatus :one
UPDATE hrms.asset_assignments
SET status = $3,
    approved_by = CASE WHEN $3 = 'approved' THEN $4 ELSE approved_by END,
    approved_at = CASE WHEN $3 = 'approved' THEN now() ELSE approved_at END,
    issued_by = CASE WHEN $3 = 'issued' THEN $4 ELSE issued_by END,
    issued_on = CASE WHEN $3 = 'issued' AND issued_on IS NULL THEN CURRENT_DATE ELSE issued_on END,
    returned_by = CASE WHEN $3 = 'returned' THEN $4 ELSE returned_by END,
    returned_on = CASE WHEN $3 = 'returned' AND returned_on IS NULL THEN CURRENT_DATE ELSE returned_on END,
    updated_by = $4,
    updated_at = now()
WHERE tenant_id = $1 AND id = $2 AND inactive = FALSE
RETURNING *;

-- name: ListAssetAssignments :many
SELECT aa.*, ai.asset_code, ai.asset_name, ai.asset_type, ai.category,
       wp.display_name AS worker_display_name, wp.worker_code
FROM hrms.asset_assignments aa
JOIN hrms.asset_items ai ON ai.tenant_id = aa.tenant_id AND ai.id = aa.asset_id
JOIN hrms.worker_profiles wp ON wp.tenant_id = aa.tenant_id AND wp.id = aa.worker_profile_id AND NOT wp.inactive
WHERE aa.tenant_id = $1
  AND aa.inactive = FALSE
  AND (sqlc.narg('worker_profile_id')::uuid IS NULL OR aa.worker_profile_id = sqlc.narg('worker_profile_id')::uuid)
  AND (sqlc.narg('asset_id')::uuid IS NULL OR aa.asset_id = sqlc.narg('asset_id')::uuid)
  AND (sqlc.narg('exit_request_id')::uuid IS NULL OR aa.exit_request_id = sqlc.narg('exit_request_id')::uuid)
  AND (sqlc.narg('status')::text IS NULL OR aa.status = sqlc.narg('status')::text)
  AND (sqlc.narg('search')::text IS NULL OR ai.asset_code ILIKE '%' || sqlc.narg('search')::text || '%' OR ai.asset_name ILIKE '%' || sqlc.narg('search')::text || '%' OR wp.display_name ILIKE '%' || sqlc.narg('search')::text || '%')
ORDER BY aa.updated_at DESC
LIMIT $2 OFFSET $3;

-- name: CreateAccessLifecycleTask :one
INSERT INTO hrms.access_lifecycle_tasks (
    tenant_id, access_item_id, worker_profile_id, employee_id, candidate_onboarding_id, exit_request_id,
    task_type, requested_by, approved_by, owner_user_id, approved_at, due_date, completed_at,
    external_reference, status, notes, metadata, created_by, updated_by
)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$18)
RETURNING *;

-- name: UpdateAccessLifecycleTask :one
UPDATE hrms.access_lifecycle_tasks
SET worker_profile_id = $3,
    employee_id = $4,
    candidate_onboarding_id = $5,
    exit_request_id = $6,
    task_type = $7,
    requested_by = $8,
    approved_by = $9,
    owner_user_id = $10,
    approved_at = $11,
    due_date = $12,
    completed_at = $13,
    external_reference = $14,
    notes = $15,
    metadata = $16,
    updated_by = $17,
    updated_at = now()
WHERE tenant_id = $1 AND id = $2 AND inactive = FALSE
RETURNING *;

-- name: UpdateAccessLifecycleTaskStatus :one
UPDATE hrms.access_lifecycle_tasks
SET status = $3,
    approved_by = CASE WHEN $3 = 'approved' THEN $4 ELSE approved_by END,
    approved_at = CASE WHEN $3 = 'approved' THEN now() ELSE approved_at END,
    completed_at = CASE WHEN $3 IN ('provisioned','revoked','reviewed') THEN now() ELSE completed_at END,
    updated_by = $4,
    updated_at = now()
WHERE tenant_id = $1 AND id = $2 AND inactive = FALSE
RETURNING *;

-- name: ListAccessLifecycleTasks :many
SELECT alt.*, ac.access_code, ac.access_name, ac.access_type, ac.system_name,
       wp.display_name AS worker_display_name, wp.worker_code
FROM hrms.access_lifecycle_tasks alt
JOIN hrms.access_catalog_items ac ON ac.tenant_id = alt.tenant_id AND ac.id = alt.access_item_id
JOIN hrms.worker_profiles wp ON wp.tenant_id = alt.tenant_id AND wp.id = alt.worker_profile_id AND NOT wp.inactive
WHERE alt.tenant_id = $1
  AND alt.inactive = FALSE
  AND (sqlc.narg('worker_profile_id')::uuid IS NULL OR alt.worker_profile_id = sqlc.narg('worker_profile_id')::uuid)
  AND (sqlc.narg('access_item_id')::uuid IS NULL OR alt.access_item_id = sqlc.narg('access_item_id')::uuid)
  AND (sqlc.narg('exit_request_id')::uuid IS NULL OR alt.exit_request_id = sqlc.narg('exit_request_id')::uuid)
  AND (sqlc.narg('status')::text IS NULL OR alt.status = sqlc.narg('status')::text)
  AND (sqlc.narg('search')::text IS NULL OR ac.access_code ILIKE '%' || sqlc.narg('search')::text || '%' OR ac.access_name ILIKE '%' || sqlc.narg('search')::text || '%' OR wp.display_name ILIKE '%' || sqlc.narg('search')::text || '%')
ORDER BY alt.due_date ASC NULLS LAST, alt.updated_at DESC
LIMIT $2 OFFSET $3;

-- name: CreateAssetAccessEvent :one
INSERT INTO hrms.asset_access_events (tenant_id, source_type, source_id, action, from_status, to_status, remarks, metadata, created_by, updated_by)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$9)
RETURNING *;

-- name: ListAssetAccessEvents :many
SELECT * FROM hrms.asset_access_events
WHERE tenant_id = $1
  AND inactive = FALSE
  AND (sqlc.narg('source_type')::text IS NULL OR source_type = sqlc.narg('source_type')::text)
  AND (sqlc.narg('source_id')::uuid IS NULL OR source_id = sqlc.narg('source_id')::uuid)
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetAssetAccessSummary :many
SELECT 'assets_total'::text AS metric, COUNT(*)::bigint AS metric_count
FROM hrms.asset_items ai
WHERE ai.tenant_id = $1 AND ai.inactive = FALSE
UNION ALL
SELECT 'assets_issued', COUNT(*)::bigint
FROM hrms.asset_items ai
WHERE ai.tenant_id = $1 AND ai.inactive = FALSE AND ai.status IN ('issued','return_due')
UNION ALL
SELECT 'assets_due_return', COUNT(*)::bigint
FROM hrms.asset_assignments aa
WHERE aa.tenant_id = $1 AND aa.inactive = FALSE AND aa.status IN ('issued','return_due') AND aa.expected_return_on IS NOT NULL AND aa.expected_return_on <= CURRENT_DATE
UNION ALL
SELECT 'open_access_tasks', COUNT(*)::bigint
FROM hrms.access_lifecycle_tasks alt
WHERE alt.tenant_id = $1 AND alt.inactive = FALSE AND alt.status IN ('requested','approved','blocked')
UNION ALL
SELECT 'revocation_tasks', COUNT(*)::bigint
FROM hrms.access_lifecycle_tasks alt
WHERE alt.tenant_id = $1 AND alt.inactive = FALSE AND alt.task_type = 'deprovision' AND alt.status NOT IN ('revoked','cancelled','rejected');
