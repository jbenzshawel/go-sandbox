package rest

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/jbenzshawel/go-sandbox/common/cerror"
	crest "github.com/jbenzshawel/go-sandbox/common/rest"
	"github.com/jbenzshawel/go-sandbox/identity/app/command"
	"github.com/jbenzshawel/go-sandbox/identity/app/query"
	"github.com/jbenzshawel/go-sandbox/identity/domain/token"
)

type sendVerificationRequest struct {
	VerificationType string `json:"verificationType"`
}

func (h *HttpHandler) SendVerification(ctx *gin.Context) {
	userUUID, ok := h.parseUUIDParam(ctx)
	if !ok {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	var r sendVerificationRequest
	var err error
	if err = ctx.BindJSON(&r); err != nil {
		h.app.Logger.Error(err)
		ctx.JSON(http.StatusBadRequest, cerror.NewValidationError("invalid JSON", nil))
		return
	}

	verificationType, ok := token.ParseVerificationType(r.VerificationType)
	if !ok {
		ctx.JSON(http.StatusBadRequest, cerror.NewValidationError("bad request",
			map[string]string{"verificationType": fmt.Sprintf("verification type %s is not supported", r.VerificationType)}))
		return
	}

	u, err := h.app.Queries.UserByUUID.Handle(ctx, query.UserByUUID{UUID: userUUID})
	if err != nil {
		crest.HandleErrorResponse(ctx, err)
		return
	}
	if u == nil {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	switch verificationType {
	case token.Email:
		err = h.app.Commands.SendVerificationEmail.Handle(ctx, command.SendVerificationEmail{
			UserUUID:  userUUID,
			FirstName: u.FirstName(),
			Email:     u.Email(),
		})
		// TODO: support additional verification methods
	}

	if err != nil {
		crest.HandleErrorResponse(ctx, err)
	} else {
		ctx.AbortWithStatus(http.StatusOK)
	}
}
