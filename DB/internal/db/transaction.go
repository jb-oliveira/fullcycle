package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func ExecuteInTransaction(ctx context.Context, conn *pgx.Conn, fn func(queries *Queries) error) error {
	tx, err := conn.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return err
	}
	queries := New(conn).WithTx(tx)
	err = fn(queries)
	if err != nil {
		errRollback := tx.Rollback(ctx)
		if errRollback != nil {
			return fmt.Errorf("rollback failed: %w; original error: %w", errRollback, err)
		}
		return err
	}
	return tx.Commit(ctx)
}
