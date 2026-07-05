-- name: UpsertAISignalLog :one
INSERT INTO hrms.ai_signal_logs (
    tenant_id, signal_key, signal_type, source_module, source_event, severity, processing_status,
    entity_type, entity_id, employee_user_id, visibility_scope, idempotency_key, correlation_id,
    payload, explainability, occurred_at, created_by, updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7,
    sqlc.narg('entity_type')::text, sqlc.narg('entity_id')::uuid, sqlc.narg('employee_user_id')::uuid,
    $8, sqlc.narg('idempotency_key')::text, sqlc.narg('correlation_id')::text,
    $9, $10, $11, sqlc.narg('actor_id')::uuid, sqlc.narg('actor_id')::uuid
)
ON CONFLICT (tenant_id, signal_key) WHERE NOT inactive DO UPDATE SET
    signal_type = EXCLUDED.signal_type,
    source_module = EXCLUDED.source_module,
    source_event = EXCLUDED.source_event,
    severity = EXCLUDED.severity,
    processing_status = EXCLUDED.processing_status,
    entity_type = EXCLUDED.entity_type,
    entity_id = EXCLUDED.entity_id,
    employee_user_id = EXCLUDED.employee_user_id,
    visibility_scope = EXCLUDED.visibility_scope,
    idempotency_key = EXCLUDED.idempotency_key,
    correlation_id = EXCLUDED.correlation_id,
    payload = EXCLUDED.payload,
    explainability = EXCLUDED.explainability,
    occurred_at = EXCLUDED.occurred_at,
    updated_at = NOW(),
    updated_by = EXCLUDED.updated_by,
    inactive = FALSE
RETURNING *;

-- name: ListAISignalLogs :many
SELECT *
FROM hrms.ai_signal_logs
WHERE tenant_id = $1
  AND NOT inactive
  AND (sqlc.narg('processing_status')::text IS NULL OR processing_status = sqlc.narg('processing_status')::text)
  AND (sqlc.narg('severity')::text IS NULL OR severity = sqlc.narg('severity')::text)
  AND (sqlc.narg('source_module')::text IS NULL OR source_module = sqlc.narg('source_module')::text)
  AND (sqlc.narg('visibility_scope')::text IS NULL OR visibility_scope = sqlc.narg('visibility_scope')::text)
ORDER BY occurred_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateAISignalStatus :one
UPDATE hrms.ai_signal_logs
SET processing_status = $3,
    processed_at = CASE WHEN $3 IN ('processed', 'ignored', 'failed') THEN NOW() ELSE processed_at END,
    error_message = sqlc.narg('error_message')::text,
    updated_at = NOW(),
    updated_by = sqlc.narg('actor_id')::uuid
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: UpsertAIAgentActionLog :one
INSERT INTO hrms.ai_agent_action_logs (
    tenant_id, action_key, agent_key, agent_name, action_type, status, severity, title, summary,
    insight_id, signal_id, entity_type, entity_id, employee_user_id, visibility_scope,
    proposed_action, input_snapshot, output_snapshot, explainability, confidence_score,
    model_version, sidecar_run_id, requires_human_review, created_by, updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9,
    sqlc.narg('insight_id')::uuid, sqlc.narg('signal_id')::uuid, sqlc.narg('entity_type')::text,
    sqlc.narg('entity_id')::uuid, sqlc.narg('employee_user_id')::uuid, $10,
    $11, $12, $13, $14, $15,
    sqlc.narg('model_version')::text, sqlc.narg('sidecar_run_id')::text, $16,
    sqlc.narg('actor_id')::uuid, sqlc.narg('actor_id')::uuid
)
ON CONFLICT (tenant_id, action_key) WHERE NOT inactive DO UPDATE SET
    agent_key = EXCLUDED.agent_key,
    agent_name = EXCLUDED.agent_name,
    action_type = EXCLUDED.action_type,
    status = CASE WHEN hrms.ai_agent_action_logs.status IN ('approved', 'rejected', 'executed', 'overridden', 'cancelled') THEN hrms.ai_agent_action_logs.status ELSE EXCLUDED.status END,
    severity = EXCLUDED.severity,
    title = EXCLUDED.title,
    summary = EXCLUDED.summary,
    insight_id = EXCLUDED.insight_id,
    signal_id = EXCLUDED.signal_id,
    entity_type = EXCLUDED.entity_type,
    entity_id = EXCLUDED.entity_id,
    employee_user_id = EXCLUDED.employee_user_id,
    visibility_scope = EXCLUDED.visibility_scope,
    proposed_action = EXCLUDED.proposed_action,
    input_snapshot = EXCLUDED.input_snapshot,
    output_snapshot = EXCLUDED.output_snapshot,
    explainability = EXCLUDED.explainability,
    confidence_score = EXCLUDED.confidence_score,
    model_version = EXCLUDED.model_version,
    sidecar_run_id = EXCLUDED.sidecar_run_id,
    requires_human_review = EXCLUDED.requires_human_review,
    updated_at = NOW(),
    updated_by = EXCLUDED.updated_by,
    inactive = FALSE
RETURNING *;

-- name: ListAIAgentActionLogs :many
SELECT *
FROM hrms.ai_agent_action_logs
WHERE tenant_id = $1
  AND NOT inactive
  AND (sqlc.narg('status')::text IS NULL OR status = sqlc.narg('status')::text)
  AND (sqlc.narg('severity')::text IS NULL OR severity = sqlc.narg('severity')::text)
  AND (sqlc.narg('agent_key')::text IS NULL OR agent_key = sqlc.narg('agent_key')::text)
  AND (sqlc.narg('visibility_scope')::text IS NULL OR visibility_scope = sqlc.narg('visibility_scope')::text)
  AND (sqlc.narg('insight_id')::uuid IS NULL OR insight_id = sqlc.narg('insight_id')::uuid)
ORDER BY
  CASE severity WHEN 'critical' THEN 1 WHEN 'high' THEN 2 WHEN 'medium' THEN 3 ELSE 4 END,
  created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetAIAgentActionLog :one
SELECT *
FROM hrms.ai_agent_action_logs
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: UpdateAIAgentActionStatus :one
UPDATE hrms.ai_agent_action_logs
SET status = $3,
    reviewed_by = CASE WHEN $3 IN ('reviewing', 'approved', 'rejected', 'overridden', 'cancelled') THEN COALESCE(sqlc.narg('reviewed_by')::uuid, reviewed_by) ELSE reviewed_by END,
    reviewed_at = CASE WHEN $3 IN ('reviewing', 'approved', 'rejected', 'overridden', 'cancelled') THEN NOW() ELSE reviewed_at END,
    executed_at = CASE WHEN $3 = 'executed' THEN NOW() ELSE executed_at END,
    failed_at = CASE WHEN $3 = 'failed' THEN NOW() ELSE failed_at END,
    failure_message = sqlc.narg('failure_message')::text,
    updated_at = NOW(),
    updated_by = sqlc.narg('actor_id')::uuid
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: CreateAIHumanOverride :one
INSERT INTO hrms.ai_human_overrides (
    tenant_id, insight_id, action_id, override_type, original_status, override_status,
    reason, decision, metadata, created_by, updated_by
) VALUES (
    $1, sqlc.narg('insight_id')::uuid, sqlc.narg('action_id')::uuid, $2,
    sqlc.narg('original_status')::text, $3, $4, $5, $6, sqlc.narg('actor_id')::uuid, sqlc.narg('actor_id')::uuid
)
RETURNING *;

-- name: ListAIHumanOverrides :many
SELECT *
FROM hrms.ai_human_overrides
WHERE tenant_id = $1
  AND NOT inactive
  AND (sqlc.narg('insight_id')::uuid IS NULL OR insight_id = sqlc.narg('insight_id')::uuid)
  AND (sqlc.narg('action_id')::uuid IS NULL OR action_id = sqlc.narg('action_id')::uuid)
  AND (sqlc.narg('decision')::text IS NULL OR decision = sqlc.narg('decision')::text)
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpsertAIEventOutbox :one
INSERT INTO hrms.ai_event_outbox (
    tenant_id, event_key, event_type, target_bus, status, payload, correlation_id, created_by, updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6, sqlc.narg('correlation_id')::text, sqlc.narg('actor_id')::uuid, sqlc.narg('actor_id')::uuid
)
ON CONFLICT (tenant_id, event_key) WHERE NOT inactive DO UPDATE SET
    event_type = EXCLUDED.event_type,
    target_bus = EXCLUDED.target_bus,
    status = EXCLUDED.status,
    payload = EXCLUDED.payload,
    correlation_id = EXCLUDED.correlation_id,
    updated_at = NOW(),
    updated_by = EXCLUDED.updated_by
RETURNING *;

-- name: ListAIEventOutbox :many
SELECT *
FROM hrms.ai_event_outbox
WHERE tenant_id = $1
  AND NOT inactive
  AND (sqlc.narg('status')::text IS NULL OR status = sqlc.narg('status')::text)
  AND (sqlc.narg('event_type')::text IS NULL OR event_type = sqlc.narg('event_type')::text)
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateAIEventOutboxStatus :one
UPDATE hrms.ai_event_outbox
SET status = $3,
    attempts = attempts + 1,
    published_at = CASE WHEN $3 = 'published' THEN NOW() ELSE published_at END,
    last_error = sqlc.narg('last_error')::text,
    updated_at = NOW(),
    updated_by = sqlc.narg('actor_id')::uuid
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;
