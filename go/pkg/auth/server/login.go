package server

import (
	"context"
	"time"

	"github.com/clstb/phi/go/pkg/auth/db"
	"github.com/clstb/phi/go/pkg/auth/pb"
	"github.com/dgrijalva/jwt-go"
	"github.com/jackc/pgx/v4"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) Login(ctx context.Context, user *pb.User) (*pb.Token, error) {
	tx, ok := ctx.Value("tx").(pgx.Tx)
	if !ok {
		return nil, status.Error(codes.Internal, "missing tx")
	}
	q := db.New(tx)

	storedUser, err := q.GetUserByName(ctx, user.Name)
	if err != nil {
		return nil, status.Error(codes.Internal, "db: getting user")
	}

	if err := bcrypt.CompareHashAndPassword(
		storedUser.Password,
		[]byte(user.Password),
	); err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid password")
	}

	now := time.Now()
	tokenExpiry := now.Add(5 * time.Minute).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, struct {
		jwt.StandardClaims
	}{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: tokenExpiry,
			Subject:   storedUser.ID.String(),
		},
	})
	tokenSigned, err := token.SignedString(s.signingSecret)
	if err != nil {
		return nil, status.Error(codes.Internal, "signing token")
	}

	return &pb.Token{
		AccessToken: tokenSigned,
		ExpiresAt:   tokenExpiry,
		TokenType:   "Bearer",
	}, nil
}
