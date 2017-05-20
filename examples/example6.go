package main

import (
	"github.com/shomali11/slacker"
	"log"
	"time"
)

func main() {
	bot := slacker.NewClient("<YOUR SLACK BOT TOKEN>")

	bot.Command("time", "Server time!", func(request *slacker.Request, response *slacker.Response) {
		response.Typing()

		time.Sleep(time.Second)

		response.Reply(time.Now().Format(time.RFC1123))
	})

	err := bot.Listen()
	if err != nil {
		log.Fatal(err)
	}
}
