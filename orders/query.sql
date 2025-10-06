-- name: CreateOrder :one
INSERT INTO orders (
  user_id, user_name, email, phone, shipping_address, status, order_id
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;