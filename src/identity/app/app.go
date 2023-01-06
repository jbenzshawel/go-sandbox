package app

import (
	"database/sql"
	"net/url"
	"os"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"

	"github.com/jbenzshawel/go-sandbox/identity/app/command"
	"github.com/jbenzshawel/go-sandbox/identity/app/query"
	"github.com/jbenzshawel/go-sandbox/identity/domain"
	"github.com/jbenzshawel/go-sandbox/identity/infrastructure/idp"
	"github.com/jbenzshawel/go-sandbox/identity/infrastructure/publisher"
	"github.com/jbenzshawel/go-sandbox/identity/infrastructure/storage"
)

type Application struct {
	Commands Commands
	Queries  Queries
	Logger   *logrus.Entry
}

type Commands struct {
	CreateUser            command.UserCreateHandler
	SendVerificationEmail command.SendVerificationEmailHandler
}

type Queries struct {
	UserByEmail query.UserByEmailHandler
	UserByUUID  query.UserByUUIDHandler
}

func NewApplication() Application {
	// TODO: Refactor how this dependency graph is built?
	logger := logrus.NewEntry(logrus.StandardLogger())

	publishers := publisher.NewNatsPublisher(os.Getenv("NATS_URL"))

	verificationTokenCache := storage.NewVerificationTokenCache()
	verificationTokenRepo := storage.NewVerificationTokenRepository(verificationTokenCache)

	userRepo := getUserRepo()
	identityProvider := getIdentityProvider()

	verificationURL, err := url.Parse("http://localhost") // TODO: pull from config
	if err != nil {
		panic(err)
	}

	return Application{
		Commands: Commands{
			CreateUser: command.NewCreateUserHandler(userRepo, identityProvider, logger),
			SendVerificationEmail: command.NewSendVerificationEmailHandler(
				verificationTokenRepo,
				verificationURL,
				publishers.NotifyVerifyEmailPublisher(),
				logger,
			),
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
