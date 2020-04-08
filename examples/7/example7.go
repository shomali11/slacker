package main

import (
	"context"
	"log"

	"github.com/shomali11/slacker"
	"github.com/slack-go/slack"
)

func main() {
	bot := slacker.NewClient("<YOUR SLACK BOT TOKEN>")

	definition := &slacker.CommandDefinition{
		Description: "Upload a word!",
		Handler: func(request slacker.Request, response slacker.ResponseWriter) {
			word := request.Param("word")
			channel := request.Event().Channel

			rtm := response.RTM()
			client := response.Client()

			rtm.SendMessage(rtm.NewOutgoingMessage("Uploading file ...", channel))
			client.UploadFile(slack.FileUploadParameters{Content: word, Channels: []string{channel}})
		},
	}

	bot.Command("upload <word>", definition)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
