package domain

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (r *MockUserRepository) InsertUser(user User) error {
	args := r.Called(user)
	return args.Error(0)
}

func (r *MockUserRepository) GetUserByEmail(email string) (*User, error) {
	args := r.Called(email)
	return args.Get(0).(*User), args.Error(1)
}

func (r *MockUserRepository) GetUserByUUID(uuid uuid.UUID) (*User, error) {
	args := r.Called(uuid)
	return args.Get(0).(*User), args.Error(1)
}
