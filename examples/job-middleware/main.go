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
	bot.AddCommand(&slacker.CommandDefinition{
		Command: "ping",
		Handler: func(ctx *slacker.CommandContext) {
			ctx.Response().Reply("pong")
		},
	})

	bot.AddJobMiddleware(LoggingJobMiddleware())

	// ┌───────────── minute (0 - 59)
	// │ ┌───────────── hour (0 - 23)
	// │ │ ┌───────────── day of the month (1 - 31)
	// │ │ │ ┌───────────── month (1 - 12)
	// │ │ │ │ ┌───────────── day of the week (0 - 6) (Sunday to Saturday)
	// │ │ │ │ │
	// │ │ │ │ │
	// │ │ │ │ │
	// * * * * * (cron expression)

	// Run every minute
	bot.AddJob(&slacker.JobDefinition{
		CronExpression: "*/1 * * * *",
		Name:           "SomeJob",
		Description:    "A cron job that runs every minute",
		Handler: func(ctx *slacker.JobContext) {
			ctx.Response().Post("#test", "Hello!")
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

func LoggingJobMiddleware() slacker.JobMiddlewareHandler {
	return func(next slacker.JobHandler) slacker.JobHandler {
		return func(ctx *slacker.JobContext) {
			ctx.Logger().Info(
				"job middleware before",
				"job_name", ctx.Definition().Name,
			)
			next(ctx)
			ctx.Logger().Info(
				"job middleware after",
				"job_name", ctx.Definition().Name,
			)
		}
	}
}
