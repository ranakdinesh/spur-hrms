WITH ranked_tokens AS (
    SELECT id,
           ROW_NUMBER() OVER (
               PARTITION BY tenant_id, user_id, device_id
               ORDER BY updated_at DESC, created_at DESC, id DESC
           ) AS row_number
    FROM hrms.device_tokens
    WHERE device_id IS NOT NULL AND NOT inactive
)
UPDATE hrms.device_tokens
SET inactive = TRUE,
    updated_at = NOW()
WHERE id IN (
    SELECT id
    FROM ranked_tokens
    WHERE row_number > 1
);

CREATE UNIQUE INDEX IF NOT EXISTS device_tokens_user_device_unique_idx
    ON hrms.device_tokens(tenant_id, user_id, device_id)
    WHERE device_id IS NOT NULL AND NOT inactive;

CREATE INDEX IF NOT EXISTS device_tokens_token_idx
    ON hrms.device_tokens(tenant_id, device_token)
    WHERE NOT inactive;
