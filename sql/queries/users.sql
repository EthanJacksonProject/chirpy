-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: GetUserByEmail :one
select * from users WHERE email = $1;

-- name: SetPasswordHashByID :one
UPDATE users
SET hashed_password = $2, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: SetEmailByID :one
UPDATE users
SET email = $2, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: SetEmailAndPWHashById :one
UPDATE users
SET email = $2, hashed_password = $3, updated_at = NOW()
WHERE id = $1
RETURNING *;
-- name: DropAllUsers :exec
DELETE FROM users;