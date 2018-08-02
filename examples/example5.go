package main

import (
	"errors"
	"github.com/shomali11/slacker"
	"log"
)

func main() {
	bot := slacker.NewClient("<YOUR SLACK BOT TOKEN>")

	bot.Command("test", "Tests something", func(request slacker.Request, response slacker.ResponseWriter) {
		response.ReportError(errors.New("Oops!"))
	})

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	if err := bot.Listen(ctx); err != nil {
		log.Fatal(err)
	}
}
