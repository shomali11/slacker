package main

import (
	"context"
	"log"
	"os"

	"github.com/shomali11/slacker/v2"
	"github.com/slack-go/slack"
)

func main() {
	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))

	bot.UnhandledInteractiveCallback(func(ctx slacker.InteractiveContext) {
		callback := ctx.Callback()
		if callback.Type != slack.InteractionTypeBlockActions {
			return
		}

		if len(callback.ActionCallback.BlockActions) != 1 {
			return
		}

		action := callback.ActionCallback.BlockActions[0]
		if action.BlockID != "mood-block" {
			return
		}

		var text string
		switch action.ActionID {
		case "happy":
			text = "I'm happy to hear you are happy!"
		case "sad":
			text = "I'm sorry to hear you are sad."
		default:
			text = "I don't understand your mood..."
		}

		_, _, _ = ctx.APIClient().PostMessage(callback.Channel.ID, slack.MsgOptionText(text, false),
			slack.MsgOptionReplaceOriginal(callback.ResponseURL))

		ctx.SocketModeClient().Ack(*ctx.Event().Request)
	})

	definition := &slacker.CommandDefinition{
		Handler: func(ctx slacker.CommandContext) {
			happyBtn := slack.NewButtonBlockElement("happy", "true", slack.NewTextBlockObject("plain_text", "Happy 🙂", true, false))
			happyBtn.Style = "primary"
			sadBtn := slack.NewButtonBlockElement("sad", "false", slack.NewTextBlockObject("plain_text", "Sad ☹️", true, false))
			sadBtn.Style = "danger"

			ctx.Response().Reply("", slacker.WithBlocks([]slack.Block{
				slack.NewSectionBlock(slack.NewTextBlockObject(slack.PlainTextType, "What is your mood today?", true, false), nil, nil),
				slack.NewActionBlock("mood-block", happyBtn, sadBtn),
			}))
		},
	}

	bot.AddCommand("mood", definition)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
