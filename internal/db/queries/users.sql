-- name: CreateUser :one
INSERT INTO users (id, email) VALUES ($1, $2) RETURNING id;

-- name: GetUserByEmail :one
SELECT id FROM users WHERE email = $1 LIMIT 1;

-- name: DeleteUser :exec
UPDATE users SET deleted_at = NOW() WHERE id = $1 AND deleted_at is NULL;
