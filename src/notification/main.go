package main

import (
	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"

	"github.com/jbenzshawel/go-sandbox/common/messaging"
	"github.com/jbenzshawel/go-sandbox/notification/app"
	"github.com/jbenzshawel/go-sandbox/notification/rest"
	"github.com/jbenzshawel/go-sandbox/notification/subscriber"
)

func main() {
	application := app.NewApplication()

	nc, err := nats.Connect(application.Config.NatsURL)
	if err != nil {
		panic(err)
	}

	subscriptionHandler := subscriber.NewSubscriptionHandler(application)
	err = messaging.NewNatsSubscriber(nc).
		WithSubscription(messaging.TOPIC_VERIFY_EMAIL, subscriptionHandler.SendVerificationEmail).
		Subscribe()

	if err != nil {
		panic(err)
	}

	httpHandler := rest.NewHttpHandler(application)
	router := gin.Default() // TODO: Update gin config for production
	router.GET("/health", httpHandler.HealthCheck)

	err = router.Run(":" + application.Config.HttpPort)
	if err != nil {
		panic(err)
	}
}
