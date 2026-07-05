CREATE TABLE IF NOT EXISTS hrms.designation_level_codes (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    code        VARCHAR(32) NOT NULL,
    label       VARCHAR(255) NOT NULL,
    description TEXT,
    sort_order  INTEGER NOT NULL DEFAULT 0,
    inactive    BOOLEAN NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by  UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by  UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT designation_level_codes_code_required CHECK (btrim(code) <> ''),
    CONSTRAINT designation_level_codes_label_required CHECK (btrim(label) <> '')
);

CREATE UNIQUE INDEX IF NOT EXISTS designation_level_codes_tenant_code_idx
ON hrms.designation_level_codes(tenant_id, code);

CREATE INDEX IF NOT EXISTS designation_level_codes_tenant_sort_idx
ON hrms.designation_level_codes(tenant_id, sort_order ASC, code ASC)
WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.designation_seniority_ranks (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    rank_value  INTEGER NOT NULL,
    label       VARCHAR(255) NOT NULL,
    description TEXT,
    sort_order  INTEGER NOT NULL DEFAULT 0,
    inactive    BOOLEAN NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by  UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by  UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT designation_seniority_ranks_value_range CHECK (rank_value BETWEEN 1 AND 9999),
    CONSTRAINT designation_seniority_ranks_label_required CHECK (btrim(label) <> '')
);

CREATE UNIQUE INDEX IF NOT EXISTS designation_seniority_ranks_tenant_rank_idx
ON hrms.designation_seniority_ranks(tenant_id, rank_value);

CREATE INDEX IF NOT EXISTS designation_seniority_ranks_tenant_sort_idx
ON hrms.designation_seniority_ranks(tenant_id, sort_order ASC, rank_value ASC)
WHERE NOT inactive;

ALTER TABLE hrms.designation_level_codes ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.designation_seniority_ranks ENABLE ROW LEVEL SECURITY;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE schemaname = 'hrms' AND tablename = 'designation_level_codes' AND policyname = 'designation_level_codes_tenant_policy') THEN
        CREATE POLICY designation_level_codes_tenant_policy ON hrms.designation_level_codes
        USING (tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::uuid OR current_setting('app.is_super_admin', true) = 'true')
        WITH CHECK (tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::uuid OR current_setting('app.is_super_admin', true) = 'true');
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE schemaname = 'hrms' AND tablename = 'designation_seniority_ranks' AND policyname = 'designation_seniority_ranks_tenant_policy') THEN
        CREATE POLICY designation_seniority_ranks_tenant_policy ON hrms.designation_seniority_ranks
        USING (tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::uuid OR current_setting('app.is_super_admin', true) = 'true')
        WITH CHECK (tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::uuid OR current_setting('app.is_super_admin', true) = 'true');
    END IF;
END $$;

INSERT INTO hrms.designation_level_codes (tenant_id, code, label, description, sort_order)
SELECT tenants.id, defaults.code, defaults.label, defaults.description, defaults.sort_order
FROM auth.tenants tenants
CROSS JOIN (VALUES
    ('L1', 'Entry / Trainee', 'Entry-level or trainee role.', 10),
    ('L2', 'Junior', 'Early-career role with guided execution.', 20),
    ('L3', 'Associate / Executive', 'Independent contributor with defined responsibilities.', 30),
    ('L4', 'Senior / Specialist', 'Experienced contributor or first-line specialist.', 40),
    ('L5', 'Lead / Manager', 'Team lead or manager with ownership of outcomes.', 50),
    ('L6', 'Senior Manager / Director', 'Senior leadership for a function or department.', 60),
    ('L7', 'Executive', 'Executive leadership with organization-wide impact.', 70),
    ('M1', 'Manager I', 'First manager level.', 80),
    ('M2', 'Manager II', 'Experienced manager level.', 90),
    ('M3', 'Senior Manager', 'Senior manager level.', 100),
    ('M4', 'Director', 'Director level.', 110),
    ('M5', 'Executive Manager', 'Executive manager level.', 120)
) AS defaults(code, label, description, sort_order)
ON CONFLICT (tenant_id, code) DO NOTHING;

INSERT INTO hrms.designation_seniority_ranks (tenant_id, rank_value, label, description, sort_order)
SELECT tenants.id, defaults.rank_value, defaults.label, defaults.description, defaults.sort_order
FROM auth.tenants tenants
CROSS JOIN (VALUES
    (100, 'Entry / Trainee', 'Lowest organization seniority.', 10),
    (200, 'Junior', 'Junior role seniority.', 20),
    (300, 'Associate / Executive', 'Standard independent contributor seniority.', 30),
    (400, 'Senior / Specialist', 'Senior contributor seniority.', 40),
    (500, 'Lead / Manager', 'Lead or manager seniority.', 50),
    (600, 'Senior Manager / Director', 'Senior leadership seniority.', 60),
    (700, 'Executive', 'Executive seniority.', 70)
) AS defaults(rank_value, label, description, sort_order)
ON CONFLICT (tenant_id, rank_value) DO NOTHING;
