-- name: GetProductByID :one
SELECT * FROM products
WHERE id = $1;
-- name: GetProducts :many
SELECT * FROM products;