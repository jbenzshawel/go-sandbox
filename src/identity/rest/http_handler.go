package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/jbenzshawel/go-sandbox/common/auth"
	"github.com/jbenzshawel/go-sandbox/identity/app"
	"github.com/jbenzshawel/go-sandbox/identity/domain/user/permission"
)

type HttpHandler struct {
	app          app.Application
	authProvider *auth.OIDCProvider
	router       *gin.Engine
}

func NewHttpHandler(application app.Application, authProvider *auth.OIDCProvider) *HttpHandler {
	return &HttpHandler{
		app:          application,
		authProvider: authProvider,
	}
}

func (h *HttpHandler) Configure() *HttpHandler {
	h.router = gin.Default() // TODO: Update gin config for production
	h.router.POST("/identity-client/callback", h.OAuthCallback)

	h.router.GET("/health", h.app.HealthCheck.Handler)

	h.router.POST("/user", h.CreateUser)
	h.router.POST("/user/:uuid/send-verification", h.SendVerification)
	h.router.POST("/user/:uuid/verify", h.VerifyUser)
	h.router.GET("/user/:uuid", h.GetUserByUUID)
	h.router.GET("/user", h.GetUsers)

	return h
}

func (h *HttpHandler) Run(addr ...string) error {
	return h.router.Run(addr...)
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

func (h *HttpHandler) authorize(ctx *gin.Context, permit permission.Type, handler authenticatedHandler) {
	h.authenticate(ctx, func(ctx *gin.Context, authUser *auth.User) {
		ok, err := h.app.Services.PermissionService.HasPermission(authUser.UserUUID, permit)
		if err != nil {
			h.app.Logger.WithError(err).Error("failed to authorize user")
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if !ok {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		handler(ctx, authUser)
	})
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
