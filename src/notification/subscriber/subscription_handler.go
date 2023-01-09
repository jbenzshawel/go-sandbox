package subscriber

import (
	"context"

	"github.com/pkg/errors"
	"github.com/vmihailenco/msgpack/v5"

	"github.com/jbenzshawel/go-sandbox/common/messaging"
	"github.com/jbenzshawel/go-sandbox/notification/app"
	"github.com/jbenzshawel/go-sandbox/notification/app/command"
)

type SubscriptionHandler struct {
	application app.Application
}

func NewSubscriptionHandler(application app.Application) *SubscriptionHandler {
	return &SubscriptionHandler{application: application}
}

func (s *SubscriptionHandler) SendVerificationEmail(msg []byte) {
	var message messaging.VerifyEmail
	err := msgpack.Unmarshal(msg, &message)
	if err != nil {
		s.application.Logger.WithError(errors.WithStack(err)).Error("failed to unmarshal VerifyEmail message")
		return
	}

	s.application.Logger.WithField("VerifyEmail", message).
		Info("VerifyEmail msg received")

	cmd := command.SendVerificationEmail{
		UserUUID:        message.UserUUID,
		FirstName:       message.FirstName,
		Email:           message.Email,
		Code:            message.Code,
		VerificationURL: message.VerificationURL,
	}

	err = s.application.Commands.SendVerificationEmail.Handle(context.Background(), cmd)
	if err != nil {
		s.application.Logger.WithError(errors.WithStack(err)).Error("failed to handle VerifyEmail message")
	}
}
