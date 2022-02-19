package main

import (
	"log"
	"os"

	"github.com/clstb/phi/go/pkg/util"
	cli "github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name: "Phi Auth Server",
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
			&cli.IntFlag{
				Name:    "port",
				Value:   8080,
				EnvVars: []string{"PORT"},
			},
			&cli.StringFlag{
				Name:     "signing-secret",
				Required: true,
				EnvVars:  []string{"SIGNING_SECRET"},
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
