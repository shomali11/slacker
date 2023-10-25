package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/shomali11/slacker/v2"
	"github.com/slack-go/slack"
)

var moodSurveyView = slack.ModalViewRequest{
	Type:       "modal",
	CallbackID: "mood-survey-callback-id",
	Title: &slack.TextBlockObject{
		Type: "plain_text",
		Text: "Which mood are you in?",
	},
	Submit: &slack.TextBlockObject{
		Type: "plain_text",
		Text: "Submit",
	},
	NotifyOnClose: true,
	Blocks: slack.Blocks{
		BlockSet: []slack.Block{
			&slack.InputBlock{
				Type:    slack.MBTInput,
				BlockID: "mood",
				Label: &slack.TextBlockObject{
					Type: "plain_text",
					Text: "Mood",
				},
				Element: &slack.SelectBlockElement{
					Type:     slack.OptTypeStatic,
					ActionID: "mood",
					Options: []*slack.OptionBlockObject{
						{
							Text: &slack.TextBlockObject{
								Type: "plain_text",
								Text: "Happy",
							},
							Value: "Happy",
						},
						{
							Text: &slack.TextBlockObject{
								Type: "plain_text",
								Text: "Sad",
							},
							Value: "Sad",
						},
					},
				},
			},
		},
	},
}

// Implements a basic interactive command with modal view.
func main() {
	bot := slacker.NewClient(
		os.Getenv("SLACK_BOT_TOKEN"),
		os.Getenv("SLACK_APP_TOKEN"),
		slacker.WithDebug(false),
	)

	bot.AddCommand(&slacker.CommandDefinition{
		Command: "mood",
		Handler: moodCmdHandler,
	})

	bot.AddInteraction(&slacker.InteractionDefinition{
		InteractionID: "mood-survey-callback-id",
		Handler:       moodViewHandler,
		Type:          slack.InteractionTypeViewSubmission,
	})

	bot.AddInteraction(&slacker.InteractionDefinition{
		InteractionID: "mood-survey-callback-id",
		Handler:       moodViewHandler,
		Type:          slack.InteractionTypeViewClosed,
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

func moodCmdHandler(ctx *slacker.CommandContext) {
	_, err := ctx.SlackClient().OpenView(
		ctx.Event().Data.(*slack.SlashCommand).TriggerID,
		moodSurveyView,
	)
	if err != nil {
		log.Printf("ERROR openEscalationModal: %v", err)
	}
}

func moodViewHandler(ctx *slacker.InteractionContext) {
	switch ctx.Callback().Type {
	case slack.InteractionTypeViewSubmission:
		{
			viewState := ctx.Callback().View.State.Values
			fmt.Printf(
				"Mood view submitted.\nMood: %s\n",
				viewState["mood"]["mood"].SelectedOption.Value,
			)
		}
	case slack.InteractionTypeViewClosed:
		{
			fmt.Print("Mood view closed.\n")
		}
	}
}
