package server

import (
	"database/sql"
	"fmt"
	"net"

	"github.com/clstb/phi/pkg/auth"
	"github.com/clstb/phi/pkg/interceptor"
	"github.com/clstb/phi/pkg/pb"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	_ "github.com/lib/pq"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func Auth(ctx *cli.Context) error {
	var unaryInterceptors []grpc.UnaryServerInterceptor
	var streamInterceptors []grpc.StreamServerInterceptor

	// create db connection
	dbStr := ctx.String("db")
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

	// create logger
	logger, err := zap.NewProduction()
	if err != nil {
		return err
	}
	defer logger.Sync()

	// setup interceptors
	unaryInterceptors = append(unaryInterceptors, []grpc.UnaryServerInterceptor{
		interceptor.TXUnary(db),
		grpc_zap.UnaryServerInterceptor(logger),
	}...)
	streamInterceptors = append(streamInterceptors, []grpc.StreamServerInterceptor{
		interceptor.TXStream(db),
		grpc_zap.StreamServerInterceptor(logger),
	}...)

	// create grpc server
	s := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			unaryInterceptors...,
		)),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			streamInterceptors...,
		)),
	)

	// create / register server implementation
	signingSecret := ctx.String("signing-secret")
	server := auth.New(
		auth.WithDB(db),
		auth.WithSigningSecret([]byte(signingSecret)),
	)
	if err != nil {
		return err
	}
	pb.RegisterAuthServer(s, server)

	// listen and serve
	port := ctx.Int("port")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
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
