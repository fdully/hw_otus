package memory

import (
	"context"
	"sync"

	"github.com/fdully/hw_otus/hw12_13_14_15_calendar/internal/calendar"
	"github.com/fdully/hw_otus/hw12_13_14_15_calendar/internal/calendar/model"
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

func (r *repo) AddEvent(ctx context.Context, e *model.Event) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.s[e.ID] = e

	return nil
}

func (r *repo) AlterEvent(ctx context.Context, e *model.Event) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.s[e.ID].Subject = e.Subject
	r.s[e.ID].OwnerID = e.OwnerID
	r.s[e.ID].Description = e.Description
	r.s[e.ID].Start = e.Start
	r.s[e.ID].End = e.End
	r.s[e.ID].NotifyPeriod = e.NotifyPeriod

	return nil
}

func (r *repo) DeleteEvent(ctx context.Context, eventID uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.s, eventID)

	return nil
}

func (r *repo) GetEvents(ctx context.Context) ([]*model.Event, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var results = make([]*model.Event, 0, len(r.s))
	for _, v := range r.s {
		results = append(results, v)
	}

	return results, nil
}
