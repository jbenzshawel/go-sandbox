package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID            int32
	UUID          uuid.UUID
	FirstName     string
	LastName      string
	Email         string
	Enabled       bool
	CreatedAt     time.Time
	LastUpdatedAt time.Time
}
