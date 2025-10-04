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
WHERE session_token = ? AND csrf_token = ?;

-- name: DeleteUser :exec
DELETE FROM users
WHERE username = ?;

-- name: GetUserBySessionToken :one
SELECT * FROM users 
WHERE session_token = ? 
LIMIT 1;