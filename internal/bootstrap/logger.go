package bootstrap

import (
	"errors"
	"memory-db/internal/configuration"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	DebugLevel = "debug"
	InfoLevel  = "info"
	WarnLevel  = "warn"
	ErrorLevel = "error"
)

var loggerLevels = map[string]zapcore.Level{
	DebugLevel: zapcore.DebugLevel,
	InfoLevel:  zapcore.InfoLevel,
	WarnLevel:  zapcore.WarnLevel,
	ErrorLevel: zapcore.ErrorLevel,
}

const defaultLoggerLevel = zapcore.InfoLevel
const defaultLoggerOutput = "stdout"
const defaultLoggerEncoding = "json"

func provideLogger(cfg *configuration.LoggerConfig) (*zap.Logger, error) {
	loggerLevel := defaultLoggerLevel

	if cfg != nil {
		var ok bool
		if loggerLevel, ok = loggerLevels[cfg.Level]; !ok {
			return nil, errors.New("invalid log level")
		}
	}

	loggerCfg := zap.Config{
		Level:         zap.NewAtomicLevelAt(loggerLevel),
		OutputPaths:   []string{defaultLoggerOutput},
		Encoding:      defaultLoggerEncoding,
		EncoderConfig: zap.NewProductionEncoderConfig(),
	}

	return loggerCfg.Build()
}
