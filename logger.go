package main

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// newLogger builds the global zap logger. Level is controlled by LOG_LEVEL:
// debug, info, warn, error (default info).
func newLogger() (*zap.Logger, error) {
	lvl := envLogLevel()

	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(lvl)
	cfg.EncoderConfig.TimeKey = "ts"
	cfg.EncoderConfig.MessageKey = "msg"
	cfg.Development = false
	return cfg.Build()
}

// envLogLevel resolves LOG_LEVEL; defaults to info.
func envLogLevel() zapcore.Level {
	switch os.Getenv("LOG_LEVEL") {
	case "debug":
		return zapcore.DebugLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}
