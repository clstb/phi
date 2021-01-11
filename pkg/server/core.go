package server

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/clstb/phi/pkg/pb"
	"github.com/doug-martin/goqu/v9"
)

type core struct {
	pb.UnimplementedCoreServer
	db *goqu.Database
}

func (s *core) CreateAccount(
	ctx context.Context,
	req *pb.Account,
) (*pb.Account, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}

	if err := tx.Wrap(func() error {
		rows, err := tx.Insert(
			"accounts",
		).Returning(goqu.C("id")).Rows(
			goqu.Record{"name": req.Name},
		).Executor().Query()
		if err != nil {
			return err
		}

		rows.Next()
		return rows.Scan(&req.Id)
	}); err != nil {
		return nil, err
	}

	return req, nil
}

func (s *core) GetAccounts(
	ctx context.Context,
	req *pb.AccountsQuery,
) (*pb.Accounts, error) {
	fields := func() []interface{} {
		fields := []interface{}{"id"}
		if req.Fields.Name {
			fields = append(fields, "name")
		}
		return fields
	}
	scan := func(rows *sql.Rows) (*pb.Account, error) {
		account := &pb.Account{}
		toScan := []interface{}{&account.Id}

		if req.Fields.Name {
			toScan = append(toScan, &account.Name)
		}

		if err := rows.Scan(toScan...); err != nil {
			return nil, err
		}

		return account, nil
	}

	sql := s.db.From(
		"accounts",
	).Select(
		fields()...,
	).Where(
		goqu.C("name").RegexpLike(req.Name),
	)

	rows, err := sql.Executor().Query()
	if err != nil {
		return nil, err
	}

	var data []*pb.Account
	byId := make(map[string]int32)
	byName := make(map[string]int32)

	var i int32
	for rows.Next() {
		account, err := scan(rows)
		if err != nil {
			return nil, err
		}
		data = append(data, account)
		byId[account.Id] = i
		byName[account.Name] = i
		i++
	}

	return &pb.Accounts{
		Data:   data,
		ById:   byId,
		ByName: byName,
	}, nil
}

func (s *core) CreateTransaction(
	ctx context.Context,
	req *pb.Transaction,
) (*pb.Transaction, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}

	if err := tx.Wrap(func() error {
		rows, err := tx.Insert(
			"transactions",
		).Returning(goqu.C("id")).Rows(
			goqu.Record{
				"date":      req.Date,
				"entity":    req.Entity,
				"reference": req.Reference,
			},
		).Executor().Query()
		if err != nil {
			return err
		}

		rows.Next()
		if err := rows.Scan(&req.Id); err != nil {
			return err
		}

		var postings []goqu.Record
		for _, posting := range req.Postings.Data {
			record := goqu.Record{
				"account":     posting.Account,
				"transaction": req.Id,
				"units":       posting.Units,
				"cost":        posting.Cost,
				"price":       posting.Price,
			}
			postings = append(postings, record)
		}

		_, err = tx.Insert(
			"postings",
		).Rows(
			postings,
		).Executor().Exec()
		fmt.Println(err)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return req, nil
}

func (s *core) GetTransactions(
	ctx context.Context,
	req *pb.TransactionsQuery,
) (*pb.Transactions, error) {
	fields := func() []interface{} {
		fields := []interface{}{"transactions.id"}
		if req.Fields.Date {
			fields = append(fields, "transactions.date")
		}
		if req.Fields.Entity {
			fields = append(fields, "transactions.entity")
		}
		if req.Fields.Reference {
			fields = append(fields, "transactions.reference")
		}
		return fields
	}
	scan := func(rows *sql.Rows) (*pb.Transaction, error) {
		transaction := &pb.Transaction{}
		toScan := []interface{}{&transaction.Id}

		if req.Fields.Date {
			toScan = append(toScan, &transaction.Date)
		}
		if req.Fields.Entity {
			toScan = append(toScan, &transaction.Entity)
		}
		if req.Fields.Reference {
			toScan = append(toScan, &transaction.Reference)
		}

		if err := rows.Scan(toScan...); err != nil {
			return nil, err
		}

		return transaction, nil
	}

	rows, err := s.db.From(
		"transactions",
	).SelectDistinct(
		fields()...,
	).Join(
		goqu.T("postings"),
		goqu.On(goqu.Ex{"transactions.id": goqu.I("postings.transaction")}),
	).Join(
		goqu.T("accounts"),
		goqu.On(
			goqu.Ex{"accounts.id": goqu.I("postings.account")},
			goqu.I("accounts.name").RegexpLike(req.AccountName),
		),
	).Where(
		goqu.C("date").Between(
			goqu.Range(req.From, req.To),
		),
	).Executor().Query()
	if err != nil {
		return nil, err
	}

	var data []*pb.Transaction
	byId := make(map[string]int32)
	var i int32
	for rows.Next() {
		transaction, err := scan(rows)
		if err != nil {
			return nil, err
		}

		data = append(data, transaction)
		byId[transaction.Id] = i
		i++
	}

	if !req.Fields.Postings {
		return &pb.Transactions{
			Data: data,
			ById: byId,
		}, nil
	}

	for _, v := range data {
		rows, err := s.db.From(
			"postings",
		).Select(
			"id",
			"account",
			"units",
			"cost",
			"price",
		).Where(
			goqu.C("transaction").Eq(v.Id),
		).Executor().Query()
		if err != nil {
			return nil, err
		}

		var dataPostings []*pb.Posting
		byIdPostings := make(map[string]int32)
		i = 0
		for rows.Next() {
			posting := &pb.Posting{
				Transaction: v.Id,
			}
			if err := rows.Scan(
				&posting.Id,
				&posting.Account,
				&posting.Units,
				&posting.Cost,
				&posting.Price,
			); err != nil {
				return nil, err
			}

			dataPostings = append(dataPostings, posting)
			byIdPostings[posting.Id] = i
			i++
		}
		v.Postings = &pb.Postings{
			Data: dataPostings,
			ById: byIdPostings,
		}
	}

	return &pb.Transactions{
		Data: data,
		ById: byId,
	}, nil
}
