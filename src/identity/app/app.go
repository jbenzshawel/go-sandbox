package app

import (
	"database/sql"
	"os"

	"github.com/jbenzshawel/go-sandbox/identity/infrastructure/idp"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"

	"github.com/jbenzshawel/go-sandbox/identity/app/command"
	"github.com/jbenzshawel/go-sandbox/identity/app/query"
	"github.com/jbenzshawel/go-sandbox/identity/domain"
	"github.com/jbenzshawel/go-sandbox/identity/infrastructure/storage"
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
	UserByUUID  query.UserByUUIDHandler
}

func NewApplication() Application {
	logger := logrus.NewEntry(logrus.StandardLogger())
	userRepo := getUserRepo()
	identityProvider := getIdentityProvider()

	return Application{
		Commands: Commands{
			RegisterUser: command.NewRegisterUserHandler(userRepo, identityProvider, logger),
		},
		Queries: Queries{
			UserByEmail: query.NewUserByEmailHandler(userRepo, logger),
			UserByUUID:  query.NewUserByUUIDHandler(userRepo, logger),
		},
		Logger: logger,
	}
}

func DbProvider() (*sql.DB, error) {
	return sql.Open("postgres", os.Getenv("IDENTITY_POSTGRES"))
}

func getUserRepo() domain.UserRepository {
	if userSqlRepo, ok := storage.TryCreateUserSqlRepository(); ok {
		return userSqlRepo
	}

	return storage.NewUserMemoryRepository()
}

func getIdentityProvider() idp.IdentityProvider {
	return idp.NewKeyCloakProvider(
		os.Getenv("IDP_BASE_PATH"),
		os.Getenv("IDP_ADMIN_USER"),
		os.Getenv("IDP_ADMIN_PASSWORD"),
		os.Getenv("IDP_REALM"),
	)
}
