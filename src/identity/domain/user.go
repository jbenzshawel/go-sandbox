package domain

import "time"

type User struct {
	ID            int
	UUID          string
	FirstName     string
	LastName      string
	Email         string
	Enabled       bool
	CreatedAt     time.Time
	LastUpdatedAt time.Time
}
