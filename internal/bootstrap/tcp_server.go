package bootstrap

import (
	"memory-db/internal/configuration"
	"memory-db/internal/network"
	"time"

	"go.uber.org/zap"
)

const defaultServerAddress = "0.0.0.0:5444"
const defaultIdleTimeout = time.Second * 15
const defaultMaxMessageSize = 4096 // 4KB
const defaultMaxConnections = 100

func provideTCPServer(cfg *configuration.NetworkConfig, handlerFunc network.RequestHandler, logger *zap.Logger) (*network.TCPServer, error) {
	if cfg == nil {
		cfg = &configuration.NetworkConfig{}
	}

	if cfg.Address == "" {
		cfg.Address = defaultServerAddress
	}
	if cfg.IdleTimeout == 0 {
		cfg.IdleTimeout = defaultIdleTimeout
	}
	if cfg.MaxMessageSize == 0 {
		cfg.MaxMessageSize = defaultMaxMessageSize
	}
	if cfg.MaxConnections == 0 {
		cfg.MaxConnections = defaultMaxConnections
	}

	return network.NewTCPServer(
		cfg.Address,
		cfg.MaxMessageSize,
		cfg.MaxConnections,
		handlerFunc,
		cfg.IdleTimeout,
		logger,
	)
}
