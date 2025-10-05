-- name: CreateOrder :one
INSERT INTO orders (
  user_id, user_name, email, phone, shipping_address, status
) VALUES (
  ?, ?, ?, ?, ?, ?
)
RETURNING *;