package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jbenzshawel/go-sandbox/common/rest"
	"github.com/jbenzshawel/go-sandbox/identity/app"
	"github.com/nats-io/nats.go"
)

type HttpHandler struct {
	app         app.Application
	healthCheck *rest.HealthCheckHandler
}

func NewHttpHandler(application app.Application, nc *nats.Conn) *HttpHandler {
	return &HttpHandler{
		app: application,
		healthCheck: rest.NewHealthCheckHandler(
			application.Logger,
			rest.GetDatabaseHealthCheck(app.DbProvider),
			rest.GetNatsHealthCheck(nc),
		),
	}
}

func (h *HttpHandler) HealthCheck(ctx *gin.Context) {
	h.healthCheck.Handler(ctx)
}

func (h *HttpHandler) parseUUIDParam(ctx *gin.Context) (uuid.UUID, bool) {
	queryParam := ctx.Param("uuid")
	if queryParam == "" {
		return uuid.Nil, false
	}

	paramUUID, err := uuid.Parse(queryParam)
	if err != nil {
		return uuid.Nil, false
	}

	return paramUUID, true
}
