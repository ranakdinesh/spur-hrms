-- Hrms module — initial migration
-- Schema: public
-- Convention: id UUID PK, tenant_id UUID NOT NULL, timestamps

CREATE TABLE IF NOT EXISTS hrmss (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID        NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()

    -- TODO: add your domain columns
    -- Example:
    -- name        TEXT        NOT NULL,
    -- status      TEXT        NOT NULL DEFAULT 'active',
    -- metadata    JSONB       NOT NULL DEFAULT '{}'
);

-- Required: index on tenant_id — every query filters by it
CREATE INDEX IF NOT EXISTS idx_hrmss_tenant ON hrmss(tenant_id);

-- Optional: auto-update updated_at
CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER LANGUAGE plpgsql AS $$
BEGIN NEW.updated_at = NOW(); RETURN NEW; END; $$;

CREATE TRIGGER hrmss_updated_at
    BEFORE UPDATE ON hrmss
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();
