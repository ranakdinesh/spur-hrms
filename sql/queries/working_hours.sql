-- name: CreateWorkingHour :one
INSERT INTO hrms.working_hours (
    tenant_id,
    branch_id,
    user_id,
    day_of_week,
    is_working_day,
    start_time,
    end_time,
    break_minutes,
    created_by,
    updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $9
)
RETURNING *;

-- name: ListWorkingHours :many
SELECT * FROM hrms.working_hours
WHERE tenant_id = $1 AND NOT inactive
ORDER BY branch_id NULLS FIRST,
         user_id NULLS FIRST,
         CASE day_of_week
             WHEN 'Monday' THEN 1 WHEN 'Tuesday' THEN 2 WHEN 'Wednesday' THEN 3
             WHEN 'Thursday' THEN 4 WHEN 'Friday' THEN 5 WHEN 'Saturday' THEN 6
             WHEN 'Sunday' THEN 7 ELSE 8 END;

-- name: GetWorkingHour :one
SELECT * FROM hrms.working_hours
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: ResolveWorkingHour :one
SELECT * FROM hrms.working_hours
WHERE tenant_id = $1
  AND day_of_week = $2
  AND NOT inactive
  AND (
      (user_id = $3)
      OR (branch_id = $4 AND user_id IS NULL)
      OR (branch_id IS NULL AND user_id IS NULL)
  )
ORDER BY CASE
    WHEN user_id = $3 THEN 1
    WHEN branch_id = $4 AND user_id IS NULL THEN 2
    WHEN branch_id IS NULL AND user_id IS NULL THEN 3
    ELSE 4 END
LIMIT 1;

-- name: UpdateWorkingHour :one
UPDATE hrms.working_hours
SET branch_id = $3,
    user_id = $4,
    day_of_week = $5,
    is_working_day = $6,
    start_time = $7,
    end_time = $8,
    break_minutes = $9,
    updated_by = $10,
    updated_at = NOW()
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: SoftDeleteWorkingHour :exec
UPDATE hrms.working_hours
SET inactive = TRUE,
    updated_by = $3,
    updated_at = NOW()
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: SoftDeleteBranchWorkingHours :exec
UPDATE hrms.working_hours
SET inactive = TRUE,
    updated_by = $3,
    updated_at = NOW()
WHERE tenant_id = $1 AND branch_id = $2 AND user_id IS NULL AND NOT inactive;

-- name: CopyTenantWorkingHoursToBranch :many
INSERT INTO hrms.working_hours (
    tenant_id,
    branch_id,
    user_id,
    day_of_week,
    is_working_day,
    start_time,
    end_time,
    break_minutes,
    created_by,
    updated_by
)
SELECT source.tenant_id,
       $2,
       NULL,
       source.day_of_week,
       source.is_working_day,
       source.start_time,
       source.end_time,
       source.break_minutes,
       $3,
       $3
FROM hrms.working_hours source
WHERE source.tenant_id = $1 AND source.branch_id IS NULL AND source.user_id IS NULL AND NOT source.inactive
ORDER BY CASE source.day_of_week
    WHEN 'Monday' THEN 1 WHEN 'Tuesday' THEN 2 WHEN 'Wednesday' THEN 3
    WHEN 'Thursday' THEN 4 WHEN 'Friday' THEN 5 WHEN 'Saturday' THEN 6
    WHEN 'Sunday' THEN 7 ELSE 8 END
RETURNING *;
