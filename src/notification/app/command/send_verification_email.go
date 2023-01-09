package command

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/jbenzshawel/go-sandbox/common/decorator"
	"github.com/jbenzshawel/go-sandbox/notification/infrastructure"
)

type SendVerificationEmail struct {
	UserUUID        uuid.UUID
	FirstName       string
	Email           string
	Code            string
	VerificationURL string
}

type SendVerificationEmailHandler decorator.CommandHandler[SendVerificationEmail]

type sendVerificationEmailHandler struct {
	emailSender infrastructure.EmailSender
}

func NewSendVerificationEmailHandler(
	emailSender infrastructure.EmailSender,
	logger *logrus.Entry,
) SendVerificationEmailHandler {
	if emailSender == nil {
		panic("nil emailSender")
	}

	return decorator.ApplyCommandDecorators[SendVerificationEmail](
		sendVerificationEmailHandler{
			emailSender: emailSender,
		},
		logger,
	)
}

const verificationEmailTemplate = `<html>
<body>
<p>Hi %s,</p>
<p>Click the link below to verify your email:<br/><a href="%s">%s</a></p>
</body>
</html>`

func (h sendVerificationEmailHandler) Handle(ctx context.Context, cmd SendVerificationEmail) error {
	email := &infrastructure.Email{
		To:      cmd.Email,
		Subject: "Let's confirm your account",
		Body:    fmt.Sprintf(verificationEmailTemplate, cmd.FirstName, cmd.VerificationURL, cmd.VerificationURL),
		IsHtml:  true,
	}

	return h.emailSender.SendMail(email)
}
