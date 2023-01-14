package rest

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/jbenzshawel/go-sandbox/common/cerror"
	crest "github.com/jbenzshawel/go-sandbox/common/rest"
	"github.com/jbenzshawel/go-sandbox/identity/app/command"
	"github.com/jbenzshawel/go-sandbox/identity/domain/token"
)

type verifyUserRequest struct {
	Code             string `json:"code"`
	VerificationType string `json:"verificationType"`
}

func (s *HttpHandler) VerifyUser(ctx *gin.Context) {
	queryParam := ctx.Param("uuid")
	if queryParam == "" {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	userUUID, err := uuid.Parse(queryParam)
	if err != nil {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	var r verifyUserRequest
	if err = ctx.BindJSON(&r); err != nil {
		s.app.Logger.Error(err)
		ctx.JSON(http.StatusBadRequest, cerror.NewValidationError("invalid JSON", nil))
		return
	}

	verificationType, ok := token.ParseVerificationType(r.VerificationType)
	if !ok {
		ctx.JSON(http.StatusBadRequest, cerror.NewValidationError("bad request",
			map[string]string{"verificationType": fmt.Sprintf("verification type %s is not supported", r.VerificationType)}))
		return
	}

	switch verificationType {
	case token.Email:
		err = s.app.Commands.VerifyEmail.Handle(ctx, command.VerifyEmail{
			UserId: userUUID,
			Code:   r.Code,
		})
		// TODO: support additional verification methods
	}

	if err != nil {
		crest.HandleErrorResponse(ctx, err)
	} else {
		ctx.AbortWithStatus(http.StatusOK)
	}
}
