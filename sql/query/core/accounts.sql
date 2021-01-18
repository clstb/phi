-- name: CreateAccount :one
INSERT INTO accounts (
  name
) VALUES (
  $1
) RETURNING *;
-- name: GetAccounts :many
SELECT
  *
FROM
  accounts
WHERE
  name ~ $1;
-- name: DeleteAccount :exec
DELETE FROM
  accounts
WHERE id = $1;
