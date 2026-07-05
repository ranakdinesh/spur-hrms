INSERT INTO hrms.employment_types (tenant_id, name)
SELECT tenants.id, defaults.name
FROM auth.tenants tenants
CROSS JOIN (VALUES
    ('Permanent'),
    ('Probation'),
    ('Contract'),
    ('Consultant'),
    ('Intern'),
    ('Part-time')
) AS defaults(name)
WHERE NOT EXISTS (
    SELECT 1
    FROM hrms.employment_types existing
    WHERE existing.tenant_id = tenants.id
      AND lower(existing.name) = lower(defaults.name)
      AND NOT existing.inactive
);

INSERT INTO hrms.marital_statuses (tenant_id, name)
SELECT tenants.id, defaults.name
FROM auth.tenants tenants
CROSS JOIN (VALUES
    ('Single'),
    ('Married'),
    ('Divorced'),
    ('Widowed'),
    ('Separated')
) AS defaults(name)
WHERE NOT EXISTS (
    SELECT 1
    FROM hrms.marital_statuses existing
    WHERE existing.tenant_id = tenants.id
      AND lower(existing.name) = lower(defaults.name)
      AND NOT existing.inactive
);
