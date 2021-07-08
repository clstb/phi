-- name: CreateTransaction :one
INSERT INTO transactions (
  date,
  entity,
  reference,
  "user",
  "from",
  "to",
  units,
  unitsCur,
  cost,
  costCur,
  price,
  priceCur,
  tink_id,
  debit
) VALUES (
  $1,
  $2,
  $3,
  $4,
  $5,
  $6,
  $7,
  $8,
  $9,
  $10,
  $11,
  $12,
  $13,
  $14
) ON CONFLICT (
  tink_id
) DO NOTHING RETURNING *;
-- name: UpdateTransaction :one
UPDATE
  transactions
SET
  reference = $1,
  "from" = $2,
  "to" = $3
WHERE
  id = $4
AND
  "user" = $5
RETURNING *;
-- name: GetTransactions :many
SELECT
  transactions.*
FROM
  transactions
JOIN
  accounts accounts_from
ON
  transactions.from = accounts_from.id
JOIN
  accounts accounts_to
ON
  transactions.to = accounts_to.id
WHERE
  date BETWEEN @from_date AND @to_date
AND
  transactions.user = @user_id
AND
  accounts_from.name ~ @account_name
OR
  accounts_to.name ~ @account_name
ORDER BY
  date
DESC;
-- name: DeleteTransaction :exec
DELETE FROM
  transactions
WHERE id = $1;
