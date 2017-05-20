package main

import (
	"github.com/nlopes/slack"
	"github.com/shomali11/slacker"
	"log"
)

func main() {
	bot := slacker.NewClient("<YOUR SLACK BOT TOKEN>")

	bot.Command("upload <word>", "Upload a word!", func(request *slacker.Request, response *slacker.Response) {
		word := request.Param("word")
		channel := request.Event.Channel
		bot.Client.UploadFile(slack.FileUploadParameters{Content: word, Channels: []string{channel}})
	})

	err := bot.Listen()
	if err != nil {
		log.Fatal(err)
	}
}
