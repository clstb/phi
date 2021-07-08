package auth

import (
	"context"
	"fmt"

	db "github.com/clstb/phi/pkg/db/auth"
	"github.com/clstb/phi/pkg/pb"
	"github.com/jackc/pgx/v4"
	"golang.org/x/crypto/bcrypt"
)

func (s *Server) Register(
	ctx context.Context,
	req *pb.User,
) (*pb.JWT, error) {
	password_hashed, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, err
	}

	tx, ok := ctx.Value("tx").(pgx.Tx)
	if !ok {
		return nil, fmt.Errorf("context: missing database transaction")
	}
	q := db.New(tx)

	_, err = q.CreateUser(ctx, db.CreateUserParams{
		Name:     req.Name,
		Password: password_hashed,
	})
	if err != nil {
		return nil, err
	}

	return s.Login(ctx, req)
}
