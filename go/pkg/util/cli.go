package util

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"net"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func Chain(funcs ...func(*cli.Context) error) func(ctx *cli.Context) error {
	return func(ctx *cli.Context) error {
		for _, f := range funcs {
			if err := f(ctx); err != nil {
				return err
			}
		}
		return nil
	}
}

func GetDB(ctx *cli.Context) error {
	dbConfig, err := pgx.ParseConfig(ctx.String("database-url"))
	if err != nil {
		return err
	}

	db := stdlib.OpenDB(*dbConfig)
	if err := db.Ping(); err != nil {
		return err
	}

	ctx.Context = context.WithValue(ctx.Context, "db", db)
	return nil
}

func CloseDB(ctx *cli.Context) error {
	db, ok := ctx.Context.Value("db").(*sql.DB)
	if !ok {
		return nil
	}

	return db.Close()
}

func GetLogger(ctx *cli.Context) error {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"}
	logger, err := config.Build()
	if err != nil {
		return err
	}

	ctx.Context = context.WithValue(ctx.Context, "logger", logger)
	return nil
}

func SyncLogger(ctx *cli.Context) error {
	logger, ok := ctx.Context.Value("logger").(*zap.Logger)
	if !ok {
		return nil
	}
	logger.Sync()

	return nil
}

func ListenGRPC(server *grpc.Server, logger *zap.Logger, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	logger.Info(
		"grpc listening",
		zap.String("host", "localhost"),
		zap.Int("port", port),
	)

	return server.Serve(lis)
}

func Migrate(
	db *sql.DB,
	migrations fs.FS,
	down bool,
) error {
	sourceDriver, err := iofs.New(migrations, "schema")
	if err != nil {
		return err
	}

	dbDriver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	mg, err := migrate.NewWithInstance("iofs", sourceDriver, "postgres", dbDriver)
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
