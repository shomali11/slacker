package main

import (
	"log"

	"context"
	"github.com/nlopes/slack"
	"github.com/shomali11/slacker"
)

func main() {
	bot := slacker.NewClient("<YOUR SLACK BOT TOKEN>")

	bot.Command("upload <word>", "Upload a word!", func(request slacker.Request, response slacker.ResponseWriter) {
		word := request.Param("word")
		channel := request.Event().Channel

		rtm := response.RTM()
		client := response.Client()

		rtm.SendMessage(rtm.NewOutgoingMessage("Uploading file ...", channel))
		client.UploadFile(slack.FileUploadParameters{Content: word, Channels: []string{channel}})
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
