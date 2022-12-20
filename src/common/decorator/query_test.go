package decorator

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

type stubQuery struct {
	stubProp string
}

type stubQueryHandler struct {
	stubHandler func(ctx context.Context, query stubQuery) (string, error)
}

func (h *stubQueryHandler) Handle(ctx context.Context, query stubQuery) (string, error) {
	if h.stubHandler != nil {
		return h.stubHandler(ctx, query)
	}
	return "", nil
}

func TestApplyQueryDecorators(t *testing.T) {
	handlerCalled := false
	stubHandler := &stubQueryHandler{
		stubHandler: func(ctx context.Context, query stubQuery) (string, error) {
			handlerCalled = true
			assert.Equal(t, "Test", query.stubProp)
			return "result", nil
		},
	}

	testLogger, hooks := test.NewNullLogger()

	commandDecorator := ApplyQueryDecorators[stubQuery, string](stubHandler, logrus.NewEntry(testLogger))
	require.NotNil(t, commandDecorator)

	res, err := commandDecorator.Handle(context.Background(), stubQuery{stubProp: "Test"})
	require.Nil(t, err)
	assert.True(t, handlerCalled)
	assert.Equal(t, "result", res)

	require.Len(t, hooks.Entries, 1)
	entry := hooks.Entries[0]
	assert.Equal(t, "Query executed successfully", entry.Message)
	assert.Equal(t, logrus.InfoLevel, entry.Level)

	require.Len(t, entry.Data, 2)
	assert.Equal(t, "stubQuery", entry.Data["query"])
	assert.Equal(t, `decorator.stubQuery{stubProp:"Test"}`, entry.Data["query_body"])
}

func TestApplyQueryDecorators_LogsErrors(t *testing.T) {
	handlerCalled := false
	fakeError := errors.New("fake query error")
	stubHandler := &stubQueryHandler{
		stubHandler: func(ctx context.Context, query stubQuery) (string, error) {
			handlerCalled = true
			assert.Equal(t, "Test", query.stubProp)
			return "", fakeError
		},
	}

	testLogger, hooks := test.NewNullLogger()

	commandDecorator := ApplyQueryDecorators[stubQuery, string](stubHandler, logrus.NewEntry(testLogger))
	require.NotNil(t, commandDecorator)

	res, err := commandDecorator.Handle(context.Background(), stubQuery{stubProp: "Test"})
	require.Error(t, fakeError, err)
	assert.True(t, handlerCalled)
	assert.Equal(t, "", res)

	require.Len(t, hooks.Entries, 1)
	entry := hooks.Entries[0]
	assert.Equal(t, "Failed to execute query", entry.Message)
	assert.Equal(t, logrus.ErrorLevel, entry.Level)

	require.Len(t, entry.Data, 3)
	assert.Equal(t, "stubQuery", entry.Data["query"])
	assert.Equal(t, `decorator.stubQuery{stubProp:"Test"}`, entry.Data["query_body"])
	assert.Equal(t, fakeError, entry.Data["error"])
}
