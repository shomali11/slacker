package main

import (
	"context"
	"log"
	"os"

	"github.com/shomali11/slacker/v2"
)

// Showcase the ability to define Cron Jobs with middleware

func main() {
	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))
	bot.AddCommand("ping", &slacker.CommandDefinition{
		Handler: func(ctx slacker.CommandContext) {
			ctx.Response().Reply("pong")
		},
	})

	bot.AddJobMiddleware(func(next slacker.JobHandler) slacker.JobHandler {
		return func(ctx slacker.JobContext) {
			ctx.Response().Post("#test", "Root Middleware!")
			next(ctx)
		}
	})

	// Run every minute
	bot.AddJob("*/1 * * * *", &slacker.JobDefinition{
		Description: "A cron job that runs every minute",
		Handler: func(jobCtx slacker.JobContext) {
			jobCtx.Response().Post("#test", "Hello!")
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
