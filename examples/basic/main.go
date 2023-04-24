package main

import (
	"context"
	"log"
	"os"

	"github.com/shomali11/slacker"
)

// Defining a command using slacker

func main() {
	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))

	definition := &slacker.CommandDefinition{
		Handler: func(botCtx slacker.CommandContext) {
			botCtx.Response().Reply("pong")
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
