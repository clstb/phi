package server

import (
	"context"
	"database/sql"

	"github.com/clstb/phi/pkg/db"
	"github.com/clstb/phi/pkg/fin"
	"github.com/clstb/phi/pkg/pb"
	"github.com/golang/protobuf/ptypes"
)

type core struct {
	pb.UnimplementedCoreServer
	db *sql.DB
}

func (s *core) CreateAccount(
	ctx context.Context,
	req *pb.Account,
) (*pb.Account, error) {
	account := fin.AccountFromPB(req)

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	q := db.New(tx)

	accountDB, err := q.CreateAccount(ctx, account.Name)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	account = fin.NewAccount(accountDB)

	return account.PB(), nil
}

func (s *core) GetAccounts(
	ctx context.Context,
	req *pb.AccountsQuery,
) (*pb.Accounts, error) {
	q := db.New(s.db)
	accountsDB, err := q.GetAccounts(ctx, req.Name)
	if err != nil {
		return nil, err
	}

	accounts := fin.NewAccounts(accountsDB...)

	return accounts.PB(), nil
}

func (s *core) CreateTransaction(
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
	res, err := transaction.PB()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *core) GetTransactions(
	ctx context.Context,
	req *pb.TransactionsQuery,
) (*pb.Transactions, error) {
	from, err := ptypes.Timestamp(req.From)
	if err != nil {
		return nil, err
	}
	to, err := ptypes.Timestamp(req.To)
	if err != nil {
		return nil, err
	}

	q := db.New(s.db)
	transactionsDB, err := q.GetTransactions(ctx, db.GetTransactionsParams{
		AccountName: req.AccountName,
		FromDate:    from,
		ToDate:      to,
	})

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

	res, err := transactions.PB()
	if err != nil {
		return nil, err
	}

	return res, nil
}
