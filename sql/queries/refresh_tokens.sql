

-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, expires_at, revoked_at, user_id)
VALUES (
    $1,
    NOW(),
    NOW(),
    NOW() + INTERVAl '60' DAY,
    NULL,
    $2

)
RETURNING *;

-- name: GetTokenDataByToken :one
select * from refresh_tokens where token = $1;

-- name: RevokeTokenByToken :exec
UPDATE refresh_tokens
SET revoked_at = NOW()
WHERE token = $1;

-- name: DropAllRefreshTokens :exec
DELETE FROM refresh_tokens;