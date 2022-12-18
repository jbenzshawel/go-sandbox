package app

import (
	"github.com/sirupsen/logrus"

	"github.com/jbenzshawel/go-sandbox/identity/app/command"
	"github.com/jbenzshawel/go-sandbox/identity/infrastructure"
)

type Application struct {
	Commands Commands
	Logger   *logrus.Entry
}

type Commands struct {
	RegisterUser command.RegisterUserHandler
}

func NewApplication() Application {
	logger := logrus.NewEntry(logrus.StandardLogger())
	userRepo := infrastructure.NewUserMemoryRepository()

	return Application{
		Commands: Commands{
			RegisterUser: command.NewRegisterUserHandler(userRepo, logger),
		},
		Logger: logger,
	}
}
