package storage

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMemoryStorage_Get(t *testing.T) {
	t.Parallel()
	storage := NewMemoryStorage()

	val, err := storage.Get("key")
	require.NoError(t, err)
	require.Equal(t, "", val)

	storage.hashMap["key"] = "value"
	val, err = storage.Get("key")
	require.NoError(t, err)
	require.Equal(t, "value", val)
}

func TestMemoryStorage_Set(t *testing.T) {
	t.Parallel()
	storage := NewMemoryStorage()
	require.Equal(t, "", storage.hashMap["key"])

	err := storage.Set("key", "value")
	require.NoError(t, err)
	require.Equal(t, "value", storage.hashMap["key"])
}

func TestMemoryStorage_Delete(t *testing.T) {
	t.Parallel()
	storage := NewMemoryStorage()
	storage.hashMap["key"] = "value"
	require.Equal(t, "value", storage.hashMap["key"])

	err := storage.Delete("key")
	require.NoError(t, err)
	require.Equal(t, "", storage.hashMap["key"])
}
