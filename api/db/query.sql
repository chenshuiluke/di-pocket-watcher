-- name: GetUser :one
SELECT id, email FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT id, email FROM users
WHERE email = $1 LIMIT 1;

-- name: ListUsers :many
SELECT id, email FROM users
ORDER BY id;

-- name: CreateUser :one
INSERT INTO users (
  email, token
) VALUES (
  $1, $2
)
RETURNING *;

-- name: UpdateUser :exec
UPDATE users
  set email = $2
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;