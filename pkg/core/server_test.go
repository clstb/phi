// +build integration

package core_test

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/cockroachdb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

var sqlDB *sql.DB

func TestMain(m *testing.M) {
	dbStr := os.Getenv("DB")
	mg, err := migrate.New(
		"file://./../../sql/schema/core",
		"crdb-"+dbStr,
	)
	if err != nil {
		log.Fatal(err)
	}

	err = mg.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	sqlDB, err = sql.Open(
		"postgres",
		dbStr,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()

	if err := sqlDB.Ping(); err != nil {
		log.Fatal(err)
	}

	status := m.Run()

	err = mg.Down()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	os.Exit(status)
}
