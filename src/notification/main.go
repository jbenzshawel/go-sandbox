package main

import (
	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"

	"github.com/jbenzshawel/go-sandbox/common/messaging"
	"github.com/jbenzshawel/go-sandbox/common/rest"
	"github.com/jbenzshawel/go-sandbox/notification/app"
	"github.com/jbenzshawel/go-sandbox/notification/subscriber"
)

func main() {
	application := app.NewApplication()

	nc, err := nats.Connect(application.Config.NATSURL)
	if err != nil {
		panic(err)
	}
	application.HealthCheck.AddCheck(rest.NatsHealthCheck(nc))

	subscriptionHandler := subscriber.NewSubscriptionHandler(application)
	err = messaging.NewNatsSubscriber(nc).
		WithSubscription(messaging.TopicVerifyEmail, subscriptionHandler.SendVerificationEmail).
		Subscribe()

	if err != nil {
		panic(err)
	}

	router := gin.Default() // TODO: Update gin config for production
	router.GET("/health", application.HealthCheck.Handler)

	err = router.Run(":" + application.Config.HTTPPort)
	if err != nil {
		panic(err)
	}
}
