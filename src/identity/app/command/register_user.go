package command

import (
	"context"
	"fmt"
	"time"

	"github.com/jbenzshawel/go-sandbox/identity/infrastructure/idp"

	"github.com/jbenzshawel/go-sandbox/common/cerror"
	"github.com/jbenzshawel/go-sandbox/common/decorator"
	"github.com/jbenzshawel/go-sandbox/identity/domain"
	"github.com/sirupsen/logrus"
)

type RegisterUser struct {
	FirstName       string
	LastName        string
	Email           string
	Password        string
	ConfirmPassword string
}

type RegisterUserHandler decorator.CommandHandler[RegisterUser]

type registerUserHandler struct {
	userRepo         domain.UserRepository
	identityProvider idp.IdentityProvider
}

func NewRegisterUserHandler(
	userRepo domain.UserRepository,
	identityProvider idp.IdentityProvider,
	logger *logrus.Entry,
) RegisterUserHandler {
	if userRepo == nil {
		panic("nil userRepo")
	}

	return decorator.ApplyCommandDecorators[RegisterUser](
		registerUserHandler{
			userRepo:         userRepo,
			identityProvider: identityProvider,
		},
		logger,
	)
}

func (h registerUserHandler) Handle(ctx context.Context, cmd RegisterUser) error {
	// TODO: Create some sort of config driven password validator
	validationErrors := map[string]string{}
	if cmd.Password != cmd.ConfirmPassword {
		validationErrors["confirmPassword"] = "password and confirm password must match"
	}

	existingUser, err := h.userRepo.GetUserByEmail(cmd.Email)
	if err != nil {
		return err
	}

	if existingUser != nil {
		validationErrors["email"] = fmt.Sprintf("user with email %s already exists", cmd.Email)
	}

	if len(validationErrors) > 0 {
		return cerror.NewValidationError("Invalid request", validationErrors)
	}

	user := domain.User{
		FirstName:     cmd.FirstName,
		LastName:      cmd.LastName,
		Email:         cmd.Email,
		Enabled:       true,
		CreatedAt:     time.Now(),
		LastUpdatedAt: time.Now(),
	}

	userUUID, err := h.identityProvider.CreateUser(ctx, user, cmd.Password)
	if err != nil {
		return err
	}
	user.UUID = userUUID

	err = h.userRepo.InsertUser(user)
	if err != nil {
		return err
	}
	return nil
}
