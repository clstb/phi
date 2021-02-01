// +build integration

package core_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/clstb/phi/pkg/core"
	db "github.com/clstb/phi/pkg/db/core"
	"github.com/clstb/phi/pkg/pb"
	"github.com/gofrs/uuid"
	"github.com/matryer/is"
)

func TestCreateAccount(t *testing.T) {
	is := is.New(t)

	server := core.New(core.WithDB(sqlDB))

	type test struct {
		do    func() (*sql.Tx, *pb.Account, error)
		check func(*sql.Tx, *pb.Account, error) *sql.Tx
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
		do: func() (*sql.Tx, *pb.Account, error) {
			tx, ctx := tx()
			res, err := server.CreateAccount(
				ctx,
				&pb.Account{
					Name: "account-test",
				},
			)
			return tx, res, err
		},
		check: func(tx *sql.Tx, a *pb.Account, e error) *sql.Tx {
			is.NoErr(e)
			is.Equal(a.Name, "account-test")
			q := db.New(tx)
			accounts, err := q.GetAccounts(context.Background(), db.GetAccountsParams{
				Name: "account-test",
			})
			is.NoErr(err)
			is.Equal(len(accounts), 1)
			return tx
		},
	})

	for _, t := range tests {
		is.NoErr(t.check(t.do()).Rollback())
	}
}

func TestGetAccounts(t *testing.T) {
	is := is.New(t)

	server := core.New(core.WithDB(sqlDB))

	type test struct {
		do    func() (*sql.Tx, *pb.Accounts, error)
		check func(*sql.Tx, *pb.Accounts, error) *sql.Tx
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
		do: func() (*sql.Tx, *pb.Accounts, error) {
			tx, ctx := tx()
			_, err := server.CreateAccount(ctx, &pb.Account{
				Name: "account-test",
			})
			is.NoErr(err)

			res, err := server.GetAccounts(
				ctx,
				&pb.AccountsQuery{
					Name: "^account-test$",
				},
			)

			return tx, res, err
		},
		check: func(tx *sql.Tx, a *pb.Accounts, e error) *sql.Tx {
			is.NoErr(e)
			is.Equal(len(a.Data), 1)
			is.Equal(a.Data[0].Name, "account-test")
			return tx
		},
	})

	for _, t := range tests {
		is.NoErr(t.check(t.do()).Rollback())
	}
}
