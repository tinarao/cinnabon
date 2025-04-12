-- name: CreateUser :one
INSERT INTO users (email, username, first_name, last_name, password)
VALUES (?, ?, ?, ?, ?)
RETURNING id;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = ?;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = ?;


