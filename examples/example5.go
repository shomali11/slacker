package main

import (
	"context"
	"errors"
	"github.com/shomali11/slacker"
	"log"
)

func main() {
	bot := slacker.NewClient("<YOUR SLACK BOT TOKEN>")

	bot.Command("test", "Tests something", func(request slacker.Request, response slacker.ResponseWriter) {
		response.ReportError(errors.New("Oops!"))
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
