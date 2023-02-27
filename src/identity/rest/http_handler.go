package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jbenzshawel/go-sandbox/common/auth"
	"github.com/jbenzshawel/go-sandbox/identity/app"
)

type HttpHandler struct {
	app          app.Application
	authProvider *auth.OIDCProvider
}

func NewHttpHandler(application app.Application, authProvider *auth.OIDCProvider) *HttpHandler {
	return &HttpHandler{
		app:          application,
		authProvider: authProvider,
	}
}

func (h *HttpHandler) OAuthCallback(ctx *gin.Context) {
	h.authProvider.CallbackHandler(ctx)
}

type authenticatedHandler func(ctx *gin.Context, authUser *auth.User)

func (h *HttpHandler) authenticate(ctx *gin.Context, handler authenticatedHandler) {
	authUser, err := h.authProvider.Authenticate(ctx)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		h.app.Logger.WithError(err).Warn("failed to authenticate user")
		return
	}
	if authUser == nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	handler(ctx, authUser)
}

func (h *HttpHandler) parseUUIDParam(ctx *gin.Context) (uuid.UUID, bool) {
	queryParam := ctx.Param("uuid")
	if queryParam == "" {
		return uuid.Nil, false
	}

	paramUUID, err := uuid.Parse(queryParam)
	if err != nil {
		return uuid.Nil, false
	}

	return paramUUID, true
}
