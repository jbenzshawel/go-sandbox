package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

type User struct {
	UserUUID      uuid.UUID
	FirstName     string
	LastName      string
	Email         string
	EmailVerified bool
}

type OIDCConfig struct {
	IssuerURL    string
	RedirectURL  string
	ClientID     string
	ClientSecret string
}

type OIDCProvider struct {
	oauthConfig oauth2.Config
	verifier    *oidc.IDTokenVerifier
}

func NewOIDCProvider(ctx context.Context, cfg OIDCConfig) (*OIDCProvider, error) {
	provider, err := oidc.NewProvider(ctx, cfg.IssuerURL)
	if err != nil {
		return nil, err
	}

	oauth2Config := oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  cfg.RedirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	oidcConfig := &oidc.Config{
		ClientID: cfg.ClientID,
	}

	return &OIDCProvider{
		oauthConfig: oauth2Config,
		verifier:    provider.Verifier(oidcConfig),
	}, nil
}

func (p *OIDCProvider) Authenticate(ctx *gin.Context) (*User, error) {
	authHeader := ctx.GetHeader("Authorization")

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 {
		return nil, nil
	}
	idToken, err := p.verifier.Verify(ctx, parts[1])
	if err != nil {
		return nil, err
	}
	return p.mapAuthenticatedUser(idToken)
}

func (p *OIDCProvider) mapAuthenticatedUser(idToken *oidc.IDToken) (*User, error) {
	var err error
	u := &User{}
	u.UserUUID, err = uuid.Parse(idToken.Subject)
	if err != nil {
		return nil, err
	}
	var claims struct {
		FirstName     string `json:"given_name"`
		LastName      string `json:"family_name"`
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
	}
	if err = idToken.Claims(&claims); err != nil {
		return nil, err
	}
	u.FirstName = claims.FirstName
	u.LastName = claims.LastName
	u.Email = claims.Email
	u.EmailVerified = claims.EmailVerified
	return u, nil
}

func (p *OIDCProvider) CallbackHandler(ctx *gin.Context) {
	oauth2Token, err := p.oauthConfig.Exchange(ctx, ctx.Param("code"))
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	idToken, err := p.verifier.Verify(ctx, rawIDToken)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	resp := struct {
		OAuth2Token   *oauth2.Token
		IDTokenClaims *json.RawMessage
	}{oauth2Token, new(json.RawMessage)}

	if err = idToken.Claims(&resp.IDTokenClaims); err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, resp)
}
