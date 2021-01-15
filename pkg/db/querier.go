// Code generated by sqlc. DO NOT EDIT.

package db

import (
	"context"

	"github.com/gofrs/uuid"
)

type Querier interface {
	CreateAccount(ctx context.Context, name string) (Account, error)
	CreatePosting(ctx context.Context, arg CreatePostingParams) (Posting, error)
	CreateTransaction(ctx context.Context, arg CreateTransactionParams) (Transaction, error)
	GetAccounts(ctx context.Context, name string) ([]Account, error)
	GetPostings(ctx context.Context, transaction uuid.UUID) ([]Posting, error)
	GetTransactions(ctx context.Context, arg GetTransactionsParams) ([]Transaction, error)
}

var _ Querier = (*Queries)(nil)