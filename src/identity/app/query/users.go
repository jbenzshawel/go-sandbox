package query

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/jbenzshawel/go-sandbox/common/decorator"
	"github.com/jbenzshawel/go-sandbox/identity/domain/user"
)

type Users struct {
	Page     int
	PageSize int
}

type UsersHandler decorator.QueryHandler[Users, []*user.User]

type usersHandler struct {
	userRepo user.Repository
}

func (h usersHandler) Handle(ctx context.Context, users Users) ([]*user.User, error) {
	u, err := h.userRepo.GetAll(users.Page, users.PageSize)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func NewUsersHandler(
	userRepo user.Repository,
	logger *logrus.Entry,
) UsersHandler {
	if userRepo == nil {
		panic("nil userRepo")
	}

	if logger == nil {
		panic("nil logger")
	}

	return decorator.ApplyQueryDecorators[Users, []*user.User](
		usersHandler{userRepo: userRepo},
		logger,
	)
}
