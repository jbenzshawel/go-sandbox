package command

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"

	"github.com/jbenzshawel/go-sandbox/common/cerror"
	"github.com/jbenzshawel/go-sandbox/common/decorator"
	"github.com/jbenzshawel/go-sandbox/identity/domain"
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
	userRepo domain.UserRepository
}

func NewRegisterUserHandler(
	userRepo domain.UserRepository,
	logger *logrus.Entry,
) RegisterUserHandler {
	if userRepo == nil {
		panic("nil userRepo")
	}

	return decorator.ApplyCommandDecorators[RegisterUser](
		registerUserHandler{userRepo: userRepo},
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

	hash, err := bcrypt.GenerateFromPassword([]byte(cmd.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := domain.User{
		UUID:          uuid.New(),
		FirstName:     cmd.FirstName,
		LastName:      cmd.LastName,
		Email:         cmd.Email,
		Enabled:       true,
		CreatedAt:     time.Now(),
		LastUpdatedAt: time.Now(),
	}
	err = h.userRepo.CreateUser(user, string(hash))
	if err != nil {
		return err
	}
	return nil
}
