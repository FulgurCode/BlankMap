-- name: CreateUser :one
INSERT INTO users (email, name, password)
VALUES ($1, $2, $3)
RETURNING id, email, name, created_at, updated_at;

-- name: GetUserByID :one
SELECT id, email, name, created_at, updated_at FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT id, email, name, password, created_at, updated_at FROM users
WHERE email = $1;
