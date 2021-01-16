// Code generated by sqlc. DO NOT EDIT.
// source: accounts.sql

package db

import (
	"context"

	"github.com/gofrs/uuid"
)

const createAccount = `-- name: CreateAccount :one
INSERT INTO accounts (
  name
) VALUES (
  $1
) RETURNING id, name
`

func (q *Queries) CreateAccount(ctx context.Context, name string) (Account, error) {
	row := q.queryRow(ctx, q.createAccountStmt, createAccount, name)
	var i Account
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const deleteAccount = `-- name: DeleteAccount :exec
DELETE FROM
  accounts
WHERE id = $1
`

func (q *Queries) DeleteAccount(ctx context.Context, id uuid.UUID) error {
	_, err := q.exec(ctx, q.deleteAccountStmt, deleteAccount, id)
	return err
}

const getAccounts = `-- name: GetAccounts :many
SELECT
  id, name
FROM
  accounts
WHERE
  name ~ $1
`

func (q *Queries) GetAccounts(ctx context.Context, name string) ([]Account, error) {
	rows, err := q.query(ctx, q.getAccountsStmt, getAccounts, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Account
	for rows.Next() {
		var i Account
		if err := rows.Scan(&i.ID, &i.Name); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
