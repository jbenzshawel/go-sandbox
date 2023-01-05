package rest

import (
	"github.com/gin-gonic/gin"

	"github.com/jbenzshawel/go-sandbox/common/rest"
	"github.com/jbenzshawel/go-sandbox/notification/app"
)

type HttpHandler struct {
	application app.Application
	healthCheck *rest.HealthCheckHandler
}

func NewHttpHandler(application app.Application) *HttpHandler {
	return &HttpHandler{
		application: application,
		healthCheck: rest.NewHealthCheckHandler(application.Logger),
	}
}

func (s *HttpHandler) HealthCheck(ctx *gin.Context) {
	s.healthCheck.Handler(ctx)
}
