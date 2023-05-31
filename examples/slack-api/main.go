package main

import (
	"context"
	"log"
	"os"

	"github.com/shomali11/slacker/v2"
	"github.com/slack-go/slack"
)

// Showcasing the ability to access the github.com/slack-go/slack API and upload a file

func main() {
	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))

	definition := &slacker.CommandDefinition{
		Description: "Upload a sentence!",
		Handler: func(ctx slacker.CommandContext) {
			sentence := ctx.Request().Param("sentence")
			apiClient := ctx.APIClient()
			event := ctx.Event()

			apiClient.PostMessage(event.ChannelID, slack.MsgOptionText("Uploading file ...", false))
			_, err := apiClient.UploadFile(slack.FileUploadParameters{Content: sentence, Channels: []string{event.ChannelID}})
			ctx.Response().ReplyError(err)
		},
	}

	bot.AddCommand("upload <sentence>", definition)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
