package cerror

import (
	"errors"

	pkgerrors "github.com/pkg/errors"

	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type stackTracer interface {
	StackTrace() pkgerrors.StackTrace
}

func TestCombineErrors(t *testing.T) {
	err := CombineErrors(nil)
	require.Nil(t, err)

	err = CombineErrors(nil, nil)
	require.Nil(t, err)

	err = CombineErrors(nil, errors.New("error 1"), nil, errors.New("error 2"), nil, nil, errors.New("error 3"))
	require.NotNil(t, err)
	assert.Equal(t, "error 3: error 2: error 1", err.Error())

	errStack, ok := err.(stackTracer)
	require.True(t, ok)
	assert.NotNil(t, errStack.StackTrace())
}
