// Package logger provides a simple logging interface for yac-p
// It implements the types.Logger interface
package logger

import (
	"log/slog"
	"os"

	"github.com/kjansson/yac-p/pkg/types"
)

type SlogLogger struct {
	Logger *slog.Logger
}

// Init initializes the logger and determines the log level
func (l *SlogLogger) Init(config types.Config) error {

	logOpts := &slog.HandlerOptions{}
	if config.Debug {
		logOpts.Level = slog.LevelDebug
	} else {
		logOpts.Level = slog.LevelInfo
	}

	var destination *os.File
	if config.LogDestination != nil {
		destination = config.LogDestination
	} else {
		destination = os.Stdout
	}

	if config.LogFormat == "json" || config.LogFormat == "JSON" {
		l.Logger = slog.New(slog.NewJSONHandler(destination, logOpts))
	} else {
		l.Logger = slog.New(slog.NewTextHandler(destination, logOpts))
	}
	return nil
}

// Log accepts generic log entry components uses the slog package to log messages
func (l *SlogLogger) Log(level string, msg string, args ...any) {
	switch level {
	case "debug":
		l.Logger.Debug(msg, args...)
	case "info":
		l.Logger.Info(msg, args...)
	case "warn":
		l.Logger.Warn(msg, args...)
	case "error":
		l.Logger.Error(msg, args...)
	default:
		l.Logger.Info(msg, args...)
	}
}
