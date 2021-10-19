package main

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"

	"github.com/clstb/phi/go/pkg/interceptor"
	pb "github.com/clstb/phi/go/pkg/services/auth/pb"
	"github.com/clstb/phi/go/pkg/services/auth/server"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

func Server(ctx *cli.Context) error {
	var unaryInterceptors []grpc.UnaryServerInterceptor
	var streamInterceptors []grpc.StreamServerInterceptor

	// create db connection
	db, err := sql.Open("pgx", ctx.String("db"))
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
	server := server.New([]byte(ctx.String("signing-secret")))
	if err != nil {
		return err
	}
	pb.RegisterAuthServer(s, server)

	port, gatewayPort := ctx.Int("port"), ctx.Int("gateway-port")

	g := errgroup.Group{}
	g.Go(func() error {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			return err
		}
		logger.Info("listening grpc", zap.Int("port", port))
		return s.Serve(lis)
	})

	g.Go(func() error {
		mux := runtime.NewServeMux()
		if err := pb.RegisterAuthHandlerFromEndpoint(
			ctx.Context,
			mux,
			fmt.Sprintf("localhost:%d", port),
			[]grpc.DialOption{grpc.WithInsecure()},
		); err != nil {
			return err
		}

		logger.Info("listening http", zap.Int("port", gatewayPort))
		return http.ListenAndServe(fmt.Sprintf(":%d", gatewayPort), mux)
	})

	return g.Wait()
}
