-- name: ListUserOTPsByUser :many
SELECT * FROM hrms.user_otps
WHERE tenant_id = $1 AND user_id = $2 AND NOT inactive
ORDER BY created_at DESC;

-- name: GetActiveUserOTP :one
SELECT * FROM hrms.user_otps
WHERE tenant_id = $1 AND user_id = $2 AND otp_for = $3 AND otp = $4
  AND NOT is_used AND expires_at > NOW() AND NOT inactive
ORDER BY created_at DESC
LIMIT 1;

-- name: MarkUserOTPUsed :exec
UPDATE hrms.user_otps
SET is_used = TRUE, updated_by = $3
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: SoftDeleteUserOTP :exec
UPDATE hrms.user_otps
SET inactive = TRUE, updated_by = $3
WHERE tenant_id = $1 AND id = $2;
