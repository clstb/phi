package core

import (
	"context"
	"fmt"
	"time"

	"github.com/clstb/phi/pkg/core/db"
	"github.com/clstb/phi/pkg/fin"
	"github.com/clstb/phi/pkg/pb"
	"github.com/gofrs/uuid"
)

func (s *Server) CreateTransaction(
	ctx context.Context,
	req *pb.Transaction,
) (*pb.Transaction, error) {
	subStr, ok := ctx.Value("sub").(string)
	if !ok {
		return nil, fmt.Errorf("context: missing subject")
	}
	sub := uuid.FromStringOrNil(subStr)

	transaction, err := fin.TransactionFromPB(req)
	if err != nil {
		return nil, err
	}

	q := db.New(s.db)

	for _, posting := range transaction.Postings {
		exists, err := q.OwnsAccount(ctx, db.OwnsAccountParams{
			Account: posting.Account,
			User:    sub,
		})
		if err != nil {
			return nil, err
		}
		if exists == 0 {
			return nil, fmt.Errorf("unauthorized: account")
		}
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	q = q.WithTx(tx)

	transactionDB, err := q.CreateTransaction(ctx, db.CreateTransactionParams{
		Date:      transaction.Date,
		Entity:    transaction.Entity,
		Reference: transaction.Reference,
		Hash:      transaction.Hash,
	})
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	var postings fin.Postings
	for _, posting := range transaction.Postings {
		postingDB, err := q.CreatePosting(ctx, db.CreatePostingParams{
			Account:     posting.Account,
			Transaction: transactionDB.ID,
			Units:       posting.Units,
			Cost:        posting.Cost,
			Price:       posting.Price,
		})
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		postings = append(postings, fin.NewPosting(postingDB))
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	transaction = fin.NewTransaction(transactionDB, postings)

	return transaction.PB(), nil
}

func (s *Server) GetTransactions(
	ctx context.Context,
	req *pb.TransactionsQuery,
) (*pb.Transactions, error) {
	subStr, ok := ctx.Value("sub").(string)
	if !ok {
		return nil, fmt.Errorf("context: missing subject")
	}
	sub := uuid.FromStringOrNil(subStr)

	from, err := time.Parse("2006-01-02", req.From)
	if err != nil {
		if req.From != "" {
			return nil, err
		}
		from = time.Unix(0, 0)
	}
	to, err := time.Parse("2006-01-02", req.To)
	if err != nil {
		if req.To != "" {
			return nil, err
		}
		to = time.Now()
	}

	q := db.New(s.db)
	transactionsDB, err := q.GetTransactions(ctx, db.GetTransactionsParams{
		UserID:      sub,
		AccountName: req.AccountName,
		FromDate:    from,
		ToDate:      to,
	})
	if err != nil {
		return nil, err
	}

	var transactions fin.Transactions
	for _, transaction := range transactionsDB {
		postingsDB, err := q.GetPostings(ctx, transaction.ID)
		if err != nil {
			return nil, err
		}
		postings := fin.NewPostings(postingsDB...)

		transactions = append(
			transactions,
			fin.NewTransaction(transaction, postings),
		)
	}

	return transactions.PB(), nil
}
