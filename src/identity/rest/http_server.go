package rest

import (
	"github.com/gin-gonic/gin"

	"github.com/jbenzshawel/go-sandbox/common/rest"
	"github.com/jbenzshawel/go-sandbox/identity/app"
)

type HttpServer struct {
	application app.Application
	healthCheck *rest.HealthCheckHandler
}

func NewHttpServer(application app.Application) *HttpServer {
	return &HttpServer{
		application: application,
		healthCheck: rest.NewHealthCheckHandler(
			application.Logger,
			rest.GetDatabaseHealthCheck(app.DbProvider),
		),
	}
}

func (s *HttpServer) HealthCheck(ctx *gin.Context) {
	s.healthCheck.Handler(ctx)
}
