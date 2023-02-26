package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"

	"github.com/jbenzshawel/go-sandbox/common/auth"
	"github.com/jbenzshawel/go-sandbox/common/rest"
	"github.com/jbenzshawel/go-sandbox/identity/app"
)

type HttpHandler struct {
	app          app.Application
	authProvider *auth.OIDCProvider
	healthCheck  *rest.HealthCheckHandler
}

func NewHttpHandler(application app.Application, nc *nats.Conn, authProvider *auth.OIDCProvider) *HttpHandler {
	return &HttpHandler{
		app:          application,
		authProvider: authProvider,
		healthCheck: rest.NewHealthCheckHandler(
			application.Logger,
			rest.GetDatabaseHealthCheck(application.DB()),
			rest.GetNatsHealthCheck(nc),
		),
	}
}

func (h *HttpHandler) HealthCheck(ctx *gin.Context) {
	h.healthCheck.Handler(ctx)
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
