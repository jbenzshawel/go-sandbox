package main

import (
	"context"
	"os"

	"github.com/jbenzshawel/go-sandbox/common/auth"
	"github.com/jbenzshawel/go-sandbox/identity/app"
	"github.com/jbenzshawel/go-sandbox/identity/rest"
)

func main() {
	authProvider, err := buildAuthProvider()
	if err != nil {
		panic(err)
	}

	err = rest.NewHttpHandler(app.NewApplication(), authProvider).
		Configure().
		Run(":" + os.Getenv("IDENTITY_HTTP_PORT"))

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
