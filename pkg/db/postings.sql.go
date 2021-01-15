// Code generated by sqlc. DO NOT EDIT.
// source: postings.sql

package db

import (
	"context"

	"github.com/gofrs/uuid"
)

const createPosting = `-- name: CreatePosting :one
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
) RETURNING id, account, transaction, units, cost, price
`

type CreatePostingParams struct {
	Account     uuid.UUID `json:"account"`
	Transaction uuid.UUID `json:"transaction"`
	Units       Amount    `json:"units"`
	Cost        Amount    `json:"cost"`
	Price       Amount    `json:"price"`
}

func (q *Queries) CreatePosting(ctx context.Context, arg CreatePostingParams) (Posting, error) {
	row := q.queryRow(ctx, q.createPostingStmt, createPosting,
		arg.Account,
		arg.Transaction,
		arg.Units,
		arg.Cost,
		arg.Price,
	)
	var i Posting
	err := row.Scan(
		&i.ID,
		&i.Account,
		&i.Transaction,
		&i.Units,
		&i.Cost,
		&i.Price,
	)
	return i, err
}

const getPostings = `-- name: GetPostings :many
SELECT
  id, account, transaction, units, cost, price
FROM
  postings
WHERE
  transaction = $1
`

func (q *Queries) GetPostings(ctx context.Context, transaction uuid.UUID) ([]Posting, error) {
	rows, err := q.query(ctx, q.getPostingsStmt, getPostings, transaction)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Posting
	for rows.Next() {
		var i Posting
		if err := rows.Scan(
			&i.ID,
			&i.Account,
			&i.Transaction,
			&i.Units,
			&i.Cost,
			&i.Price,
		); err != nil {
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