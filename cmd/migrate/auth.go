package migrate

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/cockroachdb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/urfave/cli/v2"
)

func Auth(ctx *cli.Context) error {
	down := ctx.Bool("down")

	dbStr := ctx.String("db")
	m, err := migrate.New(
		"file://sql/schema/auth",
		"crdb-"+dbStr,
	)
	if err != nil {
		return err
	}

	if down {
		return m.Down()
	}

	err = m.Up()
	if err != migrate.ErrNoChange {
		return err
	}

	return nil
}
