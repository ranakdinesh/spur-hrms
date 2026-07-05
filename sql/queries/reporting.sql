-- name: UpsertReportCatalogItem :one
INSERT INTO hrms.report_catalog (
    id, tenant_id, report_code, module, name, description, category, scope,
    permission_key, default_filters, supported_filters, output_columns,
    drilldown_contract, is_system, is_active, sort_order, created_by, updated_by
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $17)
ON CONFLICT (tenant_id, report_code) WHERE NOT inactive
DO UPDATE SET
    module = EXCLUDED.module,
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    category = EXCLUDED.category,
    scope = EXCLUDED.scope,
    permission_key = EXCLUDED.permission_key,
    default_filters = EXCLUDED.default_filters,
    supported_filters = EXCLUDED.supported_filters,
    output_columns = EXCLUDED.output_columns,
    drilldown_contract = EXCLUDED.drilldown_contract,
    is_system = EXCLUDED.is_system,
    is_active = EXCLUDED.is_active,
    sort_order = EXCLUDED.sort_order,
    updated_by = EXCLUDED.updated_by,
    updated_at = NOW()
RETURNING *;

-- name: ListReportCatalog :many
SELECT * FROM hrms.report_catalog
WHERE tenant_id = $1
  AND (sqlc.narg('module')::text IS NULL OR module = sqlc.narg('module')::text)
  AND (sqlc.narg('scope')::text IS NULL OR scope = sqlc.narg('scope')::text)
  AND NOT inactive
ORDER BY category ASC, sort_order ASC, name ASC;

-- name: GetReportCatalogItem :one
SELECT * FROM hrms.report_catalog
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: UpsertReportSavedView :one
INSERT INTO hrms.report_saved_views (
    id, tenant_id, report_id, name, description, visibility, filters, columns,
    is_favorite, owner_user_id, created_by, updated_by
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $11)
ON CONFLICT (id)
DO UPDATE SET
    report_id = EXCLUDED.report_id,
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    visibility = EXCLUDED.visibility,
    filters = EXCLUDED.filters,
    columns = EXCLUDED.columns,
    is_favorite = EXCLUDED.is_favorite,
    owner_user_id = EXCLUDED.owner_user_id,
    inactive = FALSE,
    updated_by = EXCLUDED.updated_by,
    updated_at = NOW()
WHERE hrms.report_saved_views.tenant_id = EXCLUDED.tenant_id
RETURNING *;

-- name: ListReportSavedViews :many
SELECT * FROM hrms.report_saved_views
WHERE tenant_id = $1
  AND (sqlc.narg('report_id')::uuid IS NULL OR report_id = sqlc.narg('report_id')::uuid)
  AND NOT inactive
ORDER BY is_favorite DESC, updated_at DESC, name ASC;

-- name: SoftDeleteReportSavedView :exec
UPDATE hrms.report_saved_views
SET inactive = TRUE, updated_by = $3, updated_at = NOW()
WHERE tenant_id = $1 AND id = $2;

-- name: CreateReportExportJob :one
INSERT INTO hrms.report_export_jobs (
    id, tenant_id, report_id, saved_view_id, export_format, status, filters,
    file_object_key, error_message, requested_by, started_at, completed_at, created_by, updated_by
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $10, $10)
RETURNING *;

-- name: ListReportExportJobs :many
SELECT * FROM hrms.report_export_jobs
WHERE tenant_id = $1
  AND (sqlc.narg('report_id')::uuid IS NULL OR report_id = sqlc.narg('report_id')::uuid)
  AND (sqlc.narg('status')::text IS NULL OR status = sqlc.narg('status')::text)
  AND NOT inactive
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateReportExportJobStatus :one
UPDATE hrms.report_export_jobs
SET status = $3,
    file_object_key = $4,
    error_message = $5,
    started_at = CASE WHEN $3 = 'running' AND started_at IS NULL THEN NOW() ELSE started_at END,
    completed_at = CASE WHEN $3 IN ('completed','failed') THEN NOW() ELSE completed_at END,
    updated_by = $6,
    updated_at = NOW()
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: UpsertReportSchedule :one
INSERT INTO hrms.report_schedules (
    id, tenant_id, report_id, saved_view_id, name, frequency, timezone,
    delivery_channels, recipient_user_ids, recipient_emails, next_run_at,
    last_run_at, is_active, created_by, updated_by
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $14)
ON CONFLICT (id)
DO UPDATE SET
    report_id = EXCLUDED.report_id,
    saved_view_id = EXCLUDED.saved_view_id,
    name = EXCLUDED.name,
    frequency = EXCLUDED.frequency,
    timezone = EXCLUDED.timezone,
    delivery_channels = EXCLUDED.delivery_channels,
    recipient_user_ids = EXCLUDED.recipient_user_ids,
    recipient_emails = EXCLUDED.recipient_emails,
    next_run_at = EXCLUDED.next_run_at,
    last_run_at = EXCLUDED.last_run_at,
    is_active = EXCLUDED.is_active,
    inactive = FALSE,
    updated_by = EXCLUDED.updated_by,
    updated_at = NOW()
WHERE hrms.report_schedules.tenant_id = EXCLUDED.tenant_id
RETURNING *;

-- name: ListReportSchedules :many
SELECT * FROM hrms.report_schedules
WHERE tenant_id = $1
  AND (sqlc.narg('report_id')::uuid IS NULL OR report_id = sqlc.narg('report_id')::uuid)
  AND NOT inactive
ORDER BY is_active DESC, next_run_at ASC NULLS LAST, name ASC;

-- name: SoftDeleteReportSchedule :exec
UPDATE hrms.report_schedules
SET inactive = TRUE, updated_by = $3, updated_at = NOW()
WHERE tenant_id = $1 AND id = $2;

-- name: CreateReportSnapshot :one
INSERT INTO hrms.report_snapshots (
    id, tenant_id, report_id, saved_view_id, snapshot_key, period_start,
    period_end, filters, summary, row_count, generated_at, generated_by, created_by, updated_by
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, COALESCE($11, NOW()), $12, $12, $12)
RETURNING *;

-- name: ListReportSnapshots :many
SELECT * FROM hrms.report_snapshots
WHERE tenant_id = $1
  AND (sqlc.narg('report_id')::uuid IS NULL OR report_id = sqlc.narg('report_id')::uuid)
  AND NOT inactive
ORDER BY period_end DESC, generated_at DESC
LIMIT $2 OFFSET $3;
