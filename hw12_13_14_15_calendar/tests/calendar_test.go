// +build integration

package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/fdully/hw_otus/hw12_13_14_15_calendar/internal/logging"

	"github.com/fdully/hw_otus/hw12_13_14_15_calendar/internal/calendar/model"
	"github.com/fdully/hw_otus/hw12_13_14_15_calendar/internal/rabbit"
	"github.com/gofrs/uuid"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/require"
)

type ID struct {
	ID uuid.UUID `json:"id"`
}

func TestCalendar(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := logging.InitLog(-1, "")
	require.NoError(t, err)

	var ev = `{
	"event": {
	"ID": "03fb24ea-3a81-4469-8522-7753d643dcfe",
	"subject": "test_subj",
	"description": "test_desc",
	"start": "2020-08-25T21:00:36.966Z",
	"end": "2021-05-25T23:00:36.966Z",
	"OwnerID": "test_user",
	"notifyPeriod": "300s"
	}
	}`

	res, err := http.Post("http://calendar:8080/api/v1/event/create", "application/json", bytes.NewReader([]byte(ev)))
	require.NoError(t, err)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)

	var id ID
	err = json.Unmarshal(body, &id)
	require.NoError(t, err)
	require.Equal(t, id.ID.String(), "03fb24ea-3a81-4469-8522-7753d643dcfe")

	q := rabbit.NewConnector(ctx, os.Getenv("CAL_Q_URL"), "", os.Getenv("CAL_Q_QUEUE"), 1)
	defer q.Close()

	go func(ctx context.Context, cancelFunc context.CancelFunc) {
		time.Sleep(60 * time.Second)
		cancel()
	}(ctx, cancel)

	var (
		msg amqp.Delivery
		ok  bool
	)
	select {
	case <-ctx.Done():
		cancel()
		t.Fail()
		return
	case msg, ok = <-q.Subscribe(ctx):
		if !ok {
			cancel()
			t.Fail()
			return
		}
	}

	var event model.Event
	err = json.Unmarshal(msg.Body, &event)
	require.NoError(t, err)
	require.Equal(t, id.ID.String(), event.ID.String())

	res, err = http.Post("http://calendar:8080/api/v1/event/delete", "application/json", bytes.NewReader(body))
	require.NoError(t, err)
}
