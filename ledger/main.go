package main

import (
	"fmt"
	"github.com/clstb/phi/ledger/internal/config"
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
			&cli.StringFlag{
				Name:    "TINK_GW_URI",
				EnvVars: []string{"TINK_GW_URI"},
				Value:   config.DefTinkGwAddr,
			},
			&cli.StringFlag{
				Name:    "DATA_DIR_PATH",
				EnvVars: []string{"DATA_DIR_PATH"},
				Value:   config.DefDataDirPath,
			},
		},
		Action: run,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(ctx *cli.Context) error {

	err := os.Mkdir(ctx.String("DATA_DIR_PATH"), config.DefFilePermissions)

	if err != nil {
		return err
	}

	fmt.Println("------------------------")
	fmt.Printf("TINK_GW_URI => %s\n", ctx.String("TINK_GW_URI"))
	fmt.Printf("DATA_DIR_PATH  => %s\n", ctx.String("DATA_DIR_PATH"))
	fmt.Println("------------------------")

	addr := "0.0.0.0:8082"
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	s := server.NewServer(ctx.String("TINK_GW_URI"), ctx.String("DATA_DIR_PATH"))
	s.Logger.Info("----> GRPC listeninng on ", addr)

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
