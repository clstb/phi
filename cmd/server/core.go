package server

import (
	"context"
	"fmt"
	"net"

	"github.com/clstb/phi/pkg/core"
	"github.com/clstb/phi/pkg/interceptor"
	"github.com/clstb/phi/pkg/pb"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/jackc/pgtype"
	shopspring "github.com/jackc/pgtype/ext/shopspring-numeric"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func Core(ctx *cli.Context) error {
	var unaryInterceptors []grpc.UnaryServerInterceptor
	var streamInterceptors []grpc.StreamServerInterceptor

	authServer := ctx.String("auth-host")
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

	// create db connection
	dbStr := ctx.String("db")
	connConfig, err := pgxpool.ParseConfig(dbStr)
	connConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		conn.ConnInfo().RegisterDataType(pgtype.DataType{
			Value: &shopspring.Numeric{},
			Name:  "numeric",
			OID:   pgtype.NumericOID,
		})
		return nil
	}
	if err != nil {
		return err
	}

	db, err := pgxpool.ConnectConfig(ctx.Context, connConfig)
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
	server := core.New(
		core.WithDB(db),
	)
	if err != nil {
		return err
	}
	pb.RegisterCoreServer(s, server)

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
