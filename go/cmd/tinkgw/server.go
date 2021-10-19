package main

import (
	"fmt"
	"net"
	"net/http"

	"github.com/clstb/phi/go/pkg/interceptor"
	"github.com/clstb/phi/go/pkg/services/tinkgw/pb"
	"github.com/clstb/phi/go/pkg/services/tinkgw/server"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/soheilhy/cmux"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

func Server(ctx *cli.Context) error {
	logger, err := zap.NewProduction()
	if err != nil {
		return err
	}
	defer logger.Sync()

	streamMW := []grpc.StreamServerInterceptor{
		grpc_zap.StreamServerInterceptor(
			logger,
			grpc_zap.WithLevels(grpc_zap.DefaultCodeToLevel),
		),
	}
	unaryMW := []grpc.UnaryServerInterceptor{
		grpc_zap.UnaryServerInterceptor(
			logger,
			grpc_zap.WithLevels(grpc_zap.DefaultCodeToLevel),
		),
	}

	signingSecret := ctx.String("signing-secret")
	if ctx.String("signing-secret") != "" {
		streamMW = append(streamMW,
			interceptor.ServerAuthStream([]byte(signingSecret)),
		)
		unaryMW = append(unaryMW,
			interceptor.ServerAuthUnary([]byte(signingSecret)),
		)
	}

	grpc_zap.ReplaceGrpcLoggerV2(logger)
	s := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(streamMW...)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(unaryMW...)),
	)
	server, err := server.New(
		logger,
		ctx.String("tink-client-id"),
		ctx.String("tink-client-secret"),
		ctx.String("callback-url"),
	)
	if err != nil {
		return err
	}
	pb.RegisterTinkGWServer(s, server)

	port := ctx.Int("port")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	mux := cmux.New(lis)
	grpcListener := mux.MatchWithWriters(
		cmux.HTTP2MatchHeaderFieldSendSettings(
			"content-type",
			"application/grpc",
		),
	)
	httpListener := mux.Match(cmux.Any())

	g := errgroup.Group{}
	g.Go(func() error {
		logger.Info("listening grpc", zap.Int("port", port))
		return s.Serve(grpcListener)
	})
	g.Go(func() error {
		logger.Info("listening http", zap.Int("port", port))
		return http.Serve(httpListener, server)
	})
	g.Go(func() error {
		return mux.Serve()
	})

	return g.Wait()
}
