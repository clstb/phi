package core

import (
	"context"
	"errors"
	"fmt"
	"time"

	db "github.com/clstb/phi/pkg/db/core"
	"github.com/clstb/phi/pkg/fin"
	"github.com/clstb/phi/pkg/pb"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
)

func (s *Server) CreateTransactions(
	ctx context.Context,
	req *pb.Transactions,
) (*pb.Transactions, error) {
	subStr, ok := ctx.Value("sub").(string)
	if !ok {
		return nil, fmt.Errorf("context: missing subject")
	}
	sub, err := uuid.FromString(subStr)
	if err != nil {
		return nil, fmt.Errorf("context: invalid subject id: %w", err)
	}

	transactions, err := fin.TransactionsFromPB(req)
	if err != nil {
		return nil, fmt.Errorf("parsing transactions: %w", err)
	}

	tx, ok := ctx.Value("tx").(pgx.Tx)
	if !ok {
		return nil, fmt.Errorf("context: missing database transaction")
	}
	q := db.New(tx)

	accountsDB, err := q.GetAccounts(ctx, db.GetAccountsParams{
		User: sub,
	})
	if err != nil {
		return nil, err
	}
	accounts := fin.AccountsFromDB(accountsDB...).ById()

	var res fin.Transactions
	for _, transaction := range transactions {
		_, ok := accounts[transaction.From.String()]
		if !ok {
			return nil, fmt.Errorf("invalid request: account not found: %s", transaction.From)
		}
		_, ok = accounts[transaction.To.String()]
		if !ok {
			return nil, fmt.Errorf("invalid request: account not found: %s", transaction.From)
		}

		t, err := q.CreateTransaction(ctx, db.CreateTransactionParams{
			Date:      transaction.Date,
			Entity:    transaction.Entity,
			Reference: transaction.Reference,
			User:      sub,
			From:      transaction.From,
			To:        transaction.To,
			Units:     transaction.Units.Decimal.Abs(),
			Unitscur:  transaction.Units.Currency,
			Cost:      transaction.Cost.Decimal.Abs(),
			Costcur:   transaction.Cost.Currency,
			Price:     transaction.Price.Decimal.Abs(),
			Pricecur:  transaction.Price.Currency,
			TinkID:    transaction.TinkID,
			Debit:     transaction.Debit,
		})
		if errors.Is(err, pgx.ErrNoRows) {
			continue
		}
		if err != nil {
			return nil, err
		}
		res = append(res, fin.TransactionFromDB(t))
	}

	return res.PB(), nil
}

func (s *Server) UpdateTransactions(
	ctx context.Context,
	req *pb.Transactions,
) (*pb.Transactions, error) {
	subStr, ok := ctx.Value("sub").(string)
	if !ok {
		return nil, fmt.Errorf("context: missing subject")
	}
	sub, err := uuid.FromString(subStr)
	if err != nil {
		return nil, fmt.Errorf("context: invalid subject id: %w", err)
	}

	transactions, err := fin.TransactionsFromPB(req)
	if err != nil {
		return nil, fmt.Errorf("parsing transactions: %w", err)
	}

	tx, ok := ctx.Value("tx").(pgx.Tx)
	if !ok {
		return nil, fmt.Errorf("context: missing database transaction")
	}
	q := db.New(tx)

	accountsDB, err := q.GetAccounts(ctx, db.GetAccountsParams{
		User: sub,
	})
	if err != nil {
		return nil, err
	}
	accounts := fin.AccountsFromDB(accountsDB...).ById()

	var res fin.Transactions
	for _, transaction := range transactions {
		_, ok := accounts[transaction.From.String()]
		if !ok {
			return nil, fmt.Errorf("invalid request: account not found: %s", transaction.From)
		}
		_, ok = accounts[transaction.To.String()]
		if !ok {
			return nil, fmt.Errorf("invalid request: account not found: %s", transaction.From)
		}

		t, err := q.UpdateTransaction(ctx, db.UpdateTransactionParams{
			ID:        transaction.ID,
			User:      sub,
			Reference: transaction.Reference,
			From:      transaction.From,
			To:        transaction.To,
		})
		if err != nil {
			return nil, err
		}
		res = append(res, fin.TransactionFromDB(t))
	}

	return res.PB(), nil
}

func (s *Server) GetTransactions(
	ctx context.Context,
	req *pb.TransactionsQuery,
) (*pb.Transactions, error) {
	subStr, ok := ctx.Value("sub").(string)
	if !ok {
		return nil, fmt.Errorf("context: missing subject")
	}
	sub, err := uuid.FromString(subStr)
	if err != nil {
		return nil, err
	}

	fromDate, err := time.Parse("2006-01-02", req.From)
	if err != nil {
		if req.From != "" {
			return nil, err
		}
		fromDate = time.Unix(0, 0)
	}
	toDate, err := time.Parse("2006-01-02", req.To)
	if err != nil {
		if req.To != "" {
			return nil, err
		}
		toDate = time.Now()
	}

	tx, ok := ctx.Value("tx").(pgx.Tx)
	if !ok {
		return nil, fmt.Errorf("context: missing transaction")
	}
	q := db.New(tx)

	transactionsDB, err := q.GetTransactions(ctx, db.GetTransactionsParams{
		UserID:      sub,
		FromDate:    fromDate,
		ToDate:      toDate,
		AccountName: req.AccountName,
	})
	if err != nil {
		return nil, err
	}

	return fin.TransactionsFromDB(transactionsDB).PB(), nil
}
