-- name: CreateUser :one
INSERT INTO users (
  id,
  tink_id
) VALUES (
  $1,
  $2
) ON CONFLICT DO NOTHING RETURNING id;
-- name: GetUserByID :one
SELECT
  *
FROM
  users
WHERE
  id = @id::uuid;
