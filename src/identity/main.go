package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"

	"github.com/jbenzshawel/go-sandbox/common/messaging"
	"github.com/jbenzshawel/go-sandbox/identity/app"
	"github.com/jbenzshawel/go-sandbox/identity/rest"
)

func main() {
	nc, err := nats.Connect(os.Getenv("NATS_URL"))
	if err != nil {
		panic(err)
	}
	application := app.NewApplication(messaging.NewNatsPublisher(nc))
	httpHandler := rest.NewHttpHandler(application, nc)

	router := gin.Default() // TODO: Update gin config for production
	router.GET("/health", httpHandler.HealthCheck)

	router.POST("/user", httpHandler.CreateUser)
	router.POST("/user/:uuid/verify", httpHandler.VerifyUser)
	router.GET("/user/:uuid", httpHandler.GetUserByUUID)

	err = router.Run(":" + os.Getenv("IDENTITY_HTTP_PORT"))
	if err != nil {
		panic(err)
	}
}
