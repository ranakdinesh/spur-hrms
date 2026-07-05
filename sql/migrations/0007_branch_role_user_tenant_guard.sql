CREATE OR REPLACE FUNCTION hrms.ensure_branch_role_users_same_tenant()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.branch_manager_user_id IS NOT NULL
        AND NOT EXISTS (
            SELECT 1
            FROM auth.users
            WHERE id = NEW.branch_manager_user_id
              AND tenant_id = NEW.tenant_id
        )
    THEN
        RAISE EXCEPTION 'branch manager user must belong to the branch tenant'
            USING ERRCODE = '23514';
    END IF;

    IF NEW.hr_user_id IS NOT NULL
        AND NOT EXISTS (
            SELECT 1
            FROM auth.users
            WHERE id = NEW.hr_user_id
              AND tenant_id = NEW.tenant_id
        )
    THEN
        RAISE EXCEPTION 'branch HR user must belong to the branch tenant'
            USING ERRCODE = '23514';
    END IF;

    IF NEW.accounts_user_id IS NOT NULL
        AND NOT EXISTS (
            SELECT 1
            FROM auth.users
            WHERE id = NEW.accounts_user_id
              AND tenant_id = NEW.tenant_id
        )
    THEN
        RAISE EXCEPTION 'branch accounts user must belong to the branch tenant'
            USING ERRCODE = '23514';
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS branches_role_users_same_tenant ON hrms.branches;
CREATE TRIGGER branches_role_users_same_tenant
BEFORE INSERT OR UPDATE OF tenant_id, branch_manager_user_id, hr_user_id, accounts_user_id
ON hrms.branches
FOR EACH ROW
EXECUTE FUNCTION hrms.ensure_branch_role_users_same_tenant();
