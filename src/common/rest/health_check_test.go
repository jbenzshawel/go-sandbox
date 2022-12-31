package rest

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jbenzshawel/go-sandbox/common/database"
)

type dbHealthCheckTestCase struct {
	name       string
	dbProvider database.DbProvider
	success    bool
	err        error
}

func TestGetDatabaseHealthCheck(t *testing.T) {
	testCases := []dbHealthCheckTestCase{
		{
			name: "db provider error",
			dbProvider: func() (*sql.DB, error) {
				return nil, errors.New("db provider error")
			},
			success: false,
			err:     errors.New("db provider error"),
		},
		{
			name: "db ping error",
			dbProvider: func() (*sql.DB, error) {
				return sql.Open("postgres", "")
			},
			success: false,
			err:     errors.New("pq: SSL is not enabled on the server"),
		},
	}
	if connectionString, ok := os.LookupEnv("IDENTITY_POSTGRES"); ok {
		testCases = append(testCases, dbHealthCheckTestCase{
			name: "db health check success",
			dbProvider: func() (*sql.DB, error) {
				return sql.Open("postgres", connectionString)
			},
			success: true,
			err:     nil,
		})
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			check := GetDatabaseHealthCheck(tc.dbProvider)
			require.NotNil(t, check)

			success, name, err := check()
			if tc.err == nil {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				assert.Equal(t, tc.err.Error(), err.Error())
			}

			assert.Equal(t, tc.success, success)
			assert.Equal(t, "database", name)
		})
	}
}

func TestHealthCheckHandler(t *testing.T) {
	testCases := []struct {
		name   string
		checks []HealthCheckTask
		errors []string
	}{
		{
			name: "success with no additional checks",
		},
		{
			name: "success with additional checks",
			checks: []HealthCheckTask{
				func() (bool, string, error) {
					return true, "database", nil
				},
				func() (bool, string, error) {
					return true, "nats", nil
				},
			},
		},
		{
			name: "check fails",
			checks: []HealthCheckTask{
				func() (bool, string, error) {
					return true, "database", nil
				},
				func() (bool, string, error) {
					return false, "nats", errors.New("nats connection failed")
				},
			},
			errors: []string{"nats health check failed"},
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			testLogger, hook := test.NewNullLogger()
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)

			healthCheck := NewHealthCheckHandler(logrus.NewEntry(testLogger), tc.checks...)
			healthCheck.Handler(ctx)

			resp := w.Result()
			body, _ := io.ReadAll(resp.Body)
			var healthCheckResp healthCheckResponse
			err := json.Unmarshal(body, &healthCheckResp)
			require.NoError(t, err)
			if len(tc.errors) == 0 {
				assert.Equal(t, http.StatusOK, resp.StatusCode)
				assert.Equal(t, "available", healthCheckResp.Status)
				assert.Nil(t, healthCheckResp.Errors)
			} else {
				assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
				assert.Equal(t, "unavailable", healthCheckResp.Status)
				assert.Equal(t, tc.errors, healthCheckResp.Errors)
				require.NotNil(t, hook)
				require.Len(t, hook.Entries, 1)
				logEntry := hook.Entries[0]
				require.NotNil(t, logEntry)
				assert.Equal(t, tc.errors[0], logEntry.Message)
				assert.Equal(t, errors.New("nats connection failed"), logEntry.Data["error"])
			}
		})
	}
}
