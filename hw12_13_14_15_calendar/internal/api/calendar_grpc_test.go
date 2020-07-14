package api_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes/empty"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/fdully/hw_otus/hw12_13_14_15_calendar/internal/calendar"

	"github.com/fdully/hw_otus/hw12_13_14_15_calendar/internal/api"
	"github.com/fdully/hw_otus/hw12_13_14_15_calendar/internal/pb"
	"github.com/fdully/hw_otus/hw12_13_14_15_calendar/internal/repository/memory"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCalendarGRPCApi(t *testing.T) {
	ctx := context.Background()

	m := memory.NewRepo()
	a := api.NewCalendarGRPCApi(m)

	event := calendar.CreateEvent(uuid.New(), "Test Subject", "Test Description", time.Now(), time.Now().Add(5*time.Minute), "Test Owner", time.Minute)
	pbEvent := api.MakeProtoEvent(&event)
	require.NotNil(t, pbEvent)

	err := m.AddEvent(ctx, event)
	require.NoError(t, err)

	t.Run("get event", func(t *testing.T) {
		req := &pb.GetEventRequest{Id: event.ID.String()}
		resp, err := a.GetEvent(ctx, req)
		require.NoError(t, err)

		e := resp.GetEvent()
		require.Equal(t, e, pbEvent)
	})

	t.Run("empty id event request", func(t *testing.T) {
		req := &pb.GetEventRequest{Id: ""}
		_, err := a.GetEvent(ctx, req)
		require.Equal(t, err, status.Error(codes.InvalidArgument, "id is not set"))
	})

	t.Run("delete event", func(t *testing.T) {
		req := &pb.DeleteEventRequest{Id: event.ID.String()}
		resp, err := a.DeleteEvent(ctx, req)
		require.NoError(t, err)

		require.Equal(t, resp, &empty.Empty{})
	})

	t.Run("create event", func(t *testing.T) {
		req := &pb.CreateEventRequest{Event: pbEvent}
		resp, err := a.CreateEvent(ctx, req)
		require.NoError(t, err)

		id := resp.GetId()
		require.Equal(t, id, event.ID.String())
	})

	t.Run("update event", func(t *testing.T) {
		event := calendar.CreateEvent(uuid.New(), "Test Subject", "Test Description", time.Now(), time.Now().Add(5*time.Minute), "Test Owner", time.Minute)
		pbEvent := api.MakeProtoEvent(&event)
		require.NotNil(t, pbEvent)
		req := &pb.UpdateEventRequest{Event: pbEvent}

		resp, err := a.UpdateEvent(ctx, req)
		require.NoError(t, err)

		id := resp.GetId()
		require.Equal(t, id, event.ID.String())

		event.Start = event.Start.Add(-2 * time.Minute)
		event.End = event.End.Add(1 * time.Minute)
		event.Subject = "New Subject"

		pbEvent = api.MakeProtoEvent(&event)
		require.NotNil(t, pbEvent)
		req = &pb.UpdateEventRequest{Event: pbEvent}

		resp, err = a.UpdateEvent(ctx, req)
		require.NoError(t, err)

		getReq := &pb.GetEventRequest{Id: event.ID.String()}
		getResp, err := a.GetEvent(ctx, getReq)
		require.NoError(t, err)

		e := getResp.GetEvent()
		require.Equal(t, e, pbEvent)
	})

	t.Run("get events for period", func(t *testing.T) {
		req := &empty.Empty{}
		resp, err := a.GetEventsForToday(ctx, req)
		require.NoError(t, err)

		events := resp.GetEvents()
		require.Len(t, events, 2)

		resp, err = a.GetEventsForWeek(ctx, req)
		require.NoError(t, err)
		events = resp.GetEvents()
		require.Len(t, events, 2)

		resp, err = a.GetEventsForMonth(ctx, req)
		require.NoError(t, err)
		events = resp.GetEvents()
		require.Len(t, events, 2)
	})

}
