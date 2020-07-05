package sqldb

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/fdully/hw_otus/hw12_13_14_15_calendar/internal/calendar"
	"github.com/fdully/hw_otus/hw12_13_14_15_calendar/internal/calendar/model"
	pgx "github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

var _ calendar.Repository = (*Repo)(nil)

type Repo struct {
	Pool *pgxpool.Pool
}

func (r Repo) AddEvent(ctx context.Context, e *model.Event) error {
	conn, err := r.Pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("acquiring connection: %w", err)
	}
	defer conn.Release()

	result, err := conn.Exec(ctx, `
		INSERT INTO events
			(id, subject, start_time, end_time, description, user_id)
		VALUES
			($1, $2, $3, $4, $5, $6)
		`, e.ID, e.Subject, e.Start, e.End, e.Description, e.OwnerID)
	if err != nil {
		return fmt.Errorf("inserting event: %w", err)
	}

	if result.RowsAffected() != 1 {
		return fmt.Errorf("no rows updated")
	}
	return nil
}

func (r Repo) AlterEvent(ctx context.Context, e *model.Event) error {
	conn, err := r.Pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("acquiring connection: %w", err)
	}
	defer conn.Release()

	result, err := conn.Exec(ctx, `
		UPDATE events
		SET
			subject = $1, start_time = $2, end_time = $3,
			description = $4, user_id = $5, notify_time = $6
		WHERE
			id = $7
		`, e.Subject, e.Start, e.End, e.Description, e.OwnerID, e.NotifyPeriod, e.ID)
	if err != nil {
		return fmt.Errorf("updating event: %w", err)
	}

	if result.RowsAffected() != 1 {
		return fmt.Errorf("no rows updated")
	}
	return nil
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

func (r Repo) GetEvents(ctx context.Context) ([]*model.Event, error) {
	conn, err := r.Pool.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("acquiring connection: %w", err)
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, `
		SELECT
			id, subject, start_time, end_time, description, user_id, notify_time
		FROM
			events`)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	var results []*model.Event
	for rows.Next() {
		e, err := scanOneEvent(rows)
		if err != nil {
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
		return nil, err
	}
	return &m, nil
}
