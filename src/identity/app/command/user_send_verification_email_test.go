package command

import (
	"context"
	"fmt"
	"net/url"
	"testing"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmihailenco/msgpack/v5"

	"github.com/jbenzshawel/go-sandbox/common/messaging"
	"github.com/jbenzshawel/go-sandbox/identity/infrastructure/storage"
)

type stubPublisher struct {
	handler func(topic string, msg []byte) error
}

func (p *stubPublisher) Publish(topic string, msg []byte) error {
	if p.handler != nil {
		return p.handler(topic, msg)
	}
	return nil
}

func TestSendVerificationEmailHandler(t *testing.T) {
	tokenRepo := storage.NewVerificationTokenRepository(
		storage.NewVerificationTokenCache(),
	)
	cmd := SendVerificationEmail{
		UserUUID:  uuid.New(),
		FirstName: "Test",
		Email:     "test@email.com",
	}
	publisherCall := 0
	mockPublisher := &stubPublisher{
		handler: func(topic string, msg []byte) error {
			publisherCall++
			assert.Equal(t, messaging.TopicVerifyEmail, topic)

			var payload messaging.VerifyEmail
			require.NoError(t, msgpack.Unmarshal(msg, &payload))
			require.NotNil(t, payload)
			assert.Equal(t, cmd.UserUUID, payload.UserUUID)
			assert.Equal(t, cmd.FirstName, payload.FirstName)
			assert.Equal(t, cmd.Email, payload.Email)

			token := tokenRepo.GetToken(cmd.UserUUID)
			require.NotNil(t, token)
			assert.Equal(t, token.Code(), payload.Code)
			assert.Equal(t, fmt.Sprintf("http://localhost?code=%s&id=%s", token.Code(), cmd.UserUUID),
				payload.VerificationURL)

			return nil
		},
	}
	verificationURL := &url.URL{Scheme: "http", Host: "localhost"}
	testLogger, _ := test.NewNullLogger()

	handler := NewSendVerificationEmailHandler(
		tokenRepo,
		verificationURL,
		mockPublisher,
		logrus.NewEntry(testLogger),
	)

	require.NoError(t, handler.Handle(context.Background(), cmd))
	assert.Equal(t, 1, publisherCall)
}
