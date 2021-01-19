-- name: CreateAccount :one
INSERT INTO accounts (
  name
) VALUES (
  $1
) RETURNING *;
-- name: GetAccounts :many
SELECT
  accounts.id,
  accounts.name
FROM
  accounts
JOIN
  accounts_users
ON
  accounts_users.account = accounts.id
AND
  accounts_users.user = $1
WHERE
  accounts.name ~ $2;
-- name: DeleteAccount :exec
DELETE FROM
  accounts
WHERE id = $1;
-- name: LinkAccount :one
INSERT INTO accounts_users (
  account,
  "user"
) VALUES (
  $1,
  $2
) RETURNING *;
-- name: OwnsAccount :one
SELECT
  COUNT(1)
FROM
  accounts_users
WHERE
  account = $1
AND
  "user" = $2;
