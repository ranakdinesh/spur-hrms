-- Extend the older comp-off request table into the approval-ledger model.
-- Migration 0017 already created hrms.comp_off_requests with worked_hours,
-- credit_days, multiplier, and remarks. Keep those columns for compatibility
-- and add the explicit request/review fields used by the new APIs.
ALTER TABLE hrms.comp_off_requests
    ADD COLUMN IF NOT EXISTS worked_minutes INTEGER NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS requested_days NUMERIC(6,2) NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS approved_days NUMERIC(6,2),
    ADD COLUMN IF NOT EXISTS expiry_date DATE,
    ADD COLUMN IF NOT EXISTS review_remarks TEXT,
    ADD COLUMN IF NOT EXISTS payroll_impact BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS source_attendance_id UUID REFERENCES hrms.attendances(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS source_segment_id UUID,
    ADD COLUMN IF NOT EXISTS metadata JSONB NOT NULL DEFAULT '{}'::jsonb;

UPDATE hrms.comp_off_requests
SET
    worked_minutes = CASE
        WHEN worked_minutes = 0 AND worked_hours IS NOT NULL THEN ROUND(worked_hours * 60)::INTEGER
        ELSE worked_minutes
    END,
    requested_days = CASE
        WHEN requested_days = 0 AND credit_days IS NOT NULL THEN credit_days
        ELSE requested_days
    END,
    approved_days = CASE
        WHEN approved_days IS NULL AND status = 'approved' THEN credit_days
        ELSE approved_days
    END,
    review_remarks = COALESCE(review_remarks, remarks),
    metadata = COALESCE(metadata, '{}'::jsonb);

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'comp_off_requests_approved_days_check') THEN
        ALTER TABLE hrms.comp_off_requests
            ADD CONSTRAINT comp_off_requests_approved_days_check CHECK (approved_days IS NULL OR approved_days >= 0);
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'comp_off_requests_worked_minutes_check') THEN
        ALTER TABLE hrms.comp_off_requests
            ADD CONSTRAINT comp_off_requests_worked_minutes_check CHECK (worked_minutes >= 0);
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'comp_off_requests_metadata_object_check') THEN
        ALTER TABLE hrms.comp_off_requests
            ADD CONSTRAINT comp_off_requests_metadata_object_check CHECK (jsonb_typeof(metadata) = 'object');
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'comp_off_requests_updated_at') THEN
        CREATE TRIGGER comp_off_requests_updated_at BEFORE UPDATE ON hrms.comp_off_requests FOR EACH ROW EXECUTE FUNCTION hrms.update_updated_at();
    END IF;
    IF to_regclass('hrms.attendance_workday_segments') IS NOT NULL
        AND NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'comp_off_requests_source_segment_id_fkey') THEN
        ALTER TABLE hrms.comp_off_requests
            ADD CONSTRAINT comp_off_requests_source_segment_id_fkey
            FOREIGN KEY (source_segment_id) REFERENCES hrms.attendance_workday_segments(id) ON DELETE SET NULL;
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS comp_off_requests_user_status_idx
    ON hrms.comp_off_requests(tenant_id, user_id, status, work_date DESC)
    WHERE NOT inactive;

CREATE INDEX IF NOT EXISTS comp_off_requests_status_open_idx
    ON hrms.comp_off_requests(tenant_id, status, created_at DESC)
    WHERE NOT inactive;

CREATE INDEX IF NOT EXISTS comp_off_requests_source_attendance_idx
    ON hrms.comp_off_requests(tenant_id, source_attendance_id)
    WHERE source_attendance_id IS NOT NULL AND NOT inactive;

CREATE INDEX IF NOT EXISTS comp_off_requests_source_segment_idx
    ON hrms.comp_off_requests(tenant_id, source_segment_id)
    WHERE source_segment_id IS NOT NULL AND NOT inactive;
