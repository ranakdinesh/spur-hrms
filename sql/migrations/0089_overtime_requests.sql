-- Overtime request approval foundation. Approved records are marked ready for
-- later payroll export/calculation, but this migration does not create payslip
-- earnings or payroll lock side effects.
CREATE TABLE IF NOT EXISTS hrms.overtime_requests (
    id                          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id                   UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    user_id                     UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    work_date                   DATE NOT NULL,
    requested_minutes           INTEGER NOT NULL,
    approved_minutes            INTEGER,
    reason                      TEXT,
    status                      VARCHAR(30) NOT NULL DEFAULT 'pending',
    reviewed_by                 UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    reviewed_at                 TIMESTAMPTZ,
    review_remarks             TEXT,
    calculation_type            VARCHAR(40) NOT NULL DEFAULT 'multiplier',
    rate_multiplier             NUMERIC(6,2) NOT NULL DEFAULT 1.00,
    payroll_component_code      VARCHAR(60),
    payroll_export_status       VARCHAR(30) NOT NULL DEFAULT 'not_ready',
    payroll_exported_at         TIMESTAMPTZ,
    payroll_exported_by         UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    source_attendance_id        UUID REFERENCES hrms.attendances(id) ON DELETE SET NULL,
    source_segment_id           UUID,
    metadata                    JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive                    BOOLEAN NOT NULL DEFAULT FALSE,
    created_at                  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by                  UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at                  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by                  UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT overtime_requests_requested_minutes_check CHECK (requested_minutes > 0),
    CONSTRAINT overtime_requests_approved_minutes_check CHECK (approved_minutes IS NULL OR approved_minutes >= 0),
    CONSTRAINT overtime_requests_status_check CHECK (status IN ('pending','approved','rejected','canceled')),
    CONSTRAINT overtime_requests_calculation_type_check CHECK (calculation_type IN ('fixed_rate','multiplier','payroll_component','manual')),
    CONSTRAINT overtime_requests_rate_multiplier_check CHECK (rate_multiplier >= 0),
    CONSTRAINT overtime_requests_export_status_check CHECK (payroll_export_status IN ('not_ready','ready','exported','not_applicable')),
    CONSTRAINT overtime_requests_metadata_object_check CHECK (jsonb_typeof(metadata) = 'object')
);

CREATE INDEX IF NOT EXISTS overtime_requests_user_status_idx
    ON hrms.overtime_requests(tenant_id, user_id, status, work_date DESC)
    WHERE NOT inactive;

CREATE INDEX IF NOT EXISTS overtime_requests_status_idx
    ON hrms.overtime_requests(tenant_id, status, created_at ASC)
    WHERE NOT inactive;

CREATE INDEX IF NOT EXISTS overtime_requests_payroll_export_idx
    ON hrms.overtime_requests(tenant_id, payroll_export_status, work_date)
    WHERE NOT inactive;

CREATE INDEX IF NOT EXISTS overtime_requests_source_attendance_idx
    ON hrms.overtime_requests(tenant_id, source_attendance_id)
    WHERE source_attendance_id IS NOT NULL AND NOT inactive;

CREATE INDEX IF NOT EXISTS overtime_requests_source_segment_idx
    ON hrms.overtime_requests(tenant_id, source_segment_id)
    WHERE source_segment_id IS NOT NULL AND NOT inactive;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'overtime_requests_updated_at') THEN
        CREATE TRIGGER overtime_requests_updated_at BEFORE UPDATE ON hrms.overtime_requests FOR EACH ROW EXECUTE FUNCTION hrms.update_updated_at();
    END IF;
    IF to_regclass('hrms.attendance_workday_segments') IS NOT NULL
        AND NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'overtime_requests_source_segment_id_fkey') THEN
        ALTER TABLE hrms.overtime_requests
            ADD CONSTRAINT overtime_requests_source_segment_id_fkey
            FOREIGN KEY (source_segment_id) REFERENCES hrms.attendance_workday_segments(id) ON DELETE SET NULL;
    END IF;
END $$;

ALTER TABLE hrms.overtime_requests ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS overtime_requests_tenant_isolation ON hrms.overtime_requests;
CREATE POLICY overtime_requests_tenant_isolation ON hrms.overtime_requests
USING (tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::uuid OR current_setting('app.is_super_admin', true) = 'true')
WITH CHECK (tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::uuid OR current_setting('app.is_super_admin', true) = 'true');
