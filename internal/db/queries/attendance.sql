-- name: MarkAttendance :one
INSERT INTO attendance (id, event_id, attended, date, user_id)
SELECT $1, e.id, $3, $4, $5
FROM events e
WHERE e.id = $2
  AND e.user_id = $5
ON CONFLICT (event_id, date)
DO UPDATE SET attended = EXCLUDED.attended,
              updated_at = NOW()
RETURNING *;

-- name: GetAttendanceByUserAndDate :one
SELECT id, event_id, attended, date FROM attendance
WHERE user_id = $1
  AND date = $2 
  AND event_id = $3
ORDER BY date DESC; 

-- name: GetAttendance :many
SELECT id, event_id, attended, date FROM attendance WHERE user_id = $1;

-- name: UpdateAttendance :one
UPDATE attendance SET attended = $1, updated_at = $2 WHERE event_id = $3 AND user_id = $4 RETURNING *;

-- name: DeleteAttendance :exec
DELETE FROM attendance WHERE id = $1 AND user_id = $2;

