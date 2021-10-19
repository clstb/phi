package server

import (
	"context"
	"time"

	"github.com/clstb/phi/go/pkg/interceptor"
	"github.com/clstb/phi/go/pkg/services/tinkgw/pb"
	"github.com/clstb/phi/go/pkg/services/tinkgw/tink"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetToken(
	ctx context.Context,
	req *pb.Token,
) (*pb.Token, error) {
	claims, ok := ctx.Value("claims").(*interceptor.Claims)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "unauthenticated")
	}

	oauthReq := &tink.OAuthTokenReq{
		ClientID:     s.clientID,
		ClientSecret: s.clientSecret,
	}
	if req.RefreshToken == "" {
		code, err := s.tink.AuthorizeGrant(&tink.AuthorizeGrantReq{
			ExternalUserID: claims.Subject,
			Scope:          "transactions:read,accounts:read",
		})
		if err != nil {
			return nil, status.Errorf(codes.Internal, "tink: authorize grant: %w", err)
		}
		oauthReq.GrantType = "authorization_code"
		oauthReq.Code = code
	} else {
		oauthReq.GrantType = "refresh_token"
		oauthReq.RefreshToken = req.RefreshToken
	}

	token, err := s.tink.OAuthToken(oauthReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "tink: oauth token: %w", err)
	}

	s.logger.Info(
		"requested oauth token",
		zap.String("grant_type", oauthReq.GrantType),
		zap.String("kind", "user_token"),
	)

	return &pb.Token{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		TokenType:    token.TokenType,
		Scope:        token.Scope,
		ExpiresAt: time.Now().Add(
			time.Second * time.Duration(
				token.ExpiresIn,
			)).Unix(),
	}, nil
}
