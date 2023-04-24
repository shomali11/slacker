package main

import (
	"context"
	"log"
	"os"

	"github.com/shomali11/slacker"
)

func main() {
	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))
	bot.AddCommand("ping", &slacker.CommandDefinition{
		Handler: func(botCtx slacker.CommandContext) {
			botCtx.Response().Reply("pong")
		},
	})

	bot.AddMiddleware(slacker.LoggingMiddleware())
	bot.AddMiddleware(func(next slacker.CommandHandler) slacker.CommandHandler {
		return func(botCtx slacker.CommandContext) {
			botCtx.Response().Reply("Root Middleware!")
			next(botCtx)
		}
	})

	group := bot.AddGroup("cool")
	group.AddMiddleware(func(next slacker.CommandHandler) slacker.CommandHandler {
		return func(botCtx slacker.CommandContext) {
			botCtx.Response().Reply("Group Middleware!")
			next(botCtx)
		}
	})

	commandMiddleware := func(next slacker.CommandHandler) slacker.CommandHandler {
		return func(botCtx slacker.CommandContext) {
			botCtx.Response().Reply("Command Middleware!")
			next(botCtx)
		}
	}

	group.AddCommand("weather", &slacker.CommandDefinition{
		Description: "Find me a cool weather",
		Examples:    []string{"cool weather"},
		Middlewares: []slacker.MiddlewareHandler{commandMiddleware},
		Handler: func(botCtx slacker.CommandContext) {
			botCtx.Response().Reply("San Francisco")
		},
	})

	group.AddCommand("person", &slacker.CommandDefinition{
		Description: "Find me a cool person",
		Examples:    []string{"cool person"},
		Handler: func(botCtx slacker.CommandContext) {
			botCtx.Response().Reply("Dwayne Johnson")
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
