package main

import (
	"context"
	"log"
	"os"

	"github.com/shomali11/slacker/v2"
)

// Defining a command with aliases

func main() {
	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))
	bot.AddCommand(&slacker.CommandDefinition{
		Command: "echo {word}",
		Aliases: []string{
			"repeat {word}",
			"mimic {word}",
		},
		Description: "Echo a word!",
		Examples: []string{
			"echo hello",
			"repeat hello",
			"mimic hello",
		},
		Handler: func(ctx *slacker.CommandContext) {
			word := ctx.Request().Param("word")
			ctx.Response().Reply(word)
		},
	})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
