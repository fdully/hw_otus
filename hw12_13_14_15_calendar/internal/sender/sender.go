package sender

import (
	"context"
	"encoding/json"
	"time"

	"github.com/fdully/hw_otus/hw12_13_14_15_calendar/internal/calendar/model"
	"github.com/fdully/hw_otus/hw12_13_14_15_calendar/internal/logging"
	"github.com/fdully/hw_otus/hw12_13_14_15_calendar/internal/rabbit"
)

type Sender struct {
	subscriber rabbit.Subscriber
}

func NewSender(subscriber rabbit.Subscriber) *Sender {
	return &Sender{subscriber: subscriber}
}

func (s *Sender) Run(ctx context.Context) error {
	logger := logging.FromContext(ctx)

	msgCh := s.subscriber.Subscribe(ctx)
	for {
		select {
		case <-ctx.Done():
			return nil
		case msg, ok := <-msgCh:
			if !ok {
				return nil
			}

			var event model.Event
			err := json.Unmarshal(msg.Body, &event)
			if err != nil {
				logger.Error(err)
			}

			until := time.Until(event.Start)
			logger.Infof("event %s is %f minutes", event.Subject, until.Minutes())
		}
	}
}
