package main

import (
	"log"
	"os"
	"time"

	"github.com/clstb/phi/cmd"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:        "Phi",
		Description: "Phi - Personal finance management",
		Commands: []*cli.Command{
			{
				Name:        "server",
				Description: "Runs the Phi server",
				Action:      cmd.Server,
			},
			{
				Name:        "migrate",
				Description: "Runs database migrations",
				Action:      cmd.Migrate,
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "down",
						Usage:   "toggle down migrations",
						Aliases: []string{"d"},
					},
				},
			},
			{
				Name:        "ingest",
				Description: "Parses and ingests the provided csv file",
				Action:      cmd.Ingest,
				Flags: []cli.Flag{
					&cli.PathFlag{
						Name:      "file",
						Aliases:   []string{"f"},
						Usage:     "csv file",
						Required:  true,
						TakesFile: true,
					},
					&cli.BoolFlag{
						Name:  "skip-duplicates",
						Usage: "skips duplicates based on matching hash values",
						Value: true,
					},
				},
			},
			{
				Name:        "create",
				Description: "Creates resources",
				Aliases:     []string{"c"},
				Subcommands: []*cli.Command{
					{

						Name:        "account",
						Description: "Creates a new account",
						Aliases:     []string{"a"},
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "name",
								Aliases:  []string{"n"},
								Usage:    "account name",
								Required: true,
							},
						},
						Action: cmd.CreateAccount,
					},
				},
			},
			{
				Name:        "register",
				Description: "Register a new ledger",
				Action:      cmd.Register,
			},
			{
				Name:        "balances",
				Description: "Prints trial balances",
				Action:      cmd.Balances,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "date",
						Aliases: []string{"d"},
						Usage:   "Date of trial balance",
						Value:   time.Now().Format("2006-01-02"),
					},
				},
			},
			{
				Name:        "income",
				Description: "Prints income statement for period",
				Action:      cmd.Income,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "from",
						Aliases: []string{"f"},
						Usage:   "Period start",
						Value:   "0001-01-01",
					},
					&cli.StringFlag{
						Name:    "to",
						Aliases: []string{"t"},
						Usage:   "Period end",
						Value:   time.Now().Format("2006-01-02"),
					},
				},
			},
			{
				Name:        "balsheet",
				Description: "Prints balance sheet",
				Action:      cmd.BalSheet,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "date",
						Aliases: []string{"d"},
						Usage:   "Date of balance sheet",
						Value:   time.Now().Format("2006-01-02"),
					},
				},
			},
			{
				Name:        "journal",
				Description: "Prints journal",
				Action:      cmd.Journal,
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "api-host",
				Usage: "phi server host",
				Value: "localhost:9000",
			},
		}}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
