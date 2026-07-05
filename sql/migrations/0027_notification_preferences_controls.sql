ALTER TABLE hrms.notification_preferences
    ADD COLUMN IF NOT EXISTS is_push_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    ADD COLUMN IF NOT EXISTS digest_frequency VARCHAR(20) NOT NULL DEFAULT 'instant',
    ADD COLUMN IF NOT EXISTS quiet_hours_start TIME,
    ADD COLUMN IF NOT EXISTS quiet_hours_end TIME,
    ADD COLUMN IF NOT EXISTS timezone VARCHAR(80);

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'notification_preferences_digest_frequency_check'
    ) THEN
        ALTER TABLE hrms.notification_preferences
            ADD CONSTRAINT notification_preferences_digest_frequency_check
            CHECK (digest_frequency IN ('instant','daily','weekly','muted'));
    END IF;
END $$;
