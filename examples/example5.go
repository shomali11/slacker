package main

import (
	"errors"
	"github.com/shomali11/slacker"
	"log"
)

func main() {
	bot := slacker.NewClient("<YOUR SLACK BOT TOKEN>")

	bot.Command("test", "Tests something", func(request *slacker.Request, response *slacker.Response) {
		response.ReportError(errors.New("Oops!"))
	})

	err := bot.Listen()
	if err != nil {
		log.Fatal(err)
	}
}
