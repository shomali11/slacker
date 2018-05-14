package main

import (
	"log"

	"github.com/shomali11/slacker"
)

func main() {
	bot := slacker.NewClient("<YOUR SLACK BOT TOKEN>")

	bot.Command("echo <word>", "Echo a word!", func(request slacker.Request, response slacker.ResponseWriter) {
		word := request.Param("word")
		response.Reply(word)
	})

	err := bot.Listen()
	if err != nil {
		log.Fatal(err)
	}
}
