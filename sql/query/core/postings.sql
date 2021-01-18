-- name: CreatePosting :one
INSERT INTO postings (
  account,
  transaction,
  units,
  cost,
  price
) VALUES (
  $1,
  $2,
  $3,
  $4,
  $5
) RETURNING *;
-- name: GetPostings :many
SELECT
  *
FROM
  postings
WHERE
  transaction = $1;
-- name: DeletePosting :exec
DELETE FROM
  postings
WHERE id = $1;
