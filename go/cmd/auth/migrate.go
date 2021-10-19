package main

import (
	"github.com/clstb/phi/go/pkg/services/auth/db"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/urfave/cli/v2"
)

func Migrate(ctx *cli.Context) error {
	down := ctx.Bool("down")

	d, err := iofs.New(db.Migrations, "schema")
	if err != nil {
		return err
	}

	mg, err := migrate.NewWithSourceInstance("iofs", d, ctx.String("db"))
	if err != nil {
		return err
	}

	if down {
		return mg.Down()
	}

	err = mg.Up()
	if err != migrate.ErrNoChange {
		return err
	}

	return nil
}
