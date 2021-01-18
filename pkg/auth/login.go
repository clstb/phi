package auth

import (
	"context"
	"time"

	"github.com/clstb/phi/pkg/auth/db"
	"github.com/clstb/phi/pkg/pb"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func (s *Server) Login(
	ctx context.Context,
	req *pb.User,
) (*pb.JWT, error) {
	q := db.New(s.db)
	user, err := q.GetUserByName(ctx, req.Name)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword(
		user.Password,
		[]byte(req.Password),
	); err != nil {
		return nil, err
	}

	claims := Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(30 * time.Minute).Unix(),
		},
		UserID:   user.ID.String(),
		UserName: user.Name,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenSigned, err := token.SignedString(s.signingSecret)
	if err != nil {
		return nil, err
	}

	return &pb.JWT{
		AccessToken: tokenSigned,
	}, nil
}
