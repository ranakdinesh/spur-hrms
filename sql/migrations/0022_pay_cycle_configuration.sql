-- HRMS-034: configurable payroll pay cycle foundation.
-- India payroll varies by attendance cutoff, payout date, statutory applicability,
-- arrears handling, and finance approval controls. Keep this additive so existing
-- tenant pay_cycles rows continue to work.

ALTER TABLE hrms.pay_cycles
    ADD COLUMN IF NOT EXISTS name VARCHAR(120) NOT NULL DEFAULT 'Monthly Payroll',
    ADD COLUMN IF NOT EXISTS attendance_source VARCHAR(40) NOT NULL DEFAULT 'attendance',
    ADD COLUMN IF NOT EXISTS attendance_period_type VARCHAR(40) NOT NULL DEFAULT 'current_month',
    ADD COLUMN IF NOT EXISTS attendance_cutoff_day INT NOT NULL DEFAULT 25,
    ADD COLUMN IF NOT EXISTS payout_timing VARCHAR(40) NOT NULL DEFAULT 'same_month',
    ADD COLUMN IF NOT EXISTS payout_offset_days INT NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS include_weekly_offs BOOLEAN NOT NULL DEFAULT TRUE,
    ADD COLUMN IF NOT EXISTS include_holidays BOOLEAN NOT NULL DEFAULT TRUE,
    ADD COLUMN IF NOT EXISTS prorate_joining_exit BOOLEAN NOT NULL DEFAULT TRUE,
    ADD COLUMN IF NOT EXISTS proration_basis VARCHAR(40) NOT NULL DEFAULT 'calendar_days',
    ADD COLUMN IF NOT EXISTS allow_arrears BOOLEAN NOT NULL DEFAULT TRUE,
    ADD COLUMN IF NOT EXISTS arrears_mode VARCHAR(40) NOT NULL DEFAULT 'next_cycle',
    ADD COLUMN IF NOT EXISTS allow_negative_net_pay BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS overtime_component_code VARCHAR(60),
    ADD COLUMN IF NOT EXISTS lwp_component_code VARCHAR(60) NOT NULL DEFAULT 'lwp',
    ADD COLUMN IF NOT EXISTS rounding_mode VARCHAR(40) NOT NULL DEFAULT 'nearest_rupee',
    ADD COLUMN IF NOT EXISTS payment_mode VARCHAR(40) NOT NULL DEFAULT 'bank_transfer',
    ADD COLUMN IF NOT EXISTS payment_file_format VARCHAR(40) NOT NULL DEFAULT 'bank_csv',
    ADD COLUMN IF NOT EXISTS requires_approval BOOLEAN NOT NULL DEFAULT TRUE,
    ADD COLUMN IF NOT EXISTS auto_lock_after_approval BOOLEAN NOT NULL DEFAULT TRUE,
    ADD COLUMN IF NOT EXISTS payroll_lock_day INT,
    ADD COLUMN IF NOT EXISTS pf_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    ADD COLUMN IF NOT EXISTS pf_employee_rate NUMERIC(7,4) NOT NULL DEFAULT 12.0000,
    ADD COLUMN IF NOT EXISTS pf_employer_rate NUMERIC(7,4) NOT NULL DEFAULT 12.0000,
    ADD COLUMN IF NOT EXISTS pf_wage_ceiling NUMERIC(12,2) NOT NULL DEFAULT 15000.00,
    ADD COLUMN IF NOT EXISTS pf_apply_ceiling BOOLEAN NOT NULL DEFAULT TRUE,
    ADD COLUMN IF NOT EXISTS esi_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    ADD COLUMN IF NOT EXISTS esi_employee_rate NUMERIC(7,4) NOT NULL DEFAULT 0.7500,
    ADD COLUMN IF NOT EXISTS esi_employer_rate NUMERIC(7,4) NOT NULL DEFAULT 3.2500,
    ADD COLUMN IF NOT EXISTS esi_wage_ceiling NUMERIC(12,2) NOT NULL DEFAULT 21000.00,
    ADD COLUMN IF NOT EXISTS professional_tax_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    ADD COLUMN IF NOT EXISTS tds_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    ADD COLUMN IF NOT EXISTS country_code VARCHAR(2) NOT NULL DEFAULT 'IN',
    ADD COLUMN IF NOT EXISTS state_code VARCHAR(10),
    ADD COLUMN IF NOT EXISTS notes TEXT;

ALTER TABLE hrms.pay_cycles
    DROP CONSTRAINT IF EXISTS pay_cycles_extended_check;

ALTER TABLE hrms.pay_cycles
    ADD CONSTRAINT pay_cycles_extended_check CHECK (
        cycle_type IN ('monthly','semi_monthly','weekly','bi_weekly','custom')
        AND attendance_source IN ('attendance','manual','import','none')
        AND attendance_period_type IN ('current_month','previous_month','custom_days')
        AND payout_timing IN ('same_month','next_month','fixed_offset')
        AND proration_basis IN ('calendar_days','working_days','fixed_26','fixed_30')
        AND arrears_mode IN ('same_cycle','next_cycle','manual')
        AND rounding_mode IN ('none','nearest_rupee','ceil_rupee','floor_rupee','two_decimals')
        AND payment_mode IN ('bank_transfer','cash','cheque','upi','mixed')
        AND payment_file_format IN ('bank_csv','bank_xlsx','nach','none','custom')
        AND attendance_cutoff_day BETWEEN 1 AND 31
        AND payout_offset_days BETWEEN 0 AND 31
        AND (payroll_lock_day IS NULL OR payroll_lock_day BETWEEN 1 AND 31)
        AND pf_employee_rate >= 0 AND pf_employer_rate >= 0 AND pf_wage_ceiling >= 0
        AND esi_employee_rate >= 0 AND esi_employer_rate >= 0 AND esi_wage_ceiling >= 0
    );
