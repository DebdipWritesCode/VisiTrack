-- name: CreateAppointmentStats :one
INSERT INTO appointment_stats (
  user_id, total_appointments
) VALUES (
  $1, $2
)
RETURNING *;

-- name: GetAppointmentStatsByUserID :one
SELECT * FROM appointment_stats
WHERE user_id = $1;

-- name: IncrementAppointmentCount :one
UPDATE appointment_stats
SET total_appointments = total_appointments + 1
WHERE user_id = $1
RETURNING *;

-- name: DecrementAppointmentCount :one
UPDATE appointment_stats
SET total_appointments = GREATEST(total_appointments - 1, 0)
WHERE user_id = $1
RETURNING *;

-- name: ResetAppointmentCount :one
UPDATE appointment_stats
SET total_appointments = 0
WHERE user_id = $1
RETURNING *;

-- name: DeleteAppointmentStats :exec
DELETE FROM appointment_stats
WHERE user_id = $1;
