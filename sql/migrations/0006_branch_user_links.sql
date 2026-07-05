ALTER TABLE hrms.branches
    ADD COLUMN IF NOT EXISTS branch_manager_user_id UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS hr_user_id UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS accounts_user_id UUID REFERENCES auth.users(id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS branches_manager_user_idx ON hrms.branches(tenant_id, branch_manager_user_id) WHERE branch_manager_user_id IS NOT NULL AND NOT inactive;
CREATE INDEX IF NOT EXISTS branches_hr_user_idx ON hrms.branches(tenant_id, hr_user_id) WHERE hr_user_id IS NOT NULL AND NOT inactive;
CREATE INDEX IF NOT EXISTS branches_accounts_user_idx ON hrms.branches(tenant_id, accounts_user_id) WHERE accounts_user_id IS NOT NULL AND NOT inactive;
