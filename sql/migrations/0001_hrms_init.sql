-- Hrms module — initial migration
-- HRMS owns product-specific tenant profile/settings while identity owns auth.tenants.

CREATE SCHEMA IF NOT EXISTS hrms;

CREATE TABLE IF NOT EXISTS hrms.tenant_profiles (
    tenant_id               UUID        PRIMARY KEY REFERENCES auth.tenants(id) ON DELETE CASCADE,
    subdomain               TEXT        NOT NULL,
    mobile_activation_code  TEXT        NOT NULL,
    display_name            TEXT,
    logo_object_key         TEXT,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT tenant_profiles_subdomain_format
        CHECK (subdomain ~ '^[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?$'),
    CONSTRAINT tenant_profiles_activation_code_format
        CHECK (mobile_activation_code ~ '^[A-Z0-9-]{4,32}$')
);

CREATE UNIQUE INDEX IF NOT EXISTS tenant_profiles_subdomain_idx
    ON hrms.tenant_profiles (subdomain);

CREATE UNIQUE INDEX IF NOT EXISTS tenant_profiles_activation_code_idx
    ON hrms.tenant_profiles (mobile_activation_code);

CREATE TABLE IF NOT EXISTS hrms.tenant_settings (
    tenant_id   UUID        NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    key         TEXT        NOT NULL,
    value       JSONB       NOT NULL DEFAULT '{}'::jsonb,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (tenant_id, key),
    CONSTRAINT tenant_settings_key_format
        CHECK (key ~ '^[a-z][a-z0-9_.-]*$')
);

CREATE INDEX IF NOT EXISTS tenant_settings_tenant_idx
    ON hrms.tenant_settings (tenant_id);

CREATE OR REPLACE FUNCTION hrms.update_updated_at()
RETURNS TRIGGER LANGUAGE plpgsql AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$;

CREATE OR REPLACE FUNCTION hrms.set_updated_at()
RETURNS TRIGGER LANGUAGE plpgsql AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$;

DROP TRIGGER IF EXISTS tenant_profiles_updated_at ON hrms.tenant_profiles;
CREATE TRIGGER tenant_profiles_updated_at
    BEFORE UPDATE ON hrms.tenant_profiles
    FOR EACH ROW EXECUTE FUNCTION hrms.update_updated_at();

DROP TRIGGER IF EXISTS tenant_settings_updated_at ON hrms.tenant_settings;
CREATE TRIGGER tenant_settings_updated_at
    BEFORE UPDATE ON hrms.tenant_settings
    FOR EACH ROW EXECUTE FUNCTION hrms.update_updated_at();

ALTER TABLE hrms.tenant_profiles ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.tenant_settings ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS tenant_profiles_isolation ON hrms.tenant_profiles;
CREATE POLICY tenant_profiles_isolation ON hrms.tenant_profiles
    USING (
        tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::uuid
        OR current_setting('app.is_super_admin', true) = 'true'
    )
    WITH CHECK (
        tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::uuid
        OR current_setting('app.is_super_admin', true) = 'true'
    );

DROP POLICY IF EXISTS tenant_settings_isolation ON hrms.tenant_settings;
CREATE POLICY tenant_settings_isolation ON hrms.tenant_settings
    USING (
        tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::uuid
        OR current_setting('app.is_super_admin', true) = 'true'
    )
    WITH CHECK (
        tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::uuid
        OR current_setting('app.is_super_admin', true) = 'true'
    );
