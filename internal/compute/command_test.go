package compute

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCommandFromString(t *testing.T) {
	t.Parallel()

	require.Equal(t, GetCommand, CommandFromString(GetCommandString))
	require.Equal(t, SetCommand, CommandFromString(SetCommandString))
	require.Equal(t, DeleteCommand, CommandFromString(DeleteCommandString))
}
