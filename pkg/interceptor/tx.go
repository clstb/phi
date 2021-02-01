package interceptor

import (
	"context"
	"database/sql"
	"fmt"

	"google.golang.org/grpc"
)

func TXUnary(db *sql.DB) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return nil, err
		}
		ctx = context.WithValue(ctx, "tx", tx)

		res, err := handler(ctx, req)
		if err != nil {
			return nil, fmt.Errorf(
				"err: %w; rollback err: %v",
				err,
				tx.Rollback(),
			)
		}

		if err := tx.Commit(); err != nil {
			return nil, fmt.Errorf("commit err: %w", err)
		}

		return res, nil
	}
}
