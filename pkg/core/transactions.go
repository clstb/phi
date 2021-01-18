package core

import (
	"context"
	"fmt"
	"time"

	"github.com/clstb/phi/pkg/core/db"
	"github.com/clstb/phi/pkg/fin"
	"github.com/clstb/phi/pkg/pb"
)

func (s *Server) CreateTransaction(
	ctx context.Context,
	req *pb.Transaction,
) (*pb.Transaction, error) {
	transaction, err := fin.TransactionFromPB(req)
	if err != nil {
		return nil, err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	q := db.New(tx)

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
	fmt.Println(from)
	fmt.Println(to)

	q := db.New(s.db)
	transactionsDB, err := q.GetTransactions(ctx, db.GetTransactionsParams{
		AccountName: req.AccountName,
		FromDate:    from,
		ToDate:      to,
	})
	fmt.Println(transactionsDB)

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
