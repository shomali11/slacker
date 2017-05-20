package main

import (
	"github.com/shomali11/slacker"
	"log"
)

func main() {
	bot := slacker.NewClient("<YOUR SLACK BOT TOKEN>")

	bot.Command("echo <word>", "Echo a word!", func(request *slacker.Request, response *slacker.Response) {
		word := request.Param("word")
		response.Reply(word)
	})

	err := bot.Listen()
	if err != nil {
		log.Fatal(err)
	}
}
