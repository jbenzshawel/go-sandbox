package main

import (
	"github.com/gin-gonic/gin"

	"github.com/jbenzshawel/go-sandbox/identity/app"
	"github.com/jbenzshawel/go-sandbox/identity/handlers"
)

func main() {
	application := app.NewApplication()
	httpServer := handlers.NewHttpServer(application)

	router := gin.Default()
	router.POST("/register", httpServer.RegisterUser)

	router.Run("localhost:8080") // TODO: Get port from environment variables
}
