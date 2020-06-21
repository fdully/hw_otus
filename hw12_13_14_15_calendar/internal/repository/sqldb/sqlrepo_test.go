package sqldb

import (
	"context"
	"fmt"
	"github.com/fdully/hw_otus/hw12_13_14_15_calendar/internal/calendar/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

func TestSqlDB(t *testing.T) {
	ctx := context.Background()
	db, err := OpenDB(ctx, fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable", os.Getenv("CAL_DB_HOST"), os.Getenv("CAL_DB_PORT"),
		os.Getenv("CAL_DB_LOGIN"), os.Getenv("CAL_DB_PASSWORD"), os.Getenv("CAL_DB_NAME")))

	require.NoError(t, err)

	e := &model.Event{
		ID:          uuid.New(),
		Subject:     "Test Subject",
		Description: "Test Description",
		Start:       time.Now(),
		End:         time.Now().Add(5 * time.Minute),
		OwnerID:     "Test Owner",
	}

	repo := Repo{Pool: db}

	err = repo.AddEvent(ctx, e)
	require.NoError(t, err)

	err = repo.DeleteEvent(ctx, e.ID)
	require.NoError(t, err)

}
