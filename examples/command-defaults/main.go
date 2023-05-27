package main

import (
	"context"
	"log"
	"os"

	"github.com/shomali11/slacker/v2"
)

// Defining a command with two parameters. Parsing one as a string and the other as an integer.
// (The second parameter is the default value in case no parameter was passed or could not parse the value)

func main() {
	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))

	definition := &slacker.CommandDefinition{
		Description: "Repeat a word a number of times!",
		Examples:    []string{"repeat hello 10"},
		Handler: func(ctx slacker.CommandContext) {
			word := ctx.Request().StringParam("word", "Hello!")
			number := ctx.Request().IntegerParam("number", 1)
			for i := 0; i < number; i++ {
				ctx.Response().Reply(word)
			}
		},
	}

	bot.AddCommand("repeat {word} {number}", definition)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
