package main

import (
	"context"
	"log"
	"os"

	"github.com/shomali11/slacker"
)

// Configure bot to process other bot events

func main() {
	bot := slacker.NewClient(
		os.Getenv("SLACK_BOT_TOKEN"),
		os.Getenv("SLACK_APP_TOKEN"),
		slacker.WithBotInteractionMode(slacker.BotInteractionModeIgnoreApp),
	)

	bot.AddCommand("hello", &slacker.CommandDefinition{
		Handler: func(botCtx slacker.CommandContext) {
			botCtx.Response().Reply("hai!")
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
