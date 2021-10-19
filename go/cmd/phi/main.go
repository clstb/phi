package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name: "Phi",
		Commands: []*cli.Command{
			{
				Name:   "link",
				Action: Link,
			},
			{
				Name:   "login",
				Action: Login,
			},
			{
				Name:   "register",
				Action: Register,
			},
			{
				Name:   "sync",
				Action: Sync,
			},
			{
				Name:   "categorize",
				Action: Categorize,
			},
			{
				Name:   "store",
				Action: Store,
				Flags: []cli.Flag{
					&cli.PathFlag{
						Name:     "ledger",
						Required: true,
					},
				},
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "auth-host",
				EnvVars: []string{"AUTH_HOST"},
				Value:   "localhost:9000",
			},
			&cli.StringFlag{
				Name:    "tinkgw-host",
				EnvVars: []string{"TINKGW_HOST"},
				Value:   "localhost:9001",
			},
			&cli.StringFlag{
				Name:    "bookkeeper-host",
				EnvVars: []string{"BOOKKEEPER_HOST"},
				Value:   "localhost:9002",
			},
			&cli.PathFlag{
				Name:  "config",
				Value: os.Getenv("HOME") + "/.config/phi.yaml",
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
