-- name: CreateAppointment :one
INSERT INTO appointments (
  visitor_id, host_id, appointment_date, start_time, end_time, status, qr_code
) VALUES (
  $1, $2, $3, $4, $5, COALESCE($6, 'pending'), $7
)
RETURNING *;

-- name: GetAppointmentByID :one
SELECT * FROM appointments
WHERE id = $1;

-- name: ListAppointmentsByVisitor :many
SELECT * FROM appointments
WHERE visitor_id = $1
ORDER BY appointment_date DESC, start_time DESC;

-- name: ListAppointmentsByHost :many
SELECT * FROM appointments
WHERE host_id = $1
ORDER BY appointment_date DESC, start_time DESC;

-- name: ListAppointmentsByDate :many
SELECT * FROM appointments
WHERE appointment_date = $1
ORDER BY start_time;

-- name: UpdateAppointmentStatus :one
UPDATE appointments
SET status = $2
WHERE id = $1
RETURNING *;

-- name: CancelAppointment :one
UPDATE appointments
SET status = 'cancelled'
WHERE id = $1
RETURNING *;

-- name: DeleteAppointment :exec
DELETE FROM appointments
WHERE id = $1;
