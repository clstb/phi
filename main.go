package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/clstb/phi/cmd"
	"github.com/clstb/phi/cmd/create"
	"github.com/clstb/phi/cmd/csv"
	"github.com/clstb/phi/cmd/migrate"
	"github.com/clstb/phi/cmd/server"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:        "Phi",
		Description: "Phi - Personal finance management",
		Commands: []*cli.Command{
			{
				Name:  "server",
				Usage: "runs the phi server",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:    "port",
						Usage:   "port to serve on",
						Aliases: []string{"p"},
						Value:   9000,
					},
				},
				Subcommands: []*cli.Command{
					{
						Name:   "core",
						Usage:  "starts the phi core microservice",
						Action: server.Core,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "db",
								Usage:   "postgres compatible database connection string",
								EnvVars: []string{"DB"},
							},
						},
					},
					{
						Name:   "auth",
						Usage:  "starts the phi auth microservice",
						Action: server.Auth,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "db",
								Usage:   "postgres compatible database connection string",
								EnvVars: []string{"DB"},
							},
							&cli.StringFlag{
								Name:     "signing-secret",
								Usage:    "secret to sign jwt's with",
								EnvVars:  []string{"SIGNING_SECRET"},
								Required: true,
							},
						},
					},
				},
			},
			{
				Name:  "migrate",
				Usage: "runs database migrations",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "down",
						Usage:   "toggle between up and down migrations",
						Aliases: []string{"d"},
					},
					&cli.StringFlag{
						Name:    "db",
						Usage:   "postgres compatible database connection string",
						EnvVars: []string{"DB"},
					},
				},
				Subcommands: []*cli.Command{
					{
						Name:   "core",
						Usage:  "runs migrations for the phi core microservice",
						Action: migrate.Core,
					},
					{
						Name:   "auth",
						Usage:  "runs migrations for the phi auth microservice",
						Action: migrate.Auth,
					},
				},
			},
			{
				Name: "csv",
				Subcommands: []*cli.Command{
					{
						Name:   "parse",
						Action: csv.Parse,
						Flags: []cli.Flag{
							&cli.PathFlag{
								Name:      "file",
								Aliases:   []string{"f"},
								Required:  true,
								TakesFile: true,
							},
							&cli.PathFlag{
								Name:    "output",
								Aliases: []string{"o"},
								Value:   "./parsed.csv",
							},
						},
					},
					{
						Name:   "review",
						Action: csv.Review,
						Flags: []cli.Flag{
							&cli.PathFlag{
								Name:      "file",
								Aliases:   []string{"f"},
								Required:  true,
								TakesFile: true,
							},
						},
					},
					{
						Name:   "ingest",
						Action: csv.Ingest,
						Flags: []cli.Flag{
							&cli.PathFlag{
								Name:      "file",
								Aliases:   []string{"f"},
								Required:  true,
								TakesFile: true,
							},
						},
					},
				},
			},
			{
				Name:    "create",
				Usage:   "create various resources",
				Aliases: []string{"c"},
				Subcommands: []*cli.Command{
					{

						Name:    "account",
						Usage:   "creates an account",
						Aliases: []string{"a"},
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "name",
								Aliases:  []string{"n"},
								Usage:    "account name",
								Required: true,
							},
						},
						Action: create.Account,
					},
					{
						Name:    "transaction",
						Usage:   "create a transaction",
						Aliases: []string{"t"},
						Action:  create.Transaction,
					},
				},
			},
			{
				Name:   "register",
				Usage:  "register as new user",
				Action: cmd.Register,
			},
			{
				Name:   "login",
				Usage:  "login as user",
				Action: cmd.Login,
			},
			{
				Name:   "balances",
				Usage:  "print trial balance for date",
				Action: cmd.Balances,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "date",
						Aliases: []string{"d"},
						Usage:   "date of trial balance",
						Value:   time.Now().Format("2006-01-02"),
					},
				},
			},
			{
				Name:   "income",
				Usage:  "print income statement for period",
				Action: cmd.Income,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "from",
						Aliases: []string{"f"},
						Usage:   "period start",
						Value:   "0001-01-01",
					},
					&cli.StringFlag{
						Name:    "to",
						Aliases: []string{"t"},
						Usage:   "period end",
						Value:   time.Now().Format("2006-01-02"),
					},
				},
			},
			{
				Name:   "balsheet",
				Usage:  "print balance sheet for date",
				Action: cmd.BalSheet,
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
				Name:   "journal",
				Usage:  "print journal",
				Action: cmd.Journal,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "from",
						Aliases: []string{"f"},
						Usage:   "period start",
						Value:   "0001-01-01",
					},
					&cli.StringFlag{
						Name:    "to",
						Aliases: []string{"t"},
						Usage:   "period end",
						Value:   time.Now().Format("2006-01-02"),
					},
				},
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "core-host",
				Usage:   "phi core server host",
				Value:   "localhost:9000",
				EnvVars: []string{"CORE_HOST"},
			},
			&cli.StringFlag{
				Name:    "auth-host",
				Usage:   "phi auth server host",
				Value:   "localhost:9000",
				EnvVars: []string{"AUTH_HOST"},
			},
			&cli.StringFlag{
				Name:  "config",
				Usage: "phi client config path",
				Value: fmt.Sprintf("%s/.config/phi.yaml", os.Getenv("HOME")),
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
