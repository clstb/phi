package server

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/clstb/phi/pkg/db"
	"github.com/clstb/phi/pkg/pb"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type auth struct {
	pb.UnimplementedAuthServer
	db            *sql.DB
	signingSecret []byte
}

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

func (s *auth) Register(
	ctx context.Context,
	req *pb.User,
) (*pb.JWT, error) {
	password_hashed, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	q := db.New(tx)

	_, err = q.CreateUser(ctx, db.CreateUserParams{
		Name:     req.Name,
		Password: password_hashed,
	})
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return s.Login(ctx, req)
}

func (s *auth) Login(
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
		UserID: user.ID.String(),
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

func (s *auth) Verify(
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
