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
				Name:    "api-url",
				EnvVars: []string{"API_URL"},
				Value:   "https://phi.clstb.codes",
			},
			&cli.PathFlag{
				Name:  "config",
				Value: os.Getenv("HOME") + "/.config/phi/config.yaml",
			},
			&cli.PathFlag{
				Name:  "ledger",
				Value: os.Getenv("HOME") + "/.config/phi/ledger",
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
