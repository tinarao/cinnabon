-- name: CreateUser :one
INSERT INTO users (email, username, first_name, last_name, password, foreign_id)
VALUES (?, ?, ?, ?, ?, ?)
RETURNING id;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = ?;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = ?;

-- name: CreateSession :one
INSERT INTO sessions (user_id, hash, expires_at)
VALUES (?, ?, ?)
RETURNING hash;

-- name: DeleteSessionById :exec
DELETE FROM sessions WHERE id = ?;

-- name: DeleteSessionByHash :exec
DELETE FROM sessions WHERE hash = ?;

-- name: GetSessionByID :one
SELECT * FROM sessions WHERE id = ?;

-- name: GetSessionByHash :one
SELECT * FROM sessions WHERE hash = ?;

-- name: GetSessionByUserID :one
SELECT * FROM sessions WHERE user_id = ?;


