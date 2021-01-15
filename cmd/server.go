package cmd

import (
	"database/sql"
	"net"

	"github.com/clstb/phi/pkg/pb"
	"github.com/clstb/phi/pkg/server"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	_ "github.com/lib/pq"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func Server(ctx *cli.Context) error {
	// create logger
	logger, err := zap.NewProduction()
	if err != nil {
		return err
	}
	defer logger.Sync()

	// create grpc server
	s := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_zap.StreamServerInterceptor(logger),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_zap.UnaryServerInterceptor(logger),
		)),
	)

	// create db connection
	db, err := sql.Open(
		"postgres",
		"postgres://phi@127.0.0.1:26257/phi?sslmode=disable",
	)
	if err != nil {
		return err
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		return err
	}

	// create / register server implementation
	server, err := server.NewServer(
		server.WithDB(db),
	)
	if err != nil {
		return err
	}
	pb.RegisterCoreServer(s, server.Core)

	// listen and serve
	lis, err := net.Listen("tcp", "localhost:9000")
	if err != nil {
		return err
	}

	logger.Info(
		"grpc listening",
		zap.String("host", "localhost"),
		zap.Int("port", 8080),
	)

	return s.Serve(lis)
}
