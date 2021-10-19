package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name: "Phi Auth Server",
		Commands: []*cli.Command{
			{
				Name:   "server",
				Action: Server,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "signing-secret",
						Required: true,
						EnvVars:  []string{"SIGNING_SECRET"},
					},
					&cli.IntFlag{
						Name:    "port",
						Value:   9000,
						EnvVars: []string{"PORT"},
					},
					&cli.IntFlag{
						Name:    "gateway-port",
						Value:   9090,
						EnvVars: []string{"GATEWAY_PORT"},
					},
				},
			},
			{
				Name:   "migrate",
				Action: Migrate,
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "db",
				Required: true,
				EnvVars:  []string{"DB"},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
