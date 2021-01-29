package auth

import (
	"context"

	db "github.com/clstb/phi/pkg/db/auth"
	"github.com/clstb/phi/pkg/pb"
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

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	q := db.New(tx)

	_, err = q.CreateUser(ctx, db.CreateUserParams{
		Name:     req.Name,
		Password: password_hashed,
	})
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return s.Login(ctx, req)
}
