package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/broxgit/slacker"
	"github.com/slack-go/slack"
)

func main() {
	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))

	definition := &slacker.CommandDefinition{
		Description: "Upload a sentence!",
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			sentence := request.Param("sentence")
			client := botCtx.Client()
			ev := botCtx.Event()

			if ev.Channel != "" {
				client.PostMessage(ev.Channel, slack.MsgOptionText("Uploading file ...", false))
				_, err := client.UploadFile(slack.FileUploadParameters{Content: sentence, Channels: []string{ev.Channel}})
				if err != nil {
					fmt.Printf("Error encountered when uploading file: %+v\n", err)
				}
			}
		},
	}

	bot.Command("upload <sentence>", definition)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
