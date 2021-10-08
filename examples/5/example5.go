package main

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/shomali11/slacker"
)

func main() {
	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))

	messageReplyDefinition := &slacker.CommandDefinition{
		Description: "Tests errors in new messages",
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			response.ReportError(errors.New("oops, an error occurred"))
		},
	}

	threadReplyDefinition := &slacker.CommandDefinition{
		Description: "Tests errors in threads",
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			response.ReportError(errors.New("oops, an error occurred"), slacker.WithThreadError(true))
		},
	}

	bot.Command("message", messageReplyDefinition)
	bot.Command("thread", threadReplyDefinition)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
