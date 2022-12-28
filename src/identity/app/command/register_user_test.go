package command

import (
	"context"
	"errors"
	"fmt"
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
	mockIDP.On("CreateUser", mock.Anything, mock.Anything, mock.Anything).
		Return(fakeUserID, nil).
		Once()

	handler := NewRegisterUserHandler(userRepo, mockIDP, logrus.NewEntry(testLogger))

	cmd := RegisterUser{
		FirstName:       "TestFirst",
		LastName:        "TestLast",
		Email:           fmt.Sprintf("%s@test.com", uuid.New().String()),
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
	assert.True(t, user.Enabled)

	mockIDP.AssertExpectations(t)
}

func TestRegisterUserHandler_CreateIDPUserFails(t *testing.T) {
	userRepo := getUserRepo()
	testLogger, _ := test.NewNullLogger()
	mockIDP := &idp.MockIdentityProvider{}

	mockIDP.On("CreateUser", mock.Anything, mock.Anything, mock.Anything).
		Return(uuid.Nil, errors.New("create fails")).
		Once()

	handler := NewRegisterUserHandler(userRepo, mockIDP, logrus.NewEntry(testLogger))

	cmd := RegisterUser{
		FirstName:       "TestFirst",
		LastName:        "TestLast",
		Email:           fmt.Sprintf("%s@test.com", uuid.New().String()),
		Password:        "P@ssw0RD",
		ConfirmPassword: "P@ssw0RD",
	}

	err := handler.Handle(context.Background(), cmd)
	require.Errorf(t, err, "create fails")

	mockIDP.AssertExpectations(t)
}

func TestRegisterUserHandler_CreateIDPUserPartiallyFails(t *testing.T) {
	userRepo := getUserRepo()
	testLogger, _ := test.NewNullLogger()
	mockIDP := &idp.MockIdentityProvider{}

	fakeUserID := uuid.New()
	mockIDP.On("CreateUser", mock.Anything, mock.Anything, mock.Anything).
		Return(fakeUserID, errors.New("create fails")).
		Once()
	mockIDP.On("DeleteUser", mock.Anything, fakeUserID.String()).
		Return(nil).
		Once()

	handler := NewRegisterUserHandler(userRepo, mockIDP, logrus.NewEntry(testLogger))

	cmd := RegisterUser{
		FirstName:       "TestFirst",
		LastName:        "TestLast",
		Email:           fmt.Sprintf("%s@test.com", uuid.New().String()),
		Password:        "P@ssw0RD",
		ConfirmPassword: "P@ssw0RD",
	}

	err := handler.Handle(context.Background(), cmd)
	require.Errorf(t, err, "create fails")

	mockIDP.AssertExpectations(t)
}

func TestRegisterUserHandler_RepoInsertUserFails(t *testing.T) {
	testLogger, _ := test.NewNullLogger()
	fakeUserID := uuid.New()
	fakeEmail := fmt.Sprintf("%s@test.com", uuid.New().String())

	mockIDP := &idp.MockIdentityProvider{}
	mockIDP.On("CreateUser", mock.Anything, mock.Anything, mock.Anything).
		Return(fakeUserID, nil).
		Once()
	mockIDP.On("DeleteUser", mock.Anything, fakeUserID.String()).
		Return(nil).
		Once()

	mockRepo := &domain.MockUserRepository{}
	var user *domain.User
	mockRepo.On("GetUserByEmail", fakeEmail).
		Return(user, nil).
		Once()
	mockRepo.On("InsertUser", mock.Anything).
		Return(errors.New("repo error")).
		Once()

	handler := NewRegisterUserHandler(mockRepo, mockIDP, logrus.NewEntry(testLogger))

	cmd := RegisterUser{
		FirstName:       "TestFirst",
		LastName:        "TestLast",
		Email:           fakeEmail,
		Password:        "P@ssw0RD",
		ConfirmPassword: "P@ssw0RD",
	}

	err := handler.Handle(context.Background(), cmd)
	require.Errorf(t, err, "repo error")

	mockIDP.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func getUserRepo() domain.UserRepository {
	if userSqlRepo, ok := storage.TryCreateUserSqlRepository(); ok {
		return userSqlRepo
	}

	return storage.NewUserMemoryRepository()
}
