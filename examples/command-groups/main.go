package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/shomali11/slacker/v2"
)

func main() {
	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))
	bot.AddCommand(&slacker.CommandDefinition{
		Command: "ping",
		Handler: func(ctx *slacker.CommandContext) {
			ctx.Response().Reply("pong")
		},
	})

	bot.AddCommandMiddleware(LoggingCommandMiddleware())
	bot.AddCommandMiddleware(func(next slacker.CommandHandler) slacker.CommandHandler {
		return func(ctx *slacker.CommandContext) {
			ctx.Response().Reply("Root Middleware!")
			next(ctx)
		}
	})

	group := bot.AddCommandGroup("cool")
	group.AddMiddleware(func(next slacker.CommandHandler) slacker.CommandHandler {
		return func(ctx *slacker.CommandContext) {
			ctx.Response().Reply("Group Middleware!")
			next(ctx)
		}
	})

	commandMiddleware := func(next slacker.CommandHandler) slacker.CommandHandler {
		return func(ctx *slacker.CommandContext) {
			ctx.Response().Reply("Command Middleware!")
			next(ctx)
		}
	}

	group.AddCommand(&slacker.CommandDefinition{
		Command:     "weather",
		Description: "Find me a cool weather",
		Examples:    []string{"cool weather"},
		Middlewares: []slacker.CommandMiddlewareHandler{commandMiddleware},
		Handler: func(ctx *slacker.CommandContext) {
			ctx.Response().Reply("San Francisco")
		},
	})

	group.AddCommand(&slacker.CommandDefinition{
		Command:     "person",
		Description: "Find me a cool person",
		Examples:    []string{"cool person"},
		Handler: func(ctx *slacker.CommandContext) {
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

func LoggingCommandMiddleware() slacker.CommandMiddlewareHandler {
	return func(next slacker.CommandHandler) slacker.CommandHandler {
		return func(ctx *slacker.CommandContext) {
			fmt.Printf(
				"%s executed \"%s\" with parameters %v in channel %s\n",
				ctx.Event().UserID,
				ctx.Definition().Command,
				ctx.Request().Properties(),
				ctx.Event().Channel.ID,
			)
			next(ctx)
		}
	}
}
