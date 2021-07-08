-- name: CreateToken :one
INSERT INTO tokens (
  access_token,
  refresh_token,
  token_type,
  expires_at,
  scope,
  "user"
) VALUES (
  $1,
  $2,
  $3,
  $4,
  $5,
  $6
) RETURNING *;
-- name: GetTokens :many
SELECT
  *
FROM
  tokens
WHERE
  "user" = $1;
-- name: DeleteToken :exec
DELETE FROM
  tokens
WHERE id = $1;
