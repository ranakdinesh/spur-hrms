-- name: ListTenantProfiles :many
SELECT * FROM hrms.tenant_profiles
ORDER BY created_at ASC;

-- name: AcquireJobLock :one
INSERT INTO hrms.job_locks (tenant_id, job_key, locked_until, owner_id)
VALUES ($1, $2, NOW() + ($3::text || ' seconds')::interval, $4)
ON CONFLICT (tenant_id, job_key)
DO UPDATE SET
    locked_until = EXCLUDED.locked_until,
    owner_id = EXCLUDED.owner_id,
    last_acquired_at = NOW(),
    updated_at = NOW()
WHERE hrms.job_locks.locked_until <= NOW() OR hrms.job_locks.owner_id = EXCLUDED.owner_id
RETURNING *;

-- name: ReleaseJobLock :exec
UPDATE hrms.job_locks
SET locked_until = NOW(),
    updated_at = NOW()
WHERE tenant_id = $1 AND job_key = $2 AND owner_id = $3;

-- name: GetJobRunByDate :one
SELECT * FROM hrms.job_runs
WHERE tenant_id = $1 AND job_key = $2 AND run_date = $3 AND NOT inactive;

-- name: UpsertJobRunStarted :one
INSERT INTO hrms.job_runs (tenant_id, job_key, run_date, status, owner_id, started_at, metadata, created_by)
VALUES ($1, $2, $3, 'running', $4, NOW(), $5, $6)
ON CONFLICT (tenant_id, job_key, run_date) WHERE NOT inactive
DO UPDATE SET
    status = 'running',
    owner_id = EXCLUDED.owner_id,
    started_at = NOW(),
    finished_at = NULL,
    error_message = NULL,
    metadata = EXCLUDED.metadata,
    updated_at = NOW(),
    updated_by = $6
RETURNING *;

-- name: FinishJobRun :one
UPDATE hrms.job_runs
SET status = $4,
    finished_at = NOW(),
    processed_count = $5,
    success_count = $6,
    failed_count = $7,
    skipped_count = $8,
    error_message = $9,
    metadata = $10,
    updated_at = NOW(),
    updated_by = $11
WHERE tenant_id = $1 AND job_key = $2 AND run_date = $3 AND NOT inactive
RETURNING *;

-- name: ListJobRuns :many
SELECT * FROM hrms.job_runs
WHERE tenant_id = $1 AND job_key = $2 AND NOT inactive
ORDER BY run_date DESC, started_at DESC
LIMIT $3 OFFSET $4;
