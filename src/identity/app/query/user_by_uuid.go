package query

import (
	"context"

	"github.com/google/uuid"

	"github.com/sirupsen/logrus"

	"github.com/jbenzshawel/go-sandbox/common/decorator"
	"github.com/jbenzshawel/go-sandbox/identity/domain/user"
)

type UserByUUID struct {
	UUID uuid.UUID
}

type UserByUUIDHandler decorator.QueryHandler[UserByUUID, *user.User]

type userByUUIDHandler struct {
	userRepo user.Repository
}

func (h userByUUIDHandler) Handle(ctx context.Context, userByUUID UserByUUID) (*user.User, error) {
	u, err := h.userRepo.GetByUUID(userByUUID.UUID)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func NewUserByUUIDHandler(
	userRepo user.Repository,
	logger *logrus.Entry,
) UserByUUIDHandler {
	if userRepo == nil {
		panic("nil userRepo")
	}

	if logger == nil {
		panic("nil logger")
	}

	return decorator.ApplyQueryDecorators[UserByUUID, *user.User](
		userByUUIDHandler{userRepo: userRepo},
		logger,
	)
}
