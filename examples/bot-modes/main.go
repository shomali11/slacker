package main

import (
	"context"
	"log"
	"os"

	"github.com/shomali11/slacker/v2"
)

// Configure bot to process other bot events

func main() {
	bot := slacker.NewClient(
		os.Getenv("SLACK_BOT_TOKEN"),
		os.Getenv("SLACK_APP_TOKEN"),
		slacker.WithBotMode(slacker.BotModeIgnoreApp),
	)

	bot.AddCommand(&slacker.CommandDefinition{
		Command: "hello",
		Handler: func(ctx *slacker.CommandContext) {
			ctx.Response().Reply("hai!")
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
