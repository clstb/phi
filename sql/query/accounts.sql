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
