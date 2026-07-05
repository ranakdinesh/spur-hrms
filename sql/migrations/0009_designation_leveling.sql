ALTER TABLE hrms.designations
ADD COLUMN IF NOT EXISTS level_code VARCHAR(32);

ALTER TABLE hrms.designations
ADD COLUMN IF NOT EXISTS seniority_rank INTEGER;

UPDATE hrms.designations
SET seniority_rank = 100
WHERE seniority_rank IS NULL;

UPDATE hrms.designations
SET level_code = 'L' || seniority_rank::TEXT
WHERE level_code IS NULL OR btrim(level_code) = '';

ALTER TABLE hrms.designations
ALTER COLUMN seniority_rank SET NOT NULL;

ALTER TABLE hrms.designations
ALTER COLUMN level_code SET NOT NULL;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'designations_seniority_rank_range'
          AND conrelid = 'hrms.designations'::regclass
    ) THEN
        ALTER TABLE hrms.designations
        ADD CONSTRAINT designations_seniority_rank_range
        CHECK (seniority_rank BETWEEN 1 AND 9999);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS designations_tenant_seniority_rank_idx
ON hrms.designations(tenant_id, seniority_rank DESC, name ASC)
WHERE NOT inactive;
