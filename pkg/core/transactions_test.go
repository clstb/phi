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
	t.Parallel()

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

			// TODO: fix parsing of all uuid's
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
			is.NoErr(e)
			return tx
		},
	})

	for _, t := range tests {
		is.NoErr(t.check(t.do()).Rollback())
	}
}

func TestGetTransactions(t *testing.T) {
	t.Parallel()

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
			is.NoErr(err)

			res, err := server.GetTransactions(
				ctx,
				&pb.TransactionsQuery{},
			)

			return tx, res, err
		},
		check: func(tx *sql.Tx, t *pb.Transactions, e error) *sql.Tx {
			is.NoErr(e)
			is.Equal(len(t.Data), 1)
			return tx
		},
	})

	for _, t := range tests {
		is.NoErr(t.check(t.do()).Rollback())
	}
}
