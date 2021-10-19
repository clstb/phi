package server

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/clstb/phi/go/pkg/interceptor"
	"github.com/clstb/phi/go/pkg/services/bookkeeper/db"
	"github.com/clstb/phi/go/pkg/services/bookkeeper/pb"
	"github.com/gofrs/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) Save(ctx context.Context, req *pb.Ledger) (*pb.Ledger, error) {
	claims, ok := ctx.Value("claims").(*interceptor.Claims)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "unauthenticated")
	}

	userID, err := uuid.FromString(claims.Subject)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "claims: subject: %s", err)
	}

	tx, ok := ctx.Value("tx").(*sql.Tx)
	if !ok {
		return nil, fmt.Errorf("context: missing database transaction")
	}
	q := db.New(tx)

	var ledger db.Ledger
	if req.Id == "" {
		ledger, err = q.CreateLedger(ctx, db.CreateLedgerParams{
			UserID: userID,
			Data:   req.Data,
			Dk:     req.Dk,
		})
		if err != nil {
			return nil, status.Errorf(codes.Internal, "db: creating ledger: %s", err)
		}
	} else {
		ledgerID, err := uuid.FromString(req.Id)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "id: %s", err)
		}

		ledger, err = q.UpdateLedger(ctx, db.UpdateLedgerParams{
			ID:   ledgerID,
			Data: req.Data,
			Dk:   req.Dk,
		})
		if err != nil {
			return nil, status.Errorf(codes.Internal, "db: updating ledger: %s", err)
		}

	}

	return &pb.Ledger{
		Id:   ledger.ID.String(),
		Data: ledger.Data,
		Dk:   ledger.Dk,
	}, nil
}
