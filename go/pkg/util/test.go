package util

import (
	"database/sql"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
)

func TestDB(
	url,
	name string,
) (*sql.DB, error) {
	dbConfig, err := pgx.ParseConfig(url)
	if err != nil {
		return nil, err
	}

	db := stdlib.OpenDB(*dbConfig)
	if err := db.Ping(); err != nil {
		return nil, err
	}

	_, err = db.Exec("DROP DATABASE IF EXISTS " + name)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("CREATE DATABASE " + name)
	if err != nil {
		return nil, err
	}

	if err := db.Close(); err != nil {
		return nil, err
	}

	dbConfig.Database = name

	return stdlib.OpenDB(*dbConfig), nil
}
