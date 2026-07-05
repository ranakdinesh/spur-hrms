-- Tenant attendance locations, biometric devices, and device-employee mappings.

CREATE TABLE IF NOT EXISTS hrms.attendance_locations (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id           UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    branch_id           UUID REFERENCES hrms.branches(id) ON DELETE SET NULL,
    code                VARCHAR(80) NOT NULL,
    name                VARCHAR(255) NOT NULL,
    location_type       VARCHAR(40) NOT NULL DEFAULT 'office',
    latitude            NUMERIC(10,7),
    longitude           NUMERIC(10,7),
    radius_meters       INT NOT NULL DEFAULT 100,
    address             TEXT,
    city                VARCHAR(120),
    state               VARCHAR(120),
    country             VARCHAR(120),
    pincode             VARCHAR(30),
    effective_from      DATE,
    effective_to        DATE,
    inactive            BOOLEAN NOT NULL DEFAULT FALSE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by          UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by          UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT attendance_locations_type_check CHECK (location_type IN ('office','branch','warehouse','client_site','field_zone','project_site','remote','other')),
    CONSTRAINT attendance_locations_radius_check CHECK (radius_meters >= 0),
    CONSTRAINT attendance_locations_latitude_check CHECK (latitude IS NULL OR latitude BETWEEN -90 AND 90),
    CONSTRAINT attendance_locations_longitude_check CHECK (longitude IS NULL OR longitude BETWEEN -180 AND 180),
    CONSTRAINT attendance_locations_effective_order_check CHECK (effective_to IS NULL OR effective_from IS NULL OR effective_from <= effective_to)
);
CREATE UNIQUE INDEX IF NOT EXISTS attendance_locations_code_idx ON hrms.attendance_locations(tenant_id, lower(code)) WHERE NOT inactive;
CREATE INDEX IF NOT EXISTS attendance_locations_tenant_idx ON hrms.attendance_locations(tenant_id) WHERE NOT inactive;
CREATE INDEX IF NOT EXISTS attendance_locations_branch_idx ON hrms.attendance_locations(tenant_id, branch_id) WHERE branch_id IS NOT NULL AND NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.attendance_location_assignments (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id           UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    location_id         UUID NOT NULL REFERENCES hrms.attendance_locations(id) ON DELETE CASCADE,
    branch_id           UUID REFERENCES hrms.branches(id) ON DELETE CASCADE,
    department_id       UUID REFERENCES hrms.departments(id) ON DELETE CASCADE,
    user_id             UUID REFERENCES auth.users(id) ON DELETE CASCADE,
    effective_from      DATE,
    effective_to        DATE,
    inactive            BOOLEAN NOT NULL DEFAULT FALSE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by          UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by          UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT attendance_location_assignments_scope_check CHECK (
        (CASE WHEN branch_id IS NULL THEN 0 ELSE 1 END)
        + (CASE WHEN department_id IS NULL THEN 0 ELSE 1 END)
        + (CASE WHEN user_id IS NULL THEN 0 ELSE 1 END) <= 1
    ),
    CONSTRAINT attendance_location_assignments_effective_order_check CHECK (effective_to IS NULL OR effective_from IS NULL OR effective_from <= effective_to)
);
CREATE INDEX IF NOT EXISTS attendance_location_assignments_location_idx ON hrms.attendance_location_assignments(tenant_id, location_id) WHERE NOT inactive;
CREATE INDEX IF NOT EXISTS attendance_location_assignments_user_idx ON hrms.attendance_location_assignments(tenant_id, user_id) WHERE user_id IS NOT NULL AND NOT inactive;
CREATE INDEX IF NOT EXISTS attendance_location_assignments_scope_idx ON hrms.attendance_location_assignments(tenant_id, branch_id, department_id) WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.attendance_devices (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id               UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    attendance_location_id  UUID REFERENCES hrms.attendance_locations(id) ON DELETE SET NULL,
    branch_id               UUID REFERENCES hrms.branches(id) ON DELETE SET NULL,
    code                    VARCHAR(80) NOT NULL,
    name                    VARCHAR(255) NOT NULL,
    vendor                  VARCHAR(120),
    model                   VARCHAR(120),
    serial_number           VARCHAR(160),
    device_identifier       VARCHAR(160),
    integration_type        VARCHAR(40) NOT NULL DEFAULT 'edge_agent',
    direction_mode          VARCHAR(40) NOT NULL DEFAULT 'auto',
    timezone                VARCHAR(80) NOT NULL DEFAULT 'UTC',
    status                  VARCHAR(30) NOT NULL DEFAULT 'active',
    last_seen_at            TIMESTAMPTZ,
    config                  JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive                BOOLEAN NOT NULL DEFAULT FALSE,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by              UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by              UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT attendance_devices_integration_type_check CHECK (integration_type IN ('push','poll','file_import','api','edge_agent')),
    CONSTRAINT attendance_devices_direction_mode_check CHECK (direction_mode IN ('auto','in_out','entry_exit','checkin_only','checkout_only')),
    CONSTRAINT attendance_devices_status_check CHECK (status IN ('active','inactive','maintenance')),
    CONSTRAINT attendance_devices_config_object_check CHECK (jsonb_typeof(config) = 'object')
);
CREATE UNIQUE INDEX IF NOT EXISTS attendance_devices_code_idx ON hrms.attendance_devices(tenant_id, lower(code)) WHERE NOT inactive;
CREATE UNIQUE INDEX IF NOT EXISTS attendance_devices_serial_idx ON hrms.attendance_devices(tenant_id, lower(serial_number)) WHERE serial_number IS NOT NULL AND NOT inactive;
CREATE INDEX IF NOT EXISTS attendance_devices_location_idx ON hrms.attendance_devices(tenant_id, attendance_location_id) WHERE attendance_location_id IS NOT NULL AND NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.employee_attendance_devices (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id           UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    user_id             UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    device_id           UUID NOT NULL REFERENCES hrms.attendance_devices(id) ON DELETE CASCADE,
    device_user_id      VARCHAR(160) NOT NULL,
    credential_type     VARCHAR(40) NOT NULL DEFAULT 'biometric',
    card_number         VARCHAR(160),
    effective_from      DATE,
    effective_to        DATE,
    inactive            BOOLEAN NOT NULL DEFAULT FALSE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by          UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by          UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT employee_attendance_devices_credential_type_check CHECK (credential_type IN ('biometric','fingerprint','face','card','pin','mobile','other')),
    CONSTRAINT employee_attendance_devices_effective_order_check CHECK (effective_to IS NULL OR effective_from IS NULL OR effective_from <= effective_to)
);
CREATE UNIQUE INDEX IF NOT EXISTS employee_attendance_devices_user_device_idx ON hrms.employee_attendance_devices(tenant_id, user_id, device_id) WHERE NOT inactive;
CREATE UNIQUE INDEX IF NOT EXISTS employee_attendance_devices_device_user_idx ON hrms.employee_attendance_devices(tenant_id, device_id, lower(device_user_id)) WHERE NOT inactive;
CREATE INDEX IF NOT EXISTS employee_attendance_devices_user_idx ON hrms.employee_attendance_devices(tenant_id, user_id) WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.raw_attendance_events (
    id                              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id                       UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    device_id                       UUID NOT NULL REFERENCES hrms.attendance_devices(id) ON DELETE CASCADE,
    employee_device_mapping_id      UUID REFERENCES hrms.employee_attendance_devices(id) ON DELETE SET NULL,
    attendance_id                   UUID REFERENCES hrms.attendances(id) ON DELETE SET NULL,
    external_event_id               VARCHAR(200),
    device_user_id                  VARCHAR(160),
    event_time                      TIMESTAMPTZ NOT NULL,
    event_type                      VARCHAR(60),
    attendance_type                 VARCHAR(20),
    import_batch_id                 VARCHAR(120),
    processing_status               VARCHAR(30) NOT NULL DEFAULT 'pending',
    processing_error                TEXT,
    raw_payload                     JSONB NOT NULL DEFAULT '{}'::jsonb,
    processed_at                    TIMESTAMPTZ,
    inactive                        BOOLEAN NOT NULL DEFAULT FALSE,
    created_at                      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by                      UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at                      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by                      UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT raw_attendance_events_attendance_type_check CHECK (attendance_type IS NULL OR attendance_type IN ('checkin','checkout')),
    CONSTRAINT raw_attendance_events_processing_status_check CHECK (processing_status IN ('pending','processed','ignored','failed')),
    CONSTRAINT raw_attendance_events_payload_object_check CHECK (jsonb_typeof(raw_payload) = 'object')
);
CREATE UNIQUE INDEX IF NOT EXISTS raw_attendance_events_external_idx ON hrms.raw_attendance_events(tenant_id, device_id, external_event_id) WHERE external_event_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS raw_attendance_events_pending_idx ON hrms.raw_attendance_events(tenant_id, processing_status, event_time) WHERE NOT inactive;
CREATE INDEX IF NOT EXISTS raw_attendance_events_device_user_idx ON hrms.raw_attendance_events(tenant_id, device_id, device_user_id, event_time) WHERE device_user_id IS NOT NULL;

ALTER TABLE hrms.attendances
    ADD COLUMN IF NOT EXISTS attendance_location_id UUID REFERENCES hrms.attendance_locations(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS attendance_device_id UUID REFERENCES hrms.attendance_devices(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS raw_attendance_event_id UUID REFERENCES hrms.raw_attendance_events(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS location_accuracy_meters NUMERIC(10,2),
    ADD COLUMN IF NOT EXISTS location_verification_status VARCHAR(40) NOT NULL DEFAULT 'not_checked';

ALTER TABLE hrms.attendances DROP CONSTRAINT IF EXISTS attendances_location_verification_status_check;
ALTER TABLE hrms.attendances ADD CONSTRAINT attendances_location_verification_status_check CHECK (location_verification_status IN ('not_checked','verified','outside_geofence','location_unavailable','device_trusted','pending_review'));
CREATE INDEX IF NOT EXISTS attendances_location_idx ON hrms.attendances(tenant_id, attendance_location_id, date) WHERE attendance_location_id IS NOT NULL AND NOT inactive;
CREATE INDEX IF NOT EXISTS attendances_device_idx ON hrms.attendances(tenant_id, attendance_device_id, date) WHERE attendance_device_id IS NOT NULL AND NOT inactive;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'attendance_locations_updated_at') THEN
        CREATE TRIGGER attendance_locations_updated_at BEFORE UPDATE ON hrms.attendance_locations FOR EACH ROW EXECUTE FUNCTION hrms.update_updated_at();
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'attendance_location_assignments_updated_at') THEN
        CREATE TRIGGER attendance_location_assignments_updated_at BEFORE UPDATE ON hrms.attendance_location_assignments FOR EACH ROW EXECUTE FUNCTION hrms.update_updated_at();
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'attendance_devices_updated_at') THEN
        CREATE TRIGGER attendance_devices_updated_at BEFORE UPDATE ON hrms.attendance_devices FOR EACH ROW EXECUTE FUNCTION hrms.update_updated_at();
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'employee_attendance_devices_updated_at') THEN
        CREATE TRIGGER employee_attendance_devices_updated_at BEFORE UPDATE ON hrms.employee_attendance_devices FOR EACH ROW EXECUTE FUNCTION hrms.update_updated_at();
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'raw_attendance_events_updated_at') THEN
        CREATE TRIGGER raw_attendance_events_updated_at BEFORE UPDATE ON hrms.raw_attendance_events FOR EACH ROW EXECUTE FUNCTION hrms.update_updated_at();
    END IF;
END $$;

ALTER TABLE hrms.attendance_locations ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.attendance_location_assignments ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.attendance_devices ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.employee_attendance_devices ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.raw_attendance_events ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS attendance_locations_tenant_isolation ON hrms.attendance_locations;
CREATE POLICY attendance_locations_tenant_isolation ON hrms.attendance_locations
USING (tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::uuid OR current_setting('app.is_super_admin', true) = 'true')
WITH CHECK (tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::uuid OR current_setting('app.is_super_admin', true) = 'true');

DROP POLICY IF EXISTS attendance_location_assignments_tenant_isolation ON hrms.attendance_location_assignments;
CREATE POLICY attendance_location_assignments_tenant_isolation ON hrms.attendance_location_assignments
USING (tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::uuid OR current_setting('app.is_super_admin', true) = 'true')
WITH CHECK (tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::uuid OR current_setting('app.is_super_admin', true) = 'true');

DROP POLICY IF EXISTS attendance_devices_tenant_isolation ON hrms.attendance_devices;
CREATE POLICY attendance_devices_tenant_isolation ON hrms.attendance_devices
USING (tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::uuid OR current_setting('app.is_super_admin', true) = 'true')
WITH CHECK (tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::uuid OR current_setting('app.is_super_admin', true) = 'true');

DROP POLICY IF EXISTS employee_attendance_devices_tenant_isolation ON hrms.employee_attendance_devices;
CREATE POLICY employee_attendance_devices_tenant_isolation ON hrms.employee_attendance_devices
USING (tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::uuid OR current_setting('app.is_super_admin', true) = 'true')
WITH CHECK (tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::uuid OR current_setting('app.is_super_admin', true) = 'true');

DROP POLICY IF EXISTS raw_attendance_events_tenant_isolation ON hrms.raw_attendance_events;
CREATE POLICY raw_attendance_events_tenant_isolation ON hrms.raw_attendance_events
USING (tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::uuid OR current_setting('app.is_super_admin', true) = 'true')
WITH CHECK (tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::uuid OR current_setting('app.is_super_admin', true) = 'true');
