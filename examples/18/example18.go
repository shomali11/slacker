package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/shomali11/slacker"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

func slackerCmd(actionID string) func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
	return func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
		happyBtn := slack.NewButtonBlockElement("happy", "true", slack.NewTextBlockObject("plain_text", "Happy 🙂", true, false))
		happyBtn.Style = "primary"
		sadBtn := slack.NewButtonBlockElement("sad", "false", slack.NewTextBlockObject("plain_text", "Sad ☹️", true, false))
		sadBtn.Style = "danger"

		err := response.Reply("", slacker.WithBlocks([]slack.Block{
			slack.NewSectionBlock(slack.NewTextBlockObject(slack.PlainTextType, "What is your mood today?", true, false), nil, nil),
			slack.NewActionBlock(actionID, happyBtn, sadBtn),
		}))

		if err != nil {
			fmt.Println(err)
		}
	}
}

func slackerInteractive(s *slacker.Slacker, e *socketmode.Event, callback *slack.InteractionCallback, request *socketmode.Request) {
	text := ""
	action := callback.ActionCallback.BlockActions[0]
	switch action.ActionID {
	case "happy":
		text = "I'm happy to hear you are happy!"
	case "sad":
		text = "I'm sorry to hear you are sad."
	default:
		text = "I don't understand your mood..."
	}

	_, _, _ = s.Client().PostMessage(callback.Channel.ID, slack.MsgOptionText(text, false),
		slack.MsgOptionReplaceOriginal(callback.ResponseURL))
}

func main() {
	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))
	bot.Command("slacker-cmd", &slacker.CommandDefinition{
		BlockID:     "slacker_cmd",
		Handler:     slackerCmd("slacker_cmd"),
		Interactive: slackerInteractive,
		HideHelp:    true,
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
