package sqldb

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/fdully/hw_otus/hw12_13_14_15_calendar/internal/calendar"
	"github.com/fdully/hw_otus/hw12_13_14_15_calendar/internal/calendar/model"
	"github.com/google/uuid"
	pgx "github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

var _ calendar.Repository = (*Repo)(nil)

type Repo struct {
	Pool *pgxpool.Pool
}

func (r Repo) AddEvent(ctx context.Context, e model.Event) error {
	conn, err := r.Pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("acquiring connection: %w", err)
	}
	defer conn.Release()

	result, err := conn.Exec(ctx, `
		INSERT INTO events
			(id, subject, start_time, end_time, description, owner_id, notify_period)
		VALUES
			($1, $2, $3, $4, $5, $6, $7)
		`, e.ID, e.Subject, e.Start.Format(time.RFC3339), e.End.Format(time.RFC3339), e.Description, e.OwnerID, int64(e.NotifyPeriod.Seconds()))
	if err != nil {
		return fmt.Errorf("inserting event: %w", err)
	}

	if result.RowsAffected() != 1 {
		return fmt.Errorf("no rows updated")
	}

	return nil
}

func (r Repo) UpdateEvent(ctx context.Context, e model.Event) error {
	return r.InTx(ctx, pgx.Serializable, func(tx pgx.Tx) error {
		query := `
			INSERT INTO
				events
				(id, subject, description, start_time, end_time, owner_id, notify_period)
			VALUES
				($1, $2, $3, $4, $5, $6, $7)
			ON CONFLICT
				(id)
			DO UPDATE
				SET subject = $2, description = $3, start_time = $4, end_time = $5, owner_id = $6, notify_period = $7	
		`
		_, err := tx.Exec(ctx, query, e.ID.String(), e.Subject, e.Description, e.Start.Format(time.RFC3339), e.End.Format(time.RFC3339), e.OwnerID, int64(e.NotifyPeriod.Seconds()))
		if err != nil {
			return fmt.Errorf("can't upsert event: %w", err)
		}

		return nil
	})
}

func (r Repo) DeleteEvent(ctx context.Context, eventID uuid.UUID) error {
	conn, err := r.Pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("acquiring connection: %w", err)
	}
	defer conn.Release()

	result, err := conn.Exec(ctx, `
		DELETE FROM
			events
		WHERE
			id = $1	
	`, eventID)
	if err != nil {
		return fmt.Errorf("deleting event: %w", err)
	}

	if result.RowsAffected() != 1 {
		return fmt.Errorf("no rows updated")
	}

	return nil
}

func (r Repo) GetEvent(ctx context.Context, id uuid.UUID) (*model.Event, error) {
	conn, err := r.Pool.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("acquiring connection: %w", err)
	}
	defer conn.Release()

	row := conn.QueryRow(ctx, `
		SELECT
			id, subject, start_time, end_time, description, owner_id, notify_period
		FROM
			events
	    WHERE
	        id = $1
    `, id)

	var res *model.Event
	if res, err = scanOneEvent(row); err != nil {
		if err == model.ErrNotExist {
			return nil, model.ErrNotExist
		}

		return nil, fmt.Errorf("scanning results: %w", err)
	}

	return res, nil
}

func (r Repo) GetEventsForPeriod(ctx context.Context, start, end time.Time) ([]*model.Event, error) {
	conn, err := r.Pool.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("acquiring connection: %w", err)
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, `
		SELECT
			id, subject, start_time, end_time, description, owner_id, notify_period
		FROM
			events
	    WHERE
	        (start_time BETWEEN $1 AND $2)
		OR
			(end_time BETWEEN $1 AND $2)
        OR
            (start_time <= $1 AND $2 <= end_time)
    `, start.Format(time.RFC3339), end.Format(time.RFC3339))
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	var results []*model.Event
	for rows.Next() {
		e, err := scanOneEvent(rows)
		if err != nil {
			if errors.Is(model.ErrNotExist, err) {
				return nil, model.ErrNotExist
			}

			return nil, fmt.Errorf("scaning rows: %w", err)
		}
		results = append(results, e)
	}

	return results, nil
}

func scanOneEvent(row pgx.Row) (*model.Event, error) {
	var (
		m model.Event
	)
	if err := row.Scan(&m.ID, &m.Subject, &m.Start, &m.End, &m.Description, &m.OwnerID, &m.NotifyPeriod); err != nil {
		if errors.Is(pgx.ErrNoRows, err) {
			return nil, model.ErrNotExist
		}

		return nil, err
	}

	return &m, nil
}
