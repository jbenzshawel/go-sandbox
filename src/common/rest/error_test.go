package rest

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jbenzshawel/go-sandbox/common/cerror"
)

func TestHandleErrorResponse(t *testing.T) {
	testCases := []struct {
		name               string
		err                error
		expectedStatusCode int
	}{
		{
			name: "validation error",
			err: cerror.NewValidationError(
				"Invalid request",
				map[string]string{"fieldName1": "fieldError1", "fieldName2": "fieldError2"},
			),
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "unexpected error",
			err:                errors.New("oh no - something bad has happened"),
			expectedStatusCode: http.StatusInternalServerError,
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)

			HandleErrorResponse(ctx, tc.err)

			resp := w.Result()
			assert.Equal(t, tc.expectedStatusCode, resp.StatusCode)

			body, _ := io.ReadAll(resp.Body)
			if tc.expectedStatusCode == http.StatusBadRequest {
				var respJson cerror.ValidationError
				err := json.Unmarshal(body, &respJson)
				require.NoError(t, err)

				origErr, ok := tc.err.(cerror.ValidationError)
				require.True(t, ok)
				assert.Equal(t, origErr.Message, respJson.Message)
				assert.Equal(t, origErr.FieldErrors, respJson.FieldErrors)
			} else {
				assert.Len(t, body, 0)
			}
		})
	}
}
