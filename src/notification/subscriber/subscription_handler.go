package subscriber

import (
	"github.com/vmihailenco/msgpack/v5"

	"github.com/jbenzshawel/go-sandbox/common/messaging"
	"github.com/jbenzshawel/go-sandbox/notification/app"
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
		s.application.Logger.WithError(err).Error("failed to unmarshal VerifyEmail message")
		return
	}
	s.application.Logger.
		WithField("email", message.Email).
		WithField("uuid", message.UserUUID).
		WithField("token", message.Code).
		Info("send email msg received")
}
