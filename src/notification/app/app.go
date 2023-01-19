package app

import (
	"os"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"

	"github.com/jbenzshawel/go-sandbox/notification/app/command"
	"github.com/jbenzshawel/go-sandbox/notification/infrastructure"
)

type appConfig struct {
	HTTPPort string
	NATSURL  string
	Email    infrastructure.EmailConfig
}

type Application struct {
	Commands Commands
	Logger   *logrus.Entry
	Config   appConfig
}

type Commands struct {
	SendVerificationEmail command.SendVerificationEmailHandler
}

func NewApplication() Application {
	logger := logrus.NewEntry(logrus.StandardLogger())
	config := buildConfig()
	emailClient := infrastructure.NewEmailClient(config.Email)

	return Application{
		Commands: Commands{
			SendVerificationEmail: command.NewSendVerificationEmailHandler(emailClient, logger),
		},
		Logger: logger,
		Config: config,
	}
}

func buildConfig() appConfig {
	return appConfig{
		HTTPPort: os.Getenv("NOTIFICATION_HTTP_PORT"),
		NATSURL:  os.Getenv("NATS_URL"),
		Email: infrastructure.EmailConfig{
			Addr:     os.Getenv("SMTP_URL"),
			Host:     os.Getenv("SMTP_HOST"),
			From:     os.Getenv("SMTP_FROM"),
			Username: os.Getenv("SMTP_USERNAME"),
			Password: os.Getenv("SMTP_PASSWORD"),
		},
	}
}
