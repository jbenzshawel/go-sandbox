package rest

import (
	"github.com/gin-gonic/gin"
	
	"github.com/jbenzshawel/go-sandbox/common/rest"
	"github.com/jbenzshawel/go-sandbox/identity/app"
)

type HttpServer struct {
	application app.Application
}

func NewHttpServer(application app.Application) *HttpServer {
	return &HttpServer{application: application}
}

func (s *HttpServer) HealthCheck(ctx *gin.Context) {
	healthCheck := rest.NewHealthCheckHandler(
		s.application.Logger,
		rest.GetDatabaseHealthCheck(app.DbProvider),
	)

	healthCheck.Handler(ctx)
}
