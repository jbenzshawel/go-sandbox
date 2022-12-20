package main

import (
	"os"

	"github.com/gin-gonic/gin"

	"github.com/jbenzshawel/go-sandbox/identity/app"
	"github.com/jbenzshawel/go-sandbox/identity/handlers"
)

func main() {
	application := app.NewApplication()
	httpServer := handlers.NewHttpServer(application)

	router := gin.Default()
	router.POST("/register", httpServer.RegisterUser)

	err := router.Run(":" + os.Getenv("HTTP_PORT"))
	if err != nil {
		panic(err)
	}
}
