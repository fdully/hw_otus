package memory

import (
	"context"
	"sync"
	"time"

	"github.com/fdully/hw_otus/hw12_13_14_15_calendar/internal/calendar"
	"github.com/fdully/hw_otus/hw12_13_14_15_calendar/internal/calendar/model"
	"github.com/fdully/hw_otus/hw12_13_14_15_calendar/internal/util/timeutil"
	"github.com/google/uuid"
)

var _ calendar.Repository = (*repo)(nil)

func NewRepo() calendar.Repository {
	return &repo{
		mu: sync.Mutex{},
		s:  make(map[uuid.UUID]*model.Event),
	}
}

type repo struct {
	mu sync.Mutex
	s  map[uuid.UUID]*model.Event
}

func (r *repo) AddEvent(ctx context.Context, e model.Event) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.s[e.ID]; ok {
		return model.ErrAlreadyExist
	}

	r.s[e.ID] = &e

	return nil
}

func (r *repo) UpdateEvent(ctx context.Context, e model.Event) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.s[e.ID]; ok {
		r.s[e.ID].ID = e.ID
		r.s[e.ID].Subject = e.Subject
		r.s[e.ID].Description = e.Description
		r.s[e.ID].Start = e.Start
		r.s[e.ID].End = e.End
		r.s[e.ID].OwnerID = e.OwnerID
		r.s[e.ID].NotifyPeriod = e.NotifyPeriod

		return nil
	}

	r.s[e.ID] = &e

	return nil
}

func (r *repo) DeleteEvent(ctx context.Context, eventID uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.s, eventID)

	return nil
}

func (r *repo) GetEvent(ctx context.Context, id uuid.UUID) (*model.Event, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	event, ok := r.s[id]
	if !ok {
		return nil, model.ErrNotExist
	}

	return event, nil
}

func (r *repo) GetEventsForPeriod(ctx context.Context, start, end time.Time) ([]*model.Event, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var results = make([]*model.Event, 0, len(r.s))
	for _, v := range r.s {
		// is start or end of event inside of requested range period?
		if timeutil.TimeInRange(v.Start, start, end) || timeutil.TimeInRange(v.End, start, end) || (v.Start.Before(start) && v.End.After(end)) {
			results = append(results, v)
		}
	}

	return results, nil
}
