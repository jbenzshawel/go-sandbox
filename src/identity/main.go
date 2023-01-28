package main

import (
	"context"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"

	"github.com/jbenzshawel/go-sandbox/common/auth"
	"github.com/jbenzshawel/go-sandbox/common/messaging"
	"github.com/jbenzshawel/go-sandbox/identity/app"
	"github.com/jbenzshawel/go-sandbox/identity/rest"
)

func main() {
	nc, err := natsConnection()
	if err != nil {
		panic(err)
	}
	authProvider, err := buildAuthProvider()
	if err != nil {
		panic(err)
	}

	application := app.NewApplication(messaging.NewNatsPublisher(nc))
	httpHandler := rest.NewHttpHandler(application, nc, authProvider)

	router := gin.Default() // TODO: Update gin config for production
	router.POST("/identity-client/callback", httpHandler.OAuthCallback)

	router.GET("/health", httpHandler.HealthCheck)

	router.POST("/user", httpHandler.CreateUser)
	router.POST("/user/:uuid/send-verification", httpHandler.SendVerification)
	router.POST("/user/:uuid/verify", httpHandler.VerifyUser)
	router.GET("/user/:uuid", httpHandler.GetUserByUUID)

	err = router.Run(":" + os.Getenv("IDENTITY_HTTP_PORT"))
	if err != nil {
		panic(err)
	}
}

func natsConnection() (*nats.Conn, error) {
	return nats.Connect(os.Getenv("NATS_URL"))
}

func buildAuthProvider() (*auth.OIDCProvider, error) {
	return auth.NewOIDCProvider(context.Background(), auth.OIDCConfig{
		IssuerURL:    os.Getenv("IDP_ISSUER_URL"),
		RedirectURL:  os.Getenv("IDP_REDIRECT_URL"),
		ClientID:     os.Getenv("IDP_CLIENT_ID"),
		ClientSecret: os.Getenv("IDP_CLIENT_SECRET"),
	})
}
