package user

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (r *MockUserRepository) Create(user *User) error {
	args := r.Called(user)
	return args.Error(0)
}

func (r *MockUserRepository) Update(user *User) error {
	args := r.Called(user)
	return args.Error(0)
}

func (r *MockUserRepository) GetAll(page, pageSize int) ([]*User, error) {
	args := r.Called(page, pageSize)
	return args.Get(0).([]*User), args.Error(1)
}

func (r *MockUserRepository) GetByEmail(email string) (*User, error) {
	args := r.Called(email)
	return args.Get(0).(*User), args.Error(1)
}

func (r *MockUserRepository) GetByUUID(uuid uuid.UUID) (*User, error) {
	args := r.Called(uuid)
	return args.Get(0).(*User), args.Error(1)
}
