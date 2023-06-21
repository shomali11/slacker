package main

import (
	"context"
	"log"
	"os"

	"github.com/shomali11/slacker/v2"
)

// Defining a command with a parameter. Parameters surrounded with {} will be satisfied with a word.
// Parameters surrounded with <> are "greedy" and will take as much input as fed.

func main() {
	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))

	bot.AddCommand(&slacker.CommandDefinition{
		Command:     "echo {word}",
		Description: "Echo a word!",
		Examples:    []string{"echo hello"},
		Handler: func(ctx *slacker.CommandContext) {
			word := ctx.Request().Param("word")
			ctx.Response().Reply(word)
		},
	})

	bot.AddCommand(&slacker.CommandDefinition{
		Command:     "say <sentence>",
		Description: "Say a sentence!",
		Examples:    []string{"say hello there everyone!"},
		Handler: func(ctx *slacker.CommandContext) {
			sentence := ctx.Request().Param("sentence")
			ctx.Response().Reply(sentence)
		},
	})

	// If no values were provided, the parameters will return empty strings.
	// You can define a default value in case no parameter was passed (or the value could not be parsed)
	bot.AddCommand(&slacker.CommandDefinition{
		Command:     "repeat {word} {number}",
		Description: "Repeat a word a number of times!",
		Examples:    []string{"repeat hello 10"},
		Handler: func(ctx *slacker.CommandContext) {
			word := ctx.Request().StringParam("word", "Hello!")
			number := ctx.Request().IntegerParam("number", 1)
			for i := 0; i < number; i++ {
				ctx.Response().Reply(word)
			}
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
