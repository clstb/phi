package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/clstb/phi/go/internal/tinkgw/server"
	"github.com/urfave/cli/v2"
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
				Name:     "ory-token",
				EnvVars:  []string{"ORY_TOKEN"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "jwks-url",
				EnvVars:  []string{"JWKS_URL"},
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
		log.Fatal(err)
	}
}

func run(ctx *cli.Context) error {
	server, err := server.NewServer(
		ctx.String("tink-client-id"),
		ctx.String("tink-client-secret"),
		ctx.String("ory-token"),
		ctx.String("jwks-url"),
		ctx.String("callback-url"),
	)
	if err != nil {
		return err
	}

	return http.ListenAndServe(fmt.Sprintf(":%d", ctx.Int("port")), server)
}
