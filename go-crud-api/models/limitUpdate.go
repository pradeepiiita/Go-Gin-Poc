package models

import (
	"encoding/json"
	"time"
)

type LimitUpdateJson struct {
	ID         string          `json:"id" gorm:"primarykey"`
	Details    json.RawMessage `json:"details"`
	Status     string          `json:"status"`
	StartedAt  time.Time       `json:"started_at"`
	FinishedAt time.Time       `json:"finished_at"`
}
