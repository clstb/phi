package server

import (
	"context"
	"database/sql"

	"github.com/clstb/phi/go/pkg/auth/db"
	"github.com/clstb/phi/go/pkg/auth/pb"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) Register(ctx context.Context, user *pb.User) (*pb.Token, error) {
	tx, ok := ctx.Value("tx").(*sql.Tx)
	if !ok {
		return nil, status.Error(codes.Internal, "missing tx")
	}
	q := db.New(tx)

	password_hashed, err := bcrypt.GenerateFromPassword(
		[]byte(user.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, status.Error(codes.Internal, "hashing password")
	}

	_, err = q.CreateUser(ctx, db.CreateUserParams{
		Name:     user.Name,
		Password: password_hashed,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "db: creating user")
	}

	return s.Login(ctx, user)
}
