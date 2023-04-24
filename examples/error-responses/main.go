package main

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/shomali11/slacker"
)

// Defines two commands that display sending errors to the Slack channel.
// One that replies as a new message. The other replies to the thread.

func main() {
	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))

	messageReplyDefinition := &slacker.CommandDefinition{
		Description: "Tests errors in new messages",
		Handler: func(botCtx slacker.CommandContext) {
			botCtx.Response().Error(errors.New("oops, an error occurred"))
		},
	}

	threadReplyDefinition := &slacker.CommandDefinition{
		Description: "Tests errors in threads",
		Handler: func(botCtx slacker.CommandContext) {
			botCtx.Response().Error(errors.New("oops, an error occurred"), slacker.WithThreadError(true))
		},
	}

	bot.AddCommand("message", messageReplyDefinition)
	bot.AddCommand("thread", threadReplyDefinition)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
