package main

import (
	"context"
	"log"
	"os"

	"github.com/shomali11/slacker/v2"
	"github.com/slack-go/slack"
)

// Showcasing the ability to add blocks to a `Reply`

func main() {
	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))

	definition := &slacker.CommandDefinition{
		Description: "Echo a word!",
		Handler: func(ctx slacker.CommandContext) {
			word := ctx.Request().Param("word")

			blocks := []slack.Block{}
			blocks = append(blocks, slack.NewContextBlock("1",
				slack.NewTextBlockObject("mrkdwn", word, false, false)),
			)

			// When using blocks the message argument will be thrown away and can be left blank.
			ctx.Response().Reply("", slacker.WithBlocks(blocks))
		},
	}

	bot.AddCommand("echo {word}", definition)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
