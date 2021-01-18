package auth

import (
	"context"
	"fmt"

	"github.com/clstb/phi/pkg/pb"
	"github.com/dgrijalva/jwt-go"
)

func (s *Server) Verify(
	ctx context.Context,
	req *pb.JWT,
) (*pb.Claims, error) {
	token, err := jwt.ParseWithClaims(
		req.AccessToken,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("unexpected token signing method")
			}

			return s.signingSecret, nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims.PB(), nil
}
