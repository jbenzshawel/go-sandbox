package command

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/jbenzshawel/go-sandbox/common/cerror"
	"github.com/jbenzshawel/go-sandbox/common/decorator"
	"github.com/jbenzshawel/go-sandbox/identity/domain/user"
	"github.com/jbenzshawel/go-sandbox/identity/infrastructure/idp"
)

type UserCreate struct {
	FirstName       string
	LastName        string
	Email           string
	Password        string
	ConfirmPassword string
}

type UserCreateHandler decorator.CommandHandler[UserCreate]

type userCreateHandler struct {
	userRepo         user.Repository
	identityProvider idp.IdentityProvider
}

func NewCreateUserHandler(
	userRepo user.Repository,
	identityProvider idp.IdentityProvider,
	logger *logrus.Entry,
) UserCreateHandler {
	if userRepo == nil {
		panic("nil userRepo")
	}

	if identityProvider == nil {
		panic("nil identityProvider")
	}

	if logger == nil {
		panic("nil logger")
	}

	return decorator.ApplyCommandDecorators[UserCreate](
		userCreateHandler{
			userRepo:         userRepo,
			identityProvider: identityProvider,
		},
		logger,
	)
}

func (h userCreateHandler) Handle(ctx context.Context, cmd UserCreate) error {
	// TODO: Create some sort of config driven password validator
	validationErrors := map[string]string{}
	if cmd.Password != cmd.ConfirmPassword {
		validationErrors["confirmPassword"] = "password and confirm password must match"
	}

	existingUser, err := h.userRepo.GetByEmail(cmd.Email)
	if err != nil {
		return err
	}

	if existingUser != nil {
		validationErrors["email"] = fmt.Sprintf("user with email %s already exists", cmd.Email)
	}

	if len(validationErrors) > 0 {
		return cerror.NewValidationError("invalid user", validationErrors)
	}

	u, err := user.NewUser(
		cmd.FirstName,
		cmd.LastName,
		cmd.Email,
		false,
		true,
	)
	if err != nil {
		return err
	}

	return h.createUser(ctx, cmd, u)
}

func (h userCreateHandler) createUser(ctx context.Context, cmd UserCreate, u *user.User) error {
	userUUID, err := h.identityProvider.CreateUser(ctx, u, cmd.Password)
	if err != nil {
		return h.handleCreateUserErr(ctx, userUUID, err)
	}
	err = u.SetUUID(userUUID)
	if err != nil {
		return h.handleCreateUserErr(ctx, userUUID, err)
	}

	err = h.userRepo.Create(u)
	if err != nil {
		return h.handleCreateUserErr(ctx, userUUID, err)
	}
	return nil
}

func (h userCreateHandler) handleCreateUserErr(ctx context.Context, userUUID uuid.UUID, err error) error {
	if userUUID != uuid.Nil {
		// if we created a user with our identity provider but failed with additional
		// setup attempt to delete the created user in our idp
		deleteErr := h.identityProvider.DeleteUser(ctx, userUUID.String())
		if deleteErr != nil {
			return errors.Wrap(err, "Failed to delete partially created user")
		}
	}
	return err
}
