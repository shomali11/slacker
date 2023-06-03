package main

import (
	"context"
	"log"

	//"os"

	"github.com/shomali11/slacker/v2"
)

// Showcase the ability to define Cron Jobs with middleware

func main() {
	bot := slacker.NewClient("xoxb-13360094916-2243791173942-G5ns4blGnmRL2EVdFvmwdn2r", "xapp-1-A027JMM1RV2-5374333520369-a8e99a2f622940d7cdc410f5b7129da2896e644eeb9a66d0e449ba19af3b00ff")
	//bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))
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
