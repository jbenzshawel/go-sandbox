package command

import (
	"context"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/jbenzshawel/go-sandbox/common/cerror"
	"github.com/jbenzshawel/go-sandbox/common/decorator"
	"github.com/jbenzshawel/go-sandbox/identity/domain/token"
	"github.com/jbenzshawel/go-sandbox/identity/domain/user"
	"github.com/jbenzshawel/go-sandbox/identity/infrastructure/idp"
)

type VerifyEmail struct {
	UserId uuid.UUID
	Code   string
}

type VerifyEmailHandler decorator.CommandHandler[VerifyEmail]

type verifyEmailHandler struct {
	userRepo         user.Repository
	tokenRepo        token.Repository
	identityProvider idp.IdentityProvider
}

func NewVerifyEmailHandler(
	userRepo user.Repository,
	tokenRepo token.Repository,
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
	isValid := token.Verify(h.tokenRepo, cmd.UserId, cmd.Code)

	if !isValid {
		return cerror.NewValidationError("bad request", map[string]string{"code": "email verification link expired"})
	}

	u, err := h.userRepo.GetUserByUUID(cmd.UserId)
	if err != nil {
		return err
	}

	u.SetEmailVerified(true)

	return h.updateUser(ctx, u)
}

func (h verifyEmailHandler) updateUser(ctx context.Context, u *user.User) error {
	err := h.identityProvider.UpdateUser(ctx, u)
	if err != nil {
		return err
	}
	err = h.userRepo.UpdateUser(u)
	if err != nil {
		return err
	}

	return nil
}
