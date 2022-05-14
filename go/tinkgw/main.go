package main

import (
	pb "github.com/clstb/phi/go/proto"
	"github.com/clstb/phi/go/tinkgw/pkg"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
)

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "tink-client-id",
				EnvVars:  []string{"TINK_CLIENT_ID"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "tink-client-secret",
				EnvVars:  []string{"TINK_CLIENT_SECRET"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "callback-url",
				EnvVars:  []string{"CALLBACK_URL"},
				Required: true,
			},
			&cli.IntFlag{
				Name:  "port",
				Value: 8080,
			},
		},
		Action: run,
	}

	if err := app.Run(os.Args); err != nil {
		pkg.Sugar.Fatal(err)
	}
}

type Server struct {
	pb.TinkGWServiceServer
}

func run(ctx *cli.Context) error {
	addr := "0.0.0.0:" + ctx.String("port")
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	pkg.Sugar.Info("----> GRPC listeninng on %s", addr)

	server := grpc.NewServer()
	pb.RegisterTinkGWServiceServer(server, &Server{})
	reflection.Register(server)
	if err = server.Serve(listener); err != nil {
		return err
	}
	return nil
}
