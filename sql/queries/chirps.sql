-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2

)
RETURNING *;

-- name: GetChirps :many
select * from chirps ORDER BY created_at ASC;

-- name: GetChirpByID :one
select * from chirps WHERE id = $1 ORDER BY created_at ASC;

-- name: GetChirpsByUserID :many
select * from chirps WHERE user_id = $1 ORDER BY created_at ASC;



-- name: DeleteChirpByID :exec
DELETE FROM chirps where id = $1;

-- name: DropAllChirps :exec
DELETE FROM chirps;