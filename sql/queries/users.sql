-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email)
VALUES (
    gen_random_uuid (),
    NOW(),
    NOW(),
    $1
)
RETURNING *;

-- name: AddPassword :one
UPDATE users
SET hashed_password = $2
WHERE id = $1
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;