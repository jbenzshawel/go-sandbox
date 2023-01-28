package rest

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/jbenzshawel/go-sandbox/common/auth"
	crest "github.com/jbenzshawel/go-sandbox/common/rest"
	"github.com/jbenzshawel/go-sandbox/identity/app/query"
)

type getUserResponse struct {
	ID            int       `json:"id"`
	UUID          uuid.UUID `json:"uuid"`
	FirstName     string    `json:"firstName"`
	LastName      string    `json:"lastName"`
	Email         string    `json:"email"`
	EmailVerified bool      `json:"emailVerified"`
	Enabled       bool      `json:"enabled"`
	CreatedAt     time.Time `json:"createdAt"`
	LastUpdatedAt time.Time `json:"lastUpdatedAt"`
}

func (h *HttpHandler) GetUserByUUID(ctx *gin.Context) {
	h.authenticate(ctx, h.getUserByUUID)
}

func (h *HttpHandler) getUserByUUID(ctx *gin.Context, authUser *auth.User) {
	userUUID, ok := h.parseUUIDParam(ctx)
	if !ok {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	if authUser == nil || authUser.UserUUID != userUUID {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	user, err := h.app.Queries.UserByUUID.Handle(ctx, query.UserByUUID{UUID: userUUID})
	if err != nil {
		crest.HandleErrorResponse(ctx, err)
		return
	}
	if user == nil {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	ctx.JSON(http.StatusCreated, &getUserResponse{
		ID:            user.ID(),
		UUID:          user.UUID(),
		FirstName:     user.FirstName(),
		LastName:      user.LastName(),
		Email:         user.Email(),
		EmailVerified: user.EmailVerified(),
		Enabled:       user.Enabled(),
		CreatedAt:     user.CreatedAt(),
		LastUpdatedAt: user.LastUpdatedAt(),
	})
}
