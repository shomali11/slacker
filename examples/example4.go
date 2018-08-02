package main

import (
	"github.com/shomali11/slacker"
	"log"
)

func main() {
	bot := slacker.NewClient("<YOUR SLACK BOT TOKEN>")

	bot.Command("repeat <word> <number>", "Repeat a word a number of times!", func(request slacker.Request, response slacker.ResponseWriter) {
		word := request.StringParam("word", "Hello!")
		number := request.IntegerParam("number", 1)
		for i := 0; i < number; i++ {
			response.Reply(word)
		}
	})

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	if err := bot.Listen(ctx); err != nil {
		log.Fatal(err)
	}

}
