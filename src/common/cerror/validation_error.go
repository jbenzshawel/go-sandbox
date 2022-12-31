package cerror

import (
	"errors"
	"fmt"
)

type ValidationError struct {
	err         error
	Message     string            `json:"message"`
	FieldErrors map[string]string `json:"fieldErrors"`
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
