-- name: CreateUser :one
INSERT INTO users (
  phone_number, first_name, last_name, role
) VALUES (
  $1, $2, $3, $4
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

-- name: GetTopPopularUsers :many
SELECT 
  id, 
  first_name, 
  last_name, 
  phone_number, 
  appointments_hosted, 
  appointments_visited,
  (appointments_hosted + appointments_visited) AS total_appointments
FROM users
ORDER BY total_appointments DESC
LIMIT 5;

-- name: GetTotalAppointmentsHosted :one
SELECT appointments_hosted
FROM users
WHERE id = $1;

-- name: GetTotalAppointmentsVisited :one
SELECT appointments_visited
FROM users
WHERE id = $1;
