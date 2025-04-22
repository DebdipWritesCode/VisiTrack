-- name: CreateAppointmentLog :one
INSERT INTO appointment_logs (
  appointment_id, check_in_time, check_out_time
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetAppointmentLogByAppointmentID :one
SELECT * FROM appointment_logs
WHERE appointment_id = $1;

-- name: UpdateCheckInTime :one
UPDATE appointment_logs
SET check_in_time = $2
WHERE appointment_id = $1
RETURNING *;

-- name: UpdateCheckOutTime :one
UPDATE appointment_logs
SET check_out_time = $2
WHERE appointment_id = $1
RETURNING *;

-- name: DeleteAppointmentLog :exec
DELETE FROM appointment_logs
WHERE appointment_id = $1;
