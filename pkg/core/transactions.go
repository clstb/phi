package core

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"time"

	sq "github.com/Masterminds/squirrel"
	db "github.com/clstb/phi/pkg/db/core"
	"github.com/clstb/phi/pkg/fin"
	"github.com/clstb/phi/pkg/pb"
	"github.com/gofrs/uuid"
)

func (s *Server) CreateTransactions(
	stream pb.Core_CreateTransactionsServer,
) error {
	ctx := stream.Context()

	subStr, ok := ctx.Value("sub").(string)
	if !ok {
		return fmt.Errorf("context: missing subject")
	}
	sub, err := uuid.FromString(subStr)
	if err != nil {
		return err
	}

	var transactions fin.Transactions
	for {
		transactionPB, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		transaction, err := fin.TransactionFromPB(transactionPB)
		if err != nil {
			return err
		}
		transactions = append(transactions, transaction)
	}

	tx, ok := ctx.Value("tx").(*sql.Tx)
	if !ok {
		return fmt.Errorf("context: missing transaction")
	}
	q := db.New(tx)

	accountsDB, err := q.GetAccounts(ctx, db.GetAccountsParams{
		User: sub,
	})
	if err != nil {
		return err
	}
	accounts := fin.AccountsFromDB(accountsDB...)

	transactionStmt := sq.Insert("transactions").Columns(
		"id",
		"date",
		"entity",
		"reference",
		"hash",
	)
	postingStmt := sq.Insert("postings").Columns(
		"account",
		"transaction",
		"units",
		"cost",
		"price",
	)

	for _, transaction := range transactions {
		id, err := uuid.NewV4()
		if err != nil {
			return err
		}
		transactionStmt = transactionStmt.Values(
			id,
			transaction.Date,
			transaction.Entity,
			transaction.Reference,
			transaction.Hash,
		)

		for _, posting := range transaction.Postings {
			if accounts.ById(posting.Account.String()).Empty() {
				return fmt.Errorf("unauthorized: post to account")
			}

			postingStmt = postingStmt.Values(
				posting.Account.String(),
				id,
				posting.Units.String(),
				posting.Cost.String(),
				posting.Price.String(),
			)
		}
	}

	_, err = transactionStmt.
		PlaceholderFormat(sq.Dollar).
		RunWith(tx).
		ExecContext(ctx)
	if err != nil {
		return err
	}

	_, err = postingStmt.
		PlaceholderFormat(sq.Dollar).
		RunWith(tx).
		ExecContext(ctx)
	if err != nil {
		return err
	}

	return stream.SendAndClose(transactions.PB())
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

	tx, ok := ctx.Value("tx").(*sql.Tx)
	if !ok {
		return nil, fmt.Errorf("context: missing transaction")
	}
	q := db.New(tx)

	rows, err := q.GetTransactions(ctx, db.GetTransactionsParams{
		UserID:      sub,
		AccountName: req.AccountName,
		FromDate:    from,
		ToDate:      to,
	})
	if err != nil {
		return nil, err
	}

	transactions := map[string]*pb.Transaction{}
	for _, row := range rows {
		posting := &pb.Posting{
			Id:          row.PostingID.String(),
			Account:     row.Account.String(),
			Transaction: row.ID.String(),
			Units:       row.UnitsStr,
			Cost:        row.CostStr,
			Price:       row.PriceStr,
		}

		transaction, ok := transactions[row.ID.String()]
		if !ok {
			transaction = &pb.Transaction{
				Id:        row.ID.String(),
				Date:      row.Date.Format("2006-01-02"),
				Entity:    row.Entity,
				Reference: row.Reference,
				Hash:      row.Hash,
			}
		}
		transaction.Postings = append(transaction.Postings, posting)

		transactions[row.ID.String()] = transaction
	}

	var data []*pb.Transaction
	for _, transaction := range transactions {
		data = append(data, transaction)
	}

	return &pb.Transactions{
		Data: data,
	}, nil
}
