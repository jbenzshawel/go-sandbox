package command

import (
	"context"
	"net/url"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/vmihailenco/msgpack/v5"

	"github.com/jbenzshawel/go-sandbox/common/decorator"
	"github.com/jbenzshawel/go-sandbox/common/messaging"
	"github.com/jbenzshawel/go-sandbox/identity/domain/token"
	"github.com/jbenzshawel/go-sandbox/identity/infrastructure"
)

type SendVerificationEmail struct {
	UserUUID  uuid.UUID
	FirstName string
	Email     string
}

type SendVerificationEmailHandler decorator.CommandHandler[SendVerificationEmail]

type sendVerificationEmailHandler struct {
	tokenRepo       token.Repository
	verificationURL *url.URL
	publisher       infrastructure.Publisher
}

func NewSendVerificationEmailHandler(
	tokenRepo token.Repository,
	verificationURL *url.URL,
	publisher infrastructure.Publisher,
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
	t, err := token.NewToken()
	if err != nil {
		return err
	}

	h.tokenRepo.SaveToken(cmd.UserUUID, t)

	v := url.Values{}
	v.Set("code", t.Code())
	v.Set("id", cmd.UserUUID.String())

	msgBytes, err := msgpack.Marshal(&messaging.VerifyEmail{
		UserUUID:        cmd.UserUUID,
		FirstName:       cmd.FirstName,
		Email:           cmd.Email,
		Code:            t.Code(),
		VerificationURL: h.verificationURL.String() + "?" + v.Encode(),
	})
	if err != nil {
		return err
	}

	return h.publisher.Publish(messaging.TopicVerifyEmail, msgBytes)
}
