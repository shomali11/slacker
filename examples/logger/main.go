package main

import (
	"context"
	"log"
	"os"

	"github.com/shomali11/slacker/v2"
)

// Showcasing the ability to pass your own logger

func main() {
	logger := newLogger()

	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"), slacker.WithLogger(logger))

	definition := &slacker.CommandDefinition{
		Command:     "ping",
		Description: "Ping!",
		Handler: func(ctx *slacker.CommandContext) {
			ctx.Response().Reply("pong")
		},
	}

	bot.AddCommand(definition)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

type MyLogger struct {
	debugMode bool
	logger    *log.Logger
}

func newLogger() *MyLogger {
	return &MyLogger{
		logger: log.New(os.Stdout, "something ", log.LstdFlags|log.Lshortfile|log.Lmsgprefix),
	}
}

func (l *MyLogger) Info(args ...interface{}) {
	l.logger.Println(args...)
}

func (l *MyLogger) Infof(format string, args ...interface{}) {
	l.logger.Printf(format, args...)
}

func (l *MyLogger) Debug(args ...interface{}) {
	if l.debugMode {
		l.logger.Println(args...)
	}
}

func (l *MyLogger) Debugf(format string, args ...interface{}) {
	if l.debugMode {
		l.logger.Printf(format, args...)
	}
}

func (l *MyLogger) Error(args ...interface{}) {
	l.logger.Println(args...)
}

func (l *MyLogger) Errorf(format string, args ...interface{}) {
	l.logger.Printf(format, args...)
}
