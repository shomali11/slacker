package main

import (
	"context"
	"log"
	"os"

	"github.com/shomali11/slacker/v2"
	"github.com/slack-go/slack"
)

// Implements a basic interactive command.
// This assumes that a slash command `/mood` is defined for your app.

func slackerCmd(blockID string) slacker.CommandHandler {
	return func(ctx slacker.CommandContext) {
		happyBtn := slack.NewButtonBlockElement("happy", "true", slack.NewTextBlockObject("plain_text", "Happy 🙂", true, false))
		happyBtn.Style = "primary"
		sadBtn := slack.NewButtonBlockElement("sad", "false", slack.NewTextBlockObject("plain_text", "Sad ☹️", true, false))
		sadBtn.Style = "danger"

		ctx.Response().Reply("", slacker.WithBlocks([]slack.Block{
			slack.NewSectionBlock(slack.NewTextBlockObject(slack.PlainTextType, "What is your mood today?", true, false), nil, nil),
			slack.NewActionBlock(blockID, happyBtn, sadBtn),
		}))
	}
}

func slackerInteractive(ctx slacker.InteractiveContext) {
	text := ""
	action := ctx.Callback().ActionCallback.BlockActions[0]
	switch action.ActionID {
	case "happy":
		text = "I'm happy to hear you are happy!"
	case "sad":
		text = "I'm sorry to hear you are sad."
	default:
		text = "I don't understand your mood..."
	}

	_, _, _ = ctx.APIClient().PostMessage(ctx.Callback().Channel.ID, slack.MsgOptionText(text, false),
		slack.MsgOptionReplaceOriginal(ctx.Callback().ResponseURL))
}

func main() {
	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))
	bot.AddCommand("mood", &slacker.CommandDefinition{
		BlockID:             "mood",
		Handler:             slackerCmd("mood"),
		InteractiveCallback: slackerInteractive,
		HideHelp:            true,
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
