package main

import (
	"context"
	"log"
	"log/slog"
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
	logger *slog.Logger
}

func newLogger() *MyLogger {
	return &MyLogger{
		logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}
}

func (l *MyLogger) Info(msg string, args ...any) {
	l.logger.Info(msg, args...)
}

func (l *MyLogger) Debug(msg string, args ...any) {
	l.logger.Debug(msg, args...)
}

func (l *MyLogger) Warn(msg string, args ...any) {
	l.logger.Warn(msg, args...)
}

func (l *MyLogger) Error(msg string, args ...any) {
	l.logger.Error(msg, args...)
}
