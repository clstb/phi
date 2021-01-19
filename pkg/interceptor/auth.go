package interceptor

import (
	"context"

	"github.com/clstb/phi/pkg/pb"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc"
)

func auth(client pb.AuthClient) func(ctx context.Context) (context.Context, error) {
	return func(ctx context.Context) (context.Context, error) {
		token, err := grpc_auth.AuthFromMD(ctx, "bearer")
		if err != nil {
			return nil, err
		}

		claims, err := client.Verify(
			ctx,
			&pb.JWT{
				AccessToken: token,
			},
		)
		if err != nil {
			return nil, err
		}

		ctx = context.WithValue(ctx, "sub", claims.Subject)
		return ctx, nil
	}
}

func AuthUnary(client pb.AuthClient) grpc.UnaryServerInterceptor {
	return grpc_auth.UnaryServerInterceptor(auth(client))
}

func AuthStream(client pb.AuthClient) grpc.StreamServerInterceptor {
	return grpc_auth.StreamServerInterceptor(auth(client))
}
