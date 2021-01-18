package server

import (
	"database/sql"
	"fmt"
	"net"
	"strings"

	"github.com/clstb/phi/pkg/core"
	"github.com/clstb/phi/pkg/interceptor"
	"github.com/clstb/phi/pkg/pb"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	_ "github.com/lib/pq"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func Core(ctx *cli.Context) error {
	var unaryInterceptors []grpc.UnaryServerInterceptor
	var streamInterceptors []grpc.StreamServerInterceptor

	authServer := ctx.String("auth-server")
	if authServer != "" {
		conn, err := grpc.Dial(authServer, grpc.WithInsecure())
		if err != nil {
			return err
		}
		authClient := pb.NewAuthClient(conn)

		unaryInterceptors = append(
			unaryInterceptors,
			interceptor.AuthUnary(authClient),
		)
		streamInterceptors = append(
			streamInterceptors,
			interceptor.AuthStream(authClient),
		)
	}

	// create logger
	logger, err := zap.NewProduction()
	if err != nil {
		return err
	}
	defer logger.Sync()
	unaryInterceptors = append(
		unaryInterceptors,
		grpc_zap.UnaryServerInterceptor(logger),
	)
	streamInterceptors = append(
		streamInterceptors,
		grpc_zap.StreamServerInterceptor(logger),
	)

	// create grpc server
	s := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			unaryInterceptors...,
		)),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			streamInterceptors...,
		)),
	)

	// create db connection
	dbStr := ctx.String("db")
	if strings.HasPrefix(dbStr, "cockroachdb") {
		dbStr = "postgres" + strings.TrimPrefix(dbStr, "cockroachdb")
	}
	db, err := sql.Open(
		"postgres",
		dbStr,
	)
	if err != nil {
		return err
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		return err
	}

	// create / register server implementation
	server := core.New(
		core.WithDB(db),
	)
	if err != nil {
		return err
	}
	pb.RegisterCoreServer(s, server)

	// listen and serve
	port := ctx.Int("port")
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		return err
	}

	logger.Info(
		"grpc listening",
		zap.String("host", "localhost"),
		zap.Int("port", port),
	)

	return s.Serve(lis)
}
