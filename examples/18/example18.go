package main

import (
	"context"
	"log"
	"os"

	"github.com/shomali11/slacker"
	"github.com/slack-go/slack"
)

func main() {
	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))
	bot.Command("ping", &slacker.CommandDefinition{
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			response.Reply("pong")
		},
	})

	// Run every minute
	bot.Job("0 * * * * *", &slacker.JobDefinition{
		Description: "A cron job that runs every minute",
		Handler: func(jobCtx slacker.JobContext) {
			jobCtx.APIClient().PostMessage("#test", slack.MsgOptionText("Hello!", false))
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
