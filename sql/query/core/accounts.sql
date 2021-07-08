-- name: CreateAccount :one
INSERT INTO accounts (
  name,
  "user"
) VALUES (
  $1,
  $2
) RETURNING *;
-- name: GetAccounts :many
SELECT
  accounts.id,
  accounts.name,
  accounts.user
FROM
  accounts
WHERE
  accounts.user = $1
AND
  accounts.name ~ $2;
-- name: DeleteAccount :exec
DELETE FROM
  accounts
WHERE id = $1;
