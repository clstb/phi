package server

import (
	"fmt"
	"net"
	"net/http"

	"github.com/clstb/phi/pkg/interceptor"
	"github.com/clstb/phi/pkg/pb"
	"github.com/clstb/phi/pkg/tinkgw"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/soheilhy/cmux"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

func TinkGW(ctx *cli.Context) error {
	var unaryInterceptors []grpc.UnaryServerInterceptor
	var streamInterceptors []grpc.StreamServerInterceptor

	if ctx.String("auth-host") != "" {
		conn, err := grpc.Dial(ctx.String("auth-host"), grpc.WithInsecure())
		if err != nil {
			return err
		}
		auth := pb.NewAuthClient(conn)

		unaryInterceptors = append(
			unaryInterceptors,
			interceptor.AuthUnary(auth),
		)
		streamInterceptors = append(
			streamInterceptors,
			interceptor.AuthStream(auth),
		)
	}

	conn, err := grpc.Dial(ctx.String("core-host"), grpc.WithInsecure())
	if err != nil {
		return err
	}
	core := pb.NewCoreClient(conn)

	// create db connection
	dbStr := ctx.String("db")
	db, err := pgxpool.Connect(ctx.Context, dbStr)
	if err != nil {
		return err
	}
	defer db.Close()
	if err := db.Ping(ctx.Context); err != nil {
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
	server, err := tinkgw.New(
		ctx.String("tink-client-id"),
		ctx.String("tink-client-secret"),
		ctx.String("callback-url"),
		db,
		core,
	)
	if err != nil {
		return err
	}
	pb.RegisterTinkGWServer(s, server)

	// listen and serve
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
		return s.Serve(grpcListener)
	})
	g.Go(func() error {
		return http.Serve(httpListener, server)
	})
	g.Go(func() error {
		return mux.Serve()
	})

	logger.Info(
		"grpc listening",
		zap.String("host", "localhost"),
		zap.Int("port", port),
	)
	logger.Info(
		"http listening",
		zap.String("host", "localhost"),
		zap.Int("port", port),
	)

	return g.Wait()
}
