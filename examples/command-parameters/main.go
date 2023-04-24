package main

import (
	"context"
	"log"
	"os"

	"github.com/shomali11/slacker"
)

// Defining a command with a parameter. Parameters surrounded with {} will be satisfied with a word.
// Parameters surrounded with <> are "greedy" and will take as much input as fed.

func main() {
	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))

	bot.AddCommand("echo {word}", &slacker.CommandDefinition{
		Description: "Echo a word!",
		Examples:    []string{"echo hello"},
		Handler: func(botCtx slacker.CommandContext) {
			word := botCtx.Request().Param("word")
			botCtx.Response().Reply(word)
		},
	})

	bot.AddCommand("say <sentence>", &slacker.CommandDefinition{
		Description: "Say a sentence!",
		Examples:    []string{"say hello there everyone!"},
		Handler: func(botCtx slacker.CommandContext) {
			sentence := botCtx.Request().Param("sentence")
			botCtx.Response().Reply(sentence)
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
