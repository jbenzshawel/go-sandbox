package query

import (
	"context"

	"github.com/google/uuid"

	"github.com/sirupsen/logrus"

	"github.com/jbenzshawel/go-sandbox/common/decorator"
	"github.com/jbenzshawel/go-sandbox/identity/domain"
)

type UserByUUID struct {
	UUID uuid.UUID
}

type UserByUUIDHandler decorator.QueryHandler[UserByUUID, *domain.User]

type userByUUIDHandler struct {
	userRepo domain.UserRepository
}

func (h userByUUIDHandler) Handle(ctx context.Context, userByUUID UserByUUID) (*domain.User, error) {
	user, err := h.userRepo.GetUserByUUID(userByUUID.UUID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func NewUserByUUIDHandler(
	userRepo domain.UserRepository,
	logger *logrus.Entry,
) UserByUUIDHandler {
	if userRepo == nil {
		panic("nil userRepo")
	}

	if logger == nil {
		panic("nil logger")
	}

	return decorator.ApplyQueryDecorators[UserByUUID, *domain.User](
		userByUUIDHandler{userRepo: userRepo},
		logger,
	)
}
