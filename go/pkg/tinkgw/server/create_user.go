package server

import (
	"context"
	"errors"

	"github.com/clstb/phi/go/pkg/tink"
	"github.com/clstb/phi/go/pkg/tinkgw/db"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) createUser(ctx context.Context, id uuid.UUID) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	q := db.New(tx)

	res, err := s.tink.CreateUser(&tink.CreateUserReq{
		ExternalUserID: id.String(),
		Market:         "DE",
		Locale:         "de_DE",
	})
	if err != nil && !errors.Is(err, tink.ErrUserExists) {
		s.logger.Error("tink: creating user", zap.Error(err))
		return status.Error(codes.Internal, "tink: creating user")
	}

	_, err = q.CreateUser(ctx, db.CreateUserParams{
		ID:     id,
		TinkID: res.UserID,
	})
	if err != nil {
		return status.Error(codes.Internal, "db: creating user")
	}

	return tx.Commit()
}
