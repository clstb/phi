package server

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	db "github.com/clstb/phi/go/pkg/services/auth/db"
	pb "github.com/clstb/phi/go/pkg/services/auth/pb"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) Register(
	ctx context.Context,
	req *pb.User,
) (*pb.JWT, error) {
	if req.Name == "" {
		return nil, status.Error(
			codes.InvalidArgument,
			"empty field: name",
		)
	}

	log.Println("new registration")

	password_hashed, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, err
	}

	tx, ok := ctx.Value("tx").(*sql.Tx)
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
