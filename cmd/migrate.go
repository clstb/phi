package cmd

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/cockroachdb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/urfave/cli/v2"
)

func Migrate(ctx *cli.Context) error {
	down := ctx.Bool("down")

	m, err := migrate.New(
		"file://sql/schema",
		"cockroachdb://phi@localhost:26257/phi?sslmode=disable",
	)
	if err != nil {
		return err
	}

	if down {
		return m.Down()
	}
	return m.Up()
}
