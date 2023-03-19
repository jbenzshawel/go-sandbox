package user

import (
	"github.com/google/uuid"
)

type Repository interface {
	GetAll(page, pageSize int) ([]*User, error)
	GetByEmail(email string) (*User, error)
	GetByUUID(uuid uuid.UUID) (*User, error)
	Create(user *User) error
	Update(user *User) error
}
