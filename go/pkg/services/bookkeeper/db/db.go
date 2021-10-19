// Code generated by sqlc. DO NOT EDIT.

package db

import (
	"context"
	"database/sql"
	"fmt"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

func Prepare(ctx context.Context, db DBTX) (*Queries, error) {
	q := Queries{db: db}
	var err error
	if q.createLedgerStmt, err = db.PrepareContext(ctx, createLedger); err != nil {
		return nil, fmt.Errorf("error preparing query CreateLedger: %w", err)
	}
	if q.getLedgerStmt, err = db.PrepareContext(ctx, getLedger); err != nil {
		return nil, fmt.Errorf("error preparing query GetLedger: %w", err)
	}
	if q.updateLedgerStmt, err = db.PrepareContext(ctx, updateLedger); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateLedger: %w", err)
	}
	return &q, nil
}

func (q *Queries) Close() error {
	var err error
	if q.createLedgerStmt != nil {
		if cerr := q.createLedgerStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createLedgerStmt: %w", cerr)
		}
	}
	if q.getLedgerStmt != nil {
		if cerr := q.getLedgerStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getLedgerStmt: %w", cerr)
		}
	}
	if q.updateLedgerStmt != nil {
		if cerr := q.updateLedgerStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateLedgerStmt: %w", cerr)
		}
	}
	return err
}

func (q *Queries) exec(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (sql.Result, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).ExecContext(ctx, args...)
	case stmt != nil:
		return stmt.ExecContext(ctx, args...)
	default:
		return q.db.ExecContext(ctx, query, args...)
	}
}

func (q *Queries) query(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (*sql.Rows, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryContext(ctx, args...)
	default:
		return q.db.QueryContext(ctx, query, args...)
	}
}

func (q *Queries) queryRow(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) *sql.Row {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryRowContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryRowContext(ctx, args...)
	default:
		return q.db.QueryRowContext(ctx, query, args...)
	}
}

type Queries struct {
	db               DBTX
	tx               *sql.Tx
	createLedgerStmt *sql.Stmt
	getLedgerStmt    *sql.Stmt
	updateLedgerStmt *sql.Stmt
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db:               tx,
		tx:               tx,
		createLedgerStmt: q.createLedgerStmt,
		getLedgerStmt:    q.getLedgerStmt,
		updateLedgerStmt: q.updateLedgerStmt,
	}
}
