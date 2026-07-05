ALTER TABLE hrms.salary_template_items
    DROP CONSTRAINT IF EXISTS salary_template_items_type_check;

ALTER TABLE hrms.salary_template_items
    ADD CONSTRAINT salary_template_items_type_check
    CHECK (item_type IN ('earning','deduction','employer_contribution','reimbursement'));
