-- name: CreateWorkflowDefinition :one
INSERT INTO hrms.workflow_definitions (
    tenant_id, workflow_key, name, module_key, description, status, visibility_scope, metadata, created_by, updated_by
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$9)
RETURNING *;

-- name: UpdateWorkflowDefinition :one
UPDATE hrms.workflow_definitions
SET workflow_key = $3,
    name = $4,
    module_key = $5,
    description = $6,
    status = $7,
    visibility_scope = $8,
    metadata = $9,
    updated_by = $10,
    updated_at = now()
WHERE tenant_id = $1 AND id = $2 AND inactive = FALSE
RETURNING *;

-- name: GetWorkflowDefinition :one
SELECT * FROM hrms.workflow_definitions
WHERE tenant_id = $1 AND id = $2 AND inactive = FALSE;

-- name: ListWorkflowDefinitions :many
SELECT * FROM hrms.workflow_definitions
WHERE tenant_id = $1
  AND inactive = FALSE
  AND (sqlc.narg('status')::text IS NULL OR status = sqlc.narg('status')::text)
  AND (sqlc.narg('module_key')::text IS NULL OR module_key = sqlc.narg('module_key')::text)
  AND (sqlc.narg('search')::text IS NULL OR workflow_key ILIKE '%' || sqlc.narg('search')::text || '%' OR name ILIKE '%' || sqlc.narg('search')::text || '%')
ORDER BY updated_at DESC
LIMIT $2 OFFSET $3;

-- name: CreateWorkflowDefinitionStep :one
INSERT INTO hrms.workflow_definition_steps (
    tenant_id, workflow_definition_id, step_order, step_key, name, step_type, assignment_type,
    assignment_value, required, due_offset_hours, allowed_actions, metadata, created_by, updated_by
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$13)
RETURNING *;

-- name: UpdateWorkflowDefinitionStep :one
UPDATE hrms.workflow_definition_steps
SET step_order = $3,
    step_key = $4,
    name = $5,
    step_type = $6,
    assignment_type = $7,
    assignment_value = $8,
    required = $9,
    due_offset_hours = $10,
    allowed_actions = $11,
    metadata = $12,
    updated_by = $13,
    updated_at = now()
WHERE tenant_id = $1 AND id = $2 AND inactive = FALSE
RETURNING *;

-- name: ListWorkflowDefinitionSteps :many
SELECT * FROM hrms.workflow_definition_steps
WHERE tenant_id = $1 AND workflow_definition_id = $2 AND inactive = FALSE
ORDER BY step_order ASC, created_at ASC;

-- name: CreateOperationTemplate :one
INSERT INTO hrms.operation_templates (
    tenant_id, template_key, name, category, source_module, source_type, workflow_definition_id,
    default_priority, default_severity, allowed_actions, launch_schema, is_active, metadata, created_by, updated_by
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$14)
RETURNING *;

-- name: UpdateOperationTemplate :one
UPDATE hrms.operation_templates
SET template_key = $3,
    name = $4,
    category = $5,
    source_module = $6,
    source_type = $7,
    workflow_definition_id = $8,
    default_priority = $9,
    default_severity = $10,
    allowed_actions = $11,
    launch_schema = $12,
    is_active = $13,
    metadata = $14,
    updated_by = $15,
    updated_at = now()
WHERE tenant_id = $1 AND id = $2 AND inactive = FALSE
RETURNING *;

-- name: GetOperationTemplate :one
SELECT * FROM hrms.operation_templates
WHERE tenant_id = $1 AND id = $2 AND inactive = FALSE;

-- name: ListOperationTemplates :many
SELECT * FROM hrms.operation_templates
WHERE tenant_id = $1
  AND inactive = FALSE
  AND (sqlc.narg('category')::text IS NULL OR category = sqlc.narg('category')::text)
  AND (sqlc.narg('source_module')::text IS NULL OR source_module = sqlc.narg('source_module')::text)
  AND (sqlc.narg('active_only')::boolean IS NULL OR is_active = sqlc.narg('active_only')::boolean)
  AND (sqlc.narg('search')::text IS NULL OR template_key ILIKE '%' || sqlc.narg('search')::text || '%' OR name ILIKE '%' || sqlc.narg('search')::text || '%')
ORDER BY category ASC, name ASC
LIMIT $2 OFFSET $3;

-- name: CreateWorkflowTask :one
INSERT INTO hrms.workflow_tasks (
    tenant_id, task_number, template_id, workflow_definition_id, workflow_step_id, parent_task_id,
    source_module, source_type, source_id, source_record_label, title, description, requester_user_id,
    assignee_user_id, assignee_role, assignee_team, delegated_from_user_id, status, priority,
    severity, visibility_scope, due_at, action_schema, metadata, created_by, updated_by
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$25)
RETURNING *;

-- name: UpdateWorkflowTask :one
UPDATE hrms.workflow_tasks
SET title = $3,
    description = $4,
    assignee_user_id = $5,
    assignee_role = $6,
    assignee_team = $7,
    priority = $8,
    severity = $9,
    visibility_scope = $10,
    due_at = $11,
    action_schema = $12,
    metadata = $13,
    updated_by = $14,
    updated_at = now()
WHERE tenant_id = $1 AND id = $2 AND inactive = FALSE
RETURNING *;

-- name: UpdateWorkflowTaskStatus :one
UPDATE hrms.workflow_tasks
SET status = $3,
    assignee_user_id = COALESCE($4, assignee_user_id),
    assignee_role = COALESCE($5, assignee_role),
    assignee_team = COALESCE($6, assignee_team),
    delegated_from_user_id = $7,
    completed_at = CASE WHEN $3 IN ('approved','rejected','completed','cancelled') THEN now() ELSE completed_at END,
    completed_by = CASE WHEN $3 IN ('approved','rejected','completed','cancelled') THEN $8 ELSE completed_by END,
    updated_by = $8,
    updated_at = now()
WHERE tenant_id = $1 AND id = $2 AND inactive = FALSE
RETURNING *;

-- name: GetWorkflowTask :one
SELECT * FROM hrms.workflow_tasks
WHERE tenant_id = $1 AND id = $2 AND inactive = FALSE;

-- name: ListWorkflowTasks :many
SELECT wt.*,
       (SELECT COUNT(*)::int FROM hrms.workflow_task_comments c WHERE c.tenant_id = wt.tenant_id AND c.task_id = wt.id AND c.inactive = FALSE) AS comment_count,
       (SELECT COUNT(*)::int FROM hrms.workflow_task_attachments a WHERE a.tenant_id = wt.tenant_id AND a.task_id = wt.id AND a.inactive = FALSE) AS attachment_count,
       EXISTS (
           SELECT 1 FROM hrms.workflow_task_watchers w
           WHERE w.tenant_id = wt.tenant_id AND w.task_id = wt.id AND w.watcher_user_id = sqlc.narg('viewer_user_id')::uuid AND w.inactive = FALSE
       ) AS watched_by_viewer
FROM hrms.workflow_tasks wt
WHERE wt.tenant_id = $1
  AND wt.inactive = FALSE
  AND (sqlc.narg('status')::text IS NULL OR wt.status = sqlc.narg('status')::text)
  AND (sqlc.narg('severity')::text IS NULL OR wt.severity = sqlc.narg('severity')::text)
  AND (sqlc.narg('source_module')::text IS NULL OR wt.source_module = sqlc.narg('source_module')::text)
  AND (sqlc.narg('search')::text IS NULL OR wt.task_number ILIKE '%' || sqlc.narg('search')::text || '%' OR wt.title ILIKE '%' || sqlc.narg('search')::text || '%' OR COALESCE(wt.source_record_label, '') ILIKE '%' || sqlc.narg('search')::text || '%')
  AND (
      sqlc.narg('view_key')::text IS NULL OR sqlc.narg('view_key')::text = 'all'
      OR (sqlc.narg('view_key')::text = 'inbox' AND wt.status NOT IN ('approved','rejected','completed','cancelled') AND (wt.assignee_user_id = sqlc.narg('viewer_user_id')::uuid OR sqlc.narg('viewer_user_id')::uuid IS NULL OR wt.assignee_role = sqlc.narg('viewer_role')::text OR wt.assignee_team = sqlc.narg('viewer_team')::text))
      OR (sqlc.narg('view_key')::text = 'my_requests' AND wt.requester_user_id = sqlc.narg('viewer_user_id')::uuid)
      OR (sqlc.narg('view_key')::text = 'team' AND (wt.assignee_role = sqlc.narg('viewer_role')::text OR wt.assignee_team = sqlc.narg('viewer_team')::text))
      OR (sqlc.narg('view_key')::text = 'watching' AND EXISTS (SELECT 1 FROM hrms.workflow_task_watchers w WHERE w.tenant_id = wt.tenant_id AND w.task_id = wt.id AND w.watcher_user_id = sqlc.narg('viewer_user_id')::uuid AND w.inactive = FALSE))
      OR (sqlc.narg('view_key')::text = 'completed' AND wt.status IN ('approved','rejected','completed','cancelled'))
      OR (sqlc.narg('view_key')::text = 'delegated' AND (wt.delegated_from_user_id = sqlc.narg('viewer_user_id')::uuid OR wt.status = 'delegated'))
  )
ORDER BY
  CASE WHEN wt.status IN ('approved','rejected','completed','cancelled') THEN 1 ELSE 0 END,
  wt.due_at ASC NULLS LAST,
  wt.priority ASC,
  wt.updated_at DESC
LIMIT $2 OFFSET $3;

-- name: CreateWorkflowTaskWatcher :one
INSERT INTO hrms.workflow_task_watchers (tenant_id, task_id, watcher_user_id, watch_reason, created_by, updated_by)
VALUES ($1,$2,$3,$4,$5,$5)
ON CONFLICT (tenant_id, task_id, watcher_user_id) WHERE inactive = FALSE DO UPDATE SET
    watch_reason = EXCLUDED.watch_reason,
    updated_by = EXCLUDED.updated_by,
    updated_at = now()
RETURNING *;

-- name: RemoveWorkflowTaskWatcher :exec
UPDATE hrms.workflow_task_watchers
SET inactive = TRUE, updated_by = $4, updated_at = now()
WHERE tenant_id = $1 AND task_id = $2 AND watcher_user_id = $3 AND inactive = FALSE;

-- name: ListWorkflowTaskWatchers :many
SELECT * FROM hrms.workflow_task_watchers
WHERE tenant_id = $1 AND task_id = $2 AND inactive = FALSE
ORDER BY created_at ASC;

-- name: CreateWorkflowTaskComment :one
INSERT INTO hrms.workflow_task_comments (tenant_id, task_id, visibility, body, metadata, created_by, updated_by)
VALUES ($1,$2,$3,$4,$5,$6,$6)
RETURNING *;

-- name: ListWorkflowTaskComments :many
SELECT * FROM hrms.workflow_task_comments
WHERE tenant_id = $1 AND task_id = $2 AND inactive = FALSE
ORDER BY created_at ASC;

-- name: CreateWorkflowTaskAttachment :one
INSERT INTO hrms.workflow_task_attachments (
    tenant_id, task_id, comment_id, file_name, content_type, storage_path, checksum_sha256,
    size_bytes, visibility, metadata, created_by, updated_by
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$11)
RETURNING *;

-- name: ListWorkflowTaskAttachments :many
SELECT * FROM hrms.workflow_task_attachments
WHERE tenant_id = $1 AND task_id = $2 AND inactive = FALSE
ORDER BY created_at DESC;

-- name: CreateWorkflowTaskEvent :one
INSERT INTO hrms.workflow_task_events (
    tenant_id, task_id, action, from_status, to_status, actor_user_id, remarks, metadata, created_by, updated_by
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$9)
RETURNING *;

-- name: ListWorkflowTaskEvents :many
SELECT * FROM hrms.workflow_task_events
WHERE tenant_id = $1 AND task_id = $2 AND inactive = FALSE
ORDER BY created_at ASC;

-- name: GetWorkflowTaskSummary :many
SELECT 'open_total'::text AS metric, COUNT(*)::bigint AS metric_count
FROM hrms.workflow_tasks wt
WHERE wt.tenant_id = $1 AND wt.inactive = FALSE AND wt.status NOT IN ('approved','rejected','completed','cancelled')
UNION ALL
SELECT 'overdue', COUNT(*)::bigint
FROM hrms.workflow_tasks wt
WHERE wt.tenant_id = $1 AND wt.inactive = FALSE AND wt.status NOT IN ('approved','rejected','completed','cancelled') AND wt.due_at IS NOT NULL AND wt.due_at < now()
UNION ALL
SELECT 'completed', COUNT(*)::bigint
FROM hrms.workflow_tasks wt
WHERE wt.tenant_id = $1 AND wt.inactive = FALSE AND wt.status IN ('approved','rejected','completed','cancelled')
UNION ALL
SELECT 'delegated', COUNT(*)::bigint
FROM hrms.workflow_tasks wt
WHERE wt.tenant_id = $1 AND wt.inactive = FALSE AND wt.status = 'delegated'
UNION ALL
SELECT 'confidential', COUNT(*)::bigint
FROM hrms.workflow_tasks wt
WHERE wt.tenant_id = $1 AND wt.inactive = FALSE AND wt.visibility_scope = 'confidential';
