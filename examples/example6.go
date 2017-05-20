package main

import (
	"github.com/shomali11/slacker"
	"log"
)

func main() {
	bot := slacker.NewClient("<YOUR SLACK BOT TOKEN>")

	bot.Command("ping", "Ping!", func(request *slacker.Request, response *slacker.Response) {
		response.Reply("Pong")
	})

	err := bot.Listen()
	if err != nil {
		log.Fatal(err)
	}
}
