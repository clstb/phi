package server_test

import (
	"context"
	"testing"

	"github.com/clstb/phi/pkg/pb"
	"github.com/gofrs/uuid"
	"github.com/matryer/is"
)

func TestCreateAccount(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()

	client := coreClient()

	type test struct {
		do    func() (*pb.Account, error)
		check func(*pb.Account, error)
	}
	var tests []test
	add := func(t test) {
		tests = append(tests, t)
	}

	add(test{
		do: func() (*pb.Account, error) {
			return client.CreateAccount(
				ctx,
				&pb.Account{
					Name: "account-test",
				},
			)
		},
		check: func(a *pb.Account, e error) {
			is.NoErr(e)
			is.Equal(a.Name, "account-test")
			is.NoErr(db.DeleteAccount(ctx, uuid.FromStringOrNil(a.Id)))
		},
	})

	for _, t := range tests {
		t.check(t.do())
	}
}

func TestGetAccounts(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()

	client := coreClient()

	type test struct {
		do    func() (*pb.Accounts, error)
		check func(*pb.Accounts, error)
	}
	var tests []test
	add := func(t test) {
		tests = append(tests, t)
	}

	add(test{
		do: func() (*pb.Accounts, error) {
			_, err := client.CreateAccount(
				ctx,
				&pb.Account{
					Name: "account-test",
				},
			)
			is.NoErr(err)
			return client.GetAccounts(
				ctx,
				&pb.AccountsQuery{
					Name: "^account-test$",
				},
			)
		},
		check: func(a *pb.Accounts, e error) {
			is.NoErr(e)
			is.Equal(a.Data[0].Name, "account-test")
			is.NoErr(db.DeleteAccount(ctx, uuid.FromStringOrNil(a.Data[0].Id)))
		},
	})
}
