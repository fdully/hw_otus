package api

import (
	"context"
	"errors"
	"time"

	"github.com/fdully/hw_otus/hw12_13_14_15_calendar/internal/calendar"
	"github.com/fdully/hw_otus/hw12_13_14_15_calendar/internal/calendar/model"
	"github.com/fdully/hw_otus/hw12_13_14_15_calendar/internal/logging"
	"github.com/fdully/hw_otus/hw12_13_14_15_calendar/internal/pb"
	"github.com/fdully/hw_otus/hw12_13_14_15_calendar/internal/util/timeutil"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CalGRPCApi struct {
	cc *calendar.Calendar
}

func NewCalendarGRPCApi(cc *calendar.Calendar) CalGRPCApi {
	return CalGRPCApi{cc}
}

func (c CalGRPCApi) CreateEvent(ctx context.Context, req *pb.CreateEventRequest) (*pb.CreateEventResponse, error) {
	logger := logging.FromContext(ctx)
	ev := req.GetEvent()

	if ev == nil {
		logger.Info("event is nil")

		return nil, status.Error(codes.InvalidArgument, "event must be provided")
	}

	event, err := validateRequestAndCreateEvent(ev)
	if err != nil {
		return nil, err
	}

	// adding new event to repository
	err = c.cc.AddEvent(ctx, event)
	if err != nil {
		code := codes.Internal
		if errors.Is(err, model.ErrAlreadyExist) {
			code = codes.AlreadyExists
		}

		logger.Error("adding new event to storage", err)

		return nil, status.Errorf(code, "cannot add new event: %v", err)
	}

	res := &pb.CreateEventResponse{Id: ev.ID}

	return res, nil
}

func (c CalGRPCApi) UpdateEvent(ctx context.Context, req *pb.UpdateEventRequest) (*pb.UpdateEventResponse, error) {
	logger := logging.FromContext(ctx)
	ev := req.GetEvent()
	if ev == nil {
		logger.Info("event is nil")

		return nil, status.Error(codes.InvalidArgument, "event must be provided")
	}

	event, err := validateRequestAndCreateEvent(ev)
	if err != nil {
		return nil, err
	}

	// upserting event in repository
	err = c.cc.UpdateEvent(ctx, event)
	if err != nil {
		logger.Error("upserting event in storage: ", err)

		return nil, status.Errorf(codes.Internal, "can't upsert event: %v", err)
	}

	res := &pb.UpdateEventResponse{Id: ev.ID}

	return res, nil
}

func (c CalGRPCApi) DeleteEvent(ctx context.Context, req *pb.DeleteEventRequest) (*empty.Empty, error) {
	logger := logging.FromContext(ctx)
	reqID := req.GetId()

	if reqID == "" {
		return nil, status.Errorf(codes.InvalidArgument, "id is not set")
	}

	id, err := uuid.Parse(reqID)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "id is not a valid UUID: %v", err)
	}

	err = c.cc.DeleteEvent(ctx, id)
	if err != nil {
		logger.Errorf("can't delete event %v", err)

		return nil, status.Errorf(codes.Internal, "can't delete event %v", err)
	}

	return &empty.Empty{}, nil
}

func (c CalGRPCApi) GetEvent(ctx context.Context, req *pb.GetEventRequest) (*pb.GetEventResponse, error) {
	logger := logging.FromContext(ctx)

	reqID := req.GetId()
	if reqID == "" {
		return nil, status.Errorf(codes.InvalidArgument, "id is not set")
	}
	id, err := uuid.Parse(reqID)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "id is not a valid UUID: %v", err)
	}

	event, err := c.cc.GetEvent(ctx, id)
	if err != nil {
		logger.Errorf("can't get event from repository %v", err)

		return nil, status.Errorf(codes.Internal, "can't get event %v", err)
	}

	if event == nil {
		return &pb.GetEventResponse{}, nil
	}

	start, err := ptypes.TimestampProto(event.Start)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "incorrect time format in event %v", err)
	}
	end, err := ptypes.TimestampProto(event.End)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "incorrect time format in event %v", err)
	}

	return &pb.GetEventResponse{
		Event: &pb.Event{
			ID:           event.ID.String(),
			Subject:      event.Subject,
			Description:  event.Description,
			Start:        start,
			End:          end,
			OwnerID:      event.OwnerID,
			NotifyPeriod: ptypes.DurationProto(event.NotifyPeriod),
		},
	}, nil
}

func (c CalGRPCApi) GetEventsForToday(ctx context.Context, req *pb.GetEventsForPeriodRequest) (*pb.GetEventsResponse, error) {
	logger := logging.FromContext(ctx)

	now, err := ptypes.Timestamp(req.GetSearchWithTime())
	if err != nil {
		logger.Errorf("bad request search time: %v", err)

		return nil, status.Error(codes.InvalidArgument, "bad search time provided")
	}

	start := timeutil.Bod(now)
	end := timeutil.Eod(now)

	resEvents, err := c.getEventsForPeriod(ctx, start, end)
	if err != nil {
		return nil, err
	}

	return &pb.GetEventsResponse{Events: resEvents}, nil
}

func (c CalGRPCApi) GetEventsForWeek(ctx context.Context, req *pb.GetEventsForPeriodRequest) (*pb.GetEventsResponse, error) {
	logger := logging.FromContext(ctx)

	now, err := ptypes.Timestamp(req.GetSearchWithTime())
	if err != nil {
		logger.Errorf("bad request search time: %v", err)

		return nil, status.Error(codes.InvalidArgument, "bad search time provided")
	}

	start := timeutil.Bow(now)
	end := timeutil.Eow(now)

	resEvents, err := c.getEventsForPeriod(ctx, start, end)
	if err != nil {
		return nil, err
	}

	return &pb.GetEventsResponse{Events: resEvents}, nil
}

func (c CalGRPCApi) GetEventsForMonth(ctx context.Context, req *pb.GetEventsForPeriodRequest) (*pb.GetEventsResponse, error) {
	logger := logging.FromContext(ctx)

	now, err := ptypes.Timestamp(req.GetSearchWithTime())
	if err != nil {
		logger.Errorf("bad request search time: %v", err)

		return nil, status.Error(codes.InvalidArgument, "bad search time provided")
	}

	start := timeutil.Bom(now)
	end := timeutil.Eom(now)

	resEvents, err := c.getEventsForPeriod(ctx, start, end)
	if err != nil {
		return nil, err
	}

	return &pb.GetEventsResponse{Events: resEvents}, nil
}

func MakeProtoEvent(e *model.Event) *pb.Event {
	if e == nil {
		return nil
	}
	logger := logging.FromContext(context.Background())
	start, err := ptypes.TimestampProto(e.Start)
	if err != nil {
		logger.Errorf("can't convert start time to proto timestamp event with id %s, %v", e.ID.String(), err)

		return nil
	}
	end, err := ptypes.TimestampProto(e.End)
	if err != nil {
		logger.Errorf("can't convert end time to proto timestamp event with id %s, %v", e.ID.String(), err)

		return nil
	}

	return &pb.Event{
		ID:           e.ID.String(),
		Subject:      e.Subject,
		Description:  e.Description,
		Start:        start,
		End:          end,
		OwnerID:      e.OwnerID,
		NotifyPeriod: ptypes.DurationProto(e.NotifyPeriod),
	}
}

func validateRequestAndCreateEvent(ev *pb.Event) (model.Event, error) {
	var event model.Event
	var id uuid.UUID
	if len(ev.ID) > 0 {
		// check if it's a valid UUID
		var err error
		id, err = uuid.Parse(ev.ID)
		if err != nil {
			return event, status.Errorf(codes.InvalidArgument, "Event ID is not a valid UUID: %v", err)
		}
	}
	start, err := ptypes.Timestamp(ev.Start)
	if err != nil {
		return event, status.Errorf(codes.InvalidArgument, "Start time of the event is not valid: %v", err)
	}
	end, err := ptypes.Timestamp(ev.End)
	if err != nil {
		return event, status.Errorf(codes.InvalidArgument, "End time of the event is not valid: %v", err)
	}
	notifyPeriod, err := ptypes.Duration(ev.NotifyPeriod)
	if err != nil {
		return event, status.Errorf(codes.InvalidArgument, "Notify Period of the event is not valid: %v", err)
	}

	event = calendar.CreateEvent(id, ev.Subject, ev.Description, start, end, ev.OwnerID, notifyPeriod)

	return event, nil
}

func (c CalGRPCApi) getEventsForPeriod(ctx context.Context, start, end time.Time) ([]*pb.Event, error) {
	logger := logging.FromContext(ctx)
	events, err := c.cc.GetEventsForPeriod(ctx, start, end)
	if err != nil {
		logger.Errorf("can't get events from repository %v", err)

		return nil, status.Errorf(codes.Internal, "can't get events %v", err)
	}

	if len(events) == 0 {
		return []*pb.Event{}, nil
	}

	var resEvents = make([]*pb.Event, 0, len(events))
	for _, v := range events {
		e := MakeProtoEvent(v)
		if e == nil {
			continue
		}
		resEvents = append(resEvents, e)
	}

	return resEvents, nil
}
