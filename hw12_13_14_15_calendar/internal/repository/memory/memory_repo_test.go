package memory

import (
	"context"
	"testing"
	"time"

	"github.com/fdully/hw_otus/hw12_13_14_15_calendar/internal/calendar"

	"github.com/fdully/hw_otus/hw12_13_14_15_calendar/internal/calendar/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestMemRepo(t *testing.T) {
	ctx := context.Background()

	e1 := calendar.CreateEvent(uuid.New(), "Test Subject", "Test Description", time.Now(), time.Now().Add(5*time.Minute), "Test Owner", time.Minute)

	e2 := e1
	e2.ID = uuid.New()

	repo := NewRepo()

	t.Run("add event and remove event", func(t *testing.T) {
		err := repo.AddEvent(ctx, e1)
		require.NoError(t, err)

		err = repo.AddEvent(ctx, e2)
		require.NoError(t, err)

		event, err := repo.GetEvent(ctx, e1.ID)
		require.NoError(t, err)
		require.Equal(t, e1, *event)

		_ = repo.DeleteEvent(ctx, e1.ID)
		_ = repo.DeleteEvent(ctx, e2.ID)

		event, err = repo.GetEvent(ctx, e1.ID)
		require.EqualError(t, err, model.ErrNotExist.Error())
		require.Nil(t, event)

	})

	t.Run("events inside period", func(t *testing.T) {
		err := repo.AddEvent(ctx, e1)
		require.NoError(t, err)

		events, err := repo.GetEventsForPeriod(ctx, time.Now().Add(-1*time.Minute), time.Now().Add(time.Minute*1))
		require.NoError(t, err)
		require.Equal(t, 1, len(events))

		e2.Start = time.Now().Add(-3 * time.Minute)
		e2.End = time.Now().Add(2 * time.Minute)
		err = repo.AddEvent(ctx, e2)
		require.NoError(t, err)

		events, err = repo.GetEventsForPeriod(ctx, time.Now().Add(1*time.Minute), time.Now().Add(time.Minute*2))
		require.NoError(t, err)
		require.Equal(t, 2, len(events))

		_ = repo.DeleteEvent(ctx, e1.ID)
		_ = repo.DeleteEvent(ctx, e2.ID)

	})

	t.Run("events outside of period", func(t *testing.T) {
		err := repo.AddEvent(ctx, e1)
		require.NoError(t, err)

		events, err := repo.GetEventsForPeriod(ctx, time.Now().Add(6*time.Minute), time.Now().Add(time.Minute*6))
		require.NoError(t, err)
		require.Equal(t, 0, len(events))

		events, err = repo.GetEventsForPeriod(ctx, time.Now().Add(-6*time.Minute), time.Now().Add(time.Minute*-4))
		require.NoError(t, err)
		require.Equal(t, 0, len(events))

		_ = repo.DeleteEvent(ctx, e1.ID)
	})
}
