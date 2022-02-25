package main

import (
	"fmt"

	"github.com/clstb/phi/go/pkg/interceptor"
	pb "github.com/clstb/phi/go/pkg/tinkgw/pb"
	"github.com/clstb/phi/go/pkg/tinkgw/server"
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

	signingSecret := ctx.String("signing-secret")

	// create grpc server
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			interceptor.ServerAuthUnary([]byte(signingSecret)),
			interceptor.TXUnary(db),
		)),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			interceptor.ServerAuthStream([]byte(signingSecret)),
			interceptor.TXStream(db),
		)),
	)

	// create / register server implementation
	tinkgwServer := server.New(
		ctx.Context,
		logger,
		ctx.String("tink-client-id"),
		ctx.String("tink-client-secret"),
		ctx.String("callback-url"),
		db,
	)
	pb.RegisterTinkGWServer(grpcServer, tinkgwServer)

	return util.ListenGRPC(grpcServer, logger, ctx.Int("port"))
}
