package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/jbenzshawel/go-sandbox/common/cerror"
)

func HandleErrorResponse(ctx *gin.Context, err error) {
	if _, ok := err.(cerror.ValidationError); ok {
		ctx.JSON(http.StatusBadRequest, err)
	} else {
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}
}
