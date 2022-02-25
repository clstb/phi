package main

import (
	"fmt"

	pb "github.com/clstb/phi/go/pkg/auth/pb"
	"github.com/clstb/phi/go/pkg/auth/server"
	"github.com/clstb/phi/go/pkg/interceptor"
	"github.com/clstb/phi/go/pkg/util"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func Server(ctx *cli.Context) error {
	db, ok := ctx.Context.Value("db").(*pgxpool.Pool)
	if !ok {
		return fmt.Errorf("missing db")
	}

	logger, ok := ctx.Context.Value("logger").(*zap.Logger)
	if !ok {
		return fmt.Errorf("missing logger")
	}

	// create grpc server
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			interceptor.TXUnary(db),
		)),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			interceptor.TXStream(db),
		)),
	)

	// create / register server implementation
	authServer := server.New([]byte(ctx.String("signing-secret")))
	pb.RegisterAuthServer(grpcServer, authServer)

	return util.ListenGRPC(grpcServer, logger, ctx.Int("port"))
}
