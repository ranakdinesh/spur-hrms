ALTER TABLE hrms.interview_rounds
    ADD COLUMN IF NOT EXISTS timezone VARCHAR(64) NOT NULL DEFAULT 'UTC',
    ADD COLUMN IF NOT EXISTS feedback TEXT,
    ADD COLUMN IF NOT EXISTS score NUMERIC(4,2),
    ADD COLUMN IF NOT EXISTS decision VARCHAR(50),
    ADD COLUMN IF NOT EXISTS completed_at TIMESTAMPTZ;

UPDATE hrms.interview_rounds
SET status = 'Scheduled'
WHERE status IS NULL OR BTRIM(status) = '';

ALTER TABLE hrms.interview_rounds
    ALTER COLUMN status SET DEFAULT 'Scheduled',
    ALTER COLUMN status SET NOT NULL;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'interview_rounds_status_check'
          AND conrelid = 'hrms.interview_rounds'::regclass
    ) THEN
        ALTER TABLE hrms.interview_rounds
            ADD CONSTRAINT interview_rounds_status_check
            CHECK (status IN ('Scheduled','Rescheduled','Completed','Cancelled','NoShow'));
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'interview_rounds_mode_check'
          AND conrelid = 'hrms.interview_rounds'::regclass
    ) THEN
        ALTER TABLE hrms.interview_rounds
            ADD CONSTRAINT interview_rounds_mode_check
            CHECK (mode IS NULL OR mode IN ('Phone','Video','InPerson','Panel','Assignment'));
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'interview_rounds_decision_check'
          AND conrelid = 'hrms.interview_rounds'::regclass
    ) THEN
        ALTER TABLE hrms.interview_rounds
            ADD CONSTRAINT interview_rounds_decision_check
            CHECK (decision IS NULL OR decision IN ('StrongHire','Hire','Hold','NoHire'));
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'interview_rounds_score_check'
          AND conrelid = 'hrms.interview_rounds'::regclass
    ) THEN
        ALTER TABLE hrms.interview_rounds
            ADD CONSTRAINT interview_rounds_score_check
            CHECK (score IS NULL OR (score >= 0 AND score <= 5));
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS interview_rounds_schedule_idx
    ON hrms.interview_rounds(tenant_id, scheduled_date)
    WHERE NOT inactive;

CREATE INDEX IF NOT EXISTS interview_rounds_interviewer_idx
    ON hrms.interview_rounds(tenant_id, interviewer_user_id)
    WHERE NOT inactive AND interviewer_user_id IS NOT NULL;
