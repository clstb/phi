package interceptor

import (
	"context"
	"fmt"

	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/gofrs/uuid"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type Claims struct {
	UserID uuid.UUID
	jwtgo.StandardClaims
}

func serverAuth(signingSecret []byte) func(ctx context.Context) (context.Context, error) {
	return func(ctx context.Context) (context.Context, error) {
		jwt, err := grpc_auth.AuthFromMD(ctx, "bearer")
		if err != nil {
			return nil, err
		}

		token, err := jwtgo.ParseWithClaims(
			jwt,
			&Claims{},
			func(token *jwtgo.Token) (interface{}, error) {
				_, ok := token.Method.(*jwtgo.SigningMethodHMAC)
				if !ok {
					return nil, fmt.Errorf("unexpected token signing method")
				}

				return signingSecret, nil
			},
		)
		if err != nil {
			return nil, fmt.Errorf("invalid token: %w", err)
		}

		c, ok := token.Claims.(*Claims)
		if !ok {
			return nil, fmt.Errorf("invalid token claims")
		}

		userID, err := uuid.FromString(c.Subject)
		if err != nil {
			return nil, fmt.Errorf("parsing subject")
		}
		c.UserID = userID

		ctx = context.WithValue(ctx, "claims", c)
		return ctx, nil
	}
}

func ServerAuthUnary(signingSecret []byte) grpc.UnaryServerInterceptor {
	return grpc_auth.UnaryServerInterceptor(serverAuth(signingSecret))
}

func ServerAuthStream(signingSecret []byte) grpc.StreamServerInterceptor {
	return grpc_auth.StreamServerInterceptor(serverAuth(signingSecret))
}

func ClientAuthUnary(accessToken string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctx = metadata.AppendToOutgoingContext(
			ctx,
			"authorization",
			fmt.Sprintf("Bearer %s", accessToken),
		)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func ClientAuthStream(accessToken string) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		ctx = metadata.AppendToOutgoingContext(
			ctx,
			"authorization",
			fmt.Sprintf("Bearer %s", accessToken),
		)
		return streamer(ctx, desc, cc, method, opts...)
	}
}
