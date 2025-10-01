-- name: CreateEvent :one
INSERT INTO events (id, user_id, name) VALUES ($1, $2, $3) RETURNING *;

-- name: GetEvents :many
SELECT id, name FROM events WHERE user_id = $1;

-- name: DeleteEvent :exec
DELETE FROM events e
WHERE e.id = $1
  AND EXISTS (
    SELECT 1 FROM users u
    WHERE u.id = e.user_id AND u.id = $2
  );

-- name: GetUserEventsWithAttendanceAndCounts :many
SELECT
    e.id AS event_id,
    e.name AS event_name,
    COALESCE(
        json_agg(
            json_build_object(
                'attendance_id', a.id,
                'attendance_status', a.attended,
                'attendance_date', a.date
            )
        ) FILTER (WHERE a.id IS NOT NULL),
        '[]'
    )::json AS attendance,
    COUNT(*) FILTER (WHERE a.attended = 'present') AS present_count,
    COUNT(*) FILTER (WHERE a.attended = 'absent') AS absent_count,
    COUNT(*) FILTER (WHERE a.attended = 'canceled') AS canceled_count
FROM events e
LEFT JOIN attendance a ON e.id = a.event_id
LEFT JOIN users u ON u.id = a.user_id
WHERE e.user_id = $1
GROUP BY e.id, e.name
ORDER BY e.id;

