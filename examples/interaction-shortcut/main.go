package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/shomali11/slacker/v2"
	"github.com/slack-go/slack"
)

// Implements a basic interactive command with modal view.

func main() {
	bot := slacker.NewClient(
		os.Getenv("SLACK_BOT_TOKEN"),
		os.Getenv("SLACK_APP_TOKEN"),
		slacker.WithDebug(false),
	)

	bot.AddInteraction(&slacker.InteractionDefinition{
		InteractionID: "mood-survey-message-shortcut-callback-id",
		Handler:       moodShortcutHandler,
		Type:          slack.InteractionTypeMessageAction,
	})

	bot.AddInteraction(&slacker.InteractionDefinition{
		InteractionID: "mood-survey-global-shortcut-callback-id",
		Handler:       moodShortcutHandler,
		Type:          slack.InteractionTypeShortcut,
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

func moodShortcutHandler(ctx *slacker.InteractionContext) {
	switch ctx.Callback().Type {
	case slack.InteractionTypeMessageAction:
		{
			fmt.Print("Message shortcut.\n")
		}
	case slack.InteractionTypeShortcut:
		{
			fmt.Print("Global shortcut.\n")
		}
	}
}
