-- name: CreateUser :one
INSERT INTO users (
  phone_number, first_name, last_name, role
) VALUES (
  $1, $2, $3, COALESCE($4, 'user')
)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;

-- name: GetUserByPhone :one
SELECT * FROM users
WHERE phone_number = $1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateUserName :one
UPDATE users
SET first_name = $2,
    last_name = $3
WHERE id = $1
RETURNING *;

-- name: UpdateUserRole :one
UPDATE users
SET role = $2
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
