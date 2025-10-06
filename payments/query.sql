-- name: CreatePayment :one
INSERT INTO payments (
  user_id, status, method, amount
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;