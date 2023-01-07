package rest

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"

	"github.com/jbenzshawel/go-sandbox/common/cerror"
	"github.com/jbenzshawel/go-sandbox/common/database"
)

type healthCheckResponse struct {
	Status string   `json:"status"`
	Errors []string `json:"errors,omitempty"`
}

type HealthCheckTask func() (bool, string, error)

func GetDatabaseHealthCheck(dbProvider database.DbProvider) HealthCheckTask {
	return func() (success bool, name string, err error) {
		healthCheckName := "database"
		db, err := dbProvider()
		if err != nil {
			return false, healthCheckName, err
		}
		defer func() {
			closeErr := db.Close()
			err = cerror.CombineErrors(err, closeErr)
		}()

		err = db.Ping()
		if err != nil {
			return false, healthCheckName, err
		}

		return true, healthCheckName, nil
	}
}

func GetNatsHealthCheck(nc *nats.Conn) HealthCheckTask {
	return func() (bool, string, error) {
		healthCheckName := "nats"
		if nc.IsConnected() {
			return true, healthCheckName, nil
		}
		return false, healthCheckName, nil
	}
}

type HealthCheckHandler struct {
	checks []HealthCheckTask
	logger *logrus.Entry
}

func NewHealthCheckHandler(logger *logrus.Entry, checks ...HealthCheckTask) *HealthCheckHandler {
	return &HealthCheckHandler{logger: logger, checks: checks}
}

func (h *HealthCheckHandler) Handler(ctx *gin.Context) {
	var errs []string
	for _, check := range h.checks {
		if ok, name, err := check(); !ok {
			msg := fmt.Sprintf("%s health check failed", name)
			h.logger.WithError(err).Error(msg)
			errs = append(errs, msg)
		}
	}
	if len(errs) == 0 {
		ctx.JSON(http.StatusOK, healthCheckResponse{Status: "available"})
	} else {
		ctx.JSON(http.StatusInternalServerError, healthCheckResponse{Status: "unavailable", Errors: errs})
	}
}
