package main

import (
	"database/sql"
	"fmt"

	authdb "github.com/clstb/phi/go/pkg/auth/db"
	"github.com/clstb/phi/go/pkg/util"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/urfave/cli/v2"
)

func Migrate(ctx *cli.Context) error {
	db, ok := ctx.Context.Value("db").(*sql.DB)
	if !ok {
		return fmt.Errorf("missing db")
	}
	return util.Migrate(db, authdb.Migrations, ctx.Bool("down"))
}
