-- name: CreateUser :one
INSERT INTO users (
  name,
  password
) VALUES (
  $1,
  $2
) RETURNING *;
-- name: GetUserByName :one
SELECT
  *
FROM
  users
WHERE
  name = $1;
