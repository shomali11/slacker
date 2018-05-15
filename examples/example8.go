package main

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/shomali11/slacker"
)

func main() {
	bot := slacker.NewClient("<YOUR SLACK BOT TOKEN>")

	bot.Command("process", "Process!", func(request slacker.Request, response slacker.ResponseWriter) {
		timedContext, cancel := context.WithTimeout(request.Context, time.Second)
		defer cancel()

		select {
		case <-timedContext.Done():
			response.ReportError(errors.New("Timed out"))
		case <-time.After(time.Minute):
			response.Reply("Processing done!")
		}
	})

	err := bot.Listen()
	if err != nil {
		log.Fatal(err)
	}
}
