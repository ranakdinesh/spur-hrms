-- name: CreateOKRCycle :one
INSERT INTO hrms.okr_cycles (
    tenant_id, name, cycle_code, description, start_date, end_date, status, review_cadence, metadata, created_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, COALESCE($9, '{}'::jsonb), $10
) RETURNING *;

-- name: UpdateOKRCycle :one
UPDATE hrms.okr_cycles
SET name = $3,
    cycle_code = $4,
    description = $5,
    start_date = $6,
    end_date = $7,
    status = $8,
    review_cadence = $9,
    metadata = COALESCE($10, '{}'::jsonb),
    updated_by = $11
WHERE tenant_id = $1 AND id = $2 AND inactive = FALSE
RETURNING *;

-- name: UpdateOKRCycleStatus :one
UPDATE hrms.okr_cycles
SET status = $3,
    updated_by = $4
WHERE tenant_id = $1 AND id = $2 AND inactive = FALSE
RETURNING *;

-- name: GetOKRCycle :one
SELECT * FROM hrms.okr_cycles
WHERE tenant_id = $1 AND id = $2 AND inactive = FALSE;

-- name: ListOKRCycles :many
SELECT *
FROM hrms.okr_cycles
WHERE tenant_id = $1
  AND inactive = FALSE
  AND ($2::text = '' OR status = $2)
  AND ($3::text = '' OR name ILIKE '%' || $3 || '%' OR cycle_code ILIKE '%' || $3 || '%')
ORDER BY start_date DESC, created_at DESC;

-- name: SoftDeleteOKRCycle :exec
UPDATE hrms.okr_cycles
SET inactive = TRUE, updated_by = $3
WHERE tenant_id = $1 AND id = $2 AND inactive = FALSE;

-- name: CreateObjective :one
INSERT INTO hrms.objectives (
    tenant_id, cycle_id, parent_objective_id, owner_type, owner_worker_profile_id,
    owner_department_id, owner_project_id, title, description, status, priority,
    progress_percent, weight, start_date, due_date, metadata, created_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, COALESCE($12, 0), COALESCE($13, 1), $14, $15, COALESCE($16, '{}'::jsonb), $17
) RETURNING *;

-- name: UpdateObjective :one
UPDATE hrms.objectives
SET cycle_id = $3,
    parent_objective_id = $4,
    owner_type = $5,
    owner_worker_profile_id = $6,
    owner_department_id = $7,
    owner_project_id = $8,
    title = $9,
    description = $10,
    status = $11,
    priority = $12,
    progress_percent = COALESCE($13, 0),
    weight = COALESCE($14, 1),
    start_date = $15,
    due_date = $16,
    metadata = COALESCE($17, '{}'::jsonb),
    updated_by = $18
WHERE tenant_id = $1 AND id = $2 AND inactive = FALSE
RETURNING *;

-- name: UpdateObjectiveStatus :one
UPDATE hrms.objectives
SET status = $3,
    updated_by = $4
WHERE tenant_id = $1 AND id = $2 AND inactive = FALSE
RETURNING *;

-- name: RefreshObjectiveProgress :one
UPDATE hrms.objectives obj
SET progress_percent = COALESCE((
        SELECT CASE
            WHEN SUM(kr.weight) FILTER (WHERE NOT kr.inactive) IS NULL THEN obj.progress_percent
            ELSE ROUND((SUM(kr.progress_percent * kr.weight) FILTER (WHERE NOT kr.inactive) / NULLIF(SUM(kr.weight) FILTER (WHERE NOT kr.inactive), 0))::numeric, 2)
        END
        FROM hrms.key_results kr
        WHERE kr.tenant_id = obj.tenant_id AND kr.objective_id = obj.id
    ), obj.progress_percent),
    status = CASE
        WHEN COALESCE((
            SELECT ROUND((SUM(kr.progress_percent * kr.weight) FILTER (WHERE NOT kr.inactive) / NULLIF(SUM(kr.weight) FILTER (WHERE NOT kr.inactive), 0))::numeric, 2)
            FROM hrms.key_results kr
            WHERE kr.tenant_id = obj.tenant_id AND kr.objective_id = obj.id
        ), obj.progress_percent) >= 100 THEN 'completed'
        WHEN EXISTS (SELECT 1 FROM hrms.key_results kr WHERE kr.tenant_id = obj.tenant_id AND kr.objective_id = obj.id AND kr.status IN ('at_risk', 'behind') AND NOT kr.inactive) THEN 'at_risk'
        WHEN obj.status = 'draft' THEN 'draft'
        ELSE 'active'
    END,
    updated_by = $3
WHERE obj.tenant_id = $1 AND obj.id = $2 AND obj.inactive = FALSE
RETURNING obj.*;

-- name: GetObjective :one
SELECT * FROM hrms.objectives
WHERE tenant_id = $1 AND id = $2 AND inactive = FALSE;

-- name: ListObjectives :many
SELECT obj.*,
       cycle.name AS cycle_name,
       parent.title AS parent_objective_title,
       worker.display_name AS owner_worker_name,
       worker.worker_code AS owner_worker_code,
       dept.name AS owner_department_name,
       project.name AS owner_project_name,
       project.project_code AS owner_project_code,
       COUNT(kr.id)::int AS key_result_count,
       COALESCE(AVG(kr.progress_percent) FILTER (WHERE NOT kr.inactive), 0)::numeric(5,2) AS average_key_result_progress
FROM hrms.objectives obj
JOIN hrms.okr_cycles cycle ON cycle.tenant_id = obj.tenant_id AND cycle.id = obj.cycle_id AND NOT cycle.inactive
LEFT JOIN hrms.objectives parent ON parent.tenant_id = obj.tenant_id AND parent.id = obj.parent_objective_id AND NOT parent.inactive
LEFT JOIN hrms.worker_profiles worker ON worker.tenant_id = obj.tenant_id AND worker.id = obj.owner_worker_profile_id AND NOT worker.inactive
LEFT JOIN hrms.departments dept ON dept.tenant_id = obj.tenant_id AND dept.id = obj.owner_department_id AND NOT dept.inactive
LEFT JOIN hrms.projects project ON project.tenant_id = obj.tenant_id AND project.id = obj.owner_project_id AND NOT project.inactive
LEFT JOIN hrms.key_results kr ON kr.tenant_id = obj.tenant_id AND kr.objective_id = obj.id AND NOT kr.inactive
WHERE obj.tenant_id = $1
  AND obj.inactive = FALSE
  AND ($2::uuid IS NULL OR obj.cycle_id = $2)
  AND ($3::text = '' OR obj.owner_type = $3)
  AND ($4::text = '' OR obj.status = $4)
  AND ($5::text = '' OR obj.title ILIKE '%' || $5 || '%' OR cycle.name ILIKE '%' || $5 || '%')
GROUP BY obj.id, cycle.name, parent.title, worker.display_name, worker.worker_code, dept.name, project.name, project.project_code
ORDER BY cycle.start_date DESC, obj.priority DESC, obj.updated_at DESC;

-- name: SoftDeleteObjective :exec
UPDATE hrms.objectives
SET inactive = TRUE, updated_by = $3
WHERE tenant_id = $1 AND id = $2 AND inactive = FALSE;

-- name: CreateKeyResult :one
INSERT INTO hrms.key_results (
    tenant_id, objective_id, title, description, metric_type, start_value,
    target_value, current_value, progress_percent, confidence, status, weight,
    unit_label, due_date, metadata, created_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, COALESCE($9, 0), $10, $11, COALESCE($12, 1), $13, $14, COALESCE($15, '{}'::jsonb), $16
) RETURNING *;

-- name: UpdateKeyResult :one
UPDATE hrms.key_results
SET objective_id = $3,
    title = $4,
    description = $5,
    metric_type = $6,
    start_value = $7,
    target_value = $8,
    current_value = $9,
    progress_percent = COALESCE($10, 0),
    confidence = $11,
    status = $12,
    weight = COALESCE($13, 1),
    unit_label = $14,
    due_date = $15,
    metadata = COALESCE($16, '{}'::jsonb),
    updated_by = $17
WHERE tenant_id = $1 AND id = $2 AND inactive = FALSE
RETURNING *;

-- name: UpdateKeyResultProgress :one
UPDATE hrms.key_results
SET current_value = $3,
    progress_percent = $4,
    confidence = $5,
    status = $6,
    updated_by = $7
WHERE tenant_id = $1 AND id = $2 AND inactive = FALSE
RETURNING *;

-- name: GetKeyResult :one
SELECT * FROM hrms.key_results
WHERE tenant_id = $1 AND id = $2 AND inactive = FALSE;

-- name: ListKeyResults :many
SELECT kr.*,
       obj.title AS objective_title,
       cycle.name AS cycle_name,
       latest.checkin_date AS latest_checkin_date,
       latest.note AS latest_note
FROM hrms.key_results kr
JOIN hrms.objectives obj ON obj.tenant_id = kr.tenant_id AND obj.id = kr.objective_id AND NOT obj.inactive
JOIN hrms.okr_cycles cycle ON cycle.tenant_id = obj.tenant_id AND cycle.id = obj.cycle_id AND NOT cycle.inactive
LEFT JOIN LATERAL (
    SELECT checkin_date, note
    FROM hrms.key_result_checkins c
    WHERE c.tenant_id = kr.tenant_id AND c.key_result_id = kr.id AND NOT c.inactive
    ORDER BY c.checkin_date DESC, c.created_at DESC
    LIMIT 1
) latest ON TRUE
WHERE kr.tenant_id = $1
  AND kr.inactive = FALSE
  AND ($2::uuid IS NULL OR kr.objective_id = $2)
  AND ($3::text = '' OR kr.status = $3)
ORDER BY obj.updated_at DESC, kr.created_at ASC;

-- name: SoftDeleteKeyResult :exec
UPDATE hrms.key_results
SET inactive = TRUE, updated_by = $3
WHERE tenant_id = $1 AND id = $2 AND inactive = FALSE;

-- name: CreateKeyResultCheckIn :one
INSERT INTO hrms.key_result_checkins (
    tenant_id, key_result_id, checkin_date, value, progress_percent, confidence, status, note, metadata, created_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, COALESCE($9, '{}'::jsonb), $10
) RETURNING *;

-- name: ListKeyResultCheckIns :many
SELECT c.*,
       kr.title AS key_result_title,
       obj.id AS objective_id,
       obj.title AS objective_title
FROM hrms.key_result_checkins c
JOIN hrms.key_results kr ON kr.tenant_id = c.tenant_id AND kr.id = c.key_result_id AND NOT kr.inactive
JOIN hrms.objectives obj ON obj.tenant_id = kr.tenant_id AND obj.id = kr.objective_id AND NOT obj.inactive
WHERE c.tenant_id = $1
  AND c.inactive = FALSE
  AND ($2::uuid IS NULL OR c.key_result_id = $2)
  AND ($3::uuid IS NULL OR obj.id = $3)
ORDER BY c.checkin_date DESC, c.created_at DESC;

-- name: GetOKRSummary :many
SELECT obj.owner_type,
       COUNT(DISTINCT obj.id)::int AS objective_count,
       COUNT(kr.id)::int AS key_result_count,
       COALESCE(ROUND(AVG(obj.progress_percent)::numeric, 2), 0)::numeric(5,2) AS average_progress,
       COUNT(DISTINCT obj.id) FILTER (WHERE obj.status = 'at_risk')::int AS at_risk_count,
       COUNT(DISTINCT obj.id) FILTER (WHERE obj.status = 'completed')::int AS completed_count
FROM hrms.objectives obj
LEFT JOIN hrms.key_results kr ON kr.tenant_id = obj.tenant_id AND kr.objective_id = obj.id AND NOT kr.inactive
WHERE obj.tenant_id = $1
  AND obj.inactive = FALSE
  AND ($2::uuid IS NULL OR obj.cycle_id = $2)
GROUP BY obj.owner_type
ORDER BY obj.owner_type;

-- name: CreatePerformanceCheckIn :one
INSERT INTO hrms.performance_checkins (
    tenant_id, worker_profile_id, reviewer_worker_profile_id, cycle_id, checkin_date,
    period_start, period_end, mood, status, visibility, highlights, blockers,
    next_plan, employee_comment, manager_comment, score, calibration_bucket,
    metadata, created_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17,
    COALESCE($18, '{}'::jsonb), $19
) RETURNING *;

-- name: UpdatePerformanceCheckIn :one
UPDATE hrms.performance_checkins
SET worker_profile_id = $3,
    reviewer_worker_profile_id = $4,
    cycle_id = $5,
    checkin_date = $6,
    period_start = $7,
    period_end = $8,
    mood = $9,
    status = $10,
    visibility = $11,
    highlights = $12,
    blockers = $13,
    next_plan = $14,
    employee_comment = $15,
    manager_comment = $16,
    score = $17,
    calibration_bucket = $18,
    metadata = COALESCE($19, '{}'::jsonb),
    updated_by = $20
WHERE tenant_id = $1 AND id = $2 AND inactive = FALSE
RETURNING *;

-- name: ReviewPerformanceCheckIn :one
UPDATE hrms.performance_checkins
SET status = $3,
    manager_comment = $4,
    score = $5,
    calibration_bucket = $6,
    reviewed_at = now(),
    reviewed_by = $7,
    updated_by = $7
WHERE tenant_id = $1 AND id = $2 AND inactive = FALSE
RETURNING *;

-- name: UpdatePerformanceCheckInStatus :one
UPDATE hrms.performance_checkins
SET status = $3,
    updated_by = $4
WHERE tenant_id = $1 AND id = $2 AND inactive = FALSE
RETURNING *;

-- name: GetPerformanceCheckIn :one
SELECT pc.*,
       worker.display_name AS worker_display_name,
       worker.worker_code AS worker_code,
       reviewer.display_name AS reviewer_display_name,
       cycle.name AS cycle_name,
       COALESCE(feedback.feedback_count, 0)::int AS feedback_count,
       COALESCE(feedback.average_feedback_rating, 0)::numeric(5,2) AS average_feedback_rating
FROM hrms.performance_checkins pc
JOIN hrms.worker_profiles worker ON worker.tenant_id = pc.tenant_id AND worker.id = pc.worker_profile_id AND NOT worker.inactive
LEFT JOIN hrms.worker_profiles reviewer ON reviewer.tenant_id = pc.tenant_id AND reviewer.id = pc.reviewer_worker_profile_id AND NOT reviewer.inactive
LEFT JOIN hrms.okr_cycles cycle ON cycle.tenant_id = pc.tenant_id AND cycle.id = pc.cycle_id AND NOT cycle.inactive
LEFT JOIN LATERAL (
    SELECT COUNT(frsp.id)::int AS feedback_count,
           ROUND(AVG(frsp.rating)::numeric, 2) AS average_feedback_rating
    FROM hrms.feedback_requests freq
    JOIN hrms.feedback_responses frsp ON frsp.tenant_id = freq.tenant_id AND frsp.request_id = freq.id AND NOT frsp.inactive
    WHERE freq.tenant_id = pc.tenant_id
      AND freq.subject_worker_profile_id = pc.worker_profile_id
      AND NOT freq.inactive
) feedback ON TRUE
WHERE pc.tenant_id = $1 AND pc.id = $2 AND pc.inactive = FALSE;

-- name: ListPerformanceCheckIns :many
SELECT pc.*,
       worker.display_name AS worker_display_name,
       worker.worker_code AS worker_code,
       reviewer.display_name AS reviewer_display_name,
       cycle.name AS cycle_name,
       COALESCE(feedback.feedback_count, 0)::int AS feedback_count,
       COALESCE(feedback.average_feedback_rating, 0)::numeric(5,2) AS average_feedback_rating
FROM hrms.performance_checkins pc
JOIN hrms.worker_profiles worker ON worker.tenant_id = pc.tenant_id AND worker.id = pc.worker_profile_id AND NOT worker.inactive
LEFT JOIN hrms.worker_profiles reviewer ON reviewer.tenant_id = pc.tenant_id AND reviewer.id = pc.reviewer_worker_profile_id AND NOT reviewer.inactive
LEFT JOIN hrms.okr_cycles cycle ON cycle.tenant_id = pc.tenant_id AND cycle.id = pc.cycle_id AND NOT cycle.inactive
LEFT JOIN LATERAL (
    SELECT COUNT(frsp.id)::int AS feedback_count,
           ROUND(AVG(frsp.rating)::numeric, 2) AS average_feedback_rating
    FROM hrms.feedback_requests freq
    JOIN hrms.feedback_responses frsp ON frsp.tenant_id = freq.tenant_id AND frsp.request_id = freq.id AND NOT frsp.inactive
    WHERE freq.tenant_id = pc.tenant_id
      AND freq.subject_worker_profile_id = pc.worker_profile_id
      AND NOT freq.inactive
) feedback ON TRUE
WHERE pc.tenant_id = $1
  AND pc.inactive = FALSE
  AND ($2::uuid IS NULL OR pc.worker_profile_id = $2)
  AND ($3::uuid IS NULL OR pc.reviewer_worker_profile_id = $3)
  AND ($4::uuid IS NULL OR pc.cycle_id = $4)
  AND ($5::text = '' OR pc.status = $5)
  AND ($6::text = '' OR pc.mood = $6)
ORDER BY pc.checkin_date DESC, pc.updated_at DESC;

-- name: SoftDeletePerformanceCheckIn :exec
UPDATE hrms.performance_checkins
SET inactive = TRUE, updated_by = $3
WHERE tenant_id = $1 AND id = $2 AND inactive = FALSE;

-- name: GetPerformanceCheckInSummary :many
SELECT status,
       mood,
       COUNT(*)::int AS checkin_count,
       COALESCE(ROUND(AVG(score)::numeric, 2), 0)::numeric(5,2) AS average_score
FROM hrms.performance_checkins
WHERE tenant_id = $1
  AND inactive = FALSE
  AND ($2::uuid IS NULL OR cycle_id = $2)
GROUP BY status, mood
ORDER BY status, mood;

-- name: CreateFeedbackRequest :one
INSERT INTO hrms.feedback_requests (
    tenant_id, subject_worker_profile_id, requester_worker_profile_id, objective_id,
    relationship, feedback_type, status, is_anonymous, visibility, due_date, prompt,
    metadata, created_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, COALESCE($12, '{}'::jsonb), $13
) RETURNING *;

-- name: UpdateFeedbackRequest :one
UPDATE hrms.feedback_requests
SET subject_worker_profile_id = $3,
    requester_worker_profile_id = $4,
    objective_id = $5,
    relationship = $6,
    feedback_type = $7,
    status = $8,
    is_anonymous = $9,
    visibility = $10,
    due_date = $11,
    prompt = $12,
    metadata = COALESCE($13, '{}'::jsonb),
    updated_by = $14
WHERE tenant_id = $1 AND id = $2 AND inactive = FALSE
RETURNING *;

-- name: UpdateFeedbackRequestStatus :one
UPDATE hrms.feedback_requests
SET status = $3,
    updated_by = $4
WHERE tenant_id = $1 AND id = $2 AND inactive = FALSE
RETURNING *;

-- name: GetFeedbackRequest :one
SELECT fr.*,
       subject.display_name AS subject_display_name,
       subject.worker_code AS subject_worker_code,
       requester.display_name AS requester_display_name,
       obj.title AS objective_title,
       COALESCE(response_counts.response_count, 0)::int AS response_count
FROM hrms.feedback_requests fr
JOIN hrms.worker_profiles subject ON subject.tenant_id = fr.tenant_id AND subject.id = fr.subject_worker_profile_id AND NOT subject.inactive
LEFT JOIN hrms.worker_profiles requester ON requester.tenant_id = fr.tenant_id AND requester.id = fr.requester_worker_profile_id AND NOT requester.inactive
LEFT JOIN hrms.objectives obj ON obj.tenant_id = fr.tenant_id AND obj.id = fr.objective_id AND NOT obj.inactive
LEFT JOIN LATERAL (
    SELECT COUNT(resp.id)::int AS response_count
    FROM hrms.feedback_responses resp
    WHERE resp.tenant_id = fr.tenant_id AND resp.request_id = fr.id AND NOT resp.inactive
) response_counts ON TRUE
WHERE fr.tenant_id = $1 AND fr.id = $2 AND fr.inactive = FALSE;

-- name: ListFeedbackRequests :many
SELECT fr.*,
       subject.display_name AS subject_display_name,
       subject.worker_code AS subject_worker_code,
       requester.display_name AS requester_display_name,
       obj.title AS objective_title,
       COALESCE(response_counts.response_count, 0)::int AS response_count
FROM hrms.feedback_requests fr
JOIN hrms.worker_profiles subject ON subject.tenant_id = fr.tenant_id AND subject.id = fr.subject_worker_profile_id AND NOT subject.inactive
LEFT JOIN hrms.worker_profiles requester ON requester.tenant_id = fr.tenant_id AND requester.id = fr.requester_worker_profile_id AND NOT requester.inactive
LEFT JOIN hrms.objectives obj ON obj.tenant_id = fr.tenant_id AND obj.id = fr.objective_id AND NOT obj.inactive
LEFT JOIN LATERAL (
    SELECT COUNT(resp.id)::int AS response_count
    FROM hrms.feedback_responses resp
    WHERE resp.tenant_id = fr.tenant_id AND resp.request_id = fr.id AND NOT resp.inactive
) response_counts ON TRUE
WHERE fr.tenant_id = $1
  AND fr.inactive = FALSE
  AND ($2::uuid IS NULL OR fr.subject_worker_profile_id = $2)
  AND ($3::uuid IS NULL OR fr.requester_worker_profile_id = $3)
  AND ($4::text = '' OR fr.status = $4)
  AND ($5::text = '' OR fr.feedback_type = $5)
ORDER BY fr.created_at DESC;

-- name: CreateFeedbackResponse :one
INSERT INTO hrms.feedback_responses (
    tenant_id, request_id, respondent_worker_profile_id, rating, strengths,
    improvements, comments, metadata, created_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, COALESCE($8, '{}'::jsonb), $9
) RETURNING *;

-- name: ListFeedbackResponses :many
SELECT resp.*,
       fr.subject_worker_profile_id,
       fr.is_anonymous,
       subject.display_name AS subject_display_name,
       subject.worker_code AS subject_worker_code,
       (CASE WHEN fr.is_anonymous THEN NULL ELSE respondent.display_name END)::text AS respondent_display_name,
       fr.feedback_type,
       fr.relationship
FROM hrms.feedback_responses resp
JOIN hrms.feedback_requests fr ON fr.tenant_id = resp.tenant_id AND fr.id = resp.request_id AND NOT fr.inactive
JOIN hrms.worker_profiles subject ON subject.tenant_id = resp.tenant_id AND subject.id = fr.subject_worker_profile_id AND NOT subject.inactive
LEFT JOIN hrms.worker_profiles respondent ON respondent.tenant_id = resp.tenant_id AND respondent.id = resp.respondent_worker_profile_id AND NOT respondent.inactive
WHERE resp.tenant_id = $1
  AND resp.inactive = FALSE
  AND ($2::uuid IS NULL OR resp.request_id = $2)
  AND ($3::uuid IS NULL OR fr.subject_worker_profile_id = $3)
  AND ($4::uuid IS NULL OR resp.respondent_worker_profile_id = $4)
ORDER BY resp.submitted_at DESC;

-- name: CreatePerformanceTimelineEvent :one
INSERT INTO hrms.performance_timeline_events (
    tenant_id, worker_profile_id, event_type, checkin_id, feedback_request_id,
    feedback_response_id, objective_id, actor_worker_profile_id, title, notes,
    metadata, created_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, COALESCE($11, '{}'::jsonb), $12
) RETURNING *;

-- name: ListPerformanceTimelineEvents :many
SELECT evt.*,
       worker.display_name AS worker_display_name,
       actor.display_name AS actor_display_name,
       obj.title AS objective_title
FROM hrms.performance_timeline_events evt
JOIN hrms.worker_profiles worker ON worker.tenant_id = evt.tenant_id AND worker.id = evt.worker_profile_id AND NOT worker.inactive
LEFT JOIN hrms.worker_profiles actor ON actor.tenant_id = evt.tenant_id AND actor.id = evt.actor_worker_profile_id AND NOT actor.inactive
LEFT JOIN hrms.objectives obj ON obj.tenant_id = evt.tenant_id AND obj.id = evt.objective_id AND NOT obj.inactive
WHERE evt.tenant_id = $1
  AND evt.inactive = FALSE
  AND ($2::uuid IS NULL OR evt.worker_profile_id = $2)
  AND ($3::text = '' OR evt.event_type = $3)
ORDER BY evt.created_at DESC;

-- name: ListPerformanceCalibrationRows :many
SELECT pc.worker_profile_id,
       worker.display_name AS worker_display_name,
       worker.worker_code AS worker_code,
       pc.cycle_id,
       cycle.name AS cycle_name,
       COUNT(DISTINCT pc.id)::int AS checkin_count,
       COUNT(DISTINCT pc.id) FILTER (WHERE pc.status IN ('submitted', 'reviewed', 'closed'))::int AS submitted_checkin_count,
       COALESCE(ROUND(AVG(pc.score)::numeric, 2), 0)::numeric(5,2) AS average_score,
       COALESCE(MAX(pc.calibration_bucket), '')::text AS calibration_bucket,
       COALESCE(ROUND(AVG(obj.progress_percent)::numeric, 2), 0)::numeric(5,2) AS average_okr_progress,
       COUNT(DISTINCT resp.id)::int AS feedback_count,
       COALESCE(ROUND(AVG(resp.rating)::numeric, 2), 0)::numeric(5,2) AS average_feedback_rating
FROM hrms.performance_checkins pc
JOIN hrms.worker_profiles worker ON worker.tenant_id = pc.tenant_id AND worker.id = pc.worker_profile_id AND NOT worker.inactive
LEFT JOIN hrms.okr_cycles cycle ON cycle.tenant_id = pc.tenant_id AND cycle.id = pc.cycle_id AND NOT cycle.inactive
LEFT JOIN hrms.objectives obj ON obj.tenant_id = pc.tenant_id AND obj.owner_worker_profile_id = pc.worker_profile_id AND NOT obj.inactive
LEFT JOIN hrms.feedback_requests freq ON freq.tenant_id = pc.tenant_id AND freq.subject_worker_profile_id = pc.worker_profile_id AND NOT freq.inactive
LEFT JOIN hrms.feedback_responses resp ON resp.tenant_id = freq.tenant_id AND resp.request_id = freq.id AND NOT resp.inactive
WHERE pc.tenant_id = $1
  AND pc.inactive = FALSE
  AND ($2::uuid IS NULL OR pc.cycle_id = $2)
  AND ($3::uuid IS NULL OR pc.worker_profile_id = $3)
GROUP BY pc.worker_profile_id, worker.display_name, worker.worker_code, pc.cycle_id, cycle.name
ORDER BY worker.display_name, cycle.name;
