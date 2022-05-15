package main

import (
	"github.com/clstb/phi/ledger/internal/server"
	pb "github.com/clstb/phi/proto"
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

	server := grpc.NewServer()
	pb.RegisterBeanAccountServiceServer(server, s)
	reflection.Register(server)
	if err = server.Serve(listener); err != nil {
		return err
	}
	return nil
}
