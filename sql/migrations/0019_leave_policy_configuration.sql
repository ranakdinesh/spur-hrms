ALTER TABLE hrms.leave_types ADD COLUMN IF NOT EXISTS is_enabled BOOLEAN NOT NULL DEFAULT TRUE;

ALTER TABLE hrms.leave_policy_template_rules ADD COLUMN IF NOT EXISTS annual_entitlement NUMERIC(6,2) NOT NULL DEFAULT 0;
ALTER TABLE hrms.leave_policy_template_rules ADD COLUMN IF NOT EXISTS min_request_days NUMERIC(6,2) NOT NULL DEFAULT 0;
ALTER TABLE hrms.leave_policy_template_rules ADD COLUMN IF NOT EXISTS max_request_days NUMERIC(6,2);
ALTER TABLE hrms.leave_policy_template_rules ADD COLUMN IF NOT EXISTS max_requests_per_year INT NOT NULL DEFAULT 0;
ALTER TABLE hrms.leave_policy_template_rules ADD COLUMN IF NOT EXISTS accrual_day INT NOT NULL DEFAULT 1;
ALTER TABLE hrms.leave_policy_template_rules ADD COLUMN IF NOT EXISTS lapse_unutilized BOOLEAN NOT NULL DEFAULT TRUE;
ALTER TABLE hrms.leave_policy_template_rules ADD COLUMN IF NOT EXISTS encashment_limit NUMERIC(6,2) NOT NULL DEFAULT 0;
ALTER TABLE hrms.leave_policy_template_rules ADD COLUMN IF NOT EXISTS encashment_payable_percent NUMERIC(5,2) NOT NULL DEFAULT 100;
ALTER TABLE hrms.leave_policy_template_rules ADD COLUMN IF NOT EXISTS allow_half_day BOOLEAN NOT NULL DEFAULT TRUE;
ALTER TABLE hrms.leave_policy_template_rules ADD COLUMN IF NOT EXISTS requires_approval BOOLEAN NOT NULL DEFAULT TRUE;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'leave_policy_template_rules_request_limits_check') THEN
        ALTER TABLE hrms.leave_policy_template_rules ADD CONSTRAINT leave_policy_template_rules_request_limits_check
        CHECK (
            annual_entitlement >= 0
            AND min_request_days >= 0
            AND (max_request_days IS NULL OR max_request_days >= min_request_days)
            AND max_requests_per_year >= 0
            AND accrual_day BETWEEN 1 AND 31
            AND encashment_limit >= 0
            AND encashment_payable_percent >= 0
            AND encashment_payable_percent <= 100
        );
    END IF;
END $$;
