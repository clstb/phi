package interceptor

import (
	"context"
	"fmt"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
)

func TXUnary(db *pgxpool.Pool) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		tx, err := db.BeginTx(ctx, pgx.TxOptions{})
		if err != nil {
			return nil, err
		}
		ctx = context.WithValue(ctx, "tx", tx)

		res, err := handler(ctx, req)
		if err != nil {
			return nil, fmt.Errorf(
				"err: %w; rollback err: %v",
				err,
				tx.Rollback(ctx),
			)
		}

		if err := tx.Commit(ctx); err != nil {
			return nil, fmt.Errorf("commit err: %w", err)
		}

		return res, nil
	}
}

func TXStream(db *pgxpool.Pool) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		tx, err := db.BeginTx(ss.Context(), pgx.TxOptions{})
		if err != nil {
			return err
		}
		ctx := context.WithValue(ss.Context(), "tx", tx)

		wrapped := grpc_middleware.WrapServerStream(ss)
		wrapped.WrappedContext = ctx

		if err := handler(srv, wrapped); err != nil {
			return fmt.Errorf(
				"err: %w; rollback err: %v",
				err,
				tx.Rollback(ctx),
			)
		}

		if err := tx.Commit(ctx); err != nil {
			return fmt.Errorf("commit err: %w", err)
		}

		return nil
	}
}
