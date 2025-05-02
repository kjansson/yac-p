// Package logger provides a simple logging interface for yac-p. Implements the types.Logger interface.
package logger

import (
	"fmt"
	"log/slog"
	"os"
)

type SlogLogger struct {
	Logger         *slog.Logger // slog logger instance
	LogFormat      string       // Format of the log output (json or text)
	LogDestination *os.File     // Destination of the log output
	LogLevel       string       // Logging level (debug, info, warn, error)
}

// Init initializes the logger and determines the log level
func NewLogger(LogDestination *os.File, LogFormat string, Debug bool) (*SlogLogger, error) {

	logger := &SlogLogger{}
	logOpts := &slog.HandlerOptions{}
	if Debug {
		logOpts.Level = slog.LevelDebug
	} else {
		logOpts.Level = slog.LevelInfo
	}

	var destination *os.File
	if LogDestination != nil {
		destination = LogDestination
	} else {
		destination = os.Stdout
	}

	if LogFormat != "" {
		if LogFormat != "json" && LogFormat != "JSON" && LogFormat != "text" && LogFormat != "TEXT" {
			return nil, fmt.Errorf("invalid log format: %s", LogFormat)
		}
	}

	if LogFormat == "json" || LogFormat == "JSON" {
		logger.Logger = slog.New(slog.NewJSONHandler(destination, logOpts))
	} else {
		logger.Logger = slog.New(slog.NewTextHandler(destination, logOpts))
	}
	return logger, nil
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
