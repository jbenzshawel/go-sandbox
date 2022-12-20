package rest

import (
	"github.com/jbenzshawel/go-sandbox/identity/app"
)

type HttpServer struct {
	application app.Application
}

func NewHttpServer(application app.Application) *HttpServer {
	return &HttpServer{
		application: application,
	}
}
