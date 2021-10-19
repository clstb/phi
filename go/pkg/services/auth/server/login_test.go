//go:build integration
// +build integration

package server_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/clstb/phi/go/pkg/services/auth/db"
	"github.com/clstb/phi/go/pkg/services/auth/pb"
	"github.com/clstb/phi/go/pkg/services/auth/server"
	"github.com/matryer/is"
	"golang.org/x/crypto/bcrypt"
)

func TestLogin(t *testing.T) {
	t.Parallel()

	is := is.New(t)

	server := server.New([]byte("secret"))

	tx := func(ctx context.Context) (*sql.Tx, context.Context) {
		tx, err := testdb.Begin()
		is.NoErr(err)

		return tx, context.WithValue(
			context.Background(),
			"tx",
			tx,
		)
	}

	tests := []func(){
		func() {
			tx, ctx := tx(context.Background())
			defer func() {
				is.NoErr(tx.Rollback())
			}()

			password_hashed, err := bcrypt.GenerateFromPassword(
				[]byte("password"),
				bcrypt.DefaultCost,
			)
			is.NoErr(err)

			q := db.New(tx)
			q.CreateUser(ctx, db.CreateUserParams{
				Name:     "test-user",
				Password: password_hashed,
			})

			_, err = server.Login(ctx, &pb.User{
				Name:     "test-user",
				Password: "password",
			})
			is.NoErr(err)
		},
	}

	for _, t := range tests {
		t()
	}
}
