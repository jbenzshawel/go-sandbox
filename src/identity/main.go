package main

import (
	"os"

	"github.com/gin-gonic/gin"

	"github.com/jbenzshawel/go-sandbox/identity/app"
	"github.com/jbenzshawel/go-sandbox/identity/rest"
)

func main() {
	application := app.NewApplication()
	httpHandler := rest.NewHttpHandler(application)

	router := gin.Default() // TODO: Update gin config for production
	router.GET("/health", httpHandler.HealthCheck)

	router.POST("/user", httpHandler.CreateUser)
	router.GET("/user/:uuid", httpHandler.GetUserByUUID)

	err := router.Run(":" + os.Getenv("IDENTITY_HTTP_PORT"))
	if err != nil {
		panic(err)
	}
}
