-- name: CreatePayment :one
INSERT INTO payments (
  user_id, status, method, amount
) VALUES (
  ?, ?, ?, ?
)
RETURNING *;