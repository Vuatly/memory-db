package network

import (
	"context"
	"net"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestTCPServer(t *testing.T) {
	t.Parallel()

	request := "hello server"
	response := "hello client"

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	maxMessageSize := 4096
	handler := func(ctx context.Context, buffer []byte) []byte {
		require.True(t, reflect.DeepEqual([]byte(request), buffer))
		return []byte(response)
	}

	server, err := NewTCPServer(":5444", maxMessageSize, 5, handler, time.Minute, zap.NewNop())
	require.NoError(t, err)

	go func() {
		require.NoError(t, server.Serve(ctx))
	}()

	connection, err := net.Dial("tcp", "localhost:5444")
	require.NoError(t, err)

	_, err = connection.Write([]byte(request))
	require.NoError(t, err)

	buffer := make([]byte, maxMessageSize)
	count, err := connection.Read(buffer)
	require.NoError(t, err)
	require.True(t, reflect.DeepEqual([]byte(response), buffer[:count]))
}
