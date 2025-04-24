-- name: CreateAppointment :one
WITH updated_visitor AS (
  UPDATE "users"
  SET "appointments_visited" = "appointments_visited" + 1
  WHERE "id" = $1
  RETURNING "id"
),
updated_host AS (
  UPDATE "users"
  SET "appointments_hosted" = "appointments_hosted" + 1
  WHERE "id" = $2
  RETURNING "id"
)
INSERT INTO appointments (
  visitor_id, host_id, appointment_date, start_time, end_time, status, qr_code
) 
VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: GetAppointmentByID :one
SELECT * FROM appointments
WHERE id = $1;

-- name: ListAppointmentsByVisitor :many
SELECT 
  a.*, 
  u.first_name || ' ' || u.last_name AS host_name,
  u.role AS role
FROM appointments a
JOIN users u ON a.host_id = u.id
WHERE a.visitor_id = $1
ORDER BY a.appointment_date DESC;

-- name: ListAppointmentsByHost :many
SELECT 
  a.*, 
  u.first_name || ' ' || u.last_name AS visitor_name,
  u.role AS role
FROM appointments a
JOIN users u ON a.visitor_id = u.id
WHERE a.host_id = $1
ORDER BY a.appointment_date DESC;

-- name: ListAppointmentsByDate :many
SELECT 
    a.*,
    host.first_name || ' ' || host.last_name AS host_name,
    visitor.first_name || ' ' || visitor.last_name AS visitor_name
FROM appointments a
JOIN users host ON a.host_id = host.id
JOIN users visitor ON a.visitor_id = visitor.id
WHERE a.appointment_date = $1
ORDER BY a.start_time;

-- name: GetAppointmentByQRCode :one
SELECT 
  a.*,
  host.first_name || ' ' || host.last_name AS host_name,
  visitor.first_name || ' ' || visitor.last_name AS visitor_name
FROM appointments a
JOIN users host ON a.host_id = host.id
JOIN users visitor ON a.visitor_id = visitor.id
WHERE a.qr_code = $1
LIMIT 1;

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
