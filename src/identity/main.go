package main

import (
	"context"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jbenzshawel/go-sandbox/common/auth"
	"github.com/jbenzshawel/go-sandbox/identity/app"
	"github.com/jbenzshawel/go-sandbox/identity/rest"
)

func main() {
	authProvider, err := buildAuthProvider()
	if err != nil {
		panic(err)
	}

	application := app.NewApplication()
	httpHandler := rest.NewHttpHandler(application, authProvider)

	router := gin.Default() // TODO: Update gin config for production
	router.POST("/identity-client/callback", httpHandler.OAuthCallback)

	router.GET("/health", application.HealthCheck.Handler)

	router.POST("/user", httpHandler.CreateUser)
	router.POST("/user/:uuid/send-verification", httpHandler.SendVerification)
	router.POST("/user/:uuid/verify", httpHandler.VerifyUser)
	router.GET("/user/:uuid", httpHandler.GetUserByUUID)

	err = router.Run(":" + os.Getenv("IDENTITY_HTTP_PORT"))
	if err != nil {
		panic(err)
	}
}

func buildAuthProvider() (*auth.OIDCProvider, error) {
	return auth.NewOIDCProvider(context.Background(), auth.OIDCConfig{
		IssuerURL:    os.Getenv("IDP_ISSUER_URL"),
		RedirectURL:  os.Getenv("IDP_REDIRECT_URL"),
		ClientID:     os.Getenv("IDP_CLIENT_ID"),
		ClientSecret: os.Getenv("IDP_CLIENT_SECRET"),
	})
}
