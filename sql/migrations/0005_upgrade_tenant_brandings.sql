-- Repair/upgrade tenant branding schema for databases that applied an older 0004 migration.

ALTER TABLE hrms.tenant_brandings ADD COLUMN IF NOT EXISTS tenant_id UUID;
ALTER TABLE hrms.tenant_brandings ADD COLUMN IF NOT EXISTS logo_path TEXT;
ALTER TABLE hrms.tenant_brandings ADD COLUMN IF NOT EXISTS favicon_path TEXT;
ALTER TABLE hrms.tenant_brandings ADD COLUMN IF NOT EXISTS layout VARCHAR(40) DEFAULT 'vertical';
ALTER TABLE hrms.tenant_brandings ADD COLUMN IF NOT EXISTS color_mode VARCHAR(20) DEFAULT 'light';
ALTER TABLE hrms.tenant_brandings ADD COLUMN IF NOT EXISTS sidebar_size VARCHAR(20) DEFAULT 'default';
ALTER TABLE hrms.tenant_brandings ADD COLUMN IF NOT EXISTS layout_width VARCHAR(20) DEFAULT 'fluid';
ALTER TABLE hrms.tenant_brandings ADD COLUMN IF NOT EXISTS card_layout VARCHAR(20) DEFAULT 'bordered';
ALTER TABLE hrms.tenant_brandings ADD COLUMN IF NOT EXISTS theme_color VARCHAR(7) DEFAULT '#588368';
ALTER TABLE hrms.tenant_brandings ADD COLUMN IF NOT EXISTS primary_color VARCHAR(7) DEFAULT '#588368';
ALTER TABLE hrms.tenant_brandings ADD COLUMN IF NOT EXISTS secondary_color VARCHAR(7) DEFAULT '#2f6f7d';
ALTER TABLE hrms.tenant_brandings ADD COLUMN IF NOT EXISTS tertiary_color VARCHAR(7) DEFAULT '#e87839';
ALTER TABLE hrms.tenant_brandings ADD COLUMN IF NOT EXISTS topbar_color VARCHAR(7) DEFAULT '#ffffff';
ALTER TABLE hrms.tenant_brandings ADD COLUMN IF NOT EXISTS sidebar_color VARCHAR(7) DEFAULT '#111827';
ALTER TABLE hrms.tenant_brandings ADD COLUMN IF NOT EXISTS topbar_background TEXT DEFAULT 'none';
ALTER TABLE hrms.tenant_brandings ADD COLUMN IF NOT EXISTS sidebar_background TEXT DEFAULT 'solid';
ALTER TABLE hrms.tenant_brandings ADD COLUMN IF NOT EXISTS font_family VARCHAR(120) DEFAULT 'Inter, sans-serif';
ALTER TABLE hrms.tenant_brandings ADD COLUMN IF NOT EXISTS preloader BOOLEAN DEFAULT TRUE;
ALTER TABLE hrms.tenant_brandings ADD COLUMN IF NOT EXISTS inactive BOOLEAN DEFAULT FALSE;
ALTER TABLE hrms.tenant_brandings ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ DEFAULT NOW();
ALTER TABLE hrms.tenant_brandings ADD COLUMN IF NOT EXISTS created_by UUID;
ALTER TABLE hrms.tenant_brandings ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ DEFAULT NOW();
ALTER TABLE hrms.tenant_brandings ADD COLUMN IF NOT EXISTS updated_by UUID;

UPDATE hrms.tenant_brandings SET
    layout = COALESCE(layout, 'vertical'),
    color_mode = COALESCE(color_mode, 'light'),
    sidebar_size = COALESCE(sidebar_size, 'default'),
    layout_width = COALESCE(layout_width, 'fluid'),
    card_layout = COALESCE(card_layout, 'bordered'),
    theme_color = COALESCE(theme_color, '#588368'),
    primary_color = COALESCE(primary_color, '#588368'),
    secondary_color = COALESCE(secondary_color, '#2f6f7d'),
    tertiary_color = COALESCE(tertiary_color, '#e87839'),
    topbar_color = COALESCE(topbar_color, '#ffffff'),
    sidebar_color = COALESCE(sidebar_color, '#111827'),
    topbar_background = COALESCE(topbar_background, 'none'),
    sidebar_background = COALESCE(sidebar_background, 'solid'),
    font_family = COALESCE(font_family, 'Inter, sans-serif'),
    preloader = COALESCE(preloader, TRUE),
    inactive = COALESCE(inactive, FALSE),
    created_at = COALESCE(created_at, NOW()),
    updated_at = COALESCE(updated_at, NOW());

ALTER TABLE hrms.tenant_brandings ALTER COLUMN layout SET DEFAULT 'vertical';
ALTER TABLE hrms.tenant_brandings ALTER COLUMN layout SET NOT NULL;
ALTER TABLE hrms.tenant_brandings ALTER COLUMN color_mode SET DEFAULT 'light';
ALTER TABLE hrms.tenant_brandings ALTER COLUMN color_mode SET NOT NULL;
ALTER TABLE hrms.tenant_brandings ALTER COLUMN sidebar_size SET DEFAULT 'default';
ALTER TABLE hrms.tenant_brandings ALTER COLUMN sidebar_size SET NOT NULL;
ALTER TABLE hrms.tenant_brandings ALTER COLUMN layout_width SET DEFAULT 'fluid';
ALTER TABLE hrms.tenant_brandings ALTER COLUMN layout_width SET NOT NULL;
ALTER TABLE hrms.tenant_brandings ALTER COLUMN card_layout SET DEFAULT 'bordered';
ALTER TABLE hrms.tenant_brandings ALTER COLUMN card_layout SET NOT NULL;
ALTER TABLE hrms.tenant_brandings ALTER COLUMN theme_color SET DEFAULT '#588368';
ALTER TABLE hrms.tenant_brandings ALTER COLUMN theme_color SET NOT NULL;
ALTER TABLE hrms.tenant_brandings ALTER COLUMN primary_color SET DEFAULT '#588368';
ALTER TABLE hrms.tenant_brandings ALTER COLUMN primary_color SET NOT NULL;
ALTER TABLE hrms.tenant_brandings ALTER COLUMN secondary_color SET DEFAULT '#2f6f7d';
ALTER TABLE hrms.tenant_brandings ALTER COLUMN secondary_color SET NOT NULL;
ALTER TABLE hrms.tenant_brandings ALTER COLUMN tertiary_color SET DEFAULT '#e87839';
ALTER TABLE hrms.tenant_brandings ALTER COLUMN tertiary_color SET NOT NULL;
ALTER TABLE hrms.tenant_brandings ALTER COLUMN topbar_color SET DEFAULT '#ffffff';
ALTER TABLE hrms.tenant_brandings ALTER COLUMN topbar_color SET NOT NULL;
ALTER TABLE hrms.tenant_brandings ALTER COLUMN sidebar_color SET DEFAULT '#111827';
ALTER TABLE hrms.tenant_brandings ALTER COLUMN sidebar_color SET NOT NULL;
ALTER TABLE hrms.tenant_brandings ALTER COLUMN topbar_background SET DEFAULT 'none';
ALTER TABLE hrms.tenant_brandings ALTER COLUMN topbar_background SET NOT NULL;
ALTER TABLE hrms.tenant_brandings ALTER COLUMN sidebar_background SET DEFAULT 'solid';
ALTER TABLE hrms.tenant_brandings ALTER COLUMN sidebar_background SET NOT NULL;
ALTER TABLE hrms.tenant_brandings ALTER COLUMN font_family SET DEFAULT 'Inter, sans-serif';
ALTER TABLE hrms.tenant_brandings ALTER COLUMN font_family SET NOT NULL;
ALTER TABLE hrms.tenant_brandings ALTER COLUMN preloader SET DEFAULT TRUE;
ALTER TABLE hrms.tenant_brandings ALTER COLUMN preloader SET NOT NULL;
ALTER TABLE hrms.tenant_brandings ALTER COLUMN inactive SET DEFAULT FALSE;
ALTER TABLE hrms.tenant_brandings ALTER COLUMN inactive SET NOT NULL;
ALTER TABLE hrms.tenant_brandings ALTER COLUMN created_at SET DEFAULT NOW();
ALTER TABLE hrms.tenant_brandings ALTER COLUMN created_at SET NOT NULL;
ALTER TABLE hrms.tenant_brandings ALTER COLUMN updated_at SET DEFAULT NOW();
ALTER TABLE hrms.tenant_brandings ALTER COLUMN updated_at SET NOT NULL;

CREATE UNIQUE INDEX IF NOT EXISTS tenant_brandings_tenant_id_uidx
    ON hrms.tenant_brandings (tenant_id);

CREATE INDEX IF NOT EXISTS tenant_brandings_tenant_idx
    ON hrms.tenant_brandings (tenant_id)
    WHERE NOT inactive;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'tenant_brandings_tenant_fk'
          AND conrelid = 'hrms.tenant_brandings'::regclass
    ) THEN
        ALTER TABLE hrms.tenant_brandings
            ADD CONSTRAINT tenant_brandings_tenant_fk
            FOREIGN KEY (tenant_id) REFERENCES auth.tenants(id) ON DELETE CASCADE NOT VALID;
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'tenant_brandings_layout_check' AND conrelid = 'hrms.tenant_brandings'::regclass) THEN
        ALTER TABLE hrms.tenant_brandings ADD CONSTRAINT tenant_brandings_layout_check CHECK (layout IN ('vertical','horizontal','detached','two-column')) NOT VALID;
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'tenant_brandings_color_mode_check' AND conrelid = 'hrms.tenant_brandings'::regclass) THEN
        ALTER TABLE hrms.tenant_brandings ADD CONSTRAINT tenant_brandings_color_mode_check CHECK (color_mode IN ('light','dark','system')) NOT VALID;
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'tenant_brandings_sidebar_size_check' AND conrelid = 'hrms.tenant_brandings'::regclass) THEN
        ALTER TABLE hrms.tenant_brandings ADD CONSTRAINT tenant_brandings_sidebar_size_check CHECK (sidebar_size IN ('default','compact','condensed','icon')) NOT VALID;
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'tenant_brandings_layout_width_check' AND conrelid = 'hrms.tenant_brandings'::regclass) THEN
        ALTER TABLE hrms.tenant_brandings ADD CONSTRAINT tenant_brandings_layout_width_check CHECK (layout_width IN ('fluid','boxed')) NOT VALID;
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'tenant_brandings_card_layout_check' AND conrelid = 'hrms.tenant_brandings'::regclass) THEN
        ALTER TABLE hrms.tenant_brandings ADD CONSTRAINT tenant_brandings_card_layout_check CHECK (card_layout IN ('bordered','shadow','plain')) NOT VALID;
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'tenant_brandings_theme_color_check' AND conrelid = 'hrms.tenant_brandings'::regclass) THEN
        ALTER TABLE hrms.tenant_brandings ADD CONSTRAINT tenant_brandings_theme_color_check CHECK (theme_color ~ '^#[0-9A-Fa-f]{6}$') NOT VALID;
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'tenant_brandings_primary_color_check' AND conrelid = 'hrms.tenant_brandings'::regclass) THEN
        ALTER TABLE hrms.tenant_brandings ADD CONSTRAINT tenant_brandings_primary_color_check CHECK (primary_color ~ '^#[0-9A-Fa-f]{6}$') NOT VALID;
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'tenant_brandings_secondary_color_check' AND conrelid = 'hrms.tenant_brandings'::regclass) THEN
        ALTER TABLE hrms.tenant_brandings ADD CONSTRAINT tenant_brandings_secondary_color_check CHECK (secondary_color ~ '^#[0-9A-Fa-f]{6}$') NOT VALID;
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'tenant_brandings_tertiary_color_check' AND conrelid = 'hrms.tenant_brandings'::regclass) THEN
        ALTER TABLE hrms.tenant_brandings ADD CONSTRAINT tenant_brandings_tertiary_color_check CHECK (tertiary_color ~ '^#[0-9A-Fa-f]{6}$') NOT VALID;
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'tenant_brandings_topbar_color_check' AND conrelid = 'hrms.tenant_brandings'::regclass) THEN
        ALTER TABLE hrms.tenant_brandings ADD CONSTRAINT tenant_brandings_topbar_color_check CHECK (topbar_color ~ '^#[0-9A-Fa-f]{6}$') NOT VALID;
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'tenant_brandings_sidebar_color_check' AND conrelid = 'hrms.tenant_brandings'::regclass) THEN
        ALTER TABLE hrms.tenant_brandings ADD CONSTRAINT tenant_brandings_sidebar_color_check CHECK (sidebar_color ~ '^#[0-9A-Fa-f]{6}$') NOT VALID;
    END IF;
END $$;

DROP TRIGGER IF EXISTS tenant_brandings_updated_at ON hrms.tenant_brandings;
CREATE TRIGGER tenant_brandings_updated_at
    BEFORE UPDATE ON hrms.tenant_brandings
    FOR EACH ROW EXECUTE FUNCTION hrms.update_updated_at();

ALTER TABLE hrms.tenant_brandings ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS tenant_brandings_tenant_isolation ON hrms.tenant_brandings;
CREATE POLICY tenant_brandings_tenant_isolation ON hrms.tenant_brandings
    USING (
        tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::uuid
        OR current_setting('app.is_super_admin', true) = 'true'
    )
    WITH CHECK (
        tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::uuid
        OR current_setting('app.is_super_admin', true) = 'true'
    );
