WITH ranked AS (
    SELECT id,
           ROW_NUMBER() OVER (
               PARTITION BY tenant_id,
                            COALESCE(branch_id, '00000000-0000-0000-0000-000000000000'::uuid),
                            COALESCE(user_id, '00000000-0000-0000-0000-000000000000'::uuid),
                            day_of_week
               ORDER BY updated_at DESC, created_at DESC, id DESC
           ) AS row_number
    FROM hrms.working_hours
    WHERE NOT inactive
)
UPDATE hrms.working_hours wh
SET inactive = TRUE,
    updated_at = NOW()
FROM ranked
WHERE wh.id = ranked.id AND ranked.row_number > 1;

CREATE UNIQUE INDEX IF NOT EXISTS working_hours_scope_day_active_idx
    ON hrms.working_hours (
        tenant_id,
        COALESCE(branch_id, '00000000-0000-0000-0000-000000000000'::uuid),
        COALESCE(user_id, '00000000-0000-0000-0000-000000000000'::uuid),
        day_of_week
    )
    WHERE NOT inactive;
