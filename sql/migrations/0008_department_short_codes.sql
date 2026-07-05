ALTER TABLE hrms.departments
ADD COLUMN IF NOT EXISTS short_code VARCHAR(32);

WITH generated AS (
    SELECT
        id,
        COALESCE(NULLIF(upper(left(regexp_replace(name, '[^A-Za-z0-9]+', '', 'g'), 20)), ''), 'DEPT') ||
            '-' || upper(left(replace(id::text, '-', ''), 8)) AS generated_short_code
    FROM hrms.departments
    WHERE short_code IS NULL OR btrim(short_code) = ''
)
UPDATE hrms.departments departments
SET short_code = generated.generated_short_code
FROM generated
WHERE departments.id = generated.id;

ALTER TABLE hrms.departments
ALTER COLUMN short_code SET NOT NULL;

CREATE UNIQUE INDEX IF NOT EXISTS departments_tenant_short_code_idx
ON hrms.departments(tenant_id, lower(short_code))
WHERE NOT inactive;
