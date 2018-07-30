package main

import (
	"log"

	"github.com/nlopes/slack"
	"github.com/shomali11/slacker"
	"time"
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

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	if err := bot.Listen(ctx); err != nil {
		log.Fatal(err)
	}
}
