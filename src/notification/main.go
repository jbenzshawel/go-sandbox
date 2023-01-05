package main

import (
	"os"

	"github.com/gin-gonic/gin"

	"github.com/jbenzshawel/go-sandbox/notification/app"
	"github.com/jbenzshawel/go-sandbox/notification/rest"
)

func main() {
	application := app.NewApplication()
	httpHandler := rest.NewHttpHandler(application)

	router := gin.Default() // TODO: Update gin config for production
	router.GET("/health", httpHandler.HealthCheck)

	err := router.Run(":" + os.Getenv("NOTIFICATION_HTTP_PORT"))
	if err != nil {
		panic(err)
	}
}
