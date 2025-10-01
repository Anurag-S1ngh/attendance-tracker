-- name: CreateUser :exec
INSERT INTO users (id, email, password_hash) VALUES ($1, $2, $3);

-- name: GetUserByEmail :one
SELECT id, password_hash FROM users WHERE email = $1 AND deleted_at is NULL LIMIT 1;

-- name: DeleteUser :exec
UPDATE users SET deleted_at = NOW() WHERE id = $1 AND deleted_at is NULL;
