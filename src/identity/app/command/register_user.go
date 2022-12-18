package command

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"

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
	if cmd.Password != cmd.ConfirmPassword {
		// TODO: Create validation error type and return
		return errors.New("password and confirm password must match")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(cmd.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := domain.User{
		UUID:          uuid.New().String(),
		FirstName:     cmd.FirstName,
		LastName:      cmd.LastName,
		Email:         cmd.Email,
		Enabled:       true,
		CreatedAt:     time.Now(),
		LastUpdatedAt: time.Now(),
	}
	err = h.userRepo.RegisterUser(user, string(hash))
	if err != nil {
		return err
	}
	return nil
}
