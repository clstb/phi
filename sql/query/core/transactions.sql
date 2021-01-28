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
SELECT
  transactions.id,
  transactions.date,
  transactions.entity,
  transactions.reference,
  transactions.hash,
  postings.id AS posting_id,
  postings.account,
  postings.units,
  postings.cost,
  postings.price,
  accounts.name
FROM
  transactions
INNER JOIN
  postings
ON
  transactions.id = postings.transaction
INNER JOIN
  accounts
ON
  accounts.id = postings.account
AND
  accounts.name ~ @account_name
INNER JOIN
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
