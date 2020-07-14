package calendar

import (
	"context"
	"time"

	"github.com/fdully/hw_otus/hw12_13_14_15_calendar/internal/calendar/model"
	"github.com/google/uuid"
)

type Repository interface {
	AddEvent(ctx context.Context, event model.Event) error
	UpdateEvent(ctx context.Context, event model.Event) error
	DeleteEvent(ctx context.Context, eventID uuid.UUID) error
	GetEventsForPeriod(ctx context.Context, start, end time.Time) ([]*model.Event, error)
	GetEvent(ctx context.Context, id uuid.UUID) (*model.Event, error)
}

func CreateEvent(id uuid.UUID, subj, desc string, start, end time.Time, ownerID string, notifyPeriod time.Duration) model.Event {
	return model.Event{
		ID:           id,
		Subject:      subj,
		Description:  desc,
		Start:        start,
		End:          end,
		OwnerID:      ownerID,
		NotifyPeriod: notifyPeriod,
	}
}
