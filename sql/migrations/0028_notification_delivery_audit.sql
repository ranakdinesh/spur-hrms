ALTER TABLE hrms.notification_logs
    ADD COLUMN IF NOT EXISTS idempotency_key VARCHAR(120),
    ADD COLUMN IF NOT EXISTS bulk_id UUID;

CREATE UNIQUE INDEX IF NOT EXISTS notification_logs_idempotency_idx
    ON hrms.notification_logs(tenant_id, idempotency_key)
    WHERE idempotency_key IS NOT NULL AND NOT inactive;

CREATE INDEX IF NOT EXISTS notification_logs_bulk_idx
    ON hrms.notification_logs(tenant_id, bulk_id)
    WHERE bulk_id IS NOT NULL AND NOT inactive;

DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'notification_logs_status_check'
    ) THEN
        ALTER TABLE hrms.notification_logs DROP CONSTRAINT notification_logs_status_check;
    END IF;

    ALTER TABLE hrms.notification_logs
        ADD CONSTRAINT notification_logs_status_check
        CHECK (status IN ('Pending','Sent','Failed','Suppressed'));
END $$;
