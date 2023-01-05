package app

import (
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type Application struct {
	Commands Commands
	Queries  Queries
	Logger   *logrus.Entry
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
	}
}
