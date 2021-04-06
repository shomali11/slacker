package main

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"github.com/shomali11/slacker"
)

func main() {
	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))

	definition := &slacker.CommandDefinition{
		Description: "Process!",
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			timedContext, cancel := context.WithTimeout(botCtx.Context(), time.Second)
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
