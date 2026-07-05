-- name: UpsertInsight :one
INSERT INTO hrms.insights (
    id, tenant_id, insight_key, insight_type, category, severity, status, title, summary,
    confidence_score, score, source, model_version, entity_type, entity_id, employee_user_id,
    reasons, recommendations, context, explainability, detected_at, due_at, assigned_to,
    reviewed_by, reviewed_at, resolved_at, resolution_note, created_by, updated_by
) VALUES (
    COALESCE(sqlc.narg('id')::uuid, gen_random_uuid()),
    $1, $2, $3, $4, $5, $6, $7, $8,
    $9, $10, $11, sqlc.narg('model_version')::text, sqlc.narg('entity_type')::text,
    sqlc.narg('entity_id')::uuid, sqlc.narg('employee_user_id')::uuid,
    $12, $13, $14, $15, $16, sqlc.narg('due_at')::timestamptz, sqlc.narg('assigned_to')::uuid,
    sqlc.narg('reviewed_by')::uuid, sqlc.narg('reviewed_at')::timestamptz,
    sqlc.narg('resolved_at')::timestamptz, sqlc.narg('resolution_note')::text, sqlc.narg('actor_id')::uuid, sqlc.narg('actor_id')::uuid
)
ON CONFLICT (tenant_id, insight_key) DO UPDATE SET
    insight_type = EXCLUDED.insight_type,
    category = EXCLUDED.category,
    severity = EXCLUDED.severity,
    status = CASE WHEN hrms.insights.status IN ('resolved', 'dismissed', 'overridden') THEN hrms.insights.status ELSE EXCLUDED.status END,
    title = EXCLUDED.title,
    summary = EXCLUDED.summary,
    confidence_score = EXCLUDED.confidence_score,
    score = EXCLUDED.score,
    source = EXCLUDED.source,
    model_version = EXCLUDED.model_version,
    entity_type = EXCLUDED.entity_type,
    entity_id = EXCLUDED.entity_id,
    employee_user_id = EXCLUDED.employee_user_id,
    reasons = EXCLUDED.reasons,
    recommendations = EXCLUDED.recommendations,
    context = EXCLUDED.context,
    explainability = EXCLUDED.explainability,
    detected_at = EXCLUDED.detected_at,
    due_at = EXCLUDED.due_at,
    updated_at = NOW(),
    updated_by = EXCLUDED.updated_by,
    inactive = FALSE
RETURNING *;

-- name: ListInsights :many
SELECT *
FROM hrms.insights
WHERE tenant_id = $1
  AND NOT inactive
  AND (sqlc.narg('status')::text IS NULL OR status = sqlc.narg('status')::text)
  AND (sqlc.narg('severity')::text IS NULL OR severity = sqlc.narg('severity')::text)
  AND (sqlc.narg('category')::text IS NULL OR category = sqlc.narg('category')::text)
  AND (sqlc.narg('insight_type')::text IS NULL OR insight_type = sqlc.narg('insight_type')::text)
  AND (sqlc.narg('assigned_to')::uuid IS NULL OR assigned_to = sqlc.narg('assigned_to')::uuid)
ORDER BY
  CASE severity WHEN 'critical' THEN 1 WHEN 'high' THEN 2 WHEN 'medium' THEN 3 ELSE 4 END,
  detected_at DESC
LIMIT $2 OFFSET $3;

-- name: GetInsight :one
SELECT *
FROM hrms.insights
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: UpdateInsightStatus :one
UPDATE hrms.insights
SET status = $3,
    assigned_to = COALESCE(sqlc.narg('assigned_to')::uuid, assigned_to),
    reviewed_by = COALESCE(sqlc.narg('reviewed_by')::uuid, reviewed_by),
    reviewed_at = CASE WHEN $3 IN ('reviewing', 'dismissed', 'overridden', 'resolved') THEN COALESCE(sqlc.narg('reviewed_at')::timestamptz, NOW()) ELSE reviewed_at END,
    resolved_at = CASE WHEN $3 IN ('resolved', 'dismissed', 'overridden') THEN COALESCE(sqlc.narg('resolved_at')::timestamptz, NOW()) ELSE NULL END,
    resolution_note = sqlc.narg('resolution_note')::text,
    updated_at = NOW(),
    updated_by = sqlc.narg('actor_id')::uuid
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: CreateInsightEvent :one
INSERT INTO hrms.insight_events (
    tenant_id, insight_id, action, from_status, to_status, remarks, metadata, created_by, updated_by
) VALUES (
    $1, $2, $3, sqlc.narg('from_status')::text, sqlc.narg('to_status')::text,
    sqlc.narg('remarks')::text, $4, sqlc.narg('actor_id')::uuid, sqlc.narg('actor_id')::uuid
)
RETURNING *;

-- name: ListInsightEvents :many
SELECT *
FROM hrms.insight_events
WHERE tenant_id = $1 AND insight_id = $2 AND NOT inactive
ORDER BY created_at DESC;
