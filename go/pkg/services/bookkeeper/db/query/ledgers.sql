-- name: CreateLedger :one
INSERT INTO ledgers (
  user_id,
  data,
  dk
) VALUES (
  $1,
  $2,
  $3
) RETURNING *;
-- name: UpdateLedger :one
UPDATE ledgers SET
  data = $1,
  dk = $2
WHERE
  id = $3
RETURNING *;
-- name: GetLedger :one
SELECT
  *
FROM
  ledgers
WHERE
  id = $1;