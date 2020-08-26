package model

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/stretchr/testify/require"
)

func TestEvent(t *testing.T) {
	j := `{"id":"03fb24ea-3a81-4469-8522-7753d643dcfe","subject":"test subject","description":"test description",` +
		`"start":"2020-04-04T11:02:30+03:00","end":"2020-04-04T12:02:30+03:00","owner_id":"test.user",` +
		`"notify_period":300000000000}`

	start, _ := time.Parse(time.RFC3339, "2020-04-04T11:02:30+03:00")
	end, _ := time.Parse(time.RFC3339, "2020-04-04T12:02:30+03:00")
	id, _ := uuid.Parse("03fb24ea-3a81-4469-8522-7753d643dcfe")

	e := Event{
		ID:           id,
		Subject:      "test subject",
		Description:  "test description",
		Start:        start,
		End:          end,
		OwnerID:      "test.user",
		NotifyPeriod: 300 * time.Second,
	}

	var ee Event

	err := json.Unmarshal([]byte(j), &ee)
	require.NoError(t, err)
	require.Equal(t, e, ee)

	b, err := json.Marshal(ee)
	require.NoError(t, err)
	require.Equal(t, j, string(b))
}
