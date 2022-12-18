package cerror

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ValidationError struct {
	err         error
	Message     string
	FieldErrors map[string]string
}

func (r ValidationError) Error() string {
	if len(r.FieldErrors) == 0 {
		return fmt.Sprintf("err %+v", r.err)
	}

	return fmt.Sprintf("err %+v, fieldErrors: %+v", r.err, r.FieldErrors)
}

func NewValidationError(errorMsg string, fieldErrors map[string]string) ValidationError {
	return ValidationError{err: errors.New(errorMsg), Message: errorMsg, FieldErrors: fieldErrors}
}

func HandleValidationError(c *gin.Context, err error) {
	if _, ok := err.(ValidationError); ok {
		c.IndentedJSON(http.StatusBadRequest, err)
	} else {
		c.IndentedJSON(http.StatusInternalServerError, "internal server error")
	}
}
