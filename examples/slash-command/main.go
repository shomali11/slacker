package main

import (
	"context"
	"log"
	"os"

	"github.com/shomali11/slacker/v2"
)

// Implements a simple slash command.
// In this example, we hide the command from `help`'s results.
// This assumes you have the slash command `/ping` defined for your app.

func main() {
	bot := slacker.NewClient(
		os.Getenv("SLACK_BOT_TOKEN"),
		os.Getenv("SLACK_APP_TOKEN"),
	)

	bot.AddCommand(&slacker.CommandDefinition{
		Command: "ping",
		Handler: func(ctx slacker.CommandContext) {
			ctx.Response().Reply("pong")
		},
		HideHelp: true,
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
