package bootstrap

import (
	"context"
	"fmt"
	"memory-db/internal/configuration"
	"memory-db/internal/database"
	"memory-db/internal/database/compute"
	"memory-db/internal/database/storage"
	"memory-db/internal/network"
	"os"

	"go.uber.org/zap"
)

var ConfigFileName = os.Getenv("CONFIG_FILE_NAME")

type App struct {
	logger *zap.Logger
	server *network.TCPServer
}

func NewApp() (*App, error) {
	var err error

	cfg := &configuration.Config{}
	if ConfigFileName != "" {
		cfg, err = configuration.Load(ConfigFileName)
		if err != nil {
			return nil, err
		}
	}

	logger, err := provideLogger(cfg.Logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	storageEngine := storage.NewMemoryStorage()
	db, err := provideDatabase(storageEngine, logger)
	if err != nil {
		return nil, err
	}

	handlerFunc := func(ctx context.Context, req []byte) []byte {
		return []byte(db.HandleQuery(string(req)))
	}
	server, err := provideTCPServer(cfg.Network, handlerFunc, logger)
	if err != nil {
		return nil, err
	}

	return &App{
		logger: logger,
		server: server,
	}, nil
}

func (app *App) RunTCPServer(ctx context.Context) error {
	return app.server.Serve(ctx)
}

func provideDatabase(storageEngine storage.Engine, logger *zap.Logger) (*database.Database, error) {
	analyzer := compute.NewAnalyzer(logger)
	parser := compute.NewParser(logger)

	computeEngine, err := compute.NewCompute(analyzer, parser, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create compute engine: %w", err)
	}

	return database.NewDatabase(computeEngine, storageEngine, logger), nil
}
