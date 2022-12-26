package command

import (
	"context"
	"testing"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/jbenzshawel/go-sandbox/identity/domain"
	"github.com/jbenzshawel/go-sandbox/identity/infrastructure/idp"
	"github.com/jbenzshawel/go-sandbox/identity/infrastructure/storage"
)

func TestRegisterUserHandler(t *testing.T) {
	userRepo := getUserRepo()
	testLogger, _ := test.NewNullLogger()
	mockIDP := &idp.MockIdentityProvider{}

	fakeUserID := uuid.New()
	mockIDP.On("CreateUser", mock.Anything, mock.Anything, mock.Anything).Return(fakeUserID, nil).Once()

	handler := NewRegisterUserHandler(userRepo, mockIDP, logrus.NewEntry(testLogger))

	cmd := RegisterUser{
		FirstName:       "TestFirst",
		LastName:        "TestLast",
		Email:           "test@email.com",
		Password:        "P@ssw0RD",
		ConfirmPassword: "P@ssw0RD",
	}

	err := handler.Handle(context.Background(), cmd)
	require.NoError(t, err)

	user, err := userRepo.GetUserByEmail(cmd.Email)
	require.NoError(t, err)
	require.NotNil(t, user)
	assert.Equal(t, fakeUserID, user.UUID)
	assert.Equal(t, cmd.FirstName, user.FirstName)
	assert.Equal(t, cmd.LastName, user.LastName)
	assert.Equal(t, cmd.Email, user.Email)

	mockIDP.AssertExpectations(t)
}

// TODO: Create tests for error cases

func getUserRepo() domain.UserRepository {
	if userSqlRepo, ok := storage.TryCreateUserSqlRepository(); ok {
		return userSqlRepo
	}

	return storage.NewUserMemoryRepository()
}
