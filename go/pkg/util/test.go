package util

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

func TestDB(
	ctx context.Context,
	url,
	name string,
) (*pgxpool.Pool, error) {
	dbConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, err
	}

	db, err := pgxpool.ConnectConfig(ctx, dbConfig)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(ctx); err != nil {
		return nil, err
	}

	_, err = db.Exec(ctx, "DROP DATABASE IF EXISTS "+name)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(ctx, "CREATE DATABASE "+name)
	if err != nil {
		return nil, err
	}

	db.Close()

	dbConfig.ConnConfig.Database = name
	db, err = pgxpool.ConnectConfig(ctx, dbConfig)
	if err != nil {
		return nil, err
	}

	return db, nil
}
