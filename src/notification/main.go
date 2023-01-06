package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
	"github.com/vmihailenco/msgpack/v5"

	"github.com/jbenzshawel/go-sandbox/common/messaging"
	"github.com/jbenzshawel/go-sandbox/notification/app"
	"github.com/jbenzshawel/go-sandbox/notification/rest"
)

func main() {
	application := app.NewApplication()
	httpHandler := rest.NewHttpHandler(application)

	// TODO: Refactor this out of main... (just a POC right now)
	nc, err := nats.Connect(os.Getenv("NATS_URL"))
	if err != nil {
		panic(err)
	}
	_, err = nc.Subscribe(messaging.TOPIC_VERIFY_EMAIL, func(msg *nats.Msg) {
		var message messaging.VerifyEmail
		err = msgpack.Unmarshal(msg.Data, &message)
		if err != nil {
			panic(err)
		}
		application.Logger.
			WithField("email", message.Email).
			WithField("uuid", message.UserUUID).
			WithField("token", message.Code).
			Info("send email msg received")
	})
	if err != nil {
		panic(err)
	}

	router := gin.Default() // TODO: Update gin config for production
	router.GET("/health", httpHandler.HealthCheck)

	err = router.Run(":" + os.Getenv("NOTIFICATION_HTTP_PORT"))
	if err != nil {
		panic(err)
	}
}
