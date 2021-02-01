// +build integration

package core_test

import (
	"testing"
)

func TestCreateTransaction(t *testing.T) {
	/*
		is := is.New(t)
		ctx := context.Background()

		client := client()
		type test struct {
			do    func() (*pb.Transaction, error)
			check func(*pb.Transaction, error)
		}
		var tests []test
		add := func(t test) {
			tests = append(tests, t)
		}

		account, err := client.CreateAccount(
			ctx,
			&pb.Account{
				Name: "account-test",
			},
		)
		is.NoErr(err)
		defer func() {
			is.NoErr(db.DeleteAccount(ctx, uuid.FromStringOrNil(account.Id)))
		}()

		add(test{
			do: func() (*pb.Transaction, error) {
				return client.CreateTransaction(
					ctx,
					&pb.Transaction{
						Date: time.Now().Format("2006-01-02"),
						Postings: []*pb.Posting{
							{
								Account: account.Id,
								Units:   "1 EUR",
							},
							{
								Account: account.Id,
								Units:   "-1 EUR",
							},
						},
					},
				)
			},
			check: func(t *pb.Transaction, e error) {
				is.NoErr(e)
				is.NoErr(db.DeleteTransaction(ctx, uuid.FromStringOrNil(t.Id)))
			},
		})

		for _, t := range tests {
			t.check(t.do())
		}
	*/
}

func TestGetTransactions(t *testing.T) {
	/*
		is := is.New(t)
		ctx := context.Background()

		client := client()
		type test struct {
			do    func() (*pb.Transactions, error)
			check func(*pb.Transactions, error)
		}
		var tests []test
		add := func(t test) {
			tests = append(tests, t)
		}

		account, err := client.CreateAccount(
			ctx,
			&pb.Account{
				Name: "account-test",
			},
		)
		is.NoErr(err)
		defer func() {
			is.NoErr(db.DeleteAccount(ctx, uuid.FromStringOrNil(account.Id)))
		}()

		add(test{
			do: func() (*pb.Transactions, error) {
				transactions, err := client.CreateTransaction(
					ctx,
					&pb.Transaction{
						Date: time.Now().Format("2006-01-02"),
						Postings: []*pb.Posting{
							{
								Account: account.Id,
								Units:   "1 EUR",
							},
							{
								Account: account.Id,
								Units:   "-1 EUR",
							},
						},
					},
				)
				fmt.Println(transactions)
				is.NoErr(err)
				return client.GetTransactions(
					ctx,
					&pb.TransactionsQuery{
						AccountName: "^account-.*",
					},
				)
			},
			check: func(t *pb.Transactions, e error) {
				is.NoErr(e)
				is.NoErr(db.DeleteTransaction(ctx, uuid.FromStringOrNil(t.Data[0].Id)))
			},
		})

		for _, t := range tests {
			t.check(t.do())
		}
	*/
}
