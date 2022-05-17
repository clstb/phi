package main

import (
	pb "github.com/clstb/phi/proto"
	"github.com/clstb/phi/tinkgw/internal/server"
	grpczap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	_ "github.com/motemen/go-loghttp/global"
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
			&cli.StringFlag{
				Name:     "CLIENT_ID",
				EnvVars:  []string{"CLIENT_ID"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "CLIENT_SECRET",
				EnvVars:  []string{"CLIENT_SECRET"},
				Required: true,
			},
			&cli.StringFlag{
				Name:    "CALLBACK_URL",
				EnvVars: []string{"CALLBACK_URL"},
				Value:   "http://localhost:9000/callback",
			},
		},
		Action: run,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(ctx *cli.Context) error {
	addr := "0.0.0.0:8080"
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	s := server.NewServer(ctx.String("CLIENT_ID"), ctx.String("CLIENT_SECRET"), ctx.String("CALLBACK_URL"))
	s.Logger.Info("----> GRPC listeninng on ", addr)

	_server := grpc.NewServer(
		grpc.StreamInterceptor(grpczap.StreamServerInterceptor(s.Logger.Desugar())),
		grpc.UnaryInterceptor(grpczap.UnaryServerInterceptor(s.Logger.Desugar())),
	)
	pb.RegisterTinkGWServiceServer(_server, s)
	pb.RegisterTransactionGWServiceServer(_server, s)
	reflection.Register(_server)
	if err = _server.Serve(listener); err != nil {
		return err
	}
	return nil
}
