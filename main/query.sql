-- name: GetUser :one
SELECT * FROM users
WHERE username = ? LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (
  username, hashed_password
) VALUES (
  ?, ?
)
RETURNING *;

-- name: UpdateUserTokensLogin :exec
UPDATE users
set session_token = ?,
csrf_token = ?
WHERE username = ?;

-- name: UpdateUserTokensLogout :exec
UPDATE users
set session_token = "",
csrf_token = ""
WHERE username = ?;

-- name: DeleteUser :exec
DELETE FROM users
WHERE username = ?;