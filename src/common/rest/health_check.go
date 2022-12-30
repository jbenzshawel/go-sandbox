package rest

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/jbenzshawel/go-sandbox/common/cerror"
	"github.com/jbenzshawel/go-sandbox/common/database"
)

type healthCheckResponse struct {
	Status string   `json:"status,omitempty"`
	Error  []string `json:"error,omitempty"`
}

type HealthCheckTask func() (bool, string, error)

func GetDatabaseHealthCheck(dbConn database.DbProvider) HealthCheckTask {
	return func() (bool, string, error) {
		healthCheckName := "database"
		db, err := dbConn()
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
			h.logger.WithError(err).Errorf(msg)
			errs = append(errs, msg)
		}
	}
	if len(errs) == 0 {
		ctx.JSON(http.StatusOK, healthCheckResponse{Status: "available"})
	} else {
		ctx.JSON(http.StatusInternalServerError, healthCheckResponse{Status: "unavailable", Error: errs})
	}
}
