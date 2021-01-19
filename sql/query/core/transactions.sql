-- name: CreateTransaction :one
INSERT INTO transactions (
  date,
  entity,
  reference,
  hash
) Values (
  $1,
  $2,
  $3,
  $4
) RETURNING *;
-- name: GetTransactions :many
SELECT DISTINCT
  transactions.id,
  date,
  entity,
  reference,
  hash
FROM
  transactions
JOIN
  postings
ON
  transactions.id = postings.transaction
JOIN
  accounts
ON
  accounts.id = postings.account
AND
  accounts.name ~ @account_name
JOIN
  accounts_users
ON
  accounts_users.account = accounts.id
AND
  accounts_users.user = @user_id
WHERE
  date BETWEEN @from_date AND @to_date;
-- name: DeleteTransaction :exec
DELETE FROM
  transactions
WHERE id = $1;
