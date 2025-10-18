-- name: CreateUser :one
INSERT INTO users (id, email)
VALUES ($1, $2)
RETURNING id;

-- name: GetUserByEmail :one
SELECT id FROM users
WHERE email = $1 AND deleted_at IS NULL;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 AND deleted_at IS NULL;

-- name: DeleteUser :exec
UPDATE users
SET deleted_at = NOW()
WHERE id = $1;

-- name: GetSoftDeletedUserByEmail :one
SELECT id FROM users
WHERE email = $1 AND deleted_at IS NOT NULL;

-- name: UndeleteUser :one
UPDATE users
SET deleted_at = NULL
WHERE id = $1
RETURNING id;