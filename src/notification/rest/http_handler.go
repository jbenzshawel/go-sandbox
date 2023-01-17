package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"

	"github.com/jbenzshawel/go-sandbox/common/rest"
	"github.com/jbenzshawel/go-sandbox/notification/app"
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
			rest.GetNatsHealthCheck(nc),
		),
	}
}

func (s *HttpHandler) HealthCheck(ctx *gin.Context) {
	s.healthCheck.Handler(ctx)
}
