package main

import (
	"log"
	"os"

	"github.com/clstb/phi/go/pkg/util"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name: "Phi TinkGW Server",
		Commands: []*cli.Command{
			{
				Name:   "server",
				Action: Server,
			},
			{
				Name:   "migrate",
				Action: Migrate,
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "signing-secret",
				EnvVars: []string{"SIGNING_SECRET"},
			},
			&cli.StringFlag{
				Name:    "tink-client-id",
				EnvVars: []string{"TINK_CLIENT_ID"},
			},
			&cli.StringFlag{
				Name:    "tink-client-secret",
				EnvVars: []string{"TINK_CLIENT_SECRET"},
			},
			&cli.StringFlag{
				Name:    "callback-url",
				Value:   "localhost:9000/callback",
				EnvVars: []string{"CALLBACK_URL"},
			},
			&cli.IntFlag{
				Name:    "port",
				Value:   8080,
				EnvVars: []string{"PORT"},
			},
			&cli.StringFlag{
				Name:     "database-url",
				Required: true,
				EnvVars:  []string{"DATABASE_URL"},
			},
		},
		Before: util.Chain(
			util.GetDB,
			util.GetLogger,
		),
		After: util.Chain(
			util.CloseDB,
			util.SyncLogger,
		),
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
