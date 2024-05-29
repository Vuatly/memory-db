package database

import (
	"fmt"
	"memory-db/internal/database/compute"

	"go.uber.org/zap"
)

type computeEngine interface {
	HandleQuery(query string) (compute.Query, error)
}

type storageEngine interface {
	Get(key string) (string, error)
	Set(key string, value string) error
	Delete(key string) error
}

type Database struct {
	computeEngine computeEngine
	storageLayer  storageEngine
	logger        *zap.Logger
}

func NewDatabase(compute computeEngine, storage storageEngine, logger *zap.Logger) *Database {
	return &Database{
		computeEngine: compute,
		storageLayer:  storage,
		logger:        logger,
	}
}

func (db *Database) HandleQuery(query string) string {
	db.logger.Debug("Handling query", zap.String("query", query))

	queryObj, err := db.computeEngine.HandleQuery(query)
	if err != nil {
		return fmt.Sprintf("[error] failed to handle query: %s; err: %v", query, err)
	}

	command := queryObj.GetCommand()
	args := queryObj.GetArguments()

	switch command {
	case compute.SetCommand:
		return db.handleSetCommand(args[0], args[1])
	case compute.GetCommand:
		return db.handleGetCommand(args[0])
	case compute.DeleteCommand:
		return db.handleDeleteCommand(args[0])
	}

	db.logger.Error("invalid command", zap.String("command", string(command)))
	return fmt.Sprintf("[error] failed to handle query: %s", query)
}

func (db *Database) handleSetCommand(key string, value string) string {
	err := db.storageLayer.Set(key, value)
	if err != nil {
		return fmt.Sprintf("[error] failed to set value: %s with key: %s", value, key)
	}
	return "[ok]"
}

func (db *Database) handleGetCommand(key string) string {
	val, err := db.storageLayer.Get(key)
	if err != nil {
		return fmt.Sprintf("[error] failed to get with key: %s", key)
	}

	if val == "" {
		return fmt.Sprintf("[ok] key not found")
	}

	return fmt.Sprintf("[ok] value: %s", val)
}

func (db *Database) handleDeleteCommand(key string) string {
	err := db.storageLayer.Delete(key)
	if err != nil {
		return fmt.Sprintf("[error] failed to delete key: %s", key)
	}
	return "[ok]"
}
