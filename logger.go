package slacker

import (
	"log/slog"
	"os"
)

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

type builtinLogger struct {
	debugMode bool
	logger    *slog.Logger
}

func newBuiltinLogger(debugMode bool) *builtinLogger {
	opts := &slog.HandlerOptions{
		Level: getLogLevel(debugMode),
	}

	return &builtinLogger{
		debugMode: debugMode,
		logger:    slog.New(slog.NewJSONHandler(os.Stdout, opts)),
	}
}

func (l *builtinLogger) Info(msg string, args ...any) {
	l.logger.Info(msg, args...)
}

func (l *builtinLogger) Debug(msg string, args ...any) {
	l.logger.Debug(msg, args...)
}

func (l *builtinLogger) Warn(msg string, args ...any) {
	l.logger.Warn(msg, args...)
}

func (l *builtinLogger) Error(msg string, args ...any) {
	l.logger.Error(msg, args...)
}

func getLogLevel(isDebugMode bool) slog.Level {
	if isDebugMode {
		return slog.LevelDebug
	}
	return slog.LevelInfo
}
