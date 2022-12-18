package decorator

import (
	"context"
	"errors"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type stubCommand struct {
	stubProp string
}

type stubCommandHandler struct {
	stubHandler func(ctx context.Context, command stubCommand) error
}

func (h *stubCommandHandler) Handle(ctx context.Context, cmd stubCommand) error {
	if h.stubHandler != nil {
		return h.stubHandler(ctx, cmd)
	}
	return nil
}

func TestApplyCommandDecorators(t *testing.T) {
	handlerCalled := false
	stubHandler := &stubCommandHandler{
		stubHandler: func(ctx context.Context, command stubCommand) error {
			handlerCalled = true
			assert.Equal(t, "Test", command.stubProp)
			return nil
		},
	}

	testLogger, hooks := test.NewNullLogger()

	commandDecorator := ApplyCommandDecorators[stubCommand](stubHandler, logrus.NewEntry(testLogger))
	require.NotNil(t, commandDecorator)

	err := commandDecorator.Handle(context.Background(), stubCommand{stubProp: "Test"})
	require.Nil(t, err)
	assert.True(t, handlerCalled)

	require.Len(t, hooks.Entries, 1)
	entry := hooks.Entries[0]
	assert.Equal(t, "Command executed successfully", entry.Message)

	require.Len(t, entry.Data, 1)
	assert.Equal(t, "stubCommand", entry.Data["command"])
}

func TestApplyCommandDecorators_LogsError(t *testing.T) {
	handlerCalled := false
	stubHandler := &stubCommandHandler{
		stubHandler: func(ctx context.Context, command stubCommand) error {
			handlerCalled = true
			return errors.New("fake command error")
		},
	}

	testLogger, hooks := test.NewNullLogger()

	commandDecorator := ApplyCommandDecorators[stubCommand](stubHandler, logrus.NewEntry(testLogger))
	require.NotNil(t, commandDecorator)

	err := commandDecorator.Handle(context.Background(), stubCommand{stubProp: "Test"})
	require.NotNil(t, err)
	assert.Equal(t, "fake command error", err.Error())
	assert.True(t, handlerCalled)

	require.Len(t, hooks.Entries, 1)
	entry := hooks.Entries[0]
	assert.Equal(t, "Failed to execute command", entry.Message)

	require.Len(t, entry.Data, 2)
	assert.Equal(t, "stubCommand", entry.Data["command"])
	assert.Equal(t, errors.New("fake command error"), entry.Data["error"])
}
