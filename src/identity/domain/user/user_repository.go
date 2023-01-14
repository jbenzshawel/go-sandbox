package user

import (
	"github.com/google/uuid"
)

type Repository interface {
	GetUserByEmail(email string) (*User, error)
	GetUserByUUID(uuid uuid.UUID) (*User, error)
	InsertUser(user *User) error
	UpdateUser(user *User) error
}
