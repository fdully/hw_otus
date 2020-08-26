package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/fdully/hw_otus/hw12_13_14_15_calendar/internal/calendar"
	"github.com/fdully/hw_otus/hw12_13_14_15_calendar/internal/logging"
	"github.com/fdully/hw_otus/hw12_13_14_15_calendar/internal/rabbit"
	"github.com/fdully/hw_otus/hw12_13_14_15_calendar/internal/util/timeutil"
)

type Scheduler struct {
	repo      calendar.Repository
	publisher rabbit.Publisher
}

func NewScheduler(repo calendar.Repository, publisher rabbit.Publisher) *Scheduler {
	return &Scheduler{repo: repo, publisher: publisher}
}

func (s Scheduler) notify(ctx context.Context) error {
	now := time.Now()
	events, err := s.repo.GetEventsForPeriod(ctx, timeutil.Bod(now), timeutil.Eow(now))
	if err != nil {
		return err
	}

	for _, e := range events {
		if timeutil.IsNotifyStarted(e) {
			msg, err := json.Marshal(e)
			if err != nil {
				return fmt.Errorf("failed to marshal event: %w", err)
			}

			err = s.publisher.Publish(msg)
			if err != nil {
				return fmt.Errorf("failed to publish event: %w", err)
			}
		}
	}

	return nil
}

func (s Scheduler) removeEvents(ctx context.Context) error {
	now := time.Now()
	start := now.AddDate(-10, 0, 0)
	end := now.AddDate(-1, 0, 0)

	events, err := s.repo.GetEventsForPeriod(ctx, start, end)
	if err != nil {
		return err
	}

	for _, e := range events {
		err = s.repo.DeleteEvent(ctx, e.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Scheduler) Run(ctx context.Context) error {
	logger := logging.FromContext(ctx)

	go func(ctx context.Context) {
		t := time.NewTicker(1 * time.Hour)
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				if err := s.removeEvents(ctx); err != nil {
					logger.Errorf("failed to remove old events: %w", err)
				}
			}
		}
	}(ctx)

	t := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-t.C:
			if err := s.notify(ctx); err != nil {
				logger.Errorf("failed to notify about events: %w", err)
			}
		}
	}
}
