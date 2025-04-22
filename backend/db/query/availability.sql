-- name: CreateAvailabilitySlot :one
INSERT INTO availability (
  user_id, day_of_week, start_time, end_time
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: GetAvailabilityByUser :many
SELECT * FROM availability
WHERE user_id = $1
ORDER BY day_of_week, start_time;

-- name: DeleteAvailabilitySlot :exec
DELETE FROM availability
WHERE user_id = $1
  AND day_of_week = $2
  AND start_time = $3
  AND end_time = $4;

-- name: DeleteAvailabilityByUser :exec
DELETE FROM availability
WHERE user_id = $1;
