package server

import (
	"context"
	"errors"
	"fmt"

	"github.com/clstb/phi/go/pkg/interceptor"
	"github.com/clstb/phi/go/pkg/services/tinkgw/pb"
	"github.com/clstb/phi/go/pkg/services/tinkgw/tink"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetLink(ctx context.Context, req *pb.GetLinkReq) (*pb.Link, error) {
	claims, ok := ctx.Value("claims").(*interceptor.Claims)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "unauthenticated")
	}

	_, err := s.tink.CreateUser(&tink.CreateUserReq{
		ExternalUserID: claims.Subject,
		Market:         "DE",    // TODO
		Locale:         "de_DE", // TODO
	})
	if err != nil {
		if !errors.Is(err, tink.ErrUserExists) {
			return nil, status.Errorf(codes.Internal, "tink: creating user: %s", err)
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
