package logger

// Package logger provides a simple logging interface for yac-p

import (
	"log/slog"
	"os"
	"strconv"
)

type SlogLogger struct {
	Logger *slog.Logger
}

// Init initializes the logger and determines the log level
func (l *SlogLogger) Init() error {

	debugEnv := os.Getenv("DEBUG")
	debug, _ := strconv.ParseBool(debugEnv)

	logOpts := &slog.HandlerOptions{}
	if debug {
		logOpts.Level = slog.LevelDebug
	} else {
		logOpts.Level = slog.LevelInfo
	}

	l.Logger = slog.New(slog.NewTextHandler(os.Stdout, logOpts))
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
