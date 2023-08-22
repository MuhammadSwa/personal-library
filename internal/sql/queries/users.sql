-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: CreateUser :one
INSERT INTO users (email, hashed_password, username)
VALUES ($1, $2, $3)
RETURNING id;
