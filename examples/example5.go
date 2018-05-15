package main

import (
	"errors"
	"log"

	"github.com/shomali11/slacker"
)

func main() {
	bot := slacker.NewClient("<YOUR SLACK BOT TOKEN>")

	bot.Command("test", "Tests something", func(request slacker.Request, response slacker.ResponseWriter) {
		response.ReportError(errors.New("Oops!"))
	})

	err := bot.Listen()
	if err != nil {
		log.Fatal(err)
	}
}
