package auth

import (
	"github.com/clstb/phi/pkg/pb"
	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	jwt.StandardClaims
}

func (c *Claims) PB() *pb.Claims {
	return &pb.Claims{
		Audience:  c.Audience,
		ExpiresAt: c.ExpiresAt,
		Id:        c.Id,
		IssuedAt:  c.IssuedAt,
		Issuer:    c.Issuer,
		NotBefore: c.NotBefore,
		Subject:   c.Subject,
	}
}

func ClaimsFromPB(claims *pb.Claims) Claims {
	return Claims{
		StandardClaims: jwt.StandardClaims{
			Audience:  claims.Audience,
			ExpiresAt: claims.ExpiresAt,
			Id:        claims.Id,
			IssuedAt:  claims.IssuedAt,
			Issuer:    claims.Issuer,
			NotBefore: claims.NotBefore,
			Subject:   claims.Subject,
		},
	}
}
