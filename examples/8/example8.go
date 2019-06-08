package main

import (
	"context"
	"errors"
	"github.com/shomali11/slacker"
	"log"
	"time"
)

func main() {
	bot := slacker.NewClient("<YOUR SLACK BOT TOKEN>")

	definition := &slacker.CommandDefinition{
		Description: "Process!",
		Handler: func(request slacker.Request, response slacker.ResponseWriter) {
			timedContext, cancel := context.WithTimeout(request.Context(), time.Second)
			defer cancel()

			select {
			case <-timedContext.Done():
				response.ReportError(errors.New("timed out"))
			case <-time.After(time.Minute):
				response.Reply("Processing done!")
			}
		},
	}

	bot.Command("process", definition)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
