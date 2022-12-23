package domain

import "github.com/google/uuid"

type UserRepository interface {
	GetUserByEmail(email string) (*User, error)
	GetUserByUUID(uuid uuid.UUID) (*User, error)
	InsertUser(user User) error
}
