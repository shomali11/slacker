package main

import (
	"github.com/shomali11/slacker"
	"log"
)

func main() {
	bot := slacker.NewClient("<YOUR SLACK BOT TOKEN>")

	bot.Command("echo <word>", "Echo a word!", func(request slacker.Request, response slacker.ResponseWriter) {
		word := request.Param("word")
		response.Reply(word)
	})

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	if err := bot.Listen(ctx); err != nil {
		log.Fatal(err)
	}
}
