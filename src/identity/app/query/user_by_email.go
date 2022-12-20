package query

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/jbenzshawel/go-sandbox/common/decorator"
	"github.com/jbenzshawel/go-sandbox/identity/domain"
)

type UserByEmail struct {
	Email string
}

type UserByEmailHandler decorator.QueryHandler[UserByEmail, *domain.User]

type userByEmailHandler struct {
	userRepo domain.UserRepository
}

func (h userByEmailHandler) Handle(ctx context.Context, userByEmail UserByEmail) (*domain.User, error) {
	user, err := h.userRepo.GetUserByEmail(userByEmail.Email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func NewUserByEmailHandler(
	userRepo domain.UserRepository,
	logger *logrus.Entry,
) UserByEmailHandler {
	if userRepo == nil {
		panic("nil userRepo")
	}

	return decorator.ApplyQueryDecorators[UserByEmail, *domain.User](
		userByEmailHandler{userRepo: userRepo},
		logger,
	)
}
