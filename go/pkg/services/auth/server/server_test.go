//go:build integration
// +build integration

package server_test

import (
	"database/sql"
	"log"
	"os"
	"testing"

	db "github.com/clstb/phi/go/pkg/services/auth/db"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/jackc/pgx/v4/stdlib"
)

var testdb *sql.DB

func TestMain(m *testing.M) {
	var err error

	testdb, err = sql.Open("pgx", os.Getenv("PHI_AUTH_DB"))
	if err != nil {
		log.Fatal(err)
	}

	sourceDriver, err := iofs.New(db.Migrations, "schema")
	if err != nil {
		log.Fatal(err)
	}

	dbDriver, err := postgres.WithInstance(testdb, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	mg, err := migrate.NewWithInstance("iofs", sourceDriver, "postgres", dbDriver)
	if err != nil {
		log.Fatal(err)
	}

	if err := mg.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	status := m.Run()

	if err := mg.Down(); err != nil {
		log.Fatal(err)
	}

	os.Exit(status)
}
