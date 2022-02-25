package server

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/clstb/phi/go/pkg/interceptor"
	"github.com/clstb/phi/go/pkg/tink"
	"github.com/clstb/phi/go/pkg/tinkgw/db"
	"github.com/clstb/phi/go/pkg/tinkgw/pb"
	"github.com/jackc/pgx/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetLink(ctx context.Context, req *pb.GetLinkReq) (*pb.Link, error) {
	claims, ok := ctx.Value("claims").(*interceptor.Claims)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "unauthenticated")
	}

	tx, ok := ctx.Value("tx").(pgx.Tx)
	if !ok {
		return nil, status.Error(codes.Internal, "missing tx")
	}
	q := db.New(tx)

	_, err := q.GetUserByID(ctx, claims.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			if err := s.createUser(ctx, claims.UserID); err != nil {
				return nil, err
			}
		} else {
			return nil, status.Error(codes.Internal, "db: reading user")
		}
	}

	code, err := s.tink.AuthorizeGrantDelegate(&tink.AuthorizeGrantDelegateReq{
		ResponseType:   "code",
		ActorClientID:  tink.TinkActorClientID,
		ExternalUserID: claims.Subject,
		IDHint:         claims.Subject, // TODO
		Scope:          "authorization:read,authorization:grant,credentials:refresh,credentials:read,credentials:write,providers:read,user:read",
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "tink: authorize grant delegate: %s", err)
	}

	link := fmt.Sprintf(
		"https://link.tink.com/1.0/transactions/connect-accounts?client_id=%s&redirect_uri=%s&market=%s&locale=%s&authorization_code=%s",
		s.clientID,
		s.callbackURL,
		req.Market,
		req.Locale,
		code,
	)

	return &pb.Link{
		Link: link,
	}, nil
}
