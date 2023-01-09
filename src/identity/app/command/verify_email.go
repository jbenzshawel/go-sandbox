package command

import (
	"context"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/jbenzshawel/go-sandbox/common/cerror"
	"github.com/jbenzshawel/go-sandbox/common/decorator"
	"github.com/jbenzshawel/go-sandbox/identity/domain"
	"github.com/jbenzshawel/go-sandbox/identity/infrastructure/idp"
)

type VerifyEmail struct {
	UserId uuid.UUID
	Code   string
}

type VerifyEmailHandler decorator.CommandHandler[VerifyEmail]

type verifyEmailHandler struct {
	userRepo         domain.UserRepository
	tokenRepo        domain.TokenRepository
	identityProvider idp.IdentityProvider
}

func NewVerifyEmailHandler(
	userRepo domain.UserRepository,
	tokenRepo domain.TokenRepository,
	identityProvider idp.IdentityProvider,
	logger *logrus.Entry,
) VerifyEmailHandler {
	if userRepo == nil {
		panic("nil userRepo")
	}

	if tokenRepo == nil {
		panic("nil tokenRepo")
	}

	if identityProvider == nil {
		panic("nil identityProvider")
	}

	if logger == nil {
		panic("nil logger")
	}

	return decorator.ApplyCommandDecorators[VerifyEmail](
		verifyEmailHandler{
			userRepo:         userRepo,
			tokenRepo:        tokenRepo,
			identityProvider: identityProvider,
		},
		logger,
	)
}

func (h verifyEmailHandler) Handle(ctx context.Context, cmd VerifyEmail) error {
	verificationToken := h.tokenRepo.GetToken(cmd.UserId)

	if verificationToken != cmd.Code {
		return cerror.NewValidationError("bad request", map[string]string{"code": "email verification link expired"})
	}

	h.tokenRepo.ClearToken(cmd.UserId)

	// TODO: Set email verified in db and keycloak

	return nil
}
