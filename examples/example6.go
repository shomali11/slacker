package main

import (
	"log"
	"time"

	"github.com/shomali11/slacker"
)

func main() {
	bot := slacker.NewClient("<YOUR SLACK BOT TOKEN>")

	bot.Command("time", "Server time!", func(request slacker.Request, response slacker.ResponseWriter) {
		response.Typing()

		time.Sleep(time.Second)

		response.Reply(time.Now().Format(time.RFC1123))
	})

	err := bot.Listen()
	if err != nil {
		log.Fatal(err)
	}
}
