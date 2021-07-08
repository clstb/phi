package tinkgw

import (
	"context"
	"net/http"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

func TX(db *pgxpool.Pool) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			tx, err := db.BeginTx(ctx, pgx.TxOptions{})
			if err != nil {
				http.Error(w, "starting transaction failed", http.StatusInternalServerError)
				return
			}
			ctx = context.WithValue(ctx, "tx", tx)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)

			if err := tx.Commit(ctx); err != nil {
				http.Error(w, "commiting transaction failed", http.StatusInternalServerError)
				return
			}
		})
	}
}
