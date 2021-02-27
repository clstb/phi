// +build integration

package core_test

import (
	"context"
	"database/sql"
	"io"
	"testing"
	"time"

	"github.com/clstb/phi/pkg/core"
	"github.com/clstb/phi/pkg/pb"
	"github.com/clstb/phi/pkg/util"
	"github.com/gofrs/uuid"
	"github.com/matryer/is"
	"google.golang.org/grpc/metadata"
)

type TestTransactionsServer struct {
	*util.FakeGRPCServerStream
	ctx context.Context
	req []*pb.Transaction
	res *pb.Transactions
}

func NewTestTransactionsServer(
	ctx context.Context,
	req ...*pb.Transaction,
) *TestTransactionsServer {
	s := &TestTransactionsServer{
		ctx: ctx,
		req: req,
	}

	stream := &util.FakeGRPCServerStream{
		OnSetHeader: func(m metadata.MD) error {
			return nil
		},
		OnSendHeader: func(m metadata.MD) error {
			return nil
		},
		OnSetTrailer: func(m metadata.MD) {},
		OnContext: func() context.Context {
			return s.ctx
		},
		OnSendMsg: func(m interface{}) error {
			return nil
		},
		OnRecvMsg: func(m interface{}) error {
			return nil
		},
	}
	s.FakeGRPCServerStream = stream

	return s
}

func (s *TestTransactionsServer) SendAndClose(m *pb.Transactions) error {
	s.res = m
	return nil
}

func (s *TestTransactionsServer) Recv() (*pb.Transaction, error) {
	i := len(s.req)
	if i == 0 {
		return nil, io.EOF
	}

	var v *pb.Transaction
	v, s.req = s.req[i-1], s.req[:i-1]
	return v, nil
}

func TestCreateTransaction(t *testing.T) {
	is := is.New(t)

	server := core.New(core.WithDB(sqlDB))

	type test struct {
		do    func() (*sql.Tx, *pb.Transactions, error)
		check func(*sql.Tx, *pb.Transactions, error) *sql.Tx
	}
	var tests []test
	add := func(t test) {
		tests = append(tests, t)
	}
	tx := func() (*sql.Tx, context.Context) {
		tx, err := sqlDB.Begin()
		is.NoErr(err)
		ctx := context.WithValue(
			context.Background(),
			"tx",
			tx,
		)

		return tx, context.WithValue(
			ctx,
			"sub",
			uuid.Nil.String(),
		)
	}

	add(test{
		do: func() (*sql.Tx, *pb.Transactions, error) {
			tx, ctx := tx()

			account, err := server.CreateAccount(ctx, &pb.Account{
				Name: "account-1",
			})
			is.NoErr(err)

			stream := NewTestTransactionsServer(
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
			err = server.CreateTransactions(stream)

			return tx, stream.res, err
		},
		check: func(tx *sql.Tx, t *pb.Transactions, e error) *sql.Tx {
			return tx
		},
	})

	for _, t := range tests {
		is.NoErr(t.check(t.do()).Rollback())
	}
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
