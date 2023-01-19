package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"

	"github.com/jbenzshawel/go-sandbox/common/rest"
	"github.com/jbenzshawel/go-sandbox/notification/app"
	"github.com/jbenzshawel/go-sandbox/notification/infrastructure"
)

type HttpHandler struct {
	app         app.Application
	healthCheck *rest.HealthCheckHandler
}

func NewHttpHandler(app app.Application, nc *nats.Conn) *HttpHandler {
	emailClient := infrastructure.NewEmailClient(app.Config.Email)
	return &HttpHandler{
		app: app,
		healthCheck: rest.NewHealthCheckHandler(
			app.Logger,
			rest.GetNatsHealthCheck(nc),
			emailClient.HealthCheck(),
		),
	}
}

func (s *HttpHandler) HealthCheck(ctx *gin.Context) {
	s.healthCheck.Handler(ctx)
}
