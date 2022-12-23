package app

import (
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"

	"github.com/jbenzshawel/go-sandbox/identity/app/command"
	"github.com/jbenzshawel/go-sandbox/identity/app/query"
	"github.com/jbenzshawel/go-sandbox/identity/domain"
	"github.com/jbenzshawel/go-sandbox/identity/infrastructure"
)

type Application struct {
	Commands Commands
	Queries  Queries
	Logger   *logrus.Entry
}

type Commands struct {
	RegisterUser command.RegisterUserHandler
}

type Queries struct {
	UserByEmail query.UserByEmailHandler
}

func NewApplication() Application {
	logger := logrus.NewEntry(logrus.StandardLogger())
	userRepo := getUserRepo()

	return Application{
		Commands: Commands{
			RegisterUser: command.NewRegisterUserHandler(userRepo, logger),
		},
		Queries: Queries{
			UserByEmail: query.NewUserByEmailHandler(userRepo, logger),
		},
		Logger: logger,
	}
}

func getUserRepo() domain.UserRepository {
	if userSqlRepo, ok := infrastructure.TryCreateUserSqlRepository(); ok {
		return userSqlRepo
	}

	return infrastructure.NewUserMemoryRepository()
}
