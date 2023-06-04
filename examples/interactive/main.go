package main

import (
	"context"
	"log"
	//"os"

	"github.com/shomali11/slacker/v2"
	"github.com/slack-go/slack"
)

// Implements a basic interactive command.
// This assumes that a slash command `/mood` is defined for your app.

func main() {
	bot := slacker.NewClient("xoxb-13360094916-2243791173942-Xl56AaFTAHLnNJXTV5VQ2O1A", "xapp-1-A027JMM1RV2-5371200676276-829e2afbe83227c95c67cae993c9645fd60010e2541812b87d207c781822b124")

	//bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))
	bot.AddCommand("mood", &slacker.CommandDefinition{
		Handler:  slackerCmd("mood"),
		HideHelp: true,
	})

	bot.AddInteraction("mood", &slacker.InteractionDefinition{
		Handler:  slackerInteractive,
		HideHelp: true,
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

func slackerCmd(blockID string) slacker.CommandHandler {
	return func(ctx slacker.CommandContext) {
		happyBtn := slack.NewButtonBlockElement("happy", "true", slack.NewTextBlockObject("plain_text", "Happy 🙂", true, false))
		happyBtn.Style = "primary"
		sadBtn := slack.NewButtonBlockElement("sad", "false", slack.NewTextBlockObject("plain_text", "Sad ☹️", true, false))
		sadBtn.Style = "danger"

		ctx.Response().ReplyBlocks([]slack.Block{
			slack.NewSectionBlock(slack.NewTextBlockObject(slack.PlainTextType, "What is your mood today?", true, false), nil, nil),
			slack.NewActionBlock(blockID, happyBtn, sadBtn),
		})
	}
}

func slackerInteractive(ctx slacker.InteractionContext) {
	text := ""
	action := ctx.Event().ActionCallback.BlockActions[0]
	switch action.ActionID {
	case "happy":
		text = "I'm happy to hear you are happy!"
	case "sad":
		text = "I'm sorry to hear you are sad."
	default:
		text = "I don't understand your mood..."
	}

	ctx.Response().Reply(text, slacker.WithReplace(ctx.Event().Message.Timestamp))
}
