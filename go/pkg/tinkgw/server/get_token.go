package server

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/clstb/phi/go/pkg/interceptor"
	"github.com/clstb/phi/go/pkg/tink"
	"github.com/clstb/phi/go/pkg/tinkgw/db"
	"github.com/clstb/phi/go/pkg/tinkgw/pb"
	"github.com/jackc/pgx/v4"
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

	oauthReq := &tink.OAuthTokenReq{
		ClientID:     s.clientID,
		ClientSecret: s.clientSecret,
	}
	if req.RefreshToken == "" {
		code, err := s.tink.AuthorizeGrant(&tink.AuthorizeGrantReq{
			ExternalUserID: claims.Subject,
			Scope:          "transactions:read,accounts:read,provider-consents:read",
		})
		if err != nil {
			return nil, status.Errorf(codes.Internal, "tink: authorize grant: %v", err)
		}
		oauthReq.GrantType = "authorization_code"
		oauthReq.Code = code
	} else {
		oauthReq.GrantType = "refresh_token"
		oauthReq.RefreshToken = req.RefreshToken
	}

	token, err := s.tink.OAuthToken(oauthReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "tink: oauth token: %v", err)
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
