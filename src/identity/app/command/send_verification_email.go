package command

import (
	"context"
	"crypto/rand"
	"net/url"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/vmihailenco/msgpack/v5"

	"github.com/jbenzshawel/go-sandbox/common/decorator"
	"github.com/jbenzshawel/go-sandbox/common/messaging"
	"github.com/jbenzshawel/go-sandbox/identity/domain"
)

type SendVerificationEmail struct {
	UserUUID  uuid.UUID
	FirstName string
	Email     string
}

type SendVerificationEmailHandler decorator.CommandHandler[SendVerificationEmail]

type sendVerificationEmailHandler struct {
	tokenRepo       domain.TokenRepository
	verificationURL *url.URL
	publisher       messaging.Publisher
}

func NewSendVerificationEmailHandler(
	tokenRepo domain.TokenRepository,
	verificationURL *url.URL,
	publisher messaging.Publisher,
	logger *logrus.Entry,
) SendVerificationEmailHandler {
	if tokenRepo == nil {
		panic("nil tokenRepo")
	}

	if verificationURL == nil {
		panic("nil verificationURL")
	}

	if publisher == nil {
		panic("nil publisher")
	}

	if logger == nil {
		panic("nil logger")
	}

	return decorator.ApplyCommandDecorators[SendVerificationEmail](
		sendVerificationEmailHandler{
			tokenRepo:       tokenRepo,
			verificationURL: verificationURL,
			publisher:       publisher,
		},
		logger,
	)
}

func (h sendVerificationEmailHandler) Handle(ctx context.Context, cmd SendVerificationEmail) error {
	token, err := generateToken()
	if err != nil {
		return err
	}

	h.tokenRepo.SaveToken(cmd.UserUUID, token)

	v := url.Values{}
	v.Set("token", token)
	v.Set("id", cmd.UserUUID.String())

	msgBytes, err := msgpack.Marshal(&messaging.VerifyEmail{
		UserUUID:        cmd.UserUUID,
		FirstName:       cmd.FirstName,
		Email:           cmd.Email,
		Code:            token,
		VerificationURL: h.verificationURL.String() + "?" + v.Encode(),
	})
	if err != nil {
		return err
	}

	return h.publisher.Publish(messaging.TOPIC_VERIFY_EMAIL, msgBytes)
}

const tokenChars = "1234567890"

func generateToken() (string, error) {
	buffer := make([]byte, 6)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}

	tokenCharsLength := len(tokenChars)
	for i := 0; i < 6; i++ {
		buffer[i] = tokenChars[int(buffer[i])%tokenCharsLength]
	}

	return string(buffer), nil
}
