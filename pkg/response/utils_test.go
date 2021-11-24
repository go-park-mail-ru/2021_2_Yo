package response

import (
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUtils(t *testing.T) {
	err := errors.New("")
	require.Error(t, err)
}
