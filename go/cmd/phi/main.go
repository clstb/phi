package main

import (
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "auth-host",
				EnvVars: []string{"AUTH_HOST"},
				Value:   "phi-auth.fly.dev:443",
			},
			&cli.StringFlag{
				Name:    "tinkgw-host",
				EnvVars: []string{"TINKGW_HOST"},
				Value:   "phi-tinkgw.fly.dev:443",
			},
			&cli.PathFlag{
				Name:  "config",
				Value: os.Getenv("HOME") + "/.config/phi/config.yaml",
			},
			&cli.PathFlag{
				Name:  "ledger",
				Value: os.Getenv("HOME") + "/.config/phi/ledger",
			},
			&cli.BoolFlag{
				Name:  "insecure",
				Value: false,
			},
		},
		Action: run,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(ctx *cli.Context) error {
	p := tea.NewProgram(newModel(ctx))
	return p.Start()
}
