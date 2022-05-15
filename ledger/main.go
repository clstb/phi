package main

import (
	"github.com/clstb/phi/ledger/internal/server"
	pb "github.com/clstb/phi/proto"
	grpczap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
)

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:  "port",
				Value: 8082,
			},
		},
		Action: run,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(ctx *cli.Context) error {
	addr := "0.0.0.0:" + ctx.String("port")
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	s := server.NewServer()

	s.Logger.Info("----> GRPC listeninng on %s", addr)

	_server := grpc.NewServer(
		grpc.StreamInterceptor(
			grpczap.StreamServerInterceptor(s.Logger.Desugar())),
		grpc.UnaryInterceptor(grpczap.UnaryServerInterceptor(s.Logger.Desugar())),
	)
	pb.RegisterBeanAccountServiceServer(_server, s)
	reflection.Register(_server)
	if err = _server.Serve(listener); err != nil {
		return err
	}
	return nil
}
