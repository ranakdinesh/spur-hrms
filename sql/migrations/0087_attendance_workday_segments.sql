-- Field/gig attendance timeline events. These are event-level records for
-- office/site/client/project/travel workday segments, not continuous tracking.
CREATE TABLE IF NOT EXISTS hrms.attendance_workday_segments (
    id                              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id                       UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    user_id                         UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    date                            DATE NOT NULL,
    event_time                      TIMESTAMPTZ NOT NULL,
    segment_type                    VARCHAR(40) NOT NULL,
    action                          VARCHAR(40) NOT NULL,
    work_mode                       VARCHAR(50),
    source                          VARCHAR(50),
    attendance_location_id          UUID REFERENCES hrms.attendance_locations(id) ON DELETE SET NULL,
    reference_type                  VARCHAR(40),
    reference_id                    UUID,
    reference_label                 TEXT,
    latitude                        NUMERIC(10,7),
    longitude                       NUMERIC(10,7),
    location_accuracy_meters        NUMERIC(10,2),
    location_verification_status    VARCHAR(40) NOT NULL DEFAULT 'not_checked',
    remarks                         TEXT,
    metadata                        JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive                        BOOLEAN NOT NULL DEFAULT FALSE,
    created_at                      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by                      UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at                      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by                      UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT attendance_workday_segments_type_check CHECK (segment_type IN ('office','site','client_site','project_site','travel','break','remote','other')),
    CONSTRAINT attendance_workday_segments_action_check CHECK (action IN ('start','end','checkin','checkout','arrive','depart','note')),
    CONSTRAINT attendance_workday_segments_work_mode_check CHECK (work_mode IS NULL OR work_mode IN ('office','remote','field','hybrid','client_site','project_site')),
    CONSTRAINT attendance_workday_segments_source_check CHECK (source IS NULL OR source IN ('web','mobile','kiosk','biometric','api')),
    CONSTRAINT attendance_workday_segments_reference_type_check CHECK (reference_type IS NULL OR reference_type IN ('client','project','site','route','ticket','task','other')),
    CONSTRAINT attendance_workday_segments_location_status_check CHECK (location_verification_status IN ('not_checked','verified','outside_geofence','location_unavailable','device_trusted','pending_review')),
    CONSTRAINT attendance_workday_segments_latitude_check CHECK (latitude IS NULL OR (latitude >= -90 AND latitude <= 90)),
    CONSTRAINT attendance_workday_segments_longitude_check CHECK (longitude IS NULL OR (longitude >= -180 AND longitude <= 180)),
    CONSTRAINT attendance_workday_segments_metadata_object_check CHECK (jsonb_typeof(metadata) = 'object')
);

CREATE INDEX IF NOT EXISTS attendance_workday_segments_user_date_idx
    ON hrms.attendance_workday_segments(tenant_id, user_id, date, event_time)
    WHERE NOT inactive;

CREATE INDEX IF NOT EXISTS attendance_workday_segments_location_idx
    ON hrms.attendance_workday_segments(tenant_id, attendance_location_id, date)
    WHERE attendance_location_id IS NOT NULL AND NOT inactive;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'attendance_workday_segments_updated_at') THEN
        CREATE TRIGGER attendance_workday_segments_updated_at BEFORE UPDATE ON hrms.attendance_workday_segments FOR EACH ROW EXECUTE FUNCTION hrms.update_updated_at();
    END IF;
END $$;

ALTER TABLE hrms.attendance_workday_segments ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS attendance_workday_segments_tenant_isolation ON hrms.attendance_workday_segments;
CREATE POLICY attendance_workday_segments_tenant_isolation ON hrms.attendance_workday_segments
USING (tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::uuid OR current_setting('app.is_super_admin', true) = 'true')
WITH CHECK (tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::uuid OR current_setting('app.is_super_admin', true) = 'true');
