package slacker

import "fmt"

// logger is a logger interface compatible with both stdlib and some
// 3rd party loggers. Since slack-go doesn't exposes this interface
// it is redeclared here.
type logger interface {
	Output(int, string) error
}

//A default nullLogger instance
var NullLogger = nullLogger{}

// Declares a simple no op logger to use as default.
type nullLogger struct{}

func (n nullLogger) Output(int, string) error {
	return nil
}

// ilogger represents the internal logging api we use.
type ilogger interface {
	logger
	Print(...interface{})
	Printf(string, ...interface{})
	Println(...interface{})
}

type Debug interface {
	Debug() bool

	// Debugf print a formatted debug line.
	Debugf(format string, v ...interface{})
	// Debugln print a debug line.
	Debugln(v ...interface{})
}

// internalLog implements the additional methods used by our internal logging.
type internalLog struct {
	logger
}

// Println replicates the behaviour of the standard logger.
func (t internalLog) Println(v ...interface{}) {
	t.Output(2, fmt.Sprintln(v...))
}

// Printf replicates the behaviour of the standard logger.
func (t internalLog) Printf(format string, v ...interface{}) {
	t.Output(2, fmt.Sprintf(format, v...))
}

// Print replicates the behaviour of the standard logger.
func (t internalLog) Print(v ...interface{}) {
	t.Output(2, fmt.Sprint(v...))
}

type discard struct{}

func (t discard) Debug() bool {
	return false
}

// Debugf print a formatted debug line.
func (t discard) Debugf(format string, v ...interface{}) {}

// Debugln print a debug line.
func (t discard) Debugln(v ...interface{}) {}
