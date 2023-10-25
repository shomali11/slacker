package main

import (
	"context"
	"log"
	"os"

	"github.com/shomali11/slacker/v2"
	"github.com/slack-go/slack"
)

// Implements a basic interactive command.

func main() {
	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))
	bot.AddCommand(&slacker.CommandDefinition{
		Command: "mood",
		Handler: slackerCmd("mood"),
	})

	bot.AddInteraction(&slacker.InteractionDefinition{
		InteractionID: "mood",
		Handler:       slackerInteractive,
		Type:          slack.InteractionTypeBlockActions,
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

func slackerCmd(blockID string) slacker.CommandHandler {
	return func(ctx *slacker.CommandContext) {
		happyBtn := slack.NewButtonBlockElement("happy", "true", slack.NewTextBlockObject("plain_text", "Happy 🙂", true, false))
		happyBtn.Style = slack.StylePrimary
		sadBtn := slack.NewButtonBlockElement("sad", "false", slack.NewTextBlockObject("plain_text", "Sad ☹️", true, false))
		sadBtn.Style = slack.StyleDanger

		ctx.Response().ReplyBlocks([]slack.Block{
			slack.NewSectionBlock(slack.NewTextBlockObject(slack.PlainTextType, "What is your mood today?", true, false), nil, nil),
			slack.NewActionBlock(blockID, happyBtn, sadBtn),
		})
	}
}

func slackerInteractive(ctx *slacker.InteractionContext) {
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

	ctx.Response().Reply(text, slacker.WithReplace(ctx.Callback().Message.Timestamp))
}
