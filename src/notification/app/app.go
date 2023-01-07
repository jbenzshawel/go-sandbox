package app

import (
	"os"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type appConfig struct {
	HttpPort string
	NatsURL  string
}

type Application struct {
	Commands Commands
	Queries  Queries
	Logger   *logrus.Entry
	Config   appConfig
}

type Commands struct {
}

type Queries struct {
}

func NewApplication() Application {
	logger := logrus.NewEntry(logrus.StandardLogger())

	return Application{
		Commands: Commands{},
		Queries:  Queries{},
		Logger:   logger,
		Config: appConfig{
			HttpPort: os.Getenv("NOTIFICATION_HTTP_PORT"),
			NatsURL:  os.Getenv("NATS_URL"),
		},
	}
}
