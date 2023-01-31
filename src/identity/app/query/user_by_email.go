package query

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/jbenzshawel/go-sandbox/common/decorator"
	"github.com/jbenzshawel/go-sandbox/identity/domain/user"
)

type UserByEmail struct {
	Email string
}

type UserByEmailHandler decorator.QueryHandler[UserByEmail, *user.User]

type userByEmailHandler struct {
	userRepo user.Repository
}

func (h userByEmailHandler) Handle(ctx context.Context, userByEmail UserByEmail) (*user.User, error) {
	u, err := h.userRepo.GetByEmail(userByEmail.Email)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func NewUserByEmailHandler(
	userRepo user.Repository,
	logger *logrus.Entry,
) UserByEmailHandler {
	if userRepo == nil {
		panic("nil userRepo")
	}

	if logger == nil {
		panic("nil logger")
	}

	return decorator.ApplyQueryDecorators[UserByEmail, *user.User](
		userByEmailHandler{userRepo: userRepo},
		logger,
	)
}
