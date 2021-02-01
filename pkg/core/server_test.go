// +build integration

package core_test

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var sqlDB *sql.DB

func TestMain(m *testing.M) {
	var err error
	sqlDB, err = sql.Open(
		"postgres",
		"postgresql://root@localhost:26257?sslmode=disable",
	)
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()

	if err := sqlDB.Ping(); err != nil {
		log.Fatal(err)
	}

	status := m.Run()
	os.Exit(status)
}
