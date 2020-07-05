package memory

import (
	"context"
	"github.com/fdully/hw_otus/hw12_13_14_15_calendar/internal/calendar/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestMemRepo(t *testing.T) {
	ctx := context.Background()

	e := &model.Event{
		ID:          uuid.New(),
		Subject:     "Test Subject",
		Description: "Test Description",
		Start:       time.Now(),
		End:         time.Now().Add(5 * time.Minute),
		OwnerID:     "Test Owner",
	}

	repo := NewRepo()

	err := repo.AddEvent(ctx, e)
	require.NoError(t, err)

	events, err := repo.GetEvents(ctx)
	require.NoError(t, err)
	require.Equal(t, 1, len(events))
	require.Equal(t, e, events[0])

	err = repo.DeleteEvent(ctx, e.ID)
	require.NoError(t, err)
}
