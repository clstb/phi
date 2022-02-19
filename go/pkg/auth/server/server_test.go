package server_test

import (
	"log"
	"os"
	"testing"

	authdb "github.com/clstb/phi/go/pkg/auth/db"
	"github.com/clstb/phi/go/pkg/util"
)

func TestMain(t *testing.M) {
	db, err := util.TestDB(os.Getenv("DATABASE_URL"), "phi_auth")
	if err != nil {
		log.Fatal(err)
	}

	if err := util.Migrate(db, authdb.Migrations, false); err != nil {
		log.Fatal(err)
	}

	status := t.Run()

	if err := util.Migrate(db, authdb.Migrations, true); err != nil {
		log.Fatal(err)
	}

	os.Exit(status)
}
