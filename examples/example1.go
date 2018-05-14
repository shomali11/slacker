package main

import (
	"log"

	"github.com/shomali11/slacker"
)

func main() {
	bot := slacker.NewClient("<YOUR SLACK BOT TOKEN>")

	bot.Command("ping", "Ping!", func(request slacker.Request, response slacker.ResponseWriter) {
		response.Reply("pong")
	})

	err := bot.Listen()
	if err != nil {
		log.Fatal(err)
	}
}
