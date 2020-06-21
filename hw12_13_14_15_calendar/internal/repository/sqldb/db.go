package sqldb

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

func OpenDB(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	conn, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("error on connecting to sql db %w", err)
	}
	return conn, nil
}
