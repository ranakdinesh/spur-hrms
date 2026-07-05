-- name: CreatePulseSurvey :one
INSERT INTO hrms.pulse_surveys (
    tenant_id, title, description, survey_type, status, audience_scope,
    department_id, start_date, end_date, frequency, anonymity_threshold,
    consent_required, manager_aggregate_only, critical_alerts_enabled,
    metadata, created_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, COALESCE($15, '{}'::jsonb), $16
) RETURNING *;

-- name: UpdatePulseSurvey :one
UPDATE hrms.pulse_surveys
SET title = $3,
    description = $4,
    survey_type = $5,
    status = $6,
    audience_scope = $7,
    department_id = $8,
    start_date = $9,
    end_date = $10,
    frequency = $11,
    anonymity_threshold = $12,
    consent_required = $13,
    manager_aggregate_only = $14,
    critical_alerts_enabled = $15,
    metadata = COALESCE($16, '{}'::jsonb),
    updated_by = $17
WHERE tenant_id = $1 AND id = $2 AND inactive = FALSE
RETURNING *;

-- name: UpdatePulseSurveyStatus :one
UPDATE hrms.pulse_surveys
SET status = $3,
    updated_by = $4
WHERE tenant_id = $1 AND id = $2 AND inactive = FALSE
RETURNING *;

-- name: GetPulseSurvey :one
SELECT ps.*,
       dept.name AS department_name,
       COUNT(q.id)::int AS question_count,
       COUNT(DISTINCT r.worker_profile_id)::int AS respondent_count,
       COUNT(r.id)::int AS response_count
FROM hrms.pulse_surveys ps
LEFT JOIN hrms.departments dept ON dept.tenant_id = ps.tenant_id AND dept.id = ps.department_id AND NOT dept.inactive
LEFT JOIN hrms.pulse_survey_questions q ON q.tenant_id = ps.tenant_id AND q.survey_id = ps.id AND NOT q.inactive
LEFT JOIN hrms.pulse_responses r ON r.tenant_id = ps.tenant_id AND r.survey_id = ps.id AND NOT r.inactive
WHERE ps.tenant_id = $1 AND ps.id = $2 AND ps.inactive = FALSE
GROUP BY ps.id, dept.name;

-- name: ListPulseSurveys :many
SELECT ps.*,
       dept.name AS department_name,
       COUNT(q.id)::int AS question_count,
       COUNT(DISTINCT r.worker_profile_id)::int AS respondent_count,
       COUNT(r.id)::int AS response_count
FROM hrms.pulse_surveys ps
LEFT JOIN hrms.departments dept ON dept.tenant_id = ps.tenant_id AND dept.id = ps.department_id AND NOT dept.inactive
LEFT JOIN hrms.pulse_survey_questions q ON q.tenant_id = ps.tenant_id AND q.survey_id = ps.id AND NOT q.inactive
LEFT JOIN hrms.pulse_responses r ON r.tenant_id = ps.tenant_id AND r.survey_id = ps.id AND NOT r.inactive
WHERE ps.tenant_id = $1
  AND ps.inactive = FALSE
  AND ($2::text = '' OR ps.status = $2)
  AND ($3::text = '' OR ps.survey_type = $3)
  AND ($4::text = '' OR ps.title ILIKE '%' || $4 || '%')
GROUP BY ps.id, dept.name
ORDER BY ps.start_date DESC, ps.created_at DESC;

-- name: SoftDeletePulseSurvey :exec
UPDATE hrms.pulse_surveys
SET inactive = TRUE, updated_by = $3
WHERE tenant_id = $1 AND id = $2 AND inactive = FALSE;

-- name: CreatePulseQuestion :one
INSERT INTO hrms.pulse_survey_questions (
    tenant_id, survey_id, question_text, question_type, category,
    is_required, sort_order, options, created_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, COALESCE($8, '[]'::jsonb), $9
) RETURNING *;

-- name: UpdatePulseQuestion :one
UPDATE hrms.pulse_survey_questions
SET question_text = $3,
    question_type = $4,
    category = $5,
    is_required = $6,
    sort_order = $7,
    options = COALESCE($8, '[]'::jsonb),
    updated_by = $9
WHERE tenant_id = $1 AND id = $2 AND inactive = FALSE
RETURNING *;

-- name: GetPulseQuestion :one
SELECT * FROM hrms.pulse_survey_questions
WHERE tenant_id = $1 AND id = $2 AND inactive = FALSE;

-- name: ListPulseQuestions :many
SELECT *
FROM hrms.pulse_survey_questions
WHERE tenant_id = $1
  AND inactive = FALSE
  AND ($2::uuid IS NULL OR survey_id = $2)
ORDER BY sort_order ASC, created_at ASC;

-- name: SoftDeletePulseQuestion :exec
UPDATE hrms.pulse_survey_questions
SET inactive = TRUE, updated_by = $3
WHERE tenant_id = $1 AND id = $2 AND inactive = FALSE;

-- name: CreatePulseResponse :one
INSERT INTO hrms.pulse_responses (
    tenant_id, survey_id, question_id, worker_profile_id, response_date,
    score, text_response, boolean_response, option_value, consent_given,
    is_anonymous, risk_level, critical_alert, metadata, created_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, COALESCE($14, '{}'::jsonb), $15
) RETURNING *;

-- name: ListPulseResponses :many
SELECT r.*,
       ps.title AS survey_title,
       q.question_text,
       q.category,
       worker.display_name AS worker_display_name,
       worker.worker_code AS worker_code
FROM hrms.pulse_responses r
JOIN hrms.pulse_surveys ps ON ps.tenant_id = r.tenant_id AND ps.id = r.survey_id AND NOT ps.inactive
JOIN hrms.pulse_survey_questions q ON q.tenant_id = r.tenant_id AND q.id = r.question_id AND NOT q.inactive
LEFT JOIN hrms.worker_profiles worker ON worker.tenant_id = r.tenant_id AND worker.id = r.worker_profile_id AND NOT worker.inactive
WHERE r.tenant_id = $1
  AND r.inactive = FALSE
  AND ($2::uuid IS NULL OR r.survey_id = $2)
  AND ($3::uuid IS NULL OR r.worker_profile_id = $3)
  AND ($4::text = '' OR r.risk_level = $4)
ORDER BY r.response_date DESC, r.created_at DESC;

-- name: UpsertWellbeingScore :one
INSERT INTO hrms.wellbeing_scores (
    tenant_id, worker_profile_id, score_date, source_survey_id,
    wellbeing_score, mood_score, stress_score, workload_score,
    risk_level, consent_scope, notes, metadata, created_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, COALESCE($12, '{}'::jsonb), $13
)
ON CONFLICT (tenant_id, worker_profile_id, score_date) WHERE inactive = FALSE
DO UPDATE SET source_survey_id = EXCLUDED.source_survey_id,
              wellbeing_score = EXCLUDED.wellbeing_score,
              mood_score = EXCLUDED.mood_score,
              stress_score = EXCLUDED.stress_score,
              workload_score = EXCLUDED.workload_score,
              risk_level = EXCLUDED.risk_level,
              consent_scope = EXCLUDED.consent_scope,
              notes = EXCLUDED.notes,
              metadata = EXCLUDED.metadata,
              updated_by = EXCLUDED.created_by
RETURNING *;

-- name: ListWellbeingScores :many
SELECT ws.*,
       worker.display_name AS worker_display_name,
       worker.worker_code AS worker_code,
       survey.title AS survey_title
FROM hrms.wellbeing_scores ws
JOIN hrms.worker_profiles worker ON worker.tenant_id = ws.tenant_id AND worker.id = ws.worker_profile_id AND NOT worker.inactive
LEFT JOIN hrms.pulse_surveys survey ON survey.tenant_id = ws.tenant_id AND survey.id = ws.source_survey_id AND NOT survey.inactive
WHERE ws.tenant_id = $1
  AND ws.inactive = FALSE
  AND ($2::uuid IS NULL OR ws.worker_profile_id = $2)
  AND ($3::text = '' OR ws.risk_level = $3)
ORDER BY ws.score_date DESC, ws.updated_at DESC;

-- name: CreateWellbeingAlert :one
INSERT INTO hrms.wellbeing_alerts (
    tenant_id, worker_profile_id, survey_id, response_id, alert_type,
    severity, status, visible_to, message, metadata, created_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, 'hr', $8, COALESCE($9, '{}'::jsonb), $10
) RETURNING *;

-- name: UpdateWellbeingAlertStatus :one
UPDATE hrms.wellbeing_alerts
SET status = $3,
    resolution_note = $4,
    updated_by = $5
WHERE tenant_id = $1 AND id = $2 AND inactive = FALSE
RETURNING *;

-- name: ListWellbeingAlerts :many
SELECT alert.*,
       worker.display_name AS worker_display_name,
       worker.worker_code AS worker_code,
       survey.title AS survey_title
FROM hrms.wellbeing_alerts alert
LEFT JOIN hrms.worker_profiles worker ON worker.tenant_id = alert.tenant_id AND worker.id = alert.worker_profile_id AND NOT worker.inactive
LEFT JOIN hrms.pulse_surveys survey ON survey.tenant_id = alert.tenant_id AND survey.id = alert.survey_id AND NOT survey.inactive
WHERE alert.tenant_id = $1
  AND alert.inactive = FALSE
  AND ($2::text = '' OR alert.status = $2)
  AND ($3::text = '' OR alert.severity = $3)
ORDER BY alert.created_at DESC;

-- name: ListWellbeingAggregateRows :many
SELECT ps.id AS survey_id,
       ps.title AS survey_title,
       COALESCE(dept.id, '00000000-0000-0000-0000-000000000000'::uuid) AS department_id,
       COALESCE(dept.name, 'All / Unassigned') AS department_name,
       q.category,
       COUNT(r.id)::int AS response_count,
       COUNT(DISTINCT r.worker_profile_id)::int AS respondent_count,
       (COUNT(DISTINCT r.worker_profile_id) < ps.anonymity_threshold)::bool AS suppressed,
       (CASE WHEN COUNT(DISTINCT r.worker_profile_id) >= ps.anonymity_threshold THEN ROUND(AVG(r.score)::numeric, 2) ELSE NULL END)::numeric AS average_score,
       CASE WHEN COUNT(DISTINCT r.worker_profile_id) >= ps.anonymity_threshold THEN COUNT(r.id) FILTER (WHERE r.risk_level IN ('high', 'critical'))::int ELSE 0 END AS risk_count,
       ps.anonymity_threshold
FROM hrms.pulse_surveys ps
JOIN hrms.pulse_survey_questions q ON q.tenant_id = ps.tenant_id AND q.survey_id = ps.id AND NOT q.inactive
LEFT JOIN hrms.pulse_responses r ON r.tenant_id = ps.tenant_id AND r.survey_id = ps.id AND r.question_id = q.id AND NOT r.inactive
LEFT JOIN hrms.worker_profiles worker ON worker.tenant_id = ps.tenant_id AND worker.id = r.worker_profile_id AND NOT worker.inactive
LEFT JOIN hrms.departments dept ON dept.tenant_id = ps.tenant_id AND dept.id = worker.department_id AND NOT dept.inactive
WHERE ps.tenant_id = $1
  AND ps.inactive = FALSE
  AND ($2::uuid IS NULL OR ps.id = $2)
GROUP BY ps.id, ps.title, ps.anonymity_threshold, dept.id, dept.name, q.category
ORDER BY ps.title, dept.name, q.category;
