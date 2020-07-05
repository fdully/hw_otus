package calendar

import (
	"context"

	"github.com/fdully/hw_otus/hw12_13_14_15_calendar/internal/calendar/model"
	"github.com/google/uuid"
)

type Repository interface {
	AddEvent(ctx context.Context, event *model.Event) error
	AlterEvent(ctx context.Context, event *model.Event) error
	DeleteEvent(ctx context.Context, eventID uuid.UUID) error
	GetEvents(ctx context.Context) ([]*model.Event, error)
}

type Calendar struct {
	r Repository
}

func NewCalendar(r Repository) (*Calendar, error) {
	return &Calendar{r: r}, nil
}
