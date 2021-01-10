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
				"units":       posting.Units.Value,
				"units_cur":   posting.Units.Currency,
			}
			if posting.Cost != nil {
				record["cost"] = posting.Cost.Value
				record["cost_cur"] = posting.Cost.Currency
			}
			if posting.Price != nil {
				record["price"] = posting.Price.Value
				record["price_cur"] = posting.Price.Currency
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
		fields := []interface{}{"id"}
		if req.Fields.Date {
			fields = append(fields, "date")
		}
		if req.Fields.Entity {
			fields = append(fields, "entity")
		}
		if req.Fields.Reference {
			fields = append(fields, "reference")
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
	).Select(
		fields()...,
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

	if req.Fields.Postings {
		for _, v := range data {
			postings, err := s.GetPostings(
				ctx,
				&pb.PostingsQuery{
					Fields: &pb.PostingsQuery_Fields{
						Account:     true,
						Transaction: true,
						Units:       true,
						Cost:        true,
						Price:       true,
					},
					Transaction: v.Id,
				},
			)
			if err != nil {
				return nil, err
			}
			v.Postings = postings
		}
	}

	return &pb.Transactions{
		Data: data,
		ById: byId,
	}, nil
}

func (s *core) GetPostings(
	ctx context.Context,
	req *pb.PostingsQuery,
) (*pb.Postings, error) {
	fields := func() []interface{} {
		fields := []interface{}{"id"}
		add := func(i interface{}) {
			fields = append(fields, i)
		}
		if req.Fields.Account {
			add("account")
		}
		if req.Fields.Transaction {
			add("transaction")
		}
		if req.Fields.Units {
			add("units")
			add("units_cur")
		}
		if req.Fields.Cost {
			add("cost")
			add("cost_cur")
		}
		if req.Fields.Price {
			add("price")
			add("price_cur")
		}

		return fields
	}
	scan := func(rows *sql.Rows) (*pb.Posting, error) {
		posting := &pb.Posting{}
		toScan := []interface{}{&posting.Id}
		add := func(i interface{}) {
			toScan = append(toScan, i)
		}

		if req.Fields.Account {
			add(&posting.Account)
		}
		if req.Fields.Transaction {
			add(&posting.Transaction)
		}
		if req.Fields.Units {
			units := &pb.Amount{}
			add(&units.Value)
			add(&units.Currency)
			posting.Units = units
		}
		if req.Fields.Cost {
			add(&posting.Cost.Value)
			add(&posting.Cost.Currency)
		}
		if req.Fields.Price {
			add(&posting.Price.Value)
			add(&posting.Price.Currency)
		}

		if err := rows.Scan(toScan...); err != nil {
			return nil, err
		}

		return posting, nil
	}
	filter := func() []goqu.Expression {
		var filter []goqu.Expression
		add := func(e goqu.Expression) {
			filter = append(filter, e)
		}

		if req.Account != "" {
			add(goqu.C("account").Eq(req.Account))
		}
		if req.AccountName != "" {
			add(goqu.C("account_name").RegexpLike(req.AccountName))
		}
		if req.Transaction != "" {
			add(goqu.C("transaction").Eq(req.Transaction))
		}

		from := "-infinity"
		if req.From != "" {
			from = req.From
		}
		to := "+infinity"
		if req.To != "" {
			to = req.To
		}

		add(goqu.C("date").Between(goqu.Range(from, to)))

		return filter
	}

	rows, err := s.db.From(
		"postings_joined",
	).Select(
		fields()...,
	).Where(
		filter()...,
	).Executor().Query()
	if err != nil {
		return nil, err
	}

	var data []*pb.Posting
	byId := make(map[string]int32)
	var i int32
	for rows.Next() {
		posting, err := scan(rows)
		if err != nil {
			return nil, err
		}

		data = append(data, posting)
		byId[posting.Id] = i
	}

	return &pb.Postings{
		Data: data,
		ById: byId,
	}, nil
}
