package app

import (
	"database/sql"

	"net/url"
	"os"

	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/jbenzshawel/go-sandbox/common/messaging"
	"github.com/jbenzshawel/go-sandbox/common/rest"
	"github.com/jbenzshawel/go-sandbox/identity/app/command"
	"github.com/jbenzshawel/go-sandbox/identity/app/query"
	"github.com/jbenzshawel/go-sandbox/identity/app/service"
	"github.com/jbenzshawel/go-sandbox/identity/domain/user"
	"github.com/jbenzshawel/go-sandbox/identity/infrastructure/idp"
	"github.com/jbenzshawel/go-sandbox/identity/infrastructure/storage"
)

type Application struct {
	Commands    commands
	Queries     queries
	Services    services
	HealthCheck *rest.HealthCheckHandler
	Logger      *logrus.Entry
}

type commands struct {
	CreateUser            command.UserCreateHandler
	SendVerificationEmail command.SendVerificationEmailHandler
	VerifyEmail           command.VerifyEmailHandler
}

type queries struct {
	Users       query.UsersHandler
	UserByEmail query.UserByEmailHandler
	UserByUUID  query.UserByUUIDHandler
}

type services struct {
	PermissionService *service.PermissionService
}

var verificationTokenCache = storage.NewVerificationTokenCache()

func NewApplication() Application {
	logger := logrus.NewEntry(logrus.StandardLogger())
	healthCheck := rest.NewHealthCheckHandler(logger)

	nc, err := natsConnection()
	if err != nil {
		panic(errors.Wrap(err, "failed to connect to nats"))
	}
	healthCheck.AddCheck(rest.NatsHealthCheck(nc))

	db, err := openDB()
	if err != nil {
		panic(errors.Wrap(err, "failed to connect to db"))
	}
	healthCheck.AddCheck(rest.DatabaseHealthCheck(db))

	identityProvider := buildIdentityProvider()
	userRepo := buildUserRepo(db)
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
				messaging.NewNatsPublisher(nc),
				logger,
			),
			VerifyEmail: command.NewVerifyEmailHandler(userRepo, verificationTokenRepo, identityProvider, logger),
		},
		Queries: queries{
			Users:       query.NewUsersHandler(userRepo, logger),
			UserByEmail: query.NewUserByEmailHandler(userRepo, logger),
			UserByUUID:  query.NewUserByUUIDHandler(userRepo, logger),
		},
		Services: services{
			PermissionService: service.NewPermissionService(userRepo),
		},
		HealthCheck: healthCheck,
		Logger:      logger,
	}
}

func openDB() (*sql.DB, error) {
	if connectionString, ok := os.LookupEnv("IDENTITY_POSTGRES"); ok {
		db, err := sql.Open("postgres", connectionString)
		if err != nil {
			return nil, err
		}
		return db, nil
	}
	return nil, errors.New("IDENTITY_POSTGRES env not found")
}

func natsConnection() (*nats.Conn, error) {
	return nats.Connect(os.Getenv("NATS_URL"))
}

func buildUserRepo(db *sql.DB) user.Repository {
	if db == nil {
		return storage.NewUserMemoryRepository()
	}

	return storage.NewUserSqlRepository(db)
}

func buildIdentityProvider() idp.IdentityProvider {
	return idp.NewKeyCloakProvider(
		os.Getenv("IDP_BASE_PATH"),
		os.Getenv("IDP_ADMIN_USER"),
		os.Getenv("IDP_ADMIN_PASSWORD"),
		os.Getenv("IDP_REALM"),
	)
}
