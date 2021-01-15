package server_test

import (
	"context"
	"testing"

	"github.com/clstb/phi/pkg/pb"
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
		},
	})

	for _, t := range tests {
		t.check(t.do())
	}
}
