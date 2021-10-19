package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name: "TinkGW",
		Commands: []*cli.Command{
			{
				Name: "server",
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
						Value:   9000,
						EnvVars: []string{"PORT"},
					},
				},
				Action: Server,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}
