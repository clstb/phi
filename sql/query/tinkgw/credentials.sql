-- name: CreateCredential :one
INSERT INTO credentials (
  "user",
  credentials_id
) VALUES (
  $1,
  $2
) RETURNING *;
-- name: GetCredentials :many
SELECT
  *
FROM
  credentials
WHERE
  "user" = $1;
-- name: DeleteCredential :exec
DELETE FROM
  credentials
WHERE id = $1;
