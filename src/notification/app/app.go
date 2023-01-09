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
	Email    emailConfig
}

type emailConfig struct {
	Addr string
	Host string
	From string
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
	emailClient := infrastructure.NewEmailClient(config.Email.Addr, config.Email.Host, config.Email.From)

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
		Email: emailConfig{
			Addr: os.Getenv("SMTP_URL"),
			Host: os.Getenv("SMTP_HOST"),
			From: os.Getenv("SMTP_FROM"),
		},
	}
}
