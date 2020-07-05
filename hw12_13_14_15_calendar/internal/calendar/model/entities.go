package model

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID           uuid.UUID
	Subject      string
	Description  string
	Start        time.Time
	End          time.Time
	OwnerID      string
	NotifyPeriod time.Duration
}
