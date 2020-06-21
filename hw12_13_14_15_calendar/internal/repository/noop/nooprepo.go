package noop

import (
	"context"

	"github.com/fdully/hw_otus/hw12_13_14_15_calendar/internal/calendar"
	"github.com/fdully/hw_otus/hw12_13_14_15_calendar/internal/calendar/model"
	"github.com/google/uuid"
)

var _ calendar.Repository = (*Repo)(nil)

type Repo struct{}

func (r Repo) AddEvent(ctx context.Context, event *model.Event) error {
	return nil
}

func (r Repo) AlterEvent(ctx context.Context, event *model.Event) error {
	return nil
}

func (r Repo) DeleteEvent(ctx context.Context, eventID uuid.UUID) error {
	return nil
}

func (Repo) GetEvents(ctx context.Context) ([]*model.Event, error) {
	return nil, nil
}
