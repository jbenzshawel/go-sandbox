package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"

	"github.com/jbenzshawel/go-sandbox/common/rest"
	"github.com/jbenzshawel/go-sandbox/identity/app"
)

type HttpHandler struct {
	application app.Application
	healthCheck *rest.HealthCheckHandler
}

func NewHttpHandler(application app.Application, nc *nats.Conn) *HttpHandler {
	return &HttpHandler{
		application: application,
		healthCheck: rest.NewHealthCheckHandler(
			application.Logger,
			rest.GetDatabaseHealthCheck(app.DbProvider),
			rest.GetNatsHealthCheck(nc),
		),
	}
}

func (s *HttpHandler) HealthCheck(ctx *gin.Context) {
	s.healthCheck.Handler(ctx)
}
