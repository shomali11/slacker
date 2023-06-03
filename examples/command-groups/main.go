package main

import (
	"context"
	"log"
	"os"

	"github.com/shomali11/slacker/v2"
)

func main() {
	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))
	bot.AddCommand("ping", &slacker.CommandDefinition{
		Handler: func(ctx slacker.CommandContext) {
			ctx.Response().Reply("pong")
		},
	})

	bot.AddCommandMiddleware(slacker.LoggingCommandMiddleware())
	bot.AddCommandMiddleware(func(next slacker.CommandHandler) slacker.CommandHandler {
		return func(ctx slacker.CommandContext) {
			ctx.Response().Reply("Root Middleware!")
			next(ctx)
		}
	})

	group := bot.AddCommandGroup("cool")
	group.AddMiddleware(func(next slacker.CommandHandler) slacker.CommandHandler {
		return func(ctx slacker.CommandContext) {
			ctx.Response().Reply("Group Middleware!")
			next(ctx)
		}
	})

	commandMiddleware := func(next slacker.CommandHandler) slacker.CommandHandler {
		return func(ctx slacker.CommandContext) {
			ctx.Response().Reply("Command Middleware!")
			next(ctx)
		}
	}

	group.AddCommand("weather", &slacker.CommandDefinition{
		Description: "Find me a cool weather",
		Examples:    []string{"cool weather"},
		Middlewares: []slacker.CommandMiddlewareHandler{commandMiddleware},
		Handler: func(ctx slacker.CommandContext) {
			ctx.Response().Reply("San Francisco")
		},
	})

	group.AddCommand("person", &slacker.CommandDefinition{
		Description: "Find me a cool person",
		Examples:    []string{"cool person"},
		Handler: func(ctx slacker.CommandContext) {
			ctx.Response().Reply("Dwayne Johnson")
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
