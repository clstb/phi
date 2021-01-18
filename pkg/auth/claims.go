package auth

import (
	"github.com/clstb/phi/pkg/pb"
	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	jwt.StandardClaims
	UserID   string `json:"user_id"`
	UserName string `json:"user_name"`
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
		UserId:    c.UserID,
		UserName:  c.UserName,
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
		UserID:   claims.UserId,
		UserName: claims.UserName,
	}
}
