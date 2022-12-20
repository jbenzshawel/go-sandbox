package domain

import "github.com/google/uuid"

type UserRepository interface {
	CreateUser(user User, password string) error
	GetUserByEmail(email string) (*User, error)
	GetUserByUUID(uuid uuid.UUID) (*User, error)
}
