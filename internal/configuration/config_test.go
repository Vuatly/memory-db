package configuration

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestLoadNonExistentFile(t *testing.T) {
	t.Parallel()

	cfg, err := Load("test_data/non_existent_config.yaml")
	require.Error(t, err)
	require.Nil(t, cfg)
}

func TestLoadWithEmptyFilename(t *testing.T) {
	t.Parallel()

	cfg, err := Load("")
	require.NoError(t, err)
	require.NotNil(t, cfg)
}

func TestLoadEmptyConfig(t *testing.T) {
	t.Parallel()

	cfg, err := Load("test_data/empty_config.yaml")
	require.NoError(t, err)
	require.NotNil(t, cfg)
}

func TestLoadConfig(t *testing.T) {
	t.Parallel()

	cfg, err := Load("test_data/config.yaml")
	require.NoError(t, err)

	require.Equal(t, "127.0.0.1:5444", cfg.Network.Address)
	require.Equal(t, 100, cfg.Network.MaxConnections)
	require.Equal(t, 4096, cfg.Network.MaxMessageSize)
	require.Equal(t, time.Second*15, cfg.Network.IdleTimeout)

	require.Equal(t, "info", cfg.Logger.Level)
}
