package logger

import (
	"log/slog"
	"os"
	"strconv"
)

type SlogLogger struct {
	Logger *slog.Logger
}

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
