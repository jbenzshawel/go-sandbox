package idp

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/jbenzshawel/go-sandbox/identity/domain"
)

type MockIdentityProvider struct {
	mock.Mock
}

func (m *MockIdentityProvider) CreateUser(ctx context.Context, user domain.User, password string) (uuid.UUID, error) {
	args := m.Called(ctx, user, password)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *MockIdentityProvider) DeleteUser(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}
