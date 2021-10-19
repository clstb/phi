package main

import (
	"fmt"

	"github.com/clstb/phi/go/pkg/config"
	"github.com/clstb/phi/go/pkg/interceptor"
	"github.com/clstb/phi/go/pkg/services/tinkgw/pb"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
)

func Link(ctx *cli.Context) error {
	configPath := ctx.String("config")
	config, err := config.Load(configPath)
	if err != nil {
		return err
	}

	tinkGWHost := ctx.String("tinkgw-host")
	conn, err := grpc.Dial(
		tinkGWHost,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(interceptor.ClientAuthUnary(config.AccessToken)),
		grpc.WithStreamInterceptor(interceptor.ClientAuthStream(config.AccessToken)),
	)
	if err != nil {
		return err
	}

	client := pb.NewTinkGWClient(conn)

	link, err := client.GetLink(ctx.Context, &pb.GetLinkReq{
		Market: "DE",
		Locale: "de_DE",
	})
	if err != nil {
		return err
	}

	fmt.Println(link.Link)
	return nil
}
