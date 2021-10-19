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

func (s *Server) Load(ctx context.Context, req *pb.Ledger) (*pb.Ledger, error) {
	claims, ok := ctx.Value("claims").(*interceptor.Claims)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "unauthenticated")
	}

	userID, err := uuid.FromString(claims.Subject)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "claims: subject: %s", err)
	}

	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id: empty")
	}

	ledgerID, err := uuid.FromString(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "id: %s", err)
	}

	tx, ok := ctx.Value("tx").(*sql.Tx)
	if !ok {
		return nil, fmt.Errorf("context: missing database transaction")
	}
	q := db.New(tx)

	ledger, err := q.GetLedger(ctx, ledgerID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "db: getting ledger: %s", err)
	}

	if ledger.UserID != userID {
		return nil, status.Errorf(codes.Unauthenticated, "unauthorized")
	}

	return &pb.Ledger{
		Id:   ledgerID.String(),
		Data: ledger.Data,
		Dk:   ledger.Dk,
	}, nil
}
