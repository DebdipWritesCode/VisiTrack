-- name: CreateOTP :one
INSERT INTO otps (
  phone_number, otp_code, expires_at
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetOTPByPhone :one
SELECT * FROM otps
WHERE phone_number = $1
ORDER BY created_at DESC
LIMIT 1;

-- name: DeleteOTPByPhone :exec
DELETE FROM otps
WHERE phone_number = $1;

-- name: DeleteExpiredOTPs :exec
DELETE FROM otps
WHERE expires_at < NOW();
