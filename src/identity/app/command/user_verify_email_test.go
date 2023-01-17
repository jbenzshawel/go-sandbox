package command

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/jbenzshawel/go-sandbox/identity/domain/token"
	"github.com/jbenzshawel/go-sandbox/identity/domain/user"
	"github.com/jbenzshawel/go-sandbox/identity/infrastructure/idp"
	"github.com/jbenzshawel/go-sandbox/identity/infrastructure/storage"
)

func TestVerifyEmailHandler(t *testing.T) {
	tokenRepo := storage.NewVerificationTokenRepository(
		storage.NewVerificationTokenCache(),
	)
	userRepo := storage.NewUserMemoryRepository()
	u, err := user.NewUser("First", "Last", "test@email.com", false, false)
	require.NoError(t, err)
	require.NoError(t, u.SetUUID(uuid.New()))
	require.NoError(t, userRepo.InsertUser(u))

	tkn, err := token.NewToken()
	require.NoError(t, err)
	tokenRepo.SaveToken(u.UUID(), tkn)

	mockIDP := &idp.MockIdentityProvider{}
	mockIDP.On("UpdateUser", mock.Anything, mock.Anything).
		Return(nil).
		Once()

	testLogger, _ := test.NewNullLogger()

	handler := NewVerifyEmailHandler(userRepo, tokenRepo, mockIDP, logrus.NewEntry(testLogger))

	require.NoError(t, handler.Handle(context.Background(), VerifyEmail{
		UserId: u.UUID(),
		Code:   tkn.Code(),
	}))

	u, err = userRepo.GetUserByUUID(u.UUID())
	require.NoError(t, err)
	assert.True(t, u.EmailVerified())

	assert.Nil(t, tokenRepo.GetToken(u.UUID()))

	mockIDP.AssertExpectations(t)
}
