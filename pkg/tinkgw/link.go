package tinkgw

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/clstb/phi/pkg/pb"
	"github.com/clstb/phi/pkg/tink"
	"github.com/gofrs/uuid"
)

func (s *Server) Link(
	ctx context.Context,
	req *pb.LinkReq,
) (*pb.LinkRes, error) {
	subStr, ok := ctx.Value("sub").(string)
	if !ok {
		return nil, fmt.Errorf("context: missing subject")
	}
	sub, err := uuid.FromString(subStr)
	if err != nil {
		return nil, err
	}

	_, err = s.tink.CreateUser(&tink.CreateUserReq{
		ExternalUserID: sub.String(),
		Market:         req.Market,
		Locale:         req.Locale,
	})
	if err != nil {
		if !errors.Is(err, tink.ErrUserExists) {
			return nil, err
		}
	}

	code, err := s.tink.AuthorizeGrantDelegate(&tink.AuthorizeGrantDelegateReq{
		ResponseType:   "code",
		ActorClientID:  tink.TinkActorClientID,
		ExternalUserID: sub.String(),
		Scope:          "authorization:read,authorization:grant,credentials:refresh,credentials:read,credentials:write,providers:read,user:read",
		IDHint:         sub.String(),
	})
	if err != nil {
		return nil, err
	}

	state, err := uuid.NewV4()
	if err != nil {
		return nil, fmt.Errorf("failed generating unique state: %w", err)
	}
	s.states[state.String()] = sub
	go func() {
		time.Sleep(5 * time.Minute)
		delete(s.states, state.String())
	}()

	tinkLink := fmt.Sprintf(
		"https://link.tink.com/1.0/transactions/connect-accounts?client_id=%s&state=%s&redirect_uri=%s&authorization_code=%s&market=%s&locale=%s",
		s.clientID,
		state.String(),
		s.callbackURL,
		code,
		req.Market,
		req.Locale,
	)

	return &pb.LinkRes{
		TinkLink: tinkLink,
	}, nil
}
