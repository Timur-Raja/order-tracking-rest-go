-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;
-- name: GetUsers :many
SELECT * FROM users;