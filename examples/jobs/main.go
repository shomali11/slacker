package main

import (
	"context"
	"log"
	"os"

	"github.com/shomali11/slacker"
	"github.com/slack-go/slack"
)

// Showcase the ability to define Cron Jobs

func main() {
	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))
	bot.AddCommand("ping", &slacker.CommandDefinition{
		Handler: func(botCtx slacker.CommandContext) {
			botCtx.Response().Reply("pong")
		},
	})

	// Run every minute
	bot.AddJob("0 * * * * *", &slacker.JobDefinition{
		Description: "A cron job that runs every minute",
		Handler: func(jobCtx slacker.JobContext) error {
			jobCtx.APIClient().PostMessage("#test", slack.MsgOptionText("Hello!", false))
			return nil
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
