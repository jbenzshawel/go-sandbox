package main

import (
	"os"

	"github.com/gin-gonic/gin"

	"github.com/jbenzshawel/go-sandbox/identity/app"
	"github.com/jbenzshawel/go-sandbox/identity/rest"
)

func main() {
	application := app.NewApplication()
	httpServer := rest.NewHttpServer(application)

	router := gin.Default()
	router.GET("/health", httpServer.HealthCheck)
	router.POST("/register", httpServer.RegisterUser)

	err := router.Run(":" + os.Getenv("HTTP_PORT"))
	if err != nil {
		panic(err)
	}
}
