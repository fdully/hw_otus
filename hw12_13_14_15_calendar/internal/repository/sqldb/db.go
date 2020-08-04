package sqldb

import (
	"context"
	"fmt"

	"github.com/fdully/hw_otus/hw12_13_14_15_calendar/internal/config"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

func OpenDB(ctx context.Context) (*pgxpool.Pool, error) {
	conf := config.FromContext(ctx)

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		conf.SQLDB.Host, conf.SQLDB.Port, conf.SQLDB.Login, conf.SQLDB.Password, conf.SQLDB.Database)

	conn, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("error on connecting to sql db %w", err)
	}
	return conn, nil
}

// InTx runs the given function f within a transaction with isolation level isoLevel.
func (r Repo) InTx(ctx context.Context, isoLevel pgx.TxIsoLevel, f func(tx pgx.Tx) error) error {
	conn, err := r.Pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("acquiring connection: %v", err)
	}
	defer conn.Release()

	tx, err := conn.BeginTx(ctx, pgx.TxOptions{IsoLevel: isoLevel})
	if err != nil {
		return fmt.Errorf("starting transaction: %v", err)
	}

	if err := f(tx); err != nil {
		if err1 := tx.Rollback(ctx); err1 != nil {
			return fmt.Errorf("rolling back transaction: %v (original error: %v)", err1, err)
		}
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("committing transaction: %v", err)
	}
	return nil
}
