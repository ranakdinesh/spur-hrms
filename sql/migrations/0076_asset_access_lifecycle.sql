CREATE TABLE IF NOT EXISTS hrms.asset_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    asset_code VARCHAR(80) NOT NULL,
    asset_name VARCHAR(180) NOT NULL,
    asset_type VARCHAR(60) NOT NULL DEFAULT 'hardware',
    category VARCHAR(80) NOT NULL DEFAULT 'general',
    serial_number VARCHAR(120),
    vendor VARCHAR(140),
    purchase_date DATE,
    warranty_until DATE,
    owner_user_id UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    custodian_worker_profile_id UUID REFERENCES hrms.worker_profiles(id) ON DELETE SET NULL,
    status VARCHAR(40) NOT NULL DEFAULT 'available',
    current_assignment_id UUID,
    location_label VARCHAR(160),
    notes TEXT,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT asset_items_code_check CHECK (BTRIM(asset_code) <> ''),
    CONSTRAINT asset_items_name_check CHECK (BTRIM(asset_name) <> ''),
    CONSTRAINT asset_items_status_check CHECK (status IN ('available','reserved','issued','return_due','returned','maintenance','damaged','lost','retired'))
);

CREATE UNIQUE INDEX IF NOT EXISTS asset_items_code_unique_idx
    ON hrms.asset_items(tenant_id, lower(asset_code))
    WHERE NOT inactive;

CREATE INDEX IF NOT EXISTS asset_items_tenant_status_idx
    ON hrms.asset_items(tenant_id, status, category)
    WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.access_catalog_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    access_code VARCHAR(80) NOT NULL,
    access_name VARCHAR(180) NOT NULL,
    access_type VARCHAR(60) NOT NULL DEFAULT 'software',
    system_name VARCHAR(140),
    owner_user_id UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    provisioning_method VARCHAR(60) NOT NULL DEFAULT 'manual',
    requires_approval BOOLEAN NOT NULL DEFAULT TRUE,
    default_for_onboarding BOOLEAN NOT NULL DEFAULT FALSE,
    default_for_exit_revocation BOOLEAN NOT NULL DEFAULT TRUE,
    status VARCHAR(40) NOT NULL DEFAULT 'active',
    notes TEXT,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT access_catalog_code_check CHECK (BTRIM(access_code) <> ''),
    CONSTRAINT access_catalog_name_check CHECK (BTRIM(access_name) <> ''),
    CONSTRAINT access_catalog_status_check CHECK (status IN ('active','inactive','deprecated'))
);

CREATE UNIQUE INDEX IF NOT EXISTS access_catalog_code_unique_idx
    ON hrms.access_catalog_items(tenant_id, lower(access_code))
    WHERE NOT inactive;

CREATE INDEX IF NOT EXISTS access_catalog_tenant_status_idx
    ON hrms.access_catalog_items(tenant_id, status, access_type)
    WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.asset_assignments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    asset_id UUID NOT NULL REFERENCES hrms.asset_items(id) ON DELETE CASCADE,
    worker_profile_id UUID NOT NULL REFERENCES hrms.worker_profiles(id) ON DELETE CASCADE,
    employee_id UUID REFERENCES hrms.employees(id) ON DELETE SET NULL,
    candidate_onboarding_id UUID REFERENCES hrms.candidate_onboardings(id) ON DELETE SET NULL,
    exit_request_id UUID REFERENCES hrms.employee_exit_requests(id) ON DELETE SET NULL,
    requested_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    approved_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    issued_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    returned_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    approved_at TIMESTAMPTZ,
    issued_on DATE,
    expected_return_on DATE,
    returned_on DATE,
    issue_condition VARCHAR(80) NOT NULL DEFAULT 'good',
    return_condition VARCHAR(80),
    damage_status VARCHAR(40) NOT NULL DEFAULT 'none',
    recovery_amount NUMERIC(14,2) NOT NULL DEFAULT 0,
    status VARCHAR(40) NOT NULL DEFAULT 'requested',
    notes TEXT,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT asset_assignments_status_check CHECK (status IN ('requested','approved','issued','return_due','returned','damaged','lost','cancelled')),
    CONSTRAINT asset_assignments_damage_check CHECK (damage_status IN ('none','minor','major','lost','recovered')),
    CONSTRAINT asset_assignments_recovery_check CHECK (recovery_amount >= 0)
);

CREATE UNIQUE INDEX IF NOT EXISTS asset_assignments_one_active_idx
    ON hrms.asset_assignments(tenant_id, asset_id)
    WHERE NOT inactive AND status IN ('requested','approved','issued','return_due','damaged','lost');

CREATE INDEX IF NOT EXISTS asset_assignments_worker_idx
    ON hrms.asset_assignments(tenant_id, worker_profile_id, status)
    WHERE NOT inactive;

CREATE INDEX IF NOT EXISTS asset_assignments_exit_idx
    ON hrms.asset_assignments(tenant_id, exit_request_id)
    WHERE NOT inactive AND exit_request_id IS NOT NULL;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'asset_items_current_assignment_fk'
          AND conrelid = 'hrms.asset_items'::regclass
    ) THEN
        ALTER TABLE hrms.asset_items
            ADD CONSTRAINT asset_items_current_assignment_fk
            FOREIGN KEY (current_assignment_id) REFERENCES hrms.asset_assignments(id) ON DELETE SET NULL;
    END IF;
END $$;

CREATE TABLE IF NOT EXISTS hrms.access_lifecycle_tasks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    access_item_id UUID NOT NULL REFERENCES hrms.access_catalog_items(id) ON DELETE CASCADE,
    worker_profile_id UUID NOT NULL REFERENCES hrms.worker_profiles(id) ON DELETE CASCADE,
    employee_id UUID REFERENCES hrms.employees(id) ON DELETE SET NULL,
    candidate_onboarding_id UUID REFERENCES hrms.candidate_onboardings(id) ON DELETE SET NULL,
    exit_request_id UUID REFERENCES hrms.employee_exit_requests(id) ON DELETE SET NULL,
    task_type VARCHAR(40) NOT NULL DEFAULT 'provision',
    requested_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    approved_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    owner_user_id UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    approved_at TIMESTAMPTZ,
    due_date DATE,
    completed_at TIMESTAMPTZ,
    external_reference VARCHAR(160),
    status VARCHAR(40) NOT NULL DEFAULT 'requested',
    notes TEXT,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT access_lifecycle_type_check CHECK (task_type IN ('provision','deprovision','review','change')),
    CONSTRAINT access_lifecycle_status_check CHECK (status IN ('requested','approved','provisioned','revoked','reviewed','rejected','cancelled','blocked'))
);

CREATE INDEX IF NOT EXISTS access_tasks_worker_idx
    ON hrms.access_lifecycle_tasks(tenant_id, worker_profile_id, status)
    WHERE NOT inactive;

CREATE INDEX IF NOT EXISTS access_tasks_exit_idx
    ON hrms.access_lifecycle_tasks(tenant_id, exit_request_id)
    WHERE NOT inactive AND exit_request_id IS NOT NULL;

CREATE TABLE IF NOT EXISTS hrms.asset_access_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    source_type VARCHAR(60) NOT NULL,
    source_id UUID,
    action VARCHAR(80) NOT NULL,
    from_status VARCHAR(50),
    to_status VARCHAR(50),
    remarks TEXT,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT asset_access_events_source_check CHECK (BTRIM(source_type) <> ''),
    CONSTRAINT asset_access_events_action_check CHECK (BTRIM(action) <> '')
);

CREATE INDEX IF NOT EXISTS asset_access_events_source_idx
    ON hrms.asset_access_events(tenant_id, source_type, source_id, created_at DESC)
    WHERE NOT inactive;

DROP TRIGGER IF EXISTS asset_items_set_updated_at ON hrms.asset_items;
CREATE TRIGGER asset_items_set_updated_at BEFORE UPDATE ON hrms.asset_items
    FOR EACH ROW EXECUTE FUNCTION hrms.set_updated_at();
DROP TRIGGER IF EXISTS access_catalog_items_set_updated_at ON hrms.access_catalog_items;
CREATE TRIGGER access_catalog_items_set_updated_at BEFORE UPDATE ON hrms.access_catalog_items
    FOR EACH ROW EXECUTE FUNCTION hrms.set_updated_at();
DROP TRIGGER IF EXISTS asset_assignments_set_updated_at ON hrms.asset_assignments;
CREATE TRIGGER asset_assignments_set_updated_at BEFORE UPDATE ON hrms.asset_assignments
    FOR EACH ROW EXECUTE FUNCTION hrms.set_updated_at();
DROP TRIGGER IF EXISTS access_lifecycle_tasks_set_updated_at ON hrms.access_lifecycle_tasks;
CREATE TRIGGER access_lifecycle_tasks_set_updated_at BEFORE UPDATE ON hrms.access_lifecycle_tasks
    FOR EACH ROW EXECUTE FUNCTION hrms.set_updated_at();
DROP TRIGGER IF EXISTS asset_access_events_set_updated_at ON hrms.asset_access_events;
CREATE TRIGGER asset_access_events_set_updated_at BEFORE UPDATE ON hrms.asset_access_events
    FOR EACH ROW EXECUTE FUNCTION hrms.set_updated_at();

ALTER TABLE hrms.asset_items ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.access_catalog_items ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.asset_assignments ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.access_lifecycle_tasks ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.asset_access_events ENABLE ROW LEVEL SECURITY;

DO $$
DECLARE
    table_name text;
BEGIN
    FOREACH table_name IN ARRAY ARRAY['asset_items','access_catalog_items','asset_assignments','access_lifecycle_tasks','asset_access_events']
    LOOP
        EXECUTE format('DROP POLICY IF EXISTS tenant_isolation_%1$s ON hrms.%1$s', table_name);
        EXECUTE format(
            'CREATE POLICY tenant_isolation_%1$s ON hrms.%1$s USING (current_setting(''app.is_super_admin'', true) = ''true'' OR tenant_id::text = current_setting(''app.tenant_id'', true)) WITH CHECK (current_setting(''app.is_super_admin'', true) = ''true'' OR tenant_id::text = current_setting(''app.tenant_id'', true))',
            table_name
        );
    END LOOP;
END $$;
