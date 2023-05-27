package main

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/shomali11/slacker/v2"
)

// Showcasing the ability to leverage `context.Context` to add a timeout

func main() {
	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))

	definition := &slacker.CommandDefinition{
		Description: "Process!",
		Handler: func(ctx slacker.CommandContext) {
			timedContext, cancel := context.WithTimeout(ctx.Context(), 5*time.Second)
			defer cancel()

			duration := time.Duration(rand.Int()%10+1) * time.Second

			select {
			case <-timedContext.Done():
				ctx.Response().Error(errors.New("timed out"))
			case <-time.After(duration):
				ctx.Response().Reply("Processing done!")
			}
		},
	}

	bot.AddCommand("process", definition)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
