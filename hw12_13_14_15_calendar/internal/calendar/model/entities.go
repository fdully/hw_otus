package model

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID           uuid.UUID     `json:"id"`
	Subject      string        `json:"subject"`
	Description  string        `json:"description"`
	Start        time.Time     `json:"start"`
	End          time.Time     `json:"end"`
	OwnerID      string        `json:"owner_id"`
	NotifyPeriod time.Duration `json:"notify_period"`
}
