package app

import (
	"database/sql"
	"net/url"
	"os"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/jbenzshawel/go-sandbox/identity/app/command"
	"github.com/jbenzshawel/go-sandbox/identity/app/query"
	"github.com/jbenzshawel/go-sandbox/identity/app/service"
	"github.com/jbenzshawel/go-sandbox/identity/domain/user"
	"github.com/jbenzshawel/go-sandbox/identity/infrastructure"
	"github.com/jbenzshawel/go-sandbox/identity/infrastructure/idp"
	"github.com/jbenzshawel/go-sandbox/identity/infrastructure/storage"
)

type Application struct {
	Commands commands
	Queries  queries
	Services services
	Logger   *logrus.Entry
}

type commands struct {
	CreateUser            command.UserCreateHandler
	SendVerificationEmail command.SendVerificationEmailHandler
	VerifyEmail           command.VerifyEmailHandler
}

type queries struct {
	UserByEmail query.UserByEmailHandler
	UserByUUID  query.UserByUUIDHandler
}

type services struct {
	PermissionService *service.PermissionService
}

var verificationTokenCache = storage.NewVerificationTokenCache()

func NewApplication(publisher infrastructure.Publisher) Application {
	logger := logrus.NewEntry(logrus.StandardLogger())

	identityProvider := buildIdentityProvider()

	userRepo := buildUserRepo()
	verificationTokenRepo := storage.NewVerificationTokenRepository(verificationTokenCache)

	verificationURL, err := url.Parse("http://localhost") // TODO: pull from config
	if err != nil {
		panic(errors.Wrap(err, "failed to parse verification URL"))
	}

	return Application{
		Commands: commands{
			CreateUser: command.NewCreateUserHandler(userRepo, identityProvider, logger),
			SendVerificationEmail: command.NewSendVerificationEmailHandler(
				verificationTokenRepo,
				verificationURL,
				publisher,
				logger,
			),
			VerifyEmail: command.NewVerifyEmailHandler(userRepo, verificationTokenRepo, identityProvider, logger),
		},
		Queries: queries{
			UserByEmail: query.NewUserByEmailHandler(userRepo, logger),
			UserByUUID:  query.NewUserByUUIDHandler(userRepo, logger),
		},
		Services: services{
			PermissionService: service.NewPermissionService(userRepo),
		},
		Logger: logger,
	}
}

func DbProvider() (*sql.DB, error) {
	return sql.Open("postgres", os.Getenv("IDENTITY_POSTGRES"))
}

func buildUserRepo() user.Repository {
	if userSqlRepo, ok := storage.TryCreateUserSqlRepository(); ok {
		return userSqlRepo
	}

	return storage.NewUserMemoryRepository()
}

func buildIdentityProvider() idp.IdentityProvider {
	return idp.NewKeyCloakProvider(
		os.Getenv("IDP_BASE_PATH"),
		os.Getenv("IDP_ADMIN_USER"),
		os.Getenv("IDP_ADMIN_PASSWORD"),
		os.Getenv("IDP_REALM"),
	)
}
