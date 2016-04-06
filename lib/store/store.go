package store

import (
	"time"
)

type TimeTracker struct {
	CreatedAt time.Time `json:"created-at" db:"created_at"`
	UpdatedAt time.Time `json:"updated-at" db:"updated_at"`
}

func (tracker *TimeTracker) TouchTimes() {
	if tracker.CreatedAt.IsZero() {
		tracker.CreatedAt = time.Now()
		tracker.UpdatedAt = tracker.CreatedAt
	} else {
		tracker.UpdatedAt = time.Now()
	}
}
