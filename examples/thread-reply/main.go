package main

import (
	"context"
	"log"
	"os"

	"github.com/shomali11/slacker"
)

// Defining a command with an optional description and example. The handler replies to a thread.

func main() {
	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))

	definition := &slacker.CommandDefinition{
		Description: "Ping!",
		Examples:    []string{"ping"},
		Handler: func(botCtx slacker.CommandContext) {
			botCtx.Response().Reply("pong", slacker.WithThreadReply(true))
		},
	}

	bot.AddCommand("ping", definition)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
