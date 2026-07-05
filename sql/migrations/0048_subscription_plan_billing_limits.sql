-- HRMS-015D: enforce sensible subscription pricing and employee caps.

ALTER TABLE hrms.subscription_plans
    ADD COLUMN IF NOT EXISTS minimum_amount NUMERIC(12,2) NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS included_employees INT NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS overage_amount NUMERIC(12,2) NOT NULL DEFAULT 0;

ALTER TABLE hrms.subscription_plans
    DROP CONSTRAINT IF EXISTS subscription_plans_price_basis_check;

ALTER TABLE hrms.subscription_plans
    ADD CONSTRAINT subscription_plans_price_basis_check
    CHECK (price_basis IN ('per_employee','package_plus_overage','flat','custom_quote'));

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'subscription_plans_minimum_nonnegative'
    ) THEN
        ALTER TABLE hrms.subscription_plans
            ADD CONSTRAINT subscription_plans_minimum_nonnegative
            CHECK (minimum_amount >= 0);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'subscription_plans_included_employees_nonnegative'
    ) THEN
        ALTER TABLE hrms.subscription_plans
            ADD CONSTRAINT subscription_plans_included_employees_nonnegative
            CHECK (included_employees >= 0);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'subscription_plans_overage_nonnegative'
    ) THEN
        ALTER TABLE hrms.subscription_plans
            ADD CONSTRAINT subscription_plans_overage_nonnegative
            CHECK (overage_amount >= 0);
    END IF;
END $$;

INSERT INTO hrms.subscription_plans (
    code,
    name,
    description,
    price_amount,
    price_basis,
    minimum_amount,
    included_employees,
    overage_amount,
    currency_code,
    billing_cycle,
    employee_limit,
    trial_days,
    visibility,
    is_active
)
SELECT *
FROM (
    VALUES
        ('GROWTH_100', 'Growth 100', 'Public package plan: ₹5,000/month includes 100 employees, then ₹45 per extra employee up to 150.', 5000.00, 'package_plus_overage', 5000.00, 100, 45.00, 'INR', 'monthly', 150, 14, 'public', TRUE),
        ('BUSINESS_250', 'Business 250', 'Public package plan: ₹11,000/month includes 250 employees, then ₹40 per extra employee up to 350.', 11000.00, 'package_plus_overage', 11000.00, 250, 40.00, 'INR', 'monthly', 350, 14, 'public', TRUE),
        ('ENTERPRISE_500', 'Enterprise 500', 'Public enterprise package: ₹20,000/month includes 500 employees, then ₹35 per extra employee up to 750.', 20000.00, 'package_plus_overage', 20000.00, 500, 35.00, 'INR', 'monthly', 750, 14, 'public', TRUE),
        ('NEGOTIATED_GROWTH_100', 'Negotiated Growth 100', 'Internal negotiated package: ₹4,500/month includes 100 employees, then ₹39 per extra employee up to 150.', 4500.00, 'package_plus_overage', 4500.00, 100, 39.00, 'INR', 'monthly', 150, 14, 'internal', TRUE)
) AS seed(code, name, description, price_amount, price_basis, minimum_amount, included_employees, overage_amount, currency_code, billing_cycle, employee_limit, trial_days, visibility, is_active)
WHERE NOT EXISTS (
    SELECT 1
    FROM hrms.subscription_plans existing
    WHERE lower(existing.code) = lower(seed.code)
      AND NOT existing.inactive
);

UPDATE hrms.subscription_plans
SET name = seed.name,
    description = seed.description,
    price_amount = seed.price_amount,
    price_basis = seed.price_basis,
    minimum_amount = seed.minimum_amount,
    included_employees = seed.included_employees,
    overage_amount = seed.overage_amount,
    currency_code = seed.currency_code,
    billing_cycle = seed.billing_cycle,
    employee_limit = seed.employee_limit,
    trial_days = seed.trial_days,
    visibility = seed.visibility,
    is_active = TRUE,
    updated_at = NOW()
FROM (
    VALUES
        ('STARTER', 'Starter', 'Public plan: ₹39 per employee/month with ₹999 minimum monthly commitment, capped at 50 employees.', 39.00, 'per_employee', 999.00, 0, 39.00, 'INR', 'monthly', 50, 14, 'public'),
        ('GROWTH_100', 'Growth 100', 'Public package plan: ₹5,000/month includes 100 employees, then ₹45 per extra employee up to 150.', 5000.00, 'package_plus_overage', 5000.00, 100, 45.00, 'INR', 'monthly', 150, 14, 'public'),
        ('BUSINESS_250', 'Business 250', 'Public package plan: ₹11,000/month includes 250 employees, then ₹40 per extra employee up to 350.', 11000.00, 'package_plus_overage', 11000.00, 250, 40.00, 'INR', 'monthly', 350, 14, 'public'),
        ('ENTERPRISE_500', 'Enterprise 500', 'Public enterprise package: ₹20,000/month includes 500 employees, then ₹35 per extra employee up to 750.', 20000.00, 'package_plus_overage', 20000.00, 500, 35.00, 'INR', 'monthly', 750, 14, 'public'),
        ('NEGOTIATED_STARTUP', 'Negotiated Startup', 'Internal negotiated plan: ₹29 per employee/month with ₹699 minimum monthly commitment, capped at 35 employees.', 29.00, 'per_employee', 699.00, 0, 29.00, 'INR', 'monthly', 35, 14, 'internal'),
        ('NEGOTIATED_GROWTH_100', 'Negotiated Growth 100', 'Internal negotiated package: ₹4,500/month includes 100 employees, then ₹39 per extra employee up to 150.', 4500.00, 'package_plus_overage', 4500.00, 100, 39.00, 'INR', 'monthly', 150, 14, 'internal'),
        ('ENTERPRISE_CUSTOM', 'Enterprise Custom', 'Internal quote-based enterprise plan for custom contracts, annual prepay, bundled services, or non-standard limits.', 0.00, 'custom_quote', 0.00, 0, 0.00, 'INR', 'custom', 1000, 0, 'internal')
) AS seed(code, name, description, price_amount, price_basis, minimum_amount, included_employees, overage_amount, currency_code, billing_cycle, employee_limit, trial_days, visibility)
WHERE lower(hrms.subscription_plans.code) = lower(seed.code)
  AND NOT hrms.subscription_plans.inactive;

UPDATE hrms.subscription_plans
SET is_active = FALSE,
    updated_at = NOW()
WHERE code IN ('PROFESSIONAL', 'CUSTOM', 'NEGOTIATED_GROWTH')
  AND NOT inactive;
