package cmd

import (
	"memory-db/internal"
	"memory-db/internal/compute"
	"memory-db/internal/storage"

	"go.uber.org/zap"
)

type App struct {
	Database *internal.Database
	Logger   *zap.Logger
}

func NewApp() (*App, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}

	analyzer := compute.NewAnalyzer(logger)
	parser := compute.NewParser(logger)

	computeEngine, err := compute.NewCompute(analyzer, parser, logger)
	if err != nil {
		return nil, err
	}

	storageEngine := storage.NewMemoryStorage()

	database := internal.NewDatabase(computeEngine, storageEngine, logger)

	return &App{
		Database: database,
		Logger:   logger,
	}, nil
}
