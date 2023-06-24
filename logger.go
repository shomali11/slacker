package slacker

import (
	"log"
	"os"
)

type Logger interface {
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
}

type builtinLogger struct {
	debugMode bool
	logger    *log.Logger
}

func newBuiltinLogger(debugMode bool) *builtinLogger {
	return &builtinLogger{
		debugMode: debugMode,
		logger:    log.New(os.Stdout, "", log.LstdFlags),
	}
}

func (l *builtinLogger) Info(args ...interface{}) {
	l.logger.Println(args...)
}

func (l *builtinLogger) Infof(format string, args ...interface{}) {
	l.logger.Printf(format, args...)
}

func (l *builtinLogger) Debug(args ...interface{}) {
	if l.debugMode {
		l.logger.Println(args...)
	}
}

func (l *builtinLogger) Debugf(format string, args ...interface{}) {
	if l.debugMode {
		l.logger.Printf(format, args...)
	}
}

func (l *builtinLogger) Error(args ...interface{}) {
	l.logger.Println(args...)
}

func (l *builtinLogger) Errorf(format string, args ...interface{}) {
	l.logger.Printf(format, args...)
}
