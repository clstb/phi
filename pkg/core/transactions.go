package core

import (
	"context"
	"fmt"
	"io"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/clstb/phi/pkg/core/db"
	"github.com/clstb/phi/pkg/fin"
	"github.com/clstb/phi/pkg/pb"
	"github.com/gofrs/uuid"
)

func (s *Server) CreateTransactions(
	stream pb.Core_CreateTransactionsServer,
) error {
	ctx := stream.Context()

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

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

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
			postingStmt = postingStmt.Values(
				posting.Account.String(),
				id,
				posting.Units.String(),
				posting.Cost.String(),
				posting.Price.String(),
			)
		}
	}

	_, err = transactionStmt.PlaceholderFormat(sq.Dollar).RunWith(tx).ExecContext(ctx)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = postingStmt.PlaceholderFormat(sq.Dollar).RunWith(tx).ExecContext(ctx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
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
