package main

import (
	"context"
	"log"
	"os"

	"github.com/shomali11/slacker"
	"github.com/slack-go/slack"
)

// Showcasing the ability to add blocks to a `Reply`

func main() {
	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))

	definition := &slacker.CommandDefinition{
		Description: "Echo a word!",
		Handler: func(botCtx slacker.CommandContext) {
			word := botCtx.Request().Param("word")

			attachments := []slack.Block{}
			attachments = append(attachments, slack.NewContextBlock("1",
				slack.NewTextBlockObject("mrkdwn", word, false, false)),
			)

			// When using blocks the message argument will be thrown away and can be left blank.
			botCtx.Response().Reply("", slacker.WithBlocks(attachments))
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
